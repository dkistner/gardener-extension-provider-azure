// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
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

package client

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
)

// Factory represents a factory to produce clients for various Azure services.
type Factory interface {
	Group() Group
	Storage() Storage
	StorageAccount() StorageAccount
	Vmss() Vmss
}

// Group represents an Azure group client.
type Group interface {
	CreateOrUpdate(context.Context, string, string) error
	DeleteIfExits(context.Context, string) error
}

// Storage represents an Azure storage account client.
type Storage interface {
	DeleteObjectsWithPrefix(context.Context, string, string) error
	CreateContainerIfNotExists(context.Context, string) error
	DeleteContainerIfExists(context.Context, string) error
}

// StorageAccount ...
type StorageAccount interface {
	CreateStorageAccount(context.Context, string, string, string) error
	ListStorageAccountKey(context.Context, string, string) (string, error)
}

// Vmss represents an Azure vmss client.
type Vmss interface {
	List(context.Context, string) ([]compute.VirtualMachineScaleSet, error)
	Get(context.Context, string, string) (*compute.VirtualMachineScaleSet, error)
	Create(context.Context, string, string, *compute.VirtualMachineScaleSet) (*compute.VirtualMachineScaleSet, error)
	Delete(context.Context, string, string) error
}

// AzureFactory is a implementation of Factory to produce clients for various Azure services.
type AzureFactory struct {
	group          resources.GroupsClient
	storageAccount storage.AccountsClient
	vmss           compute.VirtualMachineScaleSetsClient
	storageURL     azblob.ServiceURL
}

// StorageClient is a implementation of Storage for a storage account client.
type StorageClient struct {
	serviceURL *azblob.ServiceURL
}

// StorageAccountClient ...
type StorageAccountClient struct {
	client storage.AccountsClient
}

// GroupClient is a implementation of Group for a group client.
type GroupClient struct {
	client resources.GroupsClient
}

// VmssClient is a implementation of Vmss for a vmss client.
type VmssClient struct {
	client compute.VirtualMachineScaleSetsClient
}
