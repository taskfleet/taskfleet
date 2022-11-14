//go:build integration

package gcpinstances

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	"go.taskfleet.io/services/genesis/internal/template"
)

func TestNewDisksHelperGcp(t *testing.T) {
	ctx := context.Background()
	clients := gcputils.NewClientFactory(ctx)
	helper, err := newDisksHelper(
		ctx,
		gcpProject,
		template.GcpBootConfig{DiskSize: "10Gi"},
		[]template.InstanceDisk{},
		"pd-balanced",
		clients.DiskTypes(),
	)
	assert.Nil(t, err)

	for zone, link := range helper.diskTypeSelfLinks {
		assert.Truef(t, strings.HasSuffix(
			link,
			fmt.Sprintf("projects/%s/zones/%s/diskTypes/pd-balanced", gcpProject, zone),
		), "disk type %q for zone %q has invalid suffix", link, zone)
	}
}
