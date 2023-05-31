package gcpinstances

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
)

// InstancePromise describes an instance that is currently being created.
type instancePromise struct {
	meta      providers.InstanceMeta
	operation *compute.Operation
	client    *Client
}

// Await waits for the instance to materialize. The given context should have a sufficiently large
// timeout as this operation might take some time. If an error is returned, this usually means that
// the instance could not be created but could be caused by issues with the network.
func (p *instancePromise) Await(ctx context.Context) (providers.Instance, error) {
	// First, we need to wait for the creation operation to finish
	if err := p.operation.Wait(ctx); err != nil {
		return providers.Instance{}, fmt.Errorf("failed to wait for instance creation: %w", err)
	}
	if p.operation.Proto().Error != nil {
		return providers.Instance{}, fmt.Errorf(
			"failed to create instance: %s", p.operation.Proto().Error,
		)
	}
	// Upon success, we can query the full instance
	return p.client.Get(ctx, p.meta)
}
