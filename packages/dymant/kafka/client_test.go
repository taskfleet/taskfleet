package kafka

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPublishSubscribeSync(t *testing.T) {
	fixture := newPubsubFixture(t)

	publisher := fixture.publisher()
	group := uuid.NewString()
	subscriber := fixture.subscriber(group)

	n := 5
	publishCount := publisher.publishN(n, true)
	subscribeCount := subscriber.subscribeN(n, 0)

	fixture.await()
	assert.Equal(t, <-publishCount, <-subscribeCount)
}

func TestPublishSubscribeAsync(t *testing.T) {
	fixture := newPubsubFixture(t)

	publisher := fixture.publisher()
	group := uuid.NewString()
	subscriber := fixture.subscriber(group)

	n := 50
	publishCount := publisher.publishN(n, false)
	subscribeCount := subscriber.subscribeN(n, 0)

	fixture.await()
	assert.Equal(t, <-publishCount, <-subscribeCount)
}
