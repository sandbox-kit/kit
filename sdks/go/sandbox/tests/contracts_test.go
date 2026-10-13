package sandbox_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"math"
	"reflect"
	"testing"
)

func TestCopiedMetadataRetainsNumbersAndOwnership(t *testing.T) {
	big := uint64(math.MaxUint64)
	source := &sandbox.SandboxInfo{ID: "id", ProviderMetadata: map[string]any{"large": big, "signed": int64(9007199254740993), "nested": []any{map[string]any{"bytes": []byte{1, 2}}}}, Origins: map[string]sandbox.ValueOrigin{"id": sandbox.ValueOriginProvider}}
	copied, err := source.Clone()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(source, copied) {
		t.Fatal("types or values changed")
	}
	source.ProviderMetadata["large"] = uint64(0)
	source.ProviderMetadata["nested"].([]any)[0].(map[string]any)["bytes"].([]byte)[0] = 9
	source.Origins["id"] = sandbox.ValueOriginRequest
	if copied.ProviderMetadata["large"] != big || copied.ProviderMetadata["nested"].([]any)[0].(map[string]any)["bytes"].([]byte)[0] != 1 || copied.Origins["id"] != sandbox.ValueOriginProvider {
		t.Fatal("copy aliases original")
	}
}

func TestErrorOwnsDetailsAndPreservesCause(t *testing.T) {
	status := uint32(429)
	code := "rate_limit"
	source := "api"
	info := sandbox.ErrorInfo{Kind: sandbox.ErrorKindRateLimited, Provider: "daytona", Operation: "create", Field: "resources", Message: "limited", StatusCode: &status, ProviderCode: &code, ProviderSource: &source}
	native := errors.New("native")
	err := sandbox.NewError(info, native)
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
	original := sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindAuthentication, Provider: "modal", Operation: "create", StatusCode: sandbox.Value(uint32(401))}, native)
	if sandbox.WithErrorContext(original, "modal", "create") != original {
		t.Fatal("complete error was wrapped again")
	}
	outer := fmt.Errorf("outer: %w", original)
	if sandbox.WithErrorContext(outer, "modal", "create") != outer {
		t.Fatal("outer cause chain replaced despite complete context")
	}
	changed := sandbox.WithErrorContext(original, "daytona", "initialize")
	var detail *sandbox.Error
	if !errors.As(changed, &detail) || detail.Provider != "daytona" || detail.Operation != "initialize" || !errors.Is(changed, native) {
		t.Fatal("context enrichment lost native cause")
	}
	*detail.StatusCode = 500
	if original.Provider != "modal" || original.GetStatusCode() != 401 {
		t.Fatal("context enrichment mutated original details")
	}
}

func TestResponseOriginRequiresKnownPresentField(t *testing.T) {
	for _, tc := range []struct {
		name  string
		info  *sandbox.SandboxInfo
		valid bool
	}{
		{"absent name", &sandbox.SandboxInfo{Origins: map[string]sandbox.ValueOrigin{"name": sandbox.ValueOriginProvider}}, false},
		{"unknown path", &sandbox.SandboxInfo{Origins: map[string]sandbox.ValueOrigin{"misspelled": sandbox.ValueOriginProvider}}, false},
		{"missing nested value", &sandbox.SandboxInfo{Resources: &sandbox.Resources{}, Origins: map[string]sandbox.ValueOrigin{"resources.cpu_cores": sandbox.ValueOriginProvider}}, false},
		{"explicit empty name", &sandbox.SandboxInfo{Name: sandbox.Value(""), Origins: map[string]sandbox.ValueOrigin{"name": sandbox.ValueOriginProvider}}, true},
		{"empty map", &sandbox.SandboxInfo{Labels: map[string]string{}, Origins: map[string]sandbox.ValueOrigin{"labels": sandbox.ValueOriginProvider}}, true},
		{"explicit zero", &sandbox.SandboxInfo{Resources: &sandbox.Resources{CPUCores: sandbox.Value(0.0)}, Origins: map[string]sandbox.ValueOrigin{"resources.cpu_cores": sandbox.ValueOriginProvider}}, true},
		{"explicit false", &sandbox.SandboxInfo{Isolation: &sandbox.IsolationConfig{NestedVirtualization: sandbox.Value(false)}, Origins: map[string]sandbox.ValueOrigin{"isolation.nested_virtualization": sandbox.ValueOriginProvider}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tc.info.ID = "id"
			client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
				return &sandbox.CreateResult{Sandbox: tc.info}, nil
			}})
			_, err := client.Create(context.Background(), nil)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
		})
	}
}

func TestInvalidMetadataNeverReachesBackend(t *testing.T) {
	cyclic := map[string]any{}
	cyclic["self"] = cyclic
	for _, value := range []any{cyclic, math.NaN(), make(chan int)} {
		called := false
		client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
			called = true
			return nil, nil
		}})
		_, err := client.Create(context.Background(), &sandbox.CreateOptions{ProviderOptions: map[string]any{"value": value}})
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != "provider_options" || called {
			t.Fatalf("invalid metadata reached provider or was misclassified: %v", err)
		}
	}
}

func TestLocalErrorDetailsAndOriginalValidator(t *testing.T) {
	client, _ := newTestClient(&creationBackend{})
	_, err := client.Create(context.Background(), &sandbox.CreateOptions{Resources: &sandbox.Resources{CPUCores: sandbox.Value(-1.0)}})
	var detail *sandbox.Error
	var original validator.ValidationErrors
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != "resources.cpu_cores" || detail.Provider != "test" || detail.Operation != "create" || !errors.As(err, &original) {
		t.Fatalf("missing shared/native validation details: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.Create(ctx, nil)
	if !errors.Is(err, context.Canceled) || !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindCanceled {
		t.Fatal("cancellation cause lost")
	}
}

func TestResponseAvailabilityAndOrigins(t *testing.T) {
	response := &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "id", Name: sandbox.Value("requested"), Origins: map[string]sandbox.ValueOrigin{"name": sandbox.ValueOriginRequest}}}
	client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) { return response, nil }})
	instance, err := client.Create(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	info := instance.Info()
	if info.Resources != nil || info.Region != nil || info.Origins["name"] != sandbox.ValueOriginRequest {
		t.Fatal("missing fields became fabricated")
	}
	bad, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) { return nil, nil }})
	_, err = bad.Create(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidResponse {
		t.Fatal("invalid response not classified")
	}
}

func TestTypedMetadataCyclesAndUnknownResponseOrigins(t *testing.T) {
	value := &sandbox.MetadataValue{}
	value.List = &sandbox.MetadataList{Values: []*sandbox.MetadataValue{value}}
	if _, err := value.Clone(); err == nil {
		t.Fatal("typed metadata cycle accepted")
	}
	client, _ := newTestClient(&creationBackend{call: func(context.Context, *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
		return &sandbox.CreateResult{Sandbox: &sandbox.SandboxInfo{ID: "id", Origins: map[string]sandbox.ValueOrigin{"id": 99}}}, nil
	}})
	_, err := client.Create(context.Background(), nil)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidResponse {
		t.Fatal("unknown response origin accepted")
	}
}
