package instances

import "go.taskfleet.io/services/genesis/internal/typedefs"

// Type represents a single instance (or machine) type.
type Type struct {
	Resources
	// The provider-specific name of the instance type.
	Name string
	// The provider-specific ID of the instance type (e.g. an ARN for AWS or a link for GCP).
	UID string
	// The CPU architecture of the instance type.
	Architecture typedefs.CPUArchitecture
}

// Resources represents an instance's available resources.
type Resources struct {
	// The number of CPUs provided by the instance.
	CPUCount uint16
	// The number of MiB of memory provided by the instance.
	MemoryMiB uint32
	// The (optional) GPU configuration of the instance.
	GPU *GPUResources
}

// GPUResources describes the GPU configuration.
type GPUResources struct {
	// The type of GPU.
	Kind typedefs.GPUKind
	// The number of GPUs.
	Count uint16
}

//-------------------------------------------------------------------------------------------------
// METHODS
//-------------------------------------------------------------------------------------------------

// GPUKind returns the GPU kind when a GPU is referenced in the resources and `nil` otherwise.
func (r Resources) GPUKind() *typedefs.GPUKind {
	if r.GPU == nil {
		return nil
	}
	return &r.GPU.Kind
}

func (t Type) artificialPrice() float64 {
	// 1 GPU = 50 CPUs = 375 GiB memory
	result := float64(t.CPUCount)*7.5 + float64(t.MemoryMiB)/1024
	if t.GPU != nil {
		result += float64(t.GPU.Count) * 375
	}
	return result
}

func (r Resources) covers(other Resources) bool {
	if r.CPUCount < other.CPUCount {
		return false
	}
	if r.MemoryMiB < other.MemoryMiB {
		return false
	}
	if other.GPU != nil && r.GPU == nil {
		return false
	}
	if other.GPU != nil && r.GPU != nil {
		return r.GPU.covers(*other.GPU)
	}
	return true
}

func (g GPUResources) covers(other GPUResources) bool {
	if g.Kind != other.Kind {
		return false
	}
	return g.Count >= other.Count
}
