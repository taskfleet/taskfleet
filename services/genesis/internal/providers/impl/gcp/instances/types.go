package gcpinstances

import (
	"context"
	"fmt"
	"path"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/borchero/zeus/pkg/zeus"
	"go.taskfleet.io/packages/jack"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"go.uber.org/zap"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

func findAvailableInstanceTypes(
	ctx context.Context,
	service *compute.MachineTypesClient,
	zones *gcpzones.Client,
	projectID string,
) (map[string][]instances.Type, error) {
	// Get instance types by zone
	result := map[string][]instances.Type{}
	for _, zone := range zones.List() {
		result[zone.Name] = make([]instances.Type, 0)
	}

	// First, we fetch the aggregated list of all machine types. We need to send two requests here
	// since the response data does not provide the CPU architecture.
	processResponse := func(
		architecture typedefs.CPUArchitecture,
	) func(pair compute.MachineTypesScopedListPair) error {
		return func(pair compute.MachineTypesScopedListPair) error {
			zone := path.Base(pair.Key)

			// If the zone is not available, continue
			if _, ok := result[zone]; !ok {
				return nil
			}

			// Then iterate over available machine types
			for _, item := range pair.Value.MachineTypes {
				if instance := tryUnmarshalInstanceType(ctx, item, architecture); instance != nil {
					result[zone] = append(result[zone], *instance)
				}
			}
			return nil
		}
	}

	// 1/2) x86_64
	it := service.AggregatedList(ctx, &computepb.AggregatedListMachineTypesRequest{
		Project: projectID,
		Filter:  proto.String("architecture=\"x86_64\""),
	})
	if err := gcputils.Iterate[compute.MachineTypesScopedListPair](
		ctx, it, processResponse(typedefs.ArchitectureX86),
	); err != nil {
		return nil, fmt.Errorf("failed to fetch machine types with x86 architecture: %s", err)
	}

	// 2/2) arm64
	it = service.AggregatedList(ctx, &computepb.AggregatedListMachineTypesRequest{
		Project: projectID,
		Filter:  proto.String("architecture=\"arm64\""),
	})
	if err := gcputils.Iterate[compute.MachineTypesScopedListPair](
		ctx, it, processResponse(typedefs.ArchitectureX86),
	); err != nil {
		return nil, fmt.Errorf("failed to fetch machine types with arm architecture: %s", err)
	}

	// Then, we add all instance types which provide GPUs. At the moment, GPUs (other than the
	// A100) can only be added to N1 instances. The instances with A100 GPU were already added
	// before since the machine types have a fixed set of accelerators attached.
	for _, zone := range zones.List() {
		// Get all N1 instances for the current zone
		n1Instances := make([]instances.Type, 0)
		for _, instance := range result[zone.Name] {
			if strings.HasPrefix(instance.Name, "n1-") {
				n1Instances = append(n1Instances, instance)
			}
		}

		// For each available GPU, add all possible configurations
		for _, gpu := range zone.GPUs {
			accelerator, err := zones.GetAccelerator(zone.Name, gpu)
			if err != nil {
				continue // no possible configurations if GPU type not available
			}
			gpuInstances := explodeAvailableGpuInstanceTypes(
				zone.Name, result[zone.Name], gpu, accelerator.MaxCount(),
			)
			result[zone.Name] = append(result[zone.Name], gpuInstances...)
		}
	}

	// Now, we are done
	return result, nil
}

func tryUnmarshalInstanceType(
	ctx context.Context,
	item *computepb.MachineType,
	architecture typedefs.CPUArchitecture,
) *instances.Type {
	// When iterating over the returned machines, we exclude some sets of machines:
	// * Machines with shared cores since they provide too little resources for proper jobs
	// * Deprecated machines (even if they are still active)
	// * M2 machines since they require purchasing committed usage plans for longer usage, see
	//   https://cloud.google.com/compute/docs/memory-optimized-machines#m2_machine_types
	if item.GetIsSharedCpu() ||
		item.GetDeprecated() != nil ||
		strings.HasPrefix(item.GetName(), "m2-") {
		return nil
	}

	// Build the resources
	resources := instances.Resources{
		CPUCount:  uint16(item.GetGuestCpus()),
		MemoryMiB: uint32(item.GetMemoryMb()),
	}
	if len(item.Accelerators) > 0 {
		accelerator := item.Accelerators[0]
		gpuKind, err := typedefs.GPUKindUnmarshalProviderGcp(accelerator.GetGuestAcceleratorType())

		// Skip this instance if it cannot be parsed
		if err != nil {
			zeus.Logger(ctx).Warn(
				"skipping GCP machine type due to unknown accelerator",
				zap.String("machine_type", item.GetName()),
				zap.String("gpu", accelerator.GetGuestAcceleratorType()),
			)
			return nil
		}

		resources.GPU = &instances.GPUResources{
			Kind: gpuKind, Count: uint16(accelerator.GetGuestAcceleratorCount()),
		}
	}

	// Return the instance
	return &instances.Type{
		Name:         item.GetName(),
		UID:          item.GetSelfLink(),
		Resources:    resources,
		Architecture: architecture,
	}
}

func explodeAvailableGpuInstanceTypes(
	zone string, n1Instances []instances.Type, gpu typedefs.GPUKind, maxGpuCount uint16,
) []instances.Type {
	result := make([]instances.Type, 0)
	for count := uint16(1); count <= maxGpuCount; count *= 2 {
		maxCPU, maxMem := maxCpuAndMemoryForGpu(gpu, count, zone)
		for _, instance := range n1Instances {
			if instance.CPUCount <= maxCPU && instance.MemoryMiB <= maxMem*1024 {
				// We need to create our own name here to distinguish instance types. This is not
				// ideal but makes the design of providers much easier...
				resources := instances.GPUResources{Kind: gpu, Count: count}
				result = append(result, instances.Type{
					Name: extendedInstanceTypeName(instance.Name, &resources),
					UID:  instance.UID,
					Resources: instances.Resources{
						CPUCount:  instance.CPUCount,
						MemoryMiB: instance.MemoryMiB,
						GPU: &instances.GPUResources{
							Kind:  gpu,
							Count: count,
						},
					},
				})
			}
		}
	}
	return result
}

func extendedInstanceTypeName(name string, gpu *instances.GPUResources) string {
	if gpu == nil || !strings.HasPrefix(name, "n1-") {
		return name
	}
	// Only augment name if n1- instance and GPU attached
	return fmt.Sprintf("%s-%s-%d", name, gpu.Kind, gpu.Count)
}

//-------------------------------------------------------------------------------------------------
// GPU LIMITS (https://cloud.google.com/compute/docs/gpus#nvidia_gpus_for_compute_workloads)
//-------------------------------------------------------------------------------------------------

func maxCpuAndMemoryForGpu(
	gpuKind typedefs.GPUKind, gpuCount uint16, zone string,
) (uint16, uint32) {
	switch gpuKind {
	case typedefs.GPUNvidiaTeslaT4:
		if gpuCount <= 2 {
			return 48, 312
		}
		return 96, 624
	case typedefs.GPUNvidiaTeslaP4:
		return 24 * gpuCount, 156 * uint32(gpuCount)
	case typedefs.GPUNvidiaTeslaV100:
		return 12 * gpuCount, 78 * uint32(gpuCount)
	case typedefs.GPUNvidiaTeslaP100:
		if gpuCount <= 2 {
			return 16 * gpuCount, 104 * uint32(gpuCount)
		}
		switch zone {
		case "us-east1-c", "europe-west1-b", "europe-west1-d":
			return 64, 208
		default:
			return 96, 624
		}
	case typedefs.GPUNvidiaTeslaK80:
		switch zone {
		case "asia-east1-a", "us-east1-d":
			return 8 * gpuCount, 52 * uint32(gpuCount)
		default:
			return 8 * gpuCount, jack.Min(208, 52*uint32(gpuCount))
		}
	default:
		return 0, 0
	}
}
