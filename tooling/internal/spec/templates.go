package spec

import "path/filepath"

type DeclarationTemplate struct {
	Owner       string   `yaml:"owner"`
	Declaration Callable `yaml:"declaration"`
}
type FieldPathTemplates struct {
	Name             string            `yaml:"name"`
	Roots            map[string]string `yaml:"roots"`
	NativeLeaves     []string          `yaml:"native_leaves"`
	NativeNamespaces []string          `yaml:"native_namespaces"`
	ExcludedOrigins  []string          `yaml:"excluded_origins"`
}
type Templates struct {
	Version    int                            `yaml:"version"`
	Callables  map[string]DeclarationTemplate `yaml:"callables"`
	FieldPaths FieldPathTemplates             `yaml:"field_paths"`
}

func LoadTemplates(root string) (Templates, error) {
	var t Templates
	err := decode(filepath.Join(root, "templates.yaml"), &t)
	return t, err
}
