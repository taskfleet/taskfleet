package awsutils

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// TagFiltersFromMap converts the provided selector for tags in a valid set of filters to use in
// AWS API requests.
func TagFiltersFromMap(filters map[string]string) []types.Filter {
	result := []types.Filter{}
	for key, value := range filters {
		result = append(result, types.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", key)),
			Values: []string{value},
		})
	}
	return result
}
