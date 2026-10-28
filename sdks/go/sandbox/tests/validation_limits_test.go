package sandbox_test

import (
	"context"
	"strings"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
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
				return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{Id: "test"}}, nil
			}})
			_, err := client.Create(context.Background(), &sandbox.CreateOptions{Source: &sandbox.SandboxSource{Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: tc.reference}}}})
			if (err == nil) != tc.valid || called != tc.valid {
				t.Fatalf("valid=%v called=%v error=%v", tc.valid, called, err)
			}
		})
	}
}

