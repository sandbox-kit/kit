package spec

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TypeRef describes meaning independently of pointers, undefined, or null syntax.
type TypeRef struct {
	Arguments []TypeRef     `yaml:"arguments,omitempty"`
	Function  *FunctionType `yaml:"function,omitempty"`
	Tuple     []TypeRef     `yaml:"tuple,omitempty"`
	Ref       string        `yaml:"ref,omitempty"`
	Builtin   string        `yaml:"builtin,omitempty"`
	Optional  bool          `yaml:"optional,omitempty"`
	Nullable  bool          `yaml:"nullable,omitempty"`
	Reference bool          `yaml:"reference,omitempty"`
	Ownership string        `yaml:"ownership,omitempty"`
	Sequence  *TypeRef      `yaml:"sequence,omitempty"`
	Map       *MapType      `yaml:"map,omitempty"`
	Union     []TypeRef     `yaml:"union,omitempty"`
}
type FunctionType struct {
	Params      []APISlot `yaml:"params"`
	Results     []APISlot `yaml:"results"`
	Fallible    bool      `yaml:"fallible"`
	Cancellable bool      `yaml:"cancellable"`
}
type EnumValue struct {
	Name   string `yaml:"name"`
	Number int64  `yaml:"number"`
}

type MapType struct {
	Key   TypeRef `yaml:"key"`
	Value TypeRef `yaml:"value"`
}
type APISlot struct {
	RuntimeOnly bool    `yaml:"runtime_only"`
	Promote     bool    `yaml:"promote"`
	Doc         string  `yaml:"doc"`
	ID          string  `yaml:"id"`
	Name        string  `yaml:"name"`
	Type        TypeRef `yaml:"type"`
	Variadic    bool    `yaml:"variadic"`
}
type Callable struct {
	Kind         string    `yaml:"kind"`
	Constructs   string    `yaml:"constructs"`
	Static       bool      `yaml:"static"`
	Effect       string    `yaml:"effect"`
	Doc          string    `yaml:"doc"`
	ID           string    `yaml:"id"`
	Name         string    `yaml:"name"`
	Behavior     string    `yaml:"behavior"`
	Receiver     string    `yaml:"receiver"`
	Params       []APISlot `yaml:"params"`
	Results      []APISlot `yaml:"results"`
	Fallible     bool      `yaml:"fallible"`
	Cancellable  bool      `yaml:"cancellable"`
	NamedResults bool      `yaml:"named_results"`
	Generics     []APISlot `yaml:"generics"`
}
type APIType struct {
	Generics   []APISlot     `yaml:"generics"`
	Underlying *TypeRef      `yaml:"underlying"`
	Values     []EnumValue   `yaml:"values"`
	Signature  *FunctionType `yaml:"signature"`
	Doc        string        `yaml:"doc"`
	Visibility string        `yaml:"visibility"`
	ID         string        `yaml:"id"`
	Name       string        `yaml:"name"`
	Kind       string        `yaml:"kind"`
	Fields     []APISlot     `yaml:"fields"`
	FieldsFrom string        `yaml:"fields_from"`
	Methods    []Callable    `yaml:"methods"`
}
type APIModule struct {
	Attachments []APIAttachment `yaml:"attachments"`
	Types       []APIType       `yaml:"types"`
	Functions   []Callable      `yaml:"functions"`
}

type APIAttachment struct {
	Target string    `yaml:"target"`
	Fields []APISlot `yaml:"fields"`
}
type API struct {
	Version int                  `yaml:"version"`
	Modules map[string]APIModule `yaml:"modules"`
}
type ExternalType struct {
	Import string `yaml:"import"`
	Name   string `yaml:"name"`
	Kind   string `yaml:"kind"`
}
type LanguageProfile struct {
	Outputs         map[string]string `yaml:"outputs"`
	Version         int               `yaml:"version"`
	Language        string            `yaml:"language"`
	ObjectKind      string            `yaml:"object_kind"`
	InterfaceKind   string            `yaml:"interface_kind"`
	Reference       string            `yaml:"reference"`
	Optional        string            `yaml:"optional"`
	Nullable        string            `yaml:"nullable"`
	Fallible        string            `yaml:"fallible"`
	MultipleResults string            `yaml:"multiple_results"`
	Cancellation    struct {
		Name string       `yaml:"name"`
		Type ExternalType `yaml:"type"`
	} `yaml:"cancellation"`
	ErrorResultName string                  `yaml:"error_result_name"`
	Unions          string                  `yaml:"unions"`
	Builtins        map[string]string       `yaml:"builtins"`
	Externals       map[string]ExternalType `yaml:"externals"`
}

func LoadAPI(root string) (API, error) {
	var a API
	err := decode(filepath.Join(root, "api.yaml"), &a)
	if err == nil && (a.Version != 1 || len(a.Modules) == 0) {
		err = fmt.Errorf("unsupported or empty API declarations")
	}
	return a, err
}
func LoadLanguageProfile(root, language string) (LanguageProfile, error) {
	var p LanguageProfile
	if err := decode(filepath.Join(root, "languages", language+".yaml"), &p); err != nil {
		return p, err
	}
	if p.Version != 1 || p.Language != language || len(p.Outputs) == 0 {
		return p, fmt.Errorf("invalid language profile %s", language)
	}
	for role, suffix := range p.Outputs {
		if suffix == "" || !strings.HasPrefix(suffix, ".") || strings.ContainsAny(suffix, "/\\") {
			return p, fmt.Errorf("invalid output suffix for %s", role)
		}
	}
	return p, nil
}
