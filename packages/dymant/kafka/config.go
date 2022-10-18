package kafka

import "go.taskfleet.io/packages/eagle"

// Config allows to easily read Kafka configuration.
type Config struct {
	ID               string      `json:"id"`
	BootstrapServers []string    `json:"bootstrapServers"`
	Auth             *AuthConfig `json:"auth"`
}

// AuthConfig describes the Kafka authentication configuration.
type AuthConfig struct {
	Username  eagle.String `json:"username"`
	Password  eagle.String `json:"password"`
	Mechanism string       `json:"mechanism"`
}

// Options returns the client options that can be derived from the configuration.
func (c Config) Options() []ClientOption {
	result := []ClientOption{}
	if c.Auth != nil {
		result = append(result, WithSASLAuthentication(
			c.Auth.Username.Value(),
			c.Auth.Password.Value(),
			c.Auth.Mechanism,
		))
	}
	return result
}
