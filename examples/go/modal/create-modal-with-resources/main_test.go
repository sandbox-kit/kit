package main

import (
	"context"
	"errors"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/providers/modal"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func testClientConfig() (sandbox.Provider, *sandbox.Config) {
	environment := os.Getenv("MODAL_ENVIRONMENT")
	if environment == "" {
		environment = "main"
	}
	return modal.New(), &sandbox.Config{
		Auth: &sandbox.AuthConfig{
			Credentials: &sandbox.TokenPair{
				TokenPair: &sandbox.TokenPairCredentials{
					Id:     os.Getenv("MODAL_TOKEN_ID"),
					Secret: os.Getenv("MODAL_TOKEN_SECRET"),
				},
			},
		},
		Scope: &sandbox.Scope{
			AppName:     sandbox.Value(os.Getenv("MODAL_APP_NAME")),
			Environment: sandbox.Value(environment),
		},
		Timeout: sandbox.Duration(2 * time.Minute),
	}
}

func testCreateOptions() *sandbox.CreateOptions {
	return &sandbox.CreateOptions{
		Source: &sandbox.SandboxSource{
			Kind: &sandbox.Image{Image: &sandbox.ImageSource{Reference: "alpine:3.21"}},
		},
		Resources: &sandbox.Resources{
			CpuCores:  sandbox.Value(0.5),
			MemoryMib: sandbox.Value(uint64(512)),
		},
		Lifetime: &sandbox.LifetimePolicy{MaximumLifetime: sandbox.Duration(5 * time.Minute)},
	}
}

func TestPublicClientRejectsLossyResourcesBeforeRemoteCalls(t *testing.T) {
	var calls atomic.Int32
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Skipf("local loopback listener unavailable: %v", err)
	}
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
		http.Error(w, "unexpected remote call", 500)
	}))
	server.Listener = listener
	server.Start()
	defer server.Close()
	t.Setenv("MODAL_SERVER_URL", server.URL)
	t.Setenv("MODAL_TOKEN_ID", "test-only-id")
	t.Setenv("MODAL_TOKEN_SECRET", "test-only-secret")
	t.Setenv("MODAL_APP_NAME", "test-app")
	t.Setenv("MODAL_ENVIRONMENT", "dev")
	provider, config := testClientConfig()
	if err := sandbox.ValidateCreateOptions(testCreateOptions()); err != nil {
		t.Fatalf("example options invalid: %v", err)
	}
	client, err := sandbox.NewClient(provider, config)
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
		{"below minimum CPU", func(o *sandbox.CreateOptions) { o.Resources.CpuCores = sandbox.Value(0.0001) }, "resources.cpu_cores"},
		{"CPU overflow", func(o *sandbox.CreateOptions) { o.Resources.CpuCores = sandbox.Value(float64(math.MaxUint32+1) / 1000) }, "cpu_cores cannot be represented"},
		{"fractional lifetime", func(o *sandbox.CreateOptions) { o.Lifetime.MaximumLifetime = sandbox.Duration(time.Millisecond) }, "maximum_lifetime cannot be represented"},
		{"lifetime overflow", func(o *sandbox.CreateOptions) {
			o.Lifetime.MaximumLifetime = sandbox.Duration(time.Duration(math.MaxUint32+1) * time.Second)
		}, "maximum_lifetime cannot be represented"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := testCreateOptions()
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
	provider, config := testClientConfig()
	config.Endpoint = sandbox.Value("https://example.invalid")
	// Endpoint is rejected independently of authentication; supply explicit dummy credentials.
	config.Auth = &sandbox.AuthConfig{Credentials: &sandbox.TokenPair{TokenPair: &sandbox.TokenPairCredentials{Id: "dummy", Secret: "dummy"}}}
	config.Scope = nil
	client, err := sandbox.NewClient(provider, config)
	if client != nil || err == nil || !strings.Contains(err.Error(), "Endpoint") {
		t.Fatalf("expected unsupported endpoint, got %v", err)
	}
}
