package pkg

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ dispatchingRpcServer = (*dispatchingRpcServerImpl)(nil)

type DispatchingRpcServer interface {
	// Register binds this dispatcher to a grpc service.
	Register(grpc.ServiceRegistrar)

	// GetRegistration returns a service description and interface. This can be passed directly to a grpc registrar's RegisterService method.
	GetRegistration() (*grpc.ServiceDesc, interface{})
}

type UnaryDispatchInfo struct {
	// DecodeFunc takes the decoding function from the grpc core library, runs it, and returns either the initialized interface or an error.
	DecodeFunc func(func(interface{}) error) (interface{}, error)
	Handler grpc.UnaryHandler
}

// This is an unexported interface for ensuring safe binding without allowing custom implementations.
type dispatchingRpcServer interface {
	DispatchingRpcServer
}

type dispatchingRpcServerImpl struct {
	ServiceName string
	Metadata string

	MethodDispatchInfo map[string]*UnaryDispatchInfo
	StreamDescriptions map[string]*grpc.StreamDesc
}

func NewDispatchingRpcServer(serviceName, metadata string, methodHandlers map[string]*UnaryDispatchInfo, streamDescriptions map[string]*grpc.StreamDesc) DispatchingRpcServer {
	return &dispatchingRpcServerImpl{
		ServiceName:        serviceName,
		Metadata:           metadata,
		MethodDispatchInfo: methodHandlers,
		StreamDescriptions: streamDescriptions,
	}
}

func (s *dispatchingRpcServerImpl) Register(r grpc.ServiceRegistrar) {
	r.RegisterService(s.GetRegistration())
}

func (s *dispatchingRpcServerImpl) GetServiceDescription() *grpc.ServiceDesc {
	md := make([]grpc.MethodDesc, 0, len(s.MethodDispatchInfo))
	for name, info := range s.MethodDispatchInfo {
		md = append(md, s.unaryMethodDescription(name, info))
	}

	sd := make([]grpc.StreamDesc, 0, len(s.StreamDescriptions))
	for _, info := range s.StreamDescriptions {
		sd = append(sd, *info)
	}

	return &grpc.ServiceDesc{
		ServiceName: s.ServiceName,
		HandlerType: (*dispatchingRpcServer)(nil),
		Methods:     md,
		Streams:     sd,
		Metadata:    s.Metadata,
	}
}

func (s *dispatchingRpcServerImpl) GetRegistration() (*grpc.ServiceDesc, interface{}) {
	return s.GetServiceDescription(), s
}

func (s *dispatchingRpcServerImpl) unaryMethodDescription(methodName string, info *UnaryDispatchInfo) grpc.MethodDesc {
	return grpc.MethodDesc{
		MethodName: methodName,
		Handler:    s.unaryHandler(methodName, info),
	}
}

func (s *dispatchingRpcServerImpl) streamDescription() grpc.StreamDesc {
	return grpc.StreamDesc{
		StreamName:    "",
		Handler:       nil,
		ServerStreams: false,
		ClientStreams: false,
	}
}

func (s *dispatchingRpcServerImpl) unaryHandler(methodName string, info *UnaryDispatchInfo) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		// TODO type safety checks on decode / handler functions (former should be input to latter)
		in, err := info.DecodeFunc(dec)
		if err != nil {
			return nil, err
		}

		if interceptor == nil {
			return info.Handler(ctx, in)
		}

		// Initialize server info for interceptors
		fullMethodName := fmt.Sprintf("/%s/%s", s.ServiceName, methodName)
		grpcInfo := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: fullMethodName,
		}

		// Call the interceptor chain
		return interceptor(ctx, in, grpcInfo, info.Handler)
	}
}


func GetDefaultUnaryUnimplementedHandler(methodName string) func(context.Context, interface{}) (interface{}, error) {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		return nil, status.Errorf(codes.Unimplemented, "method %s not implemented", methodName)
	}
}

func GetDefaultStreamUnimplementedHandler(methodName string) func(interface{}, grpc.ServerStream) error {
	return func(_ interface{}, stream grpc.ServerStream) error {
		return status.Errorf(codes.Unimplemented, "streaming method %s not implemented", methodName)
	}
}
