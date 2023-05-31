package api

import (
	"context"
	"time"

	"github.com/borchero/zeus/pkg/zeus"
	v1 "go.taskfleet.io/grpc/gen/go/genesis/messages/v1"
	genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	db "go.taskfleet.io/services/genesis/db/gen"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Service) awaitCreation(
	promise providers.InstancePromise, dbInstance *db.Instance, memoryReservationMb uint32,
) {
	logger := zeus.Logger(s.ctx).With(zap.String("instance", dbInstance.ID.String()))

	// We allow this operation to take 15 minutes in total. In any other case, something must have
	// gone horribly wrong.
	ctx, cancel := context.WithTimeout(s.ctx, 15*time.Minute)
	defer cancel()

	// At the beginning, we first wait for the Google Cloud to indicate that the instance is up
	logger.Debug("waiting for instance to be created")
	instance, err := promise.Await(ctx)
	if err != nil {
		logger.Error("failed to wait for instance to be booting: initiating shutdown",
			zap.Error(err),
		)
		s.awaitShutdown(ctx, dbInstance, shutdownReasonCreationFailed)
		return
	}

	// At that point, we know that the instance is booting and we can store that info in our DB
	logger.Debug("instance was created successfully")
	if err := dbInstance.SetBooting(
		ctx, instance.Status.Compute.CPUKind, instance.Status.Network.InternalHostname,
	); err != nil {
		// We basically ignore the error
		logger.Warn("failed to write state 'booting' to database", zap.Error(err))
		return
	}

	// Afterwards, we can wait for the instance to actually be running. We assume that a gRPC
	// server will be reachable on port 5404. Since this instance is outside of Kubernetes, we
	// need to set TLS credentials here.
	logger.Debug("waiting for instance to be running")
	if err := s.grpc.AwaitHealthy(
		ctx, instance.Status.Network.InternalHostname, 5404,
	); err != nil {
		logger.Error("failed to wait for instance to be up and running: initiating shutdown",
			zap.Error(err),
		)
		s.awaitShutdown(ctx, dbInstance, shutdownReasonCreationFailed)
		return
	}

	// Now we can finally set our instance to running state
	logger.Debug("instance started up successfully")
	if err := dbInstance.SetRunning(ctx); err != nil {
		// We basically ignore the error
		logger.Warn("failed to write state 'running' to database", zap.Error(err))
		return
	}

	// Eventually, we're done so we can push the instance to Kafka
	now := time.Now()
	creationMessage := &v1.InstanceEvent{
		Instance: &genesis_v1.Instance{
			Id: instance.Ref.ID.String(),
		},
		Timestamp: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
		Event: &v1.InstanceEvent_Created{
			Created: &v1.InstanceCreatedEvent{
				Config: &genesis_v1.InstanceConfig{
					IsSpot:        instance.Spec.Compute.IsSpot,
					CloudProvider: genesis_v1.CloudProvider_CLOUD_PROVIDER_GOOGLE_CLOUD_PLATFORM,
					Zone:          instance.Ref.Zone,
				},
				Resources: &genesis_v1.InstanceResources{
					Memory:   uint32(instance.Spec.Compute.InstanceType.MemoryMiB) - memoryReservationMb,
					CpuCount: uint32(instance.Spec.Compute.InstanceType.CPUCount),
					Gpu:      instance.Spec.Compute.InstanceType.GPU.MarshalProto(),
				},
				Hostname: instance.Status.Network.InternalHostname,
			},
		},
	}

	if err := s.kafka.PublishSync(ctx, dbInstance.ID, creationMessage); err != nil {
		// Basically ignore the error - instances ought to be synchronized with the scheduler
		// periodically
		logger.Warn("failed to publish instance creation to Kafka", zap.Error(err))
	}

	logger.Info("successfully launched instance")
}
