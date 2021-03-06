// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: backtests.proto

package backtests

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

// BacktestsServiceClient is the client API for BacktestsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BacktestsServiceClient interface {
	CreateBacktest(ctx context.Context, in *CreateBacktestRequest, opts ...grpc.CallOption) (*CreateBacktestResponse, error)
	SubscribeToBacktestEvents(ctx context.Context, in *SubscribeToBacktestEventsRequest, opts ...grpc.CallOption) (*SubscribeToBacktestEventsResponse, error)
	ListenBacktest(ctx context.Context, opts ...grpc.CallOption) (BacktestsService_ListenBacktestClient, error)
	CreateBacktestOrder(ctx context.Context, in *CreateBacktestOrderRequest, opts ...grpc.CallOption) (*CreateBacktestOrderResponse, error)
	Accounts(ctx context.Context, in *AccountsRequest, opts ...grpc.CallOption) (*AccountsResponse, error)
	Orders(ctx context.Context, in *OrdersRequest, opts ...grpc.CallOption) (*OrdersResponse, error)
}

type backtestsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBacktestsServiceClient(cc grpc.ClientConnInterface) BacktestsServiceClient {
	return &backtestsServiceClient{cc}
}

func (c *backtestsServiceClient) CreateBacktest(ctx context.Context, in *CreateBacktestRequest, opts ...grpc.CallOption) (*CreateBacktestResponse, error) {
	out := new(CreateBacktestResponse)
	err := c.cc.Invoke(ctx, "/backtests.BacktestsService/CreateBacktest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtestsServiceClient) SubscribeToBacktestEvents(ctx context.Context, in *SubscribeToBacktestEventsRequest, opts ...grpc.CallOption) (*SubscribeToBacktestEventsResponse, error) {
	out := new(SubscribeToBacktestEventsResponse)
	err := c.cc.Invoke(ctx, "/backtests.BacktestsService/SubscribeToBacktestEvents", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtestsServiceClient) ListenBacktest(ctx context.Context, opts ...grpc.CallOption) (BacktestsService_ListenBacktestClient, error) {
	stream, err := c.cc.NewStream(ctx, &BacktestsService_ServiceDesc.Streams[0], "/backtests.BacktestsService/ListenBacktest", opts...)
	if err != nil {
		return nil, err
	}
	x := &backtestsServiceListenBacktestClient{stream}
	return x, nil
}

type BacktestsService_ListenBacktestClient interface {
	Send(*BacktestEventRequest) error
	Recv() (*BacktestEventResponse, error)
	grpc.ClientStream
}

type backtestsServiceListenBacktestClient struct {
	grpc.ClientStream
}

func (x *backtestsServiceListenBacktestClient) Send(m *BacktestEventRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *backtestsServiceListenBacktestClient) Recv() (*BacktestEventResponse, error) {
	m := new(BacktestEventResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backtestsServiceClient) CreateBacktestOrder(ctx context.Context, in *CreateBacktestOrderRequest, opts ...grpc.CallOption) (*CreateBacktestOrderResponse, error) {
	out := new(CreateBacktestOrderResponse)
	err := c.cc.Invoke(ctx, "/backtests.BacktestsService/CreateBacktestOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtestsServiceClient) Accounts(ctx context.Context, in *AccountsRequest, opts ...grpc.CallOption) (*AccountsResponse, error) {
	out := new(AccountsResponse)
	err := c.cc.Invoke(ctx, "/backtests.BacktestsService/Accounts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backtestsServiceClient) Orders(ctx context.Context, in *OrdersRequest, opts ...grpc.CallOption) (*OrdersResponse, error) {
	out := new(OrdersResponse)
	err := c.cc.Invoke(ctx, "/backtests.BacktestsService/Orders", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BacktestsServiceServer is the server API for BacktestsService service.
// All implementations should embed UnimplementedBacktestsServiceServer
// for forward compatibility
type BacktestsServiceServer interface {
	CreateBacktest(context.Context, *CreateBacktestRequest) (*CreateBacktestResponse, error)
	SubscribeToBacktestEvents(context.Context, *SubscribeToBacktestEventsRequest) (*SubscribeToBacktestEventsResponse, error)
	ListenBacktest(BacktestsService_ListenBacktestServer) error
	CreateBacktestOrder(context.Context, *CreateBacktestOrderRequest) (*CreateBacktestOrderResponse, error)
	Accounts(context.Context, *AccountsRequest) (*AccountsResponse, error)
	Orders(context.Context, *OrdersRequest) (*OrdersResponse, error)
}

// UnimplementedBacktestsServiceServer should be embedded to have forward compatible implementations.
type UnimplementedBacktestsServiceServer struct {
}

func (UnimplementedBacktestsServiceServer) CreateBacktest(context.Context, *CreateBacktestRequest) (*CreateBacktestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBacktest not implemented")
}
func (UnimplementedBacktestsServiceServer) SubscribeToBacktestEvents(context.Context, *SubscribeToBacktestEventsRequest) (*SubscribeToBacktestEventsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeToBacktestEvents not implemented")
}
func (UnimplementedBacktestsServiceServer) ListenBacktest(BacktestsService_ListenBacktestServer) error {
	return status.Errorf(codes.Unimplemented, "method ListenBacktest not implemented")
}
func (UnimplementedBacktestsServiceServer) CreateBacktestOrder(context.Context, *CreateBacktestOrderRequest) (*CreateBacktestOrderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBacktestOrder not implemented")
}
func (UnimplementedBacktestsServiceServer) Accounts(context.Context, *AccountsRequest) (*AccountsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Accounts not implemented")
}
func (UnimplementedBacktestsServiceServer) Orders(context.Context, *OrdersRequest) (*OrdersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Orders not implemented")
}

// UnsafeBacktestsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BacktestsServiceServer will
// result in compilation errors.
type UnsafeBacktestsServiceServer interface {
	mustEmbedUnimplementedBacktestsServiceServer()
}

func RegisterBacktestsServiceServer(s grpc.ServiceRegistrar, srv BacktestsServiceServer) {
	s.RegisterService(&BacktestsService_ServiceDesc, srv)
}

func _BacktestsService_CreateBacktest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBacktestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestsServiceServer).CreateBacktest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backtests.BacktestsService/CreateBacktest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestsServiceServer).CreateBacktest(ctx, req.(*CreateBacktestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktestsService_SubscribeToBacktestEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeToBacktestEventsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestsServiceServer).SubscribeToBacktestEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backtests.BacktestsService/SubscribeToBacktestEvents",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestsServiceServer).SubscribeToBacktestEvents(ctx, req.(*SubscribeToBacktestEventsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktestsService_ListenBacktest_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(BacktestsServiceServer).ListenBacktest(&backtestsServiceListenBacktestServer{stream})
}

type BacktestsService_ListenBacktestServer interface {
	Send(*BacktestEventResponse) error
	Recv() (*BacktestEventRequest, error)
	grpc.ServerStream
}

type backtestsServiceListenBacktestServer struct {
	grpc.ServerStream
}

func (x *backtestsServiceListenBacktestServer) Send(m *BacktestEventResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *backtestsServiceListenBacktestServer) Recv() (*BacktestEventRequest, error) {
	m := new(BacktestEventRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _BacktestsService_CreateBacktestOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBacktestOrderRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestsServiceServer).CreateBacktestOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backtests.BacktestsService/CreateBacktestOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestsServiceServer).CreateBacktestOrder(ctx, req.(*CreateBacktestOrderRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktestsService_Accounts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AccountsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestsServiceServer).Accounts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backtests.BacktestsService/Accounts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestsServiceServer).Accounts(ctx, req.(*AccountsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _BacktestsService_Orders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OrdersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BacktestsServiceServer).Orders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/backtests.BacktestsService/Orders",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BacktestsServiceServer).Orders(ctx, req.(*OrdersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BacktestsService_ServiceDesc is the grpc.ServiceDesc for BacktestsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BacktestsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "backtests.BacktestsService",
	HandlerType: (*BacktestsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBacktest",
			Handler:    _BacktestsService_CreateBacktest_Handler,
		},
		{
			MethodName: "SubscribeToBacktestEvents",
			Handler:    _BacktestsService_SubscribeToBacktestEvents_Handler,
		},
		{
			MethodName: "CreateBacktestOrder",
			Handler:    _BacktestsService_CreateBacktestOrder_Handler,
		},
		{
			MethodName: "Accounts",
			Handler:    _BacktestsService_Accounts_Handler,
		},
		{
			MethodName: "Orders",
			Handler:    _BacktestsService_Orders_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListenBacktest",
			Handler:       _BacktestsService_ListenBacktest_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "backtests.proto",
}
