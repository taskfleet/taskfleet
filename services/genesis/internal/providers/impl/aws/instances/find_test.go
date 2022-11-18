//go:build integration

package awsinstances

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindAvailableInstances(t *testing.T) {
	ctx := context.Background()
	vpcs := newVpcMocks(ctx, t)

	// Check instances found for three availability zones
	managers, err := FindAvailableInstances(ctx, vpcs)
	require.Nil(t, err)
	assert.Len(t, managers, 3)

	// Check whether a single availability zone provides many instances
	useast1a := managers["us-east-1a"]
	assert.Greater(t, len(useast1a.Types()), 100)

	// Check that not all zones provide the same instances
	assert.NotEqual(t, len(useast1a.Types()), len(managers["us-east-1b"].Types()))
	assert.NotEqual(t, len(useast1a.Types()), len(managers["us-east-2c"].Types()))
}
