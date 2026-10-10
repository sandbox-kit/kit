// Code generated from shared schema. DO NOT EDIT.
package sandbox

import (
	fmt "fmt"
)

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *CreateOptions) Clone() (*CreateOptions, error) {
	return x.clone(0)
}
func (x *CreateOptions) clone(depth int) (*CreateOptions, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Name != nil {
		value := *x.Name
		out.Name = &value
	}
	{
		value, err := x.Source.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Source = value
	}
	{
		value, err := x.Runtime.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Runtime = value
	}
	{
		value, err := x.Isolation.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Isolation = value
	}
	{
		value, err := x.Resources.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Resources = value
	}
	{
		value, err := x.Placement.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Placement = value
	}
	if x.Environment != nil {
		out.Environment = make(map[string]string, len(x.Environment))
		for key, value := range x.Environment {
			out.Environment[key] = value
		}
	}
	if x.Labels != nil {
		out.Labels = make(map[string]string, len(x.Labels))
		for key, value := range x.Labels {
			out.Labels[key] = value
		}
	}
	{
		value, err := x.Network.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Network = value
	}
	{
		value, err := x.Security.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Security = value
	}
	{
		value, err := x.Lifetime.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Lifetime = value
	}
	{
		value, err := x.Provisioning.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Provisioning = value
	}
	{
		value, err := x.Readiness.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Readiness = value
	}
	{
		value, err := x.Observability.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Observability = value
	}
	if x.ProviderOptions != nil {
		value, err := cloneMetadata(x.ProviderOptions)
		if err != nil {
			return nil, err
		}
		out.ProviderOptions = value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *SandboxSource) Clone() (*SandboxSource, error) {
	return x.clone(0)
}
func (x *SandboxSource) clone(depth int) (*SandboxSource, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	{
		value, err := x.ProviderDefault.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.ProviderDefault = value
	}
	{
		value, err := x.Image.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Image = value
	}
	{
		value, err := x.Snapshot.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Snapshot = value
	}
	{
		value, err := x.WarmPool.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.WarmPool = value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *ProviderDefaultSource) Clone() (*ProviderDefaultSource, error) {
	return x.clone(0)
}
func (x *ProviderDefaultSource) clone(depth int) (*ProviderDefaultSource, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *ImageSource) Clone() (*ImageSource, error) {
	return x.clone(0)
}
func (x *ImageSource) clone(depth int) (*ImageSource, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *SnapshotSource) Clone() (*SnapshotSource, error) {
	return x.clone(0)
}
func (x *SnapshotSource) clone(depth int) (*SnapshotSource, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.ProviderKind != nil {
		value := *x.ProviderKind
		out.ProviderKind = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *WarmPoolSource) Clone() (*WarmPoolSource, error) {
	return x.clone(0)
}
func (x *WarmPoolSource) clone(depth int) (*WarmPoolSource, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *RuntimeConfig) Clone() (*RuntimeConfig, error) {
	return x.clone(0)
}
func (x *RuntimeConfig) clone(depth int) (*RuntimeConfig, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Language != nil {
		value := *x.Language
		out.Language = &value
	}
	if x.Version != nil {
		value := *x.Version
		out.Version = &value
	}
	if x.User != nil {
		value := *x.User
		out.User = &value
	}
	if x.WorkingDirectory != nil {
		value := *x.WorkingDirectory
		out.WorkingDirectory = &value
	}
	if x.Entrypoint != nil {
		out.Entrypoint = make([]string, len(x.Entrypoint))
		for i, value := range x.Entrypoint {
			out.Entrypoint[i] = value
		}
	}
	if x.PTY != nil {
		value := *x.PTY
		out.PTY = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *IsolationConfig) Clone() (*IsolationConfig, error) {
	return x.clone(0)
}
func (x *IsolationConfig) clone(depth int) (*IsolationConfig, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.ProviderKind != nil {
		value := *x.ProviderKind
		out.ProviderKind = &value
	}
	if x.NestedVirtualization != nil {
		value := *x.NestedVirtualization
		out.NestedVirtualization = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *Resources) Clone() (*Resources, error) {
	return x.clone(0)
}
func (x *Resources) clone(depth int) (*Resources, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.CPUCores != nil {
		value := *x.CPUCores
		out.CPUCores = &value
	}
	if x.CPULimitCores != nil {
		value := *x.CPULimitCores
		out.CPULimitCores = &value
	}
	if x.MemoryMiB != nil {
		value := *x.MemoryMiB
		out.MemoryMiB = &value
	}
	if x.MemoryLimitMiB != nil {
		value := *x.MemoryLimitMiB
		out.MemoryLimitMiB = &value
	}
	if x.DiskMiB != nil {
		value := *x.DiskMiB
		out.DiskMiB = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *Placement) Clone() (*Placement, error) {
	return x.clone(0)
}
func (x *Placement) clone(depth int) (*Placement, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Cloud != nil {
		value := *x.Cloud
		out.Cloud = &value
	}
	if x.Regions != nil {
		out.Regions = make([]string, len(x.Regions))
		for i, value := range x.Regions {
			out.Regions[i] = value
		}
	}
	if x.Target != nil {
		value := *x.Target
		out.Target = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *StringAllowlist) Clone() (*StringAllowlist, error) {
	return x.clone(0)
}
func (x *StringAllowlist) clone(depth int) (*StringAllowlist, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Entries != nil {
		out.Entries = make([]string, len(x.Entries))
		for i, value := range x.Entries {
			out.Entries[i] = value
		}
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *NetworkConfig) Clone() (*NetworkConfig, error) {
	return x.clone(0)
}
func (x *NetworkConfig) clone(depth int) (*NetworkConfig, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	{
		value, err := x.OutboundCIDRs.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.OutboundCIDRs = value
	}
	{
		value, err := x.OutboundDomains.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.OutboundDomains = value
	}
	{
		value, err := x.InboundCIDRs.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.InboundCIDRs = value
	}
	if x.PrivateNetwork != nil {
		value := *x.PrivateNetwork
		out.PrivateNetwork = &value
	}
	if x.LinkedSandbox != nil {
		value := *x.LinkedSandbox
		out.LinkedSandbox = &value
	}
	if x.ProxyReference != nil {
		value := *x.ProxyReference
		out.ProxyReference = &value
	}
	if x.OutboundProxyURL != nil {
		value := *x.OutboundProxyURL
		out.OutboundProxyURL = &value
	}
	if x.CustomDomain != nil {
		value := *x.CustomDomain
		out.CustomDomain = &value
	}
	if x.Ports != nil {
		out.Ports = make([]*PortBinding, len(x.Ports))
		for i, value := range x.Ports {
			copied, err := value.clone(depth + 1)
			if err != nil {
				return nil, err
			}
			out.Ports[i] = copied
		}
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *PortBinding) Clone() (*PortBinding, error) {
	return x.clone(0)
}
func (x *PortBinding) clone(depth int) (*PortBinding, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *SecurityConfig) Clone() (*SecurityConfig, error) {
	return x.clone(0)
}
func (x *SecurityConfig) clone(depth int) (*SecurityConfig, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.PublicAccess != nil {
		value := *x.PublicAccess
		out.PublicAccess = &value
	}
	if x.WorkloadIdentity != nil {
		value := *x.WorkloadIdentity
		out.WorkloadIdentity = &value
	}
	if x.Secrets != nil {
		out.Secrets = make([]*SecretReference, len(x.Secrets))
		for i, value := range x.Secrets {
			copied, err := value.clone(depth + 1)
			if err != nil {
				return nil, err
			}
			out.Secrets[i] = copied
		}
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *SecretReference) Clone() (*SecretReference, error) {
	return x.clone(0)
}
func (x *SecretReference) clone(depth int) (*SecretReference, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.EnvironmentVariable != nil {
		value := *x.EnvironmentVariable
		out.EnvironmentVariable = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *LifetimePolicy) Clone() (*LifetimePolicy, error) {
	return x.clone(0)
}
func (x *LifetimePolicy) clone(depth int) (*LifetimePolicy, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.MaximumLifetime != nil {
		value := *x.MaximumLifetime
		out.MaximumLifetime = &value
	}
	if x.Ephemeral != nil {
		value := *x.Ephemeral
		out.Ephemeral = &value
	}
	{
		value, err := x.IdleTerminate.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.IdleTerminate = value
	}
	{
		value, err := x.IdleStop.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.IdleStop = value
	}
	{
		value, err := x.IdlePause.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.IdlePause = value
	}
	{
		value, err := x.StoppedArchive.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.StoppedArchive = value
	}
	{
		value, err := x.StoppedDelete.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.StoppedDelete = value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *AutomaticAction) Clone() (*AutomaticAction, error) {
	return x.clone(0)
}
func (x *AutomaticAction) clone(depth int) (*AutomaticAction, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.After != nil {
		value := *x.After
		out.After = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *ProvisioningOptions) Clone() (*ProvisioningOptions, error) {
	return x.clone(0)
}
func (x *ProvisioningOptions) clone(depth int) (*ProvisioningOptions, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Timeout != nil {
		value := *x.Timeout
		out.Timeout = &value
	}
	if x.QueueTimeout != nil {
		value := *x.QueueTimeout
		out.QueueTimeout = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *ReadinessProbe) Clone() (*ReadinessProbe, error) {
	return x.clone(0)
}
func (x *ReadinessProbe) clone(depth int) (*ReadinessProbe, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.TCPPort != nil {
		value := *x.TCPPort
		out.TCPPort = &value
	}
	{
		value, err := x.Command.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Command = value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *CommandProbe) Clone() (*CommandProbe, error) {
	return x.clone(0)
}
func (x *CommandProbe) clone(depth int) (*CommandProbe, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Argv != nil {
		out.Argv = make([]string, len(x.Argv))
		for i, value := range x.Argv {
			out.Argv[i] = value
		}
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *ObservabilityConfig) Clone() (*ObservabilityConfig, error) {
	return x.clone(0)
}
func (x *ObservabilityConfig) clone(depth int) (*ObservabilityConfig, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Verbose != nil {
		value := *x.Verbose
		out.Verbose = &value
	}
	if x.TelemetryEndpoint != nil {
		value := *x.TelemetryEndpoint
		out.TelemetryEndpoint = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *CreateResult) Clone() (*CreateResult, error) {
	return x.clone(0)
}
func (x *CreateResult) clone(depth int) (*CreateResult, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	{
		value, err := x.Sandbox.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Sandbox = value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *SandboxInfo) Clone() (*SandboxInfo, error) {
	return x.clone(0)
}
func (x *SandboxInfo) clone(depth int) (*SandboxInfo, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Name != nil {
		value := *x.Name
		out.Name = &value
	}
	if x.ProviderState != nil {
		value := *x.ProviderState
		out.ProviderState = &value
	}
	if x.Region != nil {
		value := *x.Region
		out.Region = &value
	}
	{
		value, err := x.Isolation.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Isolation = value
	}
	{
		value, err := x.Resources.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.Resources = value
	}
	if x.Labels != nil {
		out.Labels = make(map[string]string, len(x.Labels))
		for key, value := range x.Labels {
			out.Labels[key] = value
		}
	}
	if x.CreatedAt != nil {
		value := *x.CreatedAt
		out.CreatedAt = &value
	}
	if x.ProviderMetadata != nil {
		value, err := cloneMetadata(x.ProviderMetadata)
		if err != nil {
			return nil, err
		}
		out.ProviderMetadata = value
	}
	if x.Origins != nil {
		out.Origins = make(map[string]ValueOrigin, len(x.Origins))
		for key, value := range x.Origins {
			out.Origins[key] = value
		}
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *ErrorInfo) Clone() (*ErrorInfo, error) {
	return x.clone(0)
}
func (x *ErrorInfo) clone(depth int) (*ErrorInfo, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.StatusCode != nil {
		value := *x.StatusCode
		out.StatusCode = &value
	}
	if x.ProviderCode != nil {
		value := *x.ProviderCode
		out.ProviderCode = &value
	}
	if x.ProviderSource != nil {
		value := *x.ProviderSource
		out.ProviderSource = &value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *MetadataValue) Clone() (*MetadataValue, error) {
	return x.clone(0)
}
func (x *MetadataValue) clone(depth int) (*MetadataValue, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Boolean != nil {
		value := *x.Boolean
		out.Boolean = &value
	}
	if x.Text != nil {
		value := *x.Text
		out.Text = &value
	}
	if x.SignedInteger != nil {
		value := *x.SignedInteger
		out.SignedInteger = &value
	}
	if x.UnsignedInteger != nil {
		value := *x.UnsignedInteger
		out.UnsignedInteger = &value
	}
	if x.Number != nil {
		value := *x.Number
		out.Number = &value
	}
	if x.Binary != nil {
		var value []byte
		if *x.Binary != nil {
			value = append([]byte{}, (*x.Binary)...)
		}
		out.Binary = &value
	}
	{
		value, err := x.List.clone(depth + 1)
		if err != nil {
			return nil, err
		}
		out.List = value
	}
	if x.Object != nil {
		value, err := cloneMetadata(x.Object)
		if err != nil {
			return nil, err
		}
		out.Object = value
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *MetadataList) Clone() (*MetadataList, error) {
	return x.clone(0)
}
func (x *MetadataList) clone(depth int) (*MetadataList, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Values != nil {
		out.Values = make([]*MetadataValue, len(x.Values))
		for i, value := range x.Values {
			copied, err := value.clone(depth + 1)
			if err != nil {
				return nil, err
			}
			out.Values[i] = copied
		}
	}
	return &out, nil
}

// Clone returns an owned copy, preserving field presence and numeric types.
// It rejects unsupported metadata and excessive depth or cycles.
func (x *MetadataObject) Clone() (*MetadataObject, error) {
	return x.clone(0)
}
func (x *MetadataObject) clone(depth int) (*MetadataObject, error) {
	if x == nil {
		return nil, nil
	}
	if depth > 64 {
		return nil, fmt.Errorf("sandbox-kit: data exceeds maximum copy depth or contains a cycle")
	}
	out := *x
	if x.Values != nil {
		out.Values = make(map[string]*MetadataValue, len(x.Values))
		for key, value := range x.Values {
			copied, err := value.clone(depth + 1)
			if err != nil {
				return nil, err
			}
			out.Values[key] = copied
		}
	}
	return &out, nil
}
