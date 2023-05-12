package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	if len(service.Methods) == 0 {
		return
	}

	for _, method := range service.Methods {
		genServiceMethodAction(gen, file, g, service, method)
		genServiceMethodRequestDecode(gen, file, g, service, method)
	}
}

func genServiceMethodAction(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	interfaceName := getActionInterfaceName(service, method)
	g.P("type ", interfaceName, " interface {")
	signature := getServerSignature(g, method)
	g.P(signature)
	g.P("}")
}

func genServiceMethodRequestDecode(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	decodeMethodName := getDecodeMethodName(service, method)
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		g.P("func ", decodeMethodName, "(dec func(interface{}) error) (interface{}, error) {")
		g.P("in := new(", method.Input.GoIdent, ")")
		g.P("if err := dec(in); err != nil { return nil, err }")
		g.P("return in, nil")
		g.P("}")
		return
	}
	if !method.Desc.IsStreamingClient() {
		g.P("func ", decodeMethodName, "(dec func(interface{}) error, stream grpc.ServerStream) (interface{}, error) {")
		g.P("m := new(", method.Input.GoIdent, ")")
		g.P("if err := stream.RecvMsg(m); err != nil { return nil, err }")
		g.P("return m, nil")
		g.P("}")
	}
}

func getServerSignature(g *protogen.GeneratedFile, method *protogen.Method) string {
	var reqArgs []string
	ret := "error"
	if !method.Desc.IsStreamingClient() && !method.Desc.IsStreamingServer() {
		reqArgs = append(reqArgs, g.QualifiedGoIdent(contextPackage.Ident("Context")))
		ret = "(*" + g.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
	}
	if !method.Desc.IsStreamingClient() {
		reqArgs = append(reqArgs, "*"+g.QualifiedGoIdent(method.Input.GoIdent))
	}
	if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
		reqArgs = append(reqArgs, method.Parent.GoName+"_"+method.GoName+"Server")
	}
	return method.GoName + "(" + strings.Join(reqArgs, ", ") + ") " + ret
}

func getActionInterfaceName(service *protogen.Service, method *protogen.Method) string {
	interfaceName := fmt.Sprintf("%s%s%s", service.GoName, method.GoName, "Action")
	return interfaceName
}
