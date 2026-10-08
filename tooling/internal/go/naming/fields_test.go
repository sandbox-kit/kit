package naming

import "testing"

func TestGoInitialismsAndUnits(t *testing.T) {
	for input, want := range map[string]string{
		"cpu_cores": "CPUCores", "cpu_limit_cores": "CPULimitCores", "memory_mib": "MemoryMiB",
		"outbound_cidrs": "OutboundCIDRs", "outbound_proxy_url": "OutboundProxyURL",
		"tcp_port": "TCPPort", "pty": "PTY", "api_key": "APIKey", "oauth_refresh": "OAuthRefresh",
		"organization_id": "OrganizationID", "linux_vm": "LinuxVM", "microvm": "MicroVM",
	} {
		if got := identifier(input); got != want {
			t.Errorf("%s: got %s, want %s", input, got, want)
		}
	}
}
