package instances

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMemoryPerCPU(t *testing.T) {
	// 784 GB
	assert.False(t, validateMemoryPerCPU(16, 802816))
	assert.True(t, validateMemoryPerCPU(36, 802816))
	assert.False(t, validateMemoryPerCPU(80, 802816))

	// 480 GB
	assert.False(t, validateMemoryPerCPU(32, 491520))
	assert.True(t, validateMemoryPerCPU(64, 491520))
	assert.False(t, validateMemoryPerCPU(96, 491520))

	// 112 GB
	assert.False(t, validateMemoryPerCPU(8, 114688))
	assert.True(t, validateMemoryPerCPU(12, 114688))
	assert.True(t, validateMemoryPerCPU(50, 114688))
	assert.False(t, validateMemoryPerCPU(64, 114688))

	// 48 GB
	assert.False(t, validateMemoryPerCPU(2, 49152))
	assert.True(t, validateMemoryPerCPU(6, 49152))
	assert.True(t, validateMemoryPerCPU(48, 49152))
	assert.False(t, validateMemoryPerCPU(80, 49152))
}
