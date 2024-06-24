// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: bidderapi/v1/bidderapi.proto

package bidderapiv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Bidder_SendBid_FullMethodName           = "/bidderapi.v1.Bidder/SendBid"
	Bidder_Deposit_FullMethodName           = "/bidderapi.v1.Bidder/Deposit"
	Bidder_AutoDeposit_FullMethodName       = "/bidderapi.v1.Bidder/AutoDeposit"
	Bidder_CancelAutoDeposit_FullMethodName = "/bidderapi.v1.Bidder/CancelAutoDeposit"
	Bidder_GetDeposit_FullMethodName        = "/bidderapi.v1.Bidder/GetDeposit"
	Bidder_GetMinDeposit_FullMethodName     = "/bidderapi.v1.Bidder/GetMinDeposit"
	Bidder_Withdraw_FullMethodName          = "/bidderapi.v1.Bidder/Withdraw"
)

// BidderClient is the client API for Bidder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BidderClient interface {
	// SendBid
	//
	// Send a bid to the bidder mev-commit node.
	SendBid(ctx context.Context, in *Bid, opts ...grpc.CallOption) (Bidder_SendBidClient, error)
	// Deposit
	//
	// Deposit is called by the bidder node to add deposit in the bidder registry.
	Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error)
	// AutoDeposit
	//
	// AutoDeposit is called by the bidder node to add deposit in the bidder registry and move funds automatily from one window to another.
	AutoDeposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*AutoDepositResponse, error)
	// CancelAutoDeposit
	//
	// CancelAutoDeposit is called by the bidder node to cancel the auto deposit.
	CancelAutoDeposit(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error)
	// GetDeposit
	//
	// GetDeposit is called by the bidder to get its deposit in the bidder registry.
	GetDeposit(ctx context.Context, in *GetDepositRequest, opts ...grpc.CallOption) (*DepositResponse, error)
	// GetMinDeposit
	//
	// GetMinDeposit is called by the bidder to get the minimum deposit required in the bidder registry to make bids.
	GetMinDeposit(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*DepositResponse, error)
	// Withdraw
	//
	// Withdraw is called by the bidder to withdraw deposit from the bidder registry.
	Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error)
}

type bidderClient struct {
	cc grpc.ClientConnInterface
}

func NewBidderClient(cc grpc.ClientConnInterface) BidderClient {
	return &bidderClient{cc}
}

