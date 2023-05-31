package api

import (
	"context"

	v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	"go.taskfleet.io/packages/dymant"
	"go.taskfleet.io/packages/mercury"
	db "go.taskfleet.io/services/genesis/db/gen"
	"go.taskfleet.io/services/genesis/internal/ping"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
)

// Service implements the genesis API.
type Service struct {
	v1.UnimplementedGenesisServiceServer

	ctx      context.Context
	health   mercury.Health
	database db.Querier
	kafka    dymant.Publisher
	provider providers.Provider
	grpc     *ping.Grpc
}

// Options wraps the services that are required by the genesis service.
type Options struct {
	Health         mercury.Health
	Database       db.Querier
	KafkaPublisher dymant.Publisher
	GRPC           *ping.Grpc
	Provider       providers.Provider
}

// NewService initializes a new service with the specified clients. The given context will be
// reused as parent context for long-running background tasks.
func NewService(ctx context.Context, options Options) *Service {
	return &Service{
		ctx:      ctx,
		health:   options.Health,
		database: options.Database,
		kafka:    options.KafkaPublisher,
		provider: options.Provider,
		grpc:     options.GRPC,
	}
}
