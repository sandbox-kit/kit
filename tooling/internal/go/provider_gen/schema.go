package providergen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func sharedMessage(p *protogen.Plugin, name string) (*protogen.Message, error) {
	return schema.Message(p, name)
}
func sharedField(p *protogen.Plugin, message, name string) (*protogen.Field, error) {
	fields, err := schema.Resolve(p, message, name)
	if err != nil {
		return nil, err
	}
	return fields[len(fields)-1], nil
}

// capture copies schema fields without changing absent versus explicit values.
// Unsupported ownership shapes are rejected rather than shallow copied.
func capture(g *protogen.GeneratedFile, p *protogen.Plugin, path, typ string, optional bool, dest string) error {
	parts := strings.Split(path, ".")
	message := "Config"
	source := "config"
	guards := 0
	var field *protogen.Field
	for i, part := range parts {
		var err error
		field, err = sharedField(p, message, part)
		if err != nil {
			return err
		}
		source += "." + naming.FieldName(field)
		if i < len(parts)-1 {
			if field.Message == nil {
				return fmt.Errorf("capture path %q traverses scalar", path)
			}
			g.P("if ", source, "!=nil{")
			guards++
			message = string(field.Message.Desc.FullName())
		}
	}
	if field.Message != nil {
		if !optional || field.Message.GoIdent.GoName != typ {
			return fmt.Errorf("capture %s requires optional %s", path, field.Message.GoIdent.GoName)
		}
		g.P("if ", source, "!=nil{")
		g.P(dest, "=&", field.Message.GoIdent, "{}")
		if err := copyMessage(g, field.Message, source, dest, map[protoreflect.FullName]bool{}); err != nil {
			return err
		}
		g.P("}")
	} else {
		if field.Desc.Kind() != protoreflect.StringKind || typ != "string" || !optional || !field.Desc.HasPresence() {
			return fmt.Errorf("capture %s requires optional string state", path)
		}
		g.P("if ", source, "!=nil{value:=*", source, ";", dest, "=&value}")
	}
	for i := 0; i < guards; i++ {
		g.P("}")
	}
	return nil
}
func copyMessage(g *protogen.GeneratedFile, m *protogen.Message, source, dest string, seen map[protoreflect.FullName]bool) error {
	if seen[m.Desc.FullName()] {
		return fmt.Errorf("recursive state capture %s is unsupported", m.Desc.FullName())
	}
	seen[m.Desc.FullName()] = true
	defer delete(seen, m.Desc.FullName())
	for _, f := range m.Fields {
		s := source + "." + naming.FieldName(f)
		d := dest + "." + naming.FieldName(f)
		if f.Desc.IsList() || f.Desc.IsMap() || f.Desc.Kind() == protoreflect.BytesKind || (f.Oneof != nil && !f.Oneof.Desc.IsSynthetic()) {
			return fmt.Errorf("state capture ownership unsupported for %s", f.Desc.FullName())
		}
		if f.Message != nil {
			g.P("if ", s, "!=nil{", d, "=&", f.Message.GoIdent, "{}")
			if err := copyMessage(g, f.Message, s, d, seen); err != nil {
				return err
			}
			g.P("}")
		} else if f.Desc.HasPresence() {
			g.P("if ", s, "!=nil{value:=*", s, ";", d, "=&value}")
		} else {
			g.P(d, "=", s)
		}
	}
	return nil
}
