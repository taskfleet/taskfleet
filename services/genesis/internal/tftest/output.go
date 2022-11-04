package tftest

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/require"
)

func GetOutput[T any](ctx context.Context, t *testing.T, tf *tfexec.Terraform, key string) T {
	// Get output
	output, err := tf.Output(ctx)
	require.Nil(t, err)

	// Parse output
	var result T
	err = json.Unmarshal(output[key].Value, &result)
	require.Nil(t, err)
	return result
}
