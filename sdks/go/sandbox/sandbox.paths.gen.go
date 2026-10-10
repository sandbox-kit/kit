// Code generated from shared schema paths. DO NOT EDIT.
package sandbox

const CreateFieldName = "name"
const CreateFieldSource = "source"
const CreateFieldSourceProviderDefault = "source.provider_default"
const CreateFieldSourceImage = "source.image"
const CreateFieldSourceImageReference = "source.image.reference"
const CreateFieldSourceSnapshot = "source.snapshot"
const CreateFieldSourceSnapshotReference = "source.snapshot.reference"
const CreateFieldSourceSnapshotKind = "source.snapshot.kind"
const CreateFieldSourceSnapshotProviderKind = "source.snapshot.provider_kind"
const CreateFieldSourceWarmPool = "source.warm_pool"
const CreateFieldSourceWarmPoolReference = "source.warm_pool.reference"
const CreateFieldRuntime = "runtime"
const CreateFieldRuntimeLanguage = "runtime.language"
const CreateFieldRuntimeVersion = "runtime.version"
const CreateFieldRuntimeUser = "runtime.user"
const CreateFieldRuntimeWorkingDirectory = "runtime.working_directory"
const CreateFieldRuntimeEntrypoint = "runtime.entrypoint"
const CreateFieldRuntimePTY = "runtime.pty"
const CreateFieldIsolation = "isolation"
const CreateFieldIsolationKind = "isolation.kind"
const CreateFieldIsolationProviderKind = "isolation.provider_kind"
const CreateFieldIsolationNestedVirtualization = "isolation.nested_virtualization"
const CreateFieldResources = "resources"
const CreateFieldResourcesCPUCores = "resources.cpu_cores"
const CreateFieldResourcesCPULimitCores = "resources.cpu_limit_cores"
const CreateFieldResourcesMemoryMiB = "resources.memory_mib"
const CreateFieldResourcesMemoryLimitMiB = "resources.memory_limit_mib"
const CreateFieldResourcesDiskMiB = "resources.disk_mib"
const CreateFieldPlacement = "placement"
const CreateFieldPlacementCloud = "placement.cloud"
const CreateFieldPlacementRegions = "placement.regions"
const CreateFieldPlacementTarget = "placement.target"
const CreateFieldEnvironment = "environment"
const CreateFieldLabels = "labels"
const CreateFieldNetwork = "network"
const CreateFieldNetworkEgress = "network.egress"
const CreateFieldNetworkOutboundCIDRs = "network.outbound_cidrs"
const CreateFieldNetworkOutboundCIDRsEntries = "network.outbound_cidrs.entries"
const CreateFieldNetworkOutboundDomains = "network.outbound_domains"
const CreateFieldNetworkOutboundDomainsEntries = "network.outbound_domains.entries"
const CreateFieldNetworkInboundCIDRs = "network.inbound_cidrs"
const CreateFieldNetworkInboundCIDRsEntries = "network.inbound_cidrs.entries"
const CreateFieldNetworkPrivateNetwork = "network.private_network"
const CreateFieldNetworkLinkedSandbox = "network.linked_sandbox"
const CreateFieldNetworkProxyReference = "network.proxy_reference"
const CreateFieldNetworkOutboundProxyURL = "network.outbound_proxy_url"
const CreateFieldNetworkCustomDomain = "network.custom_domain"
const CreateFieldNetworkPorts = "network.ports"
const CreateFieldNetworkPortsPort = "network.ports.port"
const CreateFieldNetworkPortsTransport = "network.ports.transport"
const CreateFieldSecurity = "security"
const CreateFieldSecurityPublicAccess = "security.public_access"
const CreateFieldSecurityWorkloadIdentity = "security.workload_identity"
const CreateFieldSecuritySecrets = "security.secrets"
const CreateFieldSecuritySecretsReference = "security.secrets.reference"
const CreateFieldSecuritySecretsEnvironmentVariable = "security.secrets.environment_variable"
const CreateFieldSecuritySecretsInjection = "security.secrets.injection"
const CreateFieldLifetime = "lifetime"
const CreateFieldLifetimeMaximumLifetime = "lifetime.maximum_lifetime"
const CreateFieldLifetimeEphemeral = "lifetime.ephemeral"
const CreateFieldLifetimeIdleTerminate = "lifetime.idle_terminate"
const CreateFieldLifetimeIdleTerminateMode = "lifetime.idle_terminate.mode"
const CreateFieldLifetimeIdleTerminateAfter = "lifetime.idle_terminate.after"
const CreateFieldLifetimeIdleStop = "lifetime.idle_stop"
const CreateFieldLifetimeIdleStopMode = "lifetime.idle_stop.mode"
const CreateFieldLifetimeIdleStopAfter = "lifetime.idle_stop.after"
const CreateFieldLifetimeIdlePause = "lifetime.idle_pause"
const CreateFieldLifetimeIdlePauseMode = "lifetime.idle_pause.mode"
const CreateFieldLifetimeIdlePauseAfter = "lifetime.idle_pause.after"
const CreateFieldLifetimeStoppedArchive = "lifetime.stopped_archive"
const CreateFieldLifetimeStoppedArchiveMode = "lifetime.stopped_archive.mode"
const CreateFieldLifetimeStoppedArchiveAfter = "lifetime.stopped_archive.after"
const CreateFieldLifetimeStoppedDelete = "lifetime.stopped_delete"
const CreateFieldLifetimeStoppedDeleteMode = "lifetime.stopped_delete.mode"
const CreateFieldLifetimeStoppedDeleteAfter = "lifetime.stopped_delete.after"
const CreateFieldProvisioning = "provisioning"
const CreateFieldProvisioningTimeout = "provisioning.timeout"
const CreateFieldProvisioningQueueTimeout = "provisioning.queue_timeout"
const CreateFieldProvisioningWaitFor = "provisioning.wait_for"
const CreateFieldReadiness = "readiness"
const CreateFieldReadinessTCPPort = "readiness.tcp_port"
const CreateFieldReadinessCommand = "readiness.command"
const CreateFieldReadinessCommandArgv = "readiness.command.argv"
const CreateFieldObservability = "observability"
const CreateFieldObservabilityVerbose = "observability.verbose"
const CreateFieldObservabilityTelemetryEndpoint = "observability.telemetry_endpoint"
const CreateFieldProviderOptions = "provider_options"
const InfoFieldID = "id"
const InfoFieldProvider = "provider"
const InfoFieldName = "name"
const InfoFieldProviderState = "provider_state"
const InfoFieldRegion = "region"
const InfoFieldIsolation = "isolation"
const InfoFieldIsolationKind = "isolation.kind"
const InfoFieldIsolationProviderKind = "isolation.provider_kind"
const InfoFieldIsolationNestedVirtualization = "isolation.nested_virtualization"
const InfoFieldResources = "resources"
const InfoFieldResourcesCPUCores = "resources.cpu_cores"
const InfoFieldResourcesCPULimitCores = "resources.cpu_limit_cores"
const InfoFieldResourcesMemoryMiB = "resources.memory_mib"
const InfoFieldResourcesMemoryLimitMiB = "resources.memory_limit_mib"
const InfoFieldResourcesDiskMiB = "resources.disk_mib"
const InfoFieldLabels = "labels"
const InfoFieldCreatedAt = "created_at"
const InfoFieldProviderMetadata = "provider_metadata"
const InfoFieldOrigins = "origins"

