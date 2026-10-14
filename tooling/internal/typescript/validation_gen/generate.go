// Package validationgen emits TypeScript validators from portable rules.
package validationgen

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Emit writes validators for the shared request graph. The surrounding client
// runtime supplies sandboxError and ErrorKind so invalid requests have the
// same stable error shape as provider failures.
func Emit(g *protogen.GeneratedFile, schema *model.Schema, rules spec.Validation, root string) error {
	if len(rules.Messages) == 0 {
		return nil
	}
	names := make([]string, 0, len(rules.Messages))
	for name := range rules.Messages {
		if _, err := schema.Message(name); err != nil {
			return err
		}
		names = append(names, name)
	}
	sort.Strings(names)
	g.P(`function validationError(provider: string, operation: string, field: string, message: string): never { throw sandboxError({ kind: ErrorKind.InvalidArgument, provider, operation, field, message: "sandbox-kit: " + message }); }`)
	g.P(`function joinValidationPath(parent: string, field: string): string { return parent ? parent + "." + field : field; }`)
	g.P()
	for _, name := range names {
		message, err := schema.Message(name)
		if err != nil {
			return err
		}
		if err := emitMessage(g, message, rules.Messages[name], rules.Messages); err != nil {
			return fmt.Errorf("TypeScript validation %s: %w", name, err)
		}
	}
	rootMessage, err := schema.Message(root)
	if err != nil {
		return err
	}
	if _, exists := rules.Messages[root]; !exists {
		if err := emitMessage(g, rootMessage, spec.Message{}, rules.Messages); err != nil {
			return fmt.Errorf("TypeScript validation %s: %w", root, err)
		}
	}
	g.P(`export function validate`, root, `(value: `, root, `, provider = "", operation = "create"): void { validate`, root, `Value(value, provider, operation, ""); }`)
	g.P()
	return nil
}

func emitMessage(g *protogen.GeneratedFile, message protoreflect.MessageDescriptor, rule spec.Message, all map[string]spec.Message) error {
	name := string(message.Name())
	g.P(`function validate`, name, `Value(value: `, name, ` | undefined, provider: string, operation: string, path: string): void {`)
	g.P(`  if (value === undefined) return;`)
	fields := make([]string, 0, len(rule.Fields))
	for field := range rule.Fields {
		fields = append(fields, field)
	}
	sort.Strings(fields)
	for _, rawName := range fields {
		field := message.Fields().ByName(protoreflect.Name(rawName))
		if field == nil {
			return fmt.Errorf("unknown field %q", rawName)
		}
		if err := emitField(g, field, rule.Fields[rawName]); err != nil {
			return err
		}
	}
	for i, constraint := range rule.Constraints {
		if err := emitConstraint(g, message, constraint, i); err != nil {
			return err
		}
	}
	// Descend through nested messages that participate in shared validation.
	for i := 0; i < message.Fields().Len(); i++ {
		field := message.Fields().Get(i)
		if field.Message() == nil || field.IsMap() || field.Message().FullName() == "google.protobuf.Duration" {
			continue
		}
		child := string(field.Message().Name())
		if _, ok := all[child]; !ok {
			continue
		}
		fieldName := naming.Camel(string(field.Name()))
		fieldPath := `joinValidationPath(path, "` + fieldName + `")`
		if realOneof(field) {
			g.P(`  if (value.`, oneofDiscriminator(field), ` === "`, fieldName, `") validate`, child, `Value(value.`, fieldName, `, provider, operation, `, fieldPath, `);`)
			continue
		}
		if field.IsList() {
			g.P(`  for (const item of value.`, fieldName, ` ?? []) validate`, child, `Value(item, provider, operation, `, fieldPath, ` + "[]");`)
			continue
		}
		g.P(`  validate`, child, `Value(value.`, fieldName, `, provider, operation, `, fieldPath, `);`)
	}
	g.P(`}`)
	g.P()
	return nil
}

func emitField(g *protogen.GeneratedFile, field protoreflect.FieldDescriptor, rule spec.Field) error {
	name := naming.Camel(string(field.Name()))
	path := `joinValidationPath(path, "` + name + `")`
	value := `value.` + name
	if rule.Each && !field.IsList() {
		return fmt.Errorf("%s: each requires a list", field.FullName())
	}
	if rule.MinItems != nil {
		if !field.IsList() || *rule.MinItems < 0 {
			return fmt.Errorf("%s: min_items requires a list and nonnegative bound", field.FullName())
		}
		g.P(`  if (`, value, ` !== undefined && `, value, `.length < `, *rule.MinItems, `) validationError(provider, operation, `, path, `, "must contain at least `, *rule.MinItems, ` item(s)");`)
	}
	if field.IsList() {
		if rule.Nonblank || rule.Finite || rule.Minimum != nil || rule.ExclusiveMinimum != nil || rule.Maximum != nil || rule.Format != "" {
			g.P(`  for (const item of `, value, ` ?? []) {`)
			if err := emitScalarRules(g, field, rule, "item", path+` + "[]"`); err != nil {
				return err
			}
			g.P(`  }`)
		}
		return nil
	}
	if realOneof(field) {
		g.P(`  if (value.`, oneofDiscriminator(field), ` === "`, name, `") {`)
		if err := emitScalarRules(g, field, rule, value, path); err != nil {
			return err
		}
		g.P(`  }`)
		return nil
	}
	return emitScalarRules(g, field, rule, value, path)
}

