package instances

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestFindBestFit(t *testing.T) {
	instances := []Type{
		{
			Name:         "test1",
			Resources:    Resources{CPUCount: 4, MemoryMiB: 4096},
			Architecture: typedefs.ArchitectureX86,
		},
		{
			Name:         "test2",
			Resources:    Resources{CPUCount: 2, MemoryMiB: 4096},
			Architecture: typedefs.ArchitectureX86,
		},
		{
			Name:         "test3",
			Resources:    Resources{CPUCount: 8, MemoryMiB: 8192},
			Architecture: typedefs.ArchitectureX86,
		},
	}
	manager := NewManager(instances)

	// Too many CPUs
	_, err := manager.FindBestFit(
		Resources{CPUCount: 10, MemoryMiB: 1024}, typedefs.ArchitectureX86,
	)
	assert.NotNil(t, err)

	// Only one possibility
	instance, _ := manager.FindBestFit(
		Resources{CPUCount: 6, MemoryMiB: 6500}, typedefs.ArchitectureX86,
	)
	assert.Equal(t, instance.Name, "test3")

	// Choose the cheapest one
	instance, _ = manager.FindBestFit(
		Resources{CPUCount: 1, MemoryMiB: 1024}, typedefs.ArchitectureX86,
	)
	assert.Equal(t, instance.Name, "test2")
}
