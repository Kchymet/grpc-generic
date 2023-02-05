package internal

import (
	"context"
	"fmt"
	"github.com/kchymet/generic-grpc/example/api"
)

var _ api.HelloWorldServiceServer = (*Server)(nil)

type Server struct {
}

func (s Server) SayHello(ctx context.Context, request *api.HelloWorldRequest) (*api.HelloWorldResponse, error) {
	message := fmt.Sprintf("Hello, %s!", request.GetName())
	response := &api.HelloWorldResponse{Message: message}
	return response, nil
}


