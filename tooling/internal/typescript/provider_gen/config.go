package providergen

import (
	"fmt"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func uniqueSorted(values []string) []string {
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

// GenerateConfig emits portable client-configuration composition. Native SDK
// construction and provider-only prerequisites remain in the adapter.
func GenerateConfig(plugin *protogen.Plugin, file *protogen.File, provider spec.Provider, schema *model.Schema) error {
	client := provider.Client
	if client == nil {
		return nil
	}
	if err := validateConfigPlan(provider, schema); err != nil {
		return err
	}
	target, ok := client.Target.Bindings["typescript"]
	if !ok || strings.TrimSpace(target.Import) == "" || !token.IsIdentifier(target.Name) {
		return fmt.Errorf("provider %s client target requires a TypeScript binding", provider.Provider)
	}
	groups := map[string]spec.Group{}
	for _, group := range provider.Groups {
		if group.Direction == "request" {
			groups[group.Name] = group
		}
	}
	needed := []string{}
	if client.Settings != "" {
		needed = append(needed, client.Settings)
	}
	if client.Scope != nil {
		needed = append(needed, client.Scope.Group)
	}
	for _, auth := range client.Auth {
		needed = append(needed, auth.Group)
	}
	for _, name := range needed {
		if _, ok := groups[name]; !ok || !token.IsIdentifier(name) {
			return fmt.Errorf("provider %s has invalid client mapping group %q", provider.Provider, name)
		}
	}

	g := declarationgen.GeneratedFile(plugin, "provider.config.gen.ts", string(file.Desc.Path())+" + specs/providers/"+provider.Provider+".yaml")
	g.P("import type { ", target.Name, " as NativeClientConfig } from ", strconv.Quote(target.Import), ";")
	runtime := provider.Runtime["typescript"]
	sharedTypes := []string{"Config"}
	for _, state := range runtime.State {
		if state.Type != "string" {
			sharedTypes = append(sharedTypes, state.Type)
		}
	}
	g.P(`import { ErrorKind, sandboxError, validateConfig, type `, strings.Join(uniqueSorted(sharedTypes), ", type "), ` } from "@sandbox-kit/sandbox";`)
	if len(needed) > 0 {
		sort.Strings(needed)
		g.P("import { ", strings.Join(uniqueSorted(needed), ", "), ` } from "./provider.mappings.gen.js";`)
	}
	g.P()
	if len(client.Environment) > 0 {
		g.P(`export type Environment = Readonly<Record<string, string | undefined>>;`)
		g.P()
	}
	params := "config: Config"
	if len(client.Environment) > 0 {
		params += ", environment: Environment = environmentVariables()"
	}
	g.P("export function mapGeneratedClientConfig(", params, "): NativeClientConfig {")
	g.P("  let source = config;")
	for _, resolution := range client.Environment {
		field := naming.Camel(resolution.Field)
		candidates := make([]string, 0, len(resolution.Variables)+1)
		for _, variable := range resolution.Variables {
			candidates = append(candidates, "environment."+variable)
		}
		candidates = append(candidates, strconv.Quote(resolution.Default))
		g.P("  if (source.", field, " === undefined) source = { ...source, ", field, ": ", strings.Join(candidates, " ?? "), " };")
	}
	g.P("  validateConfig(source, ", strconv.Quote(provider.Provider), `, "initialize");`)
	if client.Settings != "" {
		g.P("  const target: NativeClientConfig = ", client.Settings, "(source, \"initialize\");")
	} else {
		g.P("  const target: NativeClientConfig = {};")
	}
	keys := make([]string, 0, len(client.Rejected))
	for key := range client.Rejected {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		g.P("  if (source.", naming.Camel(key), " !== undefined) unsupported(", strconv.Quote(naming.Camel(key)), ", ", strconv.Quote(client.Rejected[key]), ");")
	}
	if client.Scope != nil {
		if err := emitConfigComponent(g, *client.Scope, "source.scope", "scope", groups, client.Target, schema); err != nil {
			return err
		}
	}
	if len(client.Auth) > 0 {
		for index, auth := range client.Auth {
			branch := "if"
			if index > 0 {
				branch = "else if"
			}
			field := naming.Camel(auth.Field)
			g.P("  ", branch, " (source.auth?.credentials === ", strconv.Quote(field), ") {")
			if err := emitConfigComponentBody(g, auth, "source.auth."+field, "auth."+field, groups, client.Target, schema, "    "); err != nil {
				return err
			}
			g.P("  }")
		}
		g.P(`  else if (source.auth?.credentials !== undefined) unsupported("auth", "unsupported authentication mode");`)
	}
	g.P("  return target;")
	g.P("}")
	if len(runtime.State) > 0 {
		g.P()
		g.P("export interface ProviderState {")
		for _, state := range runtime.State {
			typeName := state.Type
			if typeName == "string" {
				typeName = "string"
			}
			optional := ""
			if state.Optional {
				optional = "?"
			}
			g.P("  readonly ", naming.Camel(state.Name), optional, ": ", typeName, ";")
		}
		g.P("}")
		g.P()
		g.P("export function captureProviderState(config: Config): ProviderState {")
		g.P("  return {")
		for _, state := range runtime.State {
			name := naming.Camel(state.Name)
			capture := naming.Camel(state.Capture)
			value := "config." + capture
			if state.Type != "string" {
				value = "structuredClone(" + value + ")"
			}
			if state.Optional {
				g.P("    ...(config.", capture, " === undefined ? {} : { ", name, ": ", value, " }),")
			} else {
				g.P("    ", name, ": ", value, ",")
			}
		}
		g.P("  };")
		g.P("}")
	}
	if len(client.Environment) > 0 {
		g.P()
		g.P(`function environmentVariables(): Environment { return (globalThis as { readonly process?: { readonly env?: Environment } }).process?.env ?? {}; }`)
	}
	g.P()
	g.P("function unsupported(field: string, message: string): never { throw sandboxError({ kind: ErrorKind.Unsupported, provider: ", strconv.Quote(provider.Provider), ", operation: \"initialize\", field, message: ", strconv.Quote("sandbox-kit "+provider.Provider+": "), " + message }); }")
	return nil
}

func validateConfigPlan(provider spec.Provider, schema *model.Schema) error {
	runtime := provider.Runtime["typescript"]
	captures := make([]string, 0, len(runtime.State))
	for _, state := range runtime.State {
		if !token.IsIdentifier(state.Name) || state.Capture == "" {
			return fmt.Errorf("provider %s has invalid TypeScript state field %q", provider.Provider, state.Name)
		}
		resolved, err := schema.Resolve("Config", state.Capture)
		if err != nil || len(resolved.Fields) != 1 {
			return fmt.Errorf("provider %s state %s must capture one Config field", provider.Provider, state.Name)
		}
		field := resolved.Fields[0]
		expected := ""
		if field.Kind() == protoreflect.StringKind {
			expected = "string"
		} else if field.Message() != nil {
			expected = string(field.Message().Name())
		}
		if expected == "" || state.Type != expected || state.Optional != field.HasPresence() {
			return fmt.Errorf("provider %s state %s disagrees with Config.%s type or presence", provider.Provider, state.Name, state.Capture)
		}
		captures = append(captures, state.Capture)
	}
	captured := func(path string) bool {
		for _, capture := range captures {
			if path == capture || strings.HasPrefix(path, capture+".") {
				return true
			}
		}
		return false
	}
	if provider.Client != nil {
		for _, retained := range provider.Client.Retained {
			if !captured(retained) {
				return fmt.Errorf("provider %s retained field %s has no TypeScript state capture", provider.Provider, retained)
			}
		}
		if provider.Client.Scope != nil {
			for _, retained := range provider.Client.Scope.Retained {
				path := provider.Client.Scope.Field + "." + retained
				if !captured(path) {
					return fmt.Errorf("provider %s retained field %s has no TypeScript state capture", provider.Provider, path)
				}
			}
		}
	}
	return nil
}

func emitConfigComponent(g *protogen.GeneratedFile, component spec.ClientComponent, source, path string, groups map[string]spec.Group, target spec.Type, schema *model.Schema) error {
	g.P("  if (", source, " !== undefined) {")
	if err := emitConfigComponentBody(g, component, source, path, groups, target, schema, "    "); err != nil {
		return err
	}
	g.P("  }")
	return nil
}

func emitConfigComponentBody(g *protogen.GeneratedFile, component spec.ClientComponent, source, path string, groups map[string]spec.Group, clientTarget spec.Type, schema *model.Schema, indent string) error {
	group := groups[component.Group]
	binding := group.Target.Bindings["typescript"]
	clientBinding := clientTarget.Bindings["typescript"]
	if binding.Import == "" || binding.Name == "" {
		return fmt.Errorf("client component %s has no TypeScript target", component.Group)
	}
	g.P(indent, "const mapped = ", component.Group, "(", source, ", \"initialize\");")
	if binding.Import == clientBinding.Import && binding.Name == clientBinding.Name {
		g.P(indent, "Object.assign(target, mapped);")
	} else {
		member := clientBinding.Fields[component.Destination]
		object := clientBinding.Objects[component.Destination]
		if !token.IsIdentifier(member) || object.Import != binding.Import || object.Name != binding.Name {
			return fmt.Errorf("client component %s has incompatible destination %q", component.Group, component.Destination)
		}
		g.P(indent, "target.", member, " = mapped;")
	}
	message, err := schema.Message(group.Source.Name)
	if err != nil {
		return err
	}
	owned := map[string]bool{}
	for _, mapping := range group.Fields {
		owned[mapping.From] = true
	}
	for _, retained := range component.Retained {
		owned[retained] = true
	}
	for index := 0; index < message.Fields().Len(); index++ {
		field := message.Fields().Get(index)
		if owned[string(field.Name())] {
			continue
		}
		member := naming.Camel(string(field.Name()))
		g.P(indent, "if (", source, ".", member, " !== undefined) unsupported(", strconv.Quote(path+"."+member), ", \"unsupported client configuration field\");")
	}
	return nil
}
