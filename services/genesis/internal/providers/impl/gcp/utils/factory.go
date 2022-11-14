package gcputils

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	"google.golang.org/api/option"
)

type ClientFactory interface {
	// AcceleratorTypes returns a client to interact with GCP accelerator types.
	AcceleratorTypes() *compute.AcceleratorTypesClient
	// DiskTypes returns a client to interact with GCP disk types.
	DiskTypes() *compute.DiskTypesClient
	// Instances returns a client to interact with GCP instances.
	Instances() *compute.InstancesClient
	// MachineTypes returns a client to interact with GCP machine types.
	MachineTypes() *compute.MachineTypesClient
	// Networks returns a client to interact with GCP networks.
	Networks() *compute.NetworksClient
	// Zones returns a client to interact with GCP zones.
	Zones() *compute.ZonesClient
}

//-------------------------------------------------------------------------------------------------
// REAL FACTORY
//-------------------------------------------------------------------------------------------------

type clientFactory struct {
	ctx     context.Context
	options []option.ClientOption

	acceleratorTypes *compute.AcceleratorTypesClient
	diskTypes        *compute.DiskTypesClient
	instances        *compute.InstancesClient
	machineTypes     *compute.MachineTypesClient
	networks         *compute.NetworksClient
	zones            *compute.ZonesClient
}

func NewClientFactory(ctx context.Context, options ...option.ClientOption) ClientFactory {
	return &clientFactory{
		ctx:     ctx,
		options: options,
	}
}

func (f *clientFactory) AcceleratorTypes() *compute.AcceleratorTypesClient {
	if f.acceleratorTypes != nil {
		return f.acceleratorTypes
	}
	acceleratorTypes, err := compute.NewAcceleratorTypesRESTClient(f.ctx, f.options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize accelerator types client: %s", err))
	}
	go deferClose(f.ctx, acceleratorTypes)
	f.acceleratorTypes = acceleratorTypes
	return f.acceleratorTypes
}

func (f *clientFactory) DiskTypes() *compute.DiskTypesClient {
	if f.diskTypes != nil {
		return f.diskTypes
	}
	diskTypes, err := compute.NewDiskTypesRESTClient(f.ctx, f.options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize disks client: %s", err))
	}
	go deferClose(f.ctx, diskTypes)
	f.diskTypes = diskTypes
	return f.diskTypes
}

func (f *clientFactory) Instances() *compute.InstancesClient {
	if f.instances != nil {
		return f.instances
	}
	instances, err := compute.NewInstancesRESTClient(f.ctx, f.options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize instances client: %s", err))
	}
	go deferClose(f.ctx, instances)
	f.instances = instances
	return f.instances
}

func (f *clientFactory) MachineTypes() *compute.MachineTypesClient {
	if f.machineTypes != nil {
		return f.machineTypes
	}
	machineTypes, err := compute.NewMachineTypesRESTClient(f.ctx, f.options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize machine types client: %s", err))
	}
	go deferClose(f.ctx, machineTypes)
	f.machineTypes = machineTypes
	return f.machineTypes
}

func (f *clientFactory) Networks() *compute.NetworksClient {
	if f.networks != nil {
		return f.networks
	}
	networks, err := compute.NewNetworksRESTClient(f.ctx, f.options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize networks client: %s", err))
	}
	go deferClose(f.ctx, networks)
	f.networks = networks
	return f.networks
}

func (f *clientFactory) Zones() *compute.ZonesClient {
	if f.zones != nil {
		return f.zones
	}
	zones, err := compute.NewZonesRESTClient(f.ctx, f.options...)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize zones client: %s", err))
	}
	go deferClose(f.ctx, zones)
	f.zones = zones
	return f.zones
}

//-------------------------------------------------------------------------------------------------
// CLOSEABLE
//-------------------------------------------------------------------------------------------------

type closeable interface {
	Close() error
}

func deferClose(ctx context.Context, objects ...closeable) {
	<-ctx.Done()
	for _, obj := range objects {
		obj.Close() // nolint:errcheck
	}
}
