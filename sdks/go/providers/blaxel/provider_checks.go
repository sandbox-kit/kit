package blaxel

import (
	"regexp"
	"strings"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var providerEnvNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Blaxel Metadata.Name: lowercase alphanumeric with single hyphens, max 49.
var providerSandboxNamePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Matches Daytona/Modal outbound domain checks.
var providerDomainPattern = regexp.MustCompile(`(?i)^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$`)

const providerSandboxNameMaxLen = 49

// validateProviderCreate rejects creation settings this provider cannot honor.
// Inventoried in specs/providers/blaxel.yaml checks; kept in sync by hand.
func validateProviderCreate(request *sandbox.CreateOptions) error {
	if request == nil {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldName, "name is required")
	}
	name := request.GetName()
	if name == "" {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldName, "name is required")
	}
	if len(name) > providerSandboxNameMaxLen || !providerSandboxNamePattern.MatchString(name) {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldName, "name must be lowercase alphanumeric with single hyphens, at most 49 characters")
	}
	if request.GetSource().GetImage() == nil || strings.TrimSpace(request.GetSource().GetImage().Reference) == "" {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldSourceImage, "an explicit image source is required")
	}
	if request.GetSource().GetSnapshot() != nil || request.GetSource().GetWarmPool() != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSource, "only an image source is supported")
	}
	if request.Placement != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldPlacement, "per-request Blaxel placement is unsupported; use Config.region")
	}
	if err := rejectUnsupportedResources(request.Resources); err != nil {
		return err
	}
	for name := range request.Environment {
		if !providerEnvNamePattern.MatchString(name) {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldEnvironment, "environment is outside supported values")
		}
	}
	if err := rejectInvalidOutboundDomains(request.Network); err != nil {
		return err
	}
	return nil
}

func rejectInvalidOutboundDomains(network *sandbox.NetworkConfig) error {
	if network == nil || network.OutboundDomains == nil {
		return nil
	}
	for _, value := range network.OutboundDomains.Entries {
		domain := strings.TrimPrefix(value, "*.")
		if !providerDomainPattern.MatchString(domain) || len(domain) > 253 {
			return providerCheckError(
				sandbox.ErrorKindInvalidArgument,
				sandbox.CreateFieldNetworkOutboundDomainsEntries,
				"network.outbound_domains.entries is outside supported values",
			)
		}
	}
	return nil
}

func rejectUnsupportedResources(resources *sandbox.Resources) error {
	if resources == nil {
		return nil
	}
	for _, check := range []struct {
		set   bool
		field string
		msg   string
	}{
		{resources.CpuCores != nil, sandbox.CreateFieldResourcesCpuCores, "Blaxel derives CPU from memory; cpu_cores is not mapped"},
		{resources.CpuLimitCores != nil, sandbox.CreateFieldResourcesCpuLimitCores, "Blaxel resource limits are not mapped"},
		{resources.MemoryLimitMib != nil, sandbox.CreateFieldResourcesMemoryLimitMib, "Blaxel resource limits are not mapped"},
		{resources.DiskMib != nil, sandbox.CreateFieldResourcesDiskMib, "Blaxel disk size is not mapped as MiB"},
	} {
		if check.set {
			return providerCheckError(sandbox.ErrorKindUnsupported, check.field, check.msg)
		}
	}
	return nil
}

func providerCheckError(kind sandbox.ErrorKind, field, message string) error {
	return kitError(kitFault{
		Kind:      kind,
		Operation: "create",
		Field:     field,
		Message:   message,
	})
}
