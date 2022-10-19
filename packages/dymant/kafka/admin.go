package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

// AdminClient allows to perform administrative tasks on the Kafka cluster.
type AdminClient struct {
	client *kafka.AdminClient
}

// Admin returns an administrative client for the Kafka cluster.
func (c *Client) Admin() (*AdminClient, error) {
	config, err := c.config.config()
	if err != nil {
		return nil, err
	}
	admin, err := kafka.NewAdminClient(&config)
	if err != nil {
		return nil, err
	}
	return &AdminClient{admin}, nil
}

// EphemeralTopic creates a new topic with a random name. The topic is created with 3 partitions
// and a replication factor of 1.
func (c *AdminClient) EphemeralTopic(ctx context.Context) (*EphemeralTopic, error) {
	id := uuid.New()
	if _, err := c.client.CreateTopics(ctx, []kafka.TopicSpecification{
		{Topic: id.String(), NumPartitions: 3, ReplicationFactor: 1},
	}); err != nil {
		return nil, err
	}

	return &EphemeralTopic{client: c.client, name: id.String()}, nil
}
