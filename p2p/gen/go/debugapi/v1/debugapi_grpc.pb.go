// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: debugapi/v1/debugapi.proto

package debugapiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DebugService_GetTopology_FullMethodName            = "/debugapi.v1.DebugService/GetTopology"
	DebugService_GetPendingTransactions_FullMethodName = "/debugapi.v1.DebugService/GetPendingTransactions"
	DebugService_CancelTransaction_FullMethodName      = "/debugapi.v1.DebugService/CancelTransaction"
)

// DebugServiceClient is the client API for DebugService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DebugServiceClient interface {
	// GetTopology
	//
	// GetTopology is called by the user to get the topology of the node. The topology
	// includes connectivity information about the node.
	GetTopology(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*TopologyResponse, error)
	// GetPendingTransactions
	//
	// GetPendingTransactions is called by the provider to get the pending transactions for the wallet.
	GetPendingTransactions(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*PendingTransactionsResponse, error)
	// CancelTransaction
	//
	// CancelTransaction is called by the provider to cancel a transaction sent from this wallet.
	CancelTransaction(ctx context.Context, in *CancelTransactionReq, opts ...grpc.CallOption) (*CancelTransactionResponse, error)
}

type debugServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDebugServiceClient(cc grpc.ClientConnInterface) DebugServiceClient {
	return &debugServiceClient{cc}
}

func (c *debugServiceClient) GetTopology(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*TopologyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TopologyResponse)
	err := c.cc.Invoke(ctx, DebugService_GetTopology_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debugServiceClient) GetPendingTransactions(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*PendingTransactionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PendingTransactionsResponse)
	err := c.cc.Invoke(ctx, DebugService_GetPendingTransactions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debugServiceClient) CancelTransaction(ctx context.Context, in *CancelTransactionReq, opts ...grpc.CallOption) (*CancelTransactionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CancelTransactionResponse)
	err := c.cc.Invoke(ctx, DebugService_CancelTransaction_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DebugServiceServer is the server API for DebugService service.
// All implementations must embed UnimplementedDebugServiceServer
// for forward compatibility.
type DebugServiceServer interface {
	// GetTopology
	//
	// GetTopology is called by the user to get the topology of the node. The topology
	// includes connectivity information about the node.
	GetTopology(context.Context, *EmptyMessage) (*TopologyResponse, error)
	// GetPendingTransactions
	//
	// GetPendingTransactions is called by the provider to get the pending transactions for the wallet.
	GetPendingTransactions(context.Context, *EmptyMessage) (*PendingTransactionsResponse, error)
	// CancelTransaction
	//
	// CancelTransaction is called by the provider to cancel a transaction sent from this wallet.
	CancelTransaction(context.Context, *CancelTransactionReq) (*CancelTransactionResponse, error)
	mustEmbedUnimplementedDebugServiceServer()
}

// UnimplementedDebugServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDebugServiceServer struct{}

func (UnimplementedDebugServiceServer) GetTopology(context.Context, *EmptyMessage) (*TopologyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopology not implemented")
}
func (UnimplementedDebugServiceServer) GetPendingTransactions(context.Context, *EmptyMessage) (*PendingTransactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPendingTransactions not implemented")
}
func (UnimplementedDebugServiceServer) CancelTransaction(context.Context, *CancelTransactionReq) (*CancelTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelTransaction not implemented")
}
func (UnimplementedDebugServiceServer) mustEmbedUnimplementedDebugServiceServer() {}
func (UnimplementedDebugServiceServer) testEmbeddedByValue()                      {}

// UnsafeDebugServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DebugServiceServer will
// result in compilation errors.
type UnsafeDebugServiceServer interface {
	mustEmbedUnimplementedDebugServiceServer()
}

func RegisterDebugServiceServer(s grpc.ServiceRegistrar, srv DebugServiceServer) {
	// If the following call pancis, it indicates UnimplementedDebugServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DebugService_ServiceDesc, srv)
}

func _DebugService_GetTopology_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugServiceServer).GetTopology(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebugService_GetTopology_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugServiceServer).GetTopology(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebugService_GetPendingTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugServiceServer).GetPendingTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebugService_GetPendingTransactions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugServiceServer).GetPendingTransactions(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebugService_CancelTransaction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CancelTransactionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugServiceServer).CancelTransaction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DebugService_CancelTransaction_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugServiceServer).CancelTransaction(ctx, req.(*CancelTransactionReq))
	}
	return interceptor(ctx, in, info, handler)
}

// DebugService_ServiceDesc is the grpc.ServiceDesc for DebugService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DebugService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "debugapi.v1.DebugService",
	HandlerType: (*DebugServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopology",
			Handler:    _DebugService_GetTopology_Handler,
		},
		{
			MethodName: "GetPendingTransactions",
			Handler:    _DebugService_GetPendingTransactions_Handler,
		},
		{
			MethodName: "CancelTransaction",
			Handler:    _DebugService_CancelTransaction_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "debugapi/v1/debugapi.proto",
}
