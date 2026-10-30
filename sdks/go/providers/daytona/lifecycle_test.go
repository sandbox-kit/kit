package daytona

import (
	"errors"
	"testing"
	"time"

	apiclient "github.com/daytona/clients/api-client-go"
	daytonasdk "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/options"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestNativeTimeout(t *testing.T) {
	if got, err := nativeTimeout("start", nil); err != nil || got != daytonaDefaultTimeout {
		t.Fatalf("omitted timeout = %v, %v", got, err)
	}
	if got, err := nativeTimeout("start", sandbox.Duration(0)); err != nil || got != 0 {
		t.Fatalf("explicit zero timeout = %v, %v", got, err)
	}
	if got, err := nativeTimeout("start", sandbox.Duration(30*time.Second)); err != nil || got != 30*time.Second {
		t.Fatalf("timeout duration = %v, %v", got, err)
	}
	_, err := nativeTimeout("start", sandbox.Duration(1500*time.Millisecond))
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.ConfigFieldTimeout || detail.Operation != "start" {
		t.Fatalf("fractional timeout = %#v", detail)
	}
}

func TestPlanCreateTimeout(t *testing.T) {
	twoSeconds := 2 * time.Second
	for _, tc := range []struct {
		name    string
		timeout *time.Duration
		want    *time.Duration
	}{
		{"omitted", nil, nil},
		{"explicit zero", new(time.Duration), nil},
		{"two seconds", &twoSeconds, &twoSeconds},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := &sandbox.CreateOptions{}
			if tc.timeout != nil {
				request.Provisioning = &sandbox.ProvisioningOptions{Timeout: sandbox.Duration(*tc.timeout)}
			}
			plan, err := planCreate(request)
			if err != nil {
				t.Fatal(err)
			}
			native := &options.CreateSandbox{}
			for _, apply := range plan.options {
				apply(native)
			}
			if tc.want == nil {
				if native.Timeout != nil {
					t.Fatalf("timeout = %v", *native.Timeout)
				}
				return
			}
			if native.Timeout == nil || *native.Timeout != *tc.want {
				t.Fatalf("timeout = %v", native.Timeout)
			}
		})
	}
}

func TestResolvedSandboxClassMapping(t *testing.T) {
	for _, tc := range []struct {
		name         string
		class        types.SandboxClass
		gpu          int32
		kind         sandbox.IsolationKind
		providerKind string
	}{
		{"container", types.SandboxClassContainer, 0, sandbox.IsolationKindContainer, ""},
		{"linux vm", types.SandboxClassLinuxVM, 0, sandbox.IsolationKindLinuxVM, ""},
		{"windows", types.SandboxClass("windows"), 0, sandbox.IsolationKindProviderSpecific, "windows"},
		{"gpu", types.SandboxClassContainer, 1, sandbox.IsolationKindContainer, ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			instance := &daytonasdk.Sandbox{ID: "id", SandboxClass: &tc.class, Gpu: tc.gpu}
			info, err := sandboxInfo(instance)
			if err != nil {
				t.Fatal(err)
			}
			if info.GetIsolation().Kind != tc.kind || info.GetIsolation().GetProviderKind() != tc.providerKind {
				t.Fatalf("unexpected isolation: %#v", info.Isolation)
			}
			if tc.gpu > 0 {
				metadata, err := sandbox.MetadataObjectAsMap(info.GetProviderMetadata())
				if err != nil || metadata["gpu_count"] != int64(1) {
					t.Fatal("GPU metadata missing")
				}
			}
		})
	}
}

func TestResolvedPortableState(t *testing.T) {
	container := types.SandboxClassContainer
	instance := &daytonasdk.Sandbox{ID: "id", SandboxClass: &container, State: apiclient.SANDBOXSTATE_STARTED}
	info, err := sandboxInfo(instance)
	if err != nil {
		t.Fatal(err)
	}
	if info.State != sandbox.SandboxStateRunning {
		t.Fatalf("unexpected state: %v", info.State)
	}
}
