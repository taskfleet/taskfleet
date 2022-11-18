package api

import (
	"context"

	"github.com/borchero/zeus/pkg/zeus"
	genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	"go.taskfleet.io/services/genesis/internal/db"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ListInstances implements the genesis interface.
func (s *Service) ListInstances(
	ctx context.Context, request *genesis_v1.ListInstancesRequest,
) (*genesis_v1.ListInstancesResponse, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Get logger
	logger := zeus.Logger(ctx).With(zap.String("owner", request.Owner))

	// List all running instances owned by the
	iterator := s.database.ListInstances(
		ctx, db.FilterOwner(request.Owner), db.InstanceStatusRunning,
	)
	instances := []*genesis_v1.RunningInstance{}
	for item := range iterator.Iter() {
		if item.Err != nil {
			logger.Error("failed to list instances", zap.Error(item.Err))
			return nil, status.Error(codes.Unknown, "failed to list owned instances")
		}

		// Create config
		config := &genesis_v1.InstanceResources{
			Memory:   uint32(item.Instance.MemoryMB),
			CpuCount: uint32(item.Instance.CPUCount),
		}
		if item.Instance.GPUKind != nil && item.Instance.GPUCount != nil {
			config.Gpu = &genesis_v1.GPUResources{
				Kind:  item.Instance.GPUKind.MarshalProto(),
				Count: uint32(*item.Instance.GPUCount),
			}
		}

		instances = append(instances, &genesis_v1.RunningInstance{
			Instance: &genesis_v1.Instance{Id: item.Instance.ID.String()},
			Config: &genesis_v1.InstanceConfig{
				IsSpot:        item.Instance.IsSpot,
				CloudProvider: item.Instance.Provider.MarshalProto(),
				Zone:          item.Instance.Zone,
			},
			Resources: config,
			Hostname:  *item.Instance.Hostname,
		})
	}

	return &genesis_v1.ListInstancesResponse{Instances: instances}, nil
}
