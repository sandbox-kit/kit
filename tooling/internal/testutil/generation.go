// Package testutil supplies real specification inputs for emitter tests.
package testutil

import (
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"testing"
)

func Generation(t testing.TB, root string) (model.Templates, spec.LanguageProfile) {
	t.Helper()
	raw, err := spec.LoadTemplates(root)
	if err != nil {
		t.Fatal(err)
	}
	templates, err := model.CompileTemplates(raw)
	if err != nil {
		t.Fatal(err)
	}
	profile, err := spec.LoadLanguageProfile(root, "go")
	if err != nil {
		t.Fatal(err)
	}
	return templates, profile
}
