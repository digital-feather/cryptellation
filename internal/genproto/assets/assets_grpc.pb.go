// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: assets.proto

package assets

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

// AssetsServiceClient is the client API for AssetsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AssetsServiceClient interface {
	CreateAssets(ctx context.Context, in *CreateAssetsRequest, opts ...grpc.CallOption) (*CreateAssetsResponse, error)
	ReadAssets(ctx context.Context, in *ReadAssetsRequest, opts ...grpc.CallOption) (*ReadAssetsResponse, error)
}

type assetsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAssetsServiceClient(cc grpc.ClientConnInterface) AssetsServiceClient {
	return &assetsServiceClient{cc}
}

func (c *assetsServiceClient) CreateAssets(ctx context.Context, in *CreateAssetsRequest, opts ...grpc.CallOption) (*CreateAssetsResponse, error) {
	out := new(CreateAssetsResponse)
	err := c.cc.Invoke(ctx, "/assets.AssetsService/CreateAssets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *assetsServiceClient) ReadAssets(ctx context.Context, in *ReadAssetsRequest, opts ...grpc.CallOption) (*ReadAssetsResponse, error) {
	out := new(ReadAssetsResponse)
	err := c.cc.Invoke(ctx, "/assets.AssetsService/ReadAssets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AssetsServiceServer is the server API for AssetsService service.
// All implementations should embed UnimplementedAssetsServiceServer
// for forward compatibility
type AssetsServiceServer interface {
	CreateAssets(context.Context, *CreateAssetsRequest) (*CreateAssetsResponse, error)
	ReadAssets(context.Context, *ReadAssetsRequest) (*ReadAssetsResponse, error)
}

// UnimplementedAssetsServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAssetsServiceServer struct {
}

func (UnimplementedAssetsServiceServer) CreateAssets(context.Context, *CreateAssetsRequest) (*CreateAssetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAssets not implemented")
}
func (UnimplementedAssetsServiceServer) ReadAssets(context.Context, *ReadAssetsRequest) (*ReadAssetsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadAssets not implemented")
}

// UnsafeAssetsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AssetsServiceServer will
// result in compilation errors.
type UnsafeAssetsServiceServer interface {
	mustEmbedUnimplementedAssetsServiceServer()
}

func RegisterAssetsServiceServer(s grpc.ServiceRegistrar, srv AssetsServiceServer) {
	s.RegisterService(&AssetsService_ServiceDesc, srv)
}

func _AssetsService_CreateAssets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAssetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetsServiceServer).CreateAssets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/assets.AssetsService/CreateAssets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetsServiceServer).CreateAssets(ctx, req.(*CreateAssetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AssetsService_ReadAssets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadAssetsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AssetsServiceServer).ReadAssets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/assets.AssetsService/ReadAssets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AssetsServiceServer).ReadAssets(ctx, req.(*ReadAssetsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AssetsService_ServiceDesc is the grpc.ServiceDesc for AssetsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AssetsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "assets.AssetsService",
	HandlerType: (*AssetsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateAssets",
			Handler:    _AssetsService_CreateAssets_Handler,
		},
		{
			MethodName: "ReadAssets",
			Handler:    _AssetsService_ReadAssets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "assets.proto",
}
