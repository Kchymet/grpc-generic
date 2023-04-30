package action

import (
	"fmt"
	"github.com/kchymet/generic-grpc/example/api"
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
			fmt.Println("closing stream")
		} else if err != nil {
			fmt.Println("failed to receive, err: %v", err)
			return err
		}
		name := request.GetName()
		message := fmt.Sprintf("Hello, %s", name)
		err = server.Send(&api.HelloWorldResponse{
			Message: message,
		})
		if err != nil {
			fmt.Println("failed to send, err: %v", err)
			return err
		}
	}
}
