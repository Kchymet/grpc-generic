package action

import (
	"fmt"
	"github.com/kchymet/grpc-generic/example/golang/api"
	"io"
)

func NewStreamHello() api.HelloWorldServiceStreamHelloAction {
	return &streamHelloAction{}
}

type streamHelloAction struct {
}

func (s *streamHelloAction) StreamHello(server api.HelloWorldService_StreamHelloServer) error {
	for {
		request, err := server.Recv()
		if err == io.EOF {
			fmt.Printf("closing stream")
		} else if err != nil {
			fmt.Printf("failed to receive, err: %v", err)
			return err
		}
		name := request.GetName()
		message := fmt.Sprintf("Hello, %s", name)
		err = server.Send(&api.HelloWorldResponse{
			Message: message,
		})
		if err != nil {
			fmt.Printf("failed to send, err: %v", err)
			return err
		}
	}
}
