package dymant

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	// NoKey initializes an empty UUID that may be used whenever a publisher does not a publisher
	// key.
	NoKey = uuid.UUID{}
)

// IsErrContext returns whether the given error originated from the context package.
func IsErrContext(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}
