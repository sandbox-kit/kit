package lifecyclegen

import (
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/model"
)

func TestUpdateOptionsUnionUsesDeclaredUpdateOperations(t *testing.T) {
	lifecycle := model.Lifecycle{Operations: []model.LifecycleOperation{
		{ID: "refresh", Outcome: "update_handle"},
		{ID: "start", Outcome: "update_handle", Options: "StartOptions"},
		{ID: "startAgain", Outcome: "update_handle", Options: "StartOptions"},
		{ID: "terminate", Outcome: "exit_result", Options: "TerminateOptions"},
		{ID: "stop", Outcome: "update_handle", Options: "StopOptions"},
	}}

	got, err := UpdateOptionsUnion(lifecycle)
	if err != nil {
		t.Fatalf("UpdateOptionsUnion() error = %v", err)
	}
	if want := "StartOptions | StopOptions"; got != want {
		t.Fatalf("UpdateOptionsUnion() = %q, want %q", got, want)
	}
}

func TestBackendSignatureUsesLifecycleDeclaration(t *testing.T) {
	operation := model.LifecycleOperation{ID: "terminate", Options: "TerminateOptions", Outcome: "exit_result"}
	if got, want := backendSignature(operation), "terminate?(id: string, request: TerminateOptions, options: OperationOptions): Promise<ExitResult>;"; got != want {
		t.Fatalf("backendSignature() = %q, want %q", got, want)
	}
}

func TestHandleMethodUsesDeclaredTerminalState(t *testing.T) {
	operation := model.LifecycleOperation{ID: "terminate", Options: "TerminateOptions", Outcome: "exit_result", TerminalState: "terminated"}
	got, err := handleMethod(operation)
	if err != nil {
		t.Fatalf("handleMethod() error = %v", err)
	}
	for _, want := range []string{"async terminate(", "this.#backend.terminate", "SandboxState.Terminated"} {
		if !strings.Contains(got, want) {
			t.Fatalf("handleMethod() missing %q:\n%s", want, got)
		}
	}
}
