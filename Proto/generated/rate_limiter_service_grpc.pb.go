// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.0
// source: Proto/source/rate_limiter_service.proto

package generated

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
	RateLimitService_GetApplicationLayerRateLimit_FullMethodName = "/rate_limit_service/GetApplicationLayerRateLimit"
	RateLimitService_GetHealth_FullMethodName                    = "/rate_limit_service/GetHealth"
)

// RateLimitServiceClient is the client API for RateLimitService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// RateLimiter service definition
type RateLimitServiceClient interface {
	GetApplicationLayerRateLimit(ctx context.Context, in *RateLimitRequest, opts ...grpc.CallOption) (*RateLimitResponse, error)
	GetHealth(ctx context.Context, in *Void, opts ...grpc.CallOption) (*HealthCheckResponse, error)
}

type rateLimitServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRateLimitServiceClient(cc grpc.ClientConnInterface) RateLimitServiceClient {
	return &rateLimitServiceClient{cc}
}

func (c *rateLimitServiceClient) GetApplicationLayerRateLimit(ctx context.Context, in *RateLimitRequest, opts ...grpc.CallOption) (*RateLimitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RateLimitResponse)
	err := c.cc.Invoke(ctx, RateLimitService_GetApplicationLayerRateLimit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rateLimitServiceClient) GetHealth(ctx context.Context, in *Void, opts ...grpc.CallOption) (*HealthCheckResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthCheckResponse)
	err := c.cc.Invoke(ctx, RateLimitService_GetHealth_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RateLimitServiceServer is the server API for RateLimitService service.
// All implementations must embed UnimplementedRateLimitServiceServer
// for forward compatibility.
//
// RateLimiter service definition
type RateLimitServiceServer interface {
	GetApplicationLayerRateLimit(context.Context, *RateLimitRequest) (*RateLimitResponse, error)
	GetHealth(context.Context, *Void) (*HealthCheckResponse, error)
	mustEmbedUnimplementedRateLimitServiceServer()
}

// UnimplementedRateLimitServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRateLimitServiceServer struct{}

func (UnimplementedRateLimitServiceServer) GetApplicationLayerRateLimit(context.Context, *RateLimitRequest) (*RateLimitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetApplicationLayerRateLimit not implemented")
}
func (UnimplementedRateLimitServiceServer) GetHealth(context.Context, *Void) (*HealthCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetHealth not implemented")
}
func (UnimplementedRateLimitServiceServer) mustEmbedUnimplementedRateLimitServiceServer() {}
func (UnimplementedRateLimitServiceServer) testEmbeddedByValue()                          {}

// UnsafeRateLimitServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RateLimitServiceServer will
// result in compilation errors.
type UnsafeRateLimitServiceServer interface {
	mustEmbedUnimplementedRateLimitServiceServer()
}

func RegisterRateLimitServiceServer(s grpc.ServiceRegistrar, srv RateLimitServiceServer) {
	// If the following call pancis, it indicates UnimplementedRateLimitServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&RateLimitService_ServiceDesc, srv)
}

func _RateLimitService_GetApplicationLayerRateLimit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RateLimitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimitServiceServer).GetApplicationLayerRateLimit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimitService_GetApplicationLayerRateLimit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimitServiceServer).GetApplicationLayerRateLimit(ctx, req.(*RateLimitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RateLimitService_GetHealth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RateLimitServiceServer).GetHealth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: RateLimitService_GetHealth_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RateLimitServiceServer).GetHealth(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

// RateLimitService_ServiceDesc is the grpc.ServiceDesc for RateLimitService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RateLimitService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rate_limit_service",
	HandlerType: (*RateLimitServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetApplicationLayerRateLimit",
			Handler:    _RateLimitService_GetApplicationLayerRateLimit_Handler,
		},
		{
			MethodName: "GetHealth",
			Handler:    _RateLimitService_GetHealth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "Proto/source/rate_limiter_service.proto",
}
