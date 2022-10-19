package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type subscriber struct {
	topic           string
	config          subscriberConfig
	logger          *zap.Logger
	consumer        *kafka.Consumer
	messageTemplate proto.Message
	buf             []proto.Message
}

type subscriberConfig struct {
	fetch            Fetch
	callbackTimeout  time.Duration
	batchBufferSize  int
	batchAggregation time.Duration
}

func newSubscriber(
	topic string,
	config kafka.ConfigMap,
	subscriberConfig subscriberConfig,
	message proto.Message,
	logger *zap.Logger,
) (*subscriber, error) {
	kafkaConsumer, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %s", err)
	}
	if err := kafkaConsumer.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("failed to initiate subscription for topic: %s", err)
	}

	// Initialize buffer
	var buf []proto.Message
	if subscriberConfig.batchAggregation <= 0 || subscriberConfig.batchBufferSize <= 0 {
		buf = make([]proto.Message, 0, 1)
	} else {
		buf = make([]proto.Message, 0, subscriberConfig.batchBufferSize)
	}

	// Create subscriber
	return &subscriber{
		topic:           topic,
		config:          subscriberConfig,
		logger:          logger,
		consumer:        kafkaConsumer,
		messageTemplate: message,
		buf:             buf,
	}, nil
}

//-------------------------------------------------------------------------------------------------

func (c *subscriber) Process(
	ctx context.Context, execute func(context.Context, []proto.Message) error,
) error {
	deadline := time.Now().Add(c.config.batchAggregation)
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		// Get next message(s). If an error occurs, we always abort.
		if err := c.next(ctx, deadline); err != nil {
			return err
		}

		// We need to set the deadline for the next iteration - otherwise, an overly long message
		// processing might further delay messages
		deadline = time.Now().Add(c.config.batchAggregation)

		// If we didn't receive messages in the batch, we don't need to return anything
		if len(c.buf) > 0 {
			// For at-most-once delivery, we commit the consumer offset before delivering to
			// the client
			if c.config.fetch == FetchAtMostOnce {
				if err := c.commit(); err != nil {
					return err
				}
			}

			callbackCtx, cancel := c.callbackContext(ctx)
			err := func() error {
				defer cancel()
				return execute(callbackCtx, c.buf)
			}()
			if err != nil {
				return err
			}

			// For at-least-once delivery, we can commit the offset as soon as the messages are
			// processed.
			if c.config.fetch == FetchAtLeastOnce {
				if err := c.commit(); err != nil {
					return err
				}
			}
		}
	}
}

func (c *subscriber) Close() {
	if err := c.consumer.Close(); err != nil {
		c.logger.Error("failed to close consumer", zap.Error(err))
	}
}

//–------------------------------------------------------------------------------------------------

var (
	errPartitionsRevoked = errors.New("partitions revoked")
	errNoEvent           = errors.New("no event")
)

func (c *subscriber) callbackContext(ctx context.Context) (context.Context, func()) {
	if c.config.callbackTimeout > 0 {
		return context.WithTimeout(ctx, c.config.callbackTimeout)
	}
	return ctx, func() {}
}

func (c *subscriber) next(ctx context.Context, deadline time.Time) error {
	c.clearBuf()

	for {
		// Poll next message - if the context expired, we don't return an error but return no
		// messages
		if ctx.Err() != nil {
			c.clearBuf()
			return nil
		}
		timeout := c.timeout(deadline)
		if timeout < 0 {
			// If the timeout is already exceeded, we can return
			return nil
		}

		// We can continue polling with a non-negative timeout
		msg, err := c.poll(int(timeout.Round(time.Millisecond).Milliseconds()))
		if err != nil {
			if err == errNoEvent {
				// In case no event occurred, we can just continue. Since the timeout should be
				// negative, nothing will happen.
				continue
			} else if err == errPartitionsRevoked {
				// If partitions are revoked, we need to clear all messages that we have received
				// in the current iteration. The messages of the partitions that will be assigned
				// again are automatically fetched.
				c.clearBuf()
			} else {
				// Otherwise, an actual error must have occurred.
				return err
			}
		}
		if msg == nil {
			// Might occur for unexpected non-error events -- also, the above errors don't
			// necessarily terminate the iteration
			continue
		}

		c.buf = append(c.buf, msg)
		if len(c.buf) == cap(c.buf) {
			return nil
		}
	}
}

func (c *subscriber) timeout(deadline time.Time) time.Duration {
	if c.config.batchAggregation <= 0 {
		return 100 * time.Millisecond
	}
	return min(time.Until(deadline), 100*time.Millisecond)
}

func (c *subscriber) poll(timeoutMs int) (proto.Message, error) {
	event := c.consumer.Poll(timeoutMs)
	if event == nil {
		// timeout exceeded
		return nil, errNoEvent
	}

	// Process the event
	switch item := event.(type) {
	case *kafka.Message:
		// First, we log the message and check for an error
		logConsumed(c.logger, item)
		if item.TopicPartition.Error != nil {
			return nil, fmt.Errorf("failed to read message: %s", item.TopicPartition.Error)
		}

		// Finally, we can parse it
		msg := proto.Clone(c.messageTemplate)
		if err := proto.Unmarshal(item.Value, msg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %s", err)
		}
		return msg, nil
	case kafka.RevokedPartitions:
		// On partition revocation, we need to remove the partitions from the consumer and tell
		// the caller that a new assignment is imminent.
		if c.logger.Core().Enabled(zap.DebugLevel) {
			c.logger.Debug("received partition revocation", logFieldPartitions(item.Partitions))
		}
		if err := c.consumer.Unassign(); err != nil {
			c.logger.Warn("failed to revoke partition", zap.Error(err))
		}
		return nil, errPartitionsRevoked
	case kafka.AssignedPartitions:
		// On partition assignment, we simply assign the consumer to the new partitions.
		if c.logger.Core().Enabled(zap.DebugLevel) {
			c.logger.Info("received partition assignment", logFieldPartitions(item.Partitions))
		}
		if err := c.consumer.Assign(item.Partitions); err != nil {
			c.logger.Warn("failed to assign partition", zap.Error(err))
		}
	case kafka.Error:
		// Errors are informational, so we only log them except if all brokers are down
		if item.Code() == kafka.ErrAllBrokersDown {
			c.logger.Error("failed to connect to all brokers", zap.Error(item))
			return nil, item
		}
		c.logger.Warn("received error from Kafka", zap.Error(item))
	case kafka.OffsetsCommitted:
		if c.logger.Core().Enabled(zap.DebugLevel) {
			c.logger.Debug("committed offsets", logFieldsOffsets(item.Offsets)...)
		}
	default:
		// We do nothing for unknown event types (this should be only kafka.PartitionEOF which
		// should never be emitted)
		c.logger.Debug("received unexpected event", zap.String("event", item.String()))
	}

	return nil, nil
}

func (c *subscriber) commit() error {
	offsets, err := c.consumer.Commit()
	if err != nil {
		return err
	}
	if c.logger.Core().Enabled(zap.DebugLevel) {
		c.logger.Debug("committed offsets", logFieldsOffsets(offsets)...)
	}
	return nil
}

func (c *subscriber) clearBuf() {
	c.buf = c.buf[:0]
}

//-------------------------------------------------------------------------------------------------

func min(lhs, rhs time.Duration) time.Duration {
	if lhs < rhs {
		return lhs
	}
	return rhs
}
