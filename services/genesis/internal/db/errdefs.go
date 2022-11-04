package db

import "errors"

var (
	// ErrNotExist indicates that a type could not be found.
	ErrNotExist = errors.New("object with specified ID could not be found")
)

// ValidationError describes an error that originated as a result of an invalid input.
type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}
