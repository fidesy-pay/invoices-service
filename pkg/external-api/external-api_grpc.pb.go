// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: api/external-api/external-api.proto

package external_api

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
	ExternalAPI_GetPrice_FullMethodName = "/external_api.ExternalAPI/GetPrice"
)

// ExternalAPIClient is the client API for ExternalAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ExternalAPIClient interface {
	GetPrice(ctx context.Context, in *GetPriceRequest, opts ...grpc.CallOption) (*GetPriceResponse, error)
}

type externalAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewExternalAPIClient(cc grpc.ClientConnInterface) ExternalAPIClient {
	return &externalAPIClient{cc}
}

func (c *externalAPIClient) GetPrice(ctx context.Context, in *GetPriceRequest, opts ...grpc.CallOption) (*GetPriceResponse, error) {
	out := new(GetPriceResponse)
	err := c.cc.Invoke(ctx, ExternalAPI_GetPrice_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExternalAPIServer is the server API for ExternalAPI service.
// All implementations must embed UnimplementedExternalAPIServer
// for forward compatibility
type ExternalAPIServer interface {
	GetPrice(context.Context, *GetPriceRequest) (*GetPriceResponse, error)
	mustEmbedUnimplementedExternalAPIServer()
}

// UnimplementedExternalAPIServer must be embedded to have forward compatible implementations.
type UnimplementedExternalAPIServer struct {
}

func (UnimplementedExternalAPIServer) GetPrice(context.Context, *GetPriceRequest) (*GetPriceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrice not implemented")
}
func (UnimplementedExternalAPIServer) mustEmbedUnimplementedExternalAPIServer() {}

// UnsafeExternalAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExternalAPIServer will
// result in compilation errors.
type UnsafeExternalAPIServer interface {
	mustEmbedUnimplementedExternalAPIServer()
}

func RegisterExternalAPIServer(s grpc.ServiceRegistrar, srv ExternalAPIServer) {
	s.RegisterService(&ExternalAPI_ServiceDesc, srv)
}

func _ExternalAPI_GetPrice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPriceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExternalAPIServer).GetPrice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ExternalAPI_GetPrice_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExternalAPIServer).GetPrice(ctx, req.(*GetPriceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ExternalAPI_ServiceDesc is the grpc.ServiceDesc for ExternalAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ExternalAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "external_api.ExternalAPI",
	HandlerType: (*ExternalAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPrice",
			Handler:    _ExternalAPI_GetPrice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/external-api/external-api.proto",
}