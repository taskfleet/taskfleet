package gcpinstances

import (
	"fmt"
	"net"
	"path"
	"time"

	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"google.golang.org/api/compute/v1"
)

func unmarshalInstance(
	instance *compute.Instance, managers map[string]*instances.Manager, project string,
) (providers.Instance, error) {
	// First, parse meta
	ref, err := unmarshalInstanceRef(instance)
	if err != nil {
		return providers.Instance{}, err
	}

	// Then, parse the spec
	spec, err := unmarshalInstanceSpec(instance, managers)
	if err != nil {
		return providers.Instance{}, err
	}

	// Finally, parse the status
	status, err := unmarshalInstanceStatus(instance, project)
	if err != nil {
		return providers.Instance{}, err
	}

	result := providers.Instance{
		Ref:    ref,
		Spec:   spec,
		Status: status,
	}
	return result, nil
}

func unmarshalInstanceRef(i *compute.Instance) (providers.InstanceRef, error) {
	return providers.InstanceRefFromCommonName(i.Name, path.Base(i.Zone))
}

func unmarshalInstanceSpec(
	i *compute.Instance, managers map[string]*instances.Manager,
) (providers.InstanceSpec, error) {
	// Get compute instance type -- for this, we need the correct manager and need to parse the
	// GPU resources
	var gpuResources *instances.GPUResources
	if len(i.GuestAccelerators) > 0 {
		gpuKind, err := typedefs.GPUKindFromProviderGcp(
			path.Base(i.GuestAccelerators[0].AcceleratorType),
		)
		if err != nil {
			return providers.InstanceSpec{}, providers.NewFatalError("unknown kind of gpu", err)
		}
		gpuResources = &instances.GPUResources{
			Kind:  gpuKind,
			Count: uint16(i.GuestAccelerators[0].AcceleratorCount),
		}
	}

	manager := managers[path.Base(i.Zone)]
	machineType, err := manager.Get(path.Base(i.MachineType), gpuResources)
	if err != nil {
		return providers.InstanceSpec{}, providers.NewFatalError("unknown machine type", err)
	}

	compute := providers.ComputeConfig{
		InstanceType: machineType,
		IsSpot:       i.Scheduling.Preemptible,
	}

	// Get metadata
	metadata := providers.MetadataConfig{
		Tags:       i.Tags.Items,
		Labels:     i.Labels,
		Attributes: map[string]string{},
	}
	for _, item := range i.Metadata.Items {
		metadata.Attributes[item.Key] = *item.Value
	}

	// Get security configuration
	security := providers.SecurityConfig{}
	if len(i.ServiceAccounts) > 0 {
		security = providers.SecurityConfig{
			ServiceAccountEmail: i.ServiceAccounts[0].Email,
		}
	}

	// Aggregate into spec
	spec := providers.InstanceSpec{
		Compute:  compute,
		Metadata: metadata,
		Security: security,
	}
	return spec, nil
}

func unmarshalInstanceStatus(
	i *compute.Instance, project string,
) (providers.InstanceStatus, error) {
	// Get creation timestamp
	createdAt, err := time.Parse(time.RFC3339, i.CreationTimestamp)
	if err != nil {
		return providers.InstanceStatus{},
			providers.NewFatalError("instance returned invalid creation timestamp", err)
	}

	// Get network configuration
	if len(i.NetworkInterfaces) == 0 {
		return providers.InstanceStatus{},
			providers.NewFatalError("instance is not attached to any networks", nil)
	}
	internalIP := net.ParseIP(i.NetworkInterfaces[0].NetworkIP)

	if len(i.NetworkInterfaces[0].AccessConfigs) == 0 {
		return providers.InstanceStatus{},
			providers.NewFatalError("instance does not have access to the internet", nil)
	}
	externalIP := net.ParseIP(i.NetworkInterfaces[0].AccessConfigs[0].NatIP)
	hostname := fmt.Sprintf("%s.%s.c.%s.internal", i.Name, path.Base(i.Zone), project)

	// Combine into status
	status := providers.InstanceStatus{
		CreationTimestamp: createdAt,
		Network: providers.InstanceNetworkStatus{
			InternalIP:       internalIP,
			ExternalIP:       externalIP,
			InternalHostname: hostname,
		},
	}
	return status, nil
}
