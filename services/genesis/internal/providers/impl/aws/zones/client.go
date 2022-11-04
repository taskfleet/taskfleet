package awszones

import (
	"context"

	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// Client represents an AWS zone client.
type Client struct {
	zones map[string][]typedefs.GPUKind
}

// NewClient creates a new AWS zone client. The client caches the available GPU kinds for each zone
// from the provided instance managers.
func NewClient(ctx context.Context, instances map[string]*instances.Manager) *Client {
	gpuKinds := map[string][]typedefs.GPUKind{}
	for zone, manager := range instances {
		gpuKinds[zone] = manager.GPUKinds()
	}
	return &Client{gpuKinds}
}

// List implements the `providers.ZoneClient` interface.
func (c *Client) List() []providers.Zone {
	zones := []providers.Zone{}
	for zone, gpus := range c.zones {
		zones = append(zones, providers.Zone{Name: zone, GPUs: gpus})
	}
	return zones
}
