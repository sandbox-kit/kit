package providergen

import (
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestCheckConditionLowersAllowedValues(t *testing.T) {
	t.Parallel()
	path, err := configTestSchema(t).Resolve("Config", "region")
	if err != nil {
		t.Fatal(err)
	}

	condition, kind, err := checkCondition(spec.Check{Allowed: []string{"us", "eu"}}, path.Fields[0], "value")
	if err != nil {
		t.Fatal(err)
	}
	if condition != `value !== "us" && value !== "eu"` || kind != "Unsupported" {
		t.Fatalf("checkCondition() = (%q, %q), want allowed-value unsupported condition", condition, kind)
	}
}

func TestCheckConditionPreservesAllowEmptyFormat(t *testing.T) {
	t.Parallel()
	path, err := configTestSchema(t).Resolve("Config", "region")
	if err != nil {
		t.Fatal(err)
	}

	condition, kind, err := checkCondition(spec.Check{Format: "domain", AllowEmpty: true}, path.Fields[0], "value")
	if err != nil {
		t.Fatal(err)
	}
	if condition != `value !== "" && !validDomain(value)` || kind != "InvalidArgument" {
		t.Fatalf("checkCondition() = (%q, %q), want allow-empty domain condition", condition, kind)
	}
}
