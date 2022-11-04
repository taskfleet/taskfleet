//go:build integration

package gcpzones

import (
	"context"
	"net/url"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	"go.taskfleet.io/services/genesis/internal/tftest"
)

func TestFetchZonesandSubnetworksGcp(t *testing.T) {
	ctx := context.Background()

	// Set up Terraform
	t.Setenv("GOOGLE_PROJECT", gcpProject)
	tf := tftest.Setup(ctx, t, "testdata/network")

	// Get the network name from Terraform
	networkName := tftest.GetOutput[string](ctx, t, tf, "network_name")

	// Check that the results are as expected
	zoneClient, err := compute.NewZonesRESTClient(ctx)
	require.Nil(t, err)
	networkClient, err := compute.NewNetworksRESTClient(ctx)
	require.Nil(t, err)

	zones, err := fetchZonesAndSubnetworks(ctx, zoneClient, networkClient, gcpProject, networkName)
	assert.Nil(t, err)

	// There should be at least 3 zones for each of the regions
	assert.GreaterOrEqual(t, len(zones), 6)

	// There should only be entries for our two regions
	for zone, subnet := range zones {
		assert.Contains(
			t, []string{"europe-west3", "europe-north1"}, gcputils.RegionFromZone(zone),
		)
		// Also check that values (= subnetworks) are full URIs
		_, err := url.ParseRequestURI(subnet)
		assert.Nil(t, err)
	}
}
