// Package schema resolves semantic paths against complete protobuf identities.
package schema

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

// Message resolves full identities; short names use the shared schema namespace.
func Message(p *protogen.Plugin, name string) (*protogen.Message, error) {
	full := name
	if !strings.Contains(full, ".") {
		full = "kit.sandbox.v1." + name
	}
	var walk func([]*protogen.Message) *protogen.Message
	walk = func(messages []*protogen.Message) *protogen.Message {
		for _, m := range messages {
			if string(m.Desc.FullName()) == full {
				return m
			}
			if found := walk(m.Messages); found != nil {
				return found
			}
		}
		return nil
	}
	for _, file := range p.Files {
		if found := walk(file.Messages); found != nil {
			return found, nil
		}
	}
	return nil, fmt.Errorf("unknown shared message %s", full)
}

// Resolve rejects paths that cross scalar, map, or native data leaves.
func Resolve(p *protogen.Plugin, root, path string) ([]*protogen.Field, error) {
	message, err := Message(p, root)
	if err != nil {
		return nil, err
	}
	var result []*protogen.Field
	parts := strings.Split(path, ".")
	for i, part := range parts {
		var found *protogen.Field
		for _, field := range message.Fields {
			if string(field.Desc.Name()) == part {
				found = field
				break
			}
		}
		if found == nil {
			return nil, fmt.Errorf("unknown path %s.%s", root, path)
		}
		result = append(result, found)
		if i < len(parts)-1 {
			if found.Message == nil || found.Desc.IsMap() {
				return nil, fmt.Errorf("path %s traverses a scalar or map", path)
			}
			if strings.HasPrefix(string(found.Message.Desc.FullName()), "google.protobuf.") || found.Message.Desc.FullName() == "kit.sandbox.v1.MetadataObject" {
				return nil, fmt.Errorf("path %s traverses a native data leaf", path)
			}
			message = found.Message
		}
	}
	return result, nil
}

// Constant derives the Go identifier for a resolved shared field path.
func Constant(root string, fields []*protogen.Field) string {
	prefix := map[string]string{"CreateOptions": "CreateField", "Config": "ConfigField", "SandboxInfo": "InfoField"}[root]
	for _, field := range fields {
		prefix += naming.FieldName(field)
	}
	return prefix
}

// Literal emits a reference to a validated, generated SDK field constant.
func Literal(g *protogen.GeneratedFile, p *protogen.Plugin, root, path string) (string, error) {
	fields, err := Resolve(p, root, path)
	if err != nil {
		return "", err
	}
	return g.QualifiedGoIdent(protogen.GoIdent{GoName: Constant(root, fields), GoImportPath: "github.com/sandbox-kit/kit/sdks/go/sandbox"}), nil
}
