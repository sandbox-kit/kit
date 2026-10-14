package declarationgen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
)

// Signature lowers a compiled callable to a TypeScript method signature. The
// callable ID is the stable semantic name; Go-specific display names from
// api.yaml are not reused as JavaScript method names.
func (r *Renderer) Signature(callable spec.Callable) (string, error) {
	params := make([]string, 0, len(callable.Params)+1)
	for _, parameter := range callable.Params {
		typeName, err := r.Type(parameter.Type)
		if err != nil {
			return "", err
		}
		optional := ""
		if parameter.Type.Optional || parameter.Type.Nullable {
			optional = "?"
			typeName = strings.TrimSuffix(strings.TrimSuffix(typeName, " | undefined"), " | null")
		}
		params = append(params, naming.Camel(parameter.Name)+optional+": "+typeName)
	}
	if callable.Cancellable {
		params = append(params, "options?: OperationOptions")
	}
	result, err := r.callableResult(callable)
	if err != nil {
		return "", err
	}
	return callableName(callable.ID) + "(" + strings.Join(params, ", ") + "): " + result, nil
}

func (r *Renderer) callableResult(callable spec.Callable) (string, error) {
	result := "void"
	if len(callable.Results) == 1 {
		var err error
		result, err = r.Type(callable.Results[0].Type)
		if err != nil {
			return "", err
		}
	} else if len(callable.Results) > 1 {
		fields := make([]string, 0, len(callable.Results))
		for _, output := range callable.Results {
			typeName, err := r.Type(output.Type)
			if err != nil {
				return "", err
			}
			if output.Name == "" {
				return "", fmt.Errorf("TypeScript multiple results require names")
			}
			fields = append(fields, "readonly "+naming.Camel(output.Name)+": "+typeName)
		}
		result = "{ " + strings.Join(fields, "; ") + " }"
	}
	if callable.Fallible || callable.Effect == "io" {
		result = "Promise<" + result + ">"
	}
	return result, nil
}

func callableName(id string) string {
	id = strings.TrimSuffix(id, "_sandbox")
	if id == "constructor" {
		return "create"
	}
	return naming.Camel(id)
}
