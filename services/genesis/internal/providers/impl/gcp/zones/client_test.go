package gcpzones

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestZoneList(t *testing.T) {
	c := &client{zones: map[string]ZoneInfo{
		"zone-1": {Accelerators: []Accelerator{{Kind: typedefs.GPUNvidiaTeslaV100}}},
		"zone-2": {Accelerators: []Accelerator{
			{Kind: typedefs.GPUNvidiaTeslaV100},
			{Kind: typedefs.GPUNvidiaTeslaK80},
		}},
	}}

	expected := []providers.Zone{
		{Name: "zone-1", GPUs: []typedefs.GPUKind{typedefs.GPUNvidiaTeslaV100}},
		{Name: "zone-2", GPUs: []typedefs.GPUKind{
			typedefs.GPUNvidiaTeslaV100, typedefs.GPUNvidiaTeslaK80,
		}},
	}
	assert.ElementsMatch(t, expected, c.List())
}

func TestZoneGetSubnetwork(t *testing.T) {
	c := &client{zones: map[string]ZoneInfo{
		"zone-1": {Subnetwork: "subnet-1"},
		"zone-2": {Subnetwork: "subnet-2"},
	}}

	// Get existing subnet
	subnet, err := c.GetSubnetwork("zone-1")
	assert.Nil(t, err)
	assert.Equal(t, subnet, "subnet-1")

	// Fail on unknown subnet
	_, err = c.GetSubnetwork("unknown")
	assert.ErrorContains(t, err, "not available")
}

func TestGetAccelerator(t *testing.T) {
	c := &client{zones: map[string]ZoneInfo{
		"zone-1": {Accelerators: []Accelerator{{Kind: typedefs.GPUNvidiaTeslaV100}}},
		"zone-2": {Accelerators: []Accelerator{
			{Kind: typedefs.GPUNvidiaTeslaV100},
			{Kind: typedefs.GPUNvidiaTeslaK80},
		}},
	}}

	// Get existing accelerator
	accelerator, err := c.GetAccelerator("zone-1", typedefs.GPUNvidiaTeslaV100)
	assert.Nil(t, err)
	assert.Equal(t, accelerator.Kind, typedefs.GPUNvidiaTeslaV100)

	// Fail on unknown accelerator
	_, err = c.GetAccelerator("zone-1", typedefs.GPUNvidiaTeslaA100)
	assert.ErrorContains(t, err, "GPU kind")
	assert.ErrorContains(t, err, "not available")

	// Fail on unkonwn zone
	_, err = c.GetAccelerator("unknown", typedefs.GPUNvidiaTeslaA100)
	assert.ErrorContains(t, err, "zone \"unknown\" is not available")
}

func TestFetchZoneInfo(t *testing.T) {
	ctx := context.Background()

	// Setup clients
	clients := gcputils.NewMockClientFactory(t)
	clients.EXPECT().AcceleratorTypes().Return(
		newAcceleratorTypesClient(ctx, t, map[string][]string{
			"region-1-a": {"nvidia-tesla-v100", "nvidia-tesla-p100"},
			"region-2-a": {"nvidia-tesla-k80"},
		}),
	)
	clients.EXPECT().Networks().Return(newNetworksClient(ctx, t, []string{
		"regions/region-1/subnetwork/subnet1",
		"regions/region-2/subnetwork/subnet2",
	}))
	clients.EXPECT().Zones().Return(newProjectZonesClient(ctx, t, []string{
		"region-1-a", "region-2-a", "region-3-a"},
	))

	// Fetch the zone info
	info, err := fetchZoneInfo(ctx, clients, "", "")
	assert.Nil(t, err)
	assert.Len(t, info, 2)

	expected := map[string]ZoneInfo{
		"region-1-a": {
			Accelerators: []Accelerator{
				{
					URI:                 "https://example.com/nvidia-tesla-v100",
					Kind:                typedefs.GPUNvidiaTeslaV100,
					MaxCountPerInstance: 4,
				},
				{
					URI:                 "https://example.com/nvidia-tesla-p100",
					Kind:                typedefs.GPUNvidiaTeslaP100,
					MaxCountPerInstance: 4,
				},
			},
			Subnetwork: "regions/region-1/subnetwork/subnet1",
		},
		"region-2-a": {
			Accelerators: []Accelerator{
				{
					URI:                 "https://example.com/nvidia-tesla-k80",
					Kind:                typedefs.GPUNvidiaTeslaK80,
					MaxCountPerInstance: 4,
				},
			},
			Subnetwork: "regions/region-2/subnetwork/subnet2",
		},
	}
	assert.True(t, reflect.DeepEqual(expected, info))
}
