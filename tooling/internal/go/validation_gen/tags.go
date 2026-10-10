package validationgen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

// Tags translates portable field rules into the Go validator dialect.
func Tags(rule spec.Field, optional bool, repeated bool) string {
	var tags []string
	if optional {
		tags = append(tags, "omitnil")
	}
	if rule.MinItems != nil {
		tags = append(tags, fmt.Sprintf("min=%d", *rule.MinItems))
	}
	if rule.Each {
		tags = append(tags, "dive")
	}
	if rule.Finite {
		tags = append(tags, "finite")
	}
	if rule.Nonblank {
		tags = append(tags, "nonblank")
	}
	if rule.Format != "" {
		tags = append(tags, rule.Format)
	}
	if rule.Minimum != nil {
		tags = append(tags, fmt.Sprintf("gte=%g", *rule.Minimum))
	}
	if rule.ExclusiveMinimum != nil {
		tags = append(tags, fmt.Sprintf("gt=%g", *rule.ExclusiveMinimum))
	}
	if rule.Maximum != nil {
		tags = append(tags, fmt.Sprintf("lte=%g", *rule.Maximum))
	}
	if repeated && !rule.Each {
		tags = append(tags, "dive", "required")
	}
	return strings.Join(tags, ",")
}
