package vsphere

import (
	"context"

	"github.com/cycloidio/terracognita/filter"
	"github.com/cycloidio/terracognita/provider"
	"github.com/vmware/govmomi/vim25/mo"
)

// ResourceType is the type used to define all the Resources
// from the Provider
type ResourceType int

//go:generate enumer -type ResourceType -addprefix vsphere_ -transform snake -linecomment
const (
	computeCluster ResourceType = iota // compute_cluster
	resourcePool                       // resource_pool

	// Inventory
	datacenter
	folder

	// Networking
	distributedPortGroup     // distributed_port_group
	distributedVirtualSwitch // distributed_virtual_switch

	// Storage
	datastoreCluster // datastore_cluster

	// Virtual Machine
	virtualMachine // virtual_machine
)

type rtFn func(ctx context.Context, vs *vsphere, vm *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error)

var (
	resources = map[ResourceType]rtFn{
		computeCluster:           getComputeClusters,
		resourcePool:             getResourcePools,
		datacenter:               getDatacenters,
		folder:                   getFolders,
		distributedPortGroup:     getDistributedPortGroups,
		distributedVirtualSwitch: getDistributedVirtualSwitches,
		datastoreCluster:         getDatastoreClusters,
		virtualMachine:           getVirtualMachines,
	}
)

func getDatastoreClusters(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"Datastore"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all Datastores
	var dss []mo.Datastore
	err = v.Retrieve(ctx, []string{"Datastore"}, []string{}, &dss)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, ds := range dss {
		r := provider.NewResource(ds.ManagedEntity.ExtensibleManagedObject.Self.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getDistributedVirtualSwitches(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"distributedVirtualSwitch"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all DistributedVirtualSwitches
	var dvss []mo.DistributedVirtualSwitch
	err = v.Retrieve(ctx, []string{"distributedVirtualSwitch"}, []string{}, &dvss)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, dvs := range dvss {
		r := provider.NewResource(dvs.ManagedEntity.ExtensibleManagedObject.Self.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getDistributedPortGroups(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"DistributedVirtualPortgroup"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all DistributedVirtualPortgroups
	var dvps []mo.DistributedVirtualPortgroup
	err = v.Retrieve(ctx, []string{"DistributedVirtualPortgroup"}, []string{}, &dvps)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, dvp := range dvps {
		r := provider.NewResource(dvp.ManagedEntity.ExtensibleManagedObject.Self.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getFolders(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"folder"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all folders
	var folders []mo.Folder
	err = v.Retrieve(ctx, []string{"folder"}, []string{}, &folders)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, f := range folders {
		r := provider.NewResource(f.ManagedEntity.ExtensibleManagedObject.Self.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getDatacenters(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"datacenter"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all datacenters
	var dss []mo.Datacenter
	err = v.Retrieve(ctx, []string{"datacenter"}, []string{}, &dss)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, ds := range dss {
		r := provider.NewResource(ds.ManagedEntity.ExtensibleManagedObject.Self.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getVirtualMachines(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"virtualMachine"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all virtual machines
	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"virtualMachine"}, []string{}, &vms)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, vm := range vms {
		r := provider.NewResource(vm.Config.InstanceUuid, resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getResourcePools(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"resourcePool"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all resource pools
	var rps []mo.ResourcePool
	err = v.Retrieve(ctx, []string{"resourcePool"}, []string{}, &rps)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, rp := range rps {
		r := provider.NewResource(rp.Config.Entity.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}

func getComputeClusters(ctx context.Context, vs *vsphere, r *reader, resourceType string, filters *filter.Filter) ([]provider.Resource, error) {
	v, err := r.CreateContainerView(ctx, r.Common.Client().ServiceContent.RootFolder, []string{"ClusterComputeResource"}, true)
	if err != nil {
		return nil, err
	}

	defer v.Destroy(ctx)

	// Retrieve summary property for all compute clusters
	var ccrs []mo.ClusterComputeResource
	err = v.Retrieve(ctx, []string{"ClusterComputeResource"}, []string{}, &ccrs)
	if err != nil {
		return nil, err
	}

	resources := make([]provider.Resource, 0)
	for _, ccr := range ccrs {
		r := provider.NewResource(ccr.ComputeResource.ManagedEntity.ExtensibleManagedObject.Self.String(), resourceType, vs)
		resources = append(resources, r)
	}

	return resources, nil
}
