package gcputils

import (
	"context"
	"fmt"

	"google.golang.org/api/iterator"
)

// Iterator describes any types returned from the API that allow iteration.
type Iterator[T any] interface {
	Next() (T, error)
}

// Iterate allows to iterator over an iterator and execute the provided function for every item.
func Iterate[T any](ctx context.Context, it Iterator[T], do func(T) error) error {
	for {
		next, err := it.Next()
		switch err {
		case nil:
			if err := do(next); err != nil {
				return err
			}
		case iterator.Done:
			return nil
		default:
			return fmt.Errorf("failed to iterate over results: %w", err)
		}
	}
}
