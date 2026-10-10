package model

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"regexp"
	"strings"
)

// Templates are declaration patterns instantiated from descriptors and mappings.
type Templates struct {
	Callables  map[string]spec.DeclarationTemplate
	FieldPaths spec.FieldPathTemplates
}

var templateToken = regexp.MustCompile(`\{([a-z_]+)\}`)

func CompileTemplates(raw spec.Templates) (Templates, error) {
	result := Templates{Callables: raw.Callables, FieldPaths: raw.FieldPaths}
	if raw.Version != 1 || len(raw.Callables) == 0 {
		return result, fmt.Errorf("unsupported or empty declaration templates")
	}
	for _, name := range []string{"getter", "clone", "clone_internal", "validator_factory", "validation", "request_mapping", "response_mapping", "provider_checks", "response_presence", "response_origins"} {
		if _, ok := raw.Callables[name]; !ok {
			return result, fmt.Errorf("missing declaration template %s", name)
		}
	}
	for id, t := range raw.Callables {
		if !declarationIdentifier.MatchString(id) {
			return result, fmt.Errorf("invalid template ID")
		}
		declaration := t.Declaration
		declaration.Name = templateToken.ReplaceAllString(declaration.Name, "Binding")
		if err := validateCallable(declaration, false); err != nil {
			return result, fmt.Errorf("template %s: %w", id, err)
		}
		if t.Owner != "" && !declarationIdentifier.MatchString(t.Owner) {
			return result, fmt.Errorf("invalid template owner")
		}
	}
	if len(raw.FieldPaths.Roots) == 0 || raw.FieldPaths.Name != "{prefix}{fields}" {
		return result, fmt.Errorf("unsupported field-path template")
	}
	for root, prefix := range raw.FieldPaths.Roots {
		if !declarationIdentifier.MatchString(root) || !declarationIdentifier.MatchString(prefix) {
			return result, fmt.Errorf("invalid field-path root or prefix")
		}
	}
	return result, nil
}
func (t Templates) Instantiate(id string, names map[string]string) (spec.DeclarationTemplate, error) {
	pattern, ok := t.Callables[id]
	if !ok {
		return pattern, fmt.Errorf("unknown declaration template %s", id)
	}
	var missing string
	replace := func(token string) string {
		name := strings.Trim(token, "{}")
		value, ok := names[name]
		if !ok {
			missing = name
		}
		return value
	}
	pattern.Declaration.Name = templateToken.ReplaceAllStringFunc(pattern.Declaration.Name, replace)
	pattern.Declaration.Doc = templateToken.ReplaceAllStringFunc(pattern.Declaration.Doc, replace)
	if missing != "" {
		return pattern, fmt.Errorf("template %s needs name binding %s", id, missing)
	}
	if err := validateCallable(pattern.Declaration, false); err != nil {
		return pattern, err
	}
	return pattern, nil
}
