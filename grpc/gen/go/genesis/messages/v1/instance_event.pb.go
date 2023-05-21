// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: genesis/messages/v1/instance_event.proto

package genesis_messages

import (
	v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InstanceCreationFailedEvent_Reason int32

const (
	InstanceCreationFailedEvent_REASON_UNSPECIFIED InstanceCreationFailedEvent_Reason = 0
	// The instance was never started up as insufficient resources were available.
	InstanceCreationFailedEvent_REASON_INSUFFICIENT_RESOURCES InstanceCreationFailedEvent_Reason = 1
	// The instance was never started up as insufficient quotas were available.
	InstanceCreationFailedEvent_REASON_QUOTA_EXCEEDED InstanceCreationFailedEvent_Reason = 2
)

// Enum value maps for InstanceCreationFailedEvent_Reason.
var (
	InstanceCreationFailedEvent_Reason_name = map[int32]string{
		0: "REASON_UNSPECIFIED",
		1: "REASON_INSUFFICIENT_RESOURCES",
		2: "REASON_QUOTA_EXCEEDED",
	}
	InstanceCreationFailedEvent_Reason_value = map[string]int32{
		"REASON_UNSPECIFIED":            0,
		"REASON_INSUFFICIENT_RESOURCES": 1,
		"REASON_QUOTA_EXCEEDED":         2,
	}
)

func (x InstanceCreationFailedEvent_Reason) Enum() *InstanceCreationFailedEvent_Reason {
	p := new(InstanceCreationFailedEvent_Reason)
	*p = x
	return p
}

func (x InstanceCreationFailedEvent_Reason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InstanceCreationFailedEvent_Reason) Descriptor() protoreflect.EnumDescriptor {
	return file_genesis_messages_v1_instance_event_proto_enumTypes[0].Descriptor()
}

func (InstanceCreationFailedEvent_Reason) Type() protoreflect.EnumType {
	return &file_genesis_messages_v1_instance_event_proto_enumTypes[0]
}

func (x InstanceCreationFailedEvent_Reason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InstanceCreationFailedEvent_Reason.Descriptor instead.
func (InstanceCreationFailedEvent_Reason) EnumDescriptor() ([]byte, []int) {
	return file_genesis_messages_v1_instance_event_proto_rawDescGZIP(), []int{2, 0}
}

type InstanceDeletedEvent_Reason int32

const (
	InstanceDeletedEvent_REASON_UNSPECIFIED InstanceDeletedEvent_Reason = 0
	// The instance was shut down gracefully as a consequence of a user's request.
	InstanceDeletedEvent_REASON_SHUTDOWN InstanceDeletedEvent_Reason = 1
	// The instance was terminated forcefully as it was deemed unhealthy.
	InstanceDeletedEvent_REASON_UNHEALTHY InstanceDeletedEvent_Reason = 2
	// The instance was terminated by the cloud provider.
	InstanceDeletedEvent_REASON_TERMINATED InstanceDeletedEvent_Reason = 3
)

// Enum value maps for InstanceDeletedEvent_Reason.
var (
	InstanceDeletedEvent_Reason_name = map[int32]string{
		0: "REASON_UNSPECIFIED",
		1: "REASON_SHUTDOWN",
		2: "REASON_UNHEALTHY",
		3: "REASON_TERMINATED",
	}
	InstanceDeletedEvent_Reason_value = map[string]int32{
		"REASON_UNSPECIFIED": 0,
		"REASON_SHUTDOWN":    1,
		"REASON_UNHEALTHY":   2,
		"REASON_TERMINATED":  3,
	}
)

func (x InstanceDeletedEvent_Reason) Enum() *InstanceDeletedEvent_Reason {
	p := new(InstanceDeletedEvent_Reason)
	*p = x
	return p
}

func (x InstanceDeletedEvent_Reason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InstanceDeletedEvent_Reason) Descriptor() protoreflect.EnumDescriptor {
	return file_genesis_messages_v1_instance_event_proto_enumTypes[1].Descriptor()
}

func (InstanceDeletedEvent_Reason) Type() protoreflect.EnumType {
	return &file_genesis_messages_v1_instance_event_proto_enumTypes[1]
}

func (x InstanceDeletedEvent_Reason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InstanceDeletedEvent_Reason.Descriptor instead.
func (InstanceDeletedEvent_Reason) EnumDescriptor() ([]byte, []int) {
	return file_genesis_messages_v1_instance_event_proto_rawDescGZIP(), []int{3, 0}
}

// InstanceEvent encapsulates the data to describe a lifetime event of a cloud instance.
// * Key: Globally unique instance ID
type InstanceEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The instance that the event refers to.
	Instance *v1.Instance `protobuf:"bytes,1,opt,name=instance,proto3" json:"instance,omitempty"`
	// The timestamp of the event.
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Types that are assignable to Event:
	//
	//	*InstanceEvent_Created
	//	*InstanceEvent_CreationFailed
	//	*InstanceEvent_Deleted
	Event isInstanceEvent_Event `protobuf_oneof:"event"`
}

func (x *InstanceEvent) Reset() {
	*x = InstanceEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceEvent) ProtoMessage() {}

func (x *InstanceEvent) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceEvent.ProtoReflect.Descriptor instead.
func (*InstanceEvent) Descriptor() ([]byte, []int) {
	return file_genesis_messages_v1_instance_event_proto_rawDescGZIP(), []int{0}
}

func (x *InstanceEvent) GetInstance() *v1.Instance {
	if x != nil {
		return x.Instance
	}
	return nil
}

func (x *InstanceEvent) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (m *InstanceEvent) GetEvent() isInstanceEvent_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *InstanceEvent) GetCreated() *InstanceCreatedEvent {
	if x, ok := x.GetEvent().(*InstanceEvent_Created); ok {
		return x.Created
	}
	return nil
}

func (x *InstanceEvent) GetCreationFailed() *InstanceCreationFailedEvent {
	if x, ok := x.GetEvent().(*InstanceEvent_CreationFailed); ok {
		return x.CreationFailed
	}
	return nil
}

func (x *InstanceEvent) GetDeleted() *InstanceDeletedEvent {
	if x, ok := x.GetEvent().(*InstanceEvent_Deleted); ok {
		return x.Deleted
	}
	return nil
}

type isInstanceEvent_Event interface {
	isInstanceEvent_Event()
}

type InstanceEvent_Created struct {
	// The event indicates that an instance was created.
	Created *InstanceCreatedEvent `protobuf:"bytes,3,opt,name=created,proto3,oneof"`
}

type InstanceEvent_CreationFailed struct {
	// The event indicates that an instance failed to be created.
	CreationFailed *InstanceCreationFailedEvent `protobuf:"bytes,4,opt,name=creation_failed,json=creationFailed,proto3,oneof"`
}

type InstanceEvent_Deleted struct {
	// The event indicates that a running instance was deleted.
	Deleted *InstanceDeletedEvent `protobuf:"bytes,5,opt,name=deleted,proto3,oneof"`
}

func (*InstanceEvent_Created) isInstanceEvent_Event() {}

func (*InstanceEvent_CreationFailed) isInstanceEvent_Event() {}

func (*InstanceEvent_Deleted) isInstanceEvent_Event() {}

// InstanceCreatedEvent wraps information about an instance when the instance was created.
type InstanceCreatedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The configuration of the created instance.
	Config *v1.InstanceConfig `protobuf:"bytes,1,opt,name=config,proto3" json:"config,omitempty"`
	// The available resources of the created instance.
	Resources *v1.InstanceResources `protobuf:"bytes,2,opt,name=resources,proto3" json:"resources,omitempty"`
	// The hostname of the created instance.
	Hostname string `protobuf:"bytes,3,opt,name=hostname,proto3" json:"hostname,omitempty"`
}

func (x *InstanceCreatedEvent) Reset() {
	*x = InstanceCreatedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceCreatedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceCreatedEvent) ProtoMessage() {}

func (x *InstanceCreatedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceCreatedEvent.ProtoReflect.Descriptor instead.
func (*InstanceCreatedEvent) Descriptor() ([]byte, []int) {
	return file_genesis_messages_v1_instance_event_proto_rawDescGZIP(), []int{1}
}

func (x *InstanceCreatedEvent) GetConfig() *v1.InstanceConfig {
	if x != nil {
		return x.Config
	}
	return nil
}

func (x *InstanceCreatedEvent) GetResources() *v1.InstanceResources {
	if x != nil {
		return x.Resources
	}
	return nil
}

func (x *InstanceCreatedEvent) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

// InstanceCreationFailedEvent wraps information about an instance that failed to start up.
type InstanceCreationFailedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The reason why instance creation failed.
	Reason InstanceCreationFailedEvent_Reason `protobuf:"varint,1,opt,name=reason,proto3,enum=genesis.messages.v1.InstanceCreationFailedEvent_Reason" json:"reason,omitempty"`
	// A message that provides more details on the failure.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *InstanceCreationFailedEvent) Reset() {
	*x = InstanceCreationFailedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceCreationFailedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceCreationFailedEvent) ProtoMessage() {}

