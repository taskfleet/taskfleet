//go:build integration

package awsinstances

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/eagle"
	awszones "go.taskfleet.io/services/genesis/internal/providers/impl/aws/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

func TestClientCreateDelete(t *testing.T) {
	ctx := context.Background()

	// Initialize client
	cfg, err := config.LoadDefaultConfig(ctx)
	require.Nil(t, err)
	ec2Client := ec2.NewFromConfig(cfg)

	// Read configuration
	var awsConfig template.AwsConfig
	err = eagle.LoadConfig(&awsConfig, eagle.WithYAMLFile("testdata/config.env.yaml", false))
	require.Nil(t, err)

	// Initialize all dependencies (without mocking)
	vpcs, err := awszones.FindVPCs(ctx, ec2Client, awsConfig.Network)
	require.Nil(t, err)

	managers, err := FindAvailableInstances(ctx, vpcs)
	require.Nil(t, err)

	// Obtain a random zone from the managers
	var zone string
	for z := range managers {
		zone = z
		break
	}

	// Create the client and launch an instance
	client := NewClient(ctx, vpcs, managers, awsConfig)

	instanceType, err := managers[zone].FindBestFit(
		instances.Resources{CPUCount: 1, MemoryMiB: 2048}, typedefs.ArchitectureX86,
	)
	require.Nil(t, err)

	id := uuid.New()
	promise, err := client.Create(ctx, providers.InstanceRef{
		ID:   id,
		Zone: zone,
	}, providers.InstanceSpec{
		InstanceType: instanceType,
		IsSpot:       false,
	})
	require.Nil(t, err)

	// Await the instance creation
	instance, err := promise.Await(ctx)
	assert.Equal(t, id, instance.Ref.ID)

	// Delete the instance
	err = client.Delete(ctx, instance.Ref)
	require.Nil(t, err)
}
