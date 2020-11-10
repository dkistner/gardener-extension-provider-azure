// Copyright (c) 2020 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controlplane

import (
	"context"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/controlplane"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	azurev1alpha1 "github.com/gardener/remedy-controller/pkg/apis/azure/v1alpha1"
)

// NewActuator creates a new Actuator that acts upon and updates the status of ControlPlane resources.
func NewActuator(
	a controlplane.Actuator,
	logger logr.Logger,
) controlplane.Actuator {
	return &actuator{
		Actuator: a,
		logger:   logger.WithName("azure-controlplane-actuator"),
	}
}

// actuator is an Actuator that acts upon and updates the status of ControlPlane resources.
type actuator struct {
	controlplane.Actuator
	client client.Client
	logger logr.Logger
}

// InjectFunc enables injecting Kubernetes dependencies into actuator's dependencies.
func (a *actuator) InjectFunc(f inject.Func) error {
	return f(a.Actuator)
}

// InjectClient injects the given client into the valuesProvider.
func (a *actuator) InjectClient(client client.Client) error {
	a.client = client
	return nil
}

// Delete reconciles the given controlplane and cluster, deleting the additional
// control plane components as needed.
// Before delegating to the composed Actuator, it ensures that all remedy controller resources have been deleted.
func (a *actuator) Delete(
	ctx context.Context,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
) error {
	// Forcefully delete all remaining remedy controller resources
	a.logger.Info("Deleting all remaining remedy controller resources")
	pubipList := &azurev1alpha1.PublicIPAddressList{}
	if err := a.client.List(ctx, pubipList, client.InNamespace(cp.Namespace)); err != nil {
		return errors.Wrap(err, "could not list publicipaddresses")
	}
	for _, pubip := range pubipList.Items {
		// Add the do-not-clean annotation to the publicipaddress resource
		// This annotation prevents attempts to clean the Azure IP address if it still exists, resulting in much faster deletion
		if err := retry.RetryOnConflict(retry.DefaultBackoff, func() error {
			pubip.Annotations = add(pubip.Annotations, "azure.remedy.gardener.cloud/do-not-clean", strconv.FormatBool(true))
			return a.client.Update(ctx, &pubip)
		}); err != nil {
			return errors.Wrap(err, "could not add do-not-clean annotation on publicipaddress")
		}
	}
	if err := a.client.DeleteAllOf(ctx, &azurev1alpha1.PublicIPAddress{}, client.InNamespace(cp.Namespace)); err != nil {
		return errors.Wrap(err, "could not delete publicipaddress resources")
	}
	if err := a.client.DeleteAllOf(ctx, &azurev1alpha1.VirtualMachine{}, client.InNamespace(cp.Namespace)); err != nil {
		return errors.Wrap(err, "could not delete virtualmachine resources")
	}

	// Wait until the remaining remedy controller resources have been deleted
	a.logger.Info("Waiting for the remaining remedy controller resources to be deleted")
	timeoutCtx1, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if err := kutil.WaitUntilResourcesDeleted(timeoutCtx1, a.client, &azurev1alpha1.PublicIPAddressList{}, 5*time.Second, client.InNamespace(cp.Namespace)); err != nil {
		return errors.Wrap(err, "could not wait for publicipaddress resources to be deleted")
	}
	timeoutCtx2, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	if err := kutil.WaitUntilResourcesDeleted(timeoutCtx2, a.client, &azurev1alpha1.VirtualMachineList{}, 5*time.Second, client.InNamespace(cp.Namespace)); err != nil {
		return errors.Wrap(err, "could not wait for virtualmachine resources to be deleted")
	}

	// Call Delete on the composed Actuator
	return a.Actuator.Delete(ctx, cp, cluster)
}

// Migrate reconciles the given controlplane and cluster, deleting the additional
// control plane components as needed.
func (a *actuator) Migrate(
	ctx context.Context,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
) error {
	return a.Delete(ctx, cp, cluster)
}

func add(m map[string]string, key, value string) map[string]string {
	if m == nil {
		m = make(map[string]string)
	}
	m[key] = value
	return m
}
