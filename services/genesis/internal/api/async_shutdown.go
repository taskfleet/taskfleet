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

type shutdownReason string

const (
	shutdownReasonCreationFailed shutdownReason = "creation-failed"
	shutdownReasonRequested      shutdownReason = "requested"
)

func (s *Service) awaitShutdown(
	ctx context.Context, instance db.Instance, reason shutdownReason,
) {
	logger := zeus.Logger(s.ctx)

	// Allow 10 minutes for deletion. Otherwise, it should be purged later, there is surely
	// something wrong
	ctx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	// Then, we set it deleted if the shutdown reason is creation failure
	if reason == shutdownReasonCreationFailed {
		logger.Debug("flagging instance as deleted")
		if err := instance.SetDeleted(ctx); err != nil {
			zeus.Logger(s.ctx).Error("failed to flag instance as deleted", zap.Error(err))
			return
		}
	}

	logger.Debug("deleting instance")
	err := s.provider.Instances().Delete(
		ctx, providers.InstanceMeta{
			ID:           instance.ID,
			ProviderID:   instance.ProviderID,
			ProviderZone: instance.Zone,
		},
	)
	errNotFound := false
	if err != nil {
		if !providers.IsErrNotFound(err) {
			logger.Error("failed to delete instance", zap.Error(err))
			return
		}
		errNotFound = true
	}

	// Also, we can publish the deletion. We may have never published a deletion message or the
	// instance may not even have been found, but we should still publish this.
	now := time.Now()
	deletionMessage := &v1.InstanceEvent{
		Instance: &genesis_v1.Instance{Id: instance.ID.String()},
		Timestamp: &timestamppb.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}
	if reason == shutdownReasonCreationFailed {
		deletionMessage.Event = &v1.InstanceEvent_CreationFailed{
			CreationFailed: &v1.InstanceCreationFailedEvent{
				Reason: v1.InstanceCreationFailedEvent_REASON_UNSPECIFIED, // TODO: update
			},
		}
	} else {
		reason := v1.InstanceDeletedEvent_REASON_SHUTDOWN
		if errNotFound {
			reason = v1.InstanceDeletedEvent_REASON_TERMINATED
		}
		deletionMessage.Event = &v1.InstanceEvent_Deleted{
			Deleted: &v1.InstanceDeletedEvent{
				Reason: reason,
			},
		}
	}

	if err := s.kafka.PublishSync(ctx, instance.ID, deletionMessage); err != nil {
		logger.Warn("failed to publish instance deletion to Kafka", zap.Error(err))
	}
	logger.Info("successfully shut down instance")
}
