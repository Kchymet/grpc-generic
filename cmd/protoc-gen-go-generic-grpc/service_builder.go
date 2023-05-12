package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strconv"
)

func genServiceBuilder(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	if len(service.Methods) == 0 {
		return
	}

	// Generate the builder interface
	//unexportedName :=
	interfaceName := fmt.Sprintf("%sBuilder", unexport(service.GoName))
	g.P("type ", interfaceName, " interface {")
	for _, method := range service.Methods {
		genServiceBindingMethodSignature(gen, file, g, service, method, interfaceName)
	}
	genWithServiceNameSignature(gen, file, g, service, interfaceName)
	genBuildMethodSignature(gen, file, g, service)
	g.P("}")

	// Generate a builder implementation
	builderName := fmt.Sprintf("%sBuilderImpl", unexport(service.GoName))
	g.P("type ", builderName, " struct {")
	g.P("ServiceName string")
	g.P("Metadata string")
	g.P("MethodDispatchInfo map[string]*", g.QualifiedGoIdent(dispatchingServerPackage.Ident("UnaryDispatchInfo")))
	g.P("StreamDescriptions map[string]", g.QualifiedGoIdent(grpcPackage.Ident("StreamHandler")))
	g.P("}")
	for _, method := range service.Methods {
		genServiceBindingMethod(gen, file, g, service, method, builderName, interfaceName)
	}
	genWithServiceName(gen, file, g, service, builderName, interfaceName)
	genBuildMethod(gen, file, g, service, builderName)

	genServiceBuilderConstructor(service, g, interfaceName, builderName, file, service.Methods)
}

func genServiceBuilderConstructor(service *protogen.Service, g *protogen.GeneratedFile, interfaceName string, builderName string, file *protogen.File, methods []*protogen.Method) {
	constructorName := fmt.Sprintf("New%sBuilder", service.GoName)
	g.P("func ", constructorName, "() ", interfaceName, " {")

	g.P("return &", builderName, "{")
	g.P("ServiceName: \"", service.GoName, "\",")
	g.P("Metadata: \"", file.Desc.Path(), "\",")

	g.P("MethodDispatchInfo: map[string]*", g.QualifiedGoIdent(dispatchingServerPackage.Ident("UnaryDispatchInfo")), "{")
	for _, method := range methods {
		if method.Desc.IsStreamingServer() || method.Desc.IsStreamingClient() {
			continue
		}
		g.P("\"", method.GoName, "\": &", g.QualifiedGoIdent(dispatchingServerPackage.Ident("UnaryDispatchInfo")), "{")
		g.P("DecodeFunc: ", getDecodeMethodName(service, method), ",")
		g.P("Handler: ", g.QualifiedGoIdent(dispatchingServerPackage.Ident("GetDefaultUnaryUnimplementedHandler")), "(", strconv.Quote(string(method.Desc.Name())), "),")
		g.P("},")
	}
	g.P("},") // End of MethodDispatchInfo initialization.

	g.P("}") // End of Builder initialization
	g.P("}") // End of Constructor method
}

func genServiceBindingMethodSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method, interfaceName string) {
	actionInterfaceName := getActionInterfaceName(service, method)
	methodName := getBindMethodName(method)
	g.P(methodName, "(", actionInterfaceName, ") ", interfaceName)
}

func genServiceBindingMethod(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method, builderName, interfaceName string) {
	actionInterfaceName := getActionInterfaceName(service, method)
	methodName := getBindMethodName(method)
	g.P("func (b *", builderName, ") ", methodName, "(a ", actionInterfaceName, ") ", interfaceName, " {")

	// Bind MethodDesc for unary (non-streaming) endpoints.
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		g.P("b.MethodDispatchInfo[", strconv.Quote(string(method.Desc.Name())), "].Handler = func(ctx ", g.QualifiedGoIdent(contextPackage.Ident("Context")), ", req interface{}) (interface{}, error) {")
		g.P("return a.", method.GoName, "(ctx, req.(*", method.Input.GoIdent, "))")
		g.P("}")
	}

	g.P("return b")

	g.P("}")
}

func genBuildMethodSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	serverName := g.QualifiedGoIdent(dispatchingServerPackage.Ident("DispatchingRpcServer"))
	g.P("Build() ", serverName)
}

func genBuildMethod(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, builderName string) {
	serverName := g.QualifiedGoIdent(dispatchingServerPackage.Ident("DispatchingRpcServer"))

	g.P("func (b *", builderName, ") Build() ", serverName, " {")

	g.P("return ", g.QualifiedGoIdent(dispatchingServerPackage.Ident("NewDispatchingRpcServer")), "(b.ServiceName, b.Metadata, b.MethodDispatchInfo, b.StreamDescriptions)")

	g.P("}")
}

func genWithServiceNameSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, interfaceName string) {
	g.P("WithServiceName(string) ", interfaceName)
}

func genWithServiceName(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, builderName, interfaceName string) {
	g.P("func (b *", builderName, ") WithServiceName(serviceName string) ", interfaceName, " {")

	g.P("b.ServiceName = serviceName")
	g.P("return b")

	g.P("}")
}

func getBindMethodName(method *protogen.Method) string {
	return fmt.Sprintf("Bind%s", method.GoName)
}
