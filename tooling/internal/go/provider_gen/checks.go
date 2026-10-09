package providergen

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/go/naming"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// GenerateChecks emits validation over schema paths; account quotas remain remote.
func GenerateChecks(p *protogen.Plugin, file *protogen.File, s spec.Provider) error {
	g := p.NewGeneratedFile(file.GeneratedFilenamePrefix+".validation.gen.go", file.GoImportPath)
	g.P("// Code generated from provider validation specs. DO NOT EDIT.")
	g.P("package ", file.GoPackageName)
	formats := map[string]bool{}
	for _, check := range s.Checks {
		formats[check.Format] = true
	}
	if formats["env_name"] {
		g.P("var providerEnvNamePattern=", protogen.GoIdent{GoName: "MustCompile", GoImportPath: "regexp"}, "(", strconv.Quote("^[A-Za-z_][A-Za-z0-9_]*$"), ")")
	}
	if formats["domain"] {
		g.P("var providerDomainPattern=", protogen.GoIdent{GoName: "MustCompile", GoImportPath: "regexp"}, "(", strconv.Quote(`(?i)^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$`), ")")
	}
	g.P("func validateProviderCreate(request *", protogen.GoIdent{GoName: "CreateOptions", GoImportPath: "github.com/sandbox-kit/kit/sdks/go/sandbox"}, ")error{if request==nil{return nil}")
	for _, check := range s.Checks {
		if len(check.ExclusivePolicies) > 0 || check.ForbidPolicyWith != "" {
			g.P("if lifetime:=request.Lifetime;lifetime!=nil{")
			policy := func(name string) (string, error) {
				fld, err := sharedField(p, "LifetimePolicy", name)
				if err != nil {
					return "", err
				}
				if fld.Message == nil || fld.Message.Desc.FullName() != "kit.sandbox.v1.AutomaticAction" {
					return "", fmt.Errorf("relationship requires policy")
				}
				n := "lifetime." + naming.FieldName(fld)
				return n + "!=nil && " + n + ".Mode==2 && " + n + ".After!=nil && *" + n + ".After>0", nil
			}
			invalid := ""
			if len(check.ExclusivePolicies) > 0 {
				if len(check.ExclusivePolicies) != 2 {
					return fmt.Errorf("exclusive policies requires two fields")
				}
				a, err := policy(check.ExclusivePolicies[0])
				if err != nil {
					return err
				}
				b, err := policy(check.ExclusivePolicies[1])
				if err != nil {
					return err
				}
				invalid = "(" + a + ") && (" + b + ")"
			} else {
				a, err := policy(check.Path)
				if err != nil {
					return err
				}
				b := ""
				switch check.ForbidPolicyWith {
				case "ephemeral":
					b = "lifetime.GetEphemeral()"
				case "stopped_delete":
					b = "lifetime.StoppedDelete!=nil && lifetime.StoppedDelete.Mode==2 && lifetime.StoppedDelete.After!=nil && *lifetime.StoppedDelete.After==0"
				default:
					return fmt.Errorf("unsupported policy relationship")
				}
				invalid = "(" + a + ") && (" + b + ")"
			}
			g.P("if ", invalid, "{return ", protogen.GoIdent{GoName: "Errorf", GoImportPath: "fmt"}, "(\"sandbox-kit: incompatible lifetime policies\")}}")
			continue
		}
		message := "CreateOptions"
		path := "request"
		guards := 0
		parts := strings.Split(check.Path, ".")
		var fld *protogen.Field
		for i, part := range parts {
			var err error
			fld, err = sharedField(p, message, part)
			if err != nil {
				return err
			}
			path += "." + naming.FieldName(fld)
			if fld.Desc.IsMap() && i < len(parts)-1 {
				return fmt.Errorf("provider checks require scalar paths")
			}
			if i < len(parts)-1 {
				if fld.Message == nil {
					return fmt.Errorf("invalid check path %s", check.Path)
				}
				if fld.Desc.IsList() {
					g.P("for _,item:=range ", path, "{if item==nil{return ", protogen.GoIdent{GoName: "Errorf", GoImportPath: "fmt"}, "(\"sandbox-kit: nil configuration entry\")}")
					path = "item"
				} else {
					g.P("if ", path, "!=nil{")
				}
				guards++
				message = string(fld.Message.Desc.FullName())
			}
		}
		value := path
		if fld.Desc.IsMap() {
			if check.Format != "env_name" {
				return fmt.Errorf("map check requires env_name keys")
			}
			g.P("for value:=range ", path, "{")
			guards++
			value = "value"
		} else if fld.Desc.IsList() {
			if fld.Desc.Kind() != protoreflect.StringKind {
				return fmt.Errorf("provider collection checks require strings")
			}
			g.P("for _,value:=range ", path, "{")
			guards++
			value = "value"
		} else if fld.Desc.HasPresence() {
			g.P("if ", path, "!=nil{")
			guards++
			value = "*" + path
		}
		invalid := []string{}
		if len(check.Allowed) > 0 {
			if fld.Desc.Kind() != protoreflect.StringKind {
				return fmt.Errorf("allowed requires string")
			}
			for _, v := range check.Allowed {
				invalid = append(invalid, value+"!="+strconv.Quote(v))
			}
		}
		for _, bound := range []struct {
			n  *float64
			op string
		}{{check.Minimum, "<"}, {check.Maximum, ">"}} {
			if bound.n == nil {
				continue
			}
			if math.IsNaN(*bound.n) || math.IsInf(*bound.n, 0) {
				return fmt.Errorf("check bounds must be finite")
			}
			switch fld.Desc.Kind() {
			case protoreflect.DoubleKind, protoreflect.Uint64Kind, protoreflect.Uint32Kind:
			default:
				return fmt.Errorf("numeric checks require numeric scalar")
			}
			invalid = append(invalid, value+bound.op+strconv.FormatFloat(*bound.n, 'g', -1, 64))
		}
		if check.Format != "" {
			if fld.Desc.Kind() != protoreflect.StringKind && !fld.Desc.IsMap() {
				return fmt.Errorf("format check requires string")
			}
			switch check.Format {
			case "env_name":
				g.P("matched:=providerEnvNamePattern.MatchString(", value, ")")
				invalid = append(invalid, "!matched")
			case "domain":
				g.P("domain:=", protogen.GoIdent{GoName: "TrimPrefix", GoImportPath: "strings"}, "(", value, ",\"*.\")")
				g.P("matched:=providerDomainPattern.MatchString(domain)")
				bad := "!matched || len(domain)>253"
				if check.AllowEmpty {
					bad = "(" + value + "!=\"\" && (!matched || len(domain)>253))"
				}
				invalid = append(invalid, bad)
			case "cidr":
				g.P("_,err:=", protogen.GoIdent{GoName: "ParsePrefix", GoImportPath: "net/netip"}, "(", value, ")")
				bad := "err!=nil"
				if check.AllowEmpty {
					bad = "(" + value + "!=\"\" && err!=nil)"
				}
				invalid = append(invalid, bad)
			case "absolute_path":
				invalid = append(invalid, "!"+g.QualifiedGoIdent(protogen.GoIdent{GoName: "HasPrefix", GoImportPath: "strings"})+"("+value+",\"/\")")
			case "http_url":
				g.P("parsed,err:=", protogen.GoIdent{GoName: "Parse", GoImportPath: "net/url"}, "(", value, ")")
				invalid = append(invalid, "err!=nil || parsed.Host==\"\" || (parsed.Scheme!=\"http\" && parsed.Scheme!=\"https\")")
			default:
				return fmt.Errorf("unknown provider format %q", check.Format)
			}
		}
		if len(invalid) == 0 {
			return fmt.Errorf("empty provider check %s", check.Path)
		}
		join := " || "
		if len(check.Allowed) > 0 {
			if check.Minimum != nil || check.Maximum != nil || check.Format != "" {
				return fmt.Errorf("allowed values cannot combine with bounds/format")
			}
			join = " && "
		}
		g.P("if ", strings.Join(invalid, join), "{return ", protogen.GoIdent{GoName: "Errorf", GoImportPath: "fmt"}, "(", strconv.Quote("sandbox-kit "+s.Provider+": "+check.Path+" is outside supported values"), ")}")
		for i := 0; i < guards; i++ {
			g.P("}")
		}
	}
	g.P("return nil}")
	return nil
}
