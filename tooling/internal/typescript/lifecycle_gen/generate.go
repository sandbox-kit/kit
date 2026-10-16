// Package lifecyclegen emits the TypeScript lifecycle surface from the
// portable lifecycle specification.
package lifecyclegen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	declarationgen "github.com/sandbox-kit/kit/tooling/internal/typescript/declaration_gen"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	"google.golang.org/protobuf/compiler/protogen"
)

// EmitBackendMethods adds optional provider operations to Backend.
func EmitBackendMethods(g *protogen.GeneratedFile, lifecycle spec.Lifecycle) {
	e := declarationgen.NewEmitter(g)
	for _, operation := range lifecycle.Operations {
		e.Line("  ", backendSignature(operation))
	}
}

// EmitHandleMethods implements Sandbox lifecycle operations from their
// declared outcomes. Unsupported optional backend methods fail consistently.
func EmitHandleMethods(g *protogen.GeneratedFile, lifecycle spec.Lifecycle) error {
	e := declarationgen.NewEmitter(g)
	for _, operation := range lifecycle.Operations {
		method, err := handleMethod(operation)
		if err != nil {
			return err
		}
		e.IndentedBlock(method, "  ")
	}
	return nil
}

func backendSignature(operation spec.LifecycleOperation) string {
	result := "SandboxInfo"
	if operation.Outcome == "exit_result" {
		result = "ExitResult"
	}
	request := ""
	if operation.Options != "" {
		request = "request: " + operation.Options + ", "
	}
	return fmt.Sprintf("%s?(id: string, %soptions: OperationOptions): Promise<%s>;", operation.ID, request, result)
}

func handleMethod(operation spec.LifecycleOperation) (string, error) {
	switch operation.Outcome {
	case "update_handle":
		if operation.Options == "" {
			return fmt.Sprintf(`
  async %s(options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#serialized(%q, options, this.#timeoutMs, async (signal) => {
      const method = this.#backend.%s;
      if (!method) unsupported(this.providerName, %q);
      const info = await awaitWithSignal(() => method.call(this.#backend, this.id, operationOptions(options, signal)), signal);
      this.#info = clone({ ...info, provider: this.providerName });
      return this.info();
    });
  }`, operation.ID, operation.ID, operation.ID, operation.ID), nil
		}
		return fmt.Sprintf(`
  async %s(request: %s = {}, options: OperationOptions = {}): Promise<SandboxInfo> {
    return this.#update(%q, this.#backend.%s, request, options);
  }`, operation.ID, operation.Options, operation.ID, operation.ID), nil
	case "exit_result":
		if operation.Options == "" || operation.TerminalState == "" {
			return "", fmt.Errorf("TypeScript lifecycle operation %s has an incomplete exit result declaration", operation.ID)
		}
		return fmt.Sprintf(`
  async %s(request: %s = {}, options: OperationOptions = {}): Promise<ExitResult> {
    return this.#serialized(%q, options, request.timeout ?? this.#timeoutMs, async (signal) => {
      const method = this.#backend.%s;
      if (!method) unsupported(this.providerName, %q);
      const result = await awaitWithSignal(() => method.call(this.#backend, this.id, clone(request), operationOptions(options, signal)), signal);
      this.#info = { ...this.#info, state: SandboxState.%s, terminal: true };
      return result;
    });
  }`, operation.ID, operation.Options, operation.ID, operation.ID, operation.ID, naming.Pascal(operation.TerminalState)), nil
	default:
		return "", fmt.Errorf("TypeScript lifecycle operation %s has unsupported outcome %s", operation.ID, operation.Outcome)
	}
}

// UpdateOptionsUnion returns the request types accepted by the shared handle
// update helper. Exit-result and requestless operations do not use it.
func UpdateOptionsUnion(lifecycle spec.Lifecycle) (string, error) {
	options := make([]string, 0, len(lifecycle.Operations))
	seen := map[string]bool{}
	for _, operation := range lifecycle.Operations {
		if operation.Outcome != "update_handle" || operation.Options == "" || seen[operation.Options] {
			continue
		}
		seen[operation.Options] = true
		options = append(options, operation.Options)
	}
	if len(options) == 0 {
		return "", fmt.Errorf("TypeScript lifecycle has no update request types")
	}
	return strings.Join(options, " | "), nil
}
