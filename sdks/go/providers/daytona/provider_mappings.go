// Provider field mappings between shared CreateOptions/Config and the native SDK.
// Inventoried in specs/providers/daytona.yaml groups; kept in sync by hand (not emitted).
package daytona

import (
	math "math"
	time "time"

	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	types "github.com/daytona/clients/sdk-go/pkg/types"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// mapLifetimePolicies maps supplied fields for mapLifetimePolicies into the
// native request and returns the fields that were not mapped.
func mapLifetimePolicies(source *sandbox.LifetimePolicy) (types.SandboxBaseParams, *sandbox.LifetimePolicy, error) {
	target := types.SandboxBaseParams{}
	if source == nil {
		return target, &sandbox.LifetimePolicy{}, nil
	}
	remaining := cloneMessage(source)
	if err := (minutePolicySpec{
		field: sandbox.CreateFieldLifetimeIdleStop, label: "idle_stop",
		allowDisable: true, forbidImmediate: true,
		set:        func(v int) { target.AutoStopInterval = &v },
		setDisable: func() { target.AutoStopInterval = sandbox.Value(0) },
	}).apply(source.IdleStop); err != nil {
		return target, remaining, err
	}
	remaining.IdleStop = nil
	if err := (minutePolicySpec{
		field: sandbox.CreateFieldLifetimeIdlePause, label: "idle_pause",
		allowDisable: true, forbidImmediate: true,
		set:        func(v int) { target.AutoPauseInterval = &v },
		setDisable: func() { target.AutoPauseInterval = sandbox.Value(0) },
	}).apply(source.IdlePause); err != nil {
		return target, remaining, err
	}
	remaining.IdlePause = nil
	if err := (minutePolicySpec{
		field: sandbox.CreateFieldLifetimeStoppedArchive, label: "stopped_archive",
		allowDisable: false, forbidImmediate: true, maxMinutes: 43200,
		set: func(v int) { target.AutoArchiveInterval = &v },
	}).apply(source.StoppedArchive); err != nil {
		return target, remaining, err
	}
	remaining.StoppedArchive = nil
	if err := (minutePolicySpec{
		field: sandbox.CreateFieldLifetimeStoppedDelete, label: "stopped_delete",
		allowDisable: false, forbidImmediate: false,
		set: func(v int) { target.AutoDeleteInterval = &v },
	}).apply(source.StoppedDelete); err != nil {
		return target, remaining, err
	}
	remaining.StoppedDelete = nil
	return target, remaining, nil
}

type minutePolicySpec struct {
	field           string
	label           string
	allowDisable    bool
	forbidImmediate bool
	maxMinutes      int
	set             func(int)
	setDisable      func()
}

func (spec minutePolicySpec) apply(policy *sandbox.AutomaticAction) error {
	if policy == nil {
		return nil
	}
	switch policy.Mode {
	case sandbox.PolicyModeDefault:
		return nil
	case sandbox.PolicyModeDisabled:
		if !spec.allowDisable {
			return sandbox.NewError(&sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create", Field: spec.field,
				Message: "sandbox-kit daytona: " + spec.label + ": explicit disabling is not representable; native zero has a different meaning",
			}, nil)
		}
		spec.setDisable()
		return nil
	case sandbox.PolicyModeAfter:
		if policy.After == nil {
			return sandbox.NewError(&sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: spec.field,
				Message: "sandbox-kit daytona: " + spec.label + ": duration is required",
			}, nil)
		}
		duration := policy.After.AsDuration()
		if spec.forbidImmediate && duration == 0 {
			return sandbox.NewError(&sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindUnsupported, Provider: "daytona", Operation: "create", Field: spec.field,
				Message: "sandbox-kit daytona: " + spec.label + ": immediate action is not representable",
			}, nil)
		}
		if duration < 0 || duration%time.Minute != 0 || duration/time.Minute > math.MaxInt32 {
			return sandbox.NewError(&sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: spec.field,
				Message: "sandbox-kit daytona: " + spec.label + ": delay cannot be represented in whole minutes",
			}, nil)
		}
		minutes := int(duration / time.Minute)
		if spec.maxMinutes > 0 && minutes > spec.maxMinutes {
			return sandbox.NewError(&sandbox.ErrorInfo{
				Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: spec.field,
				Message: "sandbox-kit daytona: " + spec.label + ": delay exceeds documented maximum",
			}, nil)
		}
		spec.set(minutes)
		return nil
	default:
		return sandbox.NewError(&sandbox.ErrorInfo{
			Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: spec.field,
			Message: "sandbox-kit daytona: " + spec.label + ": unknown policy mode",
		}, nil)
	}
}

// mapCommonCreateFields maps supplied fields for mapCommonCreateFields into
// the native request and returns the fields that were not mapped.
func mapCommonCreateFields(source *sandbox.CreateOptions) (types.SandboxBaseParams, *sandbox.CreateOptions, error) {
	target := types.SandboxBaseParams{}
	if source == nil {
		return target, &sandbox.CreateOptions{}, nil
	}
	remaining := cloneMessage(source)
	if source.Name != nil {
		value := source.GetName()
		target.Name = value
	}
	remaining.Name = nil
	value := source.Environment
	target.EnvVars = value
	remaining.Environment = nil
	value = source.Labels
	target.Labels = value
	remaining.Labels = nil
	return target, remaining, nil
}

