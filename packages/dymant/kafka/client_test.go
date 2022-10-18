package kafka

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/dymant"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var testClient = func() *Client {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	clientID := uuid.NewString()
	client, err := NewClient(clientID, []string{os.Getenv("KAFKA_BOOTSTRAP_SERVER")}, logger)
	if err != nil {
		panic(err)
	}
	return client
}()

func TestPublishSubscribe(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Create topic
	admin, err := testClient.admin()
	require.Nil(t, err)
	topic, err := admin.ephemeralTopic(ctx)
	require.Nil(t, err)
	// require.Nil(t, topic.awaitReady(ctx))
	// t.FailNow()
	defer topic.delete(ctx) // nolint:errcheck

	// Get publisher and subscriber. The subscriber uses a random group ID as otherwise, Kafka
	// might still wait for an old subscriber to come online for local tests that are repeated
	// often.
	publisher, err := testClient.Publisher(topic.name)
	require.Nil(t, err)

	group := uuid.NewString()
	subscriber, err := testClient.Subscriber(topic.name, group, &timestamppb.Timestamp{})
	require.Nil(t, err)

	publishedCount := make(chan int, 1)
	consumedCount := make(chan int, 1)

	// Produce and consume messages
	n := 5
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go publishMessages(ctx, t, publisher, n, wg, publishedCount)
	go subscribeMessages(ctx, t, subscriber, n, wg, consumedCount)
	wg.Wait()

	assert.Equal(t, <-publishedCount, <-consumedCount)
}

//-------------------------------------------------------------------------------------------------

func publishMessages(
	ctx context.Context, t *testing.T, publisher dymant.Publisher, n int,
	wg *sync.WaitGroup, result chan<- int,
) {
	defer wg.Done()

	key, err := uuid.NewRandom()
	if err != nil {
		t.Errorf("failed generating UUID: %s", err)
		result <- 0
		return
	}

	count := 0
	for i := 0; i < n; i++ {
		timestamp := timestamppb.Now()
		if err := publisher.PublishSync(ctx, key, timestamp); err != nil {
			t.Errorf("failed publishing message: %s", err)
			result <- count
			return
		}
		count++
	}

	result <- count
}

func subscribeMessages(
	ctx context.Context, t *testing.T, subscriber dymant.Subscriber, n int,
	wg *sync.WaitGroup, result chan<- int,
) {
	defer wg.Done()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	count := 0
	err := subscriber.Process(ctx, func(ctx context.Context, messages []proto.Message) error {
		count += len(messages)
		if count == n {
			cancel()
		}
		return nil
	})
	if !dymant.IsErrContext(err) {
		t.Errorf("failed to process messages: %s", err)
	}
	result <- count
}
