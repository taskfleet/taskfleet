package eagle

import (
	"os"
)

// EnvOrDefault reads the provided environment variable and, if not found, returns the default
// value.
func EnvOrDefault(env, defaultValue string) string {
	if value, ok := os.LookupEnv(env); ok {
		return value
	}
	return defaultValue
}

// LoadConfig loads the configuration from the specified file and merges this configuration with
// environment variables. Env variables take precedence over the configuration values from file.
func LoadConfig(result interface{}, sources ...ConfigSource) error {
	for _, source := range sources {
		if err := source.unmarshal(result); err != nil {
			return err
		}
	}
	return nil
}
