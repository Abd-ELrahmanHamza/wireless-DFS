// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: dk.proto

package data_keeper

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
	DataKeeperService_ReplicateFile_FullMethodName = "/dfs.DataKeeperService/ReplicateFile"
)

// DataKeeperServiceClient is the client API for DataKeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataKeeperServiceClient interface {
	ReplicateFile(ctx context.Context, in *ReplicateRequest, opts ...grpc.CallOption) (*ReplicateResponse, error)
}

type dataKeeperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataKeeperServiceClient(cc grpc.ClientConnInterface) DataKeeperServiceClient {
	return &dataKeeperServiceClient{cc}
}

func (c *dataKeeperServiceClient) ReplicateFile(ctx context.Context, in *ReplicateRequest, opts ...grpc.CallOption) (*ReplicateResponse, error) {
	out := new(ReplicateResponse)
	err := c.cc.Invoke(ctx, DataKeeperService_ReplicateFile_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataKeeperServiceServer is the server API for DataKeeperService service.
// All implementations must embed UnimplementedDataKeeperServiceServer
// for forward compatibility
type DataKeeperServiceServer interface {
	ReplicateFile(context.Context, *ReplicateRequest) (*ReplicateResponse, error)
	mustEmbedUnimplementedDataKeeperServiceServer()
}

// UnimplementedDataKeeperServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDataKeeperServiceServer struct {
}

func (UnimplementedDataKeeperServiceServer) ReplicateFile(context.Context, *ReplicateRequest) (*ReplicateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicateFile not implemented")
}
func (UnimplementedDataKeeperServiceServer) mustEmbedUnimplementedDataKeeperServiceServer() {}

// UnsafeDataKeeperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataKeeperServiceServer will
// result in compilation errors.
type UnsafeDataKeeperServiceServer interface {
	mustEmbedUnimplementedDataKeeperServiceServer()
}

func RegisterDataKeeperServiceServer(s grpc.ServiceRegistrar, srv DataKeeperServiceServer) {
	s.RegisterService(&DataKeeperService_ServiceDesc, srv)
}

func _DataKeeperService_ReplicateFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServiceServer).ReplicateFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeperService_ReplicateFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServiceServer).ReplicateFile(ctx, req.(*ReplicateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataKeeperService_ServiceDesc is the grpc.ServiceDesc for DataKeeperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataKeeperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "dfs.DataKeeperService",
	HandlerType: (*DataKeeperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicateFile",
			Handler:    _DataKeeperService_ReplicateFile_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "dk.proto",
}