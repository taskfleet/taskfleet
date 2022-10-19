package jack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceMap(t *testing.T) {
	values := []int{2, 3, 4}
	actual := SliceMap(values, func(v int) float64 {
		return float64(v) * 2.5
	})
	expected := []float64{5, 7.5, 10}
	assert.ElementsMatch(t, actual, expected)
}

func TestSliceFilterNil(t *testing.T) {
	two := 2
	three := 3
	values := []*int{&two, nil, &three}
	actual := SliceFilterNil(values)
	expected := []*int{&two, &three}
	assert.ElementsMatch(t, actual, expected)
}

func TestSliceFilterMap(t *testing.T) {
	values := []int{2, 3, 4}
	actual := SliceFilterMap(values, func(v int) (int, bool) {
		if v%2 == 0 {
			return v * 2, true
		}
		return 0, false
	})
	expected := []int{4, 8}
	assert.ElementsMatch(t, actual, expected)
}
