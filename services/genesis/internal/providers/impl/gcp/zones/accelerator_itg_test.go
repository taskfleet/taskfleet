//go:build integration

package gcpzones

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
)

func TestFetchAcceleratorsGcp(t *testing.T) {
	ctx := context.Background()
	clients := gcputils.NewClientFactory(ctx)

	accelerators, err := fetchAccelerators(
		ctx, clients.AcceleratorTypes(), gcpProject, []string{"us-central1-a", "europe-west3-b"},
	)
	assert.Nil(t, err)
	assert.Len(t, accelerators, 2)
	// At least 3 accelerators in us-central1-a
	assert.GreaterOrEqual(t, len(accelerators["us-central1-a"]), 3)
	// At least 1 accelerator in europe-west3-b
	assert.GreaterOrEqual(t, len(accelerators["europe-west3-b"]), 1)
}
