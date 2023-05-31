package api

import (
	"context"
	"fmt"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/docker/go-units"
	"github.com/google/uuid"
	v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	db "go.taskfleet.io/services/genesis/db/gen"
	gcpinstances "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/instances"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateInstance implements the genesis interface.
func (s *Service) CreateInstance(
	ctx context.Context, request *v1.CreateInstanceRequest,
) (*v1.CreateInstanceResponse, error) {
	// In this method, we assume that validation has already been performed as specified by the
	// gRPC configuration.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	logger := zeus.Logger(ctx).With(zap.String("request", "CreateInstance"))

	// As a first step, we read the instance template for the component for which the instance
	// is created
	template, err := s.store.Get(ctx, request.Component, template.InstanceDetails{
		GPUKind: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read instance template from Kubernetes: %s", err)
	}

	// At this point, we also take note of the instance's reservations
	var memoryReservation int64
	if template.Reservations.Memory != "" {
		memoryReservation, err = units.RAMInBytes(template.Reservations.Memory)
		if err != nil {
			return nil, fmt.Errorf("invalid memory reservation in template: %s", err)
		}
	}
	memoryReservationMb := memoryReservation / (1024 * 1024)

	// Afterwards, we obtain a suitable machine configuration for the request...
	gpuResources := instances.GPUResourcesUnmarshalProto(request.Resources.Gpu)
	instanceType, err := s.provider.Instances().Find(request.Config.Zone, instances.Resources{
		CPUCount:  uint16(request.Resources.CpuCount),
		MemoryMiB: request.Resources.Memory + uint32(memoryReservationMb),
		GPU:       gpuResources,
	})
	if err != nil {
		logger.Error("failed to find machine type for request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "no suitable machine type could be found")
	}

	// ...as well as the disk configuration
	disks, diskSizes, err := parseDisks(template.AdditionalDisks, request.Resources.CpuCount)
	if err != nil {
		return nil, fmt.Errorf("failed to parse additional disks: %s", err)
	}

	// Finally, we can send the request to the Google Cloud
	id, _ := uuid.Parse(request.Id)
	ref := providers.InstanceRef{ID: id, Zone: request.Config.Zone}
	spec := providers.InstanceSpec{
		Compute: providers.ComputeConfig{
			InstanceType: instanceType,
			IsSpot:       request.Config.IsSpot,
		},
		Boot: providers.BootConfig{
			ImageLink:   template.GCP.BootImage.Link,
			DiskSizeGiB: uint32(template.GCP.BootImage.SizeGiB),
		},
		Metadata: providers.MetadataConfig{
			Tags:       template.GCP.NetworkTags,
			Labels:     map[string]string{gcpinstances.LabelKeyOwnedBy: request.Owner},
			Attributes: template.Metadata,
		},
		Security: providers.SecurityConfig{
			ServiceAccountEmail: template.GCP.ServiceAccountEmail,
		},
		Disks: disks,
	}

	promise, err := s.provider.Instances().Create(ctx, ref, spec)
	if err != nil {
		logger.Error("failed to create instance", zap.Error(err))
		switch err.(type) {
		case providers.ClientError:
			return nil, status.Error(
				codes.InvalidArgument, "creating instance failed due to invalid arguments",
			)
		default:
			return nil, status.Error(
				codes.Unknown, "creating instance failed due to failed API request",
			)
		}
	}

	// As soon as we know that the initial request was successful, we can log the instance into our
	// database
	var gpuKind *typedefs.GPUKind
	var gpuCount *int32
	if gpuResources != nil {
		gpuKind = &gpuResources.Kind
		intCount := int32(gpuResources.Count)
		gpuCount = &intCount
	}

	newInstance := db.Instance{
		ID:                     id,
		Provider:               typedefs.CloudProviderUnmarshalProto(request.Config.CloudProvider),
		Zone:                   request.Config.Zone,
		Owner:                  request.Owner,
		MachineType:            instanceType.Name,
		IsSpot:                 request.Config.IsSpot,
		CPUCount:               int32(instanceType.CPUCount),
		CPUCountRequested:      int32(request.Resources.CpuCount),
		MemoryMB:               int32(instanceType.MemoryMiB),
		MemoryMBRequested:      int32(request.Resources.Memory),
		MemoryMBReserved:       int32(memoryReservationMb),
		GPUKind:                gpuKind,
		GPUCount:               gpuCount,
		BootImage:              template.GCP.BootImage.Link,
		BootDiskSizeGiB:        int32(template.GCP.BootImage.SizeGiB),
		DiskSizeHDDGiB:         int32(diskSizes.hdd),
		DiskSizeSSDStandardGiB: int32(diskSizes.ssdStandard),
		DiskSizeSSDHPGiB:       int32(diskSizes.ssdHighPerformance),
	}

	instance, err := s.database.CreateInstance(ctx, newInstance)
	if err != nil {
		logger.Error("failed to write instance to database", zap.Error(err))
		return nil, status.Error(
			codes.Unknown, "failed to write instance to database",
		)
	}

	// All else will be performed in a long-running background task ...
	go s.awaitCreation(promise, instance, uint32(memoryReservationMb))

	// Then, we only need to build the response object
	response := &v1.CreateInstanceResponse{
		Instance: &v1.Instance{Id: id.String()},
		Config: &v1.InstanceConfig{
			IsSpot:        request.Config.IsSpot,
			CloudProvider: request.Config.CloudProvider,
			Zone:          request.Config.Zone,
		},
		Resources: &v1.InstanceResources{
			Memory:   instanceType.MemoryMiB,
			CpuCount: uint32(instanceType.CPUCount),
			Gpu:      instanceType.GPU.MarshalProto(),
		},
	}
	return response, nil
}
