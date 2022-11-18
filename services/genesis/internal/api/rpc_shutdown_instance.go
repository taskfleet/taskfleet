package api

import (
	"context"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/google/uuid"
	genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	"go.taskfleet.io/services/genesis/internal/db"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ShutdownInstance implements the genesis interface.
func (s *Service) ShutdownInstance(
	ctx context.Context, request *genesis_v1.ShutdownInstanceRequest,
) (*genesis_v1.ShutdownInstanceResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// We can't get an error since the request is validated to contain a UUID
	id, _ := uuid.Parse(request.Instance.Id)
	logger := zeus.Logger(ctx).With(zap.String("instance", id.String()))

	// First, we delete the instance in our database
	instance, err := s.database.GetInstance(ctx, id)
	if err != nil {
		logger.Error("failed to find instance to delete", zap.Error(err))
		if err == db.ErrNotExist {
			return nil, status.Error(codes.NotFound, "instance does not exist or has been deleted")
		}
		return nil, status.Error(codes.Unknown, "instance could not be found")
	}

	// Then, we can set it to be deleted
	if err := instance.SetDeleted(ctx); err != nil {
		logger.Error("failed to flag instance as deleted", zap.Error(err))
		return nil, status.Error(codes.Unknown, "failed to flag the instance as deleted")
	}

	// At this point, we can assure the caller that the instance is (or will be) deleted. We
	// asynchronously schedule the deletion on Gcloud
	ctx = zeus.WithLogger(context.Background(), zeus.Logger(ctx))
	go s.awaitShutdown(ctx, instance, shutdownReasonRequested)

	return &genesis_v1.ShutdownInstanceResponse{}, nil
}
