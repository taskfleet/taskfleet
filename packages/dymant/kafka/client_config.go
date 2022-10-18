package kafka

import (
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type clientConfig struct {
	id               string
	bootstrapServers []string
	options          []ClientOption
}

//-------------------------------------------------------------------------------------------------

func (c clientConfig) config() (kafka.ConfigMap, error) {
	config := kafka.ConfigMap{
		"api.version.request": true,
		"bootstrap.servers":   strings.Join(c.bootstrapServers, ","),
		"client.id":           c.id,
		"security.protocol":   "plaintext",
	}
	for _, option := range c.options {
		if err := option.apply(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (c clientConfig) producerConfig(options []PublisherOption) (kafka.ConfigMap, error) {
	config, err := c.config()
	if err != nil {
		return nil, err
	}

	// Apply basic config
	config["partitioner"] = "consistent_random"

	// Apply defaults
	WithConsistency(ConsistencyStrong).apply(config) // nolint:errcheck

	// Apply options
	for _, option := range options {
		if err := option.apply(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (c clientConfig) consumerConfig(
	group string, options []SubscriberOption,
) (kafka.ConfigMap, error) {
	config, err := c.config()
	if err != nil {
		return nil, err
	}

	// Apply basic config
	config["group.id"] = group
	config["auto.offset.reset"] = "earliest"
	config["partition.assignment.strategy"] = "range,roundrobin"
	config["isolation.level"] = "read_committed"

	// Apply defaults
	WithFetch(FetchAtLeastOnce).apply(config) // nolint:errcheck

	// Apply options
	for _, option := range options {
		if err := option.apply(config); err != nil {
			return nil, err
		}
	}
	return config, nil
}
