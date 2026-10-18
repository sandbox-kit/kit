// Package goemit previously dispatched Kit generation phases. Emitters are
// retired; provider and shared SDK helpers are handwritten.
package goemit

import (
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"google.golang.org/protobuf/compiler/protogen"
)

// Generate is a no-op. Kit custom emitters have been removed.
func Generate(plugin *protogen.Plugin, file *protogen.File, program *model.Program) error {
	_ = plugin
	_ = file
	_ = program
	return nil
}
