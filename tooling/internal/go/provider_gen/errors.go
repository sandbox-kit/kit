package providergen

import (
	"fmt"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"go/token"
	"google.golang.org/protobuf/compiler/protogen"
	"strconv"
)

// GenerateErrors emits native matching for compiled semantic classifications.
func GenerateErrors(p *protogen.Plugin, file *protogen.File, compiled model.Provider) error {
	s := compiled.Rules
	binding := s.Errors["go"]
	if len(s.ErrorStatuses) > 0 && binding.Target.Import == "" {
		return fmt.Errorf("HTTP error rules require a native status binding")
	}
	core := protogen.GoImportPath("github.com/sandbox-kit/kit/sdks/go/sandbox")
	id := func(name string) protogen.GoIdent { return protogen.GoIdent{GoName: name, GoImportPath: core} }
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".errors.gen.go", file.GoImportPath)
	g.P("// Code generated from provider error bindings. DO NOT EDIT.")
	g.P("package ", file.GoPackageName)
	g.P("func mapProviderError(err error,operation string)error{if err==nil{return nil};var local *", id("Error"), ";if ", protogen.GoIdent{GoName: "As", GoImportPath: "errors"}, "(err,&local){return ", id("WithErrorContext"), "(err,", strconv.Quote(s.Provider), ",operation)}")
	g.P("kind:=", id("ErrorKindUnknown"), ";var statusCode *uint32;var providerCode,providerSource *string")
	kind := func(classification model.Kind) (protogen.GoIdent, error) {
		m, err := sharedMessage(p, "ErrorInfo")
		if err != nil {
			return protogen.GoIdent{}, err
		}
		for _, f := range m.Fields {
			if string(f.Desc.Name()) == "kind" {
				for _, v := range f.Enum.Values {
					if classification.Name == v.Desc.FullName() {
						return id(naming.EnumValueName(v)), nil
					}
				}
			}
		}
		return protogen.GoIdent{}, fmt.Errorf("missing Go enum metadata for %s", classification.Name)
	}
	for i, t := range binding.Types {
		if !validBinding(spec.Binding{Import: t.Import, Name: t.Name}) {
			return fmt.Errorf("invalid error type")
		}
		var classification model.Kind
		for _, rule := range compiled.Errors.Native {
			if rule.Rule == t.Rule {
				classification = rule.Kind
				break
			}
		}
		k, err := kind(classification)
		if err != nil {
			return err
		}
		typ := protogen.GoIdent{GoName: t.Name, GoImportPath: protogen.GoImportPath(t.Import)}
		if t.Pointer {
			g.P("if kind==", id("ErrorKindUnknown"), "{var native *", typ, ";if errors.As(err,&native)&&native!=nil{kind=", k, "}}")
		} else {
			g.P("if kind==", id("ErrorKindUnknown"), "{var value", i, " ", typ, ";var pointer", i, " *", typ, ";if errors.As(err,&value", i, ")||errors.As(err,&pointer", i, "){kind=", k, "}}")
		}
	}
	if binding.Target.Import != "" {
		if !validBinding(binding.Target) {
			return fmt.Errorf("invalid native error target")
		}
		for _, f := range []string{"status", "code", "source"} {
			if !token.IsIdentifier(binding.Target.Fields[f]) {
				return fmt.Errorf("missing native error member %s", f)
			}
		}
		g.P("{var native *", protogen.GoIdent{GoName: binding.Target.Name, GoImportPath: protogen.GoImportPath(binding.Target.Import)}, ";if errors.As(err,&native){")
		g.P("if native.", binding.Target.Fields["status"], ">0{statusCode=", id("Value"), "(uint32(native.", binding.Target.Fields["status"], "))};if native.", binding.Target.Fields["code"], "!=\"\"{providerCode=", id("Value"), "(native.", binding.Target.Fields["code"], ")};if native.", binding.Target.Fields["source"], "!=\"\"{providerSource=", id("Value"), "(native.", binding.Target.Fields["source"], ")}")
		g.P("if kind==", id("ErrorKindUnknown"), "{switch native.", binding.Target.Fields["status"], "{")
		for _, rule := range compiled.Errors.HTTP {
			code := rule.Code
			k, err := kind(rule.Kind)
			if err != nil {
				return err
			}
			g.P("case ", code, ":kind=", k)
		}
		g.P("}}}}")
	}
	if len(compiled.Errors.GRPC) > 0 {
		g.P("if kind==", id("ErrorKindUnknown"), "{switch ", protogen.GoIdent{GoName: "Code", GoImportPath: "google.golang.org/grpc/status"}, "(err){")
		for _, rule := range compiled.Errors.GRPC {
			key := rule.Code
			k, err := kind(rule.Kind)
			if err != nil {
				return err
			}
			g.P("case ", protogen.GoIdent{GoName: "Code", GoImportPath: "google.golang.org/grpc/codes"}, "(", key, "):kind=", k)
		}
		g.P("}}")
	}
	g.P("return ", errorgen.Expression(g, errorgen.Details{Kind: "kind", Provider: strconv.Quote(s.Provider), Operation: "operation", Message: "err.Error()", StatusCode: "statusCode", ProviderCode: "providerCode", ProviderSource: "providerSource", Cause: "err"}), "}")
	return nil
}
