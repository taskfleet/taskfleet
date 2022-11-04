package gcp

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/compute/metadata"
	"github.com/borchero/zeus/pkg/zeus"
	"golang.org/x/oauth2/google"
)

// Option is an interface that is implemented by any functions that allow customizing the GCP
// provider.
type Option interface {
	apply(opts *providerOptions)
}

type providerOptions struct {
	projectID  string
	identifier string
	network    string
}

func newOptions(options ...Option) providerOptions {
	opts := providerOptions{}
	for _, o := range options {
		o.apply(&opts)
	}
	return opts
}

func (o *providerOptions) inferMissingIfPossible(
	ctx context.Context, credentials *google.Credentials,
) error {
	if o.projectID == "" {
		id, err := o.findProjectID(credentials)
		if err != nil {
			return fmt.Errorf("failed to infer project ID: %s", err)
		}
		o.projectID = id
	}
	if o.identifier == "" {
		zeus.Logger(ctx).Warn("GCP identifier not set explicitly, defaulting to 'genesis'.")
		o.identifier = "genesis"
	}
	if o.network == "" {
		network, err := o.findNetworkName()
		if err != nil {
			return fmt.Errorf("failed to infer network name: %s", err)
		}
		o.network = network
	}
	return nil
}

func (o *providerOptions) findProjectID(credentials *google.Credentials) (string, error) {
	// First, attempt to get project from credentials
	if credentials.ProjectID != "" {
		return credentials.ProjectID, nil
	}

	// Otherwise, try to read from metadata server
	if metadata.OnGCE() {
		id, err := metadata.ProjectID()
		if err != nil {
			return "", fmt.Errorf("failed to read project ID from metadata server: %s", err)
		}
		return id, nil
	}
	return "", fmt.Errorf(
		"credentials do not provide project ID and process is not running on GCP",
	)
}

func (o *providerOptions) findNetworkName() (string, error) {
	type networkInterface struct {
		Network string `json:"network"`
	}

	// We can only infer the network name if the instance is running on GCP and attached to
	// exactly one network.
	if metadata.OnGCE() {
		response, err := metadata.Get("instance/network-interfaces/?recursive=true")
		if err != nil {
			return "", fmt.Errorf(
				"failed to query network interfaces from metadata server: %s", err,
			)
		}

		var interfaces []networkInterface
		if err := json.Unmarshal([]byte(response), &interfaces); err != nil {
			return "", fmt.Errorf("failed to parse response from metadata server: %s", err)
		}
		switch len(interfaces) {
		case 0:
			return "", fmt.Errorf("instance is not attached to a network")
		case 1:
			return interfaces[0].Network, nil
		default:
			return "", fmt.Errorf("instance is attached to more than one network")
		}
	}
	return "", fmt.Errorf("network name can only be inferred when running on GCP")
}

//-------------------------------------------------------------------------------------------------

// WithProjectID specifies the ID of the project that the provider should be interacting with.
func WithProjectID(id string) Option {
	return projectID(id)
}

type projectID string

func (p projectID) apply(opts *providerOptions) {
	opts.projectID = string(p)
}

//-------------------------------------------------------------------------------------------------

// WithIdentifier specifies a unique identifier to discern the provider from other providers
// interacting with the same project on GCP.
func WithIdentifier(id string) Option {
	return identifier(id)
}

type identifier string

func (i identifier) apply(opts *providerOptions) {
	opts.projectID = string(i)
}

//-------------------------------------------------------------------------------------------------

// WithNetworkName specifies the name of the network that the provider launches instances into.
func WithNetworkName(name string) Option {
	return networkName(name)
}

type networkName string

func (n networkName) apply(opts *providerOptions) {
	opts.projectID = string(n)
}
