package main

import (
	"context"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestPublicClientRejectsLossyResourcesBeforeRemoteCalls(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, "unexpected remote call", 500)
	}))
	defer server.Close()
	t.Setenv("MODAL_SERVER_URL", server.URL)
	t.Setenv("MODAL_TOKEN_ID", "test-only-id")
	t.Setenv("MODAL_TOKEN_SECRET", "test-only-secret")
	t.Setenv("MODAL_APP_NAME", "test-app")
	t.Setenv("MODAL_ENVIRONMENT", "dev")
	config := clientConfig()
	if err := sandbox.ValidateCreateOptions(createOptions()); err != nil {
		t.Fatalf("example options invalid: %v", err)
	}
	client, err := sandbox.NewClient(config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := client.Close(context.Background()); err != nil {
			t.Error(err)
		}
	}()
	// Mutating the input after construction must not erase the captured app name.
	*config.Scope.AppName = ""
	for _, tc := range []struct {
		name   string
		modify func(*sandbox.CreateOptions)
		want   string
	}{
		{"below minimum CPU", func(o *sandbox.CreateOptions) { o.Resources.CPUCores = sandbox.Value(0.0001) }, "resources.cpu_cores"},
		{"CPU overflow", func(o *sandbox.CreateOptions) { o.Resources.CPUCores = sandbox.Value(float64(math.MaxUint32+1) / 1000) }, "cpu_cores cannot be represented"},
		{"fractional lifetime", func(o *sandbox.CreateOptions) { o.Lifetime.MaximumLifetime = sandbox.Value(time.Millisecond) }, "maximum_lifetime cannot be represented"},
		{"lifetime overflow", func(o *sandbox.CreateOptions) {
			o.Lifetime.MaximumLifetime = sandbox.Value(time.Duration(math.MaxUint32+1) * time.Second)
		}, "maximum_lifetime cannot be represented"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := createOptions()
			tc.modify(options)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			instance, err := client.Create(ctx, options)
			if instance != nil || err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %q, got instance=%v err=%v", tc.want, instance, err)
			}
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Provider != "modal" || detail.Operation != "create" || detail.Field == "" {
				t.Fatalf("missing common error details: %v", err)
			}
		})
	}
	if calls.Load() != 0 {
		t.Fatal("invalid requests reached provider endpoint")
	}
}

func TestUnsupportedClientConfiguration(t *testing.T) {
	config := clientConfig()
	config.Endpoint = sandbox.Value("https://example.invalid")
	// Endpoint is rejected independently of authentication; supply explicit dummy credentials.
	config.Auth = &sandbox.AuthConfig{TokenPair: &sandbox.TokenPairCredentials{ID: "dummy", Secret: "dummy"}}
	config.Scope = nil
	client, err := sandbox.NewClient(config)
	if client != nil || err == nil || !strings.Contains(err.Error(), "Endpoint") {
		t.Fatalf("expected unsupported endpoint, got %v", err)
	}
}
