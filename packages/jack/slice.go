package jack

//-------------------------------------------------------------------------------------------------
// MAP
//-------------------------------------------------------------------------------------------------

// SliceMap transforms all values according to the given function and returns the result.
func SliceMap[T any, R any](values []T, f func(T) R) []R {
	result := make([]R, 0, len(values))
	for i := range values {
		result = append(result, f(values[i]))
	}
	return result
}

//-------------------------------------------------------------------------------------------------
// FILTER
//-------------------------------------------------------------------------------------------------

// SliceFilterNil removes all values which are nil from the given slice.
func SliceFilterNil[T any](values []*T) []*T {
	result := make([]*T, 0, len(values))
	for i := range values {
		if values[i] != nil {
			result = append(result, values[i])
		}
	}
	return result
}

// SliceFilterMap transforms all values according to the given function and excludes the values for
// which false is returned.
func SliceFilterMap[T any, R any](values []T, f func(T) (R, bool)) []R {
	result := make([]R, 0, len(values))
	for i := range values {
		value, use := f(values[i])
		if use {
			result = append(result, value)
		}
	}
	return result
}
