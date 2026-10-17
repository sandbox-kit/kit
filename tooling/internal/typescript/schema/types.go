// Package schema lowers protobuf descriptor types into TypeScript type syntax.
package schema

import (
	"fmt"

	"google.golang.org/protobuf/reflect/protoreflect"
)

func FieldType(field protoreflect.FieldDescriptor) (string, error) {
	var result string
	switch field.Kind() {
	case protoreflect.StringKind:
		result = "string"
	case protoreflect.BoolKind:
		result = "boolean"
	case protoreflect.DoubleKind, protoreflect.FloatKind, protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind, protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		result = "number"
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind, protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		result = "bigint"
	case protoreflect.BytesKind:
		result = "Uint8Array"
	case protoreflect.EnumKind:
		result = string(field.Enum().Name())
	case protoreflect.MessageKind:
		switch field.Message().FullName() {
		case "google.protobuf.Duration":
			result = "Duration"
		case "google.protobuf.Timestamp":
			result = "Date"
		case "kit.sandbox.v1.MetadataObject":
			result = "MetadataObject"
		default:
			result = string(field.Message().Name())
		}
	default:
		return "", fmt.Errorf("%s: unsupported field kind %s", field.FullName(), field.Kind())
	}
	if field.IsMap() {
		value, err := FieldType(field.MapValue())
		if err != nil {
			return "", err
		}
		return "Readonly<Record<string, " + value + ">>", nil
	}
	if field.IsList() {
		return "readonly " + result + "[]", nil
	}
	return result, nil
}
