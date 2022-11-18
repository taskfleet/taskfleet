package gcpinstances

import (
	"context"
	"fmt"
	"path"

	compute "cloud.google.com/go/compute/apiv1"
	gcputils "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/utils"
	gcpzones "go.taskfleet.io/services/genesis/internal/providers/impl/gcp/zones"
	"go.taskfleet.io/services/genesis/internal/providers/instances"
	providers "go.taskfleet.io/services/genesis/internal/providers/interface"
	"go.taskfleet.io/services/genesis/internal/template"
	"go.taskfleet.io/services/genesis/internal/typedefs"
	computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"
	"google.golang.org/protobuf/proto"
)

// Client represents a GCP instance client.
type Client struct {
	projectID  string
	identifier string
	network    string
	config     template.GcpConfig

	service *compute.InstancesClient
	zones   gcpzones.Client

	disksHelper        *disksHelper
	reservationsHelper *reservationsHelper
	instanceManagers   map[string]*instances.Manager
}

// NewClient initializes a new instance client which first fetch all available instance types.
func NewClient(
	ctx context.Context,
	identifier, projectID string,
	config template.GcpConfig,
	clients gcputils.ClientFactory,
	zones gcpzones.Client,
) (*Client, error) {
	// First, we initialize instance managers for each zone
	allInstances, err := findAvailableInstanceTypes(ctx, clients.MachineTypes(), zones, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available instance types: %s", err)
	}
	instanceManagers := make(map[string]*instances.Manager)
	for zone, zoneInstances := range allInstances {
		instanceManagers[zone], err = instances.NewManager(zoneInstances)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to initialize instance type manager for zone %q: %s", zone, err,
			)
		}
	}

	// Then, fetch the link of the network
	network, err := clients.Networks().Get(ctx, &computepb.GetNetworkRequest{
		Project: projectID,
		Network: config.Network.Name,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to find network %q: %s", config.Network.Name, err)
	}

	// Initialize helpers
	reservationsHelper, err := newReservationsHelper(config.Reservations)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare reservation configuration: %s", err)
	}

	disksHelper, err := newDisksHelper(
		ctx, projectID, config.Boot, config.ExtraDisks, config.Disks.Type, clients.DiskTypes(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare disk configuration: %s", err)
	}

	// Eventually, we can initialize the client
	return &Client{
		projectID:          projectID,
		identifier:         identifier,
		network:            network.GetSelfLink(),
		config:             config,
		zones:              zones,
		service:            clients.Instances(),
		disksHelper:        disksHelper,
		reservationsHelper: reservationsHelper,
		instanceManagers:   instanceManagers,
	}, nil
}

//-------------------------------------------------------------------------------------------------
// INTERFACE
//-------------------------------------------------------------------------------------------------

// Find implements the `providers.InstanceClient` interface.
func (c *Client) Find(
	zone string, resources instances.Resources, architecture typedefs.CPUArchitecture,
) (instances.Type, error) {
	// Ensure that zone is managed
	manager, ok := c.instanceManagers[zone]
	if !ok {
		return instances.Type{}, providers.NewClientError(
			fmt.Sprintf("no instances are managed in zone %q", zone), nil,
		)
	}

	// Add reservations to resource requirement if required
	resources = c.reservationsHelper.updateResources(resources)

	// Find the best instance type
	return manager.FindBestFit(resources, architecture)
}