func emitScalarRules(g *protogen.GeneratedFile, field protoreflect.FieldDescriptor, rule spec.Field, value, path string) error {
	if rule.Nonblank {
		if field.Kind() != protoreflect.StringKind {
			return fmt.Errorf("%s: nonblank requires string", field.FullName())
		}
		g.P(`  if (`, value, ` !== undefined && !`, value, `.trim()) validationError(provider, operation, `, path, `, "must not be blank");`)
	}
	if rule.Finite {
		if field.Kind() != protoreflect.DoubleKind && field.Kind() != protoreflect.FloatKind {
			return fmt.Errorf("%s: finite requires floating point", field.FullName())
		}
		g.P(`  if (`, value, ` !== undefined && !Number.isFinite(`, value, `)) validationError(provider, operation, `, path, `, "must be finite");`)
	}
	if rule.Format != "" {
		if field.Kind() != protoreflect.StringKind {
			return fmt.Errorf("%s: format requires string", field.FullName())
		}
		switch rule.Format {
		case "image_reference":
			g.P(`  if (`, value, ` !== undefined && (`, value, `.length > 512 || !/^(?:[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?(?::[0-9]+)?\/)?[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*(?:\/[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*)*(?::[A-Za-z0-9_][A-Za-z0-9_.-]{0,127})?(?:@sha256:[a-fA-F0-9]{64})?$/.test(`, value, `))) validationError(provider, operation, `, path, `, "must be a valid image reference");`)
		default:
			return fmt.Errorf("%s: unsupported format %q", field.FullName(), rule.Format)
		}
	}
	if rule.Minimum != nil && rule.ExclusiveMinimum != nil {
		return fmt.Errorf("%s: choose one minimum", field.FullName())
	}
	for _, bound := range []*float64{rule.Minimum, rule.ExclusiveMinimum, rule.Maximum} {
		if bound != nil && (math.IsNaN(*bound) || math.IsInf(*bound, 0)) {
			return fmt.Errorf("%s: bounds must be finite", field.FullName())
		}
	}
	if rule.Minimum != nil {
		if err := emitBound(g, field, value, path, *rule.Minimum, `>=`, "must be at least"); err != nil {
			return err
		}
	}
	if rule.ExclusiveMinimum != nil {
		if err := emitBound(g, field, value, path, *rule.ExclusiveMinimum, `>`, "must be greater than"); err != nil {
			return err
		}
	}
	if rule.Maximum != nil {
		if err := emitBound(g, field, value, path, *rule.Maximum, `<=`, "must be at most"); err != nil {
			return err
		}
	}
	return nil
}

func emitBound(g *protogen.GeneratedFile, field protoreflect.FieldDescriptor, value, path string, bound float64, op, words string) error {
	if !numeric(field) {
		return fmt.Errorf("%s: numeric bound on nonnumeric field", field.FullName())
	}
	literal := strconv.FormatFloat(bound, 'f', -1, 64)
	if unsigned(field) {
		literal += "n"
	}
	if field.Kind() == protoreflect.EnumKind {
		value = string(field.Enum().Name()) + "Number[" + value + "]"
	}
	g.P(`  if (`, value, ` !== undefined && !(`, value, ` `, op, ` `, literal, `)) validationError(provider, operation, `, path, `, ", `, words, ` `, strconv.FormatFloat(bound, 'f', -1, 64), `");`)
	return nil
}

