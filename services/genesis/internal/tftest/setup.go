package tftest

import (
	"context"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
)

func Setup(ctx context.Context, t testing.TB, workdir string) *tfexec.Terraform {
	// Get terraform binary
	tf, err := tfexec.NewTerraform(workdir, jack.Must(exec.LookPath("terraform")))
	require.Nil(t, err)

	// Initialize dependencies
	err = tf.Init(ctx)
	require.Nil(t, err)

	// Create resources
	err = tf.Apply(ctx)
	require.Nil(t, err)

	// Make sure that all resources are torn down after test
	t.Cleanup(func() {
		err := tf.Destroy(ctx)
		require.Nil(t, err)
	})
	return tf
}
