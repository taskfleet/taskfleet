package kafka

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/dymant"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type pubsubFixture struct {
	ctx   context.Context
	t     *testing.T
	wg    *sync.WaitGroup
	topic *EphemeralTopic
}

type testPublisher struct {
	t         *testing.T
	ctx       context.Context
	wg        *sync.WaitGroup
	publisher dymant.Publisher
}

type testSubscriber struct {
	t          *testing.T
	ctx        context.Context
	wg         *sync.WaitGroup
	subscriber dymant.Subscriber
}

func newPubsubFixture(t *testing.T) *pubsubFixture {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	topic, err := adminClient.EphemeralTopic(ctx)
	require.Nil(t, err)
	t.Cleanup(func() {
		topic.Delete(ctx) // nolint:errcheck
	})

	return &pubsubFixture{
		ctx:   ctx,
		t:     t,
		wg:    &sync.WaitGroup{},
		topic: topic,
	}
}

func (f *pubsubFixture) publisher() *testPublisher {
	f.wg.Add(1)
	publisher, err := client.Publisher(f.topic.name)
	require.Nil(f.t, err)
	return &testPublisher{f.t, f.ctx, f.wg, publisher}
}

func (f *pubsubFixture) subscriber(group string, options ...SubscriberOption) *testSubscriber {
	f.wg.Add(1)
	subscriber, err := client.Subscriber(f.topic.name, group, &timestamppb.Timestamp{}, options...)
	require.Nil(f.t, err)
	return &testSubscriber{f.t, f.ctx, f.wg, subscriber}
}

func (f *pubsubFixture) await() {
	f.wg.Wait()
}

func (p *testPublisher) publishN(n int, sync bool) <-chan int {
	ch := make(chan int, 1)
	go func() {
		defer p.wg.Done()

		count := 0
		for i := 0; i < n; i++ {
			key := uuid.New()
			timestamp := timestamppb.Now()
			var err error
			if sync {
				err = p.publisher.PublishSync(p.ctx, key, timestamp)
			} else {
				err = p.publisher.Publish(key, timestamp)
			}
			if err != nil {
				p.t.Errorf("failed publishing message: %s", err)
				ch <- count
				return
			}
			count++
		}
		ch <- count
		p.publisher.Flush(p.ctx) // nolint:errcheck
	}()
	return ch
}

func (s *testSubscriber) subscribeN(n int, wait time.Duration) <-chan int {
	ch := make(chan int, 1)
	go func() {
		defer s.wg.Done()

		ctx, cancel := context.WithTimeout(s.ctx, 20*time.Second)
		defer cancel()

		time.Sleep(wait)

		count := 0
		err := s.subscriber.Process(ctx, func(ctx context.Context, messages []proto.Message) error {
			count += len(messages)
			if count == n {
				cancel()
			}
			return nil
		})
		if !dymant.IsErrContext(err) {
			s.t.Errorf("failed to process messages: %s", err)
		}
		ch <- count
		s.subscriber.Close() // nolint: errcheck
	}()
	return ch
}
