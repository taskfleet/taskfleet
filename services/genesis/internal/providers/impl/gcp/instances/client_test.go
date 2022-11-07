package gcpinstances

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
)

func TestCreateShutdown(t *testing.T) {
	if testing.Short() {
		return
	}

	ctx := context.Background()
	service, err := compute.NewService(ctx)
	require.Nil(t, err)
	zonesClient, err := gcpzones.NewClient(ctx, service, os.Getenv("GCP_PROJECT"))
	require.Nil(t, err)

	client, err := NewClient(
		ctx,
		os.Getenv("GCP_IDENTIFIER"),
		os.Getenv("GCP_NETWORK"),
		os.Getenv("GCP_PROJECT"),
		service,
		zonesClient,
	)
	require.Nil(t, err)

	// First, create the instance
	creationContext, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// Create an ID
	id, err := uuid.NewRandom()
	require.Nil(t, err)

	// Try to start the instance in a fixed zone
	zone := "us-central1-a"

	// Get a suitable machine type
	instanceType, err := client.Find(zone, instances.Resources{
		CPUCount:  1,
		MemoryMiB: 2000,
	})
	require.Nil(t, err)

	// Check if the returned machine type is sensible
	assert.Equal(t, "n1-standard-1", instanceType.Name)
	assert.Equal(t, uint16(1), instanceType.CPUCount)
	assert.Equal(t, uint32(3840), instanceType.MemoryMiB)

	// Then create the instance specification
	meta := providers.InstanceRef{
		ID:   id,
		Zone: zone,
	}
	spec := providers.InstanceSpec{
		Compute: providers.ComputeConfig{
			InstanceType: instanceType,
			IsSpot:       false,
		},
		Boot: providers.BootConfig{
			ImageLink:   "projects/debian-cloud/global/images/family/debian-10",
			DiskSizeGiB: 10,
		},
		Metadata: providers.MetadataConfig{
			Labels: map[string]string{
				LabelKeyOwnedBy: "genesis-test",
			},
		},
	}

	promise, err := client.Create(creationContext, meta, spec)
	require.Nil(t, err)

	// Then wait for the instance to be running
	instance, err := promise.Await(creationContext)
	require.Nil(t, err)

	assert.True(t, instance.Status.CreationTimestamp.Before(time.Now()))
	assert.True(t, instance.Status.CreationTimestamp.After(time.Now().Add(-2*time.Minute)))

	hostname := fmt.Sprintf("%s.%s.c.%s.internal",
		instance.Meta.CommonName(), instance.Meta.ProviderZone, client.ProjectID,
	)
	assert.Equal(t, hostname, instance.Status.Network.InternalHostname)

	// Eventually, purge the instance
	deletionContext, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	err = client.Delete(deletionContext, instance.Ref)
	require.Nil(t, err)
}
