# Generic GRPC Server

## TODO
WIP - currently only unary services are handled fully. Client and server streaming use cases are coming soon!

Service overriding appears to be working, but there are issues with how the Go runtime loads file descriptors, so it's
not currently callable via `grpcurl`. See: https://github.com/fullstorydev/grpcurl/issues/22#issuecomment-375274465

## Overview
This package includes a code generator and some common library functionality to create gRPC servers at Runtime.

Current functionality provides a builder so RPC methods can be bound individually with a default implementation
for unspecificied actions. This guarantees compatibility when proto definitions (including those not bundled with
the implementing server code) are updated. 

## Generated Types

Assuming a service called `HelloWorldService` with unary method `SayHello`, we get the following definitions:

* `HelloWorldServiceBuilder` and `NewHelloWorldServiceBuilder()` - a builder and initializer to bind individual service actions.
* `HelloWorld_SayHello_Action` - an interface for `SayHello` handler to implement. This is method-level vs service level.

This plugs into the existing grpc network stack for Go, and exists as an alternative to the grpc-go generated interfaces.

Intended usage is:
```go
server := grpc.NewServer()
service := NewHelloWorldServiceBuilder().
	BindSayHello(/* some struct implementing HelloWorld_SayHello_Action */).
	Build().
    Register(server)

// other server initialization
server.Run()

```

## Overriding Service Name
If you want to run multiple instances of the same service, you can override the service name when binding:
```go
server := grpc.NewServer()

// Accessible via HelloWorldService.SayHello
service1 := NewHelloWorldServiceBuilder().
	BindSayHello(/* some struct implementing HelloWorld_SayHello_Action */).
	Build().
    Register(server)

// Accessible via HelloWorldService2.SayHello
service2 := NewHelloWorldServiceBuilder().
    BindSayHello(/* some struct implementing HelloWorld_SayHello_Action */).
	WithServiceName("HelloWorldService2")
    Build().
    Register(server)

// other server initialization
server.Run()
```

## Running the example

```bash
bazel run //example/cmd
```

```bash
grpcurl -d '{"name": "World"}' --plaintext localhost:9000 HelloWorldService.SayHello

# This currently fails, but the service is exposed.
# Error invoking method "HelloWorldService2.SayHello": target server does not expose service "HelloWorldService2"
grpcurl -d '{"name": "World"}' --plaintext localhost:9000 HelloWorldService2.SayHello
```

```
# should output
{
 "message": "Hello, World"
}
```