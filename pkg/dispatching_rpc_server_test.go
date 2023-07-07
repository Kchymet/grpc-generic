package pkg

import (
	"testing"

	"github.com/kchymet/grpc-generic/mocks"
	"go.uber.org/mock/gomock"
)

func Test_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	registrar := mocks.NewMockServiceRegistrar(ctrl)
	server := NewDispatchingRpcServer("test", "test", nil, nil)

	registrar.EXPECT().RegisterService(server.GetRegistration())

	server.Register(registrar)
}
