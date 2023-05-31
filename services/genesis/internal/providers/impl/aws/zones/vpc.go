package awszones

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	awsutils "go.taskfleet.io/services/genesis/internal/providers/impl/aws/utils"
	"go.taskfleet.io/services/genesis/internal/template"
)

// VPC is a type which allows to retrieve information about a single virtual private network that
// spans possibly multiple availability zones within a single region.
type VPC interface {
	// Client returns the EC2 client used for the region that this VPC belongs to.
	Client() *ec2.Client
	// Region returns the region that the VPC belongs to.
	Region() string
	// ID returns the ID of the VPC.
	ID() string
	// Subnets returns a mapping from this VPC's zones to subnet IDs of subnets belonging to the
	// VPC. Implementations may cache the (successful) return value of this method.
	Subnets(ctx context.Context) (map[string]string, error)
	// SecurityGroups returns the security groups that this VPC provides. Typically,
	// implementations choose to return a subset of available security groups based on some
	// filters.
	SecurityGroups(ctx context.Context) ([]string, error)
}

type vpc struct {
	client *ec2.Client
	region string
	id     string

	config              template.AwsNetworkConfig
	subnetsCache        map[string]string
	securityGroupsCache []string
}

// FindVPCs returns the VPCs for all regions where exactly one VPC can be found that provides tags
// that match the VPC selector of the provided configuration. If no VPC is found, a VPC for this
// region is not included in the result. If more than one VPC is found, an error is returned.
//
// The configuration is stored within the VPC such that its member functions use the selectors for
// subnets and security groups.
func FindVPCs(
	ctx context.Context, client *ec2.Client, config template.AwsNetworkConfig,
) ([]VPC, error) {
	// First, fetch all regions
	output, err := client.DescribeRegions(ctx, &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to obtain available regions: %s", err)
	}
	var regions []string
	for _, item := range output.Regions {
		regions = append(regions, *item.RegionName)
	}

	// Find all VPCs
	vpcs, err := awsutils.ParallelForEachRegion(ctx, regions,
		func(ctx context.Context, client *ec2.Client, region string) (VPC, error) {
			// Find the VPC which is tagged appropriately
			vpcs, err := client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{
				MaxResults: aws.Int32(100),
				Filters:    awsutils.TagFiltersFromMap(config.VpcSelector),
			})
			if err != nil {
				return nil, fmt.Errorf("failed to list VPCs: %s", err)
			}
			if len(vpcs.Vpcs) > 1 {
				return nil, fmt.Errorf("found more than one VPC")
			}
			if len(vpcs.Vpcs) == 0 {
				// If we do not find a VPC, this is fine. We simply indicate that there are no
				// availability zones with usable subnets
				return nil, nil
			}
			return &vpc{
				client: client,
				region: region,
				id:     *vpcs.Vpcs[0].VpcId,
				config: config,
			}, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find VPC within each region: %s", err)
	}

	// Assemble result by filtering out nil values
	result := []VPC{}
	for _, vpc := range vpcs {
		if vpc != nil {
			result = append(result, vpc)
		}
	}
	return result, nil
}

//-------------------------------------------------------------------------------------------------

func (v *vpc) Client() *ec2.Client {
	return v.client
}

func (v *vpc) Region() string {
	return v.region
}

func (v *vpc) ID() string {
	return v.id
}

func (v *vpc) Subnets(ctx context.Context) (map[string]string, error) {
	if v.subnetsCache != nil {
		return v.subnetsCache, nil
	}

	// Find all subnets
	subnets, err := v.client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		MaxResults: aws.Int32(100),
		Filters: append(
			awsutils.TagFiltersFromMap(v.config.SubnetSelector),
			types.Filter{Name: aws.String("vpc-id"), Values: []string{v.id}},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list subnets: %s", err)
	}

	// Extract subnet IDs
	zones := map[string]string{}
	for _, subnet := range subnets.Subnets {
		if _, ok := zones[*subnet.AvailabilityZone]; ok {
			return nil, fmt.Errorf(
				"duplicate subnet for zone %q", *subnet.AvailabilityZone,
			)
		}
		zones[*subnet.AvailabilityZone] = *subnet.SubnetId
	}
	v.subnetsCache = zones
	return zones, nil
}

func (v *vpc) SecurityGroups(ctx context.Context) ([]string, error) {
	if v.securityGroupsCache != nil {
		return v.securityGroupsCache, nil
	}

	// Get the security groups
	groups, err := v.client.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		MaxResults: aws.Int32(100),
		Filters: append(
			awsutils.TagFiltersFromMap(v.config.SecurityGroupSelector),
			types.Filter{Name: aws.String("vpc-id"), Values: []string{v.id}},
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list security groups: %s", err)
	}

	// Extract their IDs
	result := []string{}
	for _, group := range groups.SecurityGroups {
		result = append(result, *group.GroupId)
	}
	v.securityGroupsCache = result
	return result, nil
}
