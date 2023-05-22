package gcp

import (
	"context"
	"fmt"

	compute "cloud.google.com/go/compute/apiv1"
	gcpinstances "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/instances"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type provider struct {
	options   providerOptions
	zones     gcpzones.Client
	instances *gcpinstances.Client
}

// New creates a new provider that manages instances on the Google Cloud Platform. The context
// passed to this function should live as long as the provider is being used. Once it is cancelled,
// this provider drops important resources.
//
// For authentication, default credentials are loaded. The ID of the GCP project that the provider
// interacts with can be inferred from these credentials except for the case where the default
// credentials are obtained from a user account. In this case, the project must be passed
// explicitly.
func New(
	ctx context.Context, config template.GcpConfig, options ...Option,
) (providers.Provider, error) {
	// First, we obtain credentials for authentication
	credentials, err := google.FindDefaultCredentials(ctx, compute.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %s", err)
	}

	// Then, we find all metadata required for initializing the provider
	o := newOptions(options...)
	if err := o.inferMissingIfPossible(ctx, credentials); err != nil {
		return nil, fmt.Errorf(
			"provider options were not complete and could not be inferred: %s", err,
		)
	}

	// For initializing GCP clients, we use a dedicated factory
	clients := gcputils.NewClientFactory(ctx, option.WithCredentials(credentials))

	// ...and eventually, we can create our own higher-level clients
	zoneClient, err := gcpzones.NewClient(ctx, clients, o.projectID, config.Network.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare zone client: %s", err)
	}

	instanceClient, err := gcpinstances.NewClient(
		ctx, o.identifier, o.projectID, config, clients, zoneClient,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare instance client: %s", err)
	}

	return &provider{
		options:   o,
		zones:     zoneClient,
		instances: instanceClient,
	}, nil
}

func (p *provider) Name() typedefs.Provider {
	return typedefs.ProviderGoogleCloudPlatform
}

func (p *provider) AccountName() string {
	return p.options.projectID
}

func (p *provider) Zones() providers.ZoneClient {
	return p.zones
}

func (p *provider) Instances() providers.InstanceClient {
	return p.instances
}
