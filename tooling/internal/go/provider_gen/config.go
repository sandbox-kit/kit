package providergen

import (
	"fmt"
	"sort"
	"strconv"

	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateConfig composes typed mapping groups without embedding Go code in specs.
func GenerateConfig(plugin *protogen.Plugin, file *protogen.File, provider spec.Provider, templates model.Templates, profile spec.LanguageProfile) error {
	c := provider.Client
	if c == nil {
		return nil
	}
	if err := validateComposition(plugin, provider); err != nil {
		return err
	}
	binding := c.Target.Bindings["go"]
	if binding.Import == "" || binding.Name == "" {
		return fmt.Errorf("client target requires Go binding")
	}
	core := protogen.GoImportPath("github.com/sandbox-kit/kit/sdks/go/sandbox")
	ident := func(name string) protogen.GoIdent { return protogen.GoIdent{GoName: name, GoImportPath: core} }
	// Resolve schema names through the same naming emitter as public Go types.
	fieldName := func(message, name string) (string, error) {
		if message == "Config" && name == "provider" {
			return "Provider", nil
		}
		field, err := sharedField(plugin, message, name)
		if err != nil {
			return "", err
		}
		return naming.FieldName(field), nil
	}
	group := func(name string) (spec.Group, error) {
		for _, g := range provider.Groups {
			if g.Name == name && g.Direction == "request" {
				return g, nil
			}
		}
		return spec.Group{}, fmt.Errorf("unknown request mapping group %q", name)
	}
	g := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+".client.gen.go", file.GoImportPath)
	declarationgen.Banner(g, "source: "+string(file.Desc.Path()),
		"mapClientConfig turns shared client settings into the native SDK constructor options.",
		"native, err := mapClientConfig(config)")
	g.P("package ", file.GoPackageName)
	d, callable, err := (declarationgen.TemplateEmitter{G: g, Plugin: plugin, Templates: templates, Profile: profile}).Begin("client_mapping", declarationgen.TemplateBindings{Types: map[string]spec.TypeRef{"config": {Ref: "schema.Config"}, "target": {Ref: "native_target"}}, Externals: map[string]spec.ExternalType{"native_target": {Name: binding.Name, Import: binding.Import}}})
	if err != nil {
		return err
	}
	d.Body(callable, "var params ", protogen.GoIdent{GoName: binding.Name, GoImportPath: protogen.GoImportPath(binding.Import)})
	d.Body(callable, "if ", d.Param(callable, "config"), "==nil {", d.Param(callable, "config"), "=&", ident("Config"), "{}}")
	d.Body(callable, "if err:=", ident("ValidateConfig"), "(", d.Param(callable, "config"), ");err!=nil{return params,err}")
	keys := make([]string, 0, len(c.Rejected))
	for key := range c.Rejected {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		field, err := fieldName("Config", key)
		if err != nil {
			return err
		}
		path, err := schema.Literal(g, plugin, "Config", key, templates.FieldPaths)
		if err != nil {
			return err
		}
		d.Body(callable, "if ", d.Param(callable, "config"), ".", field, "!=nil{return params,", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindUnsupported"), Provider: strconv.Quote(provider.Provider), Operation: strconv.Quote("initialize"), Field: path, Message: strconv.Quote("sandbox-kit " + provider.Provider + ": " + c.Rejected[key])}), "}")
	}
	if c.Settings != "" {
		grp, err := group(c.Settings)
		if err != nil {
			return err
		}
		if !sameNative(grp.Target, c.Target) {
			return fmt.Errorf("settings mapping has incompatible types")
		}
		d.Body(callable, "params,remaining,err:=", c.Settings, "(", d.Param(callable, "config"), ");if err!=nil{return params,err}")
	} else {
		d.Body(callable, "remaining:=*", d.Param(callable, "config"))
	}
	owned := append(append([]string{}, c.Managed...), c.Retained...)
	owned = append(owned, "auth")
	if c.Scope != nil {
		owned = append(owned, "scope")
	}
	for _, key := range owned {
		field, err := fieldName("Config", key)
		if err != nil {
			return err
		}
		d.Body(callable, "remaining.", field, "=nil")
	}
	for _, key := range keys {
		field, _ := fieldName("Config", key)
		d.Body(callable, "remaining.", field, "=nil")
	}
	emitComponent := func(component spec.ClientComponent, parent string, message string) error {
		field, err := fieldName(message, component.Field)
		if err != nil {
			return err
		}
		grp, err := group(component.Group)
		if err != nil {
			return err
		}
		d.Body(callable, "mapped,rest,err:=", component.Group, "(", parent, ".", field, ");if err!=nil{return params,err}")
		for _, key := range component.Retained {
			n, err := fieldName(grp.Source.Name, key)
			if err != nil {
				return err
			}
			d.Body(callable, "rest.", n, "=nil")
		}
		d.Body(callable, "if err:=", ident("RejectUnmapped"), "(", strconv.Quote(provider.Provider+" client "+component.Field), ",&rest);err!=nil{return params,err}")
		if component.Destination != "" {
			dest, ok := binding.Fields[component.Destination]
			if !ok {
				return fmt.Errorf("missing destination binding %q", component.Destination)
			}
			d.Body(callable, "params.", dest, "=&mapped")
		} else {
			if !sameNative(grp.Target, c.Target) {
				return fmt.Errorf("component target type mismatch")
			}
			for _, mapping := range grp.Fields {
				member, ok := grp.Target.Bindings["go"].Fields[mapping.To]
				if !ok {
					return fmt.Errorf("missing component member binding")
				}
				d.Body(callable, "params.", member, "=mapped.", member)
			}
		}
		return nil
	}
	if c.Scope != nil {
		if c.Scope.Field != "scope" {
			return fmt.Errorf("client scope must select scope")
		}
		d.Body(callable, "if ", d.Param(callable, "config"), ".Scope!=nil{")
		if err := emitComponent(*c.Scope, d.Param(callable, "config"), "Config"); err != nil {
			return err
		}
		d.Body(callable, "}")
	}
	d.Body(callable, "if auth:=", d.Param(callable, "config"), ".GetAuth();auth!=nil{switch{")
	seen := map[string]bool{}
	for _, component := range c.Auth {
		if seen[component.Field] {
			return fmt.Errorf("duplicate auth variant %q", component.Field)
		}
		seen[component.Field] = true
		field, err := fieldName("AuthConfig", component.Field)
		if err != nil {
			return err
		}
		d.Body(callable, "case auth.", field, "!=nil:")
		if err := emitComponent(component, "auth", "AuthConfig"); err != nil {
			return err
		}
	}
	authPath, err := schema.Literal(g, plugin, "Config", "auth", templates.FieldPaths)
	if err != nil {
		return err
	}
	d.Body(callable, "default:return params,", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindUnsupported"), Provider: strconv.Quote(provider.Provider), Operation: strconv.Quote("initialize"), Field: authPath, Message: strconv.Quote("sandbox-kit " + provider.Provider + ": unsupported authentication mode")}), "}}")
	d.Body(callable, "if err:=", ident("RejectUnmapped"), "(", strconv.Quote(provider.Provider+" client"), ",&remaining);err!=nil{return params,err};return params,nil}")
	return nil
}
