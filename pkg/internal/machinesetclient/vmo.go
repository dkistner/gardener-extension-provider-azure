package machinesetclient

import (
	"context"
	"fmt"
	"net/http"

	azureapi "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"
	"github.com/gardener/gardener-extension-provider-azure/pkg/internal"
	"github.com/gardener/gardener/pkg/utils"

	azurecompute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"

	"k8s.io/utils/pointer"
)

// VirtualMachine ScaleSet Orchestration Mode VM (VMO)

// UseVMOAnnotation is an annotation for the Shoot resource which indicate
// if vmo should be used for non-zonal cluster instead of a primary AvailabilitySet.
const UseVMOAnnotation = "azure.provider.extensions.gardener.cloud/enable-vmo"

// IsVMORequired determines if VMO is required.
func IsVMORequired(infrastructureStatus *azureapi.InfrastructureStatus) bool {
	if infrastructureStatus.Zoned || len(infrastructureStatus.AvailabilitySets) > 0 {
		return false
	}
	return true
}

// ListVMOs will list all Gardener managed VMOs in a given resource group.
func (c *clients) ListVMOs(ctx context.Context, resourceGroupName string) ([]azurecompute.VirtualMachineScaleSet, error) {
	pages, err := c.vmo.List(ctx, resourceGroupName)
	if err != nil {
		return nil, err
	}

	var vmoList []azurecompute.VirtualMachineScaleSet
	for pages.NotDone() {
		for _, vmo := range pages.Values() {
			if _, hasTag := vmo.Tags[machineSetTagKey]; hasTag {
				vmoList = append(vmoList, vmo)
			}
		}
		if err := pages.NextWithContext(ctx); err != nil {
			return nil, err
		}
	}

	return vmoList, nil
}

// GetVMO will fetch a VMO based on given resourceGroupName and name.GetVMO.
func (c *clients) GetVMO(ctx context.Context, resourceGroupName, name string) (*azurecompute.VirtualMachineScaleSet, error) {
	vmo, err := c.vmo.Get(ctx, resourceGroupName, name)
	if err != nil {
		if internal.AzureAPIErrorNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return &vmo, nil
}

// CreateVMO will create a VMO with passed configuration like resource group, name, region and fault domain count.
func (c *clients) CreateVMO(ctx context.Context, resourceGroupName, name, region string, faultDomainCount int32) (*azurecompute.VirtualMachineScaleSet, error) {
	// Generate a random string as suffix for the VMO name.
	randomString, err := utils.GenerateRandomString(8)
	if err != nil {
		return nil, err
	}

	var (
		vmoName       = fmt.Sprintf("vmo-%s-%s", name, randomString)
		vmoParameters = azurecompute.VirtualMachineScaleSet{
			Location: &region,
			VirtualMachineScaleSetProperties: &azurecompute.VirtualMachineScaleSetProperties{

				SinglePlacementGroup:     pointer.BoolPtr(false),
				PlatformFaultDomainCount: &faultDomainCount,
			},
			Tags: map[string]*string{
				machineSetTagKey: pointer.StringPtr("1"),
			},
		}
	)

	future, err := c.vmo.CreateOrUpdate(ctx, resourceGroupName, vmoName, vmoParameters)
	if err != nil {
		return nil, err
	}
	if err := future.WaitForCompletionRef(ctx, c.vmo.Client); err != nil {
		return nil, err
	}
	vmo, err := future.Result(c.vmo)
	if err != nil {
		return nil, err
	}
	return &vmo, nil
}

// DeleteVMO will delete a VMO based on passed resource group and name.
func (c *clients) DeleteVMO(ctx context.Context, resourceGroupName, name string) error {
	future, err := c.vmo.Delete(ctx, resourceGroupName, name)
	if err != nil {
		return err
	}
	if err := future.WaitForCompletionRef(ctx, c.vmo.Client); err != nil {
		return err
	}
	result, err := future.Result(c.vmo)
	if err != nil {
		return err
	}
	if result.StatusCode == http.StatusOK || result.StatusCode == http.StatusAccepted || result.StatusCode == http.StatusNoContent {
		return nil
	}
	return fmt.Errorf("Deletion of vmo %s failed. statuscode=%d", name, result.StatusCode)
}
