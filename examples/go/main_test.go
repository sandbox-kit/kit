package main

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestExamplesUseCommonConfiguration(t *testing.T) {
	for _, name := range []string{"modal", "daytona"} {
		config := clientConfig(name, settings{timeout: time.Second})
		if config.Provider.Name() != name || config.Timeout == nil {
			t.Fatal("missing provider or common settings")
		}
	}
	t.Setenv("MODAL_TOKEN_ID", "test-id")
	t.Setenv("MODAL_TOKEN_SECRET", "test-secret")
	modal := clientConfig("modal", settings{app: "my-app"})
	if modal.Auth.TokenPair.ID != "test-id" || modal.Scope.GetAppName() != "my-app" {
		t.Fatal("Modal common config missing")
	}
	t.Setenv("DAYTONA_API_KEY", "test-key")
	daytona := clientConfig("daytona", settings{endpoint: "https://example.invalid", region: "region"})
	if daytona.Auth.APIKey.Key != "test-key" || daytona.GetRegion() != "region" || daytona.GetEndpoint() != "https://example.invalid" {
		t.Fatal("Daytona common config missing")
	}
}
func TestHelpDoesNotInitializeSDK(t *testing.T) {
	for _, args := range [][]string{nil, {"--help"}, {"modal", "--help"}} {
		var output bytes.Buffer
		if err := run(args, &output); err != nil {
			t.Fatal(err)
		}
		if output.Len() == 0 {
			t.Fatal("missing help")
		}
	}
}
func TestCreationDemo(t *testing.T) {
	var out bytes.Buffer
	if err := run([]string{"demo"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Simulated sandbox created: demo-sandbox (demo)") {
		t.Fatal(out.String())
	}
}
func TestInvalidArgumentsDoNotInitializeSDK(t *testing.T) {
	for _, args := range [][]string{{"unknown"}, {"modal", "daytona"}, {"modal", "--bad-flag"}} {
		var output bytes.Buffer
		if err := run(args, &output); err == nil {
			t.Fatal("invalid arguments accepted")
		}
	}
}

var _ sandbox.Provider = (*demoProvider)(nil)
var _ sandbox.Backend = (*demoBackend)(nil)
