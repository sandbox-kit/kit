package daytona

import (
	"math"
	"testing"
	"time"

	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestActionSpecificPolicySemantics(t *testing.T) {
	for _, action := range []struct {
		name  string
		set   func(*sandbox.LifetimePolicy, *sandbox.AutomaticAction)
		get   func(types.SandboxBaseParams) *int
		pause bool
	}{
		{"stop", func(p *sandbox.LifetimePolicy, a *sandbox.AutomaticAction) { p.IdleStop = a }, func(p types.SandboxBaseParams) *int { return p.AutoStopInterval }, false},
		{"pause", func(p *sandbox.LifetimePolicy, a *sandbox.AutomaticAction) { p.IdlePause = a }, func(p types.SandboxBaseParams) *int { return p.AutoPauseInterval }, true},
		{"archive", func(p *sandbox.LifetimePolicy, a *sandbox.AutomaticAction) { p.StoppedArchive = a }, func(p types.SandboxBaseParams) *int { return p.AutoArchiveInterval }, false},
		{"delete", func(p *sandbox.LifetimePolicy, a *sandbox.AutomaticAction) { p.StoppedDelete = a }, func(p types.SandboxBaseParams) *int { return p.AutoDeleteInterval }, false},
	} {
		t.Run(action.name, func(t *testing.T) {
			for _, tc := range []struct {
				name   string
				policy *sandbox.AutomaticAction
				valid  bool
				value  *int
			}{
				{"absent", nil, true, nil},
				{"default", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDefault}, true, nil},
				{"disabled", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeDisabled}, action.pause || action.name == "stop", sandbox.Value(0)},
				{"immediate", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(0))}, action.name == "delete", sandbox.Value(0)},
				{"delayed", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(3 * time.Minute)}, true, sandbox.Value(3)},
				{"fractional", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Second)}, false, nil},
				{"negative", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(-time.Minute)}, false, nil},
				{"overflow", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter, After: sandbox.Value(time.Duration(math.MaxInt64))}, false, nil},
				{"missing delay", &sandbox.AutomaticAction{Mode: sandbox.PolicyModeAfter}, false, nil},
				{"unknown", &sandbox.AutomaticAction{Mode: sandbox.PolicyMode(99)}, false, nil},
			} {
				t.Run(tc.name, func(t *testing.T) {
					lifetime := &sandbox.LifetimePolicy{}
					action.set(lifetime, tc.policy)
					mapped, _, err := mapLifetimePolicies(lifetime)
					if (err == nil) != tc.valid {
						t.Fatalf("valid=%v err=%v", tc.valid, err)
					}
					if !tc.valid {
						return
					}
					value := action.get(mapped)
					if tc.value == nil {
						if value != nil {
							t.Fatal("absent/default policy became explicit")
						}
					} else if value == nil || *value != *tc.value {
						t.Fatal("incorrect native policy")
					}
					plan, err := planCreate(&sandbox.CreateOptions{Lifetime: lifetime})
					if err != nil {
						t.Fatal(err)
					}
					actual := action.get(plan.params.(types.SnapshotParams).SandboxBaseParams)
					if (actual == nil) != (value == nil) || actual != nil && *actual != *value {
						t.Fatal("creation plan lost policy mapping")
					}
				})
			}
		})
	}
}
