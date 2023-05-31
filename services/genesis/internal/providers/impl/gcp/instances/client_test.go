package gcpinstances

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestClientFind(t *testing.T) {
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

	ctx := context.Background()
	client := testClient(ctx, t, template.GcpConfig{}, nil, nil)

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			instance, err := client.Find(tc.zone, tc.resources, tc.architecture)
			if tc.err != nil {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, *tc.err)
			} else {
				assert.Equal(t, tc.expected, instance.Name)
			}
		})
	}
}

func TestClientCreate(t *testing.T) {
	testCases := []struct {
		config template.GcpConfig
		meta   providers.InstanceMeta
		spec   providers.InstanceSpec
		err    *string
	}{
		{
			config: template.GcpConfig{
				GcpInstanceConfig: template.GcpInstanceConfig{
					Iam: template.GcpIamConfig{ServiceAccountEmail: "hi@example.com"},
				},
			},
			meta: providers.InstanceMeta{
				ID:           uuid.New(),
				ProviderZone: "zone-1",
			},
			spec: providers.InstanceSpec{
				InstanceType: instances.Type{
					Name: "n1-standard-1",
					UID:  "http://example.com/n1-standard-1",
				},
			},
		},
	}

	ctx := context.Background()
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			// Create zones mock client
			zones := gcpzones.NewMockClient(t)
			zones.EXPECT().GetSubnetwork(mock.Anything).Return("my-subnetwork", nil)
			zones.EXPECT().GetAccelerator(mock.Anything, typedefs.GPUNvidiaTeslaK80).Return(
				gcpzones.Accelerator{
					URI:                 "https://nvidia-k80",
					Kind:                typedefs.GPUNvidiaTeslaK80,
					MaxCountPerInstance: 4,
				},
				nil,
			).Maybe()

			// Create handle function
			handle := func(w http.ResponseWriter, r *http.Request) {
				result := jack.Must(protojson.Marshal(&computepb.Operation{}))
				jack.Must(w.Write(result))
			}

			// Create client and run instance creation
			client := testClient(ctx, t, tc.config, zones, handle)
			promise, err := client.Create(ctx, tc.meta, tc.spec)
			if tc.err != nil {
				assert.NotNil(t, err)
				assert.ErrorContains(t, err, *tc.err)
			} else {
				require.Nil(t, err)
				instancePromise := promise.(*instancePromise)
				assert.Equal(t, tc.meta.ID, instancePromise.meta.ID)
				assert.Equal(
					t, fmt.Sprintf("taskfleet-%s", tc.meta.ID), instancePromise.meta.ProviderID,
				)
			}
		})
	}
}
