package blaxel

import (
	"fmt"
	"math"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/blaxel-ai/sdk-go/shared"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

const (
	// Blaxel measures memory in decimal MB while Kit exposes binary MiB.
	// 1 MiB = 1.048576 MB = 131072 / 125000 MB.
	blaxelMegabytesPerMebibyteNumerator   uint64 = 131072
	blaxelMegabytesPerMebibyteDenominator uint64 = 125000
)

// mapCommonCreateFields maps name, environment, and labels into the native
// create request and returns the fields that were not mapped.
func mapCommonCreateFields(source *sandbox.CreateOptions) (blaxel.SandboxNewParams, *sandbox.CreateOptions, error) {
	target := blaxel.SandboxNewParams{}
	if source == nil {
		return target, &sandbox.CreateOptions{}, nil
	}
	remaining := cloneMessage(source)
	if source.Name != nil {
		target.Sandbox.Metadata.Name = source.GetName()
	}
	remaining.Name = nil
	if len(source.Environment) > 0 {
		envs := make([]shared.EnvParam, 0, len(source.Environment))
		for name, value := range source.Environment {
			envs = append(envs, shared.EnvParam{
				Name:  blaxel.String(name),
				Value: blaxel.String(value),
			})
		}
		target.Sandbox.Spec.Runtime.Envs = envs
	}
	remaining.Environment = nil
	if source.Labels != nil {
		target.Sandbox.Metadata.Labels = source.Labels
	}
	remaining.Labels = nil
	return target, remaining, nil
}

// mapResources maps supplied memory into the native request and returns the
// fields that were not mapped.
func mapResources(source *sandbox.Resources) (blaxel.SandboxRuntimeParam, *sandbox.Resources, error) {
	target := blaxel.SandboxRuntimeParam{}
	if source == nil {
		return target, &sandbox.Resources{}, nil
	}
	remaining := cloneMessage(source)
	if source.MemoryMib != nil {
		value := source.GetMemoryMib()
		if value == 0 || value > math.MaxInt64 {
			return target, remaining, providerCheckError(
				sandbox.ErrorKindInvalidArgument,
				sandbox.CreateFieldResourcesMemoryMib,
				"memory_mib cannot be represented",
			)
		}
		memory, ok := mibToBlaxelMB(value)
		if !ok {
			return target, remaining, providerCheckError(
				sandbox.ErrorKindInvalidArgument,
				sandbox.CreateFieldResourcesMemoryMib,
				"memory_mib cannot be represented",
			)
		}
		target.Memory = blaxel.Int(memory)
	}
	remaining.MemoryMib = nil
	return target, remaining, nil
}

// mapSandboxInfo maps native fields into shared metadata. It does not invent
// values the SDK omitted.
func mapSandboxInfo(source *blaxel.Sandbox) (*sandbox.SandboxInfo, error) {
	target := &sandbox.SandboxInfo{Provider: "blaxel"}
	if source == nil {
		return target, nil
	}
	name := source.Metadata.Name
	target.Id = name
	if name != "" {
		target.Name = sandbox.Value(name)
	}
	if source.Status != "" {
		target.ProviderState = sandbox.Value(string(source.Status))
	} else if source.State != "" {
		target.ProviderState = sandbox.Value(string(source.State))
	}
	if source.Spec.Region != "" {
		target.Region = sandbox.Value(source.Spec.Region)
	}
	if source.Metadata.Labels != nil {
		target.Labels = source.Metadata.Labels
	}
	return target, nil
}

// mapAllocatedResources maps native memory into shared metadata.
func mapAllocatedResources(source *blaxel.Sandbox) (*sandbox.Resources, error) {
	target := &sandbox.Resources{}
	if source == nil || source.Spec.Runtime.Memory <= 0 {
		return target, nil
	}
	target.MemoryMib = sandbox.Value(blaxelMBToMib(uint64(source.Spec.Runtime.Memory)))
	return target, nil
}

// mibToBlaxelMB rounds up so a Kit request is never under-provisioned.
func mibToBlaxelMB(mib uint64) (int64, bool) {
	if mib == 0 || mib > math.MaxInt64/blaxelMegabytesPerMebibyteNumerator {
		return 0, false
	}
	mb := (mib*blaxelMegabytesPerMebibyteNumerator + blaxelMegabytesPerMebibyteDenominator - 1) / blaxelMegabytesPerMebibyteDenominator
	if mb > math.MaxInt64 {
		return 0, false
	}
	return int64(mb), true
}

// blaxelMBToMib reports the whole MiB guaranteed by a native decimal-MB value.
func blaxelMBToMib(mb uint64) uint64 {
	// Divide first so a large, but valid, native int64 memory value cannot
	// overflow before conversion. The remainder product stays below
	// numerator*denominator, and the quotient product remains below uint64's
	// maximum.
	whole := mb / blaxelMegabytesPerMebibyteNumerator
	remainder := mb % blaxelMegabytesPerMebibyteNumerator
	return whole*blaxelMegabytesPerMebibyteDenominator + remainder*blaxelMegabytesPerMebibyteDenominator/blaxelMegabytesPerMebibyteNumerator
}

// formatBlaxelDuration formats a whole-second duration as a Blaxel TTL string.
// minimum is the Blaxel floor for that policy (max-age TTL: 5m, idle TTL: 1m).
// Rejection text matches TypeScript durationString.
func formatBlaxelDuration(value, minimum time.Duration) (string, error) {
	if value <= 0 || value%time.Second != 0 {
		return "", fmt.Errorf("duration must be a positive whole number of seconds")
	}
	if value < minimum {
		return "", fmt.Errorf("duration must be at least %s", blaxelDurationUnits(minimum))
	}
	return blaxelDurationUnits(value), nil
}

func blaxelDurationUnits(value time.Duration) string {
	seconds := int64(value / time.Second)
	switch {
	case seconds%86400 == 0:
		return fmt.Sprintf("%dd", seconds/86400)
	case seconds%3600 == 0:
		return fmt.Sprintf("%dh", seconds/3600)
	case seconds%60 == 0:
		return fmt.Sprintf("%dm", seconds/60)
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
