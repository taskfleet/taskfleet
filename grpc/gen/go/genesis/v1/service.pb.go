// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: genesis/v1/service.proto

package genesis

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListZonesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ListZonesRequest) Reset() {
	*x = ListZonesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListZonesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListZonesRequest) ProtoMessage() {}

func (x *ListZonesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListZonesRequest.ProtoReflect.Descriptor instead.
func (*ListZonesRequest) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{0}
}

type ListZonesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A list of available zones.
	Zones []*Zone `protobuf:"bytes,1,rep,name=zones,proto3" json:"zones,omitempty"`
}

func (x *ListZonesResponse) Reset() {
	*x = ListZonesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListZonesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListZonesResponse) ProtoMessage() {}

func (x *ListZonesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListZonesResponse.ProtoReflect.Descriptor instead.
func (*ListZonesResponse) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{1}
}

func (x *ListZonesResponse) GetZones() []*Zone {
	if x != nil {
		return x.Zones
	}
	return nil
}

type Zone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The provider to which this zone belongs.
	Provider CloudProvider `protobuf:"varint,1,opt,name=provider,proto3,enum=genesis.v1.CloudProvider" json:"provider,omitempty"`
	// The provider-specific name of the zone.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// The GPUs available in this zone.
	AvailableGpus []GPUKind `protobuf:"varint,3,rep,packed,name=available_gpus,json=availableGpus,proto3,enum=genesis.v1.GPUKind" json:"available_gpus,omitempty"`
}

func (x *Zone) Reset() {
	*x = Zone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Zone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Zone) ProtoMessage() {}

func (x *Zone) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Zone.ProtoReflect.Descriptor instead.
func (*Zone) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{2}
}

func (x *Zone) GetProvider() CloudProvider {
	if x != nil {
		return x.Provider
	}
	return CloudProvider_CLOUD_PROVIDER_UNSPECIFIED
}

func (x *Zone) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Zone) GetAvailableGpus() []GPUKind {
	if x != nil {
		return x.AvailableGpus
	}
	return nil
}

type CreateInstanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique ID of the instance to create. The ID shall be generated by the client to make it
	// easier to link created instances to those delivered via Kafka.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// An arbitrary non-empty string to identify the caller. It can be reused to fetch all
	// instances that were created by the caller. Should thus be a globally unique string.
	Owner string `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	// The component for which to create an instance. An error is returned if no configuration for
	// the component can be found.
	Component string `protobuf:"bytes,3,opt,name=component,proto3" json:"component,omitempty"`
	// The desired configuration of the instance.
	Config *InstanceConfig `protobuf:"bytes,4,opt,name=config,proto3" json:"config,omitempty"`
	// The desired amount of resources on the instance.
	Resources *InstanceResources `protobuf:"bytes,5,opt,name=resources,proto3" json:"resources,omitempty"`
	// Whether the instance is required for high-performance computing. In that case, compute-
	// optimized instances are preferably created.
	PreferHpc bool `protobuf:"varint,6,opt,name=prefer_hpc,json=preferHpc,proto3" json:"prefer_hpc,omitempty"`
}

func (x *CreateInstanceRequest) Reset() {
	*x = CreateInstanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateInstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateInstanceRequest) ProtoMessage() {}

func (x *CreateInstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateInstanceRequest.ProtoReflect.Descriptor instead.
func (*CreateInstanceRequest) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{3}
}

func (x *CreateInstanceRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateInstanceRequest) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *CreateInstanceRequest) GetComponent() string {
	if x != nil {
		return x.Component
	}
	return ""
}

func (x *CreateInstanceRequest) GetConfig() *InstanceConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *CreateInstanceRequest) GetResources() *InstanceResources {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *CreateInstanceRequest) GetPreferHpc() bool {
	if x != nil {
		return x.PreferHpc
	}
	return false
}

type CreateInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique reference to the instance that will be created.
	Instance *Instance `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	// The actual configuration of the instance. Echo'ed from the request.
	Config *InstanceConfig `protobuf:"bytes,2,opt,name=config,proto3" json:"config,omitempty"`
	// The available resources on the instance.
	Resources *InstanceResources `protobuf:"bytes,3,opt,name=resources,proto3" json:"resources,omitempty"`
}

func (x *CreateInstanceResponse) Reset() {
	*x = CreateInstanceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateInstanceResponse) ProtoMessage() {}

func (x *CreateInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateInstanceResponse.ProtoReflect.Descriptor instead.
func (*CreateInstanceResponse) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{4}
}

