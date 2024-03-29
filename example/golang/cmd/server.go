package main

import (
	"fmt"
	"github.com/kchymet/grpc-generic/example/golang/api"
	"github.com/kchymet/grpc-generic/example/golang/internal/action"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	port := 9000

	server := grpc.NewServer()

	reflection.Register(server)

	api.NewHelloWorldServiceBuilder().
		BindSayHello(action.NewSayHelloAction()).
		BindSayManyHello(action.NewSayManyHello()).
		BindStreamHello(action.NewStreamHello()).
		// Don't bind the unimplemented method from the proto.
		Build().
		Register(server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