func (c *bidderClient) SendBid(ctx context.Context, in *Bid, opts ...grpc.CallOption) (Bidder_SendBidClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Bidder_ServiceDesc.Streams[0], Bidder_SendBid_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &bidderSendBidClient{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Bidder_SendBidClient interface {
	Recv() (*Commitment, error)
	grpc.ClientStream
}

type bidderSendBidClient struct {
	grpc.ClientStream
}

func (x *bidderSendBidClient) Recv() (*Commitment, error) {
	m := new(Commitment)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *bidderClient) Deposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*DepositResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DepositResponse)
	err := c.cc.Invoke(ctx, Bidder_Deposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bidderClient) AutoDeposit(ctx context.Context, in *DepositRequest, opts ...grpc.CallOption) (*AutoDepositResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AutoDepositResponse)
	err := c.cc.Invoke(ctx, Bidder_AutoDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bidderClient) CancelAutoDeposit(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*EmptyMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyMessage)
	err := c.cc.Invoke(ctx, Bidder_CancelAutoDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bidderClient) GetDeposit(ctx context.Context, in *GetDepositRequest, opts ...grpc.CallOption) (*DepositResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DepositResponse)
	err := c.cc.Invoke(ctx, Bidder_GetDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bidderClient) GetMinDeposit(ctx context.Context, in *EmptyMessage, opts ...grpc.CallOption) (*DepositResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DepositResponse)
	err := c.cc.Invoke(ctx, Bidder_GetMinDeposit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bidderClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WithdrawResponse)
	err := c.cc.Invoke(ctx, Bidder_Withdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BidderServer is the server API for Bidder service.
// All implementations must embed UnimplementedBidderServer
// for forward compatibility
type BidderServer interface {
	// SendBid
	//
	// Send a bid to the bidder mev-commit node.
	SendBid(*Bid, Bidder_SendBidServer) error
	// Deposit
	//
	// Deposit is called by the bidder node to add deposit in the bidder registry.
	Deposit(context.Context, *DepositRequest) (*DepositResponse, error)
	// AutoDeposit
	//
	// AutoDeposit is called by the bidder node to add deposit in the bidder registry and move funds automatily from one window to another.
	AutoDeposit(context.Context, *DepositRequest) (*AutoDepositResponse, error)
	// CancelAutoDeposit
	//
	// CancelAutoDeposit is called by the bidder node to cancel the auto deposit.
	CancelAutoDeposit(context.Context, *EmptyMessage) (*EmptyMessage, error)
	// GetDeposit
	//
	// GetDeposit is called by the bidder to get its deposit in the bidder registry.
	GetDeposit(context.Context, *GetDepositRequest) (*DepositResponse, error)
	// GetMinDeposit
	//
	// GetMinDeposit is called by the bidder to get the minimum deposit required in the bidder registry to make bids.
	GetMinDeposit(context.Context, *EmptyMessage) (*DepositResponse, error)
	// Withdraw
	//
	// Withdraw is called by the bidder to withdraw deposit from the bidder registry.
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	mustEmbedUnimplementedBidderServer()
}

// UnimplementedBidderServer must be embedded to have forward compatible implementations.
type UnimplementedBidderServer struct {
}

func (UnimplementedBidderServer) SendBid(*Bid, Bidder_SendBidServer) error {
	return status.Errorf(codes.Unimplemented, "method SendBid not implemented")
}
func (UnimplementedBidderServer) Deposit(context.Context, *DepositRequest) (*DepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deposit not implemented")
}
func (UnimplementedBidderServer) AutoDeposit(context.Context, *DepositRequest) (*AutoDepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AutoDeposit not implemented")
}
func (UnimplementedBidderServer) CancelAutoDeposit(context.Context, *EmptyMessage) (*EmptyMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CancelAutoDeposit not implemented")
}
func (UnimplementedBidderServer) GetDeposit(context.Context, *GetDepositRequest) (*DepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDeposit not implemented")
}
func (UnimplementedBidderServer) GetMinDeposit(context.Context, *EmptyMessage) (*DepositResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMinDeposit not implemented")
}
func (UnimplementedBidderServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedBidderServer) mustEmbedUnimplementedBidderServer() {}

// UnsafeBidderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BidderServer will
// result in compilation errors.
type UnsafeBidderServer interface {
	mustEmbedUnimplementedBidderServer()
}

func RegisterBidderServer(s grpc.ServiceRegistrar, srv BidderServer) {
	s.RegisterService(&Bidder_ServiceDesc, srv)
}

func _Bidder_SendBid_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Bid)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BidderServer).SendBid(m, &bidderSendBidServer{ServerStream: stream})
}

type Bidder_SendBidServer interface {
	Send(*Commitment) error
	grpc.ServerStream
}

type bidderSendBidServer struct {
	grpc.ServerStream
}

func (x *bidderSendBidServer) Send(m *Commitment) error {
	return x.ServerStream.SendMsg(m)
}

func _Bidder_Deposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidderServer).Deposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bidder_Deposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidderServer).Deposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bidder_AutoDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidderServer).AutoDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bidder_AutoDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidderServer).AutoDeposit(ctx, req.(*DepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bidder_CancelAutoDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidderServer).CancelAutoDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bidder_CancelAutoDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidderServer).CancelAutoDeposit(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bidder_GetDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDepositRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidderServer).GetDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bidder_GetDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidderServer).GetDeposit(ctx, req.(*GetDepositRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bidder_GetMinDeposit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmptyMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidderServer).GetMinDeposit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bidder_GetMinDeposit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidderServer).GetMinDeposit(ctx, req.(*EmptyMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bidder_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BidderServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Bidder_Withdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BidderServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Bidder_ServiceDesc is the grpc.ServiceDesc for Bidder service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bidder_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bidderapi.v1.Bidder",
	HandlerType: (*BidderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Deposit",
			Handler:    _Bidder_Deposit_Handler,
		},
		{
			MethodName: "AutoDeposit",
			Handler:    _Bidder_AutoDeposit_Handler,
		},
		{
			MethodName: "CancelAutoDeposit",
			Handler:    _Bidder_CancelAutoDeposit_Handler,
		},
		{
			MethodName: "GetDeposit",
			Handler:    _Bidder_GetDeposit_Handler,
		},
		{
			MethodName: "GetMinDeposit",
			Handler:    _Bidder_GetMinDeposit_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _Bidder_Withdraw_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendBid",
			Handler:       _Bidder_SendBid_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "bidderapi/v1/bidderapi.proto",
}
