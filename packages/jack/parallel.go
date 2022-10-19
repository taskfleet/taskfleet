package jack

import (
	"context"

	"golang.org/x/sync/errgroup"
)

// ParallelSliceMap executes the provided function with all inputs in parallel and returns the
// mapped values on success.
func ParallelSliceMap[T any, R any](
	ctx context.Context, inputs []T, execute func(context.Context, T) (R, error),
) ([]R, error) {
	result := make([]R, len(inputs))
	eg, ctx := errgroup.WithContext(ctx)
	for i, item := range inputs {
		index := i
		input := item
		eg.Go(func() error {
			out, err := execute(ctx, input)
			if err != nil {
				return err
			}
			result[index] = out
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return result, nil
}
