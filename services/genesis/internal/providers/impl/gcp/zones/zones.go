package gcpzones

import (
	"context"
	"fmt"
	"path"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
)

// fetchZonesAndSubnetworks fetches all subnetworks from the given network and returns a mapping
// from zones to subnetworks.
func fetchZonesAndSubnetworks(
	ctx context.Context,
	zones *compute.ZonesClient,
	networks *compute.NetworksClient,
	project string,
	network string,
) (map[string]string, error) {
	// Find all subnetworks of the provided network
	regions := map[string]string{}
	gcpNetwork, err := networks.Get(ctx, &computepb.GetNetworkRequest{
		Project: project,
		Network: network,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get GCP network %q: %s", network, err)
	}
	for _, subnet := range gcpNetwork.Subnetworks {
		splits := strings.Split(subnet, "/")
		region := splits[len(splits)-3]
		if _, ok := regions[region]; ok {
			return nil, fmt.Errorf(
				"duplicate subnetwork for region %q, please remove one subnet to use this network",
				region,
			)
		}
		regions[region] = subnet
	}

	// Find all zones for which there exists a subnetwork
	result := map[string]string{}
	it := zones.List(ctx, &computepb.ListZonesRequest{Project: project})
	if err := gcputils.Iterate[*computepb.Zone](ctx, it, func(zone *computepb.Zone) error {
		if subnet, ok := regions[path.Base(*zone.Region)]; ok {
			result[*zone.Name] = subnet
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to list zones in project %q: %s", project, err)
	}
	return result, nil
}
