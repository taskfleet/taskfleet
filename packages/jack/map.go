package jack

import (
	"golang.org/x/exp/constraints"
)

// ZipMap zips the keys and values and returns a map from the zipped items. For creating the map,
// the first value is used as the key while the second value is used as the value.
func ZipMap[K constraints.Ordered, V any](keys []K, values []V) map[K]V {
	length := len(keys)
	if len(values) < length {
		length = len(values)
	}

	result := map[K]V{}
	for i := 0; i < length; i++ {
		result[keys[i]] = values[i]
	}
	return result
}

// MapKeys returns the keys from a map as a slice.
func MapKeys[K constraints.Ordered, V any](kv map[K]V) []K {
	result := []K{}
	for key := range kv {
		result = append(result, key)
	}
	return result
}
