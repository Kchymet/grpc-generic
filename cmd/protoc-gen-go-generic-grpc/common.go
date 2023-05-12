package main

import (
	"fmt"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

func getDecodeMethodName(service *protogen.Service, method *protogen.Method) string {
	return fmt.Sprintf("_%s_%s_Decode", service.GoName, method.GoName)
}


func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }