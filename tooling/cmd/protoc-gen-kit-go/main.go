// Command protoc-gen-kit-go compiles portable SDK specifications and emits Go.
package main

import (
	goemit "github.com/sandbox-kit/kit/tooling/internal/go"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		files := make([]protoreflect.FileDescriptor, 0, len(plugin.Files))
		for _, file := range plugin.Files {
			files = append(files, file.Desc)
		}
		program, err := model.Compile("specs", files)
		if err != nil {
			return err
		}
		for _, file := range plugin.Files {
			if file.Generate {
				if err := goemit.Generate(plugin, file, program); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
