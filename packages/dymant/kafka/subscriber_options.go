package kafka

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// SubscriberOption allows to update the configuration of a Kafka subscriber.
type SubscriberOption interface {
	subscriberOption()
	apply(config kafka.ConfigMap) error
	configApply(config *subscriberConfig)
}

type dummySubscriberOption struct{}

func (dummySubscriberOption) subscriberOption() {}

//-------------------------------------------------------------------------------------------------
// FETCH LEVEL
//-------------------------------------------------------------------------------------------------

// Fetch governs how message offsets are committed to Kafka for consumers.
type Fetch int

const (
	// FetchAtLeastOnce commits message offsets after a message was processed. Hence, multiple
	// messages might be replayed after failure. Note that at-least-once delivery is the best
	// option for exactly-once semantics (EOS) with external systems. This is the default.
	FetchAtLeastOnce Fetch = iota
	// FetchAtMostOnce commits message offsets right when the message was received and before it
	// was processed. This ensures that a message is never processed multiple times but may cause
	// losing the message if the application terminates while the message is being processed.
	FetchAtMostOnce
	// FetchAny commits message offsets periodically in the background. This is the most
	// efficient setting but there are no guarantees regarding message delivery.
	FetchAny
)

type subscriberOptionFetch struct {
	dummySubscriberOption
	level Fetch
}

// WithFetch specifies a particular fetch level which governs how offsets are committed to Kafka.
// If this level is not provided, `FetchAtLeastOnce` is used.
func WithFetch(level Fetch) SubscriberOption {
	return subscriberOptionFetch{level: level}
}

func (c subscriberOptionFetch) apply(config kafka.ConfigMap) error {
	switch c.level {
	case FetchAtLeastOnce, FetchAtMostOnce:
		config["enable.auto.commit"] = false
		config["go.application.rebalance.enable"] = true
	case FetchAny:
		config["enable.auto.commit"] = true
		config["go.application.rebalance.enable"] = false
	default:
		return fmt.Errorf("invalid subscriber fetch level")
	}
	return nil
}

func (c subscriberOptionFetch) configApply(config *subscriberConfig) {
	config.fetch = c.level
}

//-------------------------------------------------------------------------------------------------
// BATCH CONFIGURATION
//-------------------------------------------------------------------------------------------------

type subscriberOptionBatchConfig struct {
	dummySubscriberOption
	bufferSize  int
	aggregation time.Duration
}

// WithBatchConfig describes configuration options to set for receiving messages as batches.
// Batching is disabled by default but may greatly improve throughput.
//
// The buffer size describes how many messages are aggregated at most before being delivered to
// the client. This value may be arbitrarily large.
//
// The aggregation duration describes for how long to aggregate consumed messages. Increasing this
// duration results in larger batches. However, it also results in a greater message delay and
// more messages being replayed/lost (depending on the selected delivery guarantees) in case of
// failure. If the buffer is filled before the aggregation duration is reached, the messages are
// delivered right away. You should set the aggregation to a value which satisfies your
// requirements.  If this field is set to a negative value, messages are not aggregated, i.e.
// single messages are delivered to the client. This automatically sets the buffer size to 1
// regardless of the configuration. Processing single messages might be useful if processing a
// message is very costly.
//
// If this option is not set, messages are delivered without batching.
func WithBatchConfig(bufferSize int, aggregation time.Duration) SubscriberOption {
	return subscriberOptionBatchConfig{bufferSize: bufferSize, aggregation: aggregation}
}

func (c subscriberOptionBatchConfig) apply(config kafka.ConfigMap) error {
	return nil
}

func (c subscriberOptionBatchConfig) configApply(config *subscriberConfig) {
	config.batchBufferSize = c.bufferSize
	config.batchAggregation = c.aggregation
}

//-------------------------------------------------------------------------------------------------
// TIMEOUT
//-------------------------------------------------------------------------------------------------

type subscriberOptionTimeout struct {
	dummySubscriberOption
	timeout time.Duration
}

// WithProcessingTimeout sets a timeout on the duration a batch of messages may be processed. If
// this timeout is exceeded, message processing is considered to have failed. If this option is not
// set, messages may process for eternity.
func WithProcessingTimeout(timeout time.Duration) SubscriberOption {
	return subscriberOptionTimeout{timeout: timeout}
}

func (c subscriberOptionTimeout) apply(config kafka.ConfigMap) error {
	return nil
}

func (c subscriberOptionTimeout) configApply(config *subscriberConfig) {
	config.callbackTimeout = c.timeout
}
