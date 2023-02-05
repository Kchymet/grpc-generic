package main

import (
	"fmt"
	"github.com/kchymet/generic-grpc/example/api"
	"github.com/kchymet/generic-grpc/example/internal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main(){
	port := 9000

	server := grpc.NewServer()
	appServer := internal.Server{}
	api.RegisterHelloWorldServiceServer(server, appServer)
	reflection.Register(server)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}