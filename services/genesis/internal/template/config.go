package template

// InstanceConfig defines a template for a specific type of instance (e.g. a "Minion instance").
// The template provides the instance's configuration on various cloud providers.
type InstanceConfig struct {
	CommonInstanceConfig `json:",inline"`
	// The instance template for launching this type of instance on Amazon Web Services. If this
	// template is not provided, the instance type may not be launched on AWS.
	Aws *AwsInstanceConfig `json:"aws,omitempty"`
	// The instance template for launching this type of instance on the Google Cloud Platform. If
	// this template is not provided, the instance type may not be launched on GCP.
	Gcp *GcpInstanceConfig `json:"gcp,omitempty"`
}

// CommonInstanceConfig describes instance configuration common across cloud providers.
type CommonInstanceConfig struct {
	// Optional reservations for the launched instance. If a request comes on for an instance with
	// a memory of size "x", an actual instance with memory "x + reserved memory" is launched. The
	// caller will be informed of an instance with memory "actual memory - reserved memory".
	Reservations InstanceReservations `json:"reservations,omitempty"`
	// Configuration of additional disks to attach to the instance.
	ExtraDisks []InstanceDisk `json:"extraDisks,omitempty"`
	// A set of additional metadata values to add to the instance. This metadata must not contain
	// secret values as it is attached to instances as tags.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// InstanceReservations provides reservations for system services for created instances.
type InstanceReservations struct {
	// The reserved memory size.
	Memory *string `json:"memory,omitempty"`
}

// InstanceDisk models the configuration of additional disks mounted to instances.
type InstanceDisk struct {
	// A common name for the disk. The disk will be available at different locations for the
	// different cloud providers.
	Name string `json:"name"`
	// The size of the disk per CPU on the instance.
	SizePerCPU string `json:"sizePerCpu"`
}

//-------------------------------------------------------------------------------------------------
// AMAZON WEB SERVICES
//-------------------------------------------------------------------------------------------------

// AwsConfig is a utility type for AWS instance configurations.
type AwsConfig struct {
	CommonInstanceConfig `json:",inline"`
	AwsInstanceConfig    `json:"aws"`
}

// AwsInstanceConfig describes an instance configuration specific to Amazon Web Services.
type AwsInstanceConfig struct {
	// The boot configuration of the instance.
	Boot AwsBootConfig `json:"boot"`
	// The network configuration of the instance.
	Network AwsNetworkConfig `json:"network"`
	// The IAM configuration of the instance.
	Iam AwsIamConfig `json:"iam"`
}

// AwsBootConfig describes the boot configuration for an AWS instance.
type AwsBootConfig struct {
	// The AMIs to use. When launching an instance, the options will be iterated and the AMI of
	// the first option whose selector matches the desired instance will be used.
	Amis []Option[AwsBootAmiConfig] `json:"amis"`
	// The size of the boot disk to use for all AMIs.
	DiskSize string `json:"diskSize"`
}

// AwsBootAmiConfig defines the AMI configuration of the AWS boot configuration.
type AwsBootAmiConfig struct {
	// The owner of the AMI (e.g. an account ID). This account ID is typically provided by
	// Taskfleet.
	Owner string `json:"owner"`
	// Tags that are used to filter the AMIs provided by the owner of the AMI. This can be used,
	// for example, to select a particular version of the AMI. Filtering AMIs owned by the
	// specified owner by these tags should yield *exactly one* AMI. Otherwise, the instance
	// manager will choose an AMI at random.
	Selector map[string]string `json:"selector"`
}

// AwsNetworkConfig defines configuration associated to the networking of a created AWS instance.
type AwsNetworkConfig struct {
	// The tags that ought to be used for selecting VPCs across regions. Instances may only be
	// launched into regions where there is *exactly one* VPC which provides the defined tags.
	// Further, instances can only be launched into availability zones where a VPC has an
	// associated subnet. If there is more than one subnet, `subnetTags` will be used.
	VpcSelector map[string]string `json:"vpcSelector"`
	// An optional set of tags that can be used to select a single subnet when a VPC provides
	// several subnets in an availability zone. For example, you may choose to select a subnet
	// based on its visibility (public/private).
	SubnetSelector map[string]string `json:"subnetSelector,omitempty"`
	// The tags to use for selecting security groups that ought to be attached to created
	// instances. The instance manager will use ALL security groups that have set the provided tags
	// AND belong to the VPC specified via `vpcId`.
	SecurityGroupSelector map[string]string `json:"securityGroupSelector"`
}

// AwsIamConfig describes the IAM configuration of instances launched on AWS.
type AwsIamConfig struct {
	// The name of the instance profile to use.
	InstanceProfile string `json:"instanceProfile"`
}

//-------------------------------------------------------------------------------------------------
// GOOGLE CLOUD PLATFORM
//-------------------------------------------------------------------------------------------------

// GcpConfig is a utility type for GCP instance configurations.
type GcpConfig struct {
	CommonInstanceConfig `json:",inline"`
	GcpInstanceConfig    `json:"gcp"`
}

// GcpInstanceConfig describes an instance configuration specific to the Google Cloud Platform.
type GcpInstanceConfig struct {
	// The instance's boot configuration.
	Boot GcpBootConfig `json:"boot"`
	// The network configuration of the instance.
	Network GcpNetworkConfig `json:"network"`
	// The IAM configuration of the instance.
	Iam GcpIamConfig `json:"iam"`
	// The configuration for disks attached to the instance.
	Disks GcpDiskConfig `json:"disk"`
}

// GcpBootConfig defines the boot configuration of a GCP instance.
type GcpBootConfig struct {
	// Links to the boot images to use. When launching an instance, the options will be iterated
	// and the image link of the first option whose selector matches the desired instance will be
	// used.
	ImageLink []Option[string] `json:"imageLinks"`
	// The size of the boot image. This size applies to all image links.
	DiskSize string `json:"diskSize"`
}

// GcpNetworkConfig defines configuration associated with the networking of instances launched on
// GCP.
type GcpNetworkConfig struct {
	// The name of the GCP network into which instances should be launched. Instances can only be
	// launched into regions where this network defines exactly one subnet.
	Name string `json:"name"`
	// The network tags to apply to the instance. These tags implicitly define the firewall rules
	// that are applied to the instance.
	Tags []string `json:"tags"`
	// Whether instances should be shielded from the public internet, i.e. should not receive an
	// external IP address. Instances can still access the public internet if a cloud router is
	// configured in their network and the region they are launched in.
	Shielded bool `json:"shielded"`
}

// GcpIamConfig defines the IAM configuration of instances launched on GCP.
type GcpIamConfig struct {
	// The email of the service account to attach to instances.
	ServiceAccountEmail string `json:"serviceAccountEmail"`
}

// GcpDiskConfig defines basic properties of disks attaches to instances launched on GCP.
type GcpDiskConfig struct {
	// The disk type to use for all disks. Defaults to pd-standard.
	Type string `json:"type"`
}
