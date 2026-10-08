package main

import (
	"bytes"
	"strings"
	"testing"

	daytonaSDK "github.com/daytona/clients/sdk-go/pkg/daytona"
	modalSDK "github.com/modal-labs/modal-client/go"
	daytonaAdapter "github.com/sandbox-kit/kit/sdks/go/adapters/daytona"
	modalAdapter "github.com/sandbox-kit/kit/sdks/go/adapters/modal"
	"github.com/sandbox-kit/kit/sdks/go/core"
)

func TestSDKClientsUseTheCommonClient(t *testing.T) {
	// Construct only SDK structs: adapter/core initialization must not perform
	// authentication or network work. The application supplies these SDK types.
	for _, test := range []struct {
		name        string
		newProvider func() (core.Provider, error)
	}{
		{"modal", func() (core.Provider, error) { return modalAdapter.New(&modalSDK.Client{}) }},
		{"daytona", func() (core.Provider, error) { return daytonaAdapter.New(&daytonaSDK.Client{}) }},
	} {
		t.Run(test.name, func(t *testing.T) {
			provider, err := test.newProvider()
			if err != nil {
				t.Fatal(err)
			}
			client, err := core.NewClient(provider)
			if err != nil {
				t.Fatal(err)
			}
			if client.ProviderName() != test.name {
				t.Fatal("incorrect provider selected")
			}
		})
	}
}

func TestHelpDoesNotRequireProviderCredentials(t *testing.T) {
	for _, args := range [][]string{nil, {"-h"}, {"--help"}} {
		var output bytes.Buffer
		if err := run(args, &output); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(output.String(), "<modal|daytona>") {
			t.Fatal("missing provider usage")
		}
	}
}

func TestInvalidArgumentsDoNotInitializeAnSDK(t *testing.T) {
	for _, args := range [][]string{{"unknown"}, {"modal", "daytona"}} {
		var output bytes.Buffer
		if err := run(args, &output); err == nil {
			t.Fatal("invalid provider selection accepted")
		}
		if output.Len() != 0 {
			t.Fatal("unexpected initialization output")
		}
	}
}