func (x *CreateInstanceResponse) GetInstance() *Instance {
	if x != nil {
		return x.Instance
	}
	return nil
}

func (x *CreateInstanceResponse) GetConfig() *InstanceConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *CreateInstanceResponse) GetResources() *InstanceResources {
	if x != nil {
		return x.Resources
	}
	return nil
}

type ListInstancesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the instances' owner, i.e. the component having created the instances. Should
	// coincide with the `owner` string passed when creating instances.
	Owner string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (x *ListInstancesRequest) Reset() {
	*x = ListInstancesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListInstancesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInstancesRequest) ProtoMessage() {}

func (x *ListInstancesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInstancesRequest.ProtoReflect.Descriptor instead.
func (*ListInstancesRequest) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{5}
}

func (x *ListInstancesRequest) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

type ListInstancesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The instances owned by the owner specified in the request which are currently running.
	Instances []*RunningInstance `protobuf:"bytes,1,rep,name=instances,proto3" json:"instances,omitempty"`
}

func (x *ListInstancesResponse) Reset() {
	*x = ListInstancesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListInstancesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListInstancesResponse) ProtoMessage() {}

func (x *ListInstancesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListInstancesResponse.ProtoReflect.Descriptor instead.
func (*ListInstancesResponse) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{6}
}

func (x *ListInstancesResponse) GetInstances() []*RunningInstance {
	if x != nil {
		return x.Instances
	}
	return nil
}

type RunningInstance struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The unique reference to the instance.
	Instance *Instance `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	// The component for which the instance was created.
	Component string `protobuf:"bytes,2,opt,name=component,proto3" json:"component,omitempty"`
	// The configuration of the instance.
	Config *InstanceConfig `protobuf:"bytes,3,opt,name=config,proto3" json:"config,omitempty"`
	// The available resources on the instance.
	Resources *InstanceResources `protobuf:"bytes,4,opt,name=resources,proto3" json:"resources,omitempty"`
	// The hostname of the instance.
	Hostname string `protobuf:"bytes,5,opt,name=hostname,proto3" json:"hostname,omitempty"`
}

func (x *RunningInstance) Reset() {
	*x = RunningInstance{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RunningInstance) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RunningInstance) ProtoMessage() {}

func (x *RunningInstance) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RunningInstance.ProtoReflect.Descriptor instead.
func (*RunningInstance) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{7}
}

func (x *RunningInstance) GetInstance() *Instance {
	if x != nil {
		return x.Instance
	}
	return nil
}

func (x *RunningInstance) GetComponent() string {
	if x != nil {
		return x.Component
	}
	return ""
}

func (x *RunningInstance) GetConfig() *InstanceConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *RunningInstance) GetResources() *InstanceResources {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *RunningInstance) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

type ShutdownInstanceRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The instance to shut down.
	Instance *Instance `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
}

func (x *ShutdownInstanceRequest) Reset() {
	*x = ShutdownInstanceRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShutdownInstanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownInstanceRequest) ProtoMessage() {}

func (x *ShutdownInstanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownInstanceRequest.ProtoReflect.Descriptor instead.
func (*ShutdownInstanceRequest) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{8}
}

func (x *ShutdownInstanceRequest) GetInstance() *Instance {
	if x != nil {
		return x.Instance
	}
	return nil
}

type ShutdownInstanceResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ShutdownInstanceResponse) Reset() {
	*x = ShutdownInstanceResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_v1_service_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShutdownInstanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShutdownInstanceResponse) ProtoMessage() {}

