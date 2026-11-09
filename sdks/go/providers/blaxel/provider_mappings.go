package blaxel

import (
	"fmt"
	"math"
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/blaxel-ai/sdk-go/shared"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
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
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindInvalidArgument,
				Provider:  "blaxel",
				Operation: "create",
				Field:     sandbox.CreateFieldResourcesMemoryMib,
				Message:   "sandbox-kit blaxel: memory_mib cannot be represented",
			}, nil)
		}
		target.Memory = blaxel.Int(int64(value))
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
	target.MemoryMib = sandbox.Value(uint64(source.Spec.Runtime.Memory))
	return target, nil
}

// formatBlaxelDuration formats a whole-second duration as a Blaxel TTL string.
// minimum is the Blaxel floor for that policy (max-age TTL: 5m, idle TTL: 1m).
func formatBlaxelDuration(value, minimum time.Duration) (string, error) {
	if value <= 0 || value%time.Second != 0 {
		return "", fmt.Errorf("duration must be a positive whole number of seconds")
	}
	if value < minimum {
		return "", fmt.Errorf("duration must be at least %s", minimum)
	}
	seconds := int64(value / time.Second)
	switch {
	case seconds%86400 == 0:
		return fmt.Sprintf("%dd", seconds/86400), nil
	case seconds%3600 == 0:
		return fmt.Sprintf("%dh", seconds/3600), nil
	case seconds%60 == 0:
		return fmt.Sprintf("%dm", seconds/60), nil
	default:
		return fmt.Sprintf("%ds", seconds), nil
	}
}
