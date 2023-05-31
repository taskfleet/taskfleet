package instances

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestNewManager(t *testing.T) {
	// No issue during initialization
	instances := []Type{
		{Name: "test1"},
		{Name: "test2"},
	}
	manager, err := NewManager(instances)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []string{"test1", "test2"}, jack.MapKeys(manager.instances))

	// Ensure that it fails when using instances with the same name, even if they are different
	instances = []Type{
		{Name: "test1", Architecture: typedefs.ArchitectureArm},
		{Name: "test1", Architecture: typedefs.ArchitectureX86},
	}
	_, err = NewManager(instances)
	assert.ErrorContains(t, err, "duplicate")
}

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
	manager, err := NewManager(instances)
	require.Nil(t, err)

	// Too many CPUs
	_, err = manager.FindBestFit(
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
