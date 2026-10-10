package sandbox_test

import (
	"context"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"strings"
	"testing"
)

func TestSecureEndpoints(t *testing.T) {
	for _, tc := range []struct {
		endpoint string
		valid    bool
	}{
		{"https://api.example.test/api", true}, {"http://127.0.0.1:1234", true}, {"http://[::1]:1234", true}, {"http://localhost:1234", true},
		{"http://api.example.test", false}, {"http://127.0.0.1.evil.test", false}, {"http://0.0.0.0:1234", false}, {"https://user:secret@api.example.test", false}, {"https://api.example.test/#fragment", false},
	} {
		t.Run(tc.endpoint, func(t *testing.T) {
			err := sandbox.ValidateConfig(&sandbox.Config{Endpoint: sandbox.Value(tc.endpoint)})
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v error=%v", tc.valid, err)
			}
		})
	}
}

func TestImageReferenceValidationBeforeBackend(t *testing.T) {
	for _, tc := range []struct {
		reference string
		valid     bool
	}{
		{"alpine:3.21", true}, {"registry.example.test:5000/team/image:release", true}, {"ubuntu", true}, {"alpine@sha256:" + strings.Repeat("a", 64), true},
		{"alpine\nRUN echo injected", false}, {"alpine\rRUN echo injected", false}, {"alpine AS stage", false}, {"alpine;echo injected", false}, {"--platform=linux/amd64 alpine", false}, {"alpine@sha256:bad", false},
	} {
		t.Run(tc.reference, func(t *testing.T) {
			called := false
			client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
				called = true
				return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "test"}}, nil
			}})
			_, err := client.Create(context.Background(), &sandbox.CreateOptions{Source: &sandbox.SandboxSource{Image: &sandbox.ImageSource{Reference: tc.reference}}})
			if (err == nil) != tc.valid || called != tc.valid {
				t.Fatalf("valid=%v called=%v error=%v", tc.valid, called, err)
			}
		})
	}
}

func TestMetadataCopyBudgets(t *testing.T) {
	var graph any = "leaf"
	for i := 0; i < 20; i++ {
		graph = []any{graph, graph}
	}
	for name, value := range map[string]any{"shared graph": graph, "wide list": make([]any, 16385), "bytes": make([]byte, 1048577), "string": strings.Repeat("x", 1048577), "string list": []string{strings.Repeat("x", 1048577)}, "map key": map[string]any{strings.Repeat("x", 1048577): true}} {
		t.Run(name, func(t *testing.T) {
			_, err := (&sandbox.CreateOptions{ProviderOptions: map[string]any{"value": value}}).Clone()
			if err == nil || !strings.Contains(err.Error(), "budget") {
				t.Fatalf("expected bounded rejection, got %v", err)
			}
		})
	}
	_, err := (&sandbox.CreateOptions{ProviderOptions: map[string]any{"value": []byte{1, 2}}}).Clone()
	if err != nil {
		t.Fatal(err)
	}
}
