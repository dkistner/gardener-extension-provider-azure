package worker

import (
	"context"
	"reflect"

	azureapi "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"
	azureapihelper "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure/helper"
	"github.com/gardener/gardener-extension-provider-azure/pkg/internal/machinesetclient"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
)

func (w *workerDelegate) reconcileVmoDependencies(ctx context.Context, infrastructureStatus *azureapi.InfrastructureStatus, workerProviderStatus *azureapi.WorkerStatus) ([]azureapi.VmoDependency, error) {
	var (
		vmoDependencies  = copyVmoDependencies(workerProviderStatus)
		machineSetClient = machinesetclient.NewMachineSetClients(&w.azureClients.Vmo)
	)

	faultDomainCount, err := azureapihelper.FindDomainCountByRegion(w.cloudProfileConfig.CountFaultDomains, w.worker.Spec.Region)
	if err != nil {
		return nil, err
	}

	// Deploy workerpool dependencies and store their status to be persistent in the worker provider status.
	for _, workerPool := range w.worker.Spec.Pools {
		vmoDependencyStatus, err := w.reconcileVMO(ctx, machineSetClient, vmoDependencies, infrastructureStatus.ResourceGroup.Name, workerPool.Name, faultDomainCount)
		if err != nil {
			return vmoDependencies, err
		}
		vmoDependencies = appendVmoDependency(vmoDependencies, vmoDependencyStatus)
	}

	return vmoDependencies, nil
}

func (w *workerDelegate) reconcileVMO(ctx context.Context, client machinesetclient.MachineSetClient, dependencies []azureapi.VmoDependency, resourceGroupName, workerPoolName string, faultDomainCount int32) (*azureapi.VmoDependency, error) {
	var (
		existingDependency *azureapi.VmoDependency
		vmo                *compute.VirtualMachineScaleSet
		err                error
	)

	// Check if there is already a VMO dependency object for the workerpool in the status.
	for _, dep := range dependencies {
		if dep.PoolName == workerPoolName {
			existingDependency = &dep
			break
		}
	}

	// Try to fetch the VMO from Azure as it exists in the status.
	if existingDependency != nil {
		vmo, err = client.GetVMO(ctx, resourceGroupName, existingDependency.Name)
		if err != nil {
			return nil, err
		}
	}

	// VMO does not exists. Create it.
	if vmo == nil {
		newVMO, err := client.CreateVMO(ctx, resourceGroupName, workerPoolName, w.worker.Spec.Region, faultDomainCount)
		if err != nil {
			return nil, err
		}
		return &azureapi.VmoDependency{
			ID:       *newVMO.ID,
			Name:     *newVMO.Name,
			PoolName: workerPoolName,
		}, nil
	}

	// VMO already exists. Check if the fault domain count configuration has been changed.
	// If yes then it is required to create a new VMO with the correct configuration.
	if *vmo.PlatformFaultDomainCount != faultDomainCount {
		newVMO, err := client.CreateVMO(ctx, resourceGroupName, workerPoolName, w.worker.Spec.Region, faultDomainCount)
		if err != nil {
			return nil, err
		}
		return &azureapi.VmoDependency{
			ID:       *newVMO.ID,
			Name:     *newVMO.Name,
			PoolName: workerPoolName,
		}, nil
	}

	return &azureapi.VmoDependency{
		ID:       *vmo.ID,
		Name:     *vmo.Name,
		PoolName: workerPoolName,
	}, nil
}

func (w *workerDelegate) cleanupVmoDependencies(ctx context.Context, infrastructureStatus *azureapi.InfrastructureStatus, workerProviderStatus *azureapi.WorkerStatus) ([]azureapi.VmoDependency, error) {
	var (
		machineSetClient = machinesetclient.NewMachineSetClients(&w.azureClients.Vmo)
		vmoDependencies  = copyVmoDependencies(workerProviderStatus)
	)

	// Cleanup VMO dependencies which are not tracked in the worker provider status anymore.
	if err := cleanupOrphanVMODependencies(ctx, machineSetClient, workerProviderStatus.VmoDependencies, infrastructureStatus.ResourceGroup.Name); err != nil {
		return vmoDependencies, err
	}

	// Delete all vmo workerpool dependencies as the Worker is intended to be deleted.
	if w.worker.ObjectMeta.DeletionTimestamp != nil {
		for _, dependency := range workerProviderStatus.VmoDependencies {
			if err := machineSetClient.DeleteVMO(ctx, infrastructureStatus.ResourceGroup.Name, dependency.Name); err != nil {
				return vmoDependencies, err
			}
			vmoDependencies = removeVmoDependency(vmoDependencies, dependency)
		}
		return vmoDependencies, nil
	}

	for _, dependency := range workerProviderStatus.VmoDependencies {
		var workerPoolExists = false
		for _, pool := range w.worker.Spec.Pools {
			if pool.Name == dependency.PoolName {
				workerPoolExists = true
				break
			}
		}
		if workerPoolExists {
			continue
		}

		// Delete the dependency as no corresponding workerpool exist anymore.
		if err := machineSetClient.DeleteVMO(ctx, infrastructureStatus.ResourceGroup.Name, dependency.Name); err != nil {
			return vmoDependencies, err
		}
		vmoDependencies = removeVmoDependency(vmoDependencies, dependency)
	}
	return vmoDependencies, nil
}

func cleanupOrphanVMODependencies(ctx context.Context, client machinesetclient.MachineSetClient, dependencies []azureapi.VmoDependency, resourceGroupName string) error {
	vmoList, err := client.ListVMOs(ctx, resourceGroupName)
	if err != nil {
		return err
	}

	for _, vmo := range vmoList {
		vmoExists := false
		for _, dependency := range dependencies {
			if *vmo.ID == dependency.ID {
				vmoExists = true
				break
			}
		}
		if !vmoExists {
			if err := client.DeleteVMO(ctx, resourceGroupName, *vmo.Name); err != nil {
				return err
			}
		}
	}
	return nil
}

// VMO Helper

func copyVmoDependencies(workerStatus *azureapi.WorkerStatus) []azureapi.VmoDependency {
	statusCopy := workerStatus.DeepCopy()
	return statusCopy.VmoDependencies
}

// appendVmoDependency appends a new vmo to the dependency list.
// If the dependency list contains already a vmo for the workerpool then the
// existing vmo object will be replaced by the given vmo object.
func appendVmoDependency(dependencies []azureapi.VmoDependency, dependency *azureapi.VmoDependency) []azureapi.VmoDependency {
	var idx *int
	for i, dep := range dependencies {
		if dep.PoolName == dependency.PoolName {
			idx = &i
			break
		}
	}
	if idx != nil {
		dependencies[*idx] = *dependency
	} else {
		dependencies = append(dependencies, *dependency)
	}
	return dependencies
}

// removeVmoDependency will remove a given vmo dependency from the passed list of dependencies.
func removeVmoDependency(dependencies []azureapi.VmoDependency, dependency azureapi.VmoDependency) []azureapi.VmoDependency {
	var idx *int
	for i, dep := range dependencies {
		if reflect.DeepEqual(dependency, dep) {
			idx = &i
			break
		}
	}
	if idx != nil {
		return append(dependencies[:*idx], dependencies[*idx+1:]...)
	}
	return dependencies
}
