package awsutils

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

func TestTagFiltersFromMap(t *testing.T) {
	tags := map[string]string{
		"hello":   "world",
		"another": "tag",
	}
	expected := []types.Filter{
		{Name: aws.String("tag:hello"), Values: []string{"world"}},
		{Name: aws.String("tag:another"), Values: []string{"tag"}},
	}
	actual := TagFiltersFromMap(tags)
	assert.ElementsMatch(t, expected, actual)
}
