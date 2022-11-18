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
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"google.golang.org/api/option"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func TestAcceleratorMaxCount(t *testing.T) {
	a := Accelerator{MaxCountPerInstance: 4}
	assert.EqualValues(t, a.MaxCount(), 4)
}

func TestAcceleratorConfig(t *testing.T) {
	a := Accelerator{MaxCountPerInstance: 4}

	config, err := a.Config(1)
	assert.Nil(t, err)
	assert.EqualValues(t, *config.AcceleratorCount, 1)

	config, err = a.Config(3)
	assert.Nil(t, err)
	assert.EqualValues(t, *config.AcceleratorCount, 4)

	config, err = a.Config(4)
	assert.Nil(t, err)
	assert.EqualValues(t, *config.AcceleratorCount, 4)

	_, err = a.Config(8)
	assert.NotNil(t, err)
}

func TestFetchAccelerators(t *testing.T) {
	ctx := context.Background()

	// Parse two zones
	client := newAcceleratorTypesClient(ctx, t, map[string][]string{
		"zone-1": {"nvidia-tesla-v100", "nvidia-tesla-p100"},
		"zone-2": {"nvidia-tesla-k80"},
	})
	accelerators, _ := fetchAccelerators(ctx, client, "", []string{"zone-1", "zone-2"})
	assert.Len(t, accelerators, 2)
	assert.Len(t, accelerators["zone-1"], 2)
	assert.Len(t, accelerators["zone-2"], 1)
	assert.Equal(t, accelerators["zone-1"][0].Kind, typedefs.GPUNvidiaTeslaV100)
	assert.Equal(t, accelerators["zone-1"][1].Kind, typedefs.GPUNvidiaTeslaP100)
	assert.Equal(t, accelerators["zone-2"][0].Kind, typedefs.GPUNvidiaTeslaK80)

	// Parse only one zone
	accelerators, _ = fetchAccelerators(ctx, client, "", []string{"zone-1"})
	assert.Len(t, accelerators, 1)
	assert.Len(t, accelerators["zone-1"], 2)

	// Don't return unknown GPUs
	client = newAcceleratorTypesClient(ctx, t, map[string][]string{
		"zone-1": {"unknown-gpu"},
		"zone-2": {"nvidia-tesla-k80", "unknown-gpu"},
	})
	accelerators, _ = fetchAccelerators(ctx, client, "", []string{"zone-1", "zone-2"})
	assert.Len(t, accelerators, 2)
	assert.Len(t, accelerators["zone-1"], 0)
	assert.Len(t, accelerators["zone-2"], 1)

	// Ignore virtual workstations
	client = newAcceleratorTypesClient(ctx, t, map[string][]string{
		"zone-1": {"nvidia-tesla-v100-vws"},
	})
	accelerators, _ = fetchAccelerators(ctx, client, "", []string{"zone-1"})
	assert.Len(t, accelerators, 1)
	assert.Len(t, accelerators["zone-1"], 0)

	// Test error
	client = newAcceleratorTypesClient(ctx, t, nil)
	_, err := fetchAccelerators(ctx, client, "", []string{"zone-1"})
	assert.NotNil(t, err)
}

//-------------------------------------------------------------------------------------------------

func newAcceleratorTypesClient(
	ctx context.Context, t *testing.T, accelerators map[string][]string,
) *compute.AcceleratorTypesClient {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if accelerators == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		items := map[string]*computepb.AcceleratorTypesScopedList{}
		for zone, names := range accelerators {
			types := []*computepb.AcceleratorType{}
			for _, name := range names {
				types = append(types, &computepb.AcceleratorType{
					Name:     proto.String(name),
					SelfLink: proto.String(fmt.Sprintf("https://example.com/%s", name)),
					// Keeping this line empty for nicer formatting
					MaximumCardsPerInstance: proto.Int32(4),
				})
			}
			items[zone] = &computepb.AcceleratorTypesScopedList{AcceleratorTypes: types}
		}
		response := &computepb.AcceleratorTypeAggregatedList{Items: items}
		result := jack.Must(protojson.Marshal(response))
		jack.Must(w.Write(result))
	}))
	service, err := compute.NewAcceleratorTypesRESTClient(
		ctx, option.WithEndpoint(server.URL), option.WithoutAuthentication(),
	)
	require.Nil(t, err)
	return service
}
