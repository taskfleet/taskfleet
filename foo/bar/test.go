package kafka

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type publisher struct {
	topic    string
	logger   *zap.Logger
	producer *kafka.Producer
	flush    chan context.Context
	done     chan error
}

func newPublisher(topic string, config kafka.ConfigMap, logger *zap.Logger) (*publisher, error) {
	kafkaProducer, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %s", err)
	}

	p := &publisher{
		topic:    topic,
		logger:   logger,
		producer: kafkaProducer,
		flush:    make(chan context.Context),
		done:     make(chan error),
	}
	go p.logMessages()

	return p, nil
}

//-------------------------------------------------------------------------------------------------
// DYMANT INTERFACE
//-------------------------------------------------------------------------------------------------

func (p *publisher) Publish(key uuid.UUID, message proto.Message) error {
	msg, err := p.buildMessage(key, message)
	if err != nil {
		return err
	}
	if err := p.producer.Produce(msg, nil); err != nil {
		return fmt.Errorf("failed to initiate publishing of message: %s", err)
	}
	return nil
}

func (p *publisher) PublishSync(ctx context.Context, key uuid.UUID, message proto.Message) error {
	msg, err := p.buildMessage(key, message)
	if err != nil {
		return err
	}

	ch := make(chan kafka.Event, 1)
	if err := p.producer.Produce(msg, ch); err != nil {
		return fmt.Errorf("failed to initiate publishing of message: %s", err)
	}

	select {
	case result := <-ch:
		if msg, ok := result.(*kafka.Message); ok {
			logProduced(p.logger, msg)
			if msg.TopicPartition.Error != nil {
				return fmt.Errorf("failed to publish message: %s", msg.TopicPartition.Error)
			}
		} else {
			return fmt.Errorf("received unknown event")
		}
	case <-ctx.Done():
		// Although we return an error, we want to log the message
		go func() {
			if msg, ok := (<-ch).(*kafka.Message); ok {
				logProduced(p.logger, msg)
			}
		}()
		return ctx.Err()
	}

	return nil
}

func (p *publisher) Flush(ctx context.Context) error {
	p.flush <- ctx
	p.producer.Close()
	return <-p.done
}

//-------------------------------------------------------------------------------------------------
// UTILITIES
//-------------------------------------------------------------------------------------------------

func (p *publisher) logMessages() {
	for {
		select {
		case event := <-p.producer.Events():
			if msg, ok := event.(*kafka.Message); ok {
				logProduced(p.logger, msg)
			}
		case ctx := <-p.flush:
			if err := p.awaitRemaining(ctx); err != nil {
				p.logger.Error("failed to flush remaining messages", zap.Error(err))
				p.done <- err
			} else {
				p.done <- nil
			}
			return
		}
	}
}

func (p *publisher) awaitRemaining(ctx context.Context) error {
	for p.producer.Len() > 0 {
		if ctx.Err() != nil {
			return fmt.Errorf("context cancelled (%s), failed to publish %d messages",
				ctx.Err(), p.producer.Len(),
			)
		}
		// We flush with a timeout of 100ms -- this equals the default timeout of the Flush method.
		// Afterwards, we can check for the context to be cancelled again
		p.producer.Flush(100)
	}
	return nil
}

//–------------------------------------------------------------------------------------------------
// SERIALIZATION
//-------------------------------------------------------------------------------------------------

func (p *publisher) buildMessage(key uuid.UUID, message proto.Message) (*kafka.Message, error) {
	encoded, err := proto.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling message: %s", err)
	}

	// Need kafka.PartitionAny or it is published to partition 0
	return &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
		Key:            key[:],
		Value:          encoded,
	}, nil
}
