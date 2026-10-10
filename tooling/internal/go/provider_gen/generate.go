package providergen

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	rules "github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

// Generate emits provider declarations from YAML and native SDK orchestration.
func Generate(plugin *protogen.Plugin, file *protogen.File, runtime rules.RuntimeBinding, surface model.APIModule, profile rules.LanguageProfile) error {
	if !proto.HasExtension(file.Desc.Options(), codegen.E_Provider) {
		return nil
	}
	if proto.HasExtension(file.Desc.Options(), codegen.E_Client) {
		return fmt.Errorf("%s: declare client and provider separately", file.Desc.Path())
	}
	declaration := proto.GetExtension(file.Desc.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
	sdk := declaration.Sdks["go"]
	if strings.TrimSpace(declaration.ProviderName) == "" || sdk == nil || sdk.ImportPath == "" || !token.IsIdentifier(sdk.ClientType) {
		return fmt.Errorf("%s: provider and Go SDK binding are required", file.Desc.Path())
	}
	// Provider state is a declaration source, not duplicated API storage.
	types := append([]rules.APIType{}, surface.Types...)
	surface.TypesByID = map[string]rules.APIType{}
	for i, typ := range types {
		if typ.FieldsFrom == "provider_state" {
			typ.Fields = append([]rules.APISlot{}, typ.Fields...)
			seen := map[string]bool{}
			for _, field := range typ.Fields {
				seen[field.Name] = true
			}
			for _, field := range runtime.State {
				if !token.IsIdentifier(field.Name) || seen[field.Name] {
					return fmt.Errorf("invalid or duplicate backend state field %q", field.Name)
				}
				seen[field.Name] = true
				ref := rules.TypeRef{Ref: "schema." + field.Type, Optional: field.Optional}
				if field.Type == "string" {
					ref = rules.TypeRef{Builtin: "string", Optional: field.Optional}
				}
				typ.Fields = append(typ.Fields, rules.APISlot{ID: field.Name, Name: field.Name, Type: ref, Doc: field.Doc})
			}
		}
		types[i] = typ
		surface.TypesByID[typ.ID] = typ
	}
	surface.Types = types
	native := protogen.GoIdent{GoName: sdk.ClientType, GoImportPath: protogen.GoImportPath(sdk.ImportPath)}
	g := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+profile.Outputs["provider"], file.GoImportPath)
	declarationgen.Banner(g, "source: "+string(file.Desc.Path()),
		"New returns a provider factory. Pass it to sandbox.NewClient. The public client owns SDK cleanup.",
		"client, err := sandbox.NewClient(sandbox.Config{Provider: New()})")
	g.P("package ", file.GoPackageName)
	d, err := declarationgen.New(g, plugin, surface, profile, map[string]protogen.GoIdent{"native_client": native})
	if err != nil {
		return err
	}
	for _, typ := range surface.Types {
		if err := d.TypeDeclaration(typ.ID); err != nil {
			return err
		}
	}
	provider := surface.TypesByID["provider"].Name
	backend := surface.TypesByID["backend"].Name
	selectProvider, err := d.BeginBehavior("", "constructor", "select_provider")
	if err != nil {
		return err
	}
	d.Body(selectProvider, "return &", provider, "{}")
	g.P("}")
	for _, owner := range []string{"provider", "backend"} {
		name, err := d.BeginBehavior(owner, "name", "provider_name")
		if err != nil {
			return err
		}
		d.Body(name, "return ", strconv.Quote(declaration.ProviderName))
		g.P("}")
	}
	initialize, err := d.BeginBehavior("provider", "initialize", "initialize_provider")
	if err != nil {
		return err
	}
	d.Body(initialize, "backend,err:=", surface.FunctionsByID["initialize_backend"].Name, "(", d.Param(initialize, "config"), ");return backend,mapProviderError(err,\"initialize\")")
	g.P("}")
	if runtime.Constructor.Function != "" {
		if !token.IsIdentifier(runtime.Constructor.Function) {
			return fmt.Errorf("invalid SDK constructor %q", runtime.Constructor.Function)
		}
		initializeBackend, err := d.BeginBehavior("", "initialize_backend", "initialize_backend")
		if err != nil {
			return err
		}
		d.Body(initializeBackend, "if ", d.Param(initializeBackend, "config"), "==nil{", d.Param(initializeBackend, "config"), "=&", protogen.GoIdent{GoName: "Config", GoImportPath: "github.com/sandbox-kit/kit/sdks/go/sandbox"}, "{}}")
		d.Body(initializeBackend, "params,err:=mapClientConfig(", d.Param(initializeBackend, "config"), ");if err!=nil{return nil,err}")
		d.Body(initializeBackend, "native,err:=", protogen.GoIdent{GoName: runtime.Constructor.Function, GoImportPath: native.GoImportPath}, "(&params);if err!=nil{return nil,err}")
		d.Body(initializeBackend, "result:=&", backend, "{", d.Field("backend", "client"), ":native}")
		d.Body(initializeBackend, surface.FunctionsByID["capture_state"].Name, "(result,", d.Param(initializeBackend, "config"), ");return result,nil")
		g.P("}")
		captureState, err := d.BeginBehavior("", "capture_state", "capture_state")
		if err != nil {
			return err
		}
		d.Body(captureState, "if ", d.Param(captureState, "config"), "==nil{return}")
		params := map[string]string{}
		for _, param := range captureState.Params {
			params[param.ID] = param.Name
		}
		for _, field := range runtime.State {
			if field.Capture != "" {
				if err := capture(g, plugin, field.Capture, field.Type, field.Optional, params["result"]+"."+field.Name, params["config"]); err != nil {
					return err
				}
			}
		}
		g.P("}")
	}
	create, err := d.BeginBehavior("backend", "create", "provider_create")
	if err != nil {
		return err
	}
	d.Body(create, "response,err:=", d.Receiver(create), ".create(", d.Context(create), ",", d.Param(create, "request"), ");return response,mapProviderError(err,\"create\")")
	g.P("}")
	if !token.IsIdentifier(runtime.Cleanup.Method) {
		return fmt.Errorf("%s: SDK cleanup binding is required", file.Desc.Path())
	}
	closeMethod, err := d.BeginBehavior("backend", "close", "provider_close")
	if err != nil {
		return err
	}
	args := ""
	if runtime.Cleanup.AcceptsContext {
		args = "ctx"
	}
	if runtime.Cleanup.ReturnsError {
		d.Body(closeMethod, "return mapProviderError(", d.Receiver(closeMethod), ".", d.Field("backend", "client"), ".", runtime.Cleanup.Method, "(", args, "),\"close\")")
	} else {
		d.Body(closeMethod, d.Receiver(closeMethod), ".", d.Field("backend", "client"), ".", runtime.Cleanup.Method, "(", args, ");return nil")
	}
	g.P("}")
	g.P("var _ ", protogen.GoIdent{GoName: "Provider", GoImportPath: "github.com/sandbox-kit/kit/sdks/go/sandbox"}, "=(*", provider, ")(nil)")
	g.P("var _ ", protogen.GoIdent{GoName: "Backend", GoImportPath: "github.com/sandbox-kit/kit/sdks/go/sandbox"}, "=(*", backend, ")(nil)")
	return nil
}
