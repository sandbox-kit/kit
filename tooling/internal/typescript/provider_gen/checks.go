package providergen

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GenerateChecks emits typed, provider-specific validation from the portable
// checks vocabulary. Quotas, permissions, and remote capacity remain native.
func GenerateChecks(plugin *protogen.Plugin, file *protogen.File, provider spec.Provider, schema *model.Schema) error {
	g := declarationgen.GeneratedFile(plugin, "provider.checks.gen.ts", string(file.Desc.Path())+" + specs/providers/"+provider.Provider+".yaml")
	g.P(`import { durationMs } from "@bufbuild/protobuf/wkt";`)
	g.P(`import { ErrorKind, PolicyMode, sandboxError, type AutomaticAction, type CreateOptions } from "@sandbox-kit/sandbox";`)
	g.P()
	g.P("export function validateProviderCreate(request: CreateOptions): void {")
	for index, check := range provider.Checks {
		if len(check.ExclusivePolicies) > 0 || check.ForbidPolicyWith != "" {
			if err := emitPolicyCheck(g, check, provider.Provider); err != nil {
				return err
			}
			continue
		}
		resolved, err := schema.Resolve("CreateOptions", check.Path)
		if err != nil {
			return fmt.Errorf("provider %s check %d: %w", provider.Provider, index, err)
		}
		if err := emitPathCheck(g, check, resolved, provider.Provider, index); err != nil {
			return fmt.Errorf("provider %s check %s: %w", provider.Provider, check.Path, err)
		}
	}
	g.P("}")
	g.P()
	emitCheckHelpers(g, provider.Provider)
	return nil
}

func emitPolicyCheck(g *protogen.GeneratedFile, check spec.Check, provider string) error {
	if len(check.ExclusivePolicies) > 0 {
		if len(check.ExclusivePolicies) != 2 {
			return fmt.Errorf("exclusive policies requires exactly two fields")
		}
		a, b := naming.Camel(check.ExclusivePolicies[0]), naming.Camel(check.ExclusivePolicies[1])
		g.P("  if (delayed(request.lifetime?.", a, ") && delayed(request.lifetime?.", b, ")) invalid(\"lifetime\", \"incompatible lifetime policies\");")
		return nil
	}
	policy := naming.Camel(check.Path)
	switch check.ForbidPolicyWith {
	case "ephemeral":
		g.P("  if (delayed(request.lifetime?.", policy, ") && request.lifetime?.ephemeral === true) invalid(\"lifetime\", \"incompatible lifetime policies\");")
	case "stopped_delete":
		g.P("  if (delayed(request.lifetime?.", policy, ") && immediate(request.lifetime?.stoppedDelete)) invalid(\"lifetime\", \"incompatible lifetime policies\");")
	default:
		return fmt.Errorf("unsupported policy relationship %q", check.ForbidPolicyWith)
	}
	_ = provider
	return nil
}

func emitPathCheck(g *protogen.GeneratedFile, check spec.Check, path model.FieldPath, provider string, index int) error {
	if len(path.Fields) == 0 {
		return fmt.Errorf("empty field path")
	}
	expression := "request"
	closures := 0
	indent := "  "
	for fieldIndex, field := range path.Fields {
		member := naming.Camel(string(field.Name()))
		last := fieldIndex == len(path.Fields)-1
		next := expression + "." + member
		if !last {
			if field.IsList() {
				item := fmt.Sprintf("item%d_%d", index, fieldIndex)
				g.P(indent, "for (const ", item, " of ", next, " ?? []) {")
				expression = item
				closures++
				indent += "  "
			} else {
				value := fmt.Sprintf("value%d_%d", index, fieldIndex)
				g.P(indent, "const ", value, " = ", next, ";")
				g.P(indent, "if (", value, " !== undefined) {")
				expression = value
				closures++
				indent += "  "
			}
			continue
		}
		value := fmt.Sprintf("value%d", index)
		switch {
		case field.IsMap():
			if check.Format != "env_name" {
				return fmt.Errorf("map checks require env_name format")
			}
			g.P(indent, "for (const ", value, " of Object.keys(", next, " ?? {})) {")
			closures++
			indent += "  "
		case field.IsList():
			if field.Kind() != protoreflect.StringKind {
				return fmt.Errorf("collection checks require strings")
			}
			g.P(indent, "for (const ", value, " of ", next, " ?? []) {")
			closures++
			indent += "  "
		case field.HasPresence():
			g.P(indent, "const ", value, " = ", next, ";")
			g.P(indent, "if (", value, " !== undefined) {")
			closures++
			indent += "  "
		default:
			g.P(indent, "const ", value, " = ", next, ";")
		}
		condition, kind, err := checkCondition(check, field, value)
		if err != nil {
			return err
		}
		fieldPath := typescriptPath(path.Fields)
		g.P(indent, "if (", condition, ") fail(ErrorKind.", kind, ", ", strconv.Quote(fieldPath), ", ", strconv.Quote(check.Path+" is outside supported values"), ");")
	}
	for range closures {
		indent = strings.TrimSuffix(indent, "  ")
		g.P(indent, "}")
	}
	_ = provider
	return nil
}

