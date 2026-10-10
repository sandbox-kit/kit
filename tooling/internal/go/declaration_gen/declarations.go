package declarationgen

import (
	"fmt"
	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"strings"
)

func (e *Emitter) TypeDeclaration(id string) error {
	typ, ok := e.Module.TypesByID[id]
	if !ok {
		return fmt.Errorf("missing API type %s", id)
	}
	local, err := e.genericScope(typ.Generics)
	if err != nil {
		return err
	}
	e = local
	generic := ""
	if len(typ.Generics) > 0 {
		parameters := []string{}
		for _, parameter := range typ.Generics {
			constraint, err := e.Type(parameter.Type)
			if err != nil {
				return err
			}
			parameters = append(parameters, parameter.Name+" "+constraint)
		}
		generic = "[" + strings.Join(parameters, ",") + "]"
	}
	e.comment(typ.Name, typ.Doc)
	switch typ.Kind {
	case "alias", "enum":
		underlying, err := e.Type(*typ.Underlying)
		if err != nil {
			return err
		}
		separator := " "
		if typ.Kind == "alias" {
			separator = " = "
		}
		if typ.Kind == "enum" && (typ.Underlying.Builtin != "integer" && typ.Underlying.Builtin != "unsigned") {
			return fmt.Errorf("Go enum requires integer underlying type")
		}
		e.G.P("type ", typ.Name, generic, separator, underlying)
		if typ.Kind == "enum" {
			e.G.P("const (")
			for _, value := range typ.Values {
				e.G.P(value.Name, " ", typ.Name, " = ", value.Number)
			}
			e.G.P(")")
		}
		return nil
	case "callback":
		signature, err := e.Type(spec.TypeRef{Function: typ.Signature})
		if err != nil {
			return err
		}
		e.G.P("type ", typ.Name, generic, " ", signature)
		return nil
	}
	kind := e.Profile.ObjectKind
	if typ.Kind == "interface" {
		kind = e.Profile.InterfaceKind
	}
	e.G.P("type ", typ.Name, generic, " ", kind, "{")
	for _, field := range typ.Fields {
		e.comment(field.Name, field.Doc)
		value, err := e.Type(field.Type)
		if err != nil {
			return err
		}
		if field.Promote {
			e.G.P(value)
		} else {
			e.G.P(field.Name, " ", value)
		}
	}
	if typ.Kind == "interface" {
		for _, method := range typ.Methods {
			e.comment(method.Name, method.Doc)
			signature, err := e.Signature(id, method, true)
			if err != nil {
				return err
			}
			e.G.P(signature)
		}
	}
	e.G.P("}")
	return nil
}

func (e *Emitter) genericScope(parameters []spec.APISlot) (*Emitter, error) {
	local := *e
	local.TypeVariables = map[string]string{}
	for _, parameter := range parameters {
		if err := goIdentifier(parameter.Name); err != nil {
			return nil, err
		}
		local.TypeVariables[parameter.ID] = parameter.Name
	}
	return &local, nil
}
func (e *Emitter) Callable(owner, id string) (spec.Callable, error) {
	if owner == "" {
		if c, ok := e.Module.FunctionsByID[id]; ok {
			return c, nil
		}
	} else if typ, ok := e.Module.TypesByID[owner]; ok {
		for _, c := range typ.Methods {
			if c.ID == id {
				return c, nil
			}
		}
	}
	return spec.Callable{}, fmt.Errorf("missing API callable %s.%s", owner, id)
}
func (e *Emitter) Begin(owner, id string) (spec.Callable, error) {
	c, err := e.Callable(owner, id)
	if err != nil {
		return c, err
	}
	signature, err := e.Signature(owner, c, false)
	if err != nil {
		return c, err
	}
	e.comment(c.Name, c.Doc)
	e.G.P(signature, "{")
	return c, nil
}

func (e *Emitter) BeginBehavior(owner, id, behavior string) (spec.Callable, error) {
	c, err := e.Callable(owner, id)
	if err != nil {
		return c, err
	}
	if c.Behavior != behavior {
		return c, fmt.Errorf("unsupported behavior %s for %s.%s", c.Behavior, owner, id)
	}
	return e.Begin(owner, id)
}
func (e *Emitter) MethodName(owner, id string) string { c, _ := e.Callable(owner, id); return c.Name }
func (e *Emitter) comment(name, doc string) {
	if doc != "" {
		for i, line := range strings.Split(doc, "\n") {
			if i == 0 {
				line = name + " " + line
			}
			e.G.P("// ", line)
		}
	}
}
