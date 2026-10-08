package core_test

import (
	"testing"

	"github.com/sandbox-kit/kit/core"
)

type testProvider struct {
	name string
}

func (p *testProvider) ProviderName() string { return p.name }

func TestProvidersUseTheSameClientType(t *testing.T) {
	for _, name := range []string{"modal", "daytona"} {
		provider := &testProvider{name: name}
		var client *core.Client
		var err error
		client, err = core.NewClient(provider)
		if err != nil {
			t.Fatal(err)
		}
		if client.ProviderName() != name {
			t.Fatal("incorrect selected provider")
		}
		provider.name = "updated"
		if client.ProviderName() != "updated" {
			t.Fatal("client replaced the injected provider")
		}
	}
}

func TestClientRejectsNilProviders(t *testing.T) {
	t.Run("nil interface", func(t *testing.T) {
		var provider core.Provider
		client, err := core.NewClient(provider)
		if err == nil || client != nil {
			t.Fatal("expected constructor error and no client")
		}
	})
	t.Run("typed nil", func(t *testing.T) {
		var provider *testProvider
		client, err := core.NewClient(provider)
		if err == nil || client != nil {
			t.Fatal("expected constructor error and no client")
		}
	})
	t.Run("typed nil inside interface", func(t *testing.T) {
		var pointer *testProvider
		var provider core.Provider = pointer
		client, err := core.NewClient(provider)
		if err == nil || client != nil {
			t.Fatal("expected constructor error and no client")
		}
	})
}

type valueProvider struct{}

func (valueProvider) ProviderName() string { panic("initialization must not invoke provider methods") }

func TestClientAcceptsValueProviderWithoutInvokingIt(t *testing.T) {
	if _, err := core.NewClient(valueProvider{}); err != nil {
		t.Fatal(err)
	}
}
