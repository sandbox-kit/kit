package sandbox_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestErrorOwnsDetailsAndPreservesCause(t *testing.T) {
	status := uint32(429)
	code := "rate_limit"
	source := "api"
	info := sandbox.ErrorInfo{Kind: sandbox.ErrorKindRateLimited, Provider: "daytona", Operation: "create", Field: "resources", Message: "limited", StatusCode: &status, ProviderCode: &code, ProviderSource: &source}
	native := errors.New("native")
	err := sandbox.NewError(&info, native)
	status = 500
	code = "changed"
	source = "changed"
	info.Message = "changed"
	if err.GetStatusCode() != 429 || err.GetProviderCode() != "rate_limit" || err.GetProviderSource() != "api" || err.Message != "limited" || !errors.Is(err, native) {
		t.Fatal("constructor aliases details or loses native cause")
	}
}

func TestErrorContextAvoidsRedundantWrapping(t *testing.T) {
	native := errors.New("native")
	original := sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindAuthentication, Provider: "modal", Operation: "create", StatusCode: sandbox.Value(uint32(401))}, native)
	if sdk.WithErrorContext(original, "modal", "create") != original {
		t.Fatal("complete error was wrapped again")
	}
	outer := fmt.Errorf("outer: %w", original)
	if sdk.WithErrorContext(outer, "modal", "create") != outer {
		t.Fatal("outer cause chain replaced despite complete context")
	}
	changed := sdk.WithErrorContext(original, "daytona", "initialize")
	var detail *sandbox.Error
	if !errors.As(changed, &detail) || detail.Provider != "daytona" || detail.Operation != "initialize" || !errors.Is(changed, native) {
		t.Fatal("context enrichment lost native cause")
	}
	*detail.StatusCode = 500
	if original.Provider != "modal" || original.GetStatusCode() != 401 {
		t.Fatal("context enrichment mutated original details")
	}
}

func TestLocalErrorDetailsAndOriginalValidator(t *testing.T) {
	client, _ := newTestClient(&creationBackend{})
	_, err := client.Create(context.Background(), &sandbox.CreateOptions{Resources: &sandbox.Resources{CpuCores: sandbox.Value(-1.0)}})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != "resources.cpu_cores" || detail.Provider != "test" || detail.Operation != "create" || detail.GetRule() != "exclusive_minimum" {
		t.Fatalf("missing shared validation details: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.Create(ctx, nil)
	if !errors.Is(err, context.Canceled) || !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindCanceled {
		t.Fatal("cancellation cause lost")
	}
}

func TestResponseAvailability(t *testing.T) {
	response := &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{Id: "id", Name: sandbox.Value("requested")}}
	client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) { return response, nil }})
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	info := instance.Info()
	if info.Resources != nil || info.Region != nil || info.GetName() != "requested" {
		t.Fatal("missing fields became fabricated")
	}
	bad, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) { return nil, nil }})
	_, err = bad.Create(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidResponse {
		t.Fatal("invalid response not classified")
	}
}
