// Package naming applies Go naming conventions while preserving shared schema names.
package naming

import (
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

// identifier formats schema words using Go initialisms and conventional unit casing.
func identifier(words string) string {
	initialisms := map[string]string{"id": "ID", "api": "API", "oauth": "OAuth", "jwt": "JWT", "cpu": "CPU", "gpu": "GPU", "vm": "VM", "microvm": "MicroVM", "pty": "PTY", "cidr": "CIDR", "cidrs": "CIDRs", "url": "URL", "http": "HTTP", "http2": "HTTP2", "tls": "TLS", "tcp": "TCP", "mib": "MiB"}
	var result strings.Builder
	for _, word := range strings.Split(strings.ToLower(words), "_") {
		if word == "" {
			continue
		}
		if value, ok := initialisms[word]; ok {
			result.WriteString(value)
		} else {
			result.WriteString(strings.ToUpper(word[:1]))
			result.WriteString(word[1:])
		}
	}
	return result.String()
}
func FieldName(field *protogen.Field) string { return identifier(string(field.Desc.Name())) }
func EnumValueName(value *protogen.EnumValue) string {
	// Protobuf constants repeat the enum name in SCREAMING_SNAKE_CASE.
	enum := value.Desc.Parent()
	name := string(value.Desc.Name())
	prefix := strings.ToUpper(snake(string(enum.Name()))) + "_"
	return string(enum.Name()) + identifier(strings.TrimPrefix(name, prefix))
}
func snake(name string) string {
	var result strings.Builder
	for i, r := range name {
		if r >= 'A' && r <= 'Z' && i > 0 {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return result.String()
}
