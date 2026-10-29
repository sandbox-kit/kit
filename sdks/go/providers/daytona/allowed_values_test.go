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
