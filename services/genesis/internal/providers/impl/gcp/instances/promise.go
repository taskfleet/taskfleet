package gcpinstances

import (
	"context"

	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"google.golang.org/api/compute/v1"
)

// InstancePromise describes an instance that is currently being created.
type InstancePromise struct {
	Meta      providers.InstanceMeta
	operation *compute.Operation
	client    *Client
}

// Await waits for the instance to materialize. The given context should have a sufficiently large
// timeout as this operation might take some time. If an error is returned, this usually means that
// the instance could not be created but could be caused by issues with the network.
func (p *InstancePromise) Await(ctx context.Context) (providers.Instance, error) {
	// First, we need to wait for the creation operation to finish
	if err := p.client.requester.Poll(ctx, p.operation); err != nil {
		return providers.Instance{}, providers.NewAPIError("failed to wait for instance", err)
	}
	// Upon success, we can query the full instance
	return p.client.Get(ctx, p.Meta)
}
