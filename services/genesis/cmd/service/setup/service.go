package setup

import (
	"context"
	"time"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/kelseyhightower/envconfig"
	"go.taskfleet.io/packages/dymant"
	"go.taskfleet.io/packages/dymant/kafka"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/packages/mercury"
	"go.taskfleet.io/packages/postgres"
	"go.taskfleet.io/services/genesis/internal/api"
	"go.taskfleet.io/services/genesis/internal/db"
	"go.taskfleet.io/services/genesis/internal/ping"
	"go.taskfleet.io/services/genesis/internal/providers/impl/gcp"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.uber.org/zap"
)

type environment struct {
	DB       postgres.ConnectionConfig
	Template template.StoreConfig
}

// InitService initializes a new Genesis service from environment variables.
func InitService(
	ctx context.Context, health mercury.Health,
) genesis.InstanceManagerServiceServer {
	var env environment
	envconfig.MustProcess("", &env)

	database := func() db.Connection {
		database := postgres.MustNewConnection(env.DB)
		return db.NewConnection(database)
	}()
	kafkaPublisher := func() dymant.Publisher {
		client := kafka.MustNewClient(env.Kafka, zeus.Logger(zeus.WithName(ctx, "kafka")))
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
