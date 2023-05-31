//go:build integration

package gcpinstances

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/tftest"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestFindCreateListDelete(t *testing.T) {
	testCases := []struct {
		resources instances.Resources
		instance  string
	}{
		{
			resources: instances.Resources{
				CPUCount:  1,
				MemoryMiB: 3500,
			},
			instance: "n1-standard-1",
		},
		{
			resources: instances.Resources{
				CPUCount:  1,
				MemoryMiB: 3500,
				GPU:       &instances.GPUResources{Kind: typedefs.GPUNvidiaTeslaK80, Count: 1},
			},
			instance: "n1-standard-1-nvidia-tesla-k80-1",
		},
	}

	ctx := context.Background()

	// Set up Terraform
	tf := tftest.Setup(ctx, t, "../_testdata/terraform",
		"create_iam=true",
		fmt.Sprintf("gcp_project=%s", gcpProject),
	)

	// Get the network name from Terraform
	networkName := tftest.GetOutput[string](ctx, t, tf, "network_name")

	// Initialize dependencies
	id := fmt.Sprintf("test-%s", uuid.NewString()[:8])
	config := template.GcpConfig{
		CommonInstanceConfig: template.CommonInstanceConfig{
			Reservations: template.InstanceReservations{},
			ExtraDisks:   []template.InstanceDisk{},
			Metadata:     map[string]string{},
		},
		GcpInstanceConfig: template.GcpInstanceConfig{
			Boot: template.GcpBootConfig{
				ImageLink: []template.Option[string]{
					{Config: "https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/ubuntu-2204-jammy-v20221101a"},
				},
				DiskSize: "10Gi",
			},
			Network: template.GcpNetworkConfig{
				Name: networkName,
			},
			Iam: template.GcpIamConfig{},
			Disks: template.GcpDiskConfig{
				Type: "pd-balanced",
			},
		},
	}
	clients := gcputils.NewClientFactory(ctx)
	zones, err := gcpzones.NewClient(ctx, clients, gcpProject, config.Network.Name)
	require.Nil(t, err)

	// Set up client
	client, err := NewClient(ctx, id, gcpProject, config, clients, zones)
	require.Nil(t, err)

	// Test Group I: Find and create instances in parallel
	torig := t
	t.Run("create", func(t *testing.T) {
		for i, tc := range testCases {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				t.Parallel()

				// Find instance
				instanceType, err := client.Find(
					"us-east1-c", tc.resources, typedefs.ArchitectureX86,
				)
				require.Nil(t, err)
				assert.Equal(t, instanceType.Name, tc.instance)

				// Create instance
				instanceID := uuid.New()
				promise, err := client.Create(ctx, providers.InstanceMeta{
					ID:           instanceID,
					ProviderZone: "us-east1-c",
				}, providers.InstanceSpec{
					InstanceType: instanceType,
				})
				require.Nil(t, err)

				// Await instance
				instance, err := promise.Await(ctx)
				require.Nil(t, err)
				assert.Equal(t, instanceID, instance.Meta.ID)
				assert.Equal(t, fmt.Sprintf("taskfleet-%s", instanceID), instance.Meta.ProviderID)
				assert.Equal(t, instanceType, instance.Spec.InstanceType)

				// Ensure that instance is deleted at the end of the test function
				torig.Cleanup(func() {
					err := client.Delete(ctx, instance.Meta)
					assert.Nil(t, err)
				})
			})
		}
	})

	// List all instances and check that all are returned
	instances, err := client.List(ctx)
	assert.Nil(t, err)
	require.Len(t, instances, len(testCases))
}
