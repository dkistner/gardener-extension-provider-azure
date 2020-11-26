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

package validation

import (
	"fmt"

	apisazure "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"
	"github.com/gardener/gardener-extension-provider-azure/pkg/azure"
	cidrvalidation "github.com/gardener/gardener/pkg/utils/validation/cidr"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const (
	natGatewayMinTimeoutInMinutes int32 = 4
	natGatewayMaxTimeoutInMinutes int32 = 120
)

// ValidateInfrastructureConfig validates a InfrastructureConfig object.
func ValidateInfrastructureConfig(infra *apisazure.InfrastructureConfig, nodesCIDR, podsCIDR, servicesCIDR *string, annotations map[string]string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	var (
		nodes    cidrvalidation.CIDR
		pods     cidrvalidation.CIDR
		services cidrvalidation.CIDR
	)

	if nodesCIDR != nil {
		nodes = cidrvalidation.NewCIDR(*nodesCIDR, nil)
	}
	if podsCIDR != nil {
		pods = cidrvalidation.NewCIDR(*podsCIDR, nil)
	}
	if servicesCIDR != nil {
		services = cidrvalidation.NewCIDR(*servicesCIDR, nil)
	}

	// Currently, we will not allow deployments into existing resource groups or VNets although this functionality
	// is already implemented, because the Azure cloud provider is not cleaning up self-created resources properly.
	// This resources would be orphaned when the cluster will be deleted. We block these cases thereby that the Azure shoot
	// validation here will fail for those cases.
	// TODO: remove the following block and uncomment below blocks once deployment into existing resource groups works properly.
	if infra.ResourceGroup != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("resourceGroup"), infra.ResourceGroup, "specifying an existing resource group is not supported yet"))
	}

	annotationValue, annotationExists := annotations[azure.ShootVmoUsageAnnotation]
	if annotationExists && annotationValue == "true" && infra.Zoned {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("zoned"), infra.Zoned, fmt.Sprintf("specifying a zoned cluster and having a the %q annotation is not possible", azure.ShootVmoUsageAnnotation)))
	}

	networksPath := fldPath.Child("networks")

	workerCIDR := cidrvalidation.NewCIDR(infra.Networks.Workers, networksPath.Child("workers"))

	allErrs = append(allErrs, cidrvalidation.ValidateCIDRParse(workerCIDR)...)
	allErrs = append(allErrs, cidrvalidation.ValidateCIDRIsCanonical(networksPath.Child("workers"), infra.Networks.Workers)...)

	if (infra.Networks.VNet.Name != nil && infra.Networks.VNet.ResourceGroup == nil) || (infra.Networks.VNet.Name == nil && infra.Networks.VNet.ResourceGroup != nil) {
		allErrs = append(allErrs, field.Invalid(networksPath.Child("vnet"), infra.Networks.VNet, "specifying an existing vnet name require a vnet name and vnet resource group"))
	} else if infra.Networks.VNet.Name != nil && infra.Networks.VNet.ResourceGroup != nil {
		if infra.Networks.VNet.CIDR != nil {
			allErrs = append(allErrs, field.Invalid(networksPath.Child("vnet", "cidr"), *infra.Networks.VNet.ResourceGroup, "specifying a cidr for an existing vnet is not possible"))
		}
		if infra.ResourceGroup != nil && *infra.Networks.VNet.ResourceGroup == infra.ResourceGroup.Name {
			allErrs = append(allErrs, field.Invalid(networksPath.Child("vnet", "resourceGroup"), *infra.Networks.VNet.ResourceGroup, "the vnet resource group must not be the same as the cluster resource group"))
		}
	} else {
		cidrPath := networksPath.Child("vnet", "cidr")
		if infra.Networks.VNet.CIDR == nil {
			// Use worker/subnet cidr as cidr for the vnet.
			allErrs = append(allErrs, workerCIDR.ValidateSubset(nodes)...)
			allErrs = append(allErrs, workerCIDR.ValidateNotSubset(pods, services)...)
		} else {
			vpcCIDR := cidrvalidation.NewCIDR(*(infra.Networks.VNet.CIDR), cidrPath)
			allErrs = append(allErrs, vpcCIDR.ValidateParse()...)
			allErrs = append(allErrs, vpcCIDR.ValidateSubset(nodes)...)
			allErrs = append(allErrs, vpcCIDR.ValidateSubset(workerCIDR)...)
			allErrs = append(allErrs, vpcCIDR.ValidateNotSubset(pods, services)...)
			allErrs = append(allErrs, cidrvalidation.ValidateCIDRIsCanonical(cidrPath, *infra.Networks.VNet.CIDR)...)
		}
	}

	// TODO(dkistner) Remove once we proceed with multiple AvailabilitySet support.
	// Currently we will not offer Nat Gateway for non zoned/AvailabilitySet based
	// clusters as the NatGateway is not compatible with Basic LoadBalancer and
	// we would need Standard LoadBalancers also in combination with AvailabilitySets.
	// For the multiple AvailabilitySet approach we would always need
	// a Standard LoadBalancer and a NatGateway.
	if !infra.Zoned && !annotationExists && infra.Networks.NatGateway != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networks", "natGateway"), infra.Networks.NatGateway, "NatGateway is currently only supported for zoned cluster"))
	}

	if infra.Networks.NatGateway != nil &&
		infra.Networks.NatGateway.IdleConnectionTimeoutMinutes != nil &&
		(*infra.Networks.NatGateway.IdleConnectionTimeoutMinutes < natGatewayMinTimeoutInMinutes || *infra.Networks.NatGateway.IdleConnectionTimeoutMinutes > natGatewayMaxTimeoutInMinutes) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("networks", "natGateway", "idleConnectionTimeoutMinutes"), *infra.Networks.NatGateway.IdleConnectionTimeoutMinutes, "idleConnectionTimeoutMinutes values must range between 4 and 120"))
	}

	if infra.Identity != nil && (infra.Identity.Name == "" || infra.Identity.ResourceGroup == "") {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("identity"), infra.Identity, "specifying an identity requires the name of the identity and the resource group which hosts the identity"))
	}

	if nodes != nil {
		allErrs = append(allErrs, nodes.ValidateSubset(workerCIDR)...)
	}

	return allErrs
}

