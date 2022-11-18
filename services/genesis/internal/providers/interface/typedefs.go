package providers

import (
	"net"
	"time"

	"github.com/google/uuid"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

//-------------------------------------------------------------------------------------------------
// PROVIDER INFORMATION
//-------------------------------------------------------------------------------------------------

// Zone describes a single zone of a provider.
type Zone struct {
	// The provider-specific name of the zone.
	Name string
	// The types of GPUs available in the zone.
	GPUs []typedefs.GPUKind
}

//-------------------------------------------------------------------------------------------------
// INSTANCE MANAGEMENT
//-------------------------------------------------------------------------------------------------

// Instance represents a running instance on a cloud provider.
type Instance struct {
	// Information that uniquely identifies the instance.
	Meta InstanceMeta
	// The static specification of the instance.
	Spec InstanceSpec
	// The current status of the instance.
	Status InstanceStatus
}

// InstanceMeta describes a globally unique reference to a running instance.
type InstanceMeta struct {
	// The globally unique ID assigned by Genesis.
	ID uuid.UUID
	// The unique provider-specific ID. Must not be set when an instance is created.
	ProviderID string
	// The provider-specific zone where the instance is located.
	ProviderZone string
}

// InstanceSpec describes the full information about an instance. This specification is available
// prior to the existence of the instance.
type InstanceSpec struct {
	// The instance type that implicitly defines number of CPUs and available memory.
	InstanceType instances.Type
	// Whether the instance may be shut down at any time.
	IsSpot bool
}

// InstanceStatus represents information that is available for an instance only once the instance
// has launched.
type InstanceStatus struct {
	// The time at which the instance was created.
	CreationTimestamp time.Time
	// The instance's network status.
	Network InstanceNetworkStatus
}

// InstanceNetworkStatus describes the network status of an instance.
type InstanceNetworkStatus struct {
	// The internal IP of the instance.
	InternalIP net.IP
	// The external IP of the instance (might not exist).
	ExternalIP net.IP
	// The internal hostname of the instance.
	InternalHostname string
}