func responseFieldPresent(info *SandboxInfo, path string) bool {
	if info == nil {
		return false
	}
	switch path {
	case InfoFieldID:
		return info.ID != ""
	case InfoFieldProvider:
		return info.Provider != ""
	case InfoFieldName:
		return info.Name != nil
	case InfoFieldProviderState:
		return info.ProviderState != nil
	case InfoFieldRegion:
		return info.Region != nil
	case InfoFieldIsolation:
		return info.Isolation != nil
	case InfoFieldIsolationKind:
		if info.Isolation == nil {
			return false
		}
		return true
	case InfoFieldIsolationProviderKind:
		if info.Isolation == nil {
			return false
		}
		return info.Isolation.ProviderKind != nil
	case InfoFieldIsolationNestedVirtualization:
		if info.Isolation == nil {
			return false
		}
		return info.Isolation.NestedVirtualization != nil
	case InfoFieldResources:
		return info.Resources != nil
	case InfoFieldResourcesCPUCores:
		if info.Resources == nil {
			return false
		}
		return info.Resources.CPUCores != nil
	case InfoFieldResourcesCPULimitCores:
		if info.Resources == nil {
			return false
		}
		return info.Resources.CPULimitCores != nil
	case InfoFieldResourcesMemoryMiB:
		if info.Resources == nil {
			return false
		}
		return info.Resources.MemoryMiB != nil
	case InfoFieldResourcesMemoryLimitMiB:
		if info.Resources == nil {
			return false
		}
		return info.Resources.MemoryLimitMiB != nil
	case InfoFieldResourcesDiskMiB:
		if info.Resources == nil {
			return false
		}
		return info.Resources.DiskMiB != nil
	case InfoFieldLabels:
		return info.Labels != nil
	case InfoFieldCreatedAt:
		return info.CreatedAt != nil
	case InfoFieldProviderMetadata:
		return info.ProviderMetadata != nil
	default:
		return false
	}
}
func validateResponseOrigins(info *SandboxInfo) error {
	for path, origin := range info.Origins {
		if !origin.Valid() || !responseFieldPresent(info, path) {
			return NewError(ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: info.Provider, Operation: "create", Field: "origins." + path, Message: "sandbox-kit: origin must describe a known, present response field"}, nil)
		}
	}
	return nil
}
