package main

import (
	"os"
	"testing"
)

func TestRunHandlesLocalFailures(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("# Test-only configuration\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("DAYTONA_API_KEY", "test-only-key")
	t.Setenv("DAYTONA_API_URL", "http://127.0.0.1:1")
	if err := run(); err != nil {
		t.Fatal(err)
	}
}
