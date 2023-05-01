Generic GRPC Server

# TODO
WIP - currently only unary services are handled fully. Client and server streaming use cases are coming soon!

# Overview
This package includes a code generator and some common library functionality to create gRPC servers at Runtime.

Current functionality provides a builder so RPC methods can be bound individually with a default implementation
for unspecificied actions. This guarantees compatibility when proto definitions (including those not bundled with
the implementing server code) are updated. 

# Generated Types

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

# Running the example
```bash
bazel run //example/cmd
```

```bash
grpcurl -d '{"name": "Kyle"}' --plaintext localhost:9000 HelloWorldService.SayHello
```

```
# should output
{
 "message": "Hello, Kyle"
}
```