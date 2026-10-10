// Package errorgen provides one Go expression builder for shared errors.
package errorgen

import (
	"google.golang.org/protobuf/compiler/protogen"
	"strings"
)

const Core protogen.GoImportPath = "github.com/sandbox-kit/kit/sdks/go/sandbox"

// Details contains Go expressions resolved from semantic generator inputs.
type Details struct {
	Kind, Provider, Operation, Field, Message string
	StatusCode, ProviderCode, ProviderSource  string
	Cause                                     string
}

// Kind qualifies a schema-derived Go error-kind identifier for the output file.
func Kind(g *protogen.GeneratedFile, name string) string {
	return g.QualifiedGoIdent(protogen.GoIdent{GoName: name, GoImportPath: Core})
}

// Expression constructs shared error details and a separate native cause.
func Expression(g *protogen.GeneratedFile, d Details) string {
	fields := []string{}
	for _, item := range []struct{ name, value string }{
		{"Kind", d.Kind}, {"Provider", d.Provider}, {"Operation", d.Operation}, {"Field", d.Field}, {"Message", d.Message},
		{"StatusCode", d.StatusCode}, {"ProviderCode", d.ProviderCode}, {"ProviderSource", d.ProviderSource},
	} {
		if item.value != "" {
			fields = append(fields, item.name+":"+item.value)
		}
	}
	cause := d.Cause
	if cause == "" {
		cause = "nil"
	}
	return g.QualifiedGoIdent(protogen.GoIdent{GoName: "NewError", GoImportPath: Core}) + "(" + g.QualifiedGoIdent(protogen.GoIdent{GoName: "ErrorInfo", GoImportPath: Core}) + "{" + strings.Join(fields, ",") + "}," + cause + ")"
}
