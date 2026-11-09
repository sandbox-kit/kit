package blaxel

import (
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestBlaxelCommandArgvAndShell(t *testing.T) {
	command, err := blaxelCommand(&sandbox.ExecOptions{Argv: []string{"echo", "hello world"}})
	if err != nil || command != "'echo' 'hello world'" {
		t.Fatalf("argv command = %q err=%v", command, err)
	}
	command, err = blaxelCommand(&sandbox.ExecOptions{Shell: sandbox.Value("echo hi")})
	if err != nil || command != "echo hi" {
		t.Fatalf("shell command = %q err=%v", command, err)
	}
}

func TestCommandTimeoutWholeSeconds(t *testing.T) {
	timeout, err := commandTimeout("run", sandbox.Duration(2*time.Second))
	if err != nil || timeout != 2*time.Second {
		t.Fatalf("timeout = %v err=%v", timeout, err)
	}
	_, err = commandTimeout("run", sandbox.Duration(1500*time.Millisecond))
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Field != sandbox.ExecFieldTimeout {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}
