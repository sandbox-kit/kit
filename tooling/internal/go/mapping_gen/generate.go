// Package mappinggen emits typed mappings from portable field and conversion rules.
package mappinggen

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const sandbox = "github.com/sandbox-kit/kit/sdks/go/sandbox"

func message(p *protogen.Plugin, name string) *protogen.Message {
	for _, f := range p.Files {
		for _, m := range f.Messages {
			if m.GoIdent.GoName == name && string(f.GoImportPath) == sandbox {
				return m
			}
		}
	}
	return nil
}
func member(p *protogen.Plugin, t spec.Type, name string) (string, *protogen.Field, error) {
	if native, ok := t.Bindings["go"]; ok {
		binding := native.Fields[name]
		if !token.IsIdentifier(binding) {
			return "", nil, fmt.Errorf("%s: missing/invalid SDK field binding %q", t.Name, name)
		}
		return binding, nil, nil
	}
	m := message(p, t.Name)
	if m == nil {
		return "", nil, fmt.Errorf("unknown shared type %s", t.Name)
	}
	for _, f := range m.Fields {
		if string(f.Desc.Name()) == name {
			return naming.FieldName(f), f, nil
		}
	}
	return "", nil, fmt.Errorf("%s: unknown protobuf field %s", t.Name, name)
}
func ident(t spec.Type) protogen.GoIdent {
	path, name := sandbox, t.Name
	if native, ok := t.Bindings["go"]; ok {
		path, name = native.Import, native.Name
	}
	return protogen.GoIdent{GoName: name, GoImportPath: protogen.GoImportPath(path)}
}
func Generate(p *protogen.Plugin, f *protogen.File, rules spec.Provider) error {
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".mappings.gen.go", f.GoImportPath)
	g.P("// Code generated from specs/providers/", rules.Provider, ".yaml. DO NOT EDIT.")
	g.P("package ", f.GoPackageName)
	seen := map[string]bool{}
	for _, group := range rules.Groups {
		if !token.IsIdentifier(group.Name) || seen[group.Name] {
			return fmt.Errorf("invalid/duplicate mapping group %q", group.Name)
		}
		seen[group.Name] = true
		request := group.Direction == "request"
		if !request && group.Direction != "response" {
			return fmt.Errorf("invalid mapping direction %q", group.Direction)
		}
		if request {
			g.P("func ", group.Name, "(source *", ident(group.Source), ") (", ident(group.Target), ",", ident(group.Source), ",error){")
		} else {
			g.P("func ", group.Name, "(source *", ident(group.Source), ") (", ident(group.Target), ",error){")
		}
		g.P("target:=", ident(group.Target), "{}")
		if request {
			g.P("remaining:=", ident(group.Source), "{};if source!=nil{remaining=*source}")
		}
		returns := "target,nil"
		if request {
			returns = "target,remaining,nil"
		}
		g.P("if source==nil{return ", returns, "}")
		conditions := []string{}
		for _, condition := range group.When {
			member, _, err := member(p, group.Source, condition.Field)
			if err != nil {
				return err
			}
			conditions = append(conditions, fmt.Sprintf("source.%s>%d", member, condition.GreaterThan))
		}
		if len(conditions) > 0 {
			g.P("if !(", strings.Join(conditions, " && "), "){return ", returns, "}")
		}

		destinations := map[string]bool{}
		for _, mapping := range group.Fields {
			if destinations[mapping.To] {
				return fmt.Errorf("%s: duplicate destination %s", group.Name, mapping.To)
			}
			destinations[mapping.To] = true
			from, sf, err := member(p, group.Source, mapping.From)
			if err != nil {
				return err
			}
			to, tf, err := member(p, group.Target, mapping.To)
			if err != nil {
				return err
			}
			if mapping.Transform == "policy_minutes" {
				if !request || sf == nil || sf.Message == nil || sf.Message.Desc.FullName() != "kit.sandbox.v1.AutomaticAction" || tf != nil || mapping.Policy == nil || (mapping.Policy.Disabled != "zero" && mapping.Policy.Disabled != "reject") || mapping.Cast != "" || mapping.Factor != 0 {
					return fmt.Errorf("invalid policy_minutes mapping %s", mapping.From)
				}
				label := "sandbox-kit " + rules.Provider + ": " + mapping.From
				fail := func(message string) {
					g.P("return target,remaining,", protogen.GoIdent{GoName: "Errorf", GoImportPath: "fmt"}, "(", strconv.Quote(label+": "+message), ")")
				}
				g.P("if policy:=source.", from, ";policy!=nil{switch policy.Mode{")
				g.P("case ", protogen.GoIdent{GoName: "PolicyModeDefault", GoImportPath: sandbox}, ":")
				g.P("case ", protogen.GoIdent{GoName: "PolicyModeDisabled", GoImportPath: sandbox}, ":")
				if mapping.Policy.Disabled == "zero" {
					g.P("target.", to, "=", protogen.GoIdent{GoName: "Value", GoImportPath: sandbox}, "(0)")
				} else {
					fail("explicit disabling is not representable; zero would trigger immediately")
				}
				g.P("case ", protogen.GoIdent{GoName: "PolicyModeAfter", GoImportPath: sandbox}, ":")
				g.P("if policy.After==nil{")
				fail("duration is required")
				g.P("}")
				minute := protogen.GoIdent{GoName: "Minute", GoImportPath: "time"}
				comparison := "<0"
				if !mapping.Policy.Immediate {
					comparison = "<=0"
				}
				g.P("duration:=*policy.After;if duration", comparison, " || duration%", minute, "!=0 || duration/", minute, ">2147483647{")
				fail("delay cannot be represented in whole minutes (zero may disable this action)")
				g.P("}")
				if mapping.Maximum != nil {
					g.P("if duration/", minute, ">", *mapping.Maximum, "{")
					fail("delay exceeds documented maximum")
					g.P("}")
				}
				g.P("value:=int(duration/", minute, ");target.", to, "=&value")
				g.P("default:")
				fail("unknown policy mode")
				g.P("}}")
				g.P("remaining.", from, "=nil")
				continue
			}
			if mapping.Policy != nil {
				return fmt.Errorf("policy rules require policy_minutes transform")
			}
			duration := sf != nil && sf.Message != nil && sf.Message.Desc.FullName() == "google.protobuf.Duration"
			optional := sf != nil && sf.Desc.HasPresence() && (sf.Desc.Kind() != protoreflect.MessageKind || duration)
			if optional {
				g.P("if source.", from, "!=nil{")
			}
			value := "source." + from
			if optional {
				value = "source.Get" + from + "()"
				if duration {
					value = "*source." + from
				}
			}
			g.P("{")
			g.P("value:=", value)
			invalid := ""
			comparisonValue := ""
			switch mapping.Transform {
			case "":
			case "whole":
				invalid = "float64(value)!=" + g.QualifiedGoIdent(protogen.GoIdent{GoName: "Trunc", GoImportPath: "math"}) + "(float64(value))"
			case "scaled_integer":
				if mapping.Factor == 0 || sf == nil || sf.Desc.Kind() != protoreflect.DoubleKind || mapping.Cast != "" {
					return fmt.Errorf("scaled_integer requires a shared double, positive factor, and no cast")
				}
				g.P("scaled:=", protogen.GoIdent{GoName: "Round", GoImportPath: "math"}, "(value*", mapping.Factor, ")")
				invalid = fmt.Sprintf("%s(value) || %s(value,0) || value<0 || value!=scaled/%d", g.QualifiedGoIdent(protogen.GoIdent{GoName: "IsNaN", GoImportPath: "math"}), g.QualifiedGoIdent(protogen.GoIdent{GoName: "IsInf", GoImportPath: "math"}), mapping.Factor)
				comparisonValue = "scaled"
			case "duration_seconds":
				if sf == nil || sf.Message == nil || sf.Message.Desc.FullName() != "google.protobuf.Duration" || mapping.Cast != "" {
					return fmt.Errorf("duration_seconds requires shared duration and no cast")
				}
				second := g.QualifiedGoIdent(protogen.GoIdent{GoName: "Second", GoImportPath: "time"})
				invalid = "value<=0 || value%" + second + "!=0"
				comparisonValue = "value/" + second
			case "divide_exactly":
				if mapping.Factor == 0 {
					return fmt.Errorf("%s: divide factor must be positive", group.Name)
				}
				invalid = fmt.Sprintf("value%%%d!=0", mapping.Factor)
			case "multiply":
				if mapping.Factor == 0 {
					return fmt.Errorf("%s: multiply factor must be positive", group.Name)
				}
				invalid = fmt.Sprintf("value<0 || uint64(value)>^uint64(0)/%d", mapping.Factor)
			default:
				return fmt.Errorf("unknown mapping transform %q", mapping.Transform)
			}
			if mapping.Transform == "divide_exactly" {
				value = fmt.Sprintf("value/%d", mapping.Factor)
			} else if mapping.Transform == "multiply" {
				value = fmt.Sprintf("uint64(value)*%d", mapping.Factor)
			} else {
				value = "value"
			}
			if mapping.Maximum != nil {
				if invalid != "" {
					invalid += " || "
				}
				boundValue := value
				if comparisonValue != "" {
					boundValue = comparisonValue
				}
				invalid += fmt.Sprintf("%s>%d", boundValue, *mapping.Maximum)
			}
			if invalid != "" {
				errreturn := "target,"
				if request {
					errreturn += "remaining,"
				}
				g.P("if ", invalid, "{return ", errreturn, protogen.GoIdent{GoName: "Errorf", GoImportPath: "fmt"}, "(", strconv.Quote("sandbox-kit "+rules.Provider+": "+mapping.From+" cannot be represented"), ")}")
			}
			if mapping.Transform == "scaled_integer" {
				// Compensate for binary rounding before SDKs truncate scaled values.
				g.P("if value*", mapping.Factor, "<scaled{value=", protogen.GoIdent{GoName: "Nextafter", GoImportPath: "math"}, "(value,", protogen.GoIdent{GoName: "Inf", GoImportPath: "math"}, "(1))}")
			}
			switch mapping.Cast {
			case "":
			case "integer":
				value = "int(" + value + ")"
			case "number":
				value = "float64(" + value + ")"
			case "unsigned":
				value = "uint64(" + value + ")"
			case "string":
				value = "string(" + value + ")"
			default:
				return fmt.Errorf("unknown cast %q", mapping.Cast)
			}
			if tf != nil && tf.Desc.HasPresence() && tf.Desc.Kind() != protoreflect.MessageKind {
				value = g.QualifiedGoIdent(protogen.GoIdent{GoName: "Value", GoImportPath: sandbox}) + "(" + value + ")"
			}
			// Each field has its own scope, so local names cannot collide.
			g.P("target.", to, "=", value)
			g.P("}")
			if optional {
				g.P("}")
			}
			if request {
				zero := "nil"
				if sf != nil && !sf.Desc.HasPresence() && !sf.Desc.IsList() && !sf.Desc.IsMap() {
					zero = "0"
					if sf.Desc.Kind() == protoreflect.StringKind {
						zero = `""`
					}
					if sf.Desc.Kind() == protoreflect.BoolKind {
						zero = "false"
					}
				}
				g.P("remaining.", from, "=", zero)
			}
		}
		if !request && group.Target.Name == "SandboxInfo" {
			g.P("target.Provider=", strconv.Quote(rules.Provider))
		}
		g.P("return ", returns, "}")
	}
	return nil
}
