package awsinstances

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/borchero/zeus/pkg/zeus"
	"go.taskfleet.io/packages/jack"
	awsutils "go.taskfleet.io/services/genesis/internal/providers/impl/aws/utils"
	awszones "go.taskfleet.io/services/genesis/internal/providers/impl/aws/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"go.uber.org/zap"
)

// FindAvailableInstances finds the available instances for all provided availability zones. It
// then returns an instance manager for each availability zone.
func FindAvailableInstances(
	ctx context.Context, vpcs []awszones.VPC,
) (map[string]*instances.Manager, error) {
	// For each region, we want to fetch all instance types
	managers, err := jack.ParallelSliceMap(ctx, vpcs,
		func(ctx context.Context, vpc awszones.VPC) (map[string]*instances.Manager, error) {
			client := vpc.Client()

			// Get all instance types in this region
			awsInstanceTypes, err := awsutils.RunPaginatedRequest(
				func(nextToken *string) (*ec2.DescribeInstanceTypesOutput, error) {
					return client.DescribeInstanceTypes(ctx, &ec2.DescribeInstanceTypesInput{
						NextToken: nextToken,
						Filters: []types.Filter{
							{Name: aws.String("current-generation"), Values: []string{"true"}},
						},
					})
				},
				func(out *ec2.DescribeInstanceTypesOutput) ([]types.InstanceTypeInfo, *string) {
					return out.InstanceTypes, out.NextToken
				},
			)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to list instance types in region %q: %s", vpc.Region(), err,
				)
			}

			// Parse these instances types into our internal format
			instanceTypes := instanceTypesFromAwsInstanceTypes(
				zeus.WithFields(ctx, zap.String("region", vpc.Region())),
				awsInstanceTypes,
			)

			// For each availability zone, check whether the found instance types are actually
			// available. We cannot use this request directly since it does not provide details
			// about the instance types (only their name).
			// For this, we also need to fetch all zones within the VPC.
			subnets, err := vpc.Subnets(ctx)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to get available zones in region %q: %s", vpc.Region(), err,
				)
			}
			zones := jack.MapKeys(subnets)

			managers, err := jack.ParallelSliceMap(
				ctx, zones, managerForAvailabilityZone(client, instanceTypes),
			)
			if err != nil {
				return nil, fmt.Errorf(
					"failed to find offered instance types in zones of region %q: %s",
					vpc.Region(), err,
				)
			}
			return jack.ZipMap(zones, managers), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list available instance types: %s", err)
	}

	// Assemble result
	result := map[string]*instances.Manager{}
	for _, zones := range managers {
		for zone, manager := range zones {
			result[zone] = manager
		}
	}
	return result, nil
}

func instanceTypesFromAwsInstanceTypes(
	ctx context.Context, instanceTypes []types.InstanceTypeInfo,
) []instances.Type {
	result := make([]instances.Type, 0, len(instanceTypes))
	for _, instance := range instanceTypes {
		ctx := zeus.WithFields(ctx, zap.String("name", string(instance.InstanceType)))

		// Get the instance architecture
		architecture, err := typedefs.CPUArchitectureFromProviderAws(
			instance.ProcessorInfo.SupportedArchitectures,
		)
		if err != nil {
			zeus.Logger(ctx).Debug("skipping instance type", zap.Error(err))
			continue
		}

		// Get the GPUs (if applicable)
		var gpu *instances.GPUResources
		if instance.GpuInfo != nil {
			if len(instance.GpuInfo.Gpus) > 1 {
				zeus.Logger(ctx).Debug("skipping instance type",
					zap.Error(errors.New("more than one type of GPU is not supported")),
				)
				continue
			}
			info := instance.GpuInfo.Gpus[0]
			kind, err := typedefs.GPUKindUnmarshalProviderAws(info)
			if err != nil {
				zeus.Logger(ctx).Debug("skipping instance type", zap.Error(err))
				continue
			}
			gpu = &instances.GPUResources{Kind: kind, Count: uint16(*info.Count)}
		}

		// Build our internal type
		result = append(result, instances.Type{
			Name:         string(instance.InstanceType),
			Architecture: architecture,
			Resources: instances.Resources{
				CPUCount:  uint16(*instance.VCpuInfo.DefaultVCpus),
				MemoryMiB: uint32(*instance.MemoryInfo.SizeInMiB),
				GPU:       gpu,
			},
		})
	}
	return result
}

func managerForAvailabilityZone(
	client *ec2.Client,
	instanceTypes []instances.Type,
) func(context.Context, string) (*instances.Manager, error) {
	return func(ctx context.Context, zone string) (*instances.Manager, error) {
		// Get all offerings for the zone
		offerings, err := awsutils.RunPaginatedRequest(
			func(nextToken *string) (*ec2.DescribeInstanceTypeOfferingsOutput, error) {
				return client.DescribeInstanceTypeOfferings(
					ctx, &ec2.DescribeInstanceTypeOfferingsInput{
						Filters: []types.Filter{{
							Name:   aws.String("location"),
							Values: []string{zone},
						}},
						LocationType: types.LocationTypeAvailabilityZone,
						NextToken:    nextToken,
					},
				)
			},
			func(out *ec2.DescribeInstanceTypeOfferingsOutput) (
				[]types.InstanceTypeOffering, *string,
			) {
				return out.InstanceTypeOfferings, out.NextToken
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to list offerings in zone %q: %s", zone, err)
		}

		// Filter the instance types by the offerings and build list of instances
		offeredInstances := map[string]struct{}{}
		for _, offering := range offerings {
			offeredInstances[string(offering.InstanceType)] = struct{}{}
		}

		result := []instances.Type{}
		for _, instance := range instanceTypes {
			if _, ok := offeredInstances[instance.Name]; ok {
				result = append(result, instance)
			}
		}
		return instances.NewManager(result), nil
	}
}
