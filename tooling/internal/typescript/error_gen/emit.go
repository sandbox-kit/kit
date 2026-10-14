// Package errorgen emits TypeScript construction helpers for shared errors.
package errorgen

import (
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

func EmitMappingError(g *protogen.GeneratedFile, provider string) {
	g.P("function mappingError(kind: ErrorKind, field: string, message: string): never {")
	g.P("  throw sandboxError({ kind, provider: ", strconv.Quote(provider), ", operation: \"create\", field, message: \"sandbox-kit ", provider, ": \" + message });")
	g.P("}")
}
