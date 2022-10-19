package jack

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParallelSliceMap(t *testing.T) {
	ctx := context.Background()
	values := []int{2, 3, 4}

	// Without issues
	actual, err := ParallelSliceMap(ctx, values, func(ctx context.Context, v int) (int, error) {
		return v * 2, nil
	})
	assert.Nil(t, err)
	expected := []int{4, 6, 8}
	assert.ElementsMatch(t, actual, expected)

	// With error
	raise := errors.New("test")
	_, err = ParallelSliceMap(ctx, values, func(ctx context.Context, v int) (int, error) {
		return 0, raise
	})
	assert.Equal(t, err, raise)
}
