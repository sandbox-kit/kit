package daytona

import (
	"testing"
	"time"

	apiclient "github.com/daytona/clients/api-client-go"
	daytonasdk "github.com/daytona/clients/sdk-go/pkg/daytona"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestNativeTimeout(t *testing.T) {
	if nativeTimeout(nil) != 0 {
		t.Fatal("omitted timeout")
	}
	if nativeTimeout(sandbox.Duration(0)) != 0 {
		t.Fatal("explicit zero timeout")
	}
	if nativeTimeout(sandbox.Duration(30*time.Second)) != 30*time.Second {
		t.Fatal("timeout duration")
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
			if info.Origins[sandbox.InfoFieldIsolation] != sandbox.ValueOriginProvider {
				t.Fatal("resolved class origin missing")
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
