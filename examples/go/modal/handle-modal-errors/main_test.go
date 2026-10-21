package main

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestDiskOverrideReturnsUnsupported(t *testing.T) {
	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.TokenPair{
				TokenPair: &sandbox.TokenPairCredentials{
					Id:     "test-only-id",
					Secret: "test-only-secret",
				},
			},
		},
		Scope: &sandbox.Scope{AppName: sandbox.Value("test-only-app")},
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
		Source:    &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}}},
		Resources: &sandbox.Resources{DiskMib: sandbox.Value(uint64(1024))},
	})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != "disk_mib" {
		t.Fatalf("expected unsupported disk override, got %v", err)
	}
}

func TestMainHandlesLocalFailures(t *testing.T) {
	if os.Getenv("SANDBOX_KIT_RUN_EXAMPLE") == "1" {
		main()
		return
	}
	dir := t.TempDir()
	if err := os.WriteFile(dir+"/.env", []byte("# Test-only configuration\n"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("MODAL_TOKEN_ID", "test-only-id")
	t.Setenv("MODAL_TOKEN_SECRET", "test-only-secret")
	t.Setenv("MODAL_APP_NAME", "test-only-app")
	command := exec.Command(os.Args[0], "-test.run=^TestMainHandlesLocalFailures$")
	command.Dir = dir
	command.Env = append(os.Environ(), "SANDBOX_KIT_RUN_EXAMPLE=1")
	output, err := command.CombinedOutput()
	if err != nil {
		t.Fatalf("example failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), "Check disk_mib") {
		t.Fatalf("unexpected output: %s", output)
	}
}
