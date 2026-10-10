package typesgen

import (
	"strconv"

	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
)

// GenerateRuntime emits Go ownership and error idioms with the shared contract settings.
func GenerateRuntime(p *protogen.Plugin, file *protogen.File, rules spec.Contracts, surface model.APIModule, profile spec.LanguageProfile) error {
	if err := emitErrorRuntime(p, file, rules.Errors, surface, profile); err != nil {
		return err
	}
	if err := emitCopyRuntime(p, file, rules, surface, profile); err != nil {
		return err
	}
	return nil
}

func emitErrorRuntime(p *protogen.Plugin, file *protogen.File, policy spec.ErrorPolicy, surface model.APIModule, profile spec.LanguageProfile) error {
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+profile.Outputs["errors"], file.GoImportPath)
	declarationgen.Banner(g, "specs/contracts.yaml and specs/errors.yaml",
		"Error carries portable details and the original cause. Use errors.As, then switch on Kind.",
		"var detail *Error\nif errors.As(err, &detail) {\n    fmt.Println(detail.Kind, detail.Field)\n}")
	g.P("package ", file.GoPackageName)
	context := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "context"} }
	errors := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "errors"} }
	strings := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "strings"} }
	validator := func(member string) protogen.GoIdent {
		return protogen.GoIdent{GoName: member, GoImportPath: "github.com/go-playground/validator/v10"}
	}
	g.P()
	d, err := declarationgen.New(g, p, surface, profile, nil)
	if err != nil {
		return err
	}
	if err := d.TypeDeclaration("error"); err != nil {
		return err
	}
	current, err := d.BeginBehavior("error", "message", "error_message")
	if err != nil {
		return err
	}
	d.Body(current, "\tif ", d.Receiver(current), " == nil {")
	d.Body(current, "\t\treturn \"\"")
	d.Body(current, "\t}")
	d.Body(current, "\tif ", d.Receiver(current), ".Field != \"\" {")
	d.Body(current, "\t\treturn ", d.Receiver(current), ".Message + \" (\" + ", d.Receiver(current), ".Field + \")\"")
	d.Body(current, "\t}")
	d.Body(current, "\treturn ", d.Receiver(current), ".Message")
	g.P("}")
	current, err = d.BeginBehavior("error", "unwrap", "unwrap")
	if err != nil {
		return err
	}
	d.Body(current, "\tif ", d.Receiver(current), " == nil {")
	d.Body(current, "\t\treturn nil")
	d.Body(current, "\t}")
	d.Body(current, "\treturn ", d.Receiver(current), ".", d.Field("error", "cause"))
	g.P("}")
	current, err = d.BeginBehavior("", "constructor", "new_error")
	if err != nil {
		return err
	}
	d.Body(current, "\towned, _ := ", d.Param(current, "info"), ".Clone()")
	d.Body(current, "\tif !owned.Kind.Valid() {")
	d.Body(current, "\t\towned.Kind = ErrorKindUnknown")
	d.Body(current, "\t}")
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
		d.Body(current, prefix, " ", errors("Is"), "(", d.Param(current, "cause"), ", ", context(member), ") {")
		d.Body(current, "owned.Kind = ", kind)
	}
	if len(policy.CauseRules) > 0 {
		g.P("}")
	}
	d.Body(current, "\treturn &Error{", d.Field("error", "ErrorInfo"), ": *owned, ", d.Field("error", "cause"), ": ", d.Param(current, "cause"), "}")
	g.P("}")
	current, err = d.BeginBehavior("", "context", "error_context")
	if err != nil {
		return err
	}
	d.Body(current, "\tif ", d.Param(current, "err"), " == nil {")
	d.Body(current, "\t\treturn nil")
	d.Body(current, "\t}")
	d.Body(current, "\tvar existing *Error")
	d.Body(current, "\tif ", errors("As"), "(", d.Param(current, "err"), ", &existing) && existing != nil {")
	if policy.ExistingError.ReuseMatchingContext {
		d.Body(current, "\t\tif (", d.Param(current, "provider"), "==\"\" || existing.Provider==", d.Param(current, "provider"), ") && (", d.Param(current, "operation"), "==\"\" || existing.Operation==", d.Param(current, "operation"), "){return ", d.Param(current, "err"), "}")
	}
	d.Body(current, "\t\tinfo, _ := existing.", d.Field("error", "ErrorInfo"), ".Clone()")
	d.Body(current, "\t\tif ", d.Param(current, "provider"), " != \"\" {")
	d.Body(current, "\t\t\tinfo.Provider = ", d.Param(current, "provider"))
	d.Body(current, "\t\t}")
	d.Body(current, "\t\tif ", d.Param(current, "operation"), " != \"\" {")
	d.Body(current, "\t\t\tinfo.Operation = ", d.Param(current, "operation"))
	d.Body(current, "\t\t}")
	d.Body(current, "\t\treturn &Error{", d.Field("error", "ErrorInfo"), ": *info, ", d.Field("error", "cause"), ": ", d.Param(current, "err"), "}")
	d.Body(current, "\t}")
	d.Body(current, "\treturn ", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindUnknown"), Provider: d.Param(current, "provider"), Operation: d.Param(current, "operation"), Message: d.Param(current, "err") + ".Error()", Cause: d.Param(current, "err")}))
	g.P("}")
	d.Body(current)
	current, err = d.BeginBehavior("", "validation_error", "validation_error")
	if err != nil {
		return err
	}
	d.Body(current, "\tfield := \"\"")
	d.Body(current, "\tvar fields ", validator("ValidationErrors"))
	d.Body(current, "\tif ", errors("As"), "(", d.Param(current, "err"), ", &fields) && len(fields) > 0 {")
	d.Body(current, "\t\tfield = validationPath(fields[0].StructNamespace())")
	d.Body(current, "\t}")
	d.Body(current, "\treturn ", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, "ErrorKindInvalidArgument"), Field: "field", Message: strconv.Quote("sandbox-kit: invalid configuration"), Cause: d.Param(current, "err")}))
	g.P("}")
	d.Body(current)
	current, err = d.BeginBehavior("", "validation_path", "validation_path")
	if err != nil {
		return err
	}
	d.Body(current, "\tparts := ", strings("Split"), "(", d.Param(current, "namespace"), ", \".\")")
	d.Body(current, "\tif len(parts) > 1 {")
	d.Body(current, "\t\tparts = parts[1:]")
	d.Body(current, "\t}")
	d.Body(current, "\tfor i, part := range parts {")
	d.Body(current, "\t\tbase, suffix := part, \"\"")
	d.Body(current, "\t\tif index := ", strings("IndexByte"), "(part, '['); index >= 0 {")
	d.Body(current, "\t\t\tbase, suffix = part[:index], part[index:]")
	d.Body(current, "\t\t}")
	d.Body(current, "\t\tif native, ok := validationFieldNames[base]; ok {")
	d.Body(current, "\t\t\tparts[i] = native + suffix")
	d.Body(current, "\t\t}")
	d.Body(current, "\t\tif native, ok := configValidationFieldNames[base]; ok {")
	d.Body(current, "\t\t\tparts[i] = native + suffix")
	d.Body(current, "\t\t}")
	d.Body(current, "\t}")
	d.Body(current, "\treturn ", strings("Join"), "(parts, \".\")")
	g.P("}")
	return nil
}

