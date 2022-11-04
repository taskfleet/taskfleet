package gcpinstances

import (
	"context"
	"fmt"
	"path"
	"strings"

	"github.com/borchero/zeus/pkg/zeus"
	"go.taskfleet.io/packages/jack"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"go.uber.org/zap"
	"google.golang.org/api/compute/v1"
)

func findAvailableInstanceTypes(
	ctx context.Context, service *compute.Service, zones *gcpzones.Client, projectID string,
) (map[string][]instances.Type, error) {
	// Get instance types by zone
	result := map[string][]instances.Type{}
	for _, zone := range zones.List() {
		result[zone.Name] = make([]instances.Type, 0)
	}

	// First, we fetch the aggregated list of all machine types
	call := service.MachineTypes.AggregatedList(projectID)
	err := call.Pages(ctx, func(list *compute.MachineTypeAggregatedList) error {
		for zone, list := range list.Items {
			zone = path.Base(zone)

			// If the zone is not available, continue
			if _, ok := result[zone]; !ok {
				continue
			}

			// Then iterate over available machine types
			for _, item := range list.MachineTypes {
				if instance := tryUnmarshalInstanceType(ctx, item); instance != nil {
					result[zone] = append(result[zone], *instance)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all machine types: %s", err)
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
			accelerator, err := zones.FindAccelerator(zone.Name, gpu)
			if err != nil {
				continue // no possible configurations if GPU type not available
			}
			gpuInstances := findAvailableGpuInstanceTypes(
				zone.Name, result[zone.Name], gpu, accelerator.MaxCount(),
			)
			result[zone.Name] = append(result[zone.Name], gpuInstances...)
		}
	}

	// Now, we are done
	return result, nil
}

func tryUnmarshalInstanceType(ctx context.Context, item *compute.MachineType) *instances.Type {
	// When iterating over the returned machines, we exclude some sets of machines:
	// * Machines with shared cores
	// * M2 machines since they require purchasing committed usage plans
	if item.IsSharedCpu || strings.HasPrefix(item.Name, "m2-") {
		return nil
	}

	// Rename n1-ultramem and n1-megamem to m1-xxx, see
	// https://cloud.google.com/compute/docs/machine-types#m1_machine_types
	if strings.Contains(item.Name, "megamem") || strings.Contains(item.Name, "ultramem") {
		item.Name = "m1-" + item.Name[3:]
	}

	// Build the resources
	resources := instances.Resources{
		CPUCount:  uint16(item.GuestCpus),
		MemoryMiB: uint32(item.MemoryMb),
	}
	if len(item.Accelerators) > 0 {
		accelerator := item.Accelerators[0]
		gpuKind, err := typedefs.GPUKindFromProviderGcp(accelerator.GuestAcceleratorType)

		// Skip this instance if it cannot be parsed
		if err != nil {
			zeus.Logger(ctx).Warn(
				"skipping GCP machine type due to unknown accelerator",
				zap.String("machine_type", item.Name),
				zap.String("gpu", accelerator.GuestAcceleratorType),
			)
			return nil
		}

		resources.GPU = &instances.GPUResources{
			Kind: gpuKind, Count: uint16(accelerator.GuestAcceleratorCount),
		}
	}

	// Get the CPU architecture. Currently, the API does not allow to discern the CPU architecture.
	// Currently, only T2A instances provide ARM processors.
	architecture := typedefs.ArchitectureX86
	if strings.HasPrefix(item.Name, "t2a-") {
		architecture = typedefs.ArchitectureArm
	}

	// Return the instance
	return &instances.Type{Name: item.Name, Resources: resources, Architecture: architecture}
}

func findAvailableGpuInstanceTypes(
	zone string, n1Instances []instances.Type, gpu typedefs.GPUKind, maxGpuCount uint16,
) []instances.Type {
	result := make([]instances.Type, 0)
	for count := uint16(1); count <= maxGpuCount; count *= 2 {
		maxCPU, maxMem := maxCpuAndMemoryForGpu(gpu, count, zone)
		for _, instance := range n1Instances {
			if instance.CPUCount <= maxCPU && instance.MemoryMiB <= maxMem*1024 {
				result = append(result, instances.Type{
					Name: instance.Name,
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
