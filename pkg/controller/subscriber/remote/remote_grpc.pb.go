// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package remote

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

// SubscriberClient is the client API for Subscriber service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SubscriberClient interface {
	AddedStream(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Subscriber_AddedStreamClient, error)
	CreatedStream(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Subscriber_CreatedStreamClient, error)
}

type subscriberClient struct {
	cc grpc.ClientConnInterface
}

func NewSubscriberClient(cc grpc.ClientConnInterface) SubscriberClient {
	return &subscriberClient{cc}
}

func (c *subscriberClient) AddedStream(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Subscriber_AddedStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Subscriber_ServiceDesc.Streams[0], "/Subscriber/AddedStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &subscriberAddedStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Subscriber_AddedStreamClient interface {
	Recv() (*AddedResponse, error)
	grpc.ClientStream
}

type subscriberAddedStreamClient struct {
	grpc.ClientStream
}

func (x *subscriberAddedStreamClient) Recv() (*AddedResponse, error) {
	m := new(AddedResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *subscriberClient) CreatedStream(ctx context.Context, in *ConsumeRequest, opts ...grpc.CallOption) (Subscriber_CreatedStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &Subscriber_ServiceDesc.Streams[1], "/Subscriber/CreatedStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &subscriberCreatedStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Subscriber_CreatedStreamClient interface {
	Recv() (*CreatedResponse, error)
	grpc.ClientStream
}

type subscriberCreatedStreamClient struct {
	grpc.ClientStream
}

func (x *subscriberCreatedStreamClient) Recv() (*CreatedResponse, error) {
	m := new(CreatedResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SubscriberServer is the server API for Subscriber service.
// All implementations must embed UnimplementedSubscriberServer
// for forward compatibility
type SubscriberServer interface {
	AddedStream(*ConsumeRequest, Subscriber_AddedStreamServer) error
	CreatedStream(*ConsumeRequest, Subscriber_CreatedStreamServer) error
	mustEmbedUnimplementedSubscriberServer()
}

// UnimplementedSubscriberServer must be embedded to have forward compatible implementations.
type UnimplementedSubscriberServer struct {
}

func (UnimplementedSubscriberServer) AddedStream(*ConsumeRequest, Subscriber_AddedStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method AddedStream not implemented")
}
func (UnimplementedSubscriberServer) CreatedStream(*ConsumeRequest, Subscriber_CreatedStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method CreatedStream not implemented")
}
func (UnimplementedSubscriberServer) mustEmbedUnimplementedSubscriberServer() {}

// UnsafeSubscriberServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SubscriberServer will
// result in compilation errors.
type UnsafeSubscriberServer interface {
	mustEmbedUnimplementedSubscriberServer()
}

func RegisterSubscriberServer(s grpc.ServiceRegistrar, srv SubscriberServer) {
	s.RegisterService(&Subscriber_ServiceDesc, srv)
}

func _Subscriber_AddedStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SubscriberServer).AddedStream(m, &subscriberAddedStreamServer{stream})
}

type Subscriber_AddedStreamServer interface {
	Send(*AddedResponse) error
	grpc.ServerStream
}

type subscriberAddedStreamServer struct {
	grpc.ServerStream
}

func (x *subscriberAddedStreamServer) Send(m *AddedResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _Subscriber_CreatedStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ConsumeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(SubscriberServer).CreatedStream(m, &subscriberCreatedStreamServer{stream})
}

type Subscriber_CreatedStreamServer interface {
	Send(*CreatedResponse) error
	grpc.ServerStream
}

type subscriberCreatedStreamServer struct {
	grpc.ServerStream
}

func (x *subscriberCreatedStreamServer) Send(m *CreatedResponse) error {
	return x.ServerStream.SendMsg(m)
}

// Subscriber_ServiceDesc is the grpc.ServiceDesc for Subscriber service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Subscriber_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Subscriber",
	HandlerType: (*SubscriberServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "AddedStream",
			Handler:       _Subscriber_AddedStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "CreatedStream",
			Handler:       _Subscriber_CreatedStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "remote.proto",
}