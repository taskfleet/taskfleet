package gcpzones

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestFetchZonesAndSubnetworks(t *testing.T) {
	ctx := context.Background()

	// Correctly fetch zones and subnetworks
	zones := []string{"europe-west3-b", "europe-west3-c", "europe-north1-a"}
	zonesClient := newProjectZonesClient(ctx, t, zones)
	networksClient := newNetworksClient(ctx, t, []string{
		"regions/europe-west3/subnetworks/subnet-1",
		"regions/europe-north1/subnetworks/subnet-2",
	})
	mapping, _ := fetchZonesAndSubnetworks(ctx, zonesClient, networksClient, "", "")
	assert.ElementsMatch(t, jack.MapKeys(mapping), zones)
	assert.Equal(t, mapping["europe-west3-b"], "regions/europe-west3/subnetworks/subnet-1")
	assert.Equal(t, mapping["europe-west3-c"], "regions/europe-west3/subnetworks/subnet-1")
	assert.Equal(t, mapping["europe-north1-a"], "regions/europe-north1/subnetworks/subnet-2")

	// Check that a zone is not returned if there is not network
	networksClient = newNetworksClient(ctx, t, []string{
		"regions/europe-north1/subnetworks/subnet-2",
	})
	mapping, _ = fetchZonesAndSubnetworks(ctx, zonesClient, networksClient, "", "")
	assert.ElementsMatch(t, jack.MapKeys(mapping), []string{"europe-north1-a"})
	assert.Equal(t, mapping["europe-north1-a"], "regions/europe-north1/subnetworks/subnet-2")

	// Fail to fetch when there is more than one subnetwork per zone
	networksClient = newNetworksClient(ctx, t, []string{
		"regions/europe-west3/subnetworks/subnet1",
		"regions/europe-west3/subnetworks/subnet2",
	})
	_, err := fetchZonesAndSubnetworks(ctx, zonesClient, networksClient, "", "")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "duplicate subnetwork")

	// Fail if the networks client is broken
	networksClient = newNetworksClient(ctx, t, nil)
	_, err = fetchZonesAndSubnetworks(ctx, zonesClient, networksClient, "", "")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to get GCP network")

	// Fail if the zone client is broken
	zonesClient = newProjectZonesClient(ctx, t, nil)
	networksClient = newNetworksClient(ctx, t, []string{})
	_, err = fetchZonesAndSubnetworks(ctx, zonesClient, networksClient, "", "")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "failed to list zones")
}

//-------------------------------------------------------------------------------------------------

func newProjectZonesClient(
	ctx context.Context, t *testing.T, zones []string,
) *compute.ZonesClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if zones == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response := &computepb.ZoneList{
			Items: jack.SliceMap(zones, func(zone string) *computepb.Zone {
				return &computepb.Zone{
					Name: proto.String(zone),
					Region: proto.String(
						fmt.Sprintf("https://example.com/regions/%s",
							gcputils.RegionFromZone(zone)),
					),
				}
			}),
		}
		result := jack.Must(protojson.Marshal(response))
		jack.Must(w.Write(result))
	}))
	service, err := compute.NewZonesRESTClient(
		ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication(),
	)
	require.Nil(t, err)
	return service
}

func newNetworksClient(
	ctx context.Context, t *testing.T, subnetworks []string,
) *compute.NetworksClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if subnetworks == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		response := &computepb.Network{Subnetworks: subnetworks}
		result := jack.Must(protojson.Marshal(response))
		jack.Must(w.Write(result))
	}))
	service, err := compute.NewNetworksRESTClient(
		ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication(),
	)
	require.Nil(t, err)
	return service
}
