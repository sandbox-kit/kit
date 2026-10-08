package daytona

import (
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/core"
)

// Only the application imports the real SDK. The adapter treats it as opaque.
type exampleSDKClient struct{ token string }

func TestNewAndInjectRetainsNativeClient(t *testing.T) {
	// No SDK initialization or credentials: wrapping must not invoke the SDK.
	native := &exampleSDKClient{token: "application-owned"}
	adapter, err := New(native)
	if err != nil {
		t.Fatal(err)
	}
	client, err := core.NewClient(adapter)
	if err != nil {
		t.Fatal(err)
	}
	if adapter.client != native {
		t.Fatal("adapter replaced the borrowed SDK client")
	}
	if client.ProviderName() != "daytona" {
		t.Fatal("incorrect provider identity")
	}
}

func TestNewRejectsNilSDKClient(t *testing.T) {
	adapter, err := New[exampleSDKClient](nil)
	if err == nil || adapter != nil {
		t.Fatal("expected constructor error and no adapter")
	}
}
