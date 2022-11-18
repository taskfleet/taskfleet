package awszones

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestClientList(t *testing.T) {
	ctx := context.Background()
	instances := map[string]*instances.Manager{
		"us-east-1a": instances.NewManager([]instances.Type{
			{Name: "cpu-1"},
			{Name: "gpu-1", Resources: instances.Resources{GPU: &instances.GPUResources{
				Kind: typedefs.GPUNvidiaTeslaK80,
			}}},
		}),
		"us-east-1b": instances.NewManager([]instances.Type{
			{Name: "cpu-1"},
			{Name: "gpu-1", Resources: instances.Resources{GPU: &instances.GPUResources{
				Kind: typedefs.GPUNvidiaTeslaK80,
			}}},
			{Name: "gpu-2", Resources: instances.Resources{GPU: &instances.GPUResources{
				Kind: typedefs.GPUNvidiaTeslaM60,
			}}},
		}),
	}

	client := NewClient(ctx, instances)
	zones := client.List()
	assert.Len(t, zones, 2)
	for _, zone := range zones {
		switch zone.Name {
		case "us-east-1a":
			assert.ElementsMatch(t, zone.GPUs, []typedefs.GPUKind{typedefs.GPUNvidiaTeslaK80})
		case "us-east-1b":
			assert.ElementsMatch(t, zone.GPUs, []typedefs.GPUKind{
				typedefs.GPUNvidiaTeslaK80, typedefs.GPUNvidiaTeslaM60,
			})
		default:
			t.Fail()
		}
	}
}
