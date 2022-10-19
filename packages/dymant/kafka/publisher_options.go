package kafka

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// PublisherOption allows to update the configuration of a Kafka publisher.
type PublisherOption interface {
	publisherOption()
	apply(config kafka.ConfigMap) error
}

type dummyPublisherOption struct{}

func (dummyPublisherOption) publisherOption() {}

//-------------------------------------------------------------------------------------------------
// CONSISTENCY
//-------------------------------------------------------------------------------------------------

// Consistency describes the level of consistency to achieve when using Kafka producers.
type Consistency int

const (
	// ConsistencyStrong results in a successful write if all in-sync replicas acknowledge the
	// write. Additionally, it enables idempotent writes to ensure exactly-once delivery. This is
	// the default consistency level.
	ConsistencyStrong Consistency = iota
	// ConsistencyWeak provides waits for a single broker to acknowledge the write.
	ConsistencyWeak
	// ConsistencyMinimal provides no publishing guarantees except for TCP acknowledgements.
	ConsistencyMinimal
)

type publisherOptionConsistency struct {
	dummyPublisherOption
	level Consistency
}

// WithConsistency configures the consistency level of the publisher which governs how "safely"
// published messages will be delivered.
func WithConsistency(level Consistency) PublisherOption {
	return publisherOptionConsistency{level: level}
}

func (c publisherOptionConsistency) apply(config kafka.ConfigMap) error {
	switch c.level {
	case ConsistencyMinimal:
		config["request.required.acks"] = 0
		config["enable.idempotence"] = false
	case ConsistencyWeak:
		config["request.required.acks"] = 1
		config["enable.idempotence"] = false
	case ConsistencyStrong:
		config["request.required.acks"] = -1
		config["enable.idempotence"] = true
	default:
		return fmt.Errorf("invalid publisher consistency level")
	}
	return nil
}
