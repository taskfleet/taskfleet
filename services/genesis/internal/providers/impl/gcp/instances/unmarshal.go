package gcpinstances

import (
	"fmt"
	"net"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

func unmarshalInstance(
	instance *computepb.Instance, manager *instances.Manager, project string,
) (providers.Instance, error) {
	// First, parse meta
	meta, err := unmarshalInstanceMeta(instance)
	if err != nil {
		return providers.Instance{}, err
	}

	// Then, parse the spec
	spec, err := unmarshalInstanceSpec(instance, manager)
	if err != nil {
		return providers.Instance{}, err
	}

	// Finally, parse the status
	status, err := unmarshalInstanceStatus(instance, project)
	if err != nil {
		return providers.Instance{}, err
	}

	result := providers.Instance{
		Meta:   meta,
		Spec:   spec,
		Status: status,
	}
	return result, nil
}

func unmarshalInstanceMeta(i *computepb.Instance) (providers.InstanceMeta, error) {
	id, err := uuid.Parse(strings.TrimPrefix(i.GetName(), "taskfleet-"))
	if err != nil {
		return providers.InstanceMeta{}, fmt.Errorf("failed to parse instance ID: %s", err)
	}
	return providers.InstanceMeta{
		ID:           id,
		ProviderID:   i.GetName(),
		ProviderZone: path.Base(i.GetZone()),
	}, nil
}

func unmarshalInstanceSpec(
	i *computepb.Instance, manager *instances.Manager,
) (providers.InstanceSpec, error) {
	// Get compute instance type -- for this, we need the correct manager and need to parse the
	// GPU resources to translate them into the instance name
	var gpuResources *instances.GPUResources
	if len(i.GetGuestAccelerators()) > 0 {
		gpuKind, err := typedefs.GPUKindUnmarshalProviderGcp(
			path.Base(i.GetGuestAccelerators()[0].GetAcceleratorType()),
		)
		if err != nil {
			return providers.InstanceSpec{}, providers.NewFatalError("unknown kind of gpu", err)
		}
		gpuResources = &instances.GPUResources{
			Kind:  gpuKind,
			Count: uint16(i.GetGuestAccelerators()[0].GetAcceleratorCount()),
		}
	}

	InstanceType, err := manager.Get(
		extendedInstanceTypeName(path.Base(i.GetMachineType()), gpuResources),
	)
	if err != nil {
		// Raising a fatal error here since any instance type should be known since it was created
		// by Genesis
		return providers.InstanceSpec{}, providers.NewFatalError("unknown instance type", err)
	}

	// Spot status is determined by the instance's provisioning type
	isSpot := i.GetScheduling().GetProvisioningModel() == "SPOT"

	// Finalize the instance spec
	return providers.InstanceSpec{
		InstanceType: InstanceType,
		IsSpot:       isSpot,
	}, nil
}

func unmarshalInstanceStatus(
	i *computepb.Instance, project string,
) (providers.InstanceStatus, error) {
	// Get creation timestamp
	createdAt, err := time.Parse(time.RFC3339, i.GetCreationTimestamp())
	if err != nil {
		return providers.InstanceStatus{},
			providers.NewFatalError("instance returned invalid creation timestamp", err)
	}

	// Get network interface to read network configuration
	if len(i.GetNetworkInterfaces()) == 0 {
		return providers.InstanceStatus{},
			providers.NewFatalError("instance is not attached to any networks", nil)
	} else if len(i.GetNetworkInterfaces()) > 1 {
		return providers.InstanceStatus{},
			providers.NewFatalError("instance attached to multiple networks", nil)
	}
	iface := i.GetNetworkInterfaces()[0]

	// Get internal IP
	internalIP := net.ParseIP(iface.GetNetworkIP())
	if internalIP == nil {
		return providers.InstanceStatus{}, providers.NewFatalError("invalid network IP", nil)
	}

	// Get external IP
	var externalIP net.IP
	if len(iface.GetAccessConfigs()) > 0 {
		externalIP = net.ParseIP(iface.GetAccessConfigs()[0].GetNatIP())
	}

	// Combine into status
	status := providers.InstanceStatus{
		CreationTimestamp: createdAt,
		Network: providers.InstanceNetworkStatus{
			InternalIP:       internalIP,
			ExternalIP:       externalIP,
			InternalHostname: i.GetHostname(),
		},
	}
	return status, nil
}
