package api_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	genesis_messages "go.taskfleet.io/grpc/gen/go/genesis/messages/v1"
	"go.taskfleet.io/packages/dymant"
	"go.taskfleet.io/packages/dymant/kafka"
	"go.taskfleet.io/services/genesis/cmd/service/setup"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var errDone = errors.New("done")
var listener = bufconn.Listen(1024 * 1024)
var client = func() genesis.InstanceManagerServiceClient {
	ctx := context.Background()
	dialer := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	conn, err := grpc.DialContext(
		ctx, "bufconn",
		grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		panic(err)
	}
	return genesis.NewInstanceManagerServiceClient(conn)
}()

func init() {
	ctx := context.Background()

	grpc := setup.InitGrpcServer(ctx)
	service := setup.InitService(ctx, grpc.Health())
	genesis.RegisterInstanceManagerServiceServer(grpc.Server, service)

	go func() {
		if err := grpc.Server.Serve(listener); err != nil {
			panic(err)
		}
	}()
}

//-------------------------------------------------------------------------------------------------

func TestListAvailableGpus(t *testing.T) {
	ctx := context.Background()
	response, err := client.ListZones(ctx, &genesis.ListZonesRequest{})
	require.Nil(t, err)
	assert.GreaterOrEqual(t, len(response.Zones), 91)
}

//-------------------------------------------------------------------------------------------------

func TestCreateListShutdownInstance(t *testing.T) {
	ctx := context.Background()
	kafkaSubscriber := func() dymant.Subscriber {
		var config kafka.ClientOptions
		envconfig.MustProcess("KAFKA", &config)
		client := kafka.MustNewClient(config, zap.NewNop())
		return client.MustSubscriber(
			os.Getenv("KAFKA_TOPIC"), &genesis_messages.InstanceEvent{}, kafka.SubscriberOptions{},
		)
	}()

	// Create instance
	instance := testCreateInstance(ctx, t)

	// Wait for instance to be running, i.e. wait for a message from Kafka
	running := make(chan struct{})
	go func() {
		err := kafkaSubscriber.Process(ctx, func(ctx context.Context, batch dymant.Batch) error {
			event := batch.Messages()[0].(*genesis_messages.InstanceEvent)
			if event.Instance.Id == instance.Id && event.GetCreated() != nil {
				close(running)
				return errDone
			}
			return fmt.Errorf("received unknown created instance %s", event.Instance.Id)
		})
		if err != errDone {
			panic(err)
		}
	}()
	<-running

	// List the instances that are currently running
	testListInstances(ctx, t, instance)

	// Shutdown the instance
	testShutdownInstance(ctx, t, instance)

	// Wait for shutdown to succeed
	deleted := make(chan struct{})
	go func() {
		err := kafkaSubscriber.Process(ctx, func(ctx context.Context, batch dymant.Batch) error {
			event := batch.Messages()[0].(*genesis_messages.InstanceEvent)
			if event.Instance.Id == instance.Id && event.GetDeleted() != nil {
				close(deleted)
				return errDone
			}
			return fmt.Errorf("received unknown deleted instance %s", event.Instance.Id)
		})
		if err != errDone {
			panic(err)
		}
	}()
	<-deleted
}

func testCreateInstance(ctx context.Context, t *testing.T) *genesis.Instance {
	request := &genesis.CreateInstanceRequest{
		Id:    uuid.NewString(),
		Owner: "genesis-test",
		Config: &genesis.InstanceConfig{
			CloudProvider: genesis.CloudProvider_CLOUD_PROVIDER_GOOGLE_CLOUD_PLATFORM,
			Zone:          "us-east1-c",
		},
		Resources: &genesis.InstanceResources{
			Memory:   2500,
			CpuCount: 1,
		},
	}
	response, err := client.CreateInstance(ctx, request)
	require.Nil(t, err)

	instance := response.Instance
	assert.NotNil(t, instance)
	return instance
}

func testListInstances(ctx context.Context, t *testing.T, expected *genesis.Instance) {
	request := &genesis.ListInstancesRequest{Owner: "genesis-test"}
	instances, err := client.ListInstances(ctx, request)
	require.Nil(t, err)
	assert.Len(t, instances.Instances, 1)
	assert.Equal(t, instances.Instances[0].Instance.Id, expected.Id)
}

func testShutdownInstance(ctx context.Context, t *testing.T, instance *genesis.Instance) {
	request := &genesis.ShutdownInstanceRequest{Instance: instance}
	_, err := client.ShutdownInstance(ctx, request)
	require.Nil(t, err)
}
