package gcpinstances

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"google.golang.org/api/option"
)

func testClient(
	ctx context.Context,
	t *testing.T,
	config template.GcpConfig,
	zones gcpzones.Client,
	handle func(w http.ResponseWriter, r *http.Request),
) *Client {
	server := httptest.NewServer(http.HandlerFunc(handle))
	service, err := compute.NewInstancesRESTClient(
		ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication(),
	)
	require.Nil(t, err)

	return &Client{
		projectID:  "test-project",
		identifier: "test-identifier",
		config:     config,
		zones:      zones,
		service:    service,
		network:    "my-network",
		disksHelper: &disksHelper{
			bootDiskSizeGiB: 10,
			bootImages: []template.Option[string]{
				{Config: "boot-image-1"},
			},
			extraDisks: []disk{
				{name: "disk-1", sizePerCpuGiB: 5},
			},
			diskTypeSelfLinks: map[string]string{
				"zone-1": "zone-1/my-disk",
			},
		},
		reservationsHelper: &reservationsHelper{},
		instanceManagers: map[string]*instances.Manager{
			"zone-1": jack.Must(instances.NewManager([]instances.Type{
				{
					Name:         "instance-1",
					Resources:    instances.Resources{CPUCount: 1, MemoryMiB: 4096},
					Architecture: typedefs.ArchitectureX86,
				},
				{
					Name:         "instance-2",
					Resources:    instances.Resources{CPUCount: 2, MemoryMiB: 8192},
					Architecture: typedefs.ArchitectureX86,
				},
			})),
		},
	}
}
