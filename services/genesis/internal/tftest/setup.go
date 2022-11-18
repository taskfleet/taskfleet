package tftest

import (
	"context"
	"fmt"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/jack"
)

func Setup(
	ctx context.Context, t testing.TB, workdir string, variables ...string,
) *tfexec.Terraform {
	// Get terraform binary
	tf, err := tfexec.NewTerraform(workdir, jack.Must(exec.LookPath("terraform")))
	require.Nil(t, err)

	// Initialize dependencies
	err = tf.Init(
		ctx,
		tfexec.Reconfigure(true),
		tfexec.BackendConfig(fmt.Sprintf("path=%s/state.tfstate", t.TempDir())),
	)
	require.Nil(t, err)

	// Make sure that all resources are torn down after test. Do this even before resources are
	// created to catch any failures that create a subset of resources.
	t.Cleanup(func() {
		err := tf.Destroy(ctx)
		require.Nil(t, err)
	})

	// Create resources
	assignments := []tfexec.ApplyOption{}
	for _, v := range variables {
		assignments = append(assignments, tfexec.Var(v))
	}
	err = tf.Apply(ctx, assignments...)
	require.Nil(t, err)

	return tf
}
