package jack

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

func twoRetries(i *int) (int, error) {
	if *i == 2 {
		return 5, nil
	}
	*i += 1
	return 0, errTest
}

func TestWithRetry(t *testing.T) {
	// Should receive an error when not retrying
	i := 0
	_, err := twoRetries(&i)
	assert.NotNil(t, err)

	// Should work with retry (even with no backoff duration)
	ctx := context.Background()
	i = 0
	result, err := WithRetry(ctx, func() (int, error) {
		return twoRetries(&i)
	}, WithInitialBackoff(0))
	assert.Equal(t, result, 5)
	assert.Nil(t, err)
}

func TestWithBackoffMultiplier(t *testing.T) {
	ctx := context.Background()

	// Should not work if not specifying backoff multiplier
	tctx, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	i := 0
	_, err := WithRetry(tctx, func() (int, error) {
		return twoRetries(&i)
	}, WithInitialBackoff(4*time.Millisecond))
	assert.NotNil(t, err)

	// Should work if setting it to 1
	tctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()

	i = 0
	_, err = WithRetry(tctx, func() (int, error) {
		return twoRetries(&i)
	}, WithInitialBackoff(4*time.Millisecond), WithBackoffMultiplier(1))
	assert.Nil(t, err)
}

func TestWithCondition(t *testing.T) {
	ctx := context.Background()
	i := 0
	_, err := WithRetry(ctx, func() (int, error) {
		return twoRetries(&i)
	}, WithCondition(func(err error) bool {
		return err != errTest
	}))
	assert.NotNil(t, err)
}
