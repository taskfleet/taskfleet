package db

import (
	"context"

	"github.com/google/uuid"
)

// InstanceStatus describes the status of an instance.
type InstanceStatus int

const (
	// InstanceStatusRequested describes a state where an instance was requested from the cloud
	// provider but it is unknown whether the instance was eventually created successfully.
	InstanceStatusRequested InstanceStatus = iota

	// InstanceStatusBooting describes a state where an instance was successfully created on the
	// cloud provider but it is still starting up.
	InstanceStatusBooting

	// InstanceStatusRunning describes a state where an instance started up and is currently
	// running.
	InstanceStatusRunning

	// InstanceStatusDeleted describes a state where an instance was deleted and does surely not
	// exist anymore. It is not guaranteed that this instance was ever running.
	InstanceStatusDeleted
)

// Connection describes the public interface for interacting with the database.
type Connection interface {
	// CreateInstance creates a new instance in the database. An error is returned if the given
	// instance cannot be validated.
	CreateInstance(ctx context.Context, instance Instance) (*Instance, error)

	// Get returns the instance with the given ID or an error if it cannot be found.
	GetInstance(ctx context.Context, id uuid.UUID) (*Instance, error)

	// List returns all instances matching the specified filter. Note that, at a particular
	// point in time, the sets of instances for the different statuses are mutually distinct. If
	// an error occurs, it will only be delivered to the client as the first item of the returned
	// list is accessed.
	ListInstances(ctx context.Context, filters ...Filter) InstanceIterator
}

// Filter specifies a filter to be used for restricting the set of instances to be returned from a
// database query.
type Filter interface {
	// Statement returns the statement with %s specifiers for any placeholders.
	Statement() string

	// Values returns a list of values to be substituted for the placeholders returned by the
	// `Statement` method.
	Values() []interface{}
}
