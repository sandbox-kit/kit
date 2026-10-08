package sandbox

import (
	"fmt"
	"reflect"
	"strings"
)

// RejectUnmapped checks a copy of configuration after an adapter consumes the
// fields it maps. Reflection examines configuration presence, never SDK methods.
func RejectUnmapped(provider string, remaining any) error {
	value := reflect.ValueOf(remaining)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("sandbox-kit: invalid configuration support check")
	}
	var fields []string
	for i := 0; i < value.NumField(); i++ {
		if !value.Field(i).IsZero() {
			fields = append(fields, value.Type().Field(i).Name)
		}
	}
	if len(fields) > 0 {
		return fmt.Errorf("sandbox-kit %s: unsupported creation fields: %s", provider, strings.Join(fields, ", "))
	}
	return nil
}
