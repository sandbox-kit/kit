package sdk

import (
	"fmt"
	"reflect"
	"strings"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
)

// RejectUnmapped reports create or config fields a provider did not accept.
// Provider packages call it after copying the fields they support. The returned
// error is ErrorKindUnsupported and Field is the setting the caller should remove,
// such as "resources.cpu_cores". A nil or empty remainder returns nil.
func RejectUnmapped(provider string, remaining any) error {
	value := reflect.ValueOf(remaining)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return NewError(&sandboxv1.ErrorInfo{
			Kind:     sandboxv1.ErrorKind_ERROR_KIND_INTERNAL,
			Provider: provider,
			Message:  "sandbox-kit: invalid configuration support check",
		}, nil)
	}
	var fields []string
	var firstPath string
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		// Skip unexported and protobuf runtime fields; only contract fields matter.
		if field.PkgPath != "" || field.Name == "state" || field.Name == "sizeCache" || field.Name == "unknownFields" {
			continue
		}
		if !value.Field(i).IsZero() {
			fields = append(fields, field.Name)
			if firstPath == "" {
				firstPath = protobufFieldName(field)
			}
		}
	}
	if len(fields) > 0 {
		return NewError(&sandboxv1.ErrorInfo{
			Kind:     sandboxv1.ErrorKind_ERROR_KIND_UNSUPPORTED,
			Provider: provider,
			Field:    firstPath,
			Message:  fmt.Sprintf("sandbox-kit %s: unsupported fields: %s", provider, strings.Join(fields, ", ")),
		}, nil)
	}
	return nil
}

// protobufFieldName returns the canonical protobuf field name from a generated
// struct tag. Falls back to the Go field name when the tag is missing.
func protobufFieldName(field reflect.StructField) string {
	for _, part := range strings.Split(field.Tag.Get("protobuf"), ",") {
		if name, ok := strings.CutPrefix(part, "name="); ok && name != "" {
			return name
		}
	}
	return field.Name
}
