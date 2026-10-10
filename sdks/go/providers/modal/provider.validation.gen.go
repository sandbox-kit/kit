// Code generated from provider validation specs. DO NOT EDIT.
package modal

import (
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
				return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldResourcesCPUCores, Message: "sandbox-kit modal: resources.cpu_cores is outside supported values"}, nil)
			}
		}
	}
	if request.Resources != nil {
		if request.Resources.MemoryMiB != nil {
			if *request.Resources.MemoryMiB < 128 {
				return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldResourcesMemoryMiB, Message: "sandbox-kit modal: resources.memory_mib is outside supported values"}, nil)
			}
		}
	}
	if request.Runtime != nil {
		if request.Runtime.WorkingDirectory != nil {
			if !strings.HasPrefix(*request.Runtime.WorkingDirectory, "/") {
				return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldRuntimeWorkingDirectory, Message: "sandbox-kit modal: runtime.working_directory is outside supported values"}, nil)
			}
		}
	}
	if request.Placement != nil {
		if request.Placement.Cloud != nil {
			if *request.Placement.Cloud != "aws" && *request.Placement.Cloud != "gcp" && *request.Placement.Cloud != "oci" && *request.Placement.Cloud != "auto" {
				return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindUnsupported, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldPlacementCloud, Message: "sandbox-kit modal: placement.cloud is outside supported values"}, nil)
			}
		}
	}
	if request.Network != nil {
		if request.Network.OutboundCIDRs != nil {
			for _, value := range request.Network.OutboundCIDRs.Entries {
				_, err := netip.ParsePrefix(value)
				if value != "" && err != nil {
					return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldNetworkOutboundCIDRsEntries, Message: "sandbox-kit modal: network.outbound_cidrs.entries is outside supported values"}, nil)
				}
			}
		}
	}
	if request.Network != nil {
		if request.Network.InboundCIDRs != nil {
			for _, value := range request.Network.InboundCIDRs.Entries {
				_, err := netip.ParsePrefix(value)
				if err != nil {
					return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldNetworkInboundCIDRsEntries, Message: "sandbox-kit modal: network.inbound_cidrs.entries is outside supported values"}, nil)
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
					return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldNetworkOutboundDomainsEntries, Message: "sandbox-kit modal: network.outbound_domains.entries is outside supported values"}, nil)
				}
			}
		}
	}
	for value := range request.Environment {
		matched := providerEnvNamePattern.MatchString(value)
		if !matched {
			return sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldEnvironment, Message: "sandbox-kit modal: environment is outside supported values"}, nil)
		}
	}
	return nil
}
