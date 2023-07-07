# Generic GRPC Server


## Overview
This package includes a code generator and some common library functionality to create gRPC servers at Runtime.

Current functionality provides a builder so RPC methods can be bound individually with a default implementation
for unspecificied actions. This guarantees compatibility when proto definitions (including those not bundled with
the implementing server code) are updated. 

## TODO
The server builder and generator are in a stable state right now.

I'd like to extend the capabilities of this library to allow overriding service
names at runtime. This would let users bind multiple versions of the same service
with different names. An early attempt at this was working, but caused issues
with service reflection at runtime. This is due to how proto definitions are
loaded per-file  by the protoc-gen-go-proto generated code during initialization.

Service overriding appears to be working, but there are issues with how the Go runtime loads file descriptors, so it's
not currently callable via `grpcurl`. See: https://github.com/fullstorydev/grpcurl/issues/22#issuecomment-375274465

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

## Running the example

```bash
# example 1: native golang (uses buf for schema generation)
cd example/golang
make run

# example 2: bazel (note: this uses the git master pkg instead of the local version)
cd example/bazel
bazel run //cmd
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
