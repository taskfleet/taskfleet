package api

import (
	"context"

	v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	"go.taskfleet.io/packages/dymant"
	"go.taskfleet.io/packages/mercury"
	"go.taskfleet.io/services/genesis/internal/db"
	"go.taskfleet.io/services/genesis/internal/ping"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
)

// Service implements the genesis API.
type Service struct {
	v1.UnimplementedInstanceManagerServiceServer

	ctx      context.Context
	health   mercury.Health
	database db.Connection
	kafka    dymant.Publisher
	provider providers.Provider
	store    template.Store
	grpc     *ping.Grpc
}

// Options wraps the services that are required by the genesis service.
type Options struct {
	Health         mercury.Health
	Database       db.Connection
	KafkaPublisher dymant.Publisher
	GRPC           *ping.Grpc
	Provider       providers.Provider
	Store          template.Store
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
		store:    options.Store,
		grpc:     options.GRPC,
	}
}
