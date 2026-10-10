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
	Import  string             `yaml:"import"`
	Name    string             `yaml:"name"`
	Fields  map[string]string  `yaml:"fields"`
	Objects map[string]Binding `yaml:"objects"`
}
type Type struct {
	Name     string             `yaml:"name"`
	Bindings map[string]Binding `yaml:"bindings"`
}
type Mapping struct {
	Policy    *PolicyMapping `yaml:"policy"`
	From      string         `yaml:"from"`
	To        string         `yaml:"to"`
	Transform string         `yaml:"transform"`
	Factor    uint64         `yaml:"factor"`
	Maximum   *uint64        `yaml:"maximum"`
	Cast      string         `yaml:"cast"`
}
type PolicyMapping struct {
	Disabled  string `yaml:"disabled"`
	Immediate bool   `yaml:"immediate"`
}
type Comparison struct {
	Field       string `yaml:"field"`
	GreaterThan int    `yaml:"greater_than"`
}
type Group struct {
	ErrorPath string       `yaml:"error_path"`
	Origin    string       `yaml:"origin"`
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
	Capture  string `yaml:"capture"`
	Doc      string `yaml:"doc"`
}

// Constructor names the native constructor after common configuration mapping.
// State capture uses portable config field paths, never target-language expressions.
type Constructor struct {
	Function string `yaml:"function"`
}
type Cleanup struct {
	Method         string `yaml:"method"`
	AcceptsContext bool   `yaml:"accepts_context"`
	ReturnsError   bool   `yaml:"returns_error"`
}
type RuntimeBinding struct {
	Constructor Constructor  `yaml:"constructor"`
	State       []StateField `yaml:"state"`
	Cleanup     Cleanup      `yaml:"cleanup"`
}
type Provider struct {
	ErrorRules    map[string]string         `yaml:"error_rules"`
	ErrorStatuses map[int]string            `yaml:"error_statuses"`
	ErrorGRPC     map[int]string            `yaml:"error_grpc"`
	Errors        map[string]ErrorBinding   `yaml:"errors"`
	Checks        []Check                   `yaml:"checks"`
	Client        *ClientMapping            `yaml:"client"`
	Runtime       map[string]RuntimeBinding `yaml:"runtime"`
	Version       int                       `yaml:"version"`
	Provider      string                    `yaml:"provider"`
	Groups        []Group                   `yaml:"groups"`
}
type ErrorBinding struct {
	Types  []ErrorType `yaml:"types"`
	Target Binding     `yaml:"target"`
}
type ErrorType struct {
	Import  string `yaml:"import"`
	Name    string `yaml:"name"`
	Pointer bool   `yaml:"pointer"`
	Rule    string `yaml:"rule"`
}
type Check struct {
	AllowEmpty        bool     `yaml:"allow_empty"`
	ExclusivePolicies []string `yaml:"exclusive_policies"`
	ForbidPolicyWith  string   `yaml:"forbid_policy_with"`
	Path              string   `yaml:"path"`
	Allowed           []string `yaml:"allowed"`
	Minimum           *float64 `yaml:"minimum"`
	Maximum           *float64 `yaml:"maximum"`
	Format            string   `yaml:"format"`
}

// ClientMapping composes existing field mappings into native initialization params.
type ClientMapping struct {
	Target   Type              `yaml:"target"`
	Managed  []string          `yaml:"managed"`
	Retained []string          `yaml:"retained"`
	Rejected map[string]string `yaml:"rejected"`
	Settings string            `yaml:"settings"`
	Scope    *ClientComponent  `yaml:"scope"`
	Auth     []ClientComponent `yaml:"auth"`
}
type ClientComponent struct {
	Field       string   `yaml:"field"`
	Group       string   `yaml:"group"`
	Retained    []string `yaml:"retained"`
	Destination string   `yaml:"destination"`
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

// Client selects schema-backed settings and portable initialization validation.
type Client struct {
	Version       int        `yaml:"version"`
	ConfigMessage string     `yaml:"config_message"`
	Validation    Validation `yaml:"validation"`
}

func LoadClient(root string) (Client, error) {
	var c Client
	err := decode(filepath.Join(root, "client.yaml"), &c)
	if err == nil && (c.Version != 1 || c.ConfigMessage != "Config" || c.Validation.Version != 1 || len(c.Validation.Messages) == 0) {
		err = fmt.Errorf("invalid client specification")
	}
	return c, err
}

type Contracts struct {
	Errors          ErrorPolicy `yaml:"-"`
	Version         int         `yaml:"version"`
	CopyMaxDepth    int         `yaml:"copy_max_depth"`
	MetadataMessage string      `yaml:"metadata_message"`
	ErrorMessage    string      `yaml:"error_message"`
	OriginEnum      string      `yaml:"origin_enum"`
}

func LoadContracts(root string) (Contracts, error) {
	var c Contracts
	err := decode(filepath.Join(root, "contracts.yaml"), &c)
	if err == nil && (c.Version != 1 || c.CopyMaxDepth < 1 || c.CopyMaxDepth > 256 || c.MetadataMessage != "MetadataObject" || c.ErrorMessage != "ErrorInfo" || c.OriginEnum != "ValueOrigin") {
		err = fmt.Errorf("invalid response/error contracts spec")
	}
	if err == nil {
		c.Errors, err = LoadErrorPolicy(root)
	}
	return c, err
}
