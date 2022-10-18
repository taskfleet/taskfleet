package kafka

import (
	"fmt"

	"go.taskfleet.io/packages/dymant"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

// Client is the Kafka client that allows creating subscribers and publishers.
type Client struct {
	config clientConfig
	logger *zap.Logger
}

// NewClient creates a new client to communicate with Kafka. The client uses the specified ID to
// uniquely identify this client. The bootstrap servers define the broker address(es) of the Kafka
// cluster. Additional options may configure the basic configuration of the client such as
// authentication.
func NewClient(
	id string, bootstrapServers []string, logger *zap.Logger, options ...ClientOption,
) (*Client, error) {
	// Check inputs
	if id == "" {
		return nil, fmt.Errorf("client ID must be provided")
	}
	if len(bootstrapServers) == 0 {
		return nil, fmt.Errorf("at least one bootstrap server must be provided")
	}
	if logger == nil {
		logger = zap.NewNop()
	}

	// Initialize client
	config := clientConfig{
		id:               id,
		bootstrapServers: bootstrapServers,
		options:          options,
	}
	return &Client{config, logger}, nil
}

//-------------------------------------------------------------------------------------------------
// PUBLISHER/SUBSCRIBER
//-------------------------------------------------------------------------------------------------

// Publisher returns a new producer for the given topic, adhering to the supplied configuration.
func (c *Client) Publisher(topic string, options ...PublisherOption) (dymant.Publisher, error) {
	if topic == "" {
		return nil, fmt.Errorf("cannot publish to empty topic")
	}
	config, err := c.config.producerConfig(options)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %s", err)
	}
	return newPublisher(topic, config, c.logger.With(
		zap.String(logKeyTopic, topic),
		zap.String(logKeyComponent, "publisher"),
	))
}

// Subscriber returns a new consumer group for the given topic, expecting messages with the given
// type. Batch options and delivery guarantees are configured using the given options.
func (c *Client) Subscriber(
	topic, group string, message proto.Message, options ...SubscriberOption,
) (dymant.Subscriber, error) {
	if topic == "" {
		return nil, fmt.Errorf("cannot subscribe to empty topic")
	}
	config, err := c.config.consumerConfig(group, options)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %s", err)
	}
	subConfig := subscriberConfig{}
	for _, option := range options {
		option.configApply(&subConfig)
	}
	return newSubscriber(topic, config, subConfig, message, c.logger.With(
		zap.String(logKeyTopic, topic),
		zap.String(logKeyComponent, "subscriber"),
	))
}
