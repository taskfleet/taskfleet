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
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// Client represents a GCP instance client.
type Client struct {
	projectID  string
	identifier string
	network    string
	service    *compute.InstancesClient
	zones      *gcpzones.Client
	instances  map[string]*instances.Manager
}

// NewClient initializes a new instance client which first fetch all available instance types.
func NewClient(
	ctx context.Context,
	identifier, network, projectID string,
	service *compute.InstancesClient,
	zones *gcpzones.Client,
) (*Client, error) {
	// First, we initialize the instance manager
	allInstances, err := findAvailableInstanceTypes(ctx, service, zones, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch available instances: %s", err)
	}
	instanceManagers := make(map[string]*instances.Manager)
	for zone, zoneInstances := range allInstances {
		instanceManagers[zone] = instances.NewManager(zoneInstances)
	}

	// Then, we can initialize the client
	return &Client{
		projectID:  projectID,
		identifier: identifier,
		network:    network,
		zones:      zones,
		instances:  instanceManagers,
	}, nil
}

//-------------------------------------------------------------------------------------------------
// INTERFACE
//-------------------------------------------------------------------------------------------------

// Find implements the `providers.InstanceClient` interface.
func (c *Client) Find(
	zone string, resources instances.Resources, architecture typedefs.CPUArchitecture,
) (instances.Type, error) {
	manager, ok := c.instances[zone]
	if !ok {
		return instances.Type{}, fmt.Errorf("failed to find instance for zone %q", zone)
	}
	return manager.FindBestFit(resources, architecture)
}

// Create implements the `providers.InstanceClient` interface.
func (c *Client) Create(
	ctx context.Context, meta providers.InstanceMeta, spec providers.InstanceSpec,
) (providers.InstancePromise, error) {
	// Accelerators
	guestAccelerators := []*compute.AcceleratorConfig{}
	if spec.InstanceType.GPU != nil {
		accelerator, err := c.zones.FindAccelerator(meta.ProviderID, spec.InstanceType.GPU.Kind)
		if err != nil {
			return nil, providers.NewClientError("accelerator does not exist in zone", err)
		}
		config, err := accelerator.Config(spec.InstanceType.GPU.Count)
		if err != nil {
			return nil, providers.NewClientError("failed to obtain requested number of GPUs", err)
		}
		guestAccelerators = []*compute.AcceleratorConfig{config}
	}

	// Labels
	labels := map[string]string{
		LabelKeyCreatedBy: c.identifier,
		LabelKeyOwnedBy:   "",
	}

	// Metadata
	metadata := []*compute.MetadataItems{}

	// Service account
	serviceAccounts := []*compute.ServiceAccount{}
	// if spec.Security.ServiceAccountEmail != "" {
	// 	serviceAccounts = append(serviceAccounts, &compute.ServiceAccount{
	// 		Email:  spec.Security.ServiceAccountEmail,
	// 		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	// 	})
	// }

	// Then, we initialize the instance
	hostMaintenanceAction := "MIGRATE"
	if spec.InstanceType.GPU != nil {
		hostMaintenanceAction = "TERMINATE"
	}
	provisioning := "STANDARD"
	if spec.IsSpot {
		provisioning = "SPOT"
	}

	instance := &compute.Instance{
		Name: meta.CommonName(),
		Description: fmt.Sprintf(
			"Genesis-managed instance owned by %q", labels[LabelKeyOwnedBy],
		),
		Labels:          labels,
		Tags:            &compute.Tags{},
		Metadata:        &compute.Metadata{Items: metadata},
		ServiceAccounts: serviceAccounts,
		Scheduling: &compute.Scheduling{
			ProvisioningModel: provisioning,
			OnHostMaintenance: hostMaintenanceAction,
		},
		MachineType:       c.machineTypeLink(meta.ProviderZone, spec.InstanceType),
		GuestAccelerators: guestAccelerators,
		NetworkInterfaces: []*compute.NetworkInterface{{
			AccessConfigs: []*compute.AccessConfig{{
				Name: "External NAT",
				Type: "ONE_TO_ONE_NAT",
			}},
			Network:    c.networkLink(),
			Subnetwork: c.subnetworkLink(meta),
		}},
		Disks: c.gcpDisks(ref, spec),
	}

	// And actually make the call to create it
	call := c.Service.Instances.Insert(c.Project, ref.ProviderZone, instance)
	operation, err := c.GetOperation(ctx, call.Context(ctx))
	if err != nil {
		return nil, providers.NewAPIError("failed initiating instance creation", err)
	}

	// Finally, we can create the promise
	return &InstancePromise{Ref: ref, client: c, operation: operation}, nil
}

