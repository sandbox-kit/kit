package lifecyclegen

import (
	"testing"

	"github.com/sandbox-kit/kit/tooling/internal/spec"
)

func TestUpdateOptionsUnionUsesDeclaredUpdateOperations(t *testing.T) {
	lifecycle := spec.Lifecycle{Operations: []spec.LifecycleOperation{
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
