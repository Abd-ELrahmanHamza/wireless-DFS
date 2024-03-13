// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.3
// source: client.proto

package mp4_service

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

// MP4ServiceClient is the client API for MP4Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MP4ServiceClient interface {
	Upload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadResponse, error)
	UploadingCompletion(ctx context.Context, in *UploadingCompletionRequest, opts ...grpc.CallOption) (*UploadingCompletionResponse, error)
	Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error)
}

type mP4ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMP4ServiceClient(cc grpc.ClientConnInterface) MP4ServiceClient {
	return &mP4ServiceClient{cc}
}

func (c *mP4ServiceClient) Upload(ctx context.Context, in *UploadRequest, opts ...grpc.CallOption) (*UploadResponse, error) {
	out := new(UploadResponse)
	err := c.cc.Invoke(ctx, "/MP4Service/Upload", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mP4ServiceClient) UploadingCompletion(ctx context.Context, in *UploadingCompletionRequest, opts ...grpc.CallOption) (*UploadingCompletionResponse, error) {
	out := new(UploadingCompletionResponse)
	err := c.cc.Invoke(ctx, "/MP4Service/UploadingCompletion", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *mP4ServiceClient) Download(ctx context.Context, in *DownloadRequest, opts ...grpc.CallOption) (*DownloadResponse, error) {
	out := new(DownloadResponse)
	err := c.cc.Invoke(ctx, "/MP4Service/Download", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MP4ServiceServer is the server API for MP4Service service.
// All implementations must embed UnimplementedMP4ServiceServer
// for forward compatibility
type MP4ServiceServer interface {
	Upload(context.Context, *UploadRequest) (*UploadResponse, error)
	UploadingCompletion(context.Context, *UploadingCompletionRequest) (*UploadingCompletionResponse, error)
	Download(context.Context, *DownloadRequest) (*DownloadResponse, error)
	mustEmbedUnimplementedMP4ServiceServer()
}

// UnimplementedMP4ServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMP4ServiceServer struct {
}

func (UnimplementedMP4ServiceServer) Upload(context.Context, *UploadRequest) (*UploadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Upload not implemented")
}
func (UnimplementedMP4ServiceServer) UploadingCompletion(context.Context, *UploadingCompletionRequest) (*UploadingCompletionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadingCompletion not implemented")
}
func (UnimplementedMP4ServiceServer) Download(context.Context, *DownloadRequest) (*DownloadResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Download not implemented")
}
func (UnimplementedMP4ServiceServer) mustEmbedUnimplementedMP4ServiceServer() {}

// UnsafeMP4ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MP4ServiceServer will
// result in compilation errors.
type UnsafeMP4ServiceServer interface {
	mustEmbedUnimplementedMP4ServiceServer()
}

func RegisterMP4ServiceServer(s grpc.ServiceRegistrar, srv MP4ServiceServer) {
	s.RegisterService(&MP4Service_ServiceDesc, srv)
}

func _MP4Service_Upload_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MP4ServiceServer).Upload(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MP4Service/Upload",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MP4ServiceServer).Upload(ctx, req.(*UploadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MP4Service_UploadingCompletion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadingCompletionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MP4ServiceServer).UploadingCompletion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MP4Service/UploadingCompletion",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MP4ServiceServer).UploadingCompletion(ctx, req.(*UploadingCompletionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MP4Service_Download_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DownloadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MP4ServiceServer).Download(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/MP4Service/Download",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MP4ServiceServer).Download(ctx, req.(*DownloadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MP4Service_ServiceDesc is the grpc.ServiceDesc for MP4Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MP4Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "MP4Service",
	HandlerType: (*MP4ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Upload",
			Handler:    _MP4Service_Upload_Handler,
		},
		{
			MethodName: "UploadingCompletion",
			Handler:    _MP4Service_UploadingCompletion_Handler,
		},
		{
			MethodName: "Download",
			Handler:    _MP4Service_Download_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "client.proto",
}
