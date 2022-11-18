package gcp

import (
	"context"
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
