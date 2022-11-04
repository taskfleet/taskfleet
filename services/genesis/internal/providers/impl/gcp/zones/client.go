package gcpzones

import (
	"context"
	"fmt"
	"sync"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/borchero/zeus/pkg/zeus"
	"go.taskfleet.io/packages/jack"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"go.uber.org/zap"
)

// ZoneInfo provides information about a GCP zone.
type ZoneInfo struct {
	Accelerators []Accelerator
	Subnetwork   string
}

// Client represents a GCP zone client.
type Client struct {
	clients
	project string
	network string
	zones   map[string]ZoneInfo
	// Mutex to sync zones periodically.
	mutex sync.Mutex
}

type clients struct {
	Zones            *compute.ZonesClient
	Networks         *compute.NetworksClient
	AcceleratorTypes *compute.AcceleratorTypesClient
}

// NewClient initializes a new GCP zone client. Upon calling this function, all zone information
// is fetched. The given context is used to periodically update the zone info.
func NewClient(
	ctx context.Context,
	zonesClient *compute.ZonesClient,
	networksClient *compute.NetworksClient,
	acceleratorTypesClient *compute.AcceleratorTypesClient,
	projectID string,
	network string,
) (*Client, error) {
	c := clients{
		Zones:            zonesClient,
		Networks:         networksClient,
		AcceleratorTypes: acceleratorTypesClient,
	}

	// When initializing the zone client, we want to fetch all available zones, then fetch the
	// accelerators which are available within each of these zones. Then, the zones are refreshed
	// once a day, until the context is cancelled.
	info, err := fetchZoneInfo(ctx, c, projectID, network)
	if err != nil {
		return nil, err
	}

	client := &Client{zones: info}
	go client.updateZonesPeriodically(ctx)
	return client, nil
}

//-------------------------------------------------------------------------------------------------
// INTERFACE
//-------------------------------------------------------------------------------------------------

// List implements the `providers.ZoneClient` interface.
func (c *Client) List() []providers.Zone {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make([]providers.Zone, 0)
	for zone, info := range c.zones {
		gpuKinds := make([]typedefs.GPUKind, 0)
		for _, accelerator := range info.Accelerators {
			gpuKinds = append(gpuKinds, accelerator.kind)
		}
		result = append(result, providers.Zone{Name: zone, GPUs: gpuKinds})
	}
	return result
}

//-------------------------------------------------------------------------------------------------
// METHODS
//-------------------------------------------------------------------------------------------------

// GetSubnet returns the full path to the subnetwork
func (c *Client) GetSubnetwork(zone string) (string, error) {
	info, ok := c.zones[zone]
	if !ok {
		return "", fmt.Errorf("zone %q is not available", zone)
	}
	return info.Subnetwork, nil
}

// GetAccelerator returns an available accelerator with the specified GPU kind in the provided
// zone if there exists such an accelerator.
func (c *Client) GetAccelerator(zone string, kind typedefs.GPUKind) (Accelerator, error) {
	info, ok := c.zones[zone]
	if !ok {
		return Accelerator{}, fmt.Errorf("zone %q is not available", zone)
	}

	for _, accelerator := range info.Accelerators {
		if accelerator.kind == kind {
			return accelerator, nil
		}
	}
	return Accelerator{}, fmt.Errorf("GPU kind %q is not available in zone %q", kind, zone)
}

//-------------------------------------------------------------------------------------------------
// HELPERS
//-------------------------------------------------------------------------------------------------

func (c *Client) updateZonesPeriodically(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Hour * 24):
			info, err := fetchZoneInfo(ctx, c.clients, c.project, c.network)
			if err != nil {
				zeus.Logger(ctx).Error("failed to update zone information", zap.Error(err))
				continue
			}

			// Actually update info
			c.mutex.Lock()
			c.zones = info
			c.mutex.Unlock()
		}
	}
}

func fetchZoneInfo(
	ctx context.Context, clients clients, project, network string,
) (map[string]ZoneInfo, error) {
	// Fetch all components
	zones, err := fetchZonesAndSubnetworks(ctx, clients.Zones, clients.Networks, project, network)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available zones: %s", err)
	}

	accelerators, err := fetchAccelerators(
		ctx, clients.AcceleratorTypes, project, jack.MapKeys(zones),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available accelerators: %s", err)
	}

	// Assemble into zone info
	result := make(map[string]ZoneInfo)
	for zone, subnetwork := range zones {
		result[zone] = ZoneInfo{
			Accelerators: accelerators[zone],
			Subnetwork:   subnetwork,
		}
	}
	return result, nil
}
