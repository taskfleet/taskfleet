package awsutils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"go.taskfleet.io/packages/jack"
)

// ParallelForEachRegion executes the provided function with an EC2 client for each of the provided
// regions. The execution is performed in parallel.
func ParallelForEachRegion[T any](
	ctx context.Context,
	regions []string,
	execute func(context.Context, *ec2.Client, string) (T, error),
) ([]T, error) {
	return jack.ParallelSliceMap(ctx, regions, func(ctx context.Context, region string) (T, error) {
		var result T

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
		if err != nil {
			return result, fmt.Errorf("failed to load config for region %q: %w", region, err)
		}
		client := ec2.NewFromConfig(cfg)

		result, err = execute(ctx, client, region)
		if err != nil {
			return result, fmt.Errorf("failure in region %q: %w", region, err)
		}
		return result, nil
	})
}
