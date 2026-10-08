// Package spec loads language-neutral generation rules. It has no Go runtime SDK dependencies.
package spec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Validation struct {
	Version  int                `yaml:"version"`
	Messages map[string]Message `yaml:"messages"`
}
type Message struct {
	Fields      map[string]Field `yaml:"fields"`
	Constraints []Constraint     `yaml:"constraints"`
}
type Field struct {
	Format           string   `yaml:"format"`
	Sensitive        bool     `yaml:"sensitive"`
	Minimum          *float64 `yaml:"minimum"`
	ExclusiveMinimum *float64 `yaml:"exclusive_minimum"`
	Maximum          *float64 `yaml:"maximum"`
	Finite           bool     `yaml:"finite"`
	Nonblank         bool     `yaml:"nonblank"`
	Each             bool     `yaml:"each"`
	MinItems         *int     `yaml:"min_items"`
}
type Condition struct {
	Field     string `yaml:"field"`
	Equals    *int   `yaml:"equals"`
	NotEquals *int   `yaml:"not_equals"`
}
type Constraint struct {
	Op     string     `yaml:"op"`
	Field  string     `yaml:"field"`
	Other  string     `yaml:"other"`
	Fields []string   `yaml:"fields"`
	When   *Condition `yaml:"when"`
}
type Binding struct {
	Import string            `yaml:"import"`
	Name   string            `yaml:"name"`
	Fields map[string]string `yaml:"fields"`
}
type Type struct {
	Name     string             `yaml:"name"`
	Bindings map[string]Binding `yaml:"bindings"`
}
type Mapping struct {
	From      string  `yaml:"from"`
	To        string  `yaml:"to"`
	Transform string  `yaml:"transform"`
	Factor    uint64  `yaml:"factor"`
	Maximum   *uint64 `yaml:"maximum"`
	Cast      string  `yaml:"cast"`
}
type Comparison struct {
	Field       string `yaml:"field"`
	GreaterThan int    `yaml:"greater_than"`
}
type Group struct {
	When      []Comparison `yaml:"when"`
	Name      string       `yaml:"name"`
	Source    Type         `yaml:"source"`
	Target    Type         `yaml:"target"`
	Direction string       `yaml:"direction"`
	Fields    []Mapping    `yaml:"fields"`
}
type StateField struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Optional bool   `yaml:"optional"`
}
type Cleanup struct {
	Method         string `yaml:"method"`
	AcceptsContext bool   `yaml:"accepts_context"`
	ReturnsError   bool   `yaml:"returns_error"`
}
type RuntimeBinding struct {
	State   []StateField `yaml:"state"`
	Cleanup Cleanup      `yaml:"cleanup"`
}
type Provider struct {
	Runtime  map[string]RuntimeBinding `yaml:"runtime"`
	Version  int                       `yaml:"version"`
	Provider string                    `yaml:"provider"`
	Groups   []Group                   `yaml:"groups"`
}

func decode(path string, value any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err = decoder.Decode(value); err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	var extra any
	if err = decoder.Decode(&extra); err != io.EOF {
		return fmt.Errorf("%s: expected exactly one YAML document", path)
	}
	return nil
}
func LoadValidation(root string) (Validation, error) {
	var v Validation
	err := decode(filepath.Join(root, "validation.yaml"), &v)
	if err == nil && (v.Version != 1 || len(v.Messages) == 0) {
		err = fmt.Errorf("unsupported validation spec version %d", v.Version)
	}
	return v, err
}
func LoadProvider(root, name string) (Provider, error) {
	var p Provider
	err := decode(filepath.Join(root, "providers", name+".yaml"), &p)
	if err == nil && (p.Version != 1 || p.Provider != name || len(p.Groups) == 0) {
		err = fmt.Errorf("invalid provider spec version/identity: %s", name)
	}
	return p, err
}

// Client declares native factory attachment and portable initialization validation.
type Client struct {
	Version          int        `yaml:"version"`
	ConfigMessage    string     `yaml:"config_message"`
	ProviderField    string     `yaml:"provider_field"`
	ProviderContract string     `yaml:"provider_contract"`
	Validation       Validation `yaml:"validation"`
}

func LoadClient(root string) (Client, error) {
	var c Client
	err := decode(filepath.Join(root, "client.yaml"), &c)
	if err == nil && (c.Version != 1 || c.ConfigMessage != "Config" || c.ProviderField != "provider" || c.ProviderContract != "Provider" || c.Validation.Version != 1 || len(c.Validation.Messages) == 0) {
		err = fmt.Errorf("invalid client specification")
	}
	return c, err
}
