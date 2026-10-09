package modal

import (
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"testing"
)

func TestDocumentedProviderChecks(t *testing.T) {
	for _, tc := range []struct {
		name    string
		request *sandbox.CreateOptions
		valid   bool
	}{
		{"defaults", nil, true},
		{"CPU minimum", &sandbox.CreateOptions{Resources: &sandbox.Resources{CPUCores: sandbox.Value(0.125)}}, true},
		{"CPU too small", &sandbox.CreateOptions{Resources: &sandbox.Resources{CPUCores: sandbox.Value(0.124)}}, false},
		{"memory minimum", &sandbox.CreateOptions{Resources: &sandbox.Resources{MemoryMiB: sandbox.Value(uint64(128))}}, true},
		{"memory too small", &sandbox.CreateOptions{Resources: &sandbox.Resources{MemoryMiB: sandbox.Value(uint64(127))}}, false},
		{"cloud", &sandbox.CreateOptions{Placement: &sandbox.Placement{Cloud: sandbox.Value("aws")}}, true},
		{"unknown cloud", &sandbox.CreateOptions{Placement: &sandbox.Placement{Cloud: sandbox.Value("random")}}, false},
		{"relative path", &sandbox.CreateOptions{Runtime: &sandbox.RuntimeConfig{WorkingDirectory: sandbox.Value("relative")}}, false},
		{"absolute path", &sandbox.CreateOptions{Runtime: &sandbox.RuntimeConfig{WorkingDirectory: sandbox.Value("/app")}}, true},
		{"CIDR", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{OutboundCIDRs: &sandbox.StringAllowlist{Entries: []string{"10.0.0.0/8", "::/0"}}}}, true},
		{"invalid CIDR", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{OutboundCIDRs: &sandbox.StringAllowlist{Entries: []string{"10.0.0.0/99"}}}}, false},
		{"domain", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{OutboundDomains: &sandbox.StringAllowlist{Entries: []string{"*.example.com"}}}}, true},
		{"domain injection", &sandbox.CreateOptions{Network: &sandbox.NetworkConfig{OutboundDomains: &sandbox.StringAllowlist{Entries: []string{"good.com,bad.com"}}}}, false},
		{"invalid env name", &sandbox.CreateOptions{Environment: map[string]string{"1BAD": "value"}}, false},
		{"env name", &sandbox.CreateOptions{Environment: map[string]string{"MODE": "value"}}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := validateProviderCreate(tc.request); (err == nil) != tc.valid {
				t.Fatalf("valid=%v err=%v", tc.valid, err)
			}
		})
	}
}
