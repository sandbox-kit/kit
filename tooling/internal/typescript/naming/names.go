// Package naming translates portable protobuf names to TypeScript identifiers.
package naming

import (
	"strings"
	"unicode"
)

func Camel(value string) string {
	parts := strings.Split(value, "_")
	if len(parts) == 0 {
		return value
	}
	result := strings.ToLower(parts[0])
	for _, part := range parts[1:] {
		result += Pascal(strings.ToLower(part))
	}
	return result
}

func Pascal(value string) string {
	parts := strings.FieldsFunc(value, func(r rune) bool { return r == '_' || r == '-' })
	for index, part := range parts {
		runes := []rune(part)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		parts[index] = string(runes)
	}
	return strings.Join(parts, "")
}

func EnumPrefix(name string) string {
	var parts []string
	var current []rune
	for index, r := range name {
		if index > 0 && unicode.IsUpper(r) {
			parts = append(parts, string(current))
			current = nil
		}
		current = append(current, unicode.ToUpper(r))
	}
	parts = append(parts, string(current))
	return strings.Join(parts, "_") + "_"
}
