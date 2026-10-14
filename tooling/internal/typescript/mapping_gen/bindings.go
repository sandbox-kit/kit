// Package mappinggen lowers compiled portable mappings into TypeScript SDK
// bindings. This file validates the semantic-to-native member contract before
// any TypeScript is emitted.
package mappinggen

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/model"
)

// Binding describes the native members required by one portable mapping group.
type Binding struct {
	Import string
	Fields map[string]string
}

// ValidateBindings verifies that YAML contains every native symbol used by a
// TypeScript provider emitter. Response mappings bind their native source;
// request mappings bind their native target.
func ValidateBindings(provider model.Provider, required map[string]Binding) error {
	remaining := make(map[string]Binding, len(required))
	for name, binding := range required {
		remaining[name] = binding
	}
	for _, group := range provider.Groups {
		expected, ok := remaining[group.Rule.Name]
		if !ok {
			continue
		}
		actual, ok := group.Rule.Target.Bindings["typescript"]
		if group.Rule.Direction == "response" {
			actual, ok = group.Rule.Source.Bindings["typescript"]
		}
		if !ok || actual.Import != expected.Import {
			return fmt.Errorf("%s TypeScript: missing native binding for %s", provider.Rules.Provider, group.Rule.Name)
		}
		for semantic, native := range expected.Fields {
			if actual.Fields[semantic] != native {
				return fmt.Errorf("%s TypeScript: %s.%s must bind to %s", provider.Rules.Provider, group.Rule.Name, semantic, native)
			}
		}
		delete(remaining, group.Rule.Name)
	}
	if len(remaining) == 0 {
		return nil
	}
	names := make([]string, 0, len(remaining))
	for name := range remaining {
		names = append(names, name)
	}
	sort.Strings(names)
	return fmt.Errorf("%s TypeScript: missing mapping groups %s", provider.Rules.Provider, strings.Join(names, ", "))
}
