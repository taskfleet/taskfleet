package jack

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZipMap(t *testing.T) {
	keys := []int{2, 3, 4}
	values := []string{"2", "3", "4"}

	// Matching size
	actual := ZipMap(keys, values)
	expected := map[int]string{2: "2", 3: "3", 4: "4"}
	assert.True(t, reflect.DeepEqual(actual, expected))

	// Less keys
	keysSmall := []int{2, 3}
	actual = ZipMap(keysSmall, values)
	expected = map[int]string{2: "2", 3: "3"}
	assert.True(t, reflect.DeepEqual(actual, expected))

	// Less values
	valuesSmall := []string{"3", "4"}
	actual = ZipMap(keys, valuesSmall)
	expected = map[int]string{2: "3", 3: "4"}
	assert.True(t, reflect.DeepEqual(actual, expected))
}

func TestMapKeys(t *testing.T) {
	data := map[int]string{2: "2", 3: "3", 4: "4"}
	actual := MapKeys(data)
	expected := []int{2, 3, 4}
	assert.ElementsMatch(t, actual, expected)
}

func TestMapValues(t *testing.T) {
	data := map[int]string{2: "2", 3: "3", 4: "4"}
	actual := MapValues(data)
	expected := []string{"2", "3", "4"}
	assert.ElementsMatch(t, actual, expected)
}