func (x *ShutdownInstanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_v1_service_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShutdownInstanceResponse.ProtoReflect.Descriptor instead.
func (*ShutdownInstanceResponse) Descriptor() ([]byte, []int) {
	return file_genesis_v1_service_proto_rawDescGZIP(), []int{9}
}

var File_genesis_v1_service_proto protoreflect.FileDescriptor

var file_genesis_v1_service_proto_rawDesc = []byte{
	0x0a, 0x18, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x67, 0x65, 0x6e, 0x65,
	0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2f,
	0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x12, 0x0a, 0x10, 0x4c, 0x69, 0x73, 0x74, 0x5a,
	0x6f, 0x6e, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x3b, 0x0a, 0x11, 0x4c,
	0x69, 0x73, 0x74, 0x5a, 0x6f, 0x6e, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x26, 0x0a, 0x05, 0x7a, 0x6f, 0x6e, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x10, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x5a, 0x6f, 0x6e,
	0x65, 0x52, 0x05, 0x7a, 0x6f, 0x6e, 0x65, 0x73, 0x22, 0x8d, 0x01, 0x0a, 0x04, 0x5a, 0x6f, 0x6e,
	0x65, 0x12, 0x35, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x52, 0x08,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x3a, 0x0a, 0x0e,
	0x61, 0x76, 0x61, 0x69, 0x6c, 0x61, 0x62, 0x6c, 0x65, 0x5f, 0x67, 0x70, 0x75, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x47, 0x50, 0x55, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x0d, 0x61, 0x76, 0x61, 0x69, 0x6c,
	0x61, 0x62, 0x6c, 0x65, 0x47, 0x70, 0x75, 0x73, 0x22, 0x9b, 0x02, 0x0a, 0x15, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x72, 0x03, 0xb0, 0x01, 0x01, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x05,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x09, 0x63,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07,
	0xfa, 0x42, 0x04, 0x72, 0x02, 0x10, 0x01, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
	0x6e, 0x74, 0x12, 0x3c, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x42, 0x08,
	0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67,
	0x12, 0x45, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x73, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x09, 0x72, 0x65,
	0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x65, 0x66, 0x65,
	0x72, 0x5f, 0x68, 0x70, 0x63, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08, 0x52, 0x09, 0x70, 0x72, 0x65,
	0x66, 0x65, 0x72, 0x48, 0x70, 0x63, 0x22, 0xbb, 0x01, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x30, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3b, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x65, 0x6e,
	0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x22, 0x35, 0x0a, 0x14, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x05,
	0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x07, 0xfa, 0x42, 0x04,
	0x72, 0x02, 0x10, 0x01, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x22, 0x52, 0x0a, 0x15, 0x4c,
	0x69, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x39, 0x0a, 0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69,
	0x73, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x09, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x22,
	0xee, 0x01, 0x0a, 0x0f, 0x52, 0x75, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x12, 0x30, 0x0a, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65,
	0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e,
	0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3b, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x65, 0x6e,
	0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0x55, 0x0a, 0x17, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3a, 0x0a, 0x08, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x8a, 0x01, 0x02, 0x10, 0x01, 0x52, 0x08, 0x69,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x22, 0x1a, 0x0a, 0x18, 0x53, 0x68, 0x75, 0x74, 0x64,
	0x6f, 0x77, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xe8, 0x02, 0x0a, 0x0e, 0x47, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x48, 0x0a, 0x09, 0x4c, 0x69, 0x73, 0x74, 0x5a, 0x6f,
	0x6e, 0x65, 0x73, 0x12, 0x1c, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x5a, 0x6f, 0x6e, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c,
	0x69, 0x73, 0x74, 0x5a, 0x6f, 0x6e, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x57, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x12, 0x21, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x22, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x0d, 0x4c, 0x69, 0x73,
	0x74, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x67, 0x65, 0x6e,
	0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e, 0x73, 0x74,
	0x61, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x5d, 0x0a, 0x10, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x12, 0x23, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x24, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73,
	0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x68, 0x75, 0x74, 0x64, 0x6f, 0x77, 0x6e, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x30,
	0x5a, 0x2e, 0x67, 0x6f, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x66, 0x6c, 0x65, 0x65, 0x74, 0x2e, 0x69,
	0x6f, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67, 0x6f, 0x2f, 0x67, 0x65,
	0x6e, 0x65, 0x73, 0x69, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_genesis_v1_service_proto_rawDescOnce sync.Once
	file_genesis_v1_service_proto_rawDescData = file_genesis_v1_service_proto_rawDesc
)

func file_genesis_v1_service_proto_rawDescGZIP() []byte {
	file_genesis_v1_service_proto_rawDescOnce.Do(func() {
		file_genesis_v1_service_proto_rawDescData = protoimpl.X.CompressGZIP(file_genesis_v1_service_proto_rawDescData)
	})
	return file_genesis_v1_service_proto_rawDescData
}

var file_genesis_v1_service_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_genesis_v1_service_proto_goTypes = []interface{}{
	(*ListZonesRequest)(nil),         // 0: genesis.v1.ListZonesRequest
	(*ListZonesResponse)(nil),        // 1: genesis.v1.ListZonesResponse
	(*Zone)(nil),                     // 2: genesis.v1.Zone
	(*CreateInstanceRequest)(nil),    // 3: genesis.v1.CreateInstanceRequest
	(*CreateInstanceResponse)(nil),   // 4: genesis.v1.CreateInstanceResponse
	(*ListInstancesRequest)(nil),     // 5: genesis.v1.ListInstancesRequest
	(*ListInstancesResponse)(nil),    // 6: genesis.v1.ListInstancesResponse
	(*RunningInstance)(nil),          // 7: genesis.v1.RunningInstance
	(*ShutdownInstanceRequest)(nil),  // 8: genesis.v1.ShutdownInstanceRequest
	(*ShutdownInstanceResponse)(nil), // 9: genesis.v1.ShutdownInstanceResponse
	(CloudProvider)(0),               // 10: genesis.v1.CloudProvider
	(GPUKind)(0),                     // 11: genesis.v1.GPUKind
	(*InstanceConfig)(nil),           // 12: genesis.v1.InstanceConfig
	(*InstanceResources)(nil),        // 13: genesis.v1.InstanceResources
	(*Instance)(nil),                 // 14: genesis.v1.Instance
}
var file_genesis_v1_service_proto_depIdxs = []int32{
	2,  // 0: genesis.v1.ListZonesResponse.zones:type_name -> genesis.v1.Zone
	10, // 1: genesis.v1.Zone.provider:type_name -> genesis.v1.CloudProvider
	11, // 2: genesis.v1.Zone.available_gpus:type_name -> genesis.v1.GPUKind
	12, // 3: genesis.v1.CreateInstanceRequest.config:type_name -> genesis.v1.InstanceConfig
	13, // 4: genesis.v1.CreateInstanceRequest.resources:type_name -> genesis.v1.InstanceResources
	14, // 5: genesis.v1.CreateInstanceResponse.instance:type_name -> genesis.v1.Instance
	12, // 6: genesis.v1.CreateInstanceResponse.config:type_name -> genesis.v1.InstanceConfig
	13, // 7: genesis.v1.CreateInstanceResponse.resources:type_name -> genesis.v1.InstanceResources
	7,  // 8: genesis.v1.ListInstancesResponse.instances:type_name -> genesis.v1.RunningInstance
	14, // 9: genesis.v1.RunningInstance.instance:type_name -> genesis.v1.Instance
	12, // 10: genesis.v1.RunningInstance.config:type_name -> genesis.v1.InstanceConfig
	13, // 11: genesis.v1.RunningInstance.resources:type_name -> genesis.v1.InstanceResources
	14, // 12: genesis.v1.ShutdownInstanceRequest.instance:type_name -> genesis.v1.Instance
	0,  // 13: genesis.v1.GenesisService.ListZones:input_type -> genesis.v1.ListZonesRequest
	3,  // 14: genesis.v1.GenesisService.CreateInstance:input_type -> genesis.v1.CreateInstanceRequest
	5,  // 15: genesis.v1.GenesisService.ListInstances:input_type -> genesis.v1.ListInstancesRequest
	8,  // 16: genesis.v1.GenesisService.ShutdownInstance:input_type -> genesis.v1.ShutdownInstanceRequest
	1,  // 17: genesis.v1.GenesisService.ListZones:output_type -> genesis.v1.ListZonesResponse
	4,  // 18: genesis.v1.GenesisService.CreateInstance:output_type -> genesis.v1.CreateInstanceResponse
	6,  // 19: genesis.v1.GenesisService.ListInstances:output_type -> genesis.v1.ListInstancesResponse
	9,  // 20: genesis.v1.GenesisService.ShutdownInstance:output_type -> genesis.v1.ShutdownInstanceResponse
	17, // [17:21] is the sub-list for method output_type
	13, // [13:17] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_genesis_v1_service_proto_init() }
func file_genesis_v1_service_proto_init() {
	if File_genesis_v1_service_proto != nil {
		return
	}
	file_genesis_v1_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_genesis_v1_service_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListZonesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListZonesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Zone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateInstanceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateInstanceResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListInstancesRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListInstancesResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RunningInstance); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShutdownInstanceRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_genesis_v1_service_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShutdownInstanceResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_genesis_v1_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_genesis_v1_service_proto_goTypes,
		DependencyIndexes: file_genesis_v1_service_proto_depIdxs,
		MessageInfos:      file_genesis_v1_service_proto_msgTypes,
	}.Build()
	File_genesis_v1_service_proto = out.File
	file_genesis_v1_service_proto_rawDesc = nil
	file_genesis_v1_service_proto_goTypes = nil
	file_genesis_v1_service_proto_depIdxs = nil
}
