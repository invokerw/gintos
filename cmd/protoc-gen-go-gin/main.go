package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	showVersion     = flag.Bool("version", false, "print the version and exit")
	omitempty       = flag.Bool("omitempty", true, "omit if google.api is empty")
	omitemptyPrefix = flag.String("omitempty_prefix", "", "omit if google.api is empty")
	rbacOutPut      = flag.String("rbac_path", "", "rbac file output path")
	rbacPackageName = flag.String("rbac_package_name", "", "rbac package name")
	servicePath     = flag.String("service_path", "", "service path, gen service handler")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-go-gin %v\n", release)
		return
	}
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		var allServiceDesc []*serviceDesc
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			sds := generateFile(gen, f, *omitempty, *omitemptyPrefix)
			if sds != nil {
				allServiceDesc = append(allServiceDesc, sds...)
			}
		}
		if *rbacOutPut != "" {
			if err := rbacGenerate(gen, *rbacOutPut, *rbacPackageName, *omitempty, *omitemptyPrefix); err != nil {
				return err
			}
		}
		if *servicePath != "" {
			if err := serviceGenerate(gen, *servicePath, allServiceDesc); err != nil {
				return err
			}
		}
		return nil
	})
}