func (x *InstanceCreationFailedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceCreationFailedEvent.ProtoReflect.Descriptor instead.
func (*InstanceCreationFailedEvent) Descriptor() ([]byte, []int) {
	return file_genesis_messages_v1_instance_event_proto_rawDescGZIP(), []int{2}
}

func (x *InstanceCreationFailedEvent) GetReason() InstanceCreationFailedEvent_Reason {
	if x != nil {
		return x.Reason
	}
	return InstanceCreationFailedEvent_REASON_UNSPECIFIED
}

func (x *InstanceCreationFailedEvent) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

// InstanceDeletedEvent wraps information about an instance when the instance was deleted.
type InstanceDeletedEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The reason why the instance was deleted.
	Reason InstanceDeletedEvent_Reason `protobuf:"varint,1,opt,name=reason,proto3,enum=genesis.messages.v1.InstanceDeletedEvent_Reason" json:"reason,omitempty"`
}

func (x *InstanceDeletedEvent) Reset() {
	*x = InstanceDeletedEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InstanceDeletedEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InstanceDeletedEvent) ProtoMessage() {}

func (x *InstanceDeletedEvent) ProtoReflect() protoreflect.Message {
	mi := &file_genesis_messages_v1_instance_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InstanceDeletedEvent.ProtoReflect.Descriptor instead.
func (*InstanceDeletedEvent) Descriptor() ([]byte, []int) {
	return file_genesis_messages_v1_instance_event_proto_rawDescGZIP(), []int{3}
}

func (x *InstanceDeletedEvent) GetReason() InstanceDeletedEvent_Reason {
	if x != nil {
		return x.Reason
	}
	return InstanceDeletedEvent_REASON_UNSPECIFIED
}

var File_genesis_messages_v1_instance_event_proto protoreflect.FileDescriptor

var file_genesis_messages_v1_instance_event_proto_rawDesc = []byte{
	0x0a, 0x28, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x5f, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x67, 0x65, 0x6e, 0x65,
	0x73, 0x69, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x1a,
	0x16, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2f, 0x76, 0x31, 0x2f, 0x74, 0x79, 0x70, 0x65,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef, 0x02, 0x0a, 0x0d, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x30, 0x0a, 0x08, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x63, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x09,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x45, 0x0a, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69,
	0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x48, 0x00, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x5b, 0x0a,
	0x0f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73,
	0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73,
	0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69,
	0x6c, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x12, 0x45, 0x0a, 0x07, 0x64, 0x65,
	0x6c, 0x65, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x65,
	0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76,
	0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x07, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x64, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0xa3, 0x01, 0x0a, 0x14, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52,
	0x06, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3b, 0x0a, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x67, 0x65, 0x6e,
	0x65, 0x73, 0x69, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x73, 0x52, 0x09, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x22, 0xe8, 0x01, 0x0a, 0x1b, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x12, 0x4f, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x37, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x2e, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x5e, 0x0a, 0x06, 0x52,
	0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f,
	0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x21, 0x0a,
	0x1d, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x49, 0x4e, 0x53, 0x55, 0x46, 0x46, 0x49, 0x43,
	0x49, 0x45, 0x4e, 0x54, 0x5f, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x53, 0x10, 0x01,
	0x12, 0x19, 0x0a, 0x15, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x51, 0x55, 0x4f, 0x54, 0x41,
	0x5f, 0x45, 0x58, 0x43, 0x45, 0x45, 0x44, 0x45, 0x44, 0x10, 0x02, 0x22, 0xc4, 0x01, 0x0a, 0x14,
	0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x48, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x30, 0x2e, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2e, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x76, 0x31, 0x2e, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e,
	0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x22, 0x62,
	0x0a, 0x06, 0x52, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x45, 0x41, 0x53,
	0x4f, 0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x13, 0x0a, 0x0f, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x53, 0x48, 0x55, 0x54, 0x44,
	0x4f, 0x57, 0x4e, 0x10, 0x01, 0x12, 0x14, 0x0a, 0x10, 0x52, 0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f,
	0x55, 0x4e, 0x48, 0x45, 0x41, 0x4c, 0x54, 0x48, 0x59, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11, 0x52,
	0x45, 0x41, 0x53, 0x4f, 0x4e, 0x5f, 0x54, 0x45, 0x52, 0x4d, 0x49, 0x4e, 0x41, 0x54, 0x45, 0x44,
	0x10, 0x03, 0x42, 0x42, 0x5a, 0x40, 0x67, 0x6f, 0x2e, 0x74, 0x61, 0x73, 0x6b, 0x66, 0x6c, 0x65,
	0x65, 0x74, 0x2e, 0x69, 0x6f, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x67,
	0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x73, 0x2f, 0x76, 0x31, 0x3b, 0x67, 0x65, 0x6e, 0x65, 0x73, 0x69, 0x73, 0x5f, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_genesis_messages_v1_instance_event_proto_rawDescOnce sync.Once
	file_genesis_messages_v1_instance_event_proto_rawDescData = file_genesis_messages_v1_instance_event_proto_rawDesc
)

func file_genesis_messages_v1_instance_event_proto_rawDescGZIP() []byte {
	file_genesis_messages_v1_instance_event_proto_rawDescOnce.Do(func() {
		file_genesis_messages_v1_instance_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_genesis_messages_v1_instance_event_proto_rawDescData)
	})
	return file_genesis_messages_v1_instance_event_proto_rawDescData
}

