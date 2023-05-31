package setup

import (
	"context"
	"time"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	"go.taskfleet.io/packages/dymant"
	"go.taskfleet.io/packages/dymant/kafka"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/packages/mercury"
	db "go.taskfleet.io/services/genesis/db/gen"
	"go.taskfleet.io/services/genesis/internal/api"
	"go.taskfleet.io/services/genesis/internal/ping"
	"go.taskfleet.io/services/genesis/internal/providers/impl/gcp"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.uber.org/zap"
)

// InitService initializes a new Genesis service from environment variables.
func InitService(
	ctx context.Context, health mercury.Health,
) v1.GenesisServiceServer {
	database := func() db.Querier {
		pool := jack.Must(pgxpool.New(ctx, ""))
		return db.New(pool)
	}()
	kafkaPublisher := func() dymant.Publisher {
		client := jack.Must(kafka.NewClient(env.Kafka, zeus.Logger(zeus.WithName(ctx, "kafka"))))
		return client.MustPublisher(env.KafkaTopic, kafka.PublisherOptions{})
	}()
	gcpClient := func() providers.Provider {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		return jack.Must(gcp.NewClient(ctx, env.GCP))
	}()
	store := func() template.Store {
		return template.MustNewStoreFromConfig(env.Template)
	}()
	grpc := func() *ping.Grpc {
		pinger, err := ping.NewGrpc(env.Minion)
		if err != nil {
			zeus.Logger(ctx).Fatal("failed to initialize gRPC pinger", zap.Error(err))
		}
		return pinger
	}()

	return api.NewService(ctx, api.Options{
		Health:         health,
		Database:       database,
		KafkaPublisher: kafkaPublisher,
		Provider:       gcpClient,
		Store:          store,
		GRPC:           grpc,
	})
}
