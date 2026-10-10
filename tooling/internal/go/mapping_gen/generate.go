// Package mappinggen emits typed mappings from portable field and conversion rules.
package mappinggen

import (
	"fmt"
	"go/token"
	"strconv"
	"strings"

	declarationgen "github.com/sandbox-kit/kit/tooling/internal/go/declaration_gen"
	errorgen "github.com/sandbox-kit/kit/tooling/internal/go/error_gen"
	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/go/schema"
	"github.com/sandbox-kit/kit/tooling/internal/model"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const sandbox = "github.com/sandbox-kit/kit/sdks/go/sandbox"

func message(p *protogen.Plugin, name string) *protogen.Message {
	m, _ := schema.Message(p, name)
	return m
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

// Generate emits a validated portable mapping plan using Go SDK bindings.
func Generate(p *protogen.Plugin, f *protogen.File, compiled model.Provider, templates model.Templates, profile spec.LanguageProfile) error {
	rules := compiled.Rules
	g := p.NewGeneratedFile(f.GeneratedFilenamePrefix+".mappings.gen.go", f.GoImportPath)
	g.P("// Code generated from specs/providers/", rules.Provider, ".yaml. DO NOT EDIT.")
	g.P("package ", f.GoPackageName)
	seen := map[string]bool{}
	for _, resolvedGroup := range compiled.Groups {
		group := resolvedGroup.Rule
		if !token.IsIdentifier(group.Name) || seen[group.Name] {
			return fmt.Errorf("invalid/duplicate mapping group %q", group.Name)
		}
		seen[group.Name] = true
		request := group.Direction == "request"
		if !request && group.Direction != "response" {
			return fmt.Errorf("invalid mapping direction %q", group.Direction)
		}
		template := "response_mapping"
		if request {
			template = "request_mapping"
		}
		d, c, err := (declarationgen.TemplateEmitter{G: g, Plugin: p, Templates: templates, Profile: profile}).Begin(template, declarationgen.TemplateBindings{Names: map[string]string{"mapping": group.Name}, Types: map[string]spec.TypeRef{"source": {Ref: "native_source"}, "target": {Ref: "native_target"}}, Externals: map[string]spec.ExternalType{"native_source": declarationgen.NamedBinding(ident(group.Source)), "native_target": declarationgen.NamedBinding(ident(group.Target))}})
		if err != nil {
			return err
		}
		d.Body(c, "target:=", ident(group.Target), "{}")
		if request {
			d.Body(c, "remaining:=", ident(group.Source), "{};if ", d.Param(c, "source"), "!=nil{remaining=*", d.Param(c, "source"), "}")
		}
		returns := "target,nil"
		if request {
			returns = "target,remaining,nil"
		}
		d.Body(c, "if ", d.Param(c, "source"), "==nil{return ", returns, "}")
		conditions := []string{}
		for _, condition := range group.When {
			member, _, err := member(p, group.Source, condition.Field)
			if err != nil {
				return err
			}
			conditions = append(conditions, fmt.Sprintf("%s.%s>%d", d.Param(c, "source"), member, condition.GreaterThan))
		}
		if len(conditions) > 0 {
			d.Body(c, "if !(", strings.Join(conditions, " && "), "){return ", returns, "}")
		}

		destinations := map[string]bool{}
		if group.Origin != "" {
			if group.Direction != "response" || group.Target.Name != "SandboxInfo" || (group.Origin != "provider" && group.Origin != "request") {
				return fmt.Errorf("invalid response origin")
			}
			d.Body(c, "target.Origins=map[string]", protogen.GoIdent{GoName: "ValueOrigin", GoImportPath: sandbox}, "{}")
		}
		for _, resolvedMapping := range resolvedGroup.Fields {
			mapping := resolvedMapping.Rule
			errorPath := resolvedMapping.ErrorPath
			if destinations[mapping.To] {
				return fmt.Errorf("%s: duplicate destination %s", group.Name, mapping.To)
			}
			destinations[mapping.To] = true
			shared, err := schema.Field(p, resolvedMapping.Shared.Fields[len(resolvedMapping.Shared.Fields)-1])
			if err != nil {
				return err
			}
			var from, to string
			var sf, tf *protogen.Field
			if request {
				from, sf = naming.FieldName(shared), shared
				to, tf, err = member(p, group.Target, mapping.To)
			} else {
				to, tf = naming.FieldName(shared), shared
				from, sf, err = member(p, group.Source, mapping.From)
			}
			if err != nil {
				return err
			}
			if mapping.Transform == "policy_minutes" {
				if !request || sf == nil || sf.Message == nil || sf.Message.Desc.FullName() != "kit.sandbox.v1.AutomaticAction" || tf != nil || mapping.Policy == nil || (mapping.Policy.Disabled != "zero" && mapping.Policy.Disabled != "reject") || mapping.Cast != "" || mapping.Factor != 0 {
					return fmt.Errorf("invalid policy_minutes mapping %s", mapping.From)
				}
				label := "sandbox-kit " + rules.Provider + ": " + mapping.From
				fieldPath, err := schema.Literal(g, p, "CreateOptions", "lifetime."+mapping.From, templates.FieldPaths)
				if err != nil {
					return err
				}
				fail := func(message, kind string) {
					d.Body(c, "return target,remaining,", errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, kind), Provider: strconv.Quote(rules.Provider), Operation: strconv.Quote("create"), Field: fieldPath, Message: strconv.Quote(label + ": " + message)}))
				}
				d.Body(c, "if policy:=", d.Param(c, "source"), ".", from, ";policy!=nil{switch policy.Mode{")
				d.Body(c, "case ", protogen.GoIdent{GoName: "PolicyModeDefault", GoImportPath: sandbox}, ":")
				d.Body(c, "case ", protogen.GoIdent{GoName: "PolicyModeDisabled", GoImportPath: sandbox}, ":")
				if mapping.Policy.Disabled == "zero" {
					d.Body(c, "target.", to, "=", protogen.GoIdent{GoName: "Value", GoImportPath: sandbox}, "(0)")
				} else {
					fail("explicit disabling is not representable; native zero has a different meaning", "ErrorKindUnsupported")
				}
				d.Body(c, "case ", protogen.GoIdent{GoName: "PolicyModeAfter", GoImportPath: sandbox}, ":")
				d.Body(c, "if policy.After==nil{")
				fail("duration is required", "ErrorKindInvalidArgument")
				d.Body(c, "}")
				minute := protogen.GoIdent{GoName: "Minute", GoImportPath: "time"}
				comparison := "<0"
				if !mapping.Policy.Immediate {
					d.Body(c, "if *policy.After==0{")
					fail("immediate action is not representable", "ErrorKindUnsupported")
					d.Body(c, "}")
				}
				d.Body(c, "duration:=*policy.After;if duration", comparison, " || duration%", minute, "!=0 || duration/", minute, ">2147483647{")
				fail("delay cannot be represented in whole minutes", "ErrorKindInvalidArgument")
				d.Body(c, "}")
				if mapping.Maximum != nil {
					d.Body(c, "if duration/", minute, ">", *mapping.Maximum, "{")
					fail("delay exceeds documented maximum", "ErrorKindInvalidArgument")
					d.Body(c, "}")
				}
				d.Body(c, "value:=int(duration/", minute, ");target.", to, "=&value")
				d.Body(c, "default:")
				fail("unknown policy mode", "ErrorKindInvalidArgument")
				d.Body(c, "}}")
				d.Body(c, "remaining.", from, "=nil")
				continue
			}
			if mapping.Policy != nil {
				return fmt.Errorf("policy rules require policy_minutes transform")
			}
			duration := sf != nil && sf.Message != nil && sf.Message.Desc.FullName() == "google.protobuf.Duration"
			optional := sf != nil && sf.Desc.HasPresence() && (sf.Desc.Kind() != protoreflect.MessageKind || duration)
			if optional {
				d.Body(c, "if ", d.Param(c, "source"), ".", from, "!=nil{")
			}
			value := d.Param(c, "source") + "." + from
			if optional {
				value = d.Param(c, "source") + ".Get" + from + "()"
				if duration {
					value = "*" + d.Param(c, "source") + "." + from
				}
			}
			d.Body(c, "{")
			d.Body(c, "value:=", value)
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
				d.Body(c, "scaled:=", protogen.GoIdent{GoName: "Round", GoImportPath: "math"}, "(value*", mapping.Factor, ")")
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
				kind := "ErrorKindInvalidArgument"
				if !request {
					kind = "ErrorKindInvalidResponse"
				}
				field := strconv.Quote(errorPath)
				if request && message(p, "CreateOptions") != nil {
					resolvedFields, err := schema.Resolve(p, "CreateOptions", errorPath)
					if err != nil {
						return err
					}
					if sf != nil && resolvedFields[len(resolvedFields)-1].Desc.FullName() != sf.Desc.FullName() {
						return fmt.Errorf("error path does not match mapping source %s", mapping.From)
					}
					resolved, err := schema.Literal(g, p, "CreateOptions", errorPath, templates.FieldPaths)
					if err != nil {
						return err
					}
					field = resolved
				}
				d.Body(c, "if ", invalid, "{return ", errreturn, errorgen.Expression(g, errorgen.Details{Kind: errorgen.Kind(g, kind), Provider: strconv.Quote(rules.Provider), Operation: strconv.Quote("create"), Field: field, Message: strconv.Quote("sandbox-kit " + rules.Provider + ": " + mapping.From + " cannot be represented")}), "}")
			}
			if mapping.Transform == "scaled_integer" {
				// Compensate for binary rounding before SDKs truncate scaled values.
				d.Body(c, "if value*", mapping.Factor, "<scaled{value=", protogen.GoIdent{GoName: "Nextafter", GoImportPath: "math"}, "(value,", protogen.GoIdent{GoName: "Inf", GoImportPath: "math"}, "(1))}")
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
			d.Body(c, "target.", to, "=", value)
			if group.Origin != "" {
				origin := "ValueOriginProvider"
				if group.Origin == "request" {
					origin = "ValueOriginRequest"
				}
				guard := tf != nil && (tf.Desc.IsMap() || tf.Desc.IsList() || tf.Desc.HasPresence())
				if guard {
					d.Body(c, "if target.", to, "!=nil{")
				}
				d.Body(c, "target.Origins[", strconv.Quote(mapping.To), "]=", protogen.GoIdent{GoName: origin, GoImportPath: sandbox})
				if guard {
					d.Body(c, "}")
				}
			}
			d.Body(c, "}")
			if optional {
				d.Body(c, "}")
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
				d.Body(c, "remaining.", from, "=", zero)
			}
		}
		if !request && group.Target.Name == "SandboxInfo" {
			d.Body(c, "target.Provider=", strconv.Quote(rules.Provider))
		}
		d.Body(c, "return ", returns, "}")
	}
	return nil
}
