package jack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMin(t *testing.T) {
	// Integers
	assert.Equal(t, 5, Min(5, 6))
	assert.Equal(t, 5, Min(6, 5))
	assert.Equal(t, 5, Min(5, 5))

	// Floats
	assert.Equal(t, 5.0, Min(5.0, 6.0))
	assert.Equal(t, 5.0, Min(6.0, 5.0))
	assert.Equal(t, 5.0, Min(5.0, 5.0))
}
