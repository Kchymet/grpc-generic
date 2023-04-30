package action

import (
	"context"
	"fmt"
	"github.com/kchymet/generic-grpc/example/api"
)

func NewSayHelloAction() api.HelloWorldServiceSayHelloAction {
	return &sayHelloAction{}
}

type sayHelloAction struct {

}

func (s *sayHelloAction) SayHello(_ context.Context, request *api.HelloWorldRequest) (*api.HelloWorldResponse, error) {
	name := request.GetName()
	message := fmt.Sprintf("Hello, %s", name)
	return &api.HelloWorldResponse{
		Message: message,
	}, nil
}

