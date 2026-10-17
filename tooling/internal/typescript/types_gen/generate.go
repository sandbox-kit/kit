// Package typesgen previously emitted portable TypeScript declarations.
// Portable sandbox and client data now come from @sandbox-kit/contracts;
// the runtime aliases those messages by hand.
package typesgen

import (
	"fmt"

	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generate is retired. Callers must not emit a second parallel type surface.
func Generate(plugin *protogen.Plugin, file *protogen.File, program *model.Program, profile spec.LanguageProfile) error {
	_ = plugin
	_ = file
	_ = program
	_ = profile
	return fmt.Errorf("TypeScript portable types come from @sandbox-kit/contracts; types_gen is retired")
}
