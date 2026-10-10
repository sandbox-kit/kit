package typesgen

import (
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strconv"
	"strings"
)

// GeneratePaths emits canonical field constants and schema-based response presence.
func GeneratePaths(p *protogen.Plugin, file *protogen.File) error {
	roots := []*protogen.Message{}
	for _, m := range file.Messages {
		switch string(m.Desc.Name()) {
		case "CreateOptions", "Config", "SandboxInfo":
			roots = append(roots, m)
		}
	}
	if len(roots) == 0 {
		return nil
	}
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".paths.gen.go", file.GoImportPath)
	g.P("// Code generated from shared schema paths. DO NOT EDIT.")
	g.P("package ", file.GoPackageName)
	type pathInfo struct {
		constant string
		fields   []*protogen.Field
	}
	var response []pathInfo
	for _, root := range roots {
		var walk func(*protogen.Message, []*protogen.Field, map[protoreflect.FullName]bool)
		walk = func(m *protogen.Message, parent []*protogen.Field, seen map[protoreflect.FullName]bool) {
			if seen[m.Desc.FullName()] {
				return
			}
			seen[m.Desc.FullName()] = true
			defer delete(seen, m.Desc.FullName())
			for _, f := range m.Fields {
				fields := append(append([]*protogen.Field{}, parent...), f)
				names := []string{}
				for _, field := range fields {
					names = append(names, string(field.Desc.Name()))
				}
				path := strings.Join(names, ".")
				constant := schema.Constant(string(root.Desc.Name()), fields)
				g.P("const ", constant, " = ", strconv.Quote(path))
				if string(root.Desc.Name()) == "SandboxInfo" && string(f.Desc.Name()) != "origins" {
					response = append(response, pathInfo{constant: constant, fields: fields})
				}
				if f.Message != nil && !f.Desc.IsMap() && !strings.HasPrefix(string(f.Message.Desc.FullName()), "google.protobuf.") && f.Message.Desc.FullName() != "kit.sandbox.v1.MetadataObject" {
					walk(f.Message, fields, seen)
				}
			}
		}
		walk(root, nil, map[protoreflect.FullName]bool{})
		if string(root.Desc.Name()) == "Config" {
			g.P("const ConfigFieldProvider = \"provider\"")
		}
	}
	if len(response) > 0 {
		g.P("func responseFieldPresent(info *SandboxInfo,path string)bool{if info==nil{return false};switch path{")
		for _, item := range response {
			g.P("case ", item.constant, ":")
			expression := "info"
			valid := true
			for i, f := range item.fields {
				if f.Desc.IsList() && i < len(item.fields)-1 {
					valid = false
					break
				}
				expression += "." + naming.FieldName(f)
				if i < len(item.fields)-1 {
					g.P("if ", expression, "==nil{return false}")
				} else {
					if f.Desc.IsMap() || f.Desc.IsList() || f.Desc.HasPresence() {
						g.P("return ", expression, "!=nil")
					} else if f.Desc.Kind() == protoreflect.StringKind {
						g.P("return ", expression, "!=\"\"")
					} else {
						g.P("return true")
					}
				}
			}
			if !valid {
				g.P("return false")
			}
		}
		g.P("default:return false}}")
		g.P("func validateResponseOrigins(info *SandboxInfo)error{for path,origin:=range info.Origins{if !origin.Valid()||!responseFieldPresent(info,path){return ", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindInvalidResponse"), Provider: "info.Provider", Operation: strconv.Quote("create"), Field: strconv.Quote("origins.") + "+path", Message: strconv.Quote("sandbox-kit: origin must describe a known, present response field")}), "}};return nil}")
	}
	return nil
}