func emitCopyRuntime(p *protogen.Plugin, file *protogen.File, limits spec.Contracts, surface model.APIModule, profile spec.LanguageProfile) error {
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+profile.Outputs["copies"], file.GoImportPath)
	declarationgen.Banner(g, "specs/contracts.yaml",
		"Metadata helpers copy declarative values. They reject SDK objects, nonfinite numbers, and cycles.",
		"copied, err := cloneMetadata(info.ProviderMetadata)")
	g.P("package ", file.GoPackageName)
	fmt := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "fmt"} }
	math := func(member string) protogen.GoIdent { return protogen.GoIdent{GoName: member, GoImportPath: "math"} }
	g.P()
	d, err := declarationgen.New(g, p, surface, profile, nil)
	if err != nil {
		return err
	}
	current, err := d.BeginBehavior("", "clone_metadata", "clone_metadata")
	if err != nil {
		return err
	}
	d.Body(current, "\tif ", d.Param(current, "source"), " == nil {")
	d.Body(current, "\t\treturn nil, nil")
	d.Body(current, "\t}")
	d.Body(current, "nodes,bytes:=", limits.CopyMaxNodes, ",", limits.CopyMaxBytes)
	d.Body(current, "\tvalue, err := copyMetadataValue(", d.Param(current, "source"), ", 0, &nodes, &bytes)")
	d.Body(current, "\tif err != nil {")
	d.Body(current, "\t\treturn nil, err")
	d.Body(current, "\t}")
	d.Body(current, "\treturn value.(map[string]any), nil")
	g.P("}")
	current, err = d.BeginBehavior("", "copy_metadata", "copy_metadata")
	if err != nil {
		return err
	}
	d.Body(current, "\tif ", d.Param(current, "depth"), " > ", limits.CopyMaxDepth, " {")
	d.Body(current, "\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds maximum depth or contains a cycle\")")
	d.Body(current, "\t}")
	d.Body(current, "if *nodes<=0{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds node budget\")};*nodes--")
	d.Body(current, "\tswitch v := ", d.Param(current, "value"), ".(type) {")
	d.Body(current, "\tcase nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:")
	d.Body(current, "\t\treturn v, nil")
	d.Body(current, "case string: if len(v)>*bytes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds byte budget\")};*bytes-=len(v);return v,nil")
	d.Body(current, "\tcase float64:")
	d.Body(current, "\t\tif ", math("IsNaN"), "(v) || ", math("IsInf"), "(v, 0) {")
	d.Body(current, "\t\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: metadata requires finite numbers\")")
	d.Body(current, "\t\t}")
	d.Body(current, "\t\treturn v, nil")
	d.Body(current, "\tcase float32:")
	d.Body(current, "\t\tif ", math("IsNaN"), "(float64(v)) || ", math("IsInf"), "(float64(v), 0) {")
	d.Body(current, "\t\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: metadata requires finite numbers\")")
	d.Body(current, "\t\t}")
	d.Body(current, "\t\treturn v, nil")
	d.Body(current, "\tcase []byte:")
	d.Body(current, "\t\tif v == nil {")
	d.Body(current, "\t\t\treturn []byte(nil), nil")
	d.Body(current, "\t\t}")
	d.Body(current, "if len(v)>*bytes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds byte budget\")};*bytes-=len(v)")
	d.Body(current, "\t\treturn append([]byte{}, v...), nil")
	d.Body(current, "\tcase []string:")
	d.Body(current, "\t\tif v == nil {")
	d.Body(current, "\t\t\treturn []string(nil), nil")
	d.Body(current, "\t\t}")
	d.Body(current, "if len(v)>*nodes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds node budget\")};*nodes-=len(v)")
	d.Body(current, "for _,item:=range v{if len(item)>*bytes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds byte budget\")};*bytes-=len(item)}")
	d.Body(current, "\t\treturn append([]string{}, v...), nil")
	d.Body(current, "\tcase []any:")
	d.Body(current, "\t\tif v == nil {")
	d.Body(current, "\t\t\treturn []any(nil), nil")
	d.Body(current, "\t\t}")
	d.Body(current, "if len(v)>*nodes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds node budget\")}")
	d.Body(current, "\t\tout := make([]any, len(v))")
	d.Body(current, "\t\tfor i, item := range v {")
	d.Body(current, "\t\t\tcopied, err := copyMetadataValue(item, ", d.Param(current, "depth"), "+1, nodes, bytes)")
	d.Body(current, "\t\t\tif err != nil {")
	d.Body(current, "\t\t\t\treturn nil, err")
	d.Body(current, "\t\t\t}")
	d.Body(current, "\t\t\tout[i] = copied")
	d.Body(current, "\t\t}")
	d.Body(current, "\t\treturn out, nil")
	d.Body(current, "\tcase map[string]any:")
	d.Body(current, "\t\tif v == nil {")
	d.Body(current, "\t\t\treturn map[string]any(nil), nil")
	d.Body(current, "\t\t}")
	d.Body(current, "if len(v)>*nodes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds node budget\")}")
	d.Body(current, "for key:=range v {if len(key)>*bytes{return nil,", fmt("Errorf"), "(\"sandbox-kit: metadata exceeds byte budget\")};*bytes-=len(key)}")
	d.Body(current, "\t\tout := make(map[string]any, len(v))")
	d.Body(current, "\t\tfor key, item := range v {")
	d.Body(current, "\t\t\tcopied, err := copyMetadataValue(item, ", d.Param(current, "depth"), "+1, nodes, bytes)")
	d.Body(current, "\t\t\tif err != nil {")
	d.Body(current, "\t\t\t\treturn nil, err")
	d.Body(current, "\t\t\t}")
	d.Body(current, "\t\t\tout[key] = copied")
	d.Body(current, "\t\t}")
	d.Body(current, "\t\treturn out, nil")
	d.Body(current, "\tdefault:")
	d.Body(current, "\t\treturn nil, ", fmt("Errorf"), "(\"sandbox-kit: unsupported metadata value type %T\", ", d.Param(current, "value"), ")")
	d.Body(current, "\t}")
	g.P("}")
	return nil
}
