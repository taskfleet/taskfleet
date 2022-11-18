package providers

import (
	"context"

	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// Provider allows to interact with instances of a particular provider.
type Provider interface {
	// Provider returns the canonical name of this client's backing (cloud) provider.
	Name() typedefs.Provider
	// AccountName returns a provider-specific "account" name that the provider is interacting
	// with. For GCP, this is the project ID, for AWS, this is the account ID.
	AccountName() string
	// Zones returns a zone client that allows to interact with a provider's "zones".
	Zones() ZoneClient
	// Instances returns a client that allows to interact with the provider's instances.
	// The instance client is unusable until `Prepare` has been called on the provider.
	Instances() InstanceClient
}

// ZoneClient desccribes the public interface for interacting with a provider's "zones". For cloud
// providers, this means availability zones.
type ZoneClient interface {
	// List returns all of the provider's zones into which instances can be launched. For cloud
	// providers, this list may be a subset of all availability zones.
	List() []Zone
}

// InstanceClient describes the public interface for interacting with a provider's intances.
type InstanceClient interface {
	// Returns the best-fitting instance type in the given zone which accomodates the requested
	// resources and satisfies the given constraints. If an error is returned, this cloud provider
	// cannot fulfill the given request.
	Find(
		zone string, resources instances.Resources, architecture typedefs.CPUArchitecture,
	) (instances.Type, error)
	// Create creates an instance with the specified specification. The creation operation results
	// in an error if the initial call to the cloud provider fails. This does, however, not imply
	// that the instance was started up correctly. The returned promise can be used to await the
	// instance's boot.
	Create(ctx context.Context, ref InstanceMeta, spec InstanceSpec) (InstancePromise, error)
	// Get returns the instance that is uniquely identified by the provided reference. If the
	// instance is not currently running, an error is returned.
	Get(ctx context.Context, ref InstanceMeta) (Instance, error)
	// List returns all instances that are currently running on this cloud provider.
	List(ctx context.Context) ([]Instance, error)
	// Delete deletes the instance uniquely identified by the provided reference. Once this method
	// returns, the instance is guaranteed to be deleted. Subsequent calls to `Get` will fail
	// while `List` will not include this instance.
	Delete(ctx context.Context, ref InstanceMeta) error
}

// InstancePromise is a type which is issued to wait for an instance that is starting up.
type InstancePromise interface {
	// Await waits until the instance is booted and returns the booted instance upon success.
	Await(ctx context.Context) (Instance, error)
}
