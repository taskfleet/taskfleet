package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// EphemeralTopic is a Kafka topic that is created with a random name and which easily allows for
// deletion. Typically, this topic is used for testing.
type EphemeralTopic struct {
	name   string
	client *kafka.AdminClient
}

// Delete deletes the topic in Kafka.
func (t *EphemeralTopic) Delete(ctx context.Context) error {
	_, err := t.client.DeleteTopics(ctx, []string{t.name})
	return err
}
