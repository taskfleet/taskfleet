package eagle_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.taskfleet.io/packages/eagle"
)

func TestLoadConfigFile(t *testing.T) {
	var result credentials
	err := eagle.LoadConfig(&result,
		eagle.WithYAMLFile("testdata/config.yaml", false),
	)
	require.Nil(t, err)

	assert.Equal(t, "user", result.Username)
	assert.Equal(t, "pass", result.Password.Value())
}

func TestLoadConfigOverrideEnv(t *testing.T) {
	// Set password
	var result credentials
	err := os.Setenv("PASSWORD", "new-password")
	require.Nil(t, err)
	defer os.Unsetenv("PASSWORD") // nolint:errcheck

	err = eagle.LoadConfig(&result,
		eagle.WithYAMLFile("testdata/config.yaml", false),
		eagle.WithEnvironment(""),
	)
	require.Nil(t, err)

	assert.Equal(t, "user", result.Username)
	assert.Equal(t, "new-password", result.Password.Value())
}

func TestLoadConfigOverrideEnvFile(t *testing.T) {
	// Set password
	var result credentials
	err := os.Setenv("PREFIX__PASSWORD__FILE", "testdata/password")
	require.Nil(t, err)
	defer os.Unsetenv("PREFIX__PASSWORD__FILE") // nolint:errcheck

	err = eagle.LoadConfig(&result,
		eagle.WithYAMLFile("testdata/config.yaml", false),
		eagle.WithEnvironment("PREFIX"),
	)
	require.Nil(t, err)

	assert.Equal(t, "user", result.Username)
	assert.Equal(t, "securepass", result.Password.Value())
}
