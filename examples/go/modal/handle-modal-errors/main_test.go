package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestDiskOverrideReturnsUnsupported(t *testing.T) {
	client, err := sandbox.NewClient(sandbox.Config{
		Provider: modal.New(),
		Auth:     &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: "test-only-id", Secret: "test-only-secret"}},
		Scope:    &sandbox.Scope{AppName: sandbox.Value("test-only-app")},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := client.Close(context.Background()); err != nil {
			t.Error(err)
		}
	})
	_, err = client.Create(context.Background(), &sandbox.CreateOptions{
		Source:    &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}},
		Resources: &sandbox.Resources{DiskMiB: sandbox.Value(uint64(1024))},
	})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != "disk_mib" {
		t.Fatalf("expected unsupported disk override, got %v", err)
	}
}

func TestRunHandlesLocalFailures(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("# Test-only configuration\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MODAL_TOKEN_ID", "test-only-id")
	t.Setenv("MODAL_TOKEN_SECRET", "test-only-secret")
	t.Setenv("MODAL_APP_NAME", "test-only-app")
	if err := run(); err != nil {
		t.Fatal(err)
	}
}
