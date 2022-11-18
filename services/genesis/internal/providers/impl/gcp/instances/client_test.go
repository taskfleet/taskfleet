package gcpinstances

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.taskfleet.io/packages/jack"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestClientFind(t *testing.T) {
	client := Client{
		reservationsHelper: jack.Must(newReservationsHelper(template.InstanceReservations{})),
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

	testCases := []struct {
		zone         string
		resources    instances.Resources
		architecture typedefs.CPUArchitecture
		expected     string
		err          *string
	}{
		{
			zone:         "zone-1",
			resources:    instances.Resources{CPUCount: 1, MemoryMiB: 6144},
			architecture: typedefs.ArchitectureX86,
			expected:     "instance-2",
		},
		{
			zone:         "zone-1",
			resources:    instances.Resources{CPUCount: 3, MemoryMiB: 6144},
			architecture: typedefs.ArchitectureX86,
			err:          jack.Ptr("could not find"),
		},
		{
			zone: "zone-3",
			err:  jack.Ptr("no instances"),
		},
	}

	for _, testCase := range testCases {
		instance, err := client.Find(testCase.zone, testCase.resources, testCase.architecture)
		if testCase.err != nil {
			assert.NotNil(t, err)
			assert.ErrorContains(t, err, *testCase.err)
		} else {
			assert.Equal(t, testCase.expected, instance.Name)
		}
	}
}

// func TestCreateShutdown(t *testing.T) {
// 	if testing.Short() {
// 		return
// 	}

// 	ctx := context.Background()
// 	service, err := compute.NewService(ctx)
// 	require.Nil(t, err)
// 	zonesClient, err := gcpzones.NewClient(ctx, service, os.Getenv("GCP_PROJECT"))
// 	require.Nil(t, err)

// 	client, err := NewClient(
// 		ctx,
// 		os.Getenv("GCP_IDENTIFIER"),
// 		os.Getenv("GCP_NETWORK"),
// 		os.Getenv("GCP_PROJECT"),
// 		service,
// 		zonesClient,
// 	)
// 	require.Nil(t, err)

// 	// First, create the instance
// 	creationContext, cancel := context.WithTimeout(ctx, 2*time.Minute)
// 	defer cancel()

// 	// Create an ID
// 	id, err := uuid.NewRandom()
// 	require.Nil(t, err)

// 	// Try to start the instance in a fixed zone
// 	zone := "us-central1-a"

// 	// Get a suitable machine type
// 	instanceType, err := client.Find(zone, instances.Resources{
// 		CPUCount:  1,
// 		MemoryMiB: 2000,
// 	})
// 	require.Nil(t, err)

// 	// Check if the returned machine type is sensible
// 	assert.Equal(t, "n1-standard-1", instanceType.Name)
// 	assert.Equal(t, uint16(1), instanceType.CPUCount)
// 	assert.Equal(t, uint32(3840), instanceType.MemoryMiB)

// 	// Then create the instance specification
// 	meta := providers.InstanceRef{
// 		ID:   id,
// 		Zone: zone,
// 	}
// 	spec := providers.InstanceSpec{
// 		Compute: providers.ComputeConfig{
// 			InstanceType: instanceType,
// 			IsSpot:       false,
// 		},
// 		Boot: providers.BootConfig{
// 			ImageLink:   "projects/debian-cloud/global/images/family/debian-10",
// 			DiskSizeGiB: 10,
// 		},
// 		Metadata: providers.MetadataConfig{
// 			Labels: map[string]string{
// 				LabelKeyOwnedBy: "genesis-test",
// 			},
// 		},
// 	}

// 	promise, err := client.Create(creationContext, meta, spec)
// 	require.Nil(t, err)

// 	// Then wait for the instance to be running
// 	instance, err := promise.Await(creationContext)
// 	require.Nil(t, err)

// 	assert.True(t, instance.Status.CreationTimestamp.Before(time.Now()))
// 	assert.True(t, instance.Status.CreationTimestamp.After(time.Now().Add(-2*time.Minute)))

// 	hostname := fmt.Sprintf("%s.%s.c.%s.internal",
// 		instance.Meta.CommonName(), instance.Meta.ProviderZone, client.ProjectID,
// 	)
// 	assert.Equal(t, hostname, instance.Status.Network.InternalHostname)

// 	// Eventually, purge the instance
// 	deletionContext, cancel := context.WithTimeout(ctx, 5*time.Minute)
// 	defer cancel()

// 	err = client.Delete(deletionContext, instance.Ref)
// 	require.Nil(t, err)
// }