func emitConstraint(g *protogen.GeneratedFile, message protoreflect.MessageDescriptor, constraint spec.Constraint, index int) error {
	condition, err := conditionExpression(message, constraint.When)
	if err != nil {
		return err
	}
	fields := constraint.Fields
	if constraint.Field != "" {
		fields = []string{constraint.Field, constraint.Other}
	}
	if len(fields) == 0 {
		return fmt.Errorf("%s: %s requires fields", message.FullName(), constraint.Op)
	}
	active := make([]string, 0, len(fields))
	for _, raw := range fields {
		field := message.Fields().ByName(protoreflect.Name(raw))
		if field == nil {
			return fmt.Errorf("%s: unknown constraint field %q", message.FullName(), raw)
		}
		active = append(active, activeExpression(field))
	}
	label := constraint.Field
	if label == "" {
		label = strings.Join(fields, ",")
	}
	label = naming.Camel(label)
	switch constraint.Op {
	case "at_most_one", "exactly_one", "requires_any", "forbids":
		count := fmt.Sprintf("count%d", index)
		g.P(`  if (`, condition, `) { const `, count, ` = [`, strings.Join(active, `, `), `].filter(Boolean).length;`)
		invalid, message := "> 1", "allows at most one value"
		switch constraint.Op {
		case "exactly_one":
			invalid, message = "!== 1", "requires exactly one value"
		case "requires_any":
			invalid, message = "=== 0", "requires a value"
		case "forbids":
			invalid, message = "> 0", "forbids these values"
		}
		g.P(`    if (`, count, ` `, invalid, `) validationError(provider, operation, joinValidationPath(path, "`, label, `"), "`, message, `"); }`)
	case "limit", "requires_positive":
		if constraint.Field == "" || constraint.Other == "" {
			return fmt.Errorf("%s: %s requires field and other", message.FullName(), constraint.Op)
		}
		left := message.Fields().ByName(protoreflect.Name(constraint.Field))
		right := message.Fields().ByName(protoreflect.Name(constraint.Other))
		if left == nil || right == nil || !numeric(left) || !numeric(right) || left.Kind() != right.Kind() {
			return fmt.Errorf("%s: %s requires compatible numeric fields", message.FullName(), constraint.Op)
		}
		l, r := "value."+naming.Camel(constraint.Field), "value."+naming.Camel(constraint.Other)
		invalid, message := l+" !== undefined && "+r+" !== undefined && "+l+" > 0 && "+l+" < "+r, "must not be smaller than "+naming.Camel(constraint.Other)
		if constraint.Op == "requires_positive" {
			invalid, message = l+" !== undefined && "+l+" > 0 && !("+r+" !== undefined && "+r+" > 0)", "requires a positive "+naming.Camel(constraint.Other)
		}
		g.P(`  if (`, condition, ` && `, invalid, `) validationError(provider, operation, joinValidationPath(path, "`, label, `"), "`, message, `");`)
	default:
		return fmt.Errorf("%s: unsupported constraint %q", message.FullName(), constraint.Op)
	}
	return nil
}

func conditionExpression(message protoreflect.MessageDescriptor, condition *spec.Condition) (string, error) {
	if condition == nil {
		return "true", nil
	}
	if (condition.Equals == nil) == (condition.NotEquals == nil) {
		return "", fmt.Errorf("%s: condition needs exactly one comparison", message.FullName())
	}
	field := message.Fields().ByName(protoreflect.Name(condition.Field))
	if field == nil {
		return "", fmt.Errorf("%s: unknown condition field %q", message.FullName(), condition.Field)
	}
	value, op := condition.Equals, "==="
	if value == nil {
		value, op = condition.NotEquals, "!=="
	}
	if field.Kind() == protoreflect.EnumKind {
		enum := field.Enum().Values().ByNumber(protoreflect.EnumNumber(*value))
		if enum == nil {
			return "", fmt.Errorf("%s: unknown enum value %d", field.FullName(), *value)
		}
		prefix := naming.EnumPrefix(string(field.Enum().Name()))
		literal := strings.ToLower(strings.TrimPrefix(string(enum.Name()), prefix))
		return "value." + naming.Camel(condition.Field) + " " + op + " " + strconv.Quote(literal), nil
	}
	return "", fmt.Errorf("%s: condition requires enum field", field.FullName())
}

func activeExpression(field protoreflect.FieldDescriptor) string {
	value := "value." + naming.Camel(string(field.Name()))
	if realOneof(field) {
		return "value." + oneofDiscriminator(field) + " === " + strconv.Quote(naming.Camel(string(field.Name())))
	}
	if field.IsList() || field.IsMap() {
		return value + " !== undefined && " + value + ".length > 0"
	}
	if field.Kind() == protoreflect.BoolKind {
		return value + " === true"
	}
	if field.HasPresence() || field.Kind() == protoreflect.MessageKind {
		return value + " !== undefined"
	}
	if field.Kind() == protoreflect.StringKind {
		return value + " !== \"\""
	}
	if field.Kind() == protoreflect.EnumKind {
		return value + " !== undefined && " + value + " !== \"default\""
	}
	return value + " !== 0"
}

func realOneof(field protoreflect.FieldDescriptor) bool {
	return field.ContainingOneof() != nil && !field.ContainingOneof().IsSynthetic()
}

func oneofDiscriminator(field protoreflect.FieldDescriptor) string {
	return naming.Camel(string(field.ContainingOneof().Name()))
}

func numeric(field protoreflect.FieldDescriptor) bool {
	switch field.Kind() {
	case protoreflect.DoubleKind, protoreflect.FloatKind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.Int32Kind, protoreflect.Int64Kind:
		return true
	case protoreflect.EnumKind:
		return true
	case protoreflect.MessageKind:
		return field.Message().FullName() == "google.protobuf.Duration"
	default:
		return false
	}
}
func unsigned(field protoreflect.FieldDescriptor) bool {
	return field.Kind() == protoreflect.Uint64Kind || field.Kind() == protoreflect.Int64Kind
}
