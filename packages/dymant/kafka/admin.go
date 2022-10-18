package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
)

type adminClient struct {
	client *kafka.AdminClient
}

type topic struct {
	client *kafka.AdminClient
	name   string
}

func (c *Client) admin() (*adminClient, error) {
	config, err := c.config.config()
	if err != nil {
		return nil, err
	}
	admin, err := kafka.NewAdminClient(&config)
	if err != nil {
		return nil, err
	}
	return &adminClient{admin}, nil
}

func (c *adminClient) ephemeralTopic(
	ctx context.Context,
) (*topic, error) {
	id := uuid.New()
	if _, err := c.client.CreateTopics(ctx, []kafka.TopicSpecification{
		{Topic: id.String(), NumPartitions: 3, ReplicationFactor: 1},
	}); err != nil {
		return nil, err
	}

	return &topic{
		client: c.client,
		name:   id.String(),
	}, nil
}

func (t *topic) delete(ctx context.Context) error {
	_, err := t.client.DeleteTopics(ctx, []string{t.name})
	return err
}
