package providergen

import (
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"google.golang.org/protobuf/compiler/protogen"
)

func sharedField(p *protogen.Plugin, message, name string) (*protogen.Field, error) {
	fields, err := schema.Resolve(p, message, name)
	if err != nil {
		return nil, err
	}
	return fields[len(fields)-1], nil
}
