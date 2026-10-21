package sandbox

import (
	"errors"
	"testing"
	"time"
)

func TestCreateValidationRules(t *testing.T) {
	cases := []struct {
		name   string
		config *CreateOptions
		valid  bool
	}{
		{"absent options", &CreateOptions{}, true},
		{"explicit zero requested CPU", &CreateOptions{Resources: &Resources{CpuCores: Value(0.0)}}, false},
		{"fractional CPU", &CreateOptions{Resources: &Resources{CpuCores: Value(0.5)}}, true},
		{"zero hard limit", &CreateOptions{Resources: &Resources{CpuCores: Value(2.0), CpuLimitCores: Value(0.0)}}, true},
		{"limit below reservation", &CreateOptions{Resources: &Resources{CpuCores: Value(2.0), CpuLimitCores: Value(1.0)}}, false},
		{"blank image", &CreateOptions{Source: &SandboxSource{Kind: &Image{Image: &ImageSource{Reference: "  "}}}}, false},
		{"exclusive source", &CreateOptions{Source: &SandboxSource{Kind: &Image{Image: &ImageSource{Reference: "image"}}}}, true},
		{"after missing duration", &CreateOptions{Lifetime: &LifetimePolicy{IdleStop: &AutomaticAction{Mode: PolicyModeAfter}}}, false},
		{"immediate after", &CreateOptions{Lifetime: &LifetimePolicy{IdleStop: &AutomaticAction{Mode: PolicyModeAfter, After: Duration(0)}}}, true},
		{"disabled with duration", &CreateOptions{Lifetime: &LifetimePolicy{IdleStop: &AutomaticAction{Mode: PolicyModeDisabled, After: Duration(time.Minute)}}}, false},
		{"negative timeout", &CreateOptions{Provisioning: &ProvisioningOptions{Timeout: Duration(-time.Second)}}, false},
		{"nil port", &CreateOptions{Network: &NetworkConfig{Ports: []*PortBinding{nil}}}, false},
		{"out of range port", &CreateOptions{Network: &NetworkConfig{Ports: []*PortBinding{{Port: 65536}}}}, false},
		{"empty readiness", &CreateOptions{Readiness: &ReadinessProbe{}}, false},
		{"empty probe command", &CreateOptions{Readiness: &ReadinessProbe{Check: &Command{Command: &CommandProbe{}}}}, false},
		{"nil secret", &CreateOptions{Security: &SecurityConfig{Secrets: []*SecretReference{nil}}}, false},
		{"blank secret", &CreateOptions{Security: &SecurityConfig{Secrets: []*SecretReference{{Reference: "  "}}}}, false},
		{"deny all explicit empty allowlist", &CreateOptions{Network: &NetworkConfig{Egress: EgressModeRestricted, OutboundDomains: &StringAllowlist{}}}, true},
		{"restricted missing allowlist", &CreateOptions{Network: &NetworkConfig{Egress: EgressModeRestricted}}, false},
		{"block all conflicting allowlist", &CreateOptions{Network: &NetworkConfig{Egress: EgressModeBlockAll, OutboundDomains: &StringAllowlist{}}}, false},
		{"block all explicit false private network", &CreateOptions{Network: &NetworkConfig{Egress: EgressModeBlockAll, PrivateNetwork: Value(false)}}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateCreateOptions(tc.config)
			if (err == nil) != tc.valid {
				t.Fatalf("valid=%v, error=%v", tc.valid, err)
			}
		})
	}
}

func TestValidationErrorsExposeStableRule(t *testing.T) {
	err := ValidateCreateOptions(&CreateOptions{Resources: &Resources{CpuCores: Value(0.0)}})
	var detail *Error
	if !errors.As(err, &detail) || detail.Rule == nil || detail.GetRule() != "exclusive_minimum" {
		t.Fatalf("validation rule = %v, error = %v", detail, err)
	}
}
