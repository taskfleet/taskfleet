package eagle

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ghodss/yaml"
	"go.taskfleet.io/packages/eagle/internal/envconfig"
)

// ConfigSource is implemented by types which allow to load configuration from a particular source.
type ConfigSource interface {
	unmarshal(o interface{}) error
}

//-------------------------------------------------------------------------------------------------
// YAML FILE
//-------------------------------------------------------------------------------------------------

type configSourceYAML struct {
	file     string
	optional bool
}

// WithYAMLFile reads the configuration from the YAML file at the specified path.
func WithYAMLFile(path string, optional bool) ConfigSource {
	return configSourceYAML{path, optional}
}

func (s configSourceYAML) unmarshal(o interface{}) error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if s.optional {
			return nil
		}
		return fmt.Errorf("failed to read file %q: %s", s.file, err)
	}
	if err := yaml.Unmarshal(data, o); err != nil {
		if s.optional {
			return nil
		}
		return fmt.Errorf("failed to parse YAML file %q: %s", s.file, err)
	}
	return nil
}

//-------------------------------------------------------------------------------------------------
// JSON FILE
//-------------------------------------------------------------------------------------------------

type configSourceJSON struct {
	file     string
	optional bool
}

// WithJSONFile reads the configuration from the JSON file at the specified path.
func WithJSONFile(path string, optional bool) ConfigSource {
	return configSourceJSON{path, optional}
}

func (s configSourceJSON) unmarshal(o interface{}) error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if s.optional {
			return nil
		}
		return fmt.Errorf("failed to read file %q: %s", s.file, err)
	}
	if err := json.Unmarshal(data, o); err != nil {
		if s.optional {
			return nil
		}
		return fmt.Errorf("failed to parse JSON file %q: %s", s.file, err)
	}
	return nil
}

//-------------------------------------------------------------------------------------------------
// ENVIRONMENT VARIABLES
//-------------------------------------------------------------------------------------------------

type configSourceEnvironment struct {
	prefix string
}

// WithEnvironment reads the configuration from environment variables, using the specified prefix.
func WithEnvironment(prefix string) ConfigSource {
	return configSourceEnvironment{prefix}
}

func (s configSourceEnvironment) unmarshal(o interface{}) error {
	return envconfig.Process(s.prefix, "__", o)
}
