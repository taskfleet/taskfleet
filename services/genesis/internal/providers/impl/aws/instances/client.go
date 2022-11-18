package awsinstances

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/docker/go-units"
	awsutils "go.taskfleet.io/services/genesis/internal/providers/impl/aws/utils"
	awszones "go.taskfleet.io/services/genesis/internal/providers/impl/aws/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// Client represents an AWS instance client.
type Client struct {
	vpcs     map[string]awszones.VPC
	managers map[string]*instances.Manager
	config   template.AwsConfig
}

// NewClient initializes a new instance client.
//
// The provided VPCs must encompass all zones for which managers are provided. They are used to
// source the network configuration of launched instances.
//
// The instance client may launch instances in the zones for which a manager is provided and
// launches only those instances which are provided by the zone's manager. Typically, one would
// use `FindAvailableInstances` to initialize the mapping from zones to managers.
func NewClient(
	ctx context.Context,
	vpcs []awszones.VPC,
	managers map[string]*instances.Manager,
	config template.AwsConfig,
) *Client {
	// Map the VPCs to their regions
	vpcMap := map[string]awszones.VPC{}
	for _, vpc := range vpcs {
		vpcMap[vpc.Region()] = vpc
	}

	// Initialize the client
	return &Client{
		vpcs:     vpcMap,
		managers: managers,
		config:   config,
	}
}

// Create implements the `providers.InstanceClient` interface.
func (c *Client) Create(
	ctx context.Context, ref providers.InstanceRef, spec providers.InstanceSpec,
) (providers.InstancePromise, error) {
	// First, get a client
	region := awsutils.RegionFromZone(ref.Zone)
	vpc, ok := c.vpcs[region]
	if !ok {
		return nil, providers.NewClientError(
			fmt.Sprintf("cannot launch instance in region %q", region), nil,
		)
	}
	client := vpc.Client()

	// Then, assemble all kinds of metadata:
	// - Subnet ID
	subnets, err := vpc.Subnets(ctx)
	if err != nil {
		return nil, providers.NewAPIError("failed to get subnets for VPC in region", err)
	}
	subnet, ok := subnets[ref.Zone]
	if !ok {
		return nil, providers.NewClientError(
			fmt.Sprintf("did not find subnet for availability zone %q", ref.Zone), nil,
		)
	}

	// - Security groups
	securityGroups, err := vpc.SecurityGroups(ctx)
	if err != nil {
		return nil, providers.NewAPIError("failed to find security groups", err)
	}

	// - AMI ID
	ami, err := c.findAMI(ctx, client, spec.InstanceType)
	if err != nil {
		return nil, providers.NewAPIError("failed to find suitable AMI", err)
	}

	// - Device mappings
	diskSize, err := units.RAMInBytes(c.config.Boot.DiskSize)
	if err != nil {
		return nil, providers.NewClientError("invalid boot disk size", err)
	}
	blockDevices := []types.BlockDeviceMapping{{
		DeviceName: aws.String("/dev/sda1"),
		Ebs: &types.EbsBlockDevice{
			VolumeType:          types.VolumeTypeGp3,
			VolumeSize:          aws.Int32(int32(diskSize / units.GiB)),
			Throughput:          aws.Int32(125),
			Iops:                aws.Int32(3000),
			DeleteOnTermination: aws.Bool(true),
		},
	}}
	for i, disk := range c.config.ExtraDisks {
		sizePerCPU, err := units.RAMInBytes(disk.SizePerCPU)
		if err != nil {
			return nil, providers.NewClientError(
				fmt.Sprintf("invalid size of extra disk %d", i), nil,
			)
		}
		size := int32((sizePerCPU / units.GiB) * int64(spec.InstanceType.CPUCount))
		blockDevices = append(blockDevices, types.BlockDeviceMapping{
			DeviceName: aws.String("/dev/sd" + string('b'+rune(i))),
			Ebs: &types.EbsBlockDevice{
				VolumeType:          types.VolumeTypeGp3,
				VolumeSize:          aws.Int32(size),
				Throughput:          aws.Int32(125),
				Iops:                aws.Int32(3000),
				DeleteOnTermination: aws.Bool(true),
			},
		})
	}

	// - Spot instance
	var marketOptions *types.InstanceMarketOptionsRequest
	if spec.IsSpot {
		marketOptions = &types.InstanceMarketOptionsRequest{
			MarketType: types.MarketTypeSpot,
			SpotOptions: &types.SpotMarketOptions{
				InstanceInterruptionBehavior: types.InstanceInterruptionBehaviorTerminate,
			},
		}
	}

	// - Tags
	tags := []types.Tag{
		{Key: aws.String("Owner"), Value: aws.String("taskfleet-instance-manager")},
		{Key: aws.String("GlobalID"), Value: aws.String(ref.ID.String())},
	}
	for key, value := range c.config.Metadata {
		tags = append(tags, types.Tag{Key: aws.String(key), Value: aws.String(value)})
	}

	// Then, create the instance
	response, err := client.RunInstances(ctx, &ec2.RunInstancesInput{
		// Basic configuration
		MinCount:                          aws.Int32(1),
		MaxCount:                          aws.Int32(1),
		InstanceInitiatedShutdownBehavior: types.ShutdownBehaviorTerminate,
		InstanceMarketOptions:             marketOptions,
		// Compute
		ImageId:      aws.String(ami),
		InstanceType: types.InstanceType(spec.InstanceType.Name),
		// Storage
		BlockDeviceMappings: blockDevices,
		// Network
		SubnetId:         aws.String(subnet),
		SecurityGroupIds: securityGroups,
		// Security
		IamInstanceProfile: &types.IamInstanceProfileSpecification{
			Name: aws.String(c.config.Iam.InstanceProfile),
		},
		// Metadata
		MetadataOptions: &types.InstanceMetadataOptionsRequest{
			HttpEndpoint:         types.InstanceMetadataEndpointStateEnabled,
			InstanceMetadataTags: types.InstanceMetadataTagsStateEnabled,
		},
		TagSpecifications: []types.TagSpecification{{
			ResourceType: types.ResourceTypeInstance,
			Tags:         tags,
		}},
	})
	if err != nil {
		return nil, providers.NewAPIError("failed to create instance", err)
	}

	// Return new promise for the instance
	ref.ProviderID = *response.Instances[0].InstanceId
	return newPromise(client, ref, c.managers[ref.Zone]), nil
}

