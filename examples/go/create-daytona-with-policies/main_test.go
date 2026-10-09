package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

const createdSandbox = `{"organizationId":"org","user":"daytona","env":{},"public":false,"networkBlockAll":false,"kvm":false,"gpu":0,"id":"example-created","name":"test","state":"started","target":"test-region","cpu":2,"memory":4,"disk":8,"labels":{},"toolboxProxyUrl":"http://127.0.0.1:1"}`

func TestExampleRunUsesOfficialSDKAndMappedPolicies(t *testing.T) {
	requests := make(chan map[string]any, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/sandbox" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("Authorization") != "Bearer test-only-key" {
			t.Error("API key not mapped")
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		requests <- body
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(createdSandbox))
	}))
	defer server.Close()
	t.Setenv("DAYTONA_API_KEY", "test-only-key")
	t.Setenv("DAYTONA_API_URL", server.URL)
	t.Setenv("DAYTONA_TARGET", "test-region")
	t.Chdir(t.TempDir())
	if err := os.WriteFile(".env", []byte("# Test uses explicit dummy environment values.\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := run(); err != nil {
		t.Fatal(err)
	}
	select {
	case body := <-requests:
		if body["autoPauseInterval"] != float64(0) || body["autoDeleteInterval"] != float64(10) {
			t.Fatalf("incorrect policies sent to SDK endpoint: %v", body)
		}
		if _, ok := body["autoStopInterval"]; ok {
			t.Fatal("unspecified auto-stop became explicit")
		}
	case <-time.After(time.Second):
		t.Fatal("SDK never sent creation request")
	}
}

func TestRejectedPoliciesNeverReachCreateEndpoint(t *testing.T) {
	var creates atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			creates.Add(1)
		}
		http.NotFound(w, r)
	}))
	defer server.Close()
	t.Setenv("DAYTONA_API_KEY", "test-only-key")
	t.Setenv("DAYTONA_API_URL", server.URL)
	t.Setenv("DAYTONA_TARGET", "test-region")
	client, err := sandbox.NewClient(clientConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close(context.Background())
	for _, tc := range []struct {
		name string
		set  func(*sandbox.LifetimePolicy)
		want string
	}{
		{"immediate stop", func(p *sandbox.LifetimePolicy) {
			p.IdleStop = &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(0))}
		}, "idle_stop"},
		{"disabled archive", func(p *sandbox.LifetimePolicy) {
			p.StoppedArchive = &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}
		}, "stopped_archive"},
		{"disabled delete", func(p *sandbox.LifetimePolicy) {
			p.StoppedDelete = &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}
		}, "stopped_delete"},
		{"immediate pause", func(p *sandbox.LifetimePolicy) {
			p.IdlePause = &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(0))}
		}, "idle_pause"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			options := createOptions()
			tc.set(options.Lifetime)
			instance, err := client.Create(context.Background(), options)
			if instance != nil || err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("expected %s rejection, got %v", tc.want, err)
			}
		})
	}
	if creates.Load() != 0 {
		t.Fatal("invalid policy triggered creation")
	}
}

func TestImmediateActionsRemainExplicitOnWire(t *testing.T) {
	requests := make(chan map[string]any, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/sandbox" {
			http.NotFound(w, r)
			return
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		requests <- body
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(createdSandbox))
	}))
	defer server.Close()
	t.Setenv("DAYTONA_API_KEY", "test-only-key")
	t.Setenv("DAYTONA_API_URL", server.URL)
	t.Setenv("DAYTONA_TARGET", "test-region")
	client, err := sandbox.NewClient(clientConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close(context.Background())
	for _, tc := range []struct {
		field string
		set   func(*sandbox.LifetimePolicy, *sandbox.AutomaticAction)
	}{
		{"autoDeleteInterval", func(p *sandbox.LifetimePolicy, a *sandbox.AutomaticAction) { p.StoppedDelete = a }},
	} {
		t.Run(tc.field, func(t *testing.T) {
			lifetime := &sandbox.LifetimePolicy{}
			tc.set(lifetime, &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(0))})
			if _, err := client.Create(context.Background(), &sandbox.CreateOptions{Lifetime: lifetime}); err != nil {
				t.Fatal(err)
			}
			select {
			case body := <-requests:
				if value, ok := body[tc.field]; !ok || value != float64(0) {
					t.Fatalf("explicit zero lost: %v", body)
				}
			case <-time.After(time.Second):
				t.Fatal("missing request")
			}
		})
	}
}
