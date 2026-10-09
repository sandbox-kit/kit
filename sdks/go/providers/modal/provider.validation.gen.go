// Code generated from provider validation specs. DO NOT EDIT.
package modal

import (
	fmt "fmt"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
	netip "net/netip"
	regexp "regexp"
	strings "strings"
)

var providerEnvNamePattern = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
var providerDomainPattern = regexp.MustCompile("(?i)^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.?$")

func validateProviderCreate(request *sandbox.CreateOptions) error {
	if request == nil {
		return nil
	}
	if request.Resources != nil {
		if request.Resources.CPUCores != nil {
			if *request.Resources.CPUCores < 0.125 {
				return fmt.Errorf("sandbox-kit modal: resources.cpu_cores is outside supported values")
			}
		}
	}
	if request.Resources != nil {
		if request.Resources.MemoryMiB != nil {
			if *request.Resources.MemoryMiB < 128 {
				return fmt.Errorf("sandbox-kit modal: resources.memory_mib is outside supported values")
			}
		}
	}
	if request.Runtime != nil {
		if request.Runtime.WorkingDirectory != nil {
			if !strings.HasPrefix(*request.Runtime.WorkingDirectory, "/") {
				return fmt.Errorf("sandbox-kit modal: runtime.working_directory is outside supported values")
			}
		}
	}
	if request.Placement != nil {
		if request.Placement.Cloud != nil {
			if *request.Placement.Cloud != "aws" && *request.Placement.Cloud != "gcp" && *request.Placement.Cloud != "oci" && *request.Placement.Cloud != "auto" {
				return fmt.Errorf("sandbox-kit modal: placement.cloud is outside supported values")
			}
		}
	}
	if request.Network != nil {
		if request.Network.OutboundCIDRs != nil {
			for _, value := range request.Network.OutboundCIDRs.Entries {
				_, err := netip.ParsePrefix(value)
				if value != "" && err != nil {
					return fmt.Errorf("sandbox-kit modal: network.outbound_cidrs.entries is outside supported values")
				}
			}
		}
	}
	if request.Network != nil {
		if request.Network.InboundCIDRs != nil {
			for _, value := range request.Network.InboundCIDRs.Entries {
				_, err := netip.ParsePrefix(value)
				if err != nil {
					return fmt.Errorf("sandbox-kit modal: network.inbound_cidrs.entries is outside supported values")
				}
			}
		}
	}
	if request.Network != nil {
		if request.Network.OutboundDomains != nil {
			for _, value := range request.Network.OutboundDomains.Entries {
				domain := strings.TrimPrefix(value, "*.")
				matched := providerDomainPattern.MatchString(domain)
				if value != "" && (!matched || len(domain) > 253) {
					return fmt.Errorf("sandbox-kit modal: network.outbound_domains.entries is outside supported values")
				}
			}
		}
	}
	for value := range request.Environment {
		matched := providerEnvNamePattern.MatchString(value)
		if !matched {
			return fmt.Errorf("sandbox-kit modal: environment is outside supported values")
		}
	}
	return nil
}