// Get implements the `providers.InstanceClient` interface.
func (c *Client) Get(ctx context.Context, ref providers.InstanceMeta) (providers.Instance, error) {
	call := c.Service.Instances.Get(c.Project, ref.ProviderZone, ref.CommonName())
	i, err := gcputils.ExecuteWithRetry(ctx, func() (*compute.Instance, error) {
		return call.Context(ctx).Do()
	})
	if err != nil {
		return providers.Instance{},
			providers.NewAPIError("failed obtaining created instance", err)
	}
	return parseInstance(i, c.instances, c.Project)
}

// List implements the `providers.InstanceClient` interface.
func (c *Client) List(ctx context.Context) ([]providers.Instance, error) {
	instances := []providers.Instance{}

	call := c.Service.Instances.AggregatedList(c.Project)
	filter := fmt.Sprintf("labels.created-by=\"%s\"", c.identifier)
	err := call.Filter(filter).Pages(ctx, func(list *compute.InstanceAggregatedList) error {
		for _, item := range list.Items {
			for _, i := range item.Instances {
				instance, err := parseInstance(i, c.instances, c.Project)
				if err != nil {
					return err
				}
				instances = append(instances, instance)
			}
		}
		return nil
	})
	if err != nil {
		return nil, providers.NewAPIError("failed to fetch all instances", err)
	}

	return instances, nil
}

// Delete implements the `providers.InstanceClient` interface.
func (c *Client) Delete(ctx context.Context, ref providers.InstanceMeta) error {
	call := c.Service.Instances.Delete(c.Project, ref.ProviderZone, ref.CommonName())
	operation, err := c.GetOperation(ctx, call.Context(ctx))
	if err != nil {
		return providers.NewAPIError("failed initiating delete", err)
	}
	if err := c.PollOperation(ctx, operation); err != nil {
		return providers.NewAPIError("failed to finish deletion", err)
	}
	return nil
}

//-------------------------------------------------------------------------------------------------

func (c *Client) gcpDisks(
	ref providers.InstanceMeta, spec providers.InstanceSpec,
) []*compute.AttachedDisk {
	// First, obtain boot disk
	disks := []*compute.AttachedDisk{{
		AutoDelete: true,
		Boot:       true,
		Mode:       "READ_WRITE",
		InitializeParams: &compute.AttachedDiskInitializeParams{
			DiskName:    fmt.Sprintf("%s-boot-disk", ref.CommonName()),
			DiskSizeGb:  int64(spec.Boot.DiskSizeGiB),
			DiskType:    c.diskTypeLink(ref.ProviderZone, typedefs.DiskHDD),
			SourceImage: spec.Boot.ImageLink,
			Labels: map[string]string{
				LabelKeyCreatedBy: c.identifier,
			},
		},
	}}

	// Then, obtain additional disks
	for i, disk := range spec.Disks {
		disks = append(disks, &compute.AttachedDisk{
			AutoDelete: true,
			Boot:       false,
			DeviceName: disk.Name,
			Mode:       "READ_WRITE",
			InitializeParams: &compute.AttachedDiskInitializeParams{
				DiskName:   fmt.Sprintf("%s-disk-%d", ref.CommonName(), i),
				DiskSizeGb: int64(disk.SizeGiB),
				DiskType:   c.diskTypeLink(ref.ProviderZone, disk.Type),
				Labels: map[string]string{
					LabelKeyCreatedBy:  c.identifier,
					LabelKeyDeviceName: disk.Name,
				},
			},
		})
	}

	return disks
}

//-------------------------------------------------------------------------------------------------
// UTILS
//-------------------------------------------------------------------------------------------------

func (c *Client) networkLink() string {
	return fmt.Sprintf("projects/%s/global/networks/%s", c.projectID, c.network)
}

func (c *Client) subnetworkLink(meta providers.InstanceMeta) string {
	return fmt.Sprintf(
		"regions/%s/subnetworks/%s-%s",
		gcputils.RegionFromZone(meta.ProviderZone),
		path.Base(c.network),
		gcputils.RegionFromZone(meta.ProviderZone),
	)
}

func (c *Client) machineTypeLink(zone string, instanceType instances.Type) string {
	return fmt.Sprintf(
		"projects/%s/zones/%s/machineTypes/%s", c.projectID, zone, instanceType.Name,
	)
}

// func (c *Client) diskTypeLink(zone string, diskType typedefs.DiskType) string {
// 	return fmt.Sprintf("projects/%s/zones/%s/diskTypes/%s", c.Project, zone, gcpDiskType(diskType))
// }
