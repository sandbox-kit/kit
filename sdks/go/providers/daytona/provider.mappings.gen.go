// Code generated from specs/providers/daytona.yaml. DO NOT EDIT.
package daytona

import (
	fmt "fmt"
	daytona "github.com/daytona/clients/sdk-go/pkg/daytona"
	types "github.com/daytona/clients/sdk-go/pkg/types"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
	math "math"
)

func mapCommonCreateFields(source *sandbox.CreateOptions) (types.SandboxBaseParams, sandbox.CreateOptions, error) {
	target := types.SandboxBaseParams{}
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
		target.EnvVars = value
	}
	remaining.Environment = nil
	{
		value := source.Labels
		target.Labels = value
	}
	remaining.Labels = nil
	return target, remaining, nil
}
func mapResources(source *sandbox.Resources) (types.Resources, sandbox.Resources, error) {
	target := types.Resources{}
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
			if float64(value) != math.Trunc(float64(value)) || value > 2147483647 {
				return target, remaining, fmt.Errorf("sandbox-kit daytona: cpu_cores cannot be represented")
			}
			target.CPU = int(value)
		}
	}
	remaining.CPUCores = nil
	if source.MemoryMiB != nil {
		{
			value := source.GetMemoryMiB()
			if value%1024 != 0 || value/1024 > 2147483647 {
				return target, remaining, fmt.Errorf("sandbox-kit daytona: memory_mib cannot be represented")
			}
			target.Memory = int(value / 1024)
		}
	}
	remaining.MemoryMiB = nil
	if source.DiskMiB != nil {
		{
			value := source.GetDiskMiB()
			if value%1024 != 0 || value/1024 > 2147483647 {
				return target, remaining, fmt.Errorf("sandbox-kit daytona: disk_mib cannot be represented")
			}
			target.Disk = int(value / 1024)
		}
	}
	remaining.DiskMiB = nil
	return target, remaining, nil
}
func mapSandboxInfo(source *daytona.Sandbox) (sandbox.SandboxInfo, error) {
	target := sandbox.SandboxInfo{}
	if source == nil {
		return target, nil
	}
	{
		value := source.ID
		target.ID = value
	}
	{
		value := source.Name
		target.Name = sandbox.Value(value)
	}
	{
		value := source.State
		target.ProviderState = sandbox.Value(string(value))
	}
	{
		value := source.Target
		target.Region = sandbox.Value(value)
	}
	{
		value := source.Labels
		target.Labels = value
	}
	target.Provider = "daytona"
	return target, nil
}
func mapAllocatedResources(source *daytona.Sandbox) (sandbox.Resources, error) {
	target := sandbox.Resources{}
	if source == nil {
		return target, nil
	}
	if !(source.Cpu > 0 && source.Memory > 0 && source.Disk > 0) {
		return target, nil
	}
	{
		value := source.Cpu
		target.CPUCores = sandbox.Value(float64(value))
	}
	{
		value := source.Memory
		if value < 0 || uint64(value) > ^uint64(0)/1024 {
			return target, fmt.Errorf("sandbox-kit daytona: memory cannot be represented")
		}
		target.MemoryMiB = sandbox.Value(uint64(value) * 1024)
	}
	{
		value := source.Disk
		if value < 0 || uint64(value) > ^uint64(0)/1024 {
			return target, fmt.Errorf("sandbox-kit daytona: disk cannot be represented")
		}
		target.DiskMiB = sandbox.Value(uint64(value) * 1024)
	}
	return target, nil
}
func mapClientSettings(source *sandbox.Config) (types.DaytonaConfig, sandbox.Config, error) {
	target := types.DaytonaConfig{}
	remaining := sandbox.Config{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.Endpoint != nil {
		{
			value := source.GetEndpoint()
			target.APIUrl = value
		}
	}
	remaining.Endpoint = nil
	if source.Region != nil {
		{
			value := source.GetRegion()
			target.Target = value
		}
	}
	remaining.Region = nil
	return target, remaining, nil
}
func mapAPIKeyCredentials(source *sandbox.APIKeyCredentials) (types.DaytonaConfig, sandbox.APIKeyCredentials, error) {
	target := types.DaytonaConfig{}
	remaining := sandbox.APIKeyCredentials{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	{
		value := source.Key
		target.APIKey = value
	}
	remaining.Key = ""
	return target, remaining, nil
}
func mapBearerAuth(source *sandbox.BearerTokenCredentials) (types.DaytonaConfig, sandbox.BearerTokenCredentials, error) {
	target := types.DaytonaConfig{}
	remaining := sandbox.BearerTokenCredentials{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	{
		value := source.Token
		target.JWTToken = value
	}
	remaining.Token = ""
	return target, remaining, nil
}
func mapScope(source *sandbox.Scope) (types.DaytonaConfig, sandbox.Scope, error) {
	target := types.DaytonaConfig{}
	remaining := sandbox.Scope{}
	if source != nil {
		remaining = *source
	}
	if source == nil {
		return target, remaining, nil
	}
	if source.OrganizationID != nil {
		{
			value := source.GetOrganizationID()
			target.OrganizationID = value
		}
	}
	remaining.OrganizationID = nil
	return target, remaining, nil
}