func checkCondition(check spec.Check, field protoreflect.FieldDescriptor, value string) (string, string, error) {
	conditions := []string{}
	kind := "InvalidArgument"
	if len(check.Allowed) > 0 {
		if field.Kind() != protoreflect.StringKind || check.Minimum != nil || check.Maximum != nil || check.Format != "" {
			return "", "", fmt.Errorf("allowed values require an otherwise unconstrained string")
		}
		for _, allowed := range check.Allowed {
			conditions = append(conditions, value+" !== "+strconv.Quote(allowed))
		}
		return strings.Join(conditions, " && "), "Unsupported", nil
	}
	for _, bound := range []struct {
		value *float64
		op    string
	}{{check.Minimum, "<"}, {check.Maximum, ">"}} {
		if bound.value == nil {
			continue
		}
		if math.IsNaN(*bound.value) || math.IsInf(*bound.value, 0) {
			return "", "", fmt.Errorf("bounds must be finite")
		}
		switch field.Kind() {
		case protoreflect.DoubleKind, protoreflect.FloatKind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Int32Kind, protoreflect.Int64Kind:
		default:
			return "", "", fmt.Errorf("numeric bounds require a numeric scalar")
		}
		literal := strconv.FormatFloat(*bound.value, 'g', -1, 64)
		if field.Kind() == protoreflect.Uint64Kind || field.Kind() == protoreflect.Int64Kind {
			literal = "BigInt(" + strconv.Quote(literal) + ")"
		}
		conditions = append(conditions, value+" "+bound.op+" "+literal)
	}
	if check.Format != "" {
		if field.Kind() != protoreflect.StringKind && !field.IsMap() {
			return "", "", fmt.Errorf("format checks require strings")
		}
		format := "valid" + naming.Pascal(check.Format) + "(" + value + ")"
		condition := "!" + format
		if check.AllowEmpty {
			condition = value + ` !== "" && ` + condition
		}
		conditions = append(conditions, condition)
	}
	if len(conditions) == 0 {
		return "", "", fmt.Errorf("empty check")
	}
	return strings.Join(conditions, " || "), kind, nil
}

func typescriptPath(fields []protoreflect.FieldDescriptor) string {
	parts := make([]string, 0, len(fields))
	for _, field := range fields {
		parts = append(parts, naming.Camel(string(field.Name())))
	}
	return strings.Join(parts, ".")
}

func emitCheckHelpers(g *protogen.GeneratedFile, provider string) {
	g.P(`function delayed(policy: AutomaticAction | undefined): boolean { return (policy?.mode ?? PolicyMode.Default) === PolicyMode.After && policy?.after !== undefined && durationMs(policy.after) > 0; }`)
	g.P(`function immediate(policy: AutomaticAction | undefined): boolean { return (policy?.mode ?? PolicyMode.Default) === PolicyMode.After && policy?.after !== undefined && durationMs(policy.after) === 0; }`)
	g.P(`function validEnvName(value: string): boolean { return /^[A-Za-z_][A-Za-z0-9_]*$/.test(value); }`)
	g.P(`function validAbsolutePath(value: string): boolean { return value.startsWith("/"); }`)
	g.P(`function validHttpUrl(value: string): boolean { try { const parsed = new URL(value); return parsed.hostname !== "" && (parsed.protocol === "http:" || parsed.protocol === "https:"); } catch { return false; } }`)
	g.P(`function validDomain(value: string): boolean { const domain = value.startsWith("*.") ? value.slice(2) : value; return domain.length <= 253 && /^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$/i.test(domain); }`)
	g.P(`function validCidr(value: string): boolean { const separator = value.lastIndexOf("/"); if (separator <= 0) return false; const address = value.slice(0, separator); const prefix = Number(value.slice(separator + 1)); if (!Number.isInteger(prefix)) return false; return address.includes(":") ? prefix >= 0 && prefix <= 128 && validIpv6(address) : prefix >= 0 && prefix <= 32 && validIpv4(address); }`)
	g.P(`function validIpv4(value: string): boolean { const parts = value.split("."); return parts.length === 4 && parts.every((part) => /^(0|[1-9][0-9]{0,2})$/.test(part) && Number(part) <= 255); }`)
	g.P(`function validIpv6(value: string): boolean { if ((value.match(/::/g)?.length ?? 0) > 1) return false; const halves = value.split("::"); if (halves.length > 2) return false; const count = (part: string): number | undefined => { if (part === "") return 0; const groups = part.split(":"); let total = 0; for (const group of groups) { if (/^[0-9a-f]{1,4}$/i.test(group)) total++; else if (validIpv4(group)) total += 2; else return undefined; } return total; }; const left = count(halves[0] ?? ""); const right = count(halves[1] ?? ""); return left !== undefined && right !== undefined && (halves.length === 2 ? left + right < 8 : left === 8); }`)
	g.P("function invalid(field: string, message: string): never { fail(ErrorKind.InvalidArgument, field, message); }")
	g.P("function fail(kind: ErrorKind, field: string, message: string): never { throw sandboxError({ kind, provider: ", strconv.Quote(provider), ", operation: \"create\", field, message: ", strconv.Quote("sandbox-kit "+provider+": "), " + message }); }")
}
