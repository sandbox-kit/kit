// Package lifecyclegen emits the TypeScript lifecycle surface from the
// portable lifecycle specification.
package lifecyclegen

import (
	"fmt"
	"strings"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
	"github.com/sandbox-kit/kit/tooling/internal/typescript/naming"
	"google.golang.org/protobuf/compiler/protogen"
)

// EmitBackendMethods adds optional provider operations to Backend.
func EmitBackendMethods(g *protogen.GeneratedFile, lifecycle spec.Lifecycle) {
	for _, operation := range lifecycle.Operations {
		if operation.Options == "" {
			g.P("  ", operation.ID, "?(id: string, options: OperationOptions): Promise<SandboxInfo>;")
			continue
		}
		result := "SandboxInfo"
		if operation.Outcome == "exit_result" {
			result = "ExitResult"
		}
		g.P("  ", operation.ID, "?(id: string, request: ", operation.Options, ", options: OperationOptions): Promise<", result, ">;")
	}
}

// EmitHandleMethods implements Sandbox lifecycle operations from their
// declared outcomes. Unsupported optional backend methods fail consistently.
func EmitHandleMethods(g *protogen.GeneratedFile, lifecycle spec.Lifecycle) error {
	for _, operation := range lifecycle.Operations {
		switch operation.Outcome {
		case "update_handle":
			if operation.Options == "" {
				g.P("  async ", operation.ID, "(options: OperationOptions = {}): Promise<SandboxInfo> { this.#ensureOpen(\"", operation.ID, "\"); const method = this.#backend.", operation.ID, "; if (!method) unsupported(this.providerName, \"", operation.ID, "\"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, this.#timeoutMs); const info = await awaitWithSignal(() => method.call(this.#backend, this.id, operationOptions(options, scoped?.signal)), scoped.signal); this.#info = clone({ ...info, provider: this.providerName }); return this.info(); } catch (error) { throw withContext(error, this.providerName, \"", operation.ID, "\"); } finally { scoped?.dispose(); } }")
				continue
			}
			g.P("  async ", operation.ID, "(request: ", operation.Options, " = {}, options: OperationOptions = {}): Promise<SandboxInfo> { return this.#update(\"", operation.ID, "\", this.#backend.", operation.ID, ", request, options); }")
		case "exit_result":
			if operation.Options == "" || operation.TerminalState == "" {
				return fmt.Errorf("TypeScript lifecycle operation %s has an incomplete exit result declaration", operation.ID)
			}
			g.P("  async ", operation.ID, "(request: ", operation.Options, " = {}, options: OperationOptions = {}): Promise<ExitResult> { this.#ensureOpen(\"", operation.ID, "\"); const method = this.#backend.", operation.ID, "; if (!method) unsupported(this.providerName, \"", operation.ID, "\"); let scoped: ReturnType<typeof operationSignal> | undefined; try { scoped = operationSignal(options, request.timeout ?? this.#timeoutMs); const result = await awaitWithSignal(() => method.call(this.#backend, this.id, clone(request), operationOptions(options, scoped?.signal)), scoped.signal); this.#info = { ...this.#info, state: SandboxState.", naming.Pascal(operation.TerminalState), ", terminal: true }; return result; } catch (error) { throw withContext(error, this.providerName, \"", operation.ID, "\"); } finally { scoped?.dispose(); } }")
		default:
			return fmt.Errorf("TypeScript lifecycle operation %s has unsupported outcome %s", operation.ID, operation.Outcome)
		}
	}
	return nil
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
