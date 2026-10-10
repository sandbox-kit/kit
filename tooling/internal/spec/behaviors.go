package spec

import "path/filepath"

// BehaviorSignature describes the inputs and outputs a portable behavior uses.
// Native names and calling syntax are deliberately excluded.
type BehaviorSlot struct {
	ID       string  `yaml:"id,omitempty"`
	Type     TypeRef `yaml:"type"`
	Variadic bool    `yaml:"variadic,omitempty"`
}

type BehaviorSignature struct {
	Owner       string         `yaml:"owner"`
	Params      []BehaviorSlot `yaml:"params,omitempty"`
	Results     []BehaviorSlot `yaml:"results,omitempty"`
	Generics    []BehaviorSlot `yaml:"generics,omitempty"`
	Fallible    bool           `yaml:"fallible"`
	Cancellable bool           `yaml:"cancellable"`
}
type Behaviors struct {
	Version    int                          `yaml:"version"`
	Signatures map[string]BehaviorSignature `yaml:"signatures"`
}

func LoadBehaviors(root string) (Behaviors, error) {
	var b Behaviors
	err := decode(filepath.Join(root, "behaviors.yaml"), &b)
	return b, err
}
