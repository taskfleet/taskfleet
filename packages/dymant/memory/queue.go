package memory

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

// Queue represents a message queue that resides purely in-memory and can be used for testing
// purposes. The queue is thread-safe and not tuned for performance. The queue never delivers
// messages in batches.
type Queue struct {
	ch chan proto.Message
}

// NewQueue initializes a new message queue that resides entirely in memory. The queue is both
// a publisher and a subscriber and provides convenience methods for easily setting/getting
// messages. The queue may grow up to the specified size.
func NewQueue(size int) *Queue {
	return &Queue{
		ch: make(chan proto.Message, size),
	}
}

//-------------------------------------------------------------------------------------------------
// PUBLISHER
//-------------------------------------------------------------------------------------------------

// Publish implements the dymant.Publisher interface.
func (q *Queue) Publish(key uuid.UUID, message proto.Message) error {
	q.ch <- message
	return nil
}

// PublishSync implements the dymant.Publisher interface.
func (q *Queue) PublishSync(ctx context.Context, key uuid.UUID, message proto.Message) error {
	return q.Publish(key, message)
}

// Flush implements the dymant.Publisher interface.
func (q *Queue) Flush(ctx context.Context) error {
	return nil
}

//-------------------------------------------------------------------------------------------------
// SUBSCRIBER
//-------------------------------------------------------------------------------------------------

// Process implements the dymant.Subscriber interface.
func (q *Queue) Process(
	ctx context.Context, execute func(context.Context, []proto.Message) error,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case next, ok := <-q.ch:
			if !ok {
				return fmt.Errorf("channel closed unexpectedly")
			}
			if err := execute(ctx, []proto.Message{next}); err != nil {
				return err
			}
		}
	}
}

// Close implements the dymant.Subscriber interface.
func (q *Queue) Close() {
	close(q.ch)
}

//-------------------------------------------------------------------------------------------------
// CONVENIENCE
//-------------------------------------------------------------------------------------------------

// SetMessages is a convenience function to add the provided messages to the queue.
func (q *Queue) SetMessages(messages []proto.Message) {
	for _, msg := range messages {
		q.ch <- msg
	}
}

// GetMessages is a convenience function to get all messages published to the queue.
func (q *Queue) GetMessages() []proto.Message {
	result := make([]proto.Message, 0)
	for {
		select {
		case msg, ok := <-q.ch:
			if !ok {
				return result
			}
			result = append(result, msg)
		default:
			return result
		}
	}
}
