// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: genesis/v1/service.proto

package genesis

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	GenesisService_ListZones_FullMethodName        = "/genesis.v1.GenesisService/ListZones"
	GenesisService_CreateInstance_FullMethodName   = "/genesis.v1.GenesisService/CreateInstance"
	GenesisService_ListInstances_FullMethodName    = "/genesis.v1.GenesisService/ListInstances"
	GenesisService_ShutdownInstance_FullMethodName = "/genesis.v1.GenesisService/ShutdownInstance"
)

// GenesisServiceClient is the client API for GenesisService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GenesisServiceClient interface {
	// ListZones returns a list of all providers' available zones. For each zone, it additionally
	// returns a list of available GPUs. This call should be used to send correct zone identifiers
	// to the `CreateInstance call.`
	ListZones(ctx context.Context, in *ListZonesRequest, opts ...grpc.CallOption) (*ListZonesResponse, error)
	// CreateInstance creates a new instance according to the given requirements. If the creation
	// is successful, the instance configuration will be returned. The full instance status
	// (including its IP) will be delivered via Kafka as soon as the instance is up and running.
	CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (*CreateInstanceResponse, error)
	// ListInstances returns all the instances that are owned by a particular owner and which are
	// running at the moment. In particular, the returned set of instances does not include
	// instances which were requested successfully but are not running yet.
	ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error)
	// ShutdownInstance shuts down the instance described by the request. It does not return
	// anything if deletion was successful.
	ShutdownInstance(ctx context.Context, in *ShutdownInstanceRequest, opts ...grpc.CallOption) (*ShutdownInstanceResponse, error)
}

type genesisServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGenesisServiceClient(cc grpc.ClientConnInterface) GenesisServiceClient {
	return &genesisServiceClient{cc}
}

func (c *genesisServiceClient) ListZones(ctx context.Context, in *ListZonesRequest, opts ...grpc.CallOption) (*ListZonesResponse, error) {
	out := new(ListZonesResponse)
	err := c.cc.Invoke(ctx, GenesisService_ListZones_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genesisServiceClient) CreateInstance(ctx context.Context, in *CreateInstanceRequest, opts ...grpc.CallOption) (*CreateInstanceResponse, error) {
	out := new(CreateInstanceResponse)
	err := c.cc.Invoke(ctx, GenesisService_CreateInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genesisServiceClient) ListInstances(ctx context.Context, in *ListInstancesRequest, opts ...grpc.CallOption) (*ListInstancesResponse, error) {
	out := new(ListInstancesResponse)
	err := c.cc.Invoke(ctx, GenesisService_ListInstances_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *genesisServiceClient) ShutdownInstance(ctx context.Context, in *ShutdownInstanceRequest, opts ...grpc.CallOption) (*ShutdownInstanceResponse, error) {
	out := new(ShutdownInstanceResponse)
	err := c.cc.Invoke(ctx, GenesisService_ShutdownInstance_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GenesisServiceServer is the server API for GenesisService service.
// All implementations must embed UnimplementedGenesisServiceServer
// for forward compatibility
type GenesisServiceServer interface {
	// ListZones returns a list of all providers' available zones. For each zone, it additionally
	// returns a list of available GPUs. This call should be used to send correct zone identifiers
	// to the `CreateInstance call.`
	ListZones(context.Context, *ListZonesRequest) (*ListZonesResponse, error)
	// CreateInstance creates a new instance according to the given requirements. If the creation
	// is successful, the instance configuration will be returned. The full instance status
	// (including its IP) will be delivered via Kafka as soon as the instance is up and running.
	CreateInstance(context.Context, *CreateInstanceRequest) (*CreateInstanceResponse, error)
	// ListInstances returns all the instances that are owned by a particular owner and which are
	// running at the moment. In particular, the returned set of instances does not include
	// instances which were requested successfully but are not running yet.
	ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error)
	// ShutdownInstance shuts down the instance described by the request. It does not return
	// anything if deletion was successful.
	ShutdownInstance(context.Context, *ShutdownInstanceRequest) (*ShutdownInstanceResponse, error)
	mustEmbedUnimplementedGenesisServiceServer()
}

// UnimplementedGenesisServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGenesisServiceServer struct {
}

func (UnimplementedGenesisServiceServer) ListZones(context.Context, *ListZonesRequest) (*ListZonesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListZones not implemented")
}
func (UnimplementedGenesisServiceServer) CreateInstance(context.Context, *CreateInstanceRequest) (*CreateInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateInstance not implemented")
}
func (UnimplementedGenesisServiceServer) ListInstances(context.Context, *ListInstancesRequest) (*ListInstancesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListInstances not implemented")
}
func (UnimplementedGenesisServiceServer) ShutdownInstance(context.Context, *ShutdownInstanceRequest) (*ShutdownInstanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShutdownInstance not implemented")
}
func (UnimplementedGenesisServiceServer) mustEmbedUnimplementedGenesisServiceServer() {}

// UnsafeGenesisServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GenesisServiceServer will
// result in compilation errors.
type UnsafeGenesisServiceServer interface {
	mustEmbedUnimplementedGenesisServiceServer()
}

func RegisterGenesisServiceServer(s grpc.ServiceRegistrar, srv GenesisServiceServer) {
	s.RegisterService(&GenesisService_ServiceDesc, srv)
}

func _GenesisService_ListZones_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListZonesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenesisServiceServer).ListZones(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GenesisService_ListZones_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenesisServiceServer).ListZones(ctx, req.(*ListZonesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenesisService_CreateInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenesisServiceServer).CreateInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GenesisService_CreateInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenesisServiceServer).CreateInstance(ctx, req.(*CreateInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenesisService_ListInstances_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListInstancesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenesisServiceServer).ListInstances(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GenesisService_ListInstances_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenesisServiceServer).ListInstances(ctx, req.(*ListInstancesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GenesisService_ShutdownInstance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownInstanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GenesisServiceServer).ShutdownInstance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GenesisService_ShutdownInstance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GenesisServiceServer).ShutdownInstance(ctx, req.(*ShutdownInstanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GenesisService_ServiceDesc is the grpc.ServiceDesc for GenesisService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GenesisService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "genesis.v1.GenesisService",
	HandlerType: (*GenesisServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListZones",
			Handler:    _GenesisService_ListZones_Handler,
		},
		{
			MethodName: "CreateInstance",
			Handler:    _GenesisService_CreateInstance_Handler,
		},
		{
			MethodName: "ListInstances",
			Handler:    _GenesisService_ListInstances_Handler,
		},
		{
			MethodName: "ShutdownInstance",
			Handler:    _GenesisService_ShutdownInstance_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "genesis/v1/service.proto",
}
