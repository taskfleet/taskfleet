package jack

// Must returns the provided value if the given error is nil and panics with the error otherwise.
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Ptr returns a pointer to the provided value so that no extra variable needs to be declared.
func Ptr[T any](value T) *T {
	return &value
}
