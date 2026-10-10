package model

import (
	"fmt"

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
	API        API
	Schema     *Schema
	Generation spec.Generation
	Validation spec.Validation
	Client     spec.Client
	Contracts  spec.Contracts
	Providers  map[string]Provider
	Clients    map[string]Client
}

// Client resolves the local creation operation without constructor conventions.
type Client struct {
	Declaration *codegen.ClientDeclaration
	Config      protoreflect.MessageDescriptor
	Create      protoreflect.MethodDescriptor
	Operations  map[string][]string
	Surface     APIModule
}

func CompileClient(file protoreflect.FileDescriptor) (Client, error) {
	if !proto.HasExtension(file.Options(), codegen.E_Client) {
		return Client{}, nil
	}
	d := proto.GetExtension(file.Options(), codegen.E_Client).(*codegen.ClientDeclaration)
	result := Client{Declaration: d, Config: file.Messages().ByName(protoreflect.Name(d.ConfigMessage))}
	service := file.Services().ByName(protoreflect.Name(d.CreationService))
	if result.Config == nil || service == nil || service.Methods().Len() != 1 {
		return Client{}, fmt.Errorf("%s: client config and single creation method are required", file.Path())
	}
	result.Create = service.Methods().Get(0)
	if result.Create.IsStreamingClient() || result.Create.IsStreamingServer() {
		return Client{}, fmt.Errorf("%s: creation method must be local unary", file.Path())
	}
	return result, nil
}

func Compile(root string, files []protoreflect.FileDescriptor) (*Program, error) {
	result := &Program{Schema: NewSchema(files), Providers: map[string]Provider{}, Clients: map[string]Client{}}
	var err error
	rawAPI, err := spec.LoadAPI(root)
	if err != nil {
		return nil, err
	}
	if result.API, err = CompileAPI(rawAPI); err != nil {
		return nil, err
	}
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
	if err = result.Behaviors.ValidateAPI(result.API); err != nil {
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
			client.Surface = result.API.Modules["client"]
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
