package sandbox

import (
	"fmt"
	"reflect"
	"strings"
)

// RejectUnmapped reports fields that remain after a provider maps the ones it supports.
// Pass the remainder returned by a mapping helper. Reflection checks presence only; it
// does not call SDK methods. ErrorInfo.Field uses the protobuf field name from the
// generated struct tag (for example CpuCores → cpu_cores).
//
// Example:
//
//	native, remaining, err := mapResources(request.Resources)
//	if err != nil {
//	    return err
//	}
//	if err := sandbox.RejectUnmapped("daytona", remaining); err != nil {
//	    return err
//	}
func RejectUnmapped(provider string, remaining any) error {
	value := reflect.ValueOf(remaining)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return NewError(&ErrorInfo{
			Kind:     ErrorKindInternal,
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
		return NewError(&ErrorInfo{
			Kind:     ErrorKindUnsupported,
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
