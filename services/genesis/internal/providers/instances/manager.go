package instances

import (
	"fmt"

	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// Manager is a type which manages a set of instances for a single provider zone. Its primary
// use case is to choose the instance that fits a particular request the best.
type Manager struct {
	instances []Type
}

// NewManager initializes a new manager which manages the provided set of instances.
func NewManager(instances []Type) *Manager {
	return &Manager{instances: instances}
}

// Types returns all instance types managed by this manager.
func (m *Manager) Types() []Type {
	return m.instances
}

// Validate ensures that the provided resources do not exhibit a skew that is too large.
func (m *Manager) Validate(resources Resources) bool {
	return validateMemoryPerCPU(resources.CPUCount, resources.MemoryMiB)
}

// GPUKinds returns the (unique) GPU kinds that are provided by the instances managed by this
// manager.
func (m *Manager) GPUKinds() []typedefs.GPUKind {
	kinds := map[typedefs.GPUKind]struct{}{}
	for _, instance := range m.instances {
		if instance.GPU != nil {
			kinds[instance.GPU.Kind] = struct{}{}
		}
	}
	return jack.MapKeys(kinds)
}

// Get returns the instance type with the specified name. If it cannot be found, an error is
// returned.
func (m *Manager) Get(name string) (Type, error) {
	for _, instance := range m.instances {
		if instance.Name == name {
			return instance, nil
		}
	}
	return Type{}, fmt.Errorf("could not find suitable instance type")
}

// FindBestFit finds the instance which provides the best fit for the requested resources. If no
// instance can be found which satisfies the request, an error is returned.
func (m *Manager) FindBestFit(resources Resources, arch typedefs.CPUArchitecture) (Type, error) {
	// First, filter the instances which satisfy the constraints
	filteredInstances := make([]Type, 0)
	for _, instance := range m.instances {
		if instance.covers(resources) && instance.Architecture == arch {
			filteredInstances = append(filteredInstances, instance)
		}
	}
	if len(filteredInstances) == 0 {
		return Type{},
			fmt.Errorf("could not find instance type which covers the requested resources")
	}

	// Then, we find the instance with the minimum price... (TODO: use real price)
	choice := filteredInstances[0]
	for i := 1; i < len(filteredInstances); i++ {
		if filteredInstances[i].artificialPrice() < choice.artificialPrice() {
			choice = filteredInstances[i]
		}
	}
	// ... and return the cheapest option
	return choice, nil
}
