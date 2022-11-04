//go:build integration

package gcpzones

import (
	"context"
	"testing"

	compute "cloud.google.com/go/compute/apiv1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchAcceleratorsGcp(t *testing.T) {
	ctx := context.Background()
	client, err := compute.NewAcceleratorTypesRESTClient(ctx)
	require.Nil(t, err)

	accelerators, err := fetchAccelerators(
		ctx, client, gcpProject, []string{"us-central1-a", "europe-west3-b"},
	)
	assert.Nil(t, err)
	assert.Len(t, accelerators, 2)
	// At least 3 accelerators in us-central1-a
	assert.GreaterOrEqual(t, len(accelerators["us-central1-a"]), 3)
	// At least 1 accelerator in europe-west3-b
	assert.GreaterOrEqual(t, len(accelerators["europe-west3-b"]), 1)
}
