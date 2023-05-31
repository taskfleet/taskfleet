//go:build integration

package gcpzones

import (
	"context"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	"go.taskfleet.io/services/genesis/internal/tftest"
)

func TestFetchZonesandSubnetworksGcp(t *testing.T) {
	ctx := context.Background()

	// Set up Terraform
	tf := tftest.Setup(ctx, t, "../_testdata/terraform",
		fmt.Sprintf("gcp_project=%s", gcpProject),
	)

	// Get the network name from Terraform
	networkName := tftest.GetOutput[string](ctx, t, tf, "network_name")

	// Check that the results are as expected
	clients := gcputils.NewClientFactory(ctx)
	zones, err := fetchZonesAndSubnetworks(ctx, clients, gcpProject, networkName)
	assert.Nil(t, err)

	// There should be at least 3 zones for each of the regions
	assert.GreaterOrEqual(t, len(zones), 6)

	// There should only be entries for our two regions
	for zone, subnet := range zones {
		assert.Contains(
			t, []string{"europe-west3", "us-east1"}, gcputils.RegionFromZone(zone),
		)
		// Also check that values (= subnetworks) are full URIs
		_, err := url.ParseRequestURI(subnet)
		assert.Nil(t, err)
	}
}
