package kafka

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultipleSubscribers(t *testing.T) {
	fixture := newPubsubFixture(t)

	group := uuid.NewString()
	sub1 := fixture.subscriber(group, WithBatchConfig(100, 100*time.Millisecond))
	sub2 := fixture.subscriber(group, WithBatchConfig(100, 100*time.Millisecond))

	// Poll subscribers to ensure that both joined the consumer group
	for {
		ev := sub1.subscriber.(*subscriber).consumer.Poll(100)
		require.Nil(t, ev) // no event should be delivered here

		ev = sub2.subscriber.(*subscriber).consumer.Poll(100)
		require.Nil(t, ev) // no event should be delivered here

		assign1, err := sub1.subscriber.(*subscriber).consumer.Assignment()
		require.Nil(t, err)

		assign2, err := sub2.subscriber.(*subscriber).consumer.Assignment()
		require.Nil(t, err)

		if len(assign1) > 0 && len(assign2) > 0 {
			break
		}
	}

	// Publish messages only now
	publisher := fixture.publisher()
	n := 5
	publisher.publishN(2*n, false)

	// Read the messages (don't know how many arrive at each subscriber) -- they must arrive
	// immediately
	subscribe1Count := sub1.subscribeN(-1, time.Second)
	subscribe2Count := sub2.subscribeN(-1, time.Second)

	// Check if everything went as expectd
	s1 := <-subscribe1Count
	s2 := <-subscribe2Count

	fixture.await()
	assert.Equal(t, 2*n, s1+s2)
	assert.Greater(t, s1, 0)
	assert.Greater(t, s2, 0)
}
