package daytona

import (
	"errors"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestDocumentedCreationRules(t *testing.T) {
	delay := func(d time.Duration) *sandbox.AutomaticAction {
		return &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Duration(d)}
	}
	for _, tc := range []struct {
		name    string
		request *sandbox.CreateOptions
		valid   bool
	}{
		{"default", &sandbox.CreateOptions{}, true},
		{"language", &sandbox.CreateOptions{Runtime: &sandbox.RuntimeConfig{Language: sandbox.Value("python")}}, true},
		{"unsupported language", &sandbox.CreateOptions{Runtime: &sandbox.RuntimeConfig{Language: sandbox.Value("ruby")}}, false},
		{"both idle actions", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{IdleStop: delay(time.Minute), IdlePause: delay(time.Minute)}}, false},
		{"ephemeral pause", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true), IdlePause: delay(time.Minute)}}, false},
		{"ephemeral delayed delete", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true), StoppedDelete: delay(time.Minute)}}, false},
		{"ephemeral archive", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true), StoppedArchive: delay(time.Minute)}}, false},
		{"disabled stop", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{IdleStop: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}}}, true},
		{"archive maximum", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{StoppedArchive: delay(30 * 24 * time.Hour)}}, true},
		{"archive above maximum", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{StoppedArchive: delay(30*24*time.Hour + time.Minute)}}, false},
		{"immediate archive", &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{StoppedArchive: delay(0)}}, false},
		{"bad CIDR", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{OutboundCidrs: &sandbox.StringAllowlist{Entries: []string{"not-a-cidr"}}}}, false},
		{"linked sandbox is ephemeral", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{LinkedSandbox: sandbox.Value("parent")}, Lifetime: &sandbox.LifetimePolicy{Ephemeral: sandbox.Value(true)}}, true},
		{"linked sandbox without ephemeral lifetime", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{LinkedSandbox: sandbox.Value("parent")}}, false},
		{"placement", &sandbox.CreateOptions{Placement: &sandbox.Placement{Regions: []string{"eu"}}}, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planCreate(tc.request)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
			if tc.name == "placement" {
				var detail *sandbox.Error
				if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != sandbox.CreateFieldPlacement {
					t.Fatalf("placement error = %#v", detail)
				}
			}
		})
	}
}

func TestUnsupportedFieldsUseCanonicalPaths(t *testing.T) {
	for _, tc := range []struct {
		name    string
		request *sandbox.CreateOptions
		field   string
	}{
		{
			name:    "idle terminate",
			request: &sandbox.CreateOptions{Lifetime: &sandbox.LifetimePolicy{IdleTerminate: &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}}},
			field:   sandbox.CreateFieldLifetimeIdleTerminate,
		},
		{
			name:    "verbose",
			request: &sandbox.CreateOptions{Observability: &sandbox.ObservabilityConfig{Verbose: sandbox.Value(true)}},
			field:   sandbox.CreateFieldObservabilityVerbose,
		},
		{
			name:    "workload identity",
			request: &sandbox.CreateOptions{Security: &sandbox.SecurityConfig{WorkloadIdentity: sandbox.Value(true)}},
			field:   sandbox.CreateFieldSecurityWorkloadIdentity,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := planCreate(tc.request)
			var detail *sandbox.Error
			if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != tc.field {
				t.Fatalf("error = %v detail = %#v", err, detail)
			}
		})
	}
}
