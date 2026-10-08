// Package validationgen translates portable validation rules into Go validators.
package validationgen

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
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
func Generate(p *protogen.Plugin, f *protogen.File, rules spec.Validation, root string) error {
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
	sort.Strings(names)
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".validation.gen.go", f.GoImportPath)
	g.P("// Code generated from YAML validation specifications. DO NOT EDIT.")
	g.P("package ", f.GoPackageName)
	validator := func(name string) protogen.GoIdent {
		return protogen.GoIdent{GoName: name, GoImportPath: "github.com/go-playground/validator/v10"}
	}
	g.P("func new", root, "Validator() *", validator("Validate"), " {")
	g.P("v:=", validator("New"), "(", validator("WithRequiredStructEnabled"), "())")
	g.P("_ = v.RegisterValidation(\"finite\",func(fl ", validator("FieldLevel"), ")bool{ value:=fl.Field().Float();return !", protogen.GoIdent{GoName: "IsNaN", GoImportPath: "math"}, "(value)&&!", protogen.GoIdent{GoName: "IsInf", GoImportPath: "math"}, "(value,0) })")
	g.P("_ = v.RegisterValidation(\"nonblank\",func(fl ", validator("FieldLevel"), ")bool{return ", protogen.GoIdent{GoName: "TrimSpace", GoImportPath: "strings"}, "(fl.Field().String())!=\"\"})")
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
			if rule.Format != "" && (rule.Format != "http_url" || fld.Desc.Kind() != protoreflect.StringKind) {
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
		if len(r.Constraints) == 0 {
			continue
		}
		g.P("v.RegisterStructValidation(func(sl ", validator("StructLevel"), "){x:=sl.Current().Interface().(", name, ");_ = x")
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
			case "limit":
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
			case "at_most_one", "exactly_one", "requires_any", "forbids":
				if len(c.Fields) == 0 {
					return fmt.Errorf("%s: %s needs fields", name, c.Op)
				}
				count := fmt.Sprintf("count%d", i)
				g.P(count, ":=0")
				for _, n := range c.Fields {
					fld, err := field(m, n)
					if err != nil {
						return err
					}
					g.P("if ", active(fld), "{", count, "++}")
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
			g.P("if (", condition, ") && (", invalid, "){sl.ReportError(x,", strconv.Quote(label), ",", strconv.Quote(label), ",", strconv.Quote(c.Op), ",\"\")}")
		}
		g.P("},", name, "{})")
	}
	g.P("return v }")
	g.P("// Validate", root, " applies the generated shared rules without filling defaults.")
	g.P("func Validate", root, "(request *", root, ")error{if request==nil{return nil};if err:=new", root, "Validator().Struct(request);err!=nil{return ", protogen.GoIdent{GoName: "Errorf", GoImportPath: "fmt"}, "(\"sandbox-kit: invalid configuration: %w\",err)};return nil}")
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
