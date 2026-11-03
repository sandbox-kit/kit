package modal

import (
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestModalCommand(t *testing.T) {
	argv, err := modalCommand(&sandbox.ExecOptions{Argv: []string{"echo", "hi"}})
	if err != nil || len(argv) != 2 || argv[0] != "echo" || argv[1] != "hi" {
		t.Fatalf("argv = %#v, %v", argv, err)
	}
	shell, err := modalCommand(&sandbox.ExecOptions{Shell: sandbox.Value("echo hi")})
	if err != nil || len(shell) != 3 || shell[0] != "bash" || shell[1] != "-c" || shell[2] != "echo hi" {
		t.Fatalf("shell = %#v, %v", shell, err)
	}
}

func TestExecTimeout(t *testing.T) {
	if got, err := execTimeout("run", nil); err != nil || got != 0 {
		t.Fatalf("nil timeout = %v, %v", got, err)
	}
	if got, err := execTimeout("run", sandbox.Duration(0)); err != nil || got != 0 {
		t.Fatalf("zero timeout = %v, %v", got, err)
	}
	if got, err := execTimeout("run", sandbox.Duration(5*time.Second)); err != nil || got != 5*time.Second {
		t.Fatalf("whole timeout = %v, %v", got, err)
	}
	_, err := execTimeout("run", sandbox.Duration(1500*time.Millisecond))
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Field != sandbox.ExecFieldTimeout {
		t.Fatalf("fractional timeout = %#v", err)
	}
}
