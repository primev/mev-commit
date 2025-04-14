// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.19.1
// source: streamer.proto

package pb

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
	PayloadStreamer_Subscribe_FullMethodName = "/pb.PayloadStreamer/Subscribe"
)

// PayloadStreamerClient is the client API for PayloadStreamer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PayloadStreamerClient interface {
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (PayloadStreamer_SubscribeClient, error)
}

type payloadStreamerClient struct {
	cc grpc.ClientConnInterface
}

func NewPayloadStreamerClient(cc grpc.ClientConnInterface) PayloadStreamerClient {
	return &payloadStreamerClient{cc}
}

func (c *payloadStreamerClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (PayloadStreamer_SubscribeClient, error) {
	stream, err := c.cc.NewStream(ctx, &PayloadStreamer_ServiceDesc.Streams[0], PayloadStreamer_Subscribe_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &payloadStreamerSubscribeClient{stream}
	return x, nil
}

type PayloadStreamer_SubscribeClient interface {
	Send(*ClientMessage) error
	Recv() (*PayloadMessage, error)
	grpc.ClientStream
}

type payloadStreamerSubscribeClient struct {
	grpc.ClientStream
}

func (x *payloadStreamerSubscribeClient) Send(m *ClientMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *payloadStreamerSubscribeClient) Recv() (*PayloadMessage, error) {
	m := new(PayloadMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PayloadStreamerServer is the server API for PayloadStreamer service.
// All implementations must embed UnimplementedPayloadStreamerServer
// for forward compatibility
type PayloadStreamerServer interface {
	Subscribe(PayloadStreamer_SubscribeServer) error
	mustEmbedUnimplementedPayloadStreamerServer()
}

// UnimplementedPayloadStreamerServer must be embedded to have forward compatible implementations.
type UnimplementedPayloadStreamerServer struct {
}

func (UnimplementedPayloadStreamerServer) Subscribe(PayloadStreamer_SubscribeServer) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
}
func (UnimplementedPayloadStreamerServer) mustEmbedUnimplementedPayloadStreamerServer() {}

// UnsafePayloadStreamerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PayloadStreamerServer will
// result in compilation errors.
type UnsafePayloadStreamerServer interface {
	mustEmbedUnimplementedPayloadStreamerServer()
}

func RegisterPayloadStreamerServer(s grpc.ServiceRegistrar, srv PayloadStreamerServer) {
	s.RegisterService(&PayloadStreamer_ServiceDesc, srv)
}

func _PayloadStreamer_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PayloadStreamerServer).Subscribe(&payloadStreamerSubscribeServer{stream})
}

type PayloadStreamer_SubscribeServer interface {
	Send(*PayloadMessage) error
	Recv() (*ClientMessage, error)
	grpc.ServerStream
}

type payloadStreamerSubscribeServer struct {
	grpc.ServerStream
}

func (x *payloadStreamerSubscribeServer) Send(m *PayloadMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *payloadStreamerSubscribeServer) Recv() (*ClientMessage, error) {
	m := new(ClientMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// PayloadStreamer_ServiceDesc is the grpc.ServiceDesc for PayloadStreamer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PayloadStreamer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.PayloadStreamer",
	HandlerType: (*PayloadStreamerServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _PayloadStreamer_Subscribe_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "streamer.proto",
}
