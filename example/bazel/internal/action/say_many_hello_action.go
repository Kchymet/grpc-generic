package action

import (
	"fmt"
	"github.com/kchymet/grpc-generic/example/api"
)

func NewSayManyHello() api.HelloWorldServiceSayManyHelloAction {
	return &sayManyHelloAction{}
}

type sayManyHelloAction struct {

}

func (s *sayManyHelloAction) SayManyHello(request *api.ManyHelloRequest, server api.HelloWorldService_SayManyHelloServer) error {
	count := int(request.GetCount())
	name := request.GetName()

	for i := 0; i < count; i++{
		message := fmt.Sprintf("Hello #%d, %s", i, name)
		server.Send(&api.ManyHelloResponse{
			Message: message,
		})
	}
	return nil
}