// ValidateInfrastructureConfigUpdate validates a InfrastructureConfig object.
func ValidateInfrastructureConfigUpdate(oldConfig, newConfig *apisazure.InfrastructureConfig, providerPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, apivalidation.ValidateImmutableField(newConfig.ResourceGroup, oldConfig.ResourceGroup, providerPath.Child("resourceGroup"))...)
	allErrs = append(allErrs, apivalidation.ValidateImmutableField(newConfig.Networks.VNet, oldConfig.Networks.VNet, providerPath.Child("networks").Child("vnet"))...)
	allErrs = append(allErrs, apivalidation.ValidateImmutableField(newConfig.Networks.Workers, oldConfig.Networks.Workers, providerPath.Child("networks").Child("workers"))...)

	if oldConfig.Zoned && !newConfig.Zoned {
		allErrs = append(allErrs, field.Forbidden(providerPath.Child("zoned"), "moving a zoned cluster to a non-zoned cluster is not allowed"))
	}

	return allErrs
}

// ValidateVMOConfigurationUpdate validates the VMO configuration on update.
func ValidateVMOConfigurationUpdate(oldConfig, newConfig *apisazure.InfrastructureConfig, oldShootAnnotations, newShootAnnotation map[string]string, providerPath, metaDataPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Check if old shoot has not the vmo annotation and forbid to add it.
	if _, exists := oldShootAnnotations[azure.ShootVmoUsageAnnotation]; !exists {
		_, annotationExists := newShootAnnotation[azure.ShootVmoUsageAnnotation]
		if annotationExists {
			allErrs = append(allErrs, field.Forbidden(metaDataPath.Child("annotations"), fmt.Sprintf("not allowed to add vmo annotation %q to an already existing cluster", azure.ShootVmoUsageAnnotation)))
		}
	}

	// Check if the olf shoot has the vmo annotation and forbid to remove or modify it. Also forbid to change to move the new shoot to a zoned if annotated to be a vmo one.
	if value, exists := oldShootAnnotations[azure.ShootVmoUsageAnnotation]; exists && value == "true" {
		if newConfig.Zoned {
			allErrs = append(allErrs, field.Invalid(providerPath.Child("zoned"), newConfig.Zoned, fmt.Sprintf("not allowed to switch to a zoned cluster when already using a vmo based cluster (via annotation %q)", azure.ShootVmoUsageAnnotation)))
		}

		annotationValue, annotationExists := newShootAnnotation[azure.ShootVmoUsageAnnotation]
		if !annotationExists {

			allErrs = append(allErrs, field.Forbidden(metaDataPath.Child("annotations"), fmt.Sprintf("not allowed to remove vmo annotation %q if it is already in use", azure.ShootVmoUsageAnnotation)))
		}
		if annotationExists && annotationValue != "true" {
			allErrs = append(allErrs, field.Invalid(metaDataPath.Child("annotations"), annotationValue, fmt.Sprintf("not allowed to modify the vmo annotation %q if it is already in use", azure.ShootVmoUsageAnnotation)))
		}
	}

	return allErrs
}
