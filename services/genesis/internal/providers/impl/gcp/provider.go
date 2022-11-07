package gcp

import (
	"context"

	compute "cloud.google.com/go/compute/apiv1"
	gcpinstances "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/instances"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type provider struct {
	options   providerOptions
	zones     *gcpzones.Client
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
		return nil, providers.NewClientError("failed to find default credentials", err)
	}

	// Then, we find all metadata required for initializing the provider
	o := newOptions(options...)
	if err := o.inferMissingIfPossible(ctx, credentials); err != nil {
		return nil, providers.NewClientError(
			"provider options were not complete and could not be inferred", err,
		)
	}

	// Now, we can initialize GCP clients...
	opt := option.WithCredentials(credentials)
	gcpZoneClient, err := compute.NewZonesRESTClient(ctx, opt)
	if err != nil {
		return nil, providers.NewFatalError("failed to create new zone client: %s", err)
	}
	deferClose(ctx, gcpZoneClient)
	gcpNetworkClient, err := compute.NewNetworksRESTClient(ctx, opt)
	if err != nil {
		return nil, providers.NewFatalError("failed to create new network client: %s", err)
	}
	deferClose(ctx, gcpNetworkClient)
	gcpAcceleratorTypesClient, err := compute.NewAcceleratorTypesRESTClient(ctx, opt)
	if err != nil {
		return nil, providers.NewFatalError("failed to create new accelerator client: %s", err)
	}
	deferClose(ctx, gcpAcceleratorTypesClient)
	gcpMachineTypesClient, err := compute.NewMachineTypesRESTClient(ctx, opt)
	if err != nil {
		return nil, providers.NewFatalError("failed to create new machine types client: %s", err)
	}
	deferClose(ctx, gcpMachineTypesClient)
	gcpInstanceClient, err := compute.NewInstancesRESTClient(ctx, opt)
	if err != nil {
		return nil, providers.NewFatalError("failed to create new instance client: %s", err)
	}
	deferClose(ctx, gcpInstanceClient)
	gcpDisksClient, err := compute.NewDisksRESTClient(ctx, opt)
	if err != nil {
		return nil, providers.NewFatalError("failed to create new disks client: %s", err)
	}
	deferClose(ctx, gcpDisksClient)

	// ...and eventually, we can create our own higher-level clients
	zoneClient, err := gcpzones.NewClient(
		ctx,
		gcpZoneClient,
		gcpNetworkClient,
		gcpAcceleratorTypesClient,
		o.projectID,
		config.Network.Name,
	)
	if err != nil {
		return nil, providers.NewAPIError("failed to prepare zone client", err)
	}

	instanceClient, err := gcpinstances.NewClient(
		ctx,
		o.identifier,
		o.projectID,
		config,
		gcpMachineTypesClient,
		gcpInstanceClient,
		gcpNetworkClient,
		gcpDisksClient,
		zoneClient,
	)
	if err != nil {
		return nil, providers.NewAPIError("failed to prepare instance client", err)
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

type closeable interface {
	Close() error
}

func deferClose(ctx context.Context, objects ...closeable) {
	<-ctx.Done()
	for _, obj := range objects {
		obj.Close() // nolint:errcheck
	}
}
