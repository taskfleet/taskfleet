package dymant

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// Publisher provides a way for publishing messages to a single message queue. Typically, an
// application should only create a single publisher for a particular message queue. Publishers
// can safely be shared across threads. Consult the documentation of the individual implementations
// for details on how the publisher works with concepts of differing tools.
type Publisher interface {
	// Publish is similar to PublishSync, however, it does not wait for confirmation. Transmission
	// errors are only logged. If an error is returned, it indicates a serious problem.
	Publish(key uuid.UUID, message proto.Message) error

	// PublishSync publishes the given message to the message queue. This function waits until the
	// publisher flags the message as successful. If the context is cancelled before the message
	// has been published, an error is returned. Depending on the implementation, message delivery
	// might still continue and the result of that delivery is logged in the background.
	//
	// The key passed to this function is purely informational and is not provided to consumers. It
	// serves as metadata that the message queue can use internally to partition messages.
	PublishSync(ctx context.Context, key uuid.UUID, message proto.Message) error

	// Flush waits for all messages to be delivered to the message queue. It blocks until all
	// messages have been delivered or the given context is cancelled. If the cancellation of the
	// context causes messages to not be delivered, an error is returned. In any case, this
	// function terminates the producer.
	Flush(ctx context.Context) error
}

// Subscriber represents a consumer of a single message queue. A subscriber has a message type
// attached into which all messages are unmarshaled.
type Subscriber interface {
	// Process initiates message consumption from the subscriber's message queue. All messages
	// received are marshaled into the subscriber's associated type and subsequently batched
	// according to batch configuration. The number of messages delivered is strictly greater than
	// zero.
	//
	// The function terminates once an error occurs during consumption, the callback returns an
	// error or the context is cancelled (the latter can have a delay of up to 100ms). The callback
	// should use the provided context to terminate once the upstream context is cancelled. This
	// also allows for timeouts to be applied to the callback. The returned error is therefore
	// NEVER nil and should be handled accordingly. As a convenience, this package exposes the
	// `IsErrContext` method.
	//
	// NOTE: In order to uphold delivery guarantees, the callback must process the messages
	// synchronously, i.e. no concurrent process must be started. In order to improve throughput,
	// modify the batch configuration instead.
	//
	// ATTENTION: This method must not be called multiple times at the same time. In order to
	// uphold delivery guarantees, it should be called exactly once on a particular subscriber.
	Process(ctx context.Context, execute func(context.Context, []proto.Message) error) error

	// Close ensures that the subscriber is properly cleaned up.
	Close()
}
