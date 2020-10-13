package machinesetclient

import azurecompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"

const machineSetTagKey = "machineset.azure.extensions.gardener.cloud"

// Clients ...
type Clients struct {
	vmo azurecompute.VirtualMachineScaleSetsClient
}

// NewClients ...
func NewClients(vmo *azurecompute.VirtualMachineScaleSetsClient) *Clients {
	return &Clients{
		vmo: *vmo,
	}
}
