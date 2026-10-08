// Command protoc-gen-kit-go generates local Go SDK clients from annotations.
package main

import (
	adaptergen "github.com/sandbox-kit/kit/tooling/internal/go/adapter_gen"
	clientgen "github.com/sandbox-kit/kit/tooling/internal/go/client_gen"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(plugin *protogen.Plugin) error {
		for _, file := range plugin.Files {
			if file.Generate {
				if err := clientgen.Generate(plugin, file); err != nil {
					return err
				}
				if err := adaptergen.Generate(plugin, file); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