var file_genesis_messages_v1_instance_event_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_genesis_messages_v1_instance_event_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_genesis_messages_v1_instance_event_proto_goTypes = []interface{}{
	(InstanceCreationFailedEvent_Reason)(0), // 0: genesis.messages.v1.InstanceCreationFailedEvent.Reason
	(InstanceDeletedEvent_Reason)(0),        // 1: genesis.messages.v1.InstanceDeletedEvent.Reason
	(*InstanceEvent)(nil),                   // 2: genesis.messages.v1.InstanceEvent
	(*InstanceCreatedEvent)(nil),            // 3: genesis.messages.v1.InstanceCreatedEvent
	(*InstanceCreationFailedEvent)(nil),     // 4: genesis.messages.v1.InstanceCreationFailedEvent
	(*InstanceDeletedEvent)(nil),            // 5: genesis.messages.v1.InstanceDeletedEvent
	(*v1.Instance)(nil),                     // 6: genesis.v1.Instance
	(*timestamppb.Timestamp)(nil),           // 7: google.protobuf.Timestamp
	(*v1.InstanceConfig)(nil),               // 8: genesis.v1.InstanceConfig
	(*v1.InstanceResources)(nil),            // 9: genesis.v1.InstanceResources
}
var file_genesis_messages_v1_instance_event_proto_depIdxs = []int32{
	6, // 0: genesis.messages.v1.InstanceEvent.instance:type_name -> genesis.v1.Instance
	7, // 1: genesis.messages.v1.InstanceEvent.timestamp:type_name -> google.protobuf.Timestamp
	3, // 2: genesis.messages.v1.InstanceEvent.created:type_name -> genesis.messages.v1.InstanceCreatedEvent
	4, // 3: genesis.messages.v1.InstanceEvent.creation_failed:type_name -> genesis.messages.v1.InstanceCreationFailedEvent
	5, // 4: genesis.messages.v1.InstanceEvent.deleted:type_name -> genesis.messages.v1.InstanceDeletedEvent
	8, // 5: genesis.messages.v1.InstanceCreatedEvent.config:type_name -> genesis.v1.InstanceConfig
	9, // 6: genesis.messages.v1.InstanceCreatedEvent.resources:type_name -> genesis.v1.InstanceResources
	0, // 7: genesis.messages.v1.InstanceCreationFailedEvent.reason:type_name -> genesis.messages.v1.InstanceCreationFailedEvent.Reason
	1, // 8: genesis.messages.v1.InstanceDeletedEvent.reason:type_name -> genesis.messages.v1.InstanceDeletedEvent.Reason
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_genesis_messages_v1_instance_event_proto_init() }
func file_genesis_messages_v1_instance_event_proto_init() {
	if File_genesis_messages_v1_instance_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_genesis_messages_v1_instance_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceEvent); i {
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
		file_genesis_messages_v1_instance_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceCreatedEvent); i {
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
		file_genesis_messages_v1_instance_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceCreationFailedEvent); i {
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
		file_genesis_messages_v1_instance_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InstanceDeletedEvent); i {
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
	file_genesis_messages_v1_instance_event_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*InstanceEvent_Created)(nil),
		(*InstanceEvent_CreationFailed)(nil),
		(*InstanceEvent_Deleted)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_genesis_messages_v1_instance_event_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_genesis_messages_v1_instance_event_proto_goTypes,
		DependencyIndexes: file_genesis_messages_v1_instance_event_proto_depIdxs,
		EnumInfos:         file_genesis_messages_v1_instance_event_proto_enumTypes,
		MessageInfos:      file_genesis_messages_v1_instance_event_proto_msgTypes,
	}.Build()
	File_genesis_messages_v1_instance_event_proto = out.File
	file_genesis_messages_v1_instance_event_proto_rawDesc = nil
	file_genesis_messages_v1_instance_event_proto_goTypes = nil
	file_genesis_messages_v1_instance_event_proto_depIdxs = nil
}
