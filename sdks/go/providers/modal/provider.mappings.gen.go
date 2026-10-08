// Code generated from specs/providers/modal.yaml. DO NOT EDIT.
package modal

import (
	fmt "fmt"
	_go "github.com/modal-labs/modal-client/go"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

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
			target.CPU = value
		}
	}
	remaining.CPUCores = nil
	if source.CPULimitCores != nil {
		{
			value := source.GetCPULimitCores()
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
