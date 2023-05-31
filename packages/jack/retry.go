package jack

import (
	"context"
	"time"
)

type retryOptions struct {
	backoff    time.Duration
	multiplier time.Duration
	condition  func(error) bool
}

// WithRetry executes the provided function until the context is cancelled or the function does not
// return an error. Specifics of the retry logic can be adapted by passing additional options to
// this function.
func WithRetry[R any](
	ctx context.Context, do func() (R, error), options ...RetryOption,
) (R, error) {
	// Get full options
	o := retryOptions{
		backoff:    100 * time.Millisecond,
		multiplier: 2,
		condition: func(error) bool {
			return true
		},
	}
	for _, option := range options {
		option.apply(&o)
	}

	// Specify outcomes
	var err error
	var resource R

	// Iterate with backoff
	for {
		resource, err = do()
		if err == nil {
			return resource, nil
		}
		if !o.condition(err) {
			return resource, err
		}
		select {
		case <-ctx.Done():
			return resource, ctx.Err()
		case <-time.After(o.backoff):
			o.backoff = o.backoff * o.multiplier
		}
	}
}

// RetryOption allows to specify options for retry operations.
type RetryOption interface {
	apply(r *retryOptions)
}

//-------------------------------------------------------------------------------------------------

// WithInitialBackoff sets the initial backoff duration which defaults to 100ms.
func WithInitialBackoff(duration time.Duration) RetryOption {
	return backoff(duration)
}

type backoff time.Duration

func (b backoff) apply(r *retryOptions) {
	r.backoff = time.Duration(b)
}

//-------------------------------------------------------------------------------------------------

// WithBackoffMultiplier sets the multiplier for exponential backoff.
func WithBackoffMultiplier(mult float64) RetryOption {
	return multiplier(mult)
}

type multiplier float64

func (m multiplier) apply(r *retryOptions) {
	r.multiplier = time.Duration(m)
}

//-------------------------------------------------------------------------------------------------

// WithCondition sets the condition on when to execute a retry operation. By default, a retry is
// happening whenever an error occurs.
func WithCondition(c func(error) bool) RetryOption {
	return condition(c)
}

type condition func(error) bool

func (c condition) apply(r *retryOptions) {
	r.condition = c
}
