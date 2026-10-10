package typesgen

import (
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"google.golang.org/protobuf/compiler/protogen"
	"sort"
	"strconv"
)

// GenerateErrorPaths derives semantic validator paths from schema field names.
func GenerateErrorPaths(p *protogen.Plugin, file *protogen.File) error {
	fields := map[string]string{}
	for _, f := range p.Files {
		if f.GoImportPath != file.GoImportPath {
			continue
		}
		for _, m := range f.Messages {
			for _, field := range m.Fields {
				fields[naming.FieldName(field)] = string(field.Desc.Name())
			}
		}
	}
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".errors.gen.go", file.GoImportPath)
	declarationgen.Banner(g, "",
		"Validation field maps turn Go struct names back into canonical paths such as resources.cpu_cores.",
		"path := validationPath(\"CreateOptions.Resources.CPUCores\")")
	g.P("package ", file.GoPackageName)
	name := "validationFieldNames"
	if file.Desc.Path() == "kit/sandbox/v1/client.proto" {
		name = "configValidationFieldNames"
	}
	g.P("var ", name, "=map[string]string{")
	keys := []string{}
	for key := range fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		g.P(strconv.Quote(key), ":", strconv.Quote(fields[key]), ",")
	}
	g.P("}")
	return nil
}
