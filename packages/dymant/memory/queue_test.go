package memory

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/dymant"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestPublisher(t *testing.T) {
	queue := NewQueue(10)
	defer queue.Close()

	for i := 0; i < 5; i++ {
		err := queue.Publish(dymant.NoKey, timestamppb.Now())
		require.Nil(t, err)
	}
	assert.Len(t, queue.GetMessages(), 5)

	for i := 0; i < 10; i++ {
		err := queue.Publish(dymant.NoKey, timestamppb.Now())
		require.Nil(t, err)
	}
	assert.Len(t, queue.GetMessages(), 10)
}

func TestSubscriber(t *testing.T) {
	queue := NewQueue(5)
	defer queue.Close()

	var messages []proto.Message
	for i := 0; i < 5; i++ {
		messages = append(messages, timestamppb.Now())
	}
	queue.SetMessages(messages)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond)
	defer cancel()

	count := 0
	err := queue.Process(ctx, func(ctx context.Context, messages []proto.Message) error {
		count += len(messages)
		return nil
	})
	assert.Equal(t, 5, count)
	if !dymant.IsErrContext(err) {
		t.Error(err)
	}
}

func TestPublishSubscribe(t *testing.T) {
	queue := NewQueue(10)
	defer queue.Close()

	var messages []proto.Message
	for i := 0; i < 5; i++ {
		messages = append(messages, timestamppb.Now())
	}
	queue.SetMessages(messages)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Millisecond)
	defer cancel()

	messageCount := 0

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return queue.Process(ctx, func(ctx context.Context, messages []proto.Message) error {
			messageCount += len(messages)
			return nil
		})
	})
	eg.Go(func() error {
		for i := 0; i < 5; i++ {
			err := queue.Publish(dymant.NoKey, timestamppb.Now())
			require.Nil(t, err)
		}
		return nil
	})
	if err := eg.Wait(); !dymant.IsErrContext(err) {
		t.Error(err)
	}

	assert.Equal(t, 10, messageCount)
}
