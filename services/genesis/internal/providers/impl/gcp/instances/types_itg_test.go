//go:build integration

package gcpinstances

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/tftest"
)

func TestFindAvailableInstanceTypesGcp(t *testing.T) {
	ctx := context.Background()

	// Set up Terraform
	tf := tftest.Setup(ctx, t, "../_testdata/terraform",
		fmt.Sprintf("gcp_project=%s", gcpProject),
	)

	// Get the network name from Terraform
	networkName := tftest.GetOutput[string](ctx, t, tf, "network_name")

	// Initialize the zones client
	clients := gcputils.NewClientFactory(ctx)
	zones, err := gcpzones.NewClient(ctx, clients, gcpProject, networkName)
	require.Nil(t, err)

	// Fetch available instance types
	types, err := findAvailableInstanceTypes(ctx, clients.MachineTypes(), zones, gcpProject)
	require.Nil(t, err)

	// Check that types prove assumptions...
	// Each zone should have at least 100 instance types.
	for zone, zoneTypes := range types {
		assert.GreaterOrEqualf(
			t, len(zoneTypes), 100,
			"zone %q has less than 100 instance types", zone,
		)
	}

	// We should find N1 GPU instances, A2 GPU instances and no other instances with GPUs
	hasN1Gpu := false
	hasA2 := false
	hasNoOtherWithGpu := true
	for _, zoneTypes := range types {
		for _, t := range zoneTypes {
			if strings.HasPrefix(t.Name, "n1-") && t.Resources.GPU != nil {
				hasN1Gpu = true
			} else if strings.HasPrefix(t.Name, "a2-") {
				hasA2 = true
			} else if t.Resources.GPU != nil {
				hasNoOtherWithGpu = false
			}
		}
	}
	assert.True(t, hasN1Gpu)
	assert.True(t, hasA2)
	assert.True(t, hasNoOtherWithGpu)
}