// mapResources maps supplied fields for mapResources into the native request
// and returns the fields that were not mapped.
func mapResources(source *sandbox.Resources) (types.Resources, *sandbox.Resources, error) {
	target := types.Resources{}
	if source == nil {
		return target, &sandbox.Resources{}, nil
	}
	remaining := cloneMessage(source)
	if source.CpuCores != nil {
		value := source.GetCpuCores()
		if float64(value) != math.Trunc(float64(value)) || value > 2147483647 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: sandbox.CreateFieldResourcesCpuCores, Message: "sandbox-kit daytona: cpu_cores cannot be represented"}, nil)
		}
		target.CPU = int(value)
	}
	remaining.CpuCores = nil
	if source.MemoryMib != nil {
		value := source.GetMemoryMib()
		if value%1024 != 0 || value/1024 > 2147483647 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: sandbox.CreateFieldResourcesMemoryMib, Message: "sandbox-kit daytona: memory_mib cannot be represented"}, nil)
		}
		target.Memory = int(value / 1024)
	}
	remaining.MemoryMib = nil
	if source.DiskMib != nil {
		value := source.GetDiskMib()
		if value%1024 != 0 || value/1024 > 2147483647 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "daytona", Operation: "create", Field: sandbox.CreateFieldResourcesDiskMib, Message: "sandbox-kit daytona: disk_mib cannot be represented"}, nil)
		}
		target.Disk = int(value / 1024)
	}
	remaining.DiskMib = nil
	return target, remaining, nil
}

// mapSandboxInfo maps native fields for mapSandboxInfo into shared metadata.
// It does not invent values the SDK omitted.
func mapSandboxInfo(source *daytona.Sandbox) (*sandbox.SandboxInfo, error) {
	target := &sandbox.SandboxInfo{}
	if source == nil {
		return target, nil
	}
	target.Origins = map[string]sandbox.ValueOrigin{}
	target.Id = source.ID
	target.Origins["id"] = sandbox.ValueOriginProvider
	target.Name = sandbox.Value(source.Name)
	if target.Name != nil {
		target.Origins["name"] = sandbox.ValueOriginProvider
	}
	target.ProviderState = sandbox.Value(string(source.State))
	if target.ProviderState != nil {
		target.Origins["provider_state"] = sandbox.ValueOriginProvider
	}
	target.Region = sandbox.Value(source.Target)
	if target.Region != nil {
		target.Origins["region"] = sandbox.ValueOriginProvider
	}
	target.Labels = source.Labels
	if target.Labels != nil {
		target.Origins["labels"] = sandbox.ValueOriginProvider
	}
	target.Provider = "daytona"
	return target, nil
}

// mapAllocatedResources maps native fields for mapAllocatedResources into
// shared metadata. It does not invent values the SDK omitted.
func mapAllocatedResources(source *daytona.Sandbox) (*sandbox.Resources, error) {
	target := &sandbox.Resources{}
	if source == nil {
		return target, nil
	}
	if !(source.Cpu > 0 && source.Memory > 0 && source.Disk > 0) {
		return target, nil
	}
	value := source.Cpu
	target.CpuCores = sandbox.Value(float64(value))
	value = source.Memory
	if value < 0 || uint64(value) > ^uint64(0)/1024 {
		return target, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidResponse, Provider: "daytona", Operation: "create", Field: "memory", Message: "sandbox-kit daytona: memory cannot be represented"}, nil)
	}
	target.MemoryMib = sandbox.Value(uint64(value) * 1024)
	value = source.Disk
	if value < 0 || uint64(value) > ^uint64(0)/1024 {
		return target, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidResponse, Provider: "daytona", Operation: "create", Field: "disk", Message: "sandbox-kit daytona: disk cannot be represented"}, nil)
	}
	target.DiskMib = sandbox.Value(uint64(value) * 1024)
	return target, nil
}

// mapClientSettings maps supplied fields for mapClientSettings into the native
// request and returns the fields that were not mapped.
func mapClientSettings(source *sandbox.Config) (types.DaytonaConfig, *sandbox.Config, error) {
	target := types.DaytonaConfig{}
	if source == nil {
		return target, &sandbox.Config{}, nil
	}
	remaining := cloneMessage(source)
	if source.Endpoint != nil {
		value := source.GetEndpoint()
		target.APIUrl = value
	}
	remaining.Endpoint = nil
	if source.Region != nil {
		value := source.GetRegion()
		target.Target = value
	}
	remaining.Region = nil
	return target, remaining, nil
}

// mapAPIKeyCredentials maps supplied fields for mapAPIKeyCredentials into the
// native request and returns the fields that were not mapped.
func mapAPIKeyCredentials(source *sandbox.APIKeyCredentials) (types.DaytonaConfig, *sandbox.APIKeyCredentials, error) {
	target := types.DaytonaConfig{}
	if source == nil {
		return target, &sandbox.APIKeyCredentials{}, nil
	}
	remaining := cloneMessage(source)
	value := source.Key
	target.APIKey = value
	remaining.Key = ""
	return target, remaining, nil
}

// mapBearerTokenCredentials maps supplied fields into the native request and
// returns the fields that were not mapped.
func mapBearerTokenCredentials(source *sandbox.BearerTokenCredentials) (types.DaytonaConfig, *sandbox.BearerTokenCredentials, error) {
	target := types.DaytonaConfig{}
	if source == nil {
		return target, &sandbox.BearerTokenCredentials{}, nil
	}
	remaining := cloneMessage(source)
	value := source.Token
	target.JWTToken = value
	remaining.Token = ""
	return target, remaining, nil
}

// mapScope maps supplied fields for mapScope into the native request and
// returns the fields that were not mapped.
func mapScope(source *sandbox.Scope) (types.DaytonaConfig, *sandbox.Scope, error) {
	target := types.DaytonaConfig{}
	if source == nil {
		return target, &sandbox.Scope{}, nil
	}
	remaining := cloneMessage(source)
	if source.OrganizationId != nil {
		value := source.GetOrganizationId()
		target.OrganizationID = value
	}
	remaining.OrganizationId = nil
	return target, remaining, nil
}
