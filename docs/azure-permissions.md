# Gardener Azure operation permissions

This document lists the required resource provider operations/permissions which an Azure Shoot cluster require to be operational.
This list also contains the additional permissions which a Shoot cluster require to act as a Seed (backup bucket/entry storage account permissions).

The list of operations/permissions is split by Azure resource providers.

**Note:** This list of operations/permissions is a partial superset (not all `Microsoft.Compute/virtualMachineScaleSets` permissions are required) of the permissions which the Azure Kubernetes cloud provider (Cloud-Controller-Manager, CSI driver and in-tree disk drivers).
See here: https://kubernetes-sigs.github.io/cloud-provider-azure/topics/azure-permissions/


### Microsoft.Compute

```
# Required only for non-zonal AvailabilitySet based clusters.
Microsoft.Compute/availabilitySets/delete
Microsoft.Compute/availabilitySets/read
Microsoft.Compute/availabilitySets/write

Microsoft.Compute/disks/delete
Microsoft.Compute/disks/read
Microsoft.Compute/disks/write

Microsoft.Compute/locations/diskOperations/read
Microsoft.Compute/locations/operations/read
Microsoft.Compute/locations/vmSizes/read

Microsoft.Compute/snapshots/delete
Microsoft.Compute/snapshots/read
Microsoft.Compute/snapshots/write

Microsoft.Compute/virtualMachines/delete
Microsoft.Compute/virtualMachines/read
Microsoft.Compute/virtualMachines/start/action
Microsoft.Compute/virtualMachines/write

# Required only for non-zonal VirtualMachineScaleSet Orchestraion Mode VM based clusters.
Microsoft.Compute/virtualMachineScaleSets/delete
Microsoft.Compute/virtualMachineScaleSets/read
Microsoft.Compute/virtualMachineScaleSets/write

# Required to use images from a shared image gallery.
Microsoft.Compute/galleries/images/*/read
```

### Microsoft.ManagedIdentity
Required to assign a managed identiy to the machines of a cluster.

```
Microsoft.ManagedIdentity/userAssignedIdentities/assign/action
Microsoft.ManagedIdentity/userAssignedIdentities/read
```

### Microsoft.MarketplaceOrdering
Required to accept using Azure Marketplace images like gardenlinux.

```
Microsoft.MarketplaceOrdering/offertypes/publishers/offers/plans/agreements/read
Microsoft.MarketplaceOrdering/offertypes/publishers/offers/plans/agreements/write
```

### Microsoft.Network

```
Microsoft.Network/loadBalancers/delete
Microsoft.Network/loadBalancers/read
Microsoft.Network/loadBalancers/write
Microsoft.Network/loadBalancers/backendAddressPools/join/action

Microsoft.Network/natGateways/delete
Microsoft.Network/natGateways/read
Microsoft.Network/natGateways/write
Microsoft.Network/natGateways/join/action

Microsoft.Network/networkInterfaces/delete
Microsoft.Network/networkInterfaces/read
Microsoft.Network/networkInterfaces/write
Microsoft.Network/networkInterfaces/ipconfigurations/join/action
Microsoft.Network/networkInterfaces/ipconfigurations/read
Microsoft.Network/networkInterfaces/join/action

Microsoft.Network/networkSecurityGroups/delete
Microsoft.Network/networkSecurityGroups/read
Microsoft.Network/networkSecurityGroups/write
Microsoft.Network/networkSecurityGroups/join/action

Microsoft.Network/publicIPAddresses/delete
Microsoft.Network/publicIPAddresses/read
Microsoft.Network/publicIPAddresses/write
Microsoft.Network/publicIPAddresses/join/action

Microsoft.Network/routeTables/delete
Microsoft.Network/routeTables/read
Microsoft.Network/routeTables/write
Microsoft.Network/routeTables/join/action
Microsoft.Network/routeTables/routes/delete
Microsoft.Network/routeTables/routes/read
Microsoft.Network/routeTables/routes/write

Microsoft.Network/virtualNetworks/delete
Microsoft.Network/virtualNetworks/read
Microsoft.Network/virtualNetworks/write
Microsoft.Network/virtualNetworks/subnets/delete
Microsoft.Network/virtualNetworks/subnets/read
Microsoft.Network/virtualNetworks/subnets/write
Microsoft.Network/virtualNetworks/subnets/join/action
```

### Microsoft.Resources
Required to create, get and delete resource groups.

```
Microsoft.Resources/subscriptions/resourceGroups/delete
Microsoft.Resources/subscriptions/resourceGroups/read
Microsoft.Resources/subscriptions/resourceGroups/write
```

### Microsoft.Storage
Required for the `azurefile` storage class and for the management of the `BackupBuckets/BackupEntries` on Azure for Shoot clusters which act as Seed.

```
Microsoft.Storage/storageAccounts/read
Microsoft.Storage/storageAccounts/write
Microsoft.Storage/storageAccounts/delete

Microsoft.Storage/storageAccounts/blobServices/containers/delete
Microsoft.Storage/storageAccounts/blobServices/containers/read
Microsoft.Storage/storageAccounts/blobServices/containers/write
Microsoft.Storage/storageAccounts/blobServices/read
Microsoft.Storage/storageAccounts/listkeys/action

Microsoft.Storage/operations/read
```
