package daytona

import (
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestPosixShellJoin(t *testing.T) {
	got := posixShellJoin([]string{"echo", "hello world", "it's"})
	want := `'echo' 'hello world' 'it'"'"'s'`
	if got != want {
		t.Fatalf("posixShellJoin = %q, want %q", got, want)
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
