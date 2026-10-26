package modal

import (
	"net/netip"
	"regexp"
	"strings"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var providerEnvNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
var providerDomainPattern = regexp.MustCompile(`(?i)^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$`)

// validateProviderCreate rejects creation settings this provider cannot honor.
// It runs before the native create call.
//
// Rules are inventoried in specs/providers/modal.yaml checks and kept in sync
// here by hand (no longer emitted).
func validateProviderCreate(request *sandbox.CreateOptions) error {
	if request == nil {
		return nil
	}
	if resources := request.Resources; resources != nil {
		if resources.CpuCores != nil && *resources.CpuCores < 0.125 {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldResourcesCpuCores, "resources.cpu_cores is outside supported values")
		}
		if resources.MemoryMib != nil && *resources.MemoryMib < 128 {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldResourcesMemoryMib, "resources.memory_mib is outside supported values")
		}
	}
	if request.Runtime != nil && request.Runtime.WorkingDirectory != nil {
		if !strings.HasPrefix(*request.Runtime.WorkingDirectory, "/") {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldRuntimeWorkingDirectory, "runtime.working_directory is outside supported values")
		}
	}
	if request.Placement != nil && request.Placement.Cloud != nil {
		cloud := *request.Placement.Cloud
		if cloud != "aws" && cloud != "gcp" && cloud != "oci" && cloud != "auto" {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldPlacementCloud, "placement.cloud is outside supported values")
		}
	}
	if network := request.Network; network != nil {
		if network.OutboundCidrs != nil {
			for _, value := range network.OutboundCidrs.Entries {
				if value != "" {
					if _, err := netip.ParsePrefix(value); err != nil {
						return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkOutboundCidrsEntries, "network.outbound_cidrs.entries is outside supported values")
					}
				}
			}
		}
		if network.InboundCidrs != nil {
			for _, value := range network.InboundCidrs.Entries {
				if _, err := netip.ParsePrefix(value); err != nil {
					return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkInboundCidrsEntries, "network.inbound_cidrs.entries is outside supported values")
				}
			}
		}
		if network.OutboundDomains != nil {
			for _, value := range network.OutboundDomains.Entries {
				if value == "" {
					continue
				}
				domain := strings.TrimPrefix(value, "*.")
				if !providerDomainPattern.MatchString(domain) || len(domain) > 253 {
					return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkOutboundDomainsEntries, "network.outbound_domains.entries is outside supported values")
				}
			}
		}
	}
	for value := range request.Environment {
		if !providerEnvNamePattern.MatchString(value) {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldEnvironment, "environment is outside supported values")
		}
	}
	return nil
}

func providerCheckError(kind sandbox.ErrorKind, field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      kind,
		Provider:  "modal",
		Operation: "create",
		Field:     field,
		Message:   "sandbox-kit modal: " + message,
	}, nil)
}
