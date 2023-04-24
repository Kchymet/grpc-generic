package main

import (
	"flag"
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

const (
	contextPackage = protogen.GoImportPath("context")
	grpcPackage    = protogen.GoImportPath("google.golang.org/grpc")
	codesPackage   = protogen.GoImportPath("google.golang.org/grpc/codes")
	statusPackage  = protogen.GoImportPath("google.golang.org/grpc/status")
)

const version = "0.0.1"

// FileDescriptorProto.package field number
const fileDescriptorProtoPackageFieldNumber = 2

// FileDescriptorProto.syntax field number
const fileDescriptorProtoSyntaxFieldNumber = 12

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-generic-grpc %v\n", version)
		return
	}

	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				panic("not generating")
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

// generateFile generates a _ascii.pb.go file containing gRPC service definitions.
func generateFile(gen *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		panic("no services")
		return nil
	}

	filename := file.GeneratedFilenamePrefix + "_generic_grpc.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	// Attach all comments associated with the syntax field.
	genLeadingComments(g, file.Desc.SourceLocations().ByPath(protoreflect.SourcePath{fileDescriptorProtoSyntaxFieldNumber}))
	g.P("// Code generated by protoc-gen-go-generic-grpc. DO NOT EDIT.")
	g.P("// versions:")
	g.P("// - protoc-gen-go-generic-grpc v", version)
	g.P("// - protoc             ", protocVersion(gen))
	if file.Proto.GetOptions().GetDeprecated() {
		g.P("// ", file.Desc.Path(), " is a deprecated file.")
	} else {
		g.P("// source: ", file.Desc.Path())
	}
	g.P()
	// Attach all comments associated with the package field.
	genLeadingComments(g, file.Desc.SourceLocations().ByPath(protoreflect.SourcePath{fileDescriptorProtoPackageFieldNumber}))
	g.P("package ", file.GoPackageName)
	g.P()
	generateFileContent(gen, file, g)
	return g
}

func genLeadingComments(g *protogen.GeneratedFile, loc protoreflect.SourceLocation) {
	for _, s := range loc.LeadingDetachedComments {
		g.P(protogen.Comments(s))
		g.P()
	}
	if s := loc.LeadingComments; s != "" {
		g.P(protogen.Comments(s))
		g.P()
	}
}

func protocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

// generateFileContent generates the gRPC service definitions, excluding the package statement.
func generateFileContent(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile) {
	if len(file.Services) == 0 {
		return
	}

	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the grpc package it is being compiled against.")
	g.P("// Requires gRPC-Go v1.32.0 or later.")
	g.P("const _ = ", grpcPackage.Ident("SupportPackageIsVersion7")) // When changing, update version number above.
	g.P()
	for _, service := range file.Services {
		genService(gen, file, g, service)
		genServiceBuilder(gen, file, g, service)
	}
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	if len(service.Methods) == 0 {
		return
	}

	for _, method := range service.Methods {
		genServiceMethodAction(gen, file, g, service, method)
	}
}

func genServiceBuilder(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) {
	if len(service.Methods) == 0 {
		return
	}

	// Generate the builder interface
	//unexportedName :=
	interfaceName := fmt.Sprintf("%sBuilder", unexport(service.GoName))
	g.P("type ", interfaceName, " interface {")
	for _, method := range service.Methods {
		genServiceBindingMethodSignature(gen, file, g, service, method)
	}
	//genBuildSignature(gen, file, g, service, method)
	g.P("}")

	// Generate a builder implementation
	builderName := fmt.Sprintf("%sBuilderImpl", unexport(service.GoName))
	g.P("type ", builderName, " struct {")
	// TODO builder dependencies / embeddings
	g.P("}")
	for _, method := range service.Methods {
		genServiceBindingMethod(gen, file, g, service, method, builderName)
	}


	// Generate static builder implementation
	constructorName := fmt.Sprintf("New%sBuilder", service.GoName)
	g.P("func ", constructorName, "() ", interfaceName, " {")
	g.P("return &", builderName, "{}")
	g.P("}")

}

func genServiceMethodAction(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	interfaceName := getActionInterfaceName(service, method)
	g.P("type ", interfaceName, " interface {")
	signature := getServerSignature(g, method)
	g.P(signature)
	g.P("}")
}

func genServiceBindingMethodSignature(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method) {
	actionInterfaceName := getActionInterfaceName(service, method)
	methodName := getBindMethodName(method)
	g.P(methodName, "(",actionInterfaceName,")")
}

func genServiceBindingMethod(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service, method *protogen.Method, builderName string) {
	actionInterfaceName := getActionInterfaceName(service, method)
	methodName := getBindMethodName(method)
	g.P("func (b *", builderName, ") ", methodName, "(a ", actionInterfaceName, ") {")

	// TODO implementation

	g.P("}")
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

func getBindMethodName(method *protogen.Method) string {
	return fmt.Sprintf("Bind%s", method.GoName)
}

func getActionInterfaceName(service *protogen.Service, method *protogen.Method) string {
	interfaceName := fmt.Sprintf("%s%s%s", service.GoName, method.GoName, "Action")
	return interfaceName
}

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }