package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ClientOption allows to update the configuration of a Kafka client. Such an option applies to
// both producers and consumers (i.e. publishers/subscribers).
type ClientOption interface {
	clientOption()
	apply(config kafka.ConfigMap) error
}

type dummyClientOption struct{}

func (dummyClientOption) clientOption() {}

//-------------------------------------------------------------------------------------------------
// SASL AUTHENTICATION
//-------------------------------------------------------------------------------------------------

type configOptionSaslAuth struct {
	dummyClientOption
	username      string
	password      string
	authMechanism string
}

// WithSASLAuthentication configures a Kafka client to use SASL authentication with the specified
// username, password, and authentication mechanism. The security protocol will be set to
// `sasl_ssl`.
func WithSASLAuthentication(username, password, authMechanism string) ClientOption {
	return configOptionSaslAuth{
		username:      username,
		password:      password,
		authMechanism: authMechanism,
	}
}

func (c configOptionSaslAuth) apply(config kafka.ConfigMap) error {
	if c.username == "" {
		return fmt.Errorf("SASL username must be set to non-empty value")
	}
	if c.password == "" {
		return fmt.Errorf("SASL password must be set to non-empty value")
	}
	if c.authMechanism == "" {
		return fmt.Errorf("SASL auth mechanism must be provided")
	}

	config["security.protocol"] = "sasl_ssl"
	config["sasl.username"] = c.username
	config["sasl.password"] = c.password
	config["sasl.mechanisms"] = strings.ToUpper(c.authMechanism)
	return nil
}
