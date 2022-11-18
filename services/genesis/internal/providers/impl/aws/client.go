package aws

import (
	"context"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	awsinstances "go.taskfleet.io/services/genesis/internal/providers/impl/aws/instances"
	awszones "go.taskfleet.io/services/genesis/internal/providers/impl/aws/zones"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

type client struct {
	account   string
	zones     *awszones.Client
	instances *awsinstances.Client
}

func NewClient(ctx context.Context, config template.AwsConfig) (providers.Provider, error) {
	// First, we initialize an EC2 client
	cfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, providers.NewClientError("invalid configuration", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)

	// As a next step, we find the account ID
	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, providers.NewAPIError("failed to get caller identity", err)
	}

	// Then, find all VPCs to launch instances into...
	vpcs, err := awszones.FindVPCs(ctx, ec2Client, config.Network)
	if err != nil {
		return nil, providers.NewAPIError("failed to find usable VPCs", err)
	}

	// ...and get all the instances that can be launched
	instances, err := awsinstances.FindAvailableInstances(ctx, vpcs)
	if err != nil {
		return nil, providers.NewAPIError("failed to find all available instance types", err)
	}

	// Then, we can create the internal clients and return the newly created client
	zoneClient := awszones.NewClient(ctx, instances)
	return &client{account: *identity.Account, zones: zoneClient}, nil
}

func (c *client) CloudProvider() typedefs.CloudProvider {
	return typedefs.ProviderAmazonWebServices
}

func (c *client) AccountName() string {
	return c.account
}

func (c *client) Zones() providers.ZoneClient {
	return c.zones
}

func (c *client) Instances() providers.InstanceClient {
	return c.instances
}