// Get implements the `providers.InstanceClient` interface.
func (c *Client) Get(ctx context.Context, ref providers.InstanceRef) (providers.Instance, error) {
	// Get the clients
	region := awsutils.RegionFromZone(ref.Zone)
	vpc, ok := c.vpcs[region]
	if !ok {
		return providers.Instance{}, providers.NewClientError(
			fmt.Sprintf("cannot describe instance in region %q", region), nil,
		)
	}
	client := vpc.Client()

	// Describe the instance
	output, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{ref.ProviderID},
	})
	if err != nil {
		return providers.Instance{}, providers.NewAPIError("failed to describe instance", err)
	}
	if len(output.Reservations) == 0 || len(output.Reservations[0].Instances) == 0 {
		return providers.Instance{}, providers.NewClientError("no instance found", nil)
	}
	instance := output.Reservations[0].Instances[0]

	// Extract information from the instance
	result, err := instanceFromAwsInstance(instance, c.managers[ref.Zone], ref)
	if err != nil {
		return providers.Instance{}, providers.NewFatalError(
			"failed to parse returned instance", err,
		)
	}
	return result, nil
}

// List implements the `providers.InstanceClient` interface.
// func (c *Client) List(ctx context.Context) ([]providers.Instance, error) {
// 	// For each zone, we need to list all instances
// 	result := make([]types.Reservation, 0)
// 	var mutex sync.Mutex

// 	eg, ctx := errgroup.WithContext(ctx)
// 	for _, region := range c.vpcs {
// 		name := region
// 		eg.Go(func() error {
// 			client, err := c.factory.getClient(name)
// 			if err != nil {
// 				return fmt.Errorf("failed to get client for region %q: %s", name, err)
// 			}

// 			instances, err := awsutils.RunPaginatedRequest(
// 				func(token *string) (*ec2.DescribeInstancesOutput, error) {
// 					return client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
// 						Filters: []types.Filter{
// 							{Name: aws.String("tag:Owner"), Values: []string{"taskfleet"}},
// 						},
// 						NextToken: token,
// 					})
// 				},
// 				func(output *ec2.DescribeInstancesOutput) ([]types.Reservation, *string) {
// 					return output.Reservations, output.NextToken
// 				},
// 			)
// 			if err != nil {
// 				return fmt.Errorf("failed to list instances in region %q: %s", name, err)
// 			}

// 			mutex.Lock()
// 			defer mutex.Unlock()
// 			result = append(result, instances...)
// 			return nil
// 		})
// 	}
// 	if err := eg.Wait(); err != nil {
// 		return nil, providers.NewAPIError("failed to fetch all instances", err)
// 	}

// 	// Then, we can transform the fetched instances into our common format
// 	return nil, nil
// }

// Delete implements the `providers.InstanceClient` interface.
func (c *Client) Delete(ctx context.Context, ref providers.InstanceRef) error {
	// First, we obtain a client for the region
	region := awsutils.RegionFromZone(ref.Zone)
	vpc, ok := c.vpcs[region]
	if !ok {
		return providers.NewClientError(fmt.Sprintf("cannot use region %q", region), nil)
	}
	client := vpc.Client()

	// Then, we can issue the delete request
	response, err := client.TerminateInstances(ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{ref.ProviderID},
	})
	if err != nil {
		return providers.NewAPIError("failed to terminate instance", err)
	}
	if len(response.TerminatingInstances) != 1 {
		return providers.NewAPIError("request did not terminate exactly one instance", nil)
	}
	if *response.TerminatingInstances[0].InstanceId != ref.ProviderID {
		return providers.NewAPIError("terminated instance does not coincide with request", nil)
	}
	return nil
}

//-------------------------------------------------------------------------------------------------
// UTILS
//-------------------------------------------------------------------------------------------------

func (c *Client) findAMI(
	ctx context.Context, client *ec2.Client, instanceType instances.Type,
) (string, error) {
	// First, we need to get the appropriate tags for the instance type
	var gpuKind *typedefs.GPUKind
	if instanceType.GPU != nil {
		gpuKind = &instanceType.GPU.Kind
	}
	option := template.MatchingOption(c.config.Boot.Amis, gpuKind, instanceType.Architecture)
	if option == nil {
		return "", fmt.Errorf(
			"instance type %q does not match any template option", instanceType.Name,
		)
	}

	// Then, we try to find the AMI ID
	images, err := client.DescribeImages(ctx, &ec2.DescribeImagesInput{
		Owners: []string{option.Owner},
		Filters: append(
			awsutils.TagFiltersFromMap(option.Selector),
			types.Filter{
				Name:   aws.String("architecture"),
				Values: []string{string(instanceType.Architecture.ToProviderAws())},
			},
		),
	})
	if err != nil {
		return "", fmt.Errorf("failed to list AMIs matching filters: %s", err)
	}
	if len(images.Images) == 0 {
		return "", fmt.Errorf("did not find AMI matching filters")
	}
	if len(images.Images) > 1 {
		return "", fmt.Errorf("found more than one AMI matching filters")
	}
	return *images.Images[0].ImageId, nil
}
