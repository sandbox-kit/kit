// Package typesgen emits descriptor-driven TypeScript declarations.
package typesgen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	tsschema "github.com/sandbox-kit/kit/tooling/internal/typescript/schema"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func Generate(plugin *protogen.Plugin, file *protogen.File, profile spec.LanguageProfile) error {
	g := declarationgen.GeneratedFile(plugin, "sandbox"+profile.Outputs["shared"], string(file.Desc.Path()))
	g.P("/** A duration in milliseconds. */")
	g.P("export type DurationMs = number;")
	g.P("export type MetadataValue = null | boolean | string | number | bigint | Uint8Array | readonly MetadataValue[] | MetadataObject;")
	g.P("export interface MetadataObject { readonly [key: string]: MetadataValue; }")
	g.P()
	return EmitDeclarations(g, file, true)
}

func EmitDeclarations(g *protogen.GeneratedFile, file *protogen.File, skipMetadata bool) error {
	for _, enum := range file.Enums {
		emitEnum(g, enum.Desc)
	}
	for _, message := range file.Messages {
		if skipMetadata && (message.Desc.FullName() == "kit.sandbox.v1.MetadataValue" || message.Desc.FullName() == "kit.sandbox.v1.MetadataList" || message.Desc.FullName() == "kit.sandbox.v1.MetadataObject") {
			continue
		}
		if err := emitMessage(g, message.Desc); err != nil {
			return err
		}
	}
	return nil
}

func emitEnum(g *protogen.GeneratedFile, enum protoreflect.EnumDescriptor) {
	name := string(enum.Name())
	prefix := naming.EnumPrefix(name)
	g.P("export const ", name, " = {")
	for i := 0; i < enum.Values().Len(); i++ {
		value := enum.Values().Get(i)
		short := strings.TrimPrefix(string(value.Name()), prefix)
		g.P("  ", naming.Pascal(strings.ToLower(short)), ": ", fmt.Sprintf("%q", strings.ToLower(short)), ",")
	}
	g.P("} as const;")
	g.P("export type ", name, " = (typeof ", name, ")[keyof typeof ", name, "];")
	g.P("export const ", name, "Number: Readonly<Record<", name, ", number>> = {")
	for i := 0; i < enum.Values().Len(); i++ {
		value := enum.Values().Get(i)
		short := strings.TrimPrefix(string(value.Name()), prefix)
		key := naming.Pascal(strings.ToLower(short))
		g.P("  [", name, ".", key, "]: ", value.Number(), ",")
	}
	g.P("};")
	g.P()
}

func emitMessage(g *protogen.GeneratedFile, message protoreflect.MessageDescriptor) error {
	if realOneofs(message) > 0 {
		return emitOneofMessage(g, message)
	}
	g.P("export interface ", message.Name(), " {")
	for i := 0; i < message.Fields().Len(); i++ {
		field := message.Fields().Get(i)
		typeName, err := tsschema.FieldType(field)
		if err != nil {
			return err
		}
		optional := field.HasPresence() || field.IsList() || field.IsMap()
		marker := ""
		if optional {
			marker = "?"
		}
		g.P("  readonly ", naming.Camel(string(field.Name())), marker, ": ", typeName, ";")
	}
	g.P("}")
	g.P()
	return nil
}

func emitOneofMessage(g *protogen.GeneratedFile, message protoreflect.MessageDescriptor) error {
	if realOneofs(message) != 1 {
		return fmt.Errorf("%s: multiple oneofs are not yet representable", message.FullName())
	}
	var oneof protoreflect.OneofDescriptor
	for i := 0; i < message.Oneofs().Len(); i++ {
		if !message.Oneofs().Get(i).IsSynthetic() {
			oneof = message.Oneofs().Get(i)
			break
		}
	}
	discriminator := naming.Camel(string(oneof.Name()))
	variants := make([]string, 0, oneof.Fields().Len()+1)
	for i := 0; i < oneof.Fields().Len(); i++ {
		field := oneof.Fields().Get(i)
		typeName, err := tsschema.FieldType(field)
		if err != nil {
			return err
		}
		fieldName := naming.Camel(string(field.Name()))
		variants = append(variants, fmt.Sprintf("{ readonly %s: %q; readonly %s: %s }", discriminator, fieldName, fieldName, typeName))
	}
	variants = append(variants, fmt.Sprintf("{ readonly %s?: never }", discriminator))
	g.P("export type ", message.Name(), " =")
	for i, variant := range variants {
		end := ""
		if i == len(variants)-1 {
			end = ";"
		}
		g.P("  | ", variant, end)
	}
	g.P()
	return nil
}

func realOneofs(message protoreflect.MessageDescriptor) int {
	count := 0
	for i := 0; i < message.Oneofs().Len(); i++ {
		if !message.Oneofs().Get(i).IsSynthetic() {
			count++
		}
	}
	return count
}
