package awsinstances

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	awszones "go.taskfleet.io/services/genesis/internal/providers/impl/aws/zones"
)

func newVpcMocks(ctx context.Context, t *testing.T) []awszones.VPC {
	// us-east-1
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	require.Nil(t, err)
	client := ec2.NewFromConfig(cfg)
	useast1 := &awszones.MockVPC{}
	useast1.On("Client").Return(client)
	useast1.On("Region").Return("us-east-1")
	useast1.On("Subnets", mock.Anything).Return(map[string]string{
		"us-east-1a": "", "us-east-1b": "",
	}, nil)

	// us-east-2
	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion("us-east-2"))
	require.Nil(t, err)
	client = ec2.NewFromConfig(cfg)
	useast2 := &awszones.MockVPC{}
	useast2.On("Client").Return(client)
	useast2.On("Region").Return("us-east-2")
	useast2.On("Subnets", mock.Anything).Return(map[string]string{"us-east-2c": ""}, nil)

	return []awszones.VPC{useast1, useast2}
}
