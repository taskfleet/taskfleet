//go:build integration

package awszones

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/eagle"
	"go.taskfleet.io/services/genesis/internal/template"
)

type testdata struct {
	Tags map[string]string `json:"tags"`
	Vpc  struct {
		ID             string            `json:"id"`
		SecurityGroups []string          `json:"securityGroups"`
		Subnets        map[string]string `json:"subnets"`
	} `json:"vpc"`
}

func TestFindAvailableZones(t *testing.T) {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	require.Nil(t, err)
	client := ec2.NewFromConfig(cfg)

	var network testdata
	err = eagle.LoadConfig(&network, eagle.WithJSONFile("testdata/network.env.json", false))
	require.Nil(t, err)

	// Find VPCs
	vpcs, err := FindVPCs(ctx, client, template.AwsNetworkConfig{
		VpcSelector:           network.Tags,
		SecurityGroupSelector: network.Tags,
	})
	require.Nil(t, err)
	require.Len(t, vpcs, 1)

	// Check VPC
	vpc := vpcs[0]
	assert.Equal(t, network.Vpc.ID, vpc.ID())

	securityGroups, err := vpc.SecurityGroups(ctx)
	require.Nil(t, err)
	assert.ElementsMatch(t, network.Vpc.SecurityGroups, securityGroups)

	subnets, err := vpc.Subnets(ctx)
	require.Nil(t, err)
	assert.EqualValues(t, network.Vpc.Subnets, subnets)

	for zone := range subnets {
		assert.True(t, strings.HasPrefix(zone, vpc.Region()))
	}
}
