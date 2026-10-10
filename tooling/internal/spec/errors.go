package spec

import (
	"fmt"
	"path/filepath"
	"slices"
)

// ErrorPolicy describes portable classification and enrichment behavior.
type ErrorPolicy struct {
	Version       int         `yaml:"version"`
	FallbackKind  string      `yaml:"fallback_kind"`
	CauseRules    []CauseRule `yaml:"cause_rules"`
	ExistingError struct {
		ReuseMatchingContext bool `yaml:"reuse_matching_context"`
		PreserveDetails      bool `yaml:"preserve_details"`
	} `yaml:"existing_error"`
	ClassificationPrecedence []string `yaml:"classification_precedence"`
}
type CauseRule struct {
	Match string `yaml:"match"`
	Kind  string `yaml:"kind"`
}

func LoadErrorPolicy(root string) (ErrorPolicy, error) {
	var p ErrorPolicy
	if err := decode(filepath.Join(root, "errors.yaml"), &p); err != nil {
		return p, err
	}
	if p.Version != 1 || p.FallbackKind != "unknown" || !p.ExistingError.PreserveDetails || !slices.Equal(p.ClassificationPrecedence, []string{"existing_error", "native_type", "http_status", "grpc_code", "fallback"}) {
		return p, fmt.Errorf("unsupported error policy")
	}
	seen := map[string]bool{}
	for _, r := range p.CauseRules {
		if seen[r.Match] || (r.Match != "canceled" && r.Match != "deadline_exceeded") || (r.Kind != "canceled" && r.Kind != "timeout") {
			return p, fmt.Errorf("invalid cause rule %q", r.Match)
		}
		seen[r.Match] = true
	}
	return p, nil
}
