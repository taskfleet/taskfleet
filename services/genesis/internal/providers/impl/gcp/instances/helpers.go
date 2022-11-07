package gcpinstances

import (
	"context"
	"fmt"
	"path"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/docker/go-units"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

//-------------------------------------------------------------------------------------------------
// DISKS
//-------------------------------------------------------------------------------------------------

type disksHelper struct {
	bootDiskSizeGiB   int64
	bootImages        []template.Option[string]
	extraDisks        []disk
	diskTypeSelfLinks map[string]string
}

type disk struct {
	name          string
	sizePerCpuGiB int64
}

func newDisksHelper(
	ctx context.Context,
	projectID string,
	bootConfig template.GcpBootConfig,
	extraDisksConfig []template.InstanceDisk,
	diskType string,
	diskClient *compute.DisksClient,
) (*disksHelper, error) {
	// Parse disk sizes
	bootDiskSize, err := units.RAMInBytes(bootConfig.DiskSize)
	if err != nil {
		return nil, fmt.Errorf("failed to parse boot disk size %q: %s", bootConfig.DiskSize, err)
	}

	extraDisks := []disk{}
	for i, config := range extraDisksConfig {
		diskSize, err := units.RAMInBytes(config.SizePerCPU)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to parse size of extra disk %d %q: %s", i, config.SizePerCPU, err,
			)
		}
		extraDisks = append(extraDisks, disk{
			name:          config.Name,
			sizePerCpuGiB: int64(diskSize / units.GiB),
		})
	}

	// Get links to requested disk type
	it := diskClient.AggregatedList(ctx, &computepb.AggregatedListDisksRequest{
		Project: projectID,
		Filter:  proto.String(fmt.Sprintf("name=\"%s\"", diskType)),
	})
	diskTypes := map[string]string{}
	if err := gcputils.Iterate[compute.DisksScopedListPair](
		ctx, it, func(pair compute.DisksScopedListPair) error {
			region := path.Base(pair.Key)
			if len(pair.Value.GetDisks()) != 1 {
				return fmt.Errorf(
					"unexpectedly found %d instead of one disk type %q in region %s",
					len(pair.Value.Disks), diskType, region,
				)
			}
			diskTypes[region] = pair.Value.GetDisks()[0].GetSelfLink()
			return nil
		},
	); err != nil {
		return nil, fmt.Errorf("failed to find disk type %q: %s", diskType, err)
	}

	return &disksHelper{
		bootDiskSizeGiB:   int64(bootDiskSize / units.GiB),
		bootImages:        bootConfig.ImageLink,
		extraDisks:        extraDisks,
		diskTypeSelfLinks: diskTypes,
	}, nil
}

func (h *disksHelper) diskConfig(
	instanceID string,
	zone string,
	resources instances.Resources,
	architecture typedefs.CPUArchitecture,
) []*computepb.AttachedDisk {
	region := gcputils.RegionFromZone(zone)

	// First, we need to get the correct source image
	bootImage := template.MatchingOption(h.bootImages, resources.GPUKind(), architecture)

	// Then, we obtain boot disk
	disks := []*computepb.AttachedDisk{{
		AutoDelete: proto.Bool(true),
		Boot:       proto.Bool(true),
		Mode:       proto.String("READ_WRITE"),
		InitializeParams: &computepb.AttachedDiskInitializeParams{
			DiskName:    proto.String(fmt.Sprintf("%s-boot-disk", instanceID)),
			DiskSizeGb:  proto.Int64(h.bootDiskSizeGiB),
			DiskType:    proto.String(h.diskTypeSelfLinks[region]),
			SourceImage: bootImage,
		},
	}}

	// Optionally, we add additional disks
	for i, disk := range h.extraDisks {
		disks = append(disks, &computepb.AttachedDisk{
			AutoDelete: proto.Bool(true),
			Boot:       proto.Bool(false),
			DeviceName: proto.String(disk.name),
			Mode:       proto.String("READ_WRITE"),
			InitializeParams: &computepb.AttachedDiskInitializeParams{
				DiskName:   proto.String(fmt.Sprintf("%s-extra-disk-%d", instanceID, i)),
				DiskSizeGb: proto.Int64(disk.sizePerCpuGiB * int64(resources.CPUCount)),
				DiskType:   proto.String(h.diskTypeSelfLinks[region]),
			},
		})
	}

	return disks
}

//-------------------------------------------------------------------------------------------------
// RESERVATIONS
//-------------------------------------------------------------------------------------------------

type reservationsHelper struct {
	memoryMiB uint32
}

func newReservationsHelper(
	reservations template.InstanceReservations,
) (*reservationsHelper, error) {
	if reservations.Memory == nil {
		return &reservationsHelper{memoryMiB: 0}, nil
	}
	memoryReservation, err := units.RAMInBytes(*reservations.Memory)
	if err != nil {
		return nil, fmt.Errorf("invalid memory reservation %q", *reservations.Memory)
	}
	return &reservationsHelper{memoryMiB: uint32(memoryReservation / units.MiB)}, nil
}

func (h *reservationsHelper) updateResources(resources instances.Resources) instances.Resources {
	resources.MemoryMiB += h.memoryMiB
	return resources
}
