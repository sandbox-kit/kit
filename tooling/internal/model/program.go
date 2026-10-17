package model

import (
	"fmt"
	"strings"

	codegen "github.com/sandbox-kit/kit/tooling/internal/gen/codegen/v1"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Program is compiled once before language emission. Descriptors preserve
// source identity and presence; bindings retain target-language SDK symbols.
type Program struct {
	Behaviors  Behaviors
	Templates  Templates
	Schema     *Schema
	Generation spec.Generation
	Validation spec.Validation
	Client     spec.Client
	Contracts  spec.Contracts
	Lifecycle  Lifecycle
	Providers  map[string]Provider
	Clients    map[string]Client
}

// Lifecycle is compiled solely from annotated protobuf service methods. It is
// deliberately not a YAML specification: service descriptors own the public
// operation graph while emitters own native syntax.
type Lifecycle struct {
	Operations []LifecycleOperation
}

type LifecycleOperation struct {
	ID            string
	Capability    string
	Interface     string
	Method        string
	Options       string
	Outcome       string
	Timeout       bool
	TerminalState string
}

// Client resolves the local creation operation without constructor conventions.
type Client struct {
	Declaration       *codegen.ClientDeclaration
	Config            protoreflect.MessageDescriptor
	Service           protoreflect.ServiceDescriptor
	Create            protoreflect.MethodDescriptor
	Lifecycle         Lifecycle
	ServiceOperations []ServiceOperation
	Operations        map[string][]string
}

type ServiceOperation struct {
	ID, Interface, Method, Options, Outcome string
	Handle                                  bool
}

func CompileClient(file protoreflect.FileDescriptor) (Client, error) {
	if !proto.HasExtension(file.Options(), codegen.E_Client) {
		return Client{}, nil
	}
	d := proto.GetExtension(file.Options(), codegen.E_Client).(*codegen.ClientDeclaration)
	result := Client{Declaration: d, Config: file.Messages().ByName(protoreflect.Name(d.ConfigMessage))}
	service := file.Services().ByName(protoreflect.Name(d.CreationService))
	if result.Config == nil || service == nil {
		return Client{}, fmt.Errorf("%s: client config and declared SDK service are required", file.Path())
	}
	result.Service = service
	result.Create = service.Methods().ByName("Create")
	if result.Create == nil || result.Create.IsStreamingClient() || result.Create.IsStreamingServer() {
		return Client{}, fmt.Errorf("%s: SandboxService.Create must be local unary", file.Path())
	}
	lifecycle, operations, err := compileLifecycle(service)
	if err != nil {
		return Client{}, fmt.Errorf("%s: %w", file.Path(), err)
	}
	result.Lifecycle = lifecycle
	result.ServiceOperations = operations
	return result, nil
}

// compileLifecycle is descriptor-first: public operation topology comes from
// SandboxService method annotations, never a parallel YAML registry. Native
// interface names remain an emitter concern and are derived from the portable
// operation identity below.
func compileLifecycle(service protoreflect.ServiceDescriptor) (Lifecycle, []ServiceOperation, error) {
	result := Lifecycle{}
	var serviceOperations []ServiceOperation
	seen := map[string]bool{}
	for i := 0; i < service.Methods().Len(); i++ {
		method := service.Methods().Get(i)
		if !proto.HasExtension(method.Options(), codegen.E_Operation) {
			return result, nil, fmt.Errorf("service method %s is missing operation metadata", method.FullName())
		}
		declaration := proto.GetExtension(method.Options(), codegen.E_Operation).(*codegen.OperationDeclaration)
		if method.IsStreamingClient() || method.IsStreamingServer() || strings.TrimSpace(declaration.Id) == "" || seen[declaration.Id] {
			return result, nil, fmt.Errorf("invalid operation metadata on %s", method.FullName())
		}
		seen[declaration.Id] = true
		if !declaration.Handle {
			if declaration.Id != "create" && declaration.Id != "get" && declaration.Id != "attach" {
				return result, nil, fmt.Errorf("client operation %s is not supported", declaration.Id)
			}
			if strings.TrimSpace(declaration.AdapterInterface) == "" || strings.TrimSpace(declaration.AdapterMethod) == "" {
				return result, nil, fmt.Errorf("client operation %s has no adapter binding", declaration.Id)
			}
			serviceOperations = append(serviceOperations, ServiceOperation{ID: declaration.Id, Interface: declaration.AdapterInterface, Method: declaration.AdapterMethod, Outcome: declaration.Outcome})
			continue
		}
		if declaration.Capability == "" || declaration.Outcome == "" || strings.TrimSpace(declaration.AdapterInterface) == "" || strings.TrimSpace(declaration.AdapterMethod) == "" {
			return result, nil, fmt.Errorf("invalid handle operation %s", declaration.Id)
		}
		result.Operations = append(result.Operations, LifecycleOperation{
			ID: declaration.Id, Capability: declaration.Capability, Interface: declaration.AdapterInterface,
			Method: declaration.AdapterMethod, Options: declaration.OptionsMessage,
			Outcome: declaration.Outcome, Timeout: declaration.Timeout, TerminalState: declaration.TerminalState,
		})
		serviceOperations = append(serviceOperations, ServiceOperation{ID: declaration.Id, Interface: declaration.AdapterInterface, Method: declaration.AdapterMethod, Options: declaration.OptionsMessage, Outcome: declaration.Outcome, Handle: true})
	}
	if len(result.Operations) == 0 {
		return result, nil, fmt.Errorf("service has no handle operations")
	}
	return result, serviceOperations, nil
}

func Compile(root string, files []protoreflect.FileDescriptor) (*Program, error) {
	result := &Program{Schema: NewSchema(files), Providers: map[string]Provider{}, Clients: map[string]Client{}}
	var err error
	rawTemplates, err := spec.LoadTemplates(root)
	if err != nil {
		return nil, err
	}
	if result.Templates, err = CompileTemplates(rawTemplates); err != nil {
		return nil, err
	}
	rawBehaviors, err := spec.LoadBehaviors(root)
	if err != nil {
		return nil, err
	}
	if result.Behaviors, err = CompileBehaviors(rawBehaviors); err != nil {
		return nil, err
	}
	if err = result.Behaviors.ValidateTemplates(result.Templates); err != nil {
		return nil, err
	}

	if result.Generation, err = spec.LoadGeneration(root); err != nil {
		return nil, err
	}
	if result.Validation, err = spec.LoadValidation(root); err != nil {
		return nil, err
	}
	if result.Client, err = spec.LoadClient(root); err != nil {
		return nil, err
	}
	if result.Contracts, err = spec.LoadContracts(root); err != nil {
		return nil, err
	}
	if _, lookupErr := result.Schema.Message(result.Client.ConfigMessage); lookupErr == nil {
		if err = result.Schema.ValidateRepresentation(result.Client.Representation); err != nil {
			return nil, err
		}
	}
	// protoc can invoke the plugin separately for shared and client files.
	// Validate only specifications whose root exists in this descriptor graph.
	for _, item := range []struct {
		root  string
		rules spec.Validation
	}{{"CreateOptions", result.Validation}, {result.Client.ConfigMessage, result.Client.Validation}} {
		if _, lookupErr := result.Schema.Message(item.root); lookupErr == nil {
			if err = result.Schema.ValidateRules(item.rules); err != nil {
				return nil, err
			}
		}
	}
	for _, file := range files {
		if proto.HasExtension(file.Options(), codegen.E_Client) && proto.HasExtension(file.Options(), codegen.E_Provider) {
			return nil, fmt.Errorf("%s: declare client and provider separately", file.Path())
		}
		if proto.HasExtension(file.Options(), codegen.E_Client) {
			client, err := CompileClient(file)
			if err != nil {
				return nil, err
			}
			client.Operations = result.Generation.Operations
			for _, operation := range client.Lifecycle.Operations {
				if operation.Options != "" {
					if _, lookupErr := result.Schema.Message(operation.Options); lookupErr != nil {
						return nil, fmt.Errorf("lifecycle %s: %w", operation.ID, lookupErr)
					}
				}
			}
			result.Lifecycle = client.Lifecycle
			result.Clients[file.Path()] = client
		}
		if !proto.HasExtension(file.Options(), codegen.E_Provider) {
			continue
		}
		d := proto.GetExtension(file.Options(), codegen.E_Provider).(*codegen.ProviderDeclaration)
		if _, ok := result.Providers[d.ProviderName]; ok {
			return nil, fmt.Errorf("duplicate provider %s", d.ProviderName)
		}
		rules, err := spec.LoadProvider(root, d.ProviderName)
		if err != nil {
			return nil, err
		}
		compiled, err := CompileProvider(result.Schema, rules)
		if err != nil {
			return nil, fmt.Errorf("provider %s: %w", d.ProviderName, err)
		}
		result.Providers[d.ProviderName] = compiled
	}
	return result, nil
}

func compileServiceContracts(surface APIModule, client Client) (APIModule, error) {
	add := func(typ spec.APIType) error {
		if _, exists := surface.TypesByID[typ.ID]; exists {
			return fmt.Errorf("service contract %s duplicates runtime API declaration", typ.ID)
		}
		surface.Types = append(surface.Types, typ)
		surface.TypesByID[typ.ID] = typ
		return nil
	}
	contract := func(id, name, doc string, method spec.Callable) error {
		return add(spec.APIType{ID: id, Name: name, Doc: doc, Kind: "interface", Visibility: "public", Methods: []spec.Callable{method}})
	}
	ref := func(name string) spec.TypeRef { return spec.TypeRef{Ref: "schema." + name, Reference: true} }
	for _, operation := range client.ServiceOperations {
		if operation.Handle || operation.ID == "get" {
			continue
		}
		params := []spec.APISlot{{ID: "id", Name: "id", Type: spec.TypeRef{Builtin: "string"}}}
		resultName := "SandboxInfo"
		if operation.ID == "create" {
			params = []spec.APISlot{{ID: "request", Name: "request", Type: spec.TypeRef{Ref: "schema.CreateOptions", Reference: true, Nullable: true}}}
			resultName = "CreateResult"
		}
		if err := contract(operation.Interface, pascalIdentifier(operation.Interface), "implements the "+operation.ID+" operation declared by SandboxService.", spec.Callable{
			ID: operation.Method, Name: pascalIdentifier(operation.Method), Cancellable: true, Fallible: true,
			Params: params, Results: []spec.APISlot{{Type: ref(resultName)}},
		}); err != nil {
			return surface, err
		}
	}
	for _, operation := range client.Lifecycle.Operations {
		option := []spec.APISlot{{ID: "id", Name: "id", Type: spec.TypeRef{Builtin: "string"}}}
		if operation.Options != "" {
			option = append(option, spec.APISlot{ID: "options", Name: "options", Type: spec.TypeRef{Ref: "schema." + operation.Options, Reference: true, Nullable: true}})
		}
		resultName := "SandboxInfo"
		if operation.Outcome == "exit_result" {
			resultName = "ExitResult"
		}
		name := pascalIdentifier(operation.Interface)
		methodName := pascalIdentifier(operation.Method)
		if err := contract(operation.Interface, name, "implements the "+operation.ID+" sandbox operation declared by SandboxService.", spec.Callable{
			ID: operation.Method, Name: methodName, Cancellable: true, Fallible: true, Params: option, Results: []spec.APISlot{{Type: ref(resultName)}},
		}); err != nil {
			return surface, err
		}
	}
	return surface, nil
}

func pascalIdentifier(value string) string {
	parts := strings.Split(value, "_")
	for i := range parts {
		if parts[i] != "" {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}
