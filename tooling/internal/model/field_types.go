package model

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// FieldType projects a schema field into a portable native type. Getter mode
// removes scalar presence while retaining message, map, and sequence shapes.
func FieldType(field protoreflect.FieldDescriptor, getter bool) (spec.TypeRef, error) {
	if field.IsMap() {
		key, err := FieldType(field.MapKey(), true)
		if err != nil {
			return spec.TypeRef{}, err
		}
		value, err := FieldType(field.MapValue(), false)
		return spec.TypeRef{Map: &spec.MapType{Key: key, Value: value}}, err
	}
	var t spec.TypeRef
	switch field.Kind() {
	case protoreflect.StringKind:
		t.Builtin = "string"
	case protoreflect.BoolKind:
		t.Builtin = "boolean"
	case protoreflect.DoubleKind:
		t.Builtin = "number"
	case protoreflect.FloatKind:
		t.Builtin = "float32"
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		t.Builtin = "int32"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		t.Builtin = "integer"
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		t.Builtin = "uint32"
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		t.Builtin = "unsigned"
	case protoreflect.BytesKind:
		t.Builtin = "bytes"
	case protoreflect.EnumKind:
		t.Ref = "schema_enum." + string(field.Enum().FullName())
	case protoreflect.MessageKind:
		switch field.Message().FullName() {
		case "google.protobuf.Duration":
			t = spec.TypeRef{Ref: "duration", Reference: true}
		case "google.protobuf.Timestamp":
			t = spec.TypeRef{Ref: "timestamp", Reference: true}
		case "google.protobuf.Struct", "kit.sandbox.v1.MetadataObject":
			t.Map = &spec.MapType{Key: spec.TypeRef{Builtin: "string"}, Value: spec.TypeRef{Builtin: "any"}}
		default:
			t = spec.TypeRef{Ref: "schema." + string(field.Message().FullName()), Reference: true}
		}
	default:
		return t, fmt.Errorf("unsupported field type %s", field.FullName())
	}
	if field.IsList() {
		return spec.TypeRef{Sequence: &t}, nil
	}
	if !getter && field.HasPresence() && field.Kind() != protoreflect.MessageKind {
		t.Optional = true
	}
	return t, nil
}
