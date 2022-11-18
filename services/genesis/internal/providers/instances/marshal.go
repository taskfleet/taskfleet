package instances

import (
	genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"
	"go.taskfleet.io/services/genesis/internal/typedefs"
)

// GPUResourcesUnmarshalProto converts a message into GPU resources.
func GPUResourcesUnmarshalProto(message *genesis_v1.GPUResources) *GPUResources {
	if message == nil {
		return nil
	}
	return &GPUResources{
		Kind:  typedefs.GPUKindUnmarshalProto(message.Kind),
		Count: uint16(message.Count),
	}
}

// MarshalProto converts the GPU resources into a message.
func (r *GPUResources) MarshalProto() *genesis_v1.GPUResources {
	if r == nil {
		return nil
	}
	return &genesis_v1.GPUResources{
		Kind:  r.Kind.MarshalProto(),
		Count: uint32(r.Count),
	}
}
