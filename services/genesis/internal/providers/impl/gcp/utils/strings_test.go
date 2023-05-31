package gcputils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegionFromZone(t *testing.T) {
	assert.Equal(t, RegionFromZone("europe-west3-b"), "europe-west3")
}
