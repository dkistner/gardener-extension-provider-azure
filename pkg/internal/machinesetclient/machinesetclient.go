package machinesetclient

import (
	"context"

	azurecompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
)

const machineSetTagKey = "machineset.azure.extensions.gardener.cloud"

// MachineSetClient is an interface for a Azure MachineSet client.
type MachineSetClient interface {
	// ListVMOs list all Gardener managed VMOs in a given resource group.
	ListVMOs(context.Context, string) ([]azurecompute.VirtualMachineScaleSet, error)

	// GetVMO fetches a VMO based on given resourceGroupName and name.
	GetVMO(context.Context, string, string) (*azurecompute.VirtualMachineScaleSet, error)

	// CreateVMO creates a VMO with passed configuration like resource group, name, region and fault domain count.
	CreateVMO(context.Context, string, string, string, int32) (*azurecompute.VirtualMachineScaleSet, error)

	// DeleteVMO deletes a VMO based on passed resource group and name.
	DeleteVMO(context.Context, string, string) error
}

type clients struct {
	vmo azurecompute.VirtualMachineScaleSetsClient
}

// NewMachineSetClients will return a new MachineSetClient.
func NewMachineSetClients(vmo *azurecompute.VirtualMachineScaleSetsClient) MachineSetClient {
	return &clients{
		vmo: *vmo,
	}
}
