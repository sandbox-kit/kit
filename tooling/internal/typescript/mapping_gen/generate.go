package mappinggen

import (
	"fmt"
	"go/token"
	"sort"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/typescript/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	validationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/validation_gen"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Generate emits portable provider field mappings from the compiled YAML model.
// Provider call ordering and source selection deliberately remain handwritten.
func Generate(plugin *protogen.Plugin, file *protogen.File, provider model.Provider) error {
	if err := validationgen.ValidateProvider(provider); err != nil {
		return err
	}
	g := declarationgen.GeneratedFile(plugin, "provider.mappings.gen.ts", string(file.Desc.Path())+" + specs/providers/"+provider.Rules.Provider+".yaml")

	shared, native := mappingImports(provider)
	g.P("import { durationMs } from \"@bufbuild/protobuf/wkt\";")
	g.P("import { ErrorKind, PolicyMode, ValueOrigin, sandboxError, type ", strings.Join(shared, ", type "), " } from \"@sandbox-kit/sandbox\";")
	for _, path := range sortedKeys(native) {
		g.P("import type { ", strings.Join(native[path], ", "), " } from ", strconv.Quote(path), ";")
	}
	g.P()
	g.P("type Mutable<T> = { -readonly [K in keyof T]: T[K] };")
	g.P()

	for _, group := range provider.Groups {
		if err := emitGroup(g, provider, group); err != nil {
			return err
		}
	}
	errorgen.EmitMappingError(g, provider.Rules.Provider)
	return nil
}

func emitGroup(g *protogen.GeneratedFile, provider model.Provider, group model.Group) error {
	rule := group.Rule
	request := rule.Direction == "request"
	if !request && rule.Direction != "response" {
		return fmt.Errorf("%s: invalid TypeScript mapping direction %q", rule.Name, rule.Direction)
	}
	source := typeName(rule.Source)
	target := typeName(rule.Target)
	result := target
	if !request {
		result = "Partial<" + target + ">"
	}
	operation := "create"
	if provider.Rules.Client != nil && rule.Target.Name == provider.Rules.Client.Target.Name {
		operation = "initialize"
	}
	g.P("export function ", rule.Name, "(source: ", source, " | undefined, operation: string = ", strconv.Quote(operation), "): ", result, " {")
	g.P("  const target: Mutable<", result, "> = {};")
	g.P("  if (source === undefined) return target;")
	if len(rule.When) > 0 {
		conditions := make([]string, 0, len(rule.When))
		for _, condition := range rule.When {
			member, err := nativeMember(rule.Source, condition.Field)
			if err != nil {
				return fmt.Errorf("%s: %w", rule.Name, err)
			}
			conditions = append(conditions, "(source."+member+" ?? 0) > "+strconv.Itoa(condition.GreaterThan))
		}
		g.P("  if (!(" + strings.Join(conditions, " && ") + ")) return target;")
	}
	if rule.Origin != "" {
		g.P("  const origins: Record<string, ValueOrigin> = {};")
		g.P("  target.origins = origins;")
	}
	for _, resolved := range group.Fields {
		mapping := resolved.Rule
		if request {
			if binding, ok := rule.Target.Bindings["typescript"]; ok {
				if representation, ok := binding.Representations[mapping.To]; ok {
					mapping.Transform = representation.Transform
					mapping.Factor = representation.Factor
					mapping.Maximum = representation.Maximum
					mapping.Cast = representation.Cast
				}
			}
		}
		sharedField := resolved.Shared.Fields[len(resolved.Shared.Fields)-1]
		sharedName := naming.Camel(string(sharedField.Name()))
		from, to := sharedName, mapping.To
		if request {
			member, err := nativeMember(rule.Target, mapping.To)
			if err != nil {
				return fmt.Errorf("%s: %w", rule.Name, err)
			}
			to = member
		} else {
			member, err := nativeMember(rule.Source, mapping.From)
			if err != nil {
				return fmt.Errorf("%s: %w", rule.Name, err)
			}
			from, to = member, sharedName
		}
		if mapping.Transform == "policy_minutes" {
			if err := emitPolicy(g, resolved.ErrorPath, from, to, mapping); err != nil {
				return fmt.Errorf("%s: %w", rule.Name, err)
			}
			continue
		}
		optional := sharedField.HasPresence() || sharedField.IsList() || sharedField.IsMap()
		if !request {
			optional = true // Native SDK response fields may be absent even when protobuf requires them.
		}
		g.P("  {")
		indent := "    "
		if optional {
			g.P("    if (source.", from, " !== undefined) {")
			indent = "      "
		}
		if isDurationField(sharedField) {
			g.P(indent, "const value = durationMs(source.", from, ");")
		} else {
			g.P(indent, "const value = source.", from, ";")
		}
		expression, validation, err := transformExpression("value", mapping.Transform, mapping.Factor, mapping.Maximum, mapping.Cast, sharedField)
		if err != nil {
			return fmt.Errorf("%s.%s: %w", rule.Name, mapping.From, err)
		}
		if validation != "" {
			kind := "ErrorKind.InvalidArgument"
			if !request {
				kind = "ErrorKind.InvalidResponse"
			}
			g.P(indent, "if (", validation, ") mappingError(", kind, ", ", strconv.Quote(resolved.ErrorPath), ", ", strconv.Quote(mapping.From+" cannot be represented"), ", operation);")
		}
		if mapping.Transform == "" && mapping.Cast == "" {
			if sharedField.IsMap() {
				expression = "{ ...value }"
			} else if sharedField.IsList() {
				expression = "[...value]"
			}
		}
		g.P(indent, "target.", to, " = ", expression, ";")
		if rule.Origin != "" {
			origin := "ValueOrigin.Provider"
			if rule.Origin == "request" {
				origin = "ValueOrigin.Request"
			}
			g.P(indent, "origins[", strconv.Quote(mapping.To), "] = ", origin, ";")
		}
		if optional {
			g.P("    }")
		}
		g.P("  }")
	}
	if !request && rule.Target.Name == "SandboxInfo" {
		g.P("  target.provider = ", strconv.Quote(provider.Rules.Provider), ";")
	}
	g.P("  return target;")
	g.P("}")
	g.P()
	return nil
}

func emitPolicy(g *protogen.GeneratedFile, field, from, to string, mapping spec.Mapping) error {
	if mapping.Policy == nil || (mapping.Policy.Disabled != "zero" && mapping.Policy.Disabled != "reject") || mapping.Factor != 0 || mapping.Cast != "" {
		return fmt.Errorf("invalid policy_minutes declaration")
	}
	g.P("  if (source.", from, " !== undefined) {")
	g.P("    const policy = source.", from, ";")
	g.P("    const mode = policy.mode ?? PolicyMode.Default;")
	g.P("    if (mode === PolicyMode.Disabled) {")
	if mapping.Policy.Disabled == "zero" {
		g.P("      target.", to, " = 0;")
	} else {
		g.P("      mappingError(ErrorKind.Unsupported, ", strconv.Quote(field), ", \"explicit disabling is not representable\", operation);")
	}
	g.P("    } else if (mode === PolicyMode.After) {")
	g.P("      if (policy.after === undefined) mappingError(ErrorKind.InvalidArgument, ", strconv.Quote(field), ", \"duration is required\", operation);")
	g.P("      const afterMs = durationMs(policy.after);")
	if !mapping.Policy.Immediate {
		g.P("      if (afterMs === 0) mappingError(ErrorKind.Unsupported, ", strconv.Quote(field), ", \"immediate action is not representable\", operation);")
	}
	g.P("      if (!Number.isFinite(afterMs) || afterMs < 0 || afterMs % 60_000 !== 0) mappingError(ErrorKind.InvalidArgument, ", strconv.Quote(field), ", \"delay cannot be represented in whole minutes\", operation);")
	g.P("      const minutes = afterMs / 60_000;")
	if mapping.Maximum != nil {
		g.P("      if (minutes > ", *mapping.Maximum, ") mappingError(ErrorKind.InvalidArgument, ", strconv.Quote(field), ", \"delay exceeds documented maximum\", operation);")
	}
	g.P("      target.", to, " = minutes;")
	g.P("    } else if (mode !== PolicyMode.Default) {")
	g.P("      mappingError(ErrorKind.InvalidArgument, ", strconv.Quote(field), ", \"unknown policy mode\", operation);")
	g.P("    }")
	g.P("  }")
	return nil
}

func transformExpression(value, transform string, factor uint64, maximum *uint64, cast string, field protoreflect.FieldDescriptor) (string, string, error) {
	expression, invalid := value, ""
	switch transform {
	case "":
	case "whole":
		invalid = "!Number.isInteger(" + value + ")"
	case "scaled_integer":
		if factor == 0 || field.Kind() != protoreflect.DoubleKind {
			return "", "", fmt.Errorf("scaled_integer requires a double and positive factor")
		}
		expression = "Math.round(" + value + " * " + strconv.FormatUint(factor, 10) + ")"
		invalid = "!Number.isFinite(" + value + ") || " + value + " < 0 || " + value + " !== " + expression + " / " + strconv.FormatUint(factor, 10)
	case "duration_seconds":
		expression = value + " / 1_000"
		invalid = value + " <= 0 || " + value + " % 1_000 !== 0"
	case "divide_exactly":
		if factor == 0 {
			return "", "", fmt.Errorf("divide_exactly requires a positive factor")
		}
		divisor := strconv.FormatUint(factor, 10) + "n"
		expression = value + " / " + divisor
		invalid = value + " % " + divisor + " !== 0n"
	case "multiply":
		if factor == 0 {
			return "", "", fmt.Errorf("multiply requires a positive factor")
		}
		expression = "BigInt(" + value + ") * " + strconv.FormatUint(factor, 10) + "n"
		invalid = "!Number.isSafeInteger(" + value + ") || " + value + " < 0"
	default:
		return "", "", fmt.Errorf("unknown transform %q", transform)
	}
	if maximum != nil {
		bound := expression
		suffix := ""
		if field.Kind() == protoreflect.Int64Kind || field.Kind() == protoreflect.Uint64Kind || transform == "multiply" || transform == "divide_exactly" {
			suffix = "n"
		}
		condition := bound + " > " + strconv.FormatUint(*maximum, 10) + suffix
		if invalid != "" {
			invalid += " || "
		}
		invalid += condition
	}
	switch cast {
	case "", "string":
		if cast == "string" {
			expression = "String(" + expression + ")"
		}
	case "integer", "number":
		expression = "Number(" + expression + ")"
	case "unsigned":
	default:
		return "", "", fmt.Errorf("unknown cast %q", cast)
	}
	return expression, invalid, nil
}

func isDurationField(field protoreflect.FieldDescriptor) bool {
	return field.Kind() == protoreflect.MessageKind && field.Message() != nil && field.Message().FullName() == "google.protobuf.Duration"
}

func mappingImports(provider model.Provider) ([]string, map[string][]string) {
	sharedSet := map[string]bool{}
	nativeSet := map[string]map[string]bool{}
	for _, group := range provider.Groups {
		collectTypeImport(sharedSet, nativeSet, group.Rule.Source)
		collectTypeImport(sharedSet, nativeSet, group.Rule.Target)
	}
	shared := sortedKeys(sharedSet)
	native := map[string][]string{}
	for path, names := range nativeSet {
		for name := range names {
			native[path] = append(native[path], name+" as Native"+name)
		}
		sort.Strings(native[path])
	}
	return shared, native
}

func collectTypeImport(shared map[string]bool, native map[string]map[string]bool, value spec.Type) {
	binding, ok := value.Bindings["typescript"]
	if !ok {
		shared[value.Name] = true
		return
	}
	if native[binding.Import] == nil {
		native[binding.Import] = map[string]bool{}
	}
	native[binding.Import][binding.Name] = true
}

func typeName(value spec.Type) string {
	if binding, ok := value.Bindings["typescript"]; ok {
		return "Native" + binding.Name
	}
	return value.Name
}

func nativeMember(value spec.Type, field string) (string, error) {
	binding, ok := value.Bindings["typescript"]
	if !ok {
		return naming.Camel(field), nil
	}
	member := binding.Fields[field]
	if !token.IsIdentifier(member) {
		return "", fmt.Errorf("%s: missing/invalid TypeScript SDK field binding %q", value.Name, field)
	}
	return member, nil
}

func sortedKeys[V any](values map[string]V) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
