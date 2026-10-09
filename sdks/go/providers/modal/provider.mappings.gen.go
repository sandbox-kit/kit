// Code generated from specs/providers/modal.yaml. DO NOT EDIT.
package modal

import (
	fmt "fmt"
	_go "github.com/modal-labs/modal-client/go"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
	math "math"
	time "time"
)

func mapClientScope(source *sandbox.Scope) (_go.ClientParams, sandbox.Scope, error) {
	target := _go.ClientParams{}
	remaining := sandbox.Scope{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.Environment != nil {
		{
			value := source.GetEnvironment()
			target.Environment = value
		}
	}
	remaining.Environment = nil
	return target, remaining, nil
}
func mapCommonCreateFields(source *sandbox.CreateOptions) (_go.SandboxCreateParams, sandbox.CreateOptions, error) {
	target := _go.SandboxCreateParams{}
	remaining := sandbox.CreateOptions{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.Name != nil {
		{
			value := source.GetName()
			target.Name = value
		}
	}
	remaining.Name = nil
	{
		value := source.Environment
		target.Env = value
	}
	remaining.Environment = nil
	{
		value := source.Labels
		target.Tags = value
	}
	remaining.Labels = nil
	return target, remaining, nil
}
func mapResources(source *sandbox.Resources) (_go.SandboxCreateParams, sandbox.Resources, error) {
	target := _go.SandboxCreateParams{}
	remaining := sandbox.Resources{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.CPUCores != nil {
		{
			value := source.GetCPUCores()
			scaled := math.Round(value * 1000)
			if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value != scaled/1000 || scaled > 4294967295 {
				return target, remaining, fmt.Errorf("sandbox-kit modal: cpu_cores cannot be represented")
			}
			if value*1000 < scaled {
				value = math.Nextafter(value, math.Inf(1))
			}
			target.CPU = value
		}
	}
	remaining.CPUCores = nil
	if source.CPULimitCores != nil {
		{
			value := source.GetCPULimitCores()
			scaled := math.Round(value * 1000)
			if math.IsNaN(value) || math.IsInf(value, 0) || value < 0 || value != scaled/1000 || scaled > 4294967295 {
				return target, remaining, fmt.Errorf("sandbox-kit modal: cpu_limit_cores cannot be represented")
			}
			if value*1000 < scaled {
				value = math.Nextafter(value, math.Inf(1))
			}
			target.CPULimit = value
		}
	}
	remaining.CPULimitCores = nil
	if source.MemoryMiB != nil {
		{
			value := source.GetMemoryMiB()
			if value > 2147483647 {
				return target, remaining, fmt.Errorf("sandbox-kit modal: memory_mib cannot be represented")
			}
			target.MemoryMiB = int(value)
		}
	}
	remaining.MemoryMiB = nil
	if source.MemoryLimitMiB != nil {
		{
			value := source.GetMemoryLimitMiB()
			if value > 2147483647 {
				return target, remaining, fmt.Errorf("sandbox-kit modal: memory_limit_mib cannot be represented")
			}
			target.MemoryLimitMiB = int(value)
		}
	}
	remaining.MemoryLimitMiB = nil
	return target, remaining, nil
}
func mapLifetimeDuration(source *sandbox.LifetimePolicy) (_go.SandboxCreateParams, sandbox.LifetimePolicy, error) {
	target := _go.SandboxCreateParams{}
	remaining := sandbox.LifetimePolicy{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.MaximumLifetime != nil {
		{
			value := *source.MaximumLifetime
			if value <= 0 || value%time.Second != 0 || value/time.Second > 86400 {
				return target, remaining, fmt.Errorf("sandbox-kit modal: maximum_lifetime cannot be represented")
			}
			target.Timeout = value
		}
	}
	remaining.MaximumLifetime = nil
	return target, remaining, nil
}
func mapIdleDuration(source *sandbox.AutomaticAction) (_go.SandboxCreateParams, sandbox.AutomaticAction, error) {
	target := _go.SandboxCreateParams{}
	remaining := sandbox.AutomaticAction{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.After != nil {
		{
			value := *source.After
			if value <= 0 || value%time.Second != 0 || value/time.Second > 4294967295 {
				return target, remaining, fmt.Errorf("sandbox-kit modal: after cannot be represented")
			}
			target.IdleTimeout = value
		}
	}
	remaining.After = nil
	return target, remaining, nil
}
func mapTokenPairCredentials(source *sandbox.TokenPairCredentials) (_go.ClientParams, sandbox.TokenPairCredentials, error) {
	target := _go.ClientParams{}
	remaining := sandbox.TokenPairCredentials{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	{
		value := source.ID
		target.TokenID = value
	}
	remaining.ID = ""
	{
		value := source.Secret
		target.TokenSecret = value
	}
	remaining.Secret = ""
	return target, remaining, nil
}
func mapOAuthAuth(source *sandbox.OAuthCredentials) (_go.OAuthCredentialsParams, sandbox.OAuthCredentials, error) {
	target := _go.OAuthCredentialsParams{}
	remaining := sandbox.OAuthCredentials{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	{
		value := source.RefreshToken
		target.RefreshToken = value
	}
	remaining.RefreshToken = ""
	{
		value := source.ClientID
		target.ClientID = value
	}
	remaining.ClientID = ""
	if source.ClientSecret != nil {
		{
			value := source.GetClientSecret()
			target.ClientSecret = value
		}
	}
	remaining.ClientSecret = nil
	if source.JWTKey != nil {
		{
			value := source.GetJWTKey()
			target.JWTKey = value
		}
	}
	remaining.JWTKey = nil
	return target, remaining, nil
}
