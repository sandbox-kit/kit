// Package validationgen translates portable validation rules into Go validators.
package validationgen

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func field(m *protogen.Message, name string) (*protogen.Field, error) {
	for _, f := range m.Fields {
		if string(f.Desc.Name()) == name {
			return f, nil
		}
	}
	return nil, fmt.Errorf("%s: unknown validation field %q", m.Desc.FullName(), name)
}
func active(f *protogen.Field) string {
	path := "x." + naming.FieldName(f)
	if f.Desc.IsList() || f.Desc.IsMap() {
		return "len(" + path + ")>0"
	}
	if f.Desc.Kind() == protoreflect.BoolKind {
		return "x.Get" + naming.FieldName(f) + "()"
	}
	if f.Desc.HasPresence() {
		return path + "!=nil"
	}
	if f.Desc.Kind() == protoreflect.StringKind {
		return path + `!=""`
	}
	return path + "!=0"
}
func Generate(p *protogen.Plugin, f *protogen.File, rules spec.Validation, root string, templates model.Templates, profile spec.LanguageProfile) error {
	if len(rules.Messages) == 0 {
		return nil
	}
	messages := map[string]*protogen.Message{}
	for _, m := range f.Messages {
		messages[m.GoIdent.GoName] = m
	}
	names := make([]string, 0, len(rules.Messages))
	for name := range rules.Messages {
		names = append(names, name)
	}
	for _, m := range f.Messages {
		if _, ok := rules.Messages[m.GoIdent.GoName]; ok {
			continue
		}
		for _, fld := range m.Fields {
			if fld.Desc.Kind() == protoreflect.EnumKind {
				names = append(names, m.GoIdent.GoName)
				break
			}
		}
	}
	sort.Strings(names)
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".validation.gen.go", f.GoImportPath)
	declarationgen.Banner(g, "source: "+string(f.Desc.Path()),
		"Validators reject values the shared contract cannot represent. They run before a provider call.",
		"if err := Validate"+root+"(value); err != nil { return err }")
	g.P("package ", f.GoPackageName)
	validator := func(name string) protogen.GoIdent {
		return protogen.GoIdent{GoName: name, GoImportPath: "github.com/go-playground/validator/v10"}
	}
	emitter := declarationgen.TemplateEmitter{G: g, Plugin: p, Templates: templates, Profile: profile}
	d, callable, err := emitter.Begin("validator_factory", declarationgen.TemplateBindings{Names: map[string]string{"subject": root}, Types: map[string]spec.TypeRef{"validator": {Ref: "native_validator"}}, Externals: map[string]spec.ExternalType{"native_validator": declarationgen.NamedBinding(validator("Validate"))}})
	if err != nil {
		return err
	}
	d.Body(callable, "v:=", validator("New"), "(", validator("WithRequiredStructEnabled"), "())")
	d.Body(callable, "_ = v.RegisterValidation(\"finite\",func(fl ", validator("FieldLevel"), ")bool{ value:=fl.Field().Float();return !", protogen.GoIdent{GoName: "IsNaN", GoImportPath: "math"}, "(value)&&!", protogen.GoIdent{GoName: "IsInf", GoImportPath: "math"}, "(value,0) })")
	d.Body(callable, "_ = v.RegisterValidation(\"nonblank\",func(fl ", validator("FieldLevel"), ")bool{return ", protogen.GoIdent{GoName: "TrimSpace", GoImportPath: "strings"}, "(fl.Field().String())!=\"\"})")
	d.Body(callable, "_ = v.RegisterValidation(\"secure_endpoint\",func(fl ", validator("FieldLevel"), ")bool{ u,err:=", protogen.GoIdent{GoName: "Parse", GoImportPath: "net/url"}, "(fl.Field().String());if err!=nil || u.Hostname()==\"\" || u.User!=nil || u.Fragment!=\"\"{return false};if u.Scheme==\"https\"{return true};if u.Scheme!=\"http\"{return false};ip,err:=", protogen.GoIdent{GoName: "ParseAddr", GoImportPath: "net/netip"}, "(u.Hostname());return u.Hostname()==\"localhost\" || (err==nil && ip.IsLoopback()) })")
	d.Body(callable, "_ = v.RegisterValidation(\"image_reference\",func(fl ", validator("FieldLevel"), ")bool{return len(fl.Field().String())<=512 && ", protogen.GoIdent{GoName: "MustCompile", GoImportPath: "regexp"}, "(", strconv.Quote(`^(?:[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?(?::[0-9]+)?/)?[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*(?:/[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*)*(?::[A-Za-z0-9_][A-Za-z0-9_.-]{0,127})?(?:@sha256:[a-fA-F0-9]{64})?$`), ").MatchString(fl.Field().String())})")
	for _, name := range names {
		m := messages[name]
		if m == nil {
			return fmt.Errorf("unknown validation message %q", name)
		}
		r := rules.Messages[name]
		for n, rule := range r.Fields {
			fld, err := field(m, n)
			if err != nil {
				return err
			}
			if rule.Each && !fld.Desc.IsList() {
				return fmt.Errorf("%s.%s: each requires a list", name, n)
			}
			if rule.Finite && fld.Desc.Kind() != protoreflect.DoubleKind {
				return fmt.Errorf("%s.%s: finite requires double", name, n)
			}
			if rule.Nonblank && fld.Desc.Kind() != protoreflect.StringKind {
				return fmt.Errorf("%s.%s: nonblank requires string", name, n)
			}
			if rule.Format != "" && ((rule.Format != "http_url" && rule.Format != "secure_endpoint" && rule.Format != "image_reference") || fld.Desc.Kind() != protoreflect.StringKind) {
				return fmt.Errorf("%s.%s: unsupported field format %q", name, n, rule.Format)
			}
			if fld.Desc.IsList() && !rule.Each && (rule.Nonblank || rule.Finite || rule.Minimum != nil || rule.Maximum != nil || rule.ExclusiveMinimum != nil) {
				return fmt.Errorf("%s.%s: element rules require each", name, n)
			}
			if rule.MinItems != nil && (!fld.Desc.IsList() || *rule.MinItems < 0) {
				return fmt.Errorf("%s.%s: min_items requires a list and nonnegative bound", name, n)
			}
			for _, bound := range []*float64{rule.Minimum, rule.ExclusiveMinimum, rule.Maximum} {
				if bound == nil {
					continue
				}
				if math.IsNaN(*bound) || math.IsInf(*bound, 0) {
					return fmt.Errorf("%s.%s: bounds must be finite", name, n)
				}
				if !numeric(fld) {
					return fmt.Errorf("%s.%s: numeric bound on nonnumeric field", name, n)
				}
			}
			if rule.Minimum != nil && rule.ExclusiveMinimum != nil {
				return fmt.Errorf("%s.%s: choose one minimum", name, n)
			}
		}
		hasEnums := false
		for _, fld := range m.Fields {
			if fld.Desc.Kind() == protoreflect.EnumKind {
				hasEnums = true
			}
		}
		if len(r.Constraints) == 0 && !hasEnums {
			continue
		}
		d.Body(callable, "v.RegisterStructValidation(func(sl ", validator("StructLevel"), "){x:=sl.Current().Interface().(", name, ");_ = x")
		for _, fld := range m.Fields {
			if fld.Desc.Kind() != protoreflect.EnumKind {
				continue
			}
			if fld.Desc.IsList() {
				return fmt.Errorf("repeated enum validation not implemented")
			}
			d.Body(callable, "switch x.Get", naming.FieldName(fld), "(){")
			values := []string{}
			seen := map[int32]bool{}
			for _, v := range fld.Enum.Values {
				n := int32(v.Desc.Number())
				if !seen[n] {
					values = append(values, strconv.Itoa(int(n)))
					seen[n] = true
				}
			}
			d.Body(callable, "case ", strings.Join(values, ","), ": default: sl.ReportError(x,", strconv.Quote(naming.FieldName(fld)), ",", strconv.Quote(naming.FieldName(fld)), ",\"known_enum\",\"\")}")
		}
		for i, c := range r.Constraints {
			condition := "true"
			if c.When != nil {
				fld, err := field(m, c.When.Field)
				if err != nil {
					return err
				}
				if (c.When.Equals == nil) == (c.When.NotEquals == nil) {
					return fmt.Errorf("%s: condition needs exactly one comparison", name)
				}
				if !numeric(fld) || fld.Desc.Kind() == protoreflect.MessageKind || fld.Desc.IsList() || fld.Desc.IsMap() {
					return fmt.Errorf("%s: condition requires numeric scalar", name)
				}
				op := "=="
				value := c.When.Equals
				if value == nil {
					value = c.When.NotEquals
					op = "!="
				}
				condition = "x.Get" + naming.FieldName(fld) + "()" + op + strconv.Itoa(*value)
			}
			invalid := ""
			switch c.Op {
			case "limit", "requires_positive":
				limit, err := field(m, c.Field)
				if err != nil {
					return err
				}
				other, err := field(m, c.Other)
				if err != nil {
					return err
				}
				if !numeric(limit) || !numeric(other) || limit.Desc.Kind() == protoreflect.MessageKind || limit.Desc.Kind() != other.Desc.Kind() || limit.Desc.IsList() || other.Desc.IsList() {
					return fmt.Errorf("%s: limit requires compatible numeric scalars", name)
				}
				invalid = "x.Get" + naming.FieldName(limit) + "()>0 && x.Get" + naming.FieldName(limit) + "()<x.Get" + naming.FieldName(other) + "()"
				if c.Op == "requires_positive" {
					invalid = "x.Get" + naming.FieldName(limit) + "()>0 && x.Get" + naming.FieldName(other) + "()==0"
				}
			case "at_most_one", "exactly_one", "requires_any", "forbids":
				if len(c.Fields) == 0 {
					return fmt.Errorf("%s: %s needs fields", name, c.Op)
				}
				count := fmt.Sprintf("count%d", i)
				d.Body(callable, count, ":=0")
				for _, n := range c.Fields {
					fld, err := field(m, n)
					if err != nil {
						return err
					}
					d.Body(callable, "if ", active(fld), "{", count, "++}")
				}
				switch c.Op {
				case "at_most_one":
					invalid = count + ">1"
				case "exactly_one":
					invalid = count + "!=1"
				case "requires_any":
					invalid = count + "==0"
				case "forbids":
					invalid = count + ">0"
				}
			default:
				return fmt.Errorf("%s: unknown constraint %q", name, c.Op)
			}
			label := c.Field
			if label == "" {
				label = strings.Join(c.Fields, ",")
			}
			d.Body(callable, "if (", condition, ") && (", invalid, "){sl.ReportError(x,", strconv.Quote(label), ",", strconv.Quote(label), ",", strconv.Quote(c.Op), ",\"\")}")
		}
		d.Body(callable, "},", name, "{})")
	}
	d.Body(callable, "return v }")
	d.Body(callable, "// Validate", root, " applies the generated shared rules without filling defaults.")
	d, callable, err = emitter.Begin("validation", declarationgen.TemplateBindings{Names: map[string]string{"subject": root}, Types: map[string]spec.TypeRef{"subject": {Ref: "schema." + root}}})
	if err != nil {
		return err
	}
	factory, err := templates.Instantiate("validator_factory", map[string]string{"subject": root})
	if err != nil {
		return err
	}
	d.Body(callable, "if ", d.Param(callable, "request"), "==nil{return nil};if err:=", factory.Declaration.Name, "().Struct(", d.Param(callable, "request"), ");err!=nil{return validationError(err)};return nil")
	g.P("}")
	return nil
}

func numeric(field *protogen.Field) bool {
	switch field.Desc.Kind() {
	case protoreflect.DoubleKind, protoreflect.Uint32Kind, protoreflect.Uint64Kind, protoreflect.EnumKind:
		return true
	case protoreflect.MessageKind:
		return field.Message.Desc.FullName() == "google.protobuf.Duration"
	default:
		return false
	}
}
