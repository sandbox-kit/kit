package typesgen

import (
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"strconv"
)

// GenerateRuntime emits Go ownership and error idioms with the shared contract settings.
func GenerateRuntime(p *protogen.Plugin, file *protogen.File, rules spec.Contracts) error {
	emitErrorRuntime(p, file, rules.Errors)
	emitCopyRuntime(p, file, rules.CopyMaxDepth)
	return nil
}

func emitErrorRuntime(p *protogen.Plugin, file *protogen.File, policy spec.ErrorPolicy) {
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".errors.runtime.gen.go", file.GoImportPath)
	g.P("// Code generated from contracts.yaml and errors.yaml by the Go emitter. DO NOT EDIT.")
	g.P("package ", file.GoPackageName)
	context := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "context"} }
	errors := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "errors"} }
	strings := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "strings"} }
	validator := func(member string) protogen.GoIdent {
		return protogen.GoIdent{GoName: member, GoImportPath: "github.com/go-playground/validator/v10"}
	}
	g.P()
	g.P("// Error exposes portable details and retains the original error for ", errors("Is"), "/As.")
	g.P("type Error struct {")
	g.P("\tErrorInfo")
	g.P("\tcause error")
	g.P("}")
	g.P()
	g.P("// Error includes a semantic field path when the error concerns a field.")
	g.P("func (e *Error) Error() string {")
	g.P("\tif e == nil {")
	g.P("\t\treturn \"\"")
	g.P("\t}")
	g.P("\tif e.Field != \"\" {")
	g.P("\t\treturn e.Message + \" (\" + e.Field + \")\"")
	g.P("\t}")
	g.P("\treturn e.Message")
	g.P("}")
	g.P()
	g.P("// Unwrap preserves native SDK and context matching through ", errors("Is"), "/As.")
	g.P("func (e *Error) Unwrap() error {")
	g.P("\tif e == nil {")
	g.P("\t\treturn nil")
	g.P("\t}")
	g.P("\treturn e.cause")
	g.P("}")
	g.P()
	g.P("// NewError owns a copy of portable details and retains the original cause.")
	g.P("func NewError(info ErrorInfo, cause error) *Error {")
	g.P("\towned, _ := info.Clone()")
	g.P("\tif !owned.Kind.Valid() {")
	g.P("\t\towned.Kind = ErrorKindUnknown")
	g.P("\t}")
	for i, rule := range policy.CauseRules {
		member := "Canceled"
		if rule.Match == "deadline_exceeded" {
			member = "DeadlineExceeded"
		}
		prefix := "if"
		if i > 0 {
			prefix = "} else if"
		}
		kind := "ErrorKindCanceled"
		if rule.Kind == "timeout" {
			kind = "ErrorKindTimeout"
		}
		g.P(prefix, " ", errors("Is"), "(cause, ", context(member), ") {")
		g.P("owned.Kind = ", kind)
	}
	if len(policy.CauseRules) > 0 {
		g.P("}")
	}
	g.P("\treturn &Error{ErrorInfo: *owned, cause: cause}")
	g.P("}")
	g.P()
	g.P("// WithErrorContext adds provider/operation details without mutating an error.")
	g.P("func WithErrorContext(err error, provider, operation string) error {")
	g.P("\tif err == nil {")
	g.P("\t\treturn nil")
	g.P("\t}")
	g.P("\tvar existing *Error")
	g.P("\tif ", errors("As"), "(err, &existing) && existing != nil {")
	if policy.ExistingError.ReuseMatchingContext {
		g.P("\t\tif (provider==\"\" || existing.Provider==provider) && (operation==\"\" || existing.Operation==operation){return err}")
	}
	g.P("\t\tinfo, _ := existing.ErrorInfo.Clone()")
	g.P("\t\tif provider != \"\" {")
	g.P("\t\t\tinfo.Provider = provider")
	g.P("\t\t}")
	g.P("\t\tif operation != \"\" {")
	g.P("\t\t\tinfo.Operation = operation")
	g.P("\t\t}")
	g.P("\t\treturn &Error{ErrorInfo: *info, cause: err}")
	g.P("\t}")
	g.P("\treturn ", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindUnknown"), Provider: "provider", Operation: "operation", Message: "err.Error()", Cause: "err"}))
	g.P("}")
	g.P()
	g.P("func validationError(err error) error {")
	g.P("\tfield := \"\"")
	g.P("\tvar fields ", validator("ValidationErrors"))
	g.P("\tif ", errors("As"), "(err, &fields) && len(fields) > 0 {")
	g.P("\t\tfield = validationPath(fields[0].StructNamespace())")
	g.P("\t}")
	g.P("\treturn ", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindInvalidArgument"), Field: "field", Message: strconv.Quote("sandbox-kit: invalid configuration"), Cause: "err"}))
	g.P("}")
	g.P()
	g.P("func validationPath(namespace string) string {")
	g.P("\tparts := ", strings("Split"), "(namespace, \".\")")
	g.P("\tif len(parts) > 1 {")
	g.P("\t\tparts = parts[1:]")
	g.P("\t}")
	g.P("\tfor i, part := range parts {")
	g.P("\t\tbase, suffix := part, \"\"")
	g.P("\t\tif index := ", strings("IndexByte"), "(part, '['); index >= 0 {")
	g.P("\t\t\tbase, suffix = part[:index], part[index:]")
	g.P("\t\t}")
	g.P("\t\tif native, ok := validationFieldNames[base]; ok {")
	g.P("\t\t\tparts[i] = native + suffix")
	g.P("\t\t}")
	g.P("\t\tif native, ok := configValidationFieldNames[base]; ok {")
	g.P("\t\t\tparts[i] = native + suffix")
	g.P("\t\t}")
	g.P("\t}")
	g.P("\treturn ", strings("Join"), "(parts, \".\")")
	g.P("}")
}

func emitCopyRuntime(p *protogen.Plugin, file *protogen.File, maxDepth int) {
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".clone.runtime.gen.go", file.GoImportPath)
	g.P("// Code generated from contracts.yaml and errors.yaml by the Go emitter. DO NOT EDIT.")
	g.P("package ", file.GoPackageName)
	fmt := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "fmt"} }
	math := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "math"} }
	g.P()
	g.P("// cloneMetadata copies declarative metadata, never SDK objects.")
	g.P("// It preserves optional values and ownership boundaries.")
	g.P("func cloneMetadata(source map[string]any) (map[string]any, error) {")
	g.P("\tif source == nil {")
	g.P("\t\treturn nil, nil")
	g.P("\t}")
	g.P("\tvalue, err := copyMetadataValue(source, 0)")
	g.P("\tif err != nil {")
	g.P("\t\treturn nil, err")
	g.P("\t}")
	g.P("\treturn value.(map[string]any), nil")
	g.P("}")
	g.P()
	g.P("// Metadata is declarative: finite primitive values, bytes, lists, and string-key")
	g.P("// objects. Depth bounds reject cycles without retaining caller-owned objects.")
	g.P("func copyMetadataValue(value any, depth int) (any, error) {")
	g.P("\tif depth > ", maxDepth, " {")
	g.P("\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds maximum depth or contains a cycle\")")
	g.P("\t}")
	g.P("\tswitch v := value.(type) {")
	g.P("\tcase nil, bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:")
	g.P("\t\treturn v, nil")
	g.P("\tcase float64:")
	g.P("\t\tif ", math("IsNaN"), "(v) || ", math("IsInf"), "(v, 0) {")
	g.P("\t\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: metadata requires finite numbers\")")
	g.P("\t\t}")
	g.P("\t\treturn v, nil")
	g.P("\tcase float32:")
	g.P("\t\tif ", math("IsNaN"), "(float64(v)) || ", math("IsInf"), "(float64(v), 0) {")
	g.P("\t\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: metadata requires finite numbers\")")
	g.P("\t\t}")
	g.P("\t\treturn v, nil")
	g.P("\tcase []byte:")
	g.P("\t\tif v == nil {")
	g.P("\t\t\treturn []byte(nil), nil")
	g.P("\t\t}")
	g.P("\t\treturn append([]byte{}, v...), nil")
	g.P("\tcase []string:")
	g.P("\t\tif v == nil {")
	g.P("\t\t\treturn []string(nil), nil")
	g.P("\t\t}")
	g.P("\t\treturn append([]string{}, v...), nil")
	g.P("\tcase []any:")
	g.P("\t\tif v == nil {")
	g.P("\t\t\treturn []any(nil), nil")
	g.P("\t\t}")
	g.P("\t\tout := make([]any, len(v))")
	g.P("\t\tfor i, item := range v {")
	g.P("\t\t\tcopied, err := copyMetadataValue(item, depth+1)")
	g.P("\t\t\tif err != nil {")
	g.P("\t\t\t\treturn nil, err")
	g.P("\t\t\t}")
	g.P("\t\t\tout[i] = copied")
	g.P("\t\t}")
	g.P("\t\treturn out, nil")
	g.P("\tcase map[string]any:")
	g.P("\t\tif v == nil {")
	g.P("\t\t\treturn map[string]any(nil), nil")
	g.P("\t\t}")
	g.P("\t\tout := make(map[string]any, len(v))")
	g.P("\t\tfor key, item := range v {")
	g.P("\t\t\tcopied, err := copyMetadataValue(item, depth+1)")
	g.P("\t\t\tif err != nil {")
	g.P("\t\t\t\treturn nil, err")
	g.P("\t\t\t}")
	g.P("\t\t\tout[key] = copied")
	g.P("\t\t}")
	g.P("\t\treturn out, nil")
	g.P("\tdefault:")
	g.P("\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: unsupported metadata value type %T\", value)")
	g.P("\t}")
	g.P("}")
}
