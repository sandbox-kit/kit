package blaxel

import (
	"regexp"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var providerEnvNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// Blaxel Metadata.Name: lowercase alphanumeric with single hyphens, max 49.
var providerSandboxNamePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

const providerSandboxNameMaxLen = 49

// validateProviderCreate rejects creation settings this provider cannot honor.
// It runs before the native create call.
//
// Rules are inventoried in specs/providers/blaxel.yaml checks and kept in sync
// here by hand (not emitted).
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
	if request.GetSource().GetImage() == nil || request.GetSource().GetImage().Reference == "" {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldSourceImage, "an explicit image source is required")
	}
	if request.GetSource().GetSnapshot() != nil || request.GetSource().GetWarmPool() != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSource, "only an image source is supported")
	}
	if request.Placement != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldPlacement, "per-request Blaxel placement is unsupported; use Config.region")
	}
	if resources := request.Resources; resources != nil {
		if resources.CpuCores != nil {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldResourcesCpuCores, "Blaxel derives CPU from memory; cpu_cores is not mapped")
		}
		if resources.CpuLimitCores != nil {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldResourcesCpuLimitCores, "Blaxel resource limits are not mapped")
		}
		if resources.MemoryLimitMib != nil {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldResourcesMemoryLimitMib, "Blaxel resource limits are not mapped")
		}
		if resources.DiskMib != nil {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldResourcesDiskMib, "Blaxel disk size is not mapped as MiB")
		}
	}
	for name := range request.Environment {
		if !providerEnvNamePattern.MatchString(name) {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldEnvironment, "environment is outside supported values")
		}
	}
	return nil
}

func providerCheckError(kind sandbox.ErrorKind, field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      kind,
		Provider:  "blaxel",
		Operation: "create",
		Field:     field,
		Message:   "sandbox-kit blaxel: " + message,
	}, nil)
}
