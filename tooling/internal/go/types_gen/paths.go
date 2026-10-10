package typesgen

import (
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"strconv"
	"strings"
)

// GeneratePaths emits canonical field constants and schema-based response presence.
func GeneratePaths(p *protogen.Plugin, file *protogen.File, templates model.Templates, profile spec.LanguageProfile) error {
	roots := []*protogen.Message{}
	for _, m := range file.Messages {
		if _, ok := templates.FieldPaths.Roots[string(m.Desc.Name())]; ok {
			roots = append(roots, m)
		}
	}
	if len(roots) == 0 {
		return nil
	}
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".paths.gen.go", file.GoImportPath)
	declarationgen.Banner(g, "",
		"Field-path constants identify settings in errors. Their values are canonical protobuf paths.",
		"fmt.Println(CreateFieldResourcesCPUCores) // resources.cpu_cores")
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
				constant := schema.Constant(string(root.Desc.Name()), fields, templates.FieldPaths)
				declarationgen.WriteComment(g, constant, "is the canonical path "+strconv.Quote(path)+".")
				g.P("const ", constant, " = ", strconv.Quote(path))
				if string(root.Desc.Name()) == "SandboxInfo" && !contains(templates.FieldPaths.ExcludedOrigins, string(f.Desc.Name())) {
					response = append(response, pathInfo{constant: constant, fields: fields})
				}
				if f.Message != nil && !f.Desc.IsMap() && !nativePathLeaf(string(f.Message.Desc.FullName()), templates.FieldPaths) {
					walk(f.Message, fields, seen)
				}
			}
		}
		walk(root, nil, map[protoreflect.FullName]bool{})
		if string(root.Desc.Name()) == "Config" {
			declarationgen.WriteComment(g, templates.FieldPaths.Roots["Config"]+"Provider", "is the canonical path \"provider\".")
			g.P("const ", templates.FieldPaths.Roots["Config"], "Provider = \"provider\"")
		}
	}
	if len(response) > 0 {
		emitter := declarationgen.TemplateEmitter{G: g, Plugin: p, Templates: templates, Profile: profile}
		d, c, err := emitter.Begin("response_presence", declarationgen.TemplateBindings{Types: map[string]spec.TypeRef{"response": {Ref: "schema.SandboxInfo"}}})
		if err != nil {
			return err
		}
		d.Body(c, "if ", d.Param(c, "info"), "==nil{return false};switch ", d.Param(c, "path"), "{")
		for _, item := range response {
			d.Body(c, "case ", item.constant, ":")
			expression := d.Param(c, "info")
			valid := true
			for i, f := range item.fields {
				if f.Desc.IsList() && i < len(item.fields)-1 {
					valid = false
					break
				}
				expression += "." + naming.FieldName(f)
				if i < len(item.fields)-1 {
					d.Body(c, "if ", expression, "==nil{return false}")
				} else {
					if f.Desc.IsMap() || f.Desc.IsList() || f.Desc.HasPresence() {
						d.Body(c, "return ", expression, "!=nil")
					} else if f.Desc.Kind() == protoreflect.StringKind {
						d.Body(c, "return ", expression, "!=\"\"")
					} else {
						d.Body(c, "return true")
					}
				}
			}
			if !valid {
				d.Body(c, "return false")
			}
		}
		d.Body(c, "default:return false}}")
		d, c, err = emitter.Begin("response_origins", declarationgen.TemplateBindings{Types: map[string]spec.TypeRef{"response": {Ref: "schema.SandboxInfo"}}})
		if err != nil {
			return err
		}
		presence, err := templates.Instantiate("response_presence", nil)
		if err != nil {
			return err
		}
		d.Body(c, "for path,origin:=range ", d.Param(c, "info"), ".Origins{if !origin.Valid()||!", presence.Declaration.Name, "(", d.Param(c, "info"), ",path){return ", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindInvalidResponse"), Provider: d.Param(c, "info") + ".Provider", Operation: strconv.Quote("create"), Field: strconv.Quote("origins.") + "+path", Message: strconv.Quote("sandbox-kit: origin must describe a known, present response field")}), "}};return nil")
		g.P("}")
	}
	return nil
}

func contains(values []string, value string) bool {
	for _, candidate := range values {
		if candidate == value {
			return true
		}
	}
	return false
}
func nativePathLeaf(name string, rules spec.FieldPathTemplates) bool {
	if contains(rules.NativeLeaves, name) {
		return true
	}
	for _, namespace := range rules.NativeNamespaces {
		if strings.HasPrefix(name, namespace) {
			return true
		}
	}
	return false
}
