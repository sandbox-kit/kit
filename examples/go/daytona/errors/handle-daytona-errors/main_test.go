package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMainHandlesLocalFailures(t *testing.T) {
	if os.Getenv("SANDBOX_KIT_RUN_EXAMPLE") == "1" {
		main()
		return
	}
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/.env", []byte("# Test-only configuration\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DAYTONA_API_KEY", "test-only-key")
	t.Setenv("DAYTONA_API_URL", "http://127.0.0.1:1")
	command := exec.Command(os.Args[0], "-test.run=^TestMainHandlesLocalFailures$")
	command.Dir = dir
	command.Env = append(os.Environ(), "SANDBOX_KIT_RUN_EXAMPLE=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("example failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "Check resources.cpu_cores") {
		t.Fatalf("unexpected output: %s", output)
	}
}
