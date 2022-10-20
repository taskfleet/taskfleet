package kafka

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMultipleSubscribers(t *testing.T) {
	fixture := newPubsubFixture(t)
	defer fixture.teardown()

	// Round 1: only one subscriber
	publisher := fixture.publisher()
	group := uuid.NewString()
	subscriber := fixture.subscriber(group, WithBatchConfig(100, 100*time.Millisecond))

	n := 500 // Need to publish many messages to ensure that subscribers use consumer group
	publisher.publishN(5*n, false)
	subscribeCount := subscriber.subscribeN(n) // only read n of 5 * n

	fixture.await()
	assert.Equal(t, n, <-subscribeCount)

	// Round 2: two subscribers
	subscriber1 := fixture.subscriber(group, WithBatchConfig(100, 100*time.Millisecond))
	subscriber2 := fixture.subscriber(group, WithBatchConfig(100, 100*time.Millisecond))

	subscribe1Count := subscriber1.subscribeN(2 * n)
	subscribe2Count := subscriber2.subscribeN(2 * n)

	s1 := <-subscribe1Count
	s2 := <-subscribe2Count

	fixture.await()
	assert.Equal(t, 4*n, s1+s2)
	assert.Greater(t, s1, 0)
	assert.Greater(t, s2, 0)
}
