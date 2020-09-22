package worker

import (
	"context"

	api "github.com/gardener/gardener-extension-provider-azure/pkg/apis/azure"
	"github.com/pkg/errors"

	"github.com/gardener/gardener-extension-provider-azure/pkg/internal/machineset"
)

func (w *workerDelegate) DeployMachineDependencies(ctx context.Context) error {
	infrastructureStatus, err := w.decodeAzureInfrastructureStatus()
	if err != nil {
		return err
	}

	workerProviderStatus, err := w.decodeWorkerProviderStatus()
	if err != nil {
		return err
	}

	if machineset.IsVMORequired(infrastructureStatus) {
		newVmoDependencies, err := w.reconcileVmoDependencies(ctx, infrastructureStatus, workerProviderStatus)
		if err != nil {
			return w.updateMachineDependenciesStatus(ctx, workerProviderStatus, newVmoDependencies, err)
		}
		return w.updateMachineDependenciesStatus(ctx, workerProviderStatus, newVmoDependencies, nil)
	}
	return nil
}

func (w *workerDelegate) CleanupMachineDependencies(ctx context.Context) error {
	infrastructureStatus, err := w.decodeAzureInfrastructureStatus()
	if err != nil {
		return err
	}

	workerProviderStatus, err := w.decodeWorkerProviderStatus()
	if err != nil {
		return err
	}

	if machineset.IsVMORequired(infrastructureStatus) {
		newVmoDependencies, err := w.cleanupVmoDependencies(ctx, infrastructureStatus, workerProviderStatus)
		if err != nil {
			return w.updateMachineDependenciesStatus(ctx, workerProviderStatus, newVmoDependencies, err)
		}
		return w.updateMachineDependenciesStatus(ctx, workerProviderStatus, newVmoDependencies, nil)
	}

	return nil
}

// Helper

func (w workerDelegate) updateMachineDependenciesStatus(ctx context.Context, workerStatus *api.WorkerStatus, vmoDependencies []api.VmoDependency, err error) error {
	workerStatus.VmoDependencies = vmoDependencies

	if statusUpdateErr := w.updateWorkerProviderStatus(ctx, workerStatus); statusUpdateErr != nil {
		err = errors.Wrapf(statusUpdateErr, err.Error())
	}
	if err != nil {
		return err
	}
	return nil
}
