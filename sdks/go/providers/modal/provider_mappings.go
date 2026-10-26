package modal

import (
	math "math"
	time "time"

	sdk "github.com/modal-labs/modal-client/go"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// mapSandboxIdentity maps native fields for mapSandboxIdentity into shared
// metadata. It does not invent values the SDK omitted.
func mapSandboxIdentity(source *sdk.Sandbox) (*sandbox.SandboxInfo, error) {
	target := &sandbox.SandboxInfo{}
	if source == nil {
		return target, nil
	}
	target.Origins = map[string]sandbox.ValueOrigin{}
	value := source.SandboxID
	target.Id = value
	target.Origins["id"] = sandbox.ValueOriginProvider
	target.Provider = "modal"
	return target, nil
}

// mapScope maps supplied fields into the native request and returns the
// fields that were not mapped.
func mapScope(source *sandbox.Scope) (sdk.ClientParams, *sandbox.Scope, error) {
	target := sdk.ClientParams{}
	if source == nil {
		return target, &sandbox.Scope{}, nil
	}
	remaining := cloneMessage(source)
	if source.Environment != nil {
		value := source.GetEnvironment()
		target.Environment = value
	}
	remaining.Environment = nil
	return target, remaining, nil
}

// mapCommonCreateFields maps supplied fields for mapCommonCreateFields into
// the native request and returns the fields that were not mapped.
func mapCommonCreateFields(source *sandbox.CreateOptions) (sdk.SandboxCreateParams, *sandbox.CreateOptions, error) {
	target := sdk.SandboxCreateParams{}
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
	target.Env = value
	remaining.Environment = nil
	value = source.Labels
	target.Tags = value
	remaining.Labels = nil
	return target, remaining, nil
}

// mapResources maps supplied fields for mapResources into the native request
// and returns the fields that were not mapped.
func mapResources(source *sandbox.Resources) (sdk.SandboxCreateParams, *sandbox.Resources, error) {
	target := sdk.SandboxCreateParams{}
	if source == nil {
		return target, &sandbox.Resources{}, nil
	}
	remaining := cloneMessage(source)
	if source.CpuCores != nil {
		value := source.GetCpuCores()
		scaled := math.Round(value * 1000)
		if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value != scaled/1000 || scaled > 4294967295 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldResourcesCpuCores, Message: "sandbox-kit modal: cpu_cores cannot be represented"}, nil)
		}
		if value*1000 < scaled {
			value = math.Nextafter(value, math.Inf(1))
		}
		target.CPU = value
	}
	remaining.CpuCores = nil
	if source.CpuLimitCores != nil {
		value := source.GetCpuLimitCores()
		scaled := math.Round(value * 1000)
		if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value != scaled/1000 || scaled > 4294967295 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldResourcesCpuLimitCores, Message: "sandbox-kit modal: cpu_limit_cores cannot be represented"}, nil)
		}
		if value*1000 < scaled {
			value = math.Nextafter(value, math.Inf(1))
		}
		target.CPULimit = value
	}
	remaining.CpuLimitCores = nil
	if source.MemoryMib != nil {
		value := source.GetMemoryMib()
		if value > 2147483647 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldResourcesMemoryMib, Message: "sandbox-kit modal: memory_mib cannot be represented"}, nil)
		}
		target.MemoryMiB = int(value)
	}
	remaining.MemoryMib = nil
	if source.MemoryLimitMib != nil {
		value := source.GetMemoryLimitMib()
		if value > 2147483647 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldResourcesMemoryLimitMib, Message: "sandbox-kit modal: memory_limit_mib cannot be represented"}, nil)
		}
		target.MemoryLimitMiB = int(value)
	}
	remaining.MemoryLimitMib = nil
	return target, remaining, nil
}

// mapLifetimeDuration maps supplied fields for mapLifetimeDuration into the
// native request and returns the fields that were not mapped.
func mapLifetimeDuration(source *sandbox.LifetimePolicy) (sdk.SandboxCreateParams, *sandbox.LifetimePolicy, error) {
	target := sdk.SandboxCreateParams{}
	if source == nil {
		return target, &sandbox.LifetimePolicy{}, nil
	}
	remaining := cloneMessage(source)
	if source.MaximumLifetime != nil {
		value := source.MaximumLifetime.AsDuration()
		if value <= 0 || value%time.Second != 0 || value/time.Second > 86400 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldLifetimeMaximumLifetime, Message: "sandbox-kit modal: maximum_lifetime cannot be represented"}, nil)
		}
		target.Timeout = value
	}
	remaining.MaximumLifetime = nil
	return target, remaining, nil
}

// mapIdleDuration maps supplied fields for mapIdleDuration into the native
// request and returns the fields that were not mapped.
func mapIdleDuration(source *sandbox.AutomaticAction) (sdk.SandboxCreateParams, *sandbox.AutomaticAction, error) {
	target := sdk.SandboxCreateParams{}
	if source == nil {
		return target, &sandbox.AutomaticAction{}, nil
	}
	remaining := cloneMessage(source)
	if source.After != nil {
		value := source.After.AsDuration()
		if value <= 0 || value%time.Second != 0 || value/time.Second > 4294967295 {
			return target, remaining, sandbox.NewError(&sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidArgument, Provider: "modal", Operation: "create", Field: sandbox.CreateFieldLifetimeIdleTerminateAfter, Message: "sandbox-kit modal: after cannot be represented"}, nil)
		}
		target.IdleTimeout = value
	}
	remaining.After = nil
	return target, remaining, nil
}

// mapTokenPairCredentials maps supplied fields for mapTokenPairCredentials
// into the native request and returns the fields that were not mapped.
func mapTokenPairCredentials(source *sandbox.TokenPairCredentials) (sdk.ClientParams, *sandbox.TokenPairCredentials, error) {
	target := sdk.ClientParams{}
	if source == nil {
		return target, &sandbox.TokenPairCredentials{}, nil
	}
	remaining := cloneMessage(source)
	value := source.Id
	target.TokenID = value
	remaining.Id = ""
	value = source.Secret
	target.TokenSecret = value
	remaining.Secret = ""
	return target, remaining, nil
}

// mapOAuthCredentials maps supplied fields into the native request and
// returns the fields that were not mapped.
func mapOAuthCredentials(source *sandbox.OAuthCredentials) (sdk.OAuthCredentialsParams, *sandbox.OAuthCredentials, error) {
	target := sdk.OAuthCredentialsParams{}
	if source == nil {
		return target, &sandbox.OAuthCredentials{}, nil
	}
	remaining := cloneMessage(source)
	value := source.RefreshToken
	target.RefreshToken = value
	remaining.RefreshToken = ""
	value = source.ClientId
	target.ClientID = value
	remaining.ClientId = ""
	if source.ClientSecret != nil {
		value = source.GetClientSecret()
		target.ClientSecret = value
	}
	remaining.ClientSecret = nil
	if source.JwtKey != nil {
		value = source.GetJwtKey()
		target.JWTKey = value
	}
	remaining.JwtKey = nil
	return target, remaining, nil
}
