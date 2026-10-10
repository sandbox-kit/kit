package sandbox

import (
	"fmt"
	"reflect"
	"strings"
)

// RejectUnmapped reports fields that remain after a provider maps the ones it supports.
// Pass the remainder returned by a generated mapping helper. Reflection checks presence
// only; it does not call SDK methods.
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
		return NewError(ErrorInfo{
			Kind:     ErrorKindInternal,
			Provider: provider,
			Message:  "sandbox-kit: invalid configuration support check",
		}, nil)
	}
	var fields []string
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			fields = append(fields, value.Type().Field(i).Name)
		}
	}
	if len(fields) > 0 {
		return NewError(ErrorInfo{
			Kind:     ErrorKindUnsupported,
			Provider: provider,
			Field:    validationPath("Root." + fields[0]),
			Message:  fmt.Sprintf("sandbox-kit %s: unsupported fields: %s", provider, strings.Join(fields, ", ")),
		}, nil)
	}
	return nil
}