// Create implements the `providers.InstanceClient` interface.
func (c *Client) Create(
	ctx context.Context, meta providers.InstanceMeta, spec providers.InstanceSpec,
) (providers.InstancePromise, error) {
	// We successively build all of the properties that need to be passed for instance creation...
	meta.ProviderID = fmt.Sprintf("taskfleet-%s", meta.ID)

	// 1) Metadata
	labels := map[string]string{
		LabelID:           meta.ID.String(),
		LabelKeyCreatedBy: c.identifier,
	}
	metadataItems := []*computepb.Items{}
	for k, v := range c.config.Metadata {
		metadataItems = append(metadataItems, &computepb.Items{
			Key:   proto.String(k),
			Value: proto.String(v),
		})
	}
	metadata := &computepb.Metadata{Items: metadataItems}

	// 2) Security
	serviceAccounts := []*computepb.ServiceAccount{{
		Email:  proto.String(c.config.Iam.ServiceAccountEmail),
		Scopes: compute.DefaultAuthScopes(),
	}}

	// 3) Network
	tags := &computepb.Tags{Items: c.config.Network.Tags}
	subnetwork, err := c.zones.GetSubnetwork(meta.ProviderZone)
	if err != nil {
		return nil, providers.NewClientError(
			fmt.Sprintf("failed to find subnetwork for zone %q", meta.ProviderZone), err,
		)
	}
	networkInterface := computepb.NetworkInterface{
		Network:    proto.String(c.network),
		Subnetwork: proto.String(subnetwork),
	}
	if !c.config.Network.Shielded {
		networkInterface.AccessConfigs = []*computepb.AccessConfig{{
			Name: proto.String("External NAT"),
			Type: proto.String("ONE_TO_ONE_NAT"),
		}}
	}
	interfaces := []*computepb.NetworkInterface{&networkInterface}

	// 4) Resources
	guestAccelerators := []*computepb.AcceleratorConfig{}
	if spec.InstanceType.GPU != nil {
		accelerator, err := c.zones.GetAccelerator(meta.ProviderID, spec.InstanceType.GPU.Kind)
		if err != nil {
			return nil, providers.NewClientError("accelerator does not exist in zone", err)
		}
		config, err := accelerator.Config(spec.InstanceType.GPU.Count)
		if err != nil {
			return nil, providers.NewClientError("failed to obtain requested number of GPUs", err)
		}
		guestAccelerators = []*computepb.AcceleratorConfig{config}
	}

	// 5) Scheduling
	scheduling := &computepb.Scheduling{
		OnHostMaintenance: proto.String("MIGRATE"),
		ProvisioningModel: proto.String("STANDARD"),
	}
	if spec.InstanceType.GPU != nil {
		scheduling.OnHostMaintenance = proto.String("TERMINATE")
	}
	if spec.IsSpot {
		scheduling.ProvisioningModel = proto.String("SPOT")
	}

	// 6) Disks
	disks := c.disksHelper.diskConfig(
		meta.ProviderID,
		meta.ProviderZone,
		spec.InstanceType.Resources,
		spec.InstanceType.Architecture,
	)

	// Now, we can finally declare the instance
	instance := computepb.Instance{
		Name: proto.String(meta.ProviderID),
		Description: proto.String(fmt.Sprintf(
			"Genesis-managed instance owned by %q", labels[LabelKeyOwnedBy],
		)),
		Labels:            labels,
		Tags:              tags,
		Metadata:          metadata,
		ServiceAccounts:   serviceAccounts,
		Scheduling:        scheduling,
		MachineType:       proto.String(spec.InstanceType.UID),
		GuestAccelerators: guestAccelerators,
		NetworkInterfaces: interfaces,
		Disks:             disks,
	}

	// And actually make the call to create it
	operation, err := c.service.Insert(ctx, &computepb.InsertInstanceRequest{
		Project:          c.projectID,
		Zone:             meta.ProviderZone,
		InstanceResource: &instance,
	})
	if err != nil {
		return nil, providers.NewAPIError("failed to initiate instance creation", err)
	}

	// Finally, we can create the promise
	return &InstancePromise{meta: meta, client: c, operation: operation}, nil
}

// Get implements the `providers.InstanceClient` interface.
func (c *Client) Get(
	ctx context.Context, meta providers.InstanceMeta,
) (providers.Instance, error) {
	// Before getting the instance, ensure that we manage the zone for which the instance is
	// requested
	manager, ok := c.instanceManagers[meta.ProviderZone]
	if !ok {
		return providers.Instance{}, providers.NewClientError(
			fmt.Sprintf("requested instance from unmanaged zone %q", meta.ProviderZone), nil,
		)
	}

	// Then, we can send the request
	instance, err := c.service.Get(ctx, &computepb.GetInstanceRequest{
		Project:  c.projectID,
		Instance: meta.ProviderID,
		Zone:     meta.ProviderZone,
	})
	if err != nil {
		return providers.Instance{},
			providers.NewAPIError(fmt.Sprintf("failed to obtain instance %q", meta.ID), err)
	}

	// And eventually unmarshal the returned instance into our own type
	return unmarshalInstance(instance, manager, c.projectID)
}

// List implements the `providers.InstanceClient` interface.
func (c *Client) List(ctx context.Context) ([]providers.Instance, error) {
	instances := []providers.Instance{}

	it := c.service.AggregatedList(ctx, &computepb.AggregatedListInstancesRequest{
		Project: c.projectID,
		Filter:  proto.String(fmt.Sprintf("labels.%s=\"%s\"", LabelKeyCreatedBy, c.identifier)),
	})
	if err := gcputils.Iterate[compute.InstancesScopedListPair](
		ctx, it, func(pair compute.InstancesScopedListPair) error {
			// Check for zone
			zone := path.Base(pair.Key)
			manager, ok := c.instanceManagers[zone]
			if !ok && len(pair.Value.GetInstances()) > 0 {
				return fmt.Errorf("found instances in zone %q which is not managed", zone)
			}

			// Parse instance
			for _, i := range pair.Value.GetInstances() {
				instance, err := unmarshalInstance(i, manager, c.projectID)
				if err != nil {
					return fmt.Errorf("failed to unmarshal instance %q: %s", i.GetName(), err)
				}
				instances = append(instances, instance)
			}
			return nil
		},
	); err != nil {
		return nil, providers.NewAPIError("failed to fetch all instances", err)
	}
	return instances, nil
}

// Delete implements the `providers.InstanceClient` interface.
func (c *Client) Delete(ctx context.Context, meta providers.InstanceMeta) error {
	operation, err := c.service.Delete(ctx, &computepb.DeleteInstanceRequest{
		Instance: meta.ProviderID,
		Project:  c.projectID,
		Zone:     meta.ProviderZone,
	})
	if err != nil {
		return providers.NewAPIError("failed to initiate instance deletion", err)
	}
	if err := operation.Wait(ctx); err != nil {
		return providers.NewAPIError("failed to await instance deletion", err)
	}
	if operation.Proto().Error != nil {
		return providers.NewClientError(
			fmt.Sprintf("instance deletion failed: %s", operation.Proto().Error), nil,
		)
	}
	return nil
}
