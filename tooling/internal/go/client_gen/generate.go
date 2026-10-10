package clientgen

import (
	"fmt"
	"strconv"

	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generate applies Go constructor and naming conventions to a resolved client.
func Generate(plugin *protogen.Plugin, file *protogen.File, compiled model.Client, profile spec.LanguageProfile) error {
	if compiled.Declaration == nil {
		return nil
	}
	if err := spec.ValidateOperations(compiled.Operations); err != nil {
		return err
	}
	name := compiled.Surface.TypesByID["client"].Name
	var create *protogen.Method
	for _, service := range file.Services {
		for _, method := range service.Methods {
			if method.Desc.FullName() == compiled.Create.FullName() {
				create = method
			}
		}
	}
	var configMessage *protogen.Message
	for _, message := range file.Messages {
		if message.Desc.FullName() == compiled.Config.FullName() {
			configMessage = message
		}
	}
	if create == nil || configMessage == nil {
		return fmt.Errorf("%s: missing Go metadata for client", file.Desc.Path())
	}
	config := configMessage.GoIdent
	g := plugin.NewGeneratedFile(file.GeneratedFilenamePrefix+profile.Outputs["client"], file.GoImportPath)
	declarationgen.Banner(g, "source: "+string(file.Desc.Path()),
		"Client owns an initialized provider SDK. Create returns a sandbox handle. Close releases the SDK and leaves the cloud sandbox running.",
		"client, err := NewClient(Config{\n    Provider: provider,\n})\ninstance, err := client.Create(ctx, nil)\nerr = client.Close(ctx)")
	g.P("package ", file.GoPackageName)
	errors := func(name string) protogen.GoIdent { return protogen.GoIdent{GoName: name, GoImportPath: "errors"} }
	localError := func(kind, provider, operation, field, message string) string {
		return errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, kind), Provider: provider, Operation: strconv.Quote(operation), Field: strconv.Quote(field), Message: message})
	}
	d, err := declarationgen.New(g, plugin, compiled.Surface, profile, nil)
	if err != nil {
		return err
	}
	for _, id := range []string{"backend", "provider", "client"} {
		if err := d.TypeDeclaration(id); err != nil {
			return err
		}
	}
	binding, err := d.BeginBehavior("", "missing_binding", "missing_binding")
	if err != nil {
		return err
	}
	d.Body(binding, "value:=", protogen.GoIdent{GoName: "ValueOf", GoImportPath: "reflect"}, "(", d.Param(binding, "binding"), ");if !value.IsValid(){return true};switch value.Kind(){case reflect.Chan,reflect.Func,reflect.Interface,reflect.Map,reflect.Pointer,reflect.Slice:return value.IsNil()};return false")
	g.P("}")
	constructor, err := d.BeginBehavior("", "constructor", "initialize")
	if err != nil {
		return err
	}
	d.Body(constructor, "provider:=\"\";defer func(){", d.ErrorResult(constructor), "=WithErrorContext(", d.ErrorResult(constructor), ",provider,\"initialize\")}()")
	for _, step := range compiled.Operations["initialize"] {
		switch step {
		case "require_provider":
			d.Body(constructor, "if ", compiled.Surface.FunctionsByID["missing_binding"].Name, "(", d.Param(constructor, "config"), ".Provider){return nil,", localError("ErrorKindInvalidArgument", "", "initialize", "provider", strconv.Quote("sandbox-kit: provider is required")), "}")
		case "identify_provider":
			d.Body(constructor, "provider=", d.Param(constructor, "config"), ".Provider.", d.MethodName("provider", "name"), "()")
		case "validate_configuration":
			d.Body(constructor, "if err:=Validate", config.GoName, "(&", d.Param(constructor, "config"), ");err!=nil{return nil,err}")
		case "require_provider_name":
			d.Body(constructor, "expected:=provider;if ", protogen.GoIdent{GoName: "TrimSpace", GoImportPath: "strings"}, "(expected)==\"\"{return nil,", localError("ErrorKindInvalidArgument", "", "initialize", "provider", strconv.Quote("sandbox-kit: provider name is required")), "}")
		case "capture_timeout":
			d.Body(constructor, "var timeout *time.Duration;if ", d.Param(constructor, "config"), ".Timeout!=nil{value:=*", d.Param(constructor, "config"), ".Timeout;timeout=&value}")
		case "initialize_provider":
			d.Body(constructor, "backend,err:=", d.Param(constructor, "config"), ".Provider.", d.MethodName("provider", "initialize"), "(&", d.Param(constructor, "config"), ");if err!=nil{return nil,err}")
		case "require_backend":
			d.Body(constructor, "if ", compiled.Surface.FunctionsByID["missing_binding"].Name, "(backend){return nil,", localError("ErrorKindInvalidResponse", "expected", "initialize", "backend", strconv.Quote("sandbox-kit: provider returned no initialized backend")), "}")
		case "match_backend_identity":
			mismatch := g.QualifiedGoIdent(protogen.GoIdent{GoName: "Sprintf", GoImportPath: "fmt"}) + "(" + strconv.Quote("sandbox-kit: initialized backend %q differs from selected provider %q") + ",actual,expected)"
			d.Body(constructor, "if actual:=backend.", d.MethodName("backend", "name"), "();actual!=expected{return nil,", errors("Join"), "(", localError("ErrorKindInvalidResponse", "expected", "initialize", "provider", mismatch), ",backend.", d.MethodName("backend", "close"), "(context.Background()))}")
		case "return_client":
			d.Body(constructor, "return &", name, "{", d.Field("client", "backend"), ":backend,", d.Field("client", "timeout"), ":timeout},nil")
		}
	}
	g.P("}")
	accessor, err := d.BeginBehavior("client", "provider_name", "provider_name")
	if err != nil {
		return err
	}
	d.Body(accessor, "return ", d.Receiver(accessor), ".", d.Field("client", "backend"), ".", d.MethodName("backend", "name"), "()")
	g.P("}")
	closeMethod, err := d.BeginBehavior("client", "close", "close")
	if err != nil {
		return err
	}
	d.Body(closeMethod, "provider:=\"\";if "+d.Receiver(closeMethod)+"!=nil&&"+d.Receiver(closeMethod)+"."+d.Field("client", "backend")+"!=nil{provider="+d.Receiver(closeMethod)+"."+d.Field("client", "backend")+"."+d.MethodName("backend", "name")+"()};defer func(){"+d.ErrorResult(closeMethod)+"=WithErrorContext("+d.ErrorResult(closeMethod)+",provider,\"close\")}()")
	for _, step := range compiled.Operations["close"] {
		switch step {
		case "require_context":
			d.Body(closeMethod, "if ", d.Context(closeMethod), "==nil{return ", localError("ErrorKindInvalidArgument", "provider", "close", "context", strconv.Quote("sandbox-kit: context is required")), "}")
		case "close_backend":
			d.Body(closeMethod, "if ", d.Receiver(closeMethod), "==nil||", d.Receiver(closeMethod), ".", d.Field("client", "backend"), "==nil{return nil};return ", d.Receiver(closeMethod), ".", d.Field("client", "backend"), ".", d.MethodName("backend", "close"), "(", d.Context(closeMethod), ")")
		}
	}
	g.P("}")
	createMethod, err := d.BeginBehavior("client", "create", "create")
	if err != nil {
		return err
	}
	d.Body(createMethod, "provider:=\"\";if "+d.Receiver(createMethod)+"!=nil&&"+d.Receiver(createMethod)+"."+d.Field("client", "backend")+"!=nil{provider="+d.Receiver(createMethod)+"."+d.Field("client", "backend")+"."+d.MethodName("backend", "name")+"()};defer func(){"+d.ErrorResult(createMethod)+"=WithErrorContext("+d.ErrorResult(createMethod)+",provider,\"create\")}()")
	for _, step := range compiled.Operations["create"] {
		switch step {
		case "require_context":
			d.Body(createMethod, "if ", d.Context(createMethod), "==nil{return nil,", localError("ErrorKindInvalidArgument", "provider", "create", "context", strconv.Quote("sandbox-kit: context is required")), "};if err:=", d.Context(createMethod), ".Err();err!=nil{return nil,err}")
		case "require_client":
			d.Body(createMethod, "if ", d.Receiver(createMethod), "==nil||", d.Receiver(createMethod), ".", d.Field("client", "backend"), "==nil{return nil,", localError("ErrorKindInvalidArgument", "provider", "create", "client", strconv.Quote("sandbox-kit: client is not initialized")), "}")
		case "prepare_request":
			d.Body(createMethod, d.Param(createMethod, "request"), ",err:=prepareCreateRequest(", d.Param(createMethod, "request"), ");if err!=nil{return nil,err}")
		case "resolve_deadline":
			d.Body(createMethod, "timeout:=", d.Param(createMethod, "request"), ".GetProvisioning().GetTimeout();if timeout==nil{timeout=", d.Receiver(createMethod), ".", d.Field("client", "timeout"), "};if timeout!=nil&&*timeout>0{var cancel context.CancelFunc;", d.Context(createMethod), ",cancel=context.WithTimeout(", d.Context(createMethod), ",*timeout);defer cancel()}")
			// Materialize the resolved timeout for backend readiness mappings too.
			d.Body(createMethod, "if timeout!=nil{if ", d.Param(createMethod, "request"), ".Provisioning==nil{", d.Param(createMethod, "request"), ".Provisioning=&ProvisioningOptions{}};value:=*timeout;", d.Param(createMethod, "request"), ".Provisioning.Timeout=&value}")
		case "invoke_creation":
			d.Body(createMethod, "response,err:=", d.Receiver(createMethod), ".", d.Field("client", "backend"), ".", d.MethodName("backend", "create"), "(", d.Context(createMethod), ",", d.Param(createMethod, "request"), ");if err!=nil{return nil,err}")
		case "return_handle":
			d.Body(createMethod, "return sandboxFromResponse(", d.Receiver(createMethod), ".", d.Field("client", "backend"), ",response)")
		}
	}
	g.P("}")
	return nil
}
