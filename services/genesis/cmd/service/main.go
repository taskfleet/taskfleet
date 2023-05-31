package main

import (
	"context"

	"github.com/borchero/zeus/pkg/zeus"

	"go.taskfleet.io/packages/eagle"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/packages/mercury"
	"go.taskfleet.io/packages/postgres"
	"go.taskfleet.io/services/genesis/cmd"
	"go.taskfleet.io/services/genesis/cmd/service/setup"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.uber.org/zap"
)

type config struct {
	Kafka    cmd.KafkaConfig         `json:"kafka"`
	Postgres postgres.Config         `json:"postgres"`
	Minion   cmd.MinionConfig        `json:"minion"`
	Template template.InstanceConfig `json:"template"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Load configuration
	var config config
	if err := eagle.LoadConfig(&config,
		eagle.WithYAMLFile("/etc/instance-manager/config.yaml", true),
		eagle.WithEnvironment(""),
	); err != nil {
		zeus.Logger(ctx).Fatal("failed to load configuration", zap.Error(err))
	}

	// Setup servers
	zeus.Logger(ctx).Info("running initialization")
	grpc := jack.Must(mercury.NewGrpc(5404,
		mercury.WithPrometheusMetrics(),
		mercury.WithRequestValidation(),
		mercury.WithHealthService(),
		mercury.WithLogger(zeus.Logger(zeus.WithName(ctx, "grpc")), true),
	))
	prometheusInstance := mercury.NewPrometheus(9090)

	// Setup service
	service := setup.InitService(ctx, grpc.Health())
	genesis.RegisterInstanceManagerServiceServer(grpc.Server, service)
	zeus.Logger(ctx).Info("finished initialization")

	// Run servers
	runtime := mercury.NewRuntime(ctx).
		Schedule("grpc", grpc).
		Schedule("prometheus", prometheusInstance)
	if err := runtime.Await(); err != nil {
		zeus.Logger(ctx).Fatal("server encountered failure", zap.Error(err))
	}
}
