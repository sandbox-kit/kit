// Package schema resolves semantic paths against complete protobuf identities.
package schema

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strings"
)

// Index adapts the language-neutral descriptor index to protogen input.
func Index(p *protogen.Plugin) *model.Schema {
	files := make([]protoreflect.FileDescriptor, 0, len(p.Files))
	for _, file := range p.Files {
		files = append(files, file.Desc)
	}
	return model.NewSchema(files)
}

// Message adapts a resolved descriptor to its Go emission metadata.
func Message(p *protogen.Plugin, name string) (*protogen.Message, error) {
	descriptor, err := Index(p).Message(name)
	if err != nil {
		return nil, err
	}
	var walk func([]*protogen.Message) *protogen.Message
	walk = func(messages []*protogen.Message) *protogen.Message {
		for _, m := range messages {
			if m.Desc.FullName() == descriptor.FullName() {
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
	return nil, fmt.Errorf("missing Go metadata for %s", descriptor.FullName())
}

// Resolve uses portable path semantics and attaches Go names afterward.
func Resolve(p *protogen.Plugin, root, path string) ([]*protogen.Field, error) {
	resolved, err := Index(p).Resolve(root, path)
	if err != nil {
		return nil, err
	}
	result := make([]*protogen.Field, 0, len(resolved.Fields))
	for _, descriptor := range resolved.Fields {
		field, err := Field(p, descriptor)
		if err != nil {
			return nil, err
		}
		result = append(result, field)
	}
	return result, nil
}

// Field attaches Go emission metadata to an already-resolved portable field.
func Field(p *protogen.Plugin, descriptor protoreflect.FieldDescriptor) (*protogen.Field, error) {
	var walk func([]*protogen.Message) *protogen.Field
	walk = func(messages []*protogen.Message) *protogen.Field {
		for _, m := range messages {
			if m.Desc.FullName() == descriptor.ContainingMessage().FullName() {
				for _, field := range m.Fields {
					if field.Desc.Number() == descriptor.Number() {
						return field
					}
				}
			}
			if found := walk(m.Messages); found != nil {
				return found
			}
		}
		return nil
	}
	if file := p.FilesByPath[descriptor.ParentFile().Path()]; file != nil {
		if found := walk(file.Messages); found != nil {
			return found, nil
		}
	}
	return nil, fmt.Errorf("missing Go metadata for %s", descriptor.FullName())
}

// Constant derives the Go identifier for a resolved shared field path.
func Constant(root string, fields []*protogen.Field, rules spec.FieldPathTemplates) string {
	prefix := rules.Roots[root]
	names := ""
	for _, field := range fields {
		names += naming.FieldName(field)
	}
	return strings.NewReplacer("{prefix}", prefix, "{fields}", names).Replace(rules.Name)
}

// Literal emits a reference to a validated, generated SDK field constant.
func Literal(g *protogen.GeneratedFile, p *protogen.Plugin, root, path string, rules spec.FieldPathTemplates) (string, error) {
	if rules.Roots[root] == "" {
		return "", fmt.Errorf("missing field-path root %s", root)
	}
	fields, err := Resolve(p, root, path)
	if err != nil {
		return "", err
	}
	return g.QualifiedGoIdent(protogen.GoIdent{GoName: Constant(root, fields, rules), GoImportPath: "github.com/sandbox-kit/kit/sdks/go/sandbox"}), nil
}
