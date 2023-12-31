// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.23.2
// source: service_seq.proto

package service_seq

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

// SequenceServiceClient is the client API for SequenceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SequenceServiceClient interface {
	First(ctx context.Context, in *Seq, opts ...grpc.CallOption) (*Seq, error)
	Second(ctx context.Context, in *Seq, opts ...grpc.CallOption) (*Seq, error)
	Third(ctx context.Context, in *Seq, opts ...grpc.CallOption) (*Seq, error)
}

type sequenceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSequenceServiceClient(cc grpc.ClientConnInterface) SequenceServiceClient {
	return &sequenceServiceClient{cc}
}

func (c *sequenceServiceClient) First(ctx context.Context, in *Seq, opts ...grpc.CallOption) (*Seq, error) {
	out := new(Seq)
	err := c.cc.Invoke(ctx, "/kalisto.tests.examples.service_seq.SequenceService/First", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sequenceServiceClient) Second(ctx context.Context, in *Seq, opts ...grpc.CallOption) (*Seq, error) {
	out := new(Seq)
	err := c.cc.Invoke(ctx, "/kalisto.tests.examples.service_seq.SequenceService/Second", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sequenceServiceClient) Third(ctx context.Context, in *Seq, opts ...grpc.CallOption) (*Seq, error) {
	out := new(Seq)
	err := c.cc.Invoke(ctx, "/kalisto.tests.examples.service_seq.SequenceService/Third", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SequenceServiceServer is the server API for SequenceService service.
// All implementations must embed UnimplementedSequenceServiceServer
// for forward compatibility
type SequenceServiceServer interface {
	First(context.Context, *Seq) (*Seq, error)
	Second(context.Context, *Seq) (*Seq, error)
	Third(context.Context, *Seq) (*Seq, error)
	mustEmbedUnimplementedSequenceServiceServer()
}

// UnimplementedSequenceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSequenceServiceServer struct {
}

func (UnimplementedSequenceServiceServer) First(context.Context, *Seq) (*Seq, error) {
	return nil, status.Errorf(codes.Unimplemented, "method First not implemented")
}
func (UnimplementedSequenceServiceServer) Second(context.Context, *Seq) (*Seq, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Second not implemented")
}
func (UnimplementedSequenceServiceServer) Third(context.Context, *Seq) (*Seq, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Third not implemented")
}
func (UnimplementedSequenceServiceServer) mustEmbedUnimplementedSequenceServiceServer() {}

// UnsafeSequenceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SequenceServiceServer will
// result in compilation errors.
type UnsafeSequenceServiceServer interface {
	mustEmbedUnimplementedSequenceServiceServer()
}

func RegisterSequenceServiceServer(s grpc.ServiceRegistrar, srv SequenceServiceServer) {
	s.RegisterService(&SequenceService_ServiceDesc, srv)
}

func _SequenceService_First_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Seq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SequenceServiceServer).First(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kalisto.tests.examples.service_seq.SequenceService/First",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SequenceServiceServer).First(ctx, req.(*Seq))
	}
	return interceptor(ctx, in, info, handler)
}

func _SequenceService_Second_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Seq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SequenceServiceServer).Second(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kalisto.tests.examples.service_seq.SequenceService/Second",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SequenceServiceServer).Second(ctx, req.(*Seq))
	}
	return interceptor(ctx, in, info, handler)
}

func _SequenceService_Third_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Seq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SequenceServiceServer).Third(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kalisto.tests.examples.service_seq.SequenceService/Third",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SequenceServiceServer).Third(ctx, req.(*Seq))
	}
	return interceptor(ctx, in, info, handler)
}

// SequenceService_ServiceDesc is the grpc.ServiceDesc for SequenceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SequenceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kalisto.tests.examples.service_seq.SequenceService",
	HandlerType: (*SequenceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "First",
			Handler:    _SequenceService_First_Handler,
		},
		{
			MethodName: "Second",
			Handler:    _SequenceService_Second_Handler,
		},
		{
			MethodName: "Third",
			Handler:    _SequenceService_Third_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_seq.proto",
}
