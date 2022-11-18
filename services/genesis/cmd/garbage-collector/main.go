package main

import (
	"context"
	"time"

	"github.com/borchero/zeus/pkg/zeus"
	"github.com/kelseyhightower/envconfig"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/packages/postgres"
	"go.taskfleet.io/services/genesis/internal/db"
	"go.taskfleet.io/services/genesis/internal/gc"
	"go.taskfleet.io/services/genesis/internal/providers/impl/gcp"
	"go.uber.org/zap"
)

type environment struct {
	GCP gcp.ClientOptions
	DB  postgres.ConnectionConfig
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	// Setup
	var env environment
	envconfig.MustProcess("", &env)

	dbConnection := postgres.MustNewConnection(env.DB)
	database := db.NewConnection(dbConnection)
	gcpClient := jack.Must(gcp.NewClient(ctx, env.GCP))

	// Run garbage collection
	collector := gc.NewGarbageCollector(gc.Options{
		Database: database,
		Client:   gcpClient,
	})
	if err := collector.Collect(ctx); err != nil {
		zeus.Logger(ctx).Fatal("failed to run garbage collection", zap.Error(err))
	}
	zeus.Logger(ctx).Info("successfully ran garbage collection")
}
