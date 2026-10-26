// Field-path constants identify create/info settings in errors. Values are canonical protobuf paths.
//
// Example:
//
//	fmt.Println(CreateFieldResourcesCPUCores) // resources.cpu_cores
package sandbox

import (
	"reflect"
	"strings"

	"google.golang.org/protobuf/reflect/protoreflect"
)

// CreateFieldName is the canonical path "name".
const CreateFieldName = "name"

// CreateFieldSource is the canonical path "source".
const CreateFieldSource = "source"

// CreateFieldSourceProviderDefault is the canonical path
// "source.provider_default".
const CreateFieldSourceProviderDefault = "source.provider_default"

// CreateFieldSourceImage is the canonical path "source.image".
const CreateFieldSourceImage = "source.image"

// CreateFieldSourceImageReference is the canonical path
// "source.image.reference".
const CreateFieldSourceImageReference = "source.image.reference"

// CreateFieldSourceSnapshot is the canonical path "source.snapshot".
const CreateFieldSourceSnapshot = "source.snapshot"

// CreateFieldSourceSnapshotReference is the canonical path
// "source.snapshot.reference".
const CreateFieldSourceSnapshotReference = "source.snapshot.reference"

// CreateFieldSourceSnapshotKind is the canonical path "source.snapshot.kind".
const CreateFieldSourceSnapshotKind = "source.snapshot.kind"

// CreateFieldSourceSnapshotProviderKind is the canonical path
// "source.snapshot.provider_kind".
const CreateFieldSourceSnapshotProviderKind = "source.snapshot.provider_kind"

// CreateFieldSourceWarmPool is the canonical path "source.warm_pool".
const CreateFieldSourceWarmPool = "source.warm_pool"

// CreateFieldSourceWarmPoolReference is the canonical path
// "source.warm_pool.reference".
const CreateFieldSourceWarmPoolReference = "source.warm_pool.reference"

// CreateFieldRuntime is the canonical path "runtime".
const CreateFieldRuntime = "runtime"

// CreateFieldRuntimeLanguage is the canonical path "runtime.language".
const CreateFieldRuntimeLanguage = "runtime.language"

// CreateFieldRuntimeVersion is the canonical path "runtime.version".
const CreateFieldRuntimeVersion = "runtime.version"

// CreateFieldRuntimeUser is the canonical path "runtime.user".
const CreateFieldRuntimeUser = "runtime.user"

// CreateFieldRuntimeWorkingDirectory is the canonical path
// "runtime.working_directory".
const CreateFieldRuntimeWorkingDirectory = "runtime.working_directory"

// CreateFieldRuntimeEntrypoint is the canonical path "runtime.entrypoint".
const CreateFieldRuntimeEntrypoint = "runtime.entrypoint"

// CreateFieldRuntimePty is the canonical path "runtime.pty".
const CreateFieldRuntimePty = "runtime.pty"

// CreateFieldIsolation is the canonical path "isolation".
const CreateFieldIsolation = "isolation"

// CreateFieldIsolationKind is the canonical path "isolation.kind".
const CreateFieldIsolationKind = "isolation.kind"

// CreateFieldIsolationProviderKind is the canonical path
// "isolation.provider_kind".
const CreateFieldIsolationProviderKind = "isolation.provider_kind"

// CreateFieldIsolationNestedVirtualization is the canonical path
// "isolation.nested_virtualization".
const CreateFieldIsolationNestedVirtualization = "isolation.nested_virtualization"

// CreateFieldResources is the canonical path "resources".
const CreateFieldResources = "resources"

// CreateFieldResourcesCpuCores is the canonical path "resources.cpu_cores".
const CreateFieldResourcesCpuCores = "resources.cpu_cores"

// CreateFieldResourcesCpuLimitCores is the canonical path
// "resources.cpu_limit_cores".
const CreateFieldResourcesCpuLimitCores = "resources.cpu_limit_cores"

// CreateFieldResourcesMemoryMib is the canonical path "resources.memory_mib".
const CreateFieldResourcesMemoryMib = "resources.memory_mib"

// CreateFieldResourcesMemoryLimitMib is the canonical path
// "resources.memory_limit_mib".
const CreateFieldResourcesMemoryLimitMib = "resources.memory_limit_mib"

// CreateFieldResourcesDiskMib is the canonical path "resources.disk_mib".
const CreateFieldResourcesDiskMib = "resources.disk_mib"

// CreateFieldPlacement is the canonical path "placement".
const CreateFieldPlacement = "placement"

// CreateFieldPlacementCloud is the canonical path "placement.cloud".
const CreateFieldPlacementCloud = "placement.cloud"

// CreateFieldPlacementRegions is the canonical path "placement.regions".
const CreateFieldPlacementRegions = "placement.regions"

// CreateFieldPlacementTarget is the canonical path "placement.target".
const CreateFieldPlacementTarget = "placement.target"

// CreateFieldEnvironment is the canonical path "environment".
const CreateFieldEnvironment = "environment"

// CreateFieldLabels is the canonical path "labels".
const CreateFieldLabels = "labels"

// CreateFieldNetwork is the canonical path "network".
const CreateFieldNetwork = "network"

// CreateFieldNetworkEgress is the canonical path "network.egress".
const CreateFieldNetworkEgress = "network.egress"

// CreateFieldNetworkOutboundCidrs is the canonical path
// "network.outbound_cidrs".
const CreateFieldNetworkOutboundCidrs = "network.outbound_cidrs"

// CreateFieldNetworkOutboundCidrsEntries is the canonical path
// "network.outbound_cidrs.entries".
const CreateFieldNetworkOutboundCidrsEntries = "network.outbound_cidrs.entries"

// CreateFieldNetworkOutboundDomains is the canonical path
// "network.outbound_domains".
const CreateFieldNetworkOutboundDomains = "network.outbound_domains"

// CreateFieldNetworkOutboundDomainsEntries is the canonical path
// "network.outbound_domains.entries".
const CreateFieldNetworkOutboundDomainsEntries = "network.outbound_domains.entries"

// CreateFieldNetworkInboundCidrs is the canonical path
// "network.inbound_cidrs".
const CreateFieldNetworkInboundCidrs = "network.inbound_cidrs"

// CreateFieldNetworkInboundCidrsEntries is the canonical path
// "network.inbound_cidrs.entries".
const CreateFieldNetworkInboundCidrsEntries = "network.inbound_cidrs.entries"

// CreateFieldNetworkPrivateNetwork is the canonical path
// "network.private_network".
const CreateFieldNetworkPrivateNetwork = "network.private_network"

// CreateFieldNetworkLinkedSandbox is the canonical path
// "network.linked_sandbox".
const CreateFieldNetworkLinkedSandbox = "network.linked_sandbox"

// CreateFieldNetworkProxyReference is the canonical path
// "network.proxy_reference".
const CreateFieldNetworkProxyReference = "network.proxy_reference"

// CreateFieldNetworkOutboundProxyUrl is the canonical path
// "network.outbound_proxy_url".
const CreateFieldNetworkOutboundProxyUrl = "network.outbound_proxy_url"

// CreateFieldNetworkCustomDomain is the canonical path
// "network.custom_domain".
const CreateFieldNetworkCustomDomain = "network.custom_domain"

// CreateFieldNetworkPorts is the canonical path "network.ports".
const CreateFieldNetworkPorts = "network.ports"

// CreateFieldNetworkPortsPort is the canonical path "network.ports.port".
const CreateFieldNetworkPortsPort = "network.ports.port"

// CreateFieldNetworkPortsTransport is the canonical path
// "network.ports.transport".
const CreateFieldNetworkPortsTransport = "network.ports.transport"

// CreateFieldSecurity is the canonical path "security".
const CreateFieldSecurity = "security"

// CreateFieldSecurityPublicAccess is the canonical path
// "security.public_access".
const CreateFieldSecurityPublicAccess = "security.public_access"

// CreateFieldSecurityWorkloadIdentity is the canonical path
// "security.workload_identity".
const CreateFieldSecurityWorkloadIdentity = "security.workload_identity"

// CreateFieldSecuritySecrets is the canonical path "security.secrets".
const CreateFieldSecuritySecrets = "security.secrets"

// CreateFieldSecuritySecretsReference is the canonical path
// "security.secrets.reference".
const CreateFieldSecuritySecretsReference = "security.secrets.reference"

// CreateFieldSecuritySecretsEnvironmentVariable is the canonical path
// "security.secrets.environment_variable".
const CreateFieldSecuritySecretsEnvironmentVariable = "security.secrets.environment_variable"

// CreateFieldSecuritySecretsInjection is the canonical path
// "security.secrets.injection".
const CreateFieldSecuritySecretsInjection = "security.secrets.injection"

// CreateFieldLifetime is the canonical path "lifetime".
const CreateFieldLifetime = "lifetime"

// CreateFieldLifetimeMaximumLifetime is the canonical path
// "lifetime.maximum_lifetime".
const CreateFieldLifetimeMaximumLifetime = "lifetime.maximum_lifetime"

// CreateFieldLifetimeEphemeral is the canonical path "lifetime.ephemeral".
const CreateFieldLifetimeEphemeral = "lifetime.ephemeral"

// CreateFieldLifetimeIdleTerminate is the canonical path
// "lifetime.idle_terminate".
const CreateFieldLifetimeIdleTerminate = "lifetime.idle_terminate"

// CreateFieldLifetimeIdleTerminateMode is the canonical path
// "lifetime.idle_terminate.mode".
const CreateFieldLifetimeIdleTerminateMode = "lifetime.idle_terminate.mode"

// CreateFieldLifetimeIdleTerminateAfter is the canonical path
// "lifetime.idle_terminate.after".
const CreateFieldLifetimeIdleTerminateAfter = "lifetime.idle_terminate.after"

// CreateFieldLifetimeIdleStop is the canonical path "lifetime.idle_stop".
const CreateFieldLifetimeIdleStop = "lifetime.idle_stop"

// CreateFieldLifetimeIdleStopMode is the canonical path
// "lifetime.idle_stop.mode".
const CreateFieldLifetimeIdleStopMode = "lifetime.idle_stop.mode"

// CreateFieldLifetimeIdleStopAfter is the canonical path
// "lifetime.idle_stop.after".
const CreateFieldLifetimeIdleStopAfter = "lifetime.idle_stop.after"

// CreateFieldLifetimeIdlePause is the canonical path "lifetime.idle_pause".
const CreateFieldLifetimeIdlePause = "lifetime.idle_pause"

// CreateFieldLifetimeIdlePauseMode is the canonical path
// "lifetime.idle_pause.mode".
const CreateFieldLifetimeIdlePauseMode = "lifetime.idle_pause.mode"

// CreateFieldLifetimeIdlePauseAfter is the canonical path
// "lifetime.idle_pause.after".
const CreateFieldLifetimeIdlePauseAfter = "lifetime.idle_pause.after"

// CreateFieldLifetimeStoppedArchive is the canonical path
// "lifetime.stopped_archive".
const CreateFieldLifetimeStoppedArchive = "lifetime.stopped_archive"

// CreateFieldLifetimeStoppedArchiveMode is the canonical path
// "lifetime.stopped_archive.mode".
const CreateFieldLifetimeStoppedArchiveMode = "lifetime.stopped_archive.mode"

// CreateFieldLifetimeStoppedArchiveAfter is the canonical path
// "lifetime.stopped_archive.after".
const CreateFieldLifetimeStoppedArchiveAfter = "lifetime.stopped_archive.after"

// CreateFieldLifetimeStoppedDelete is the canonical path
// "lifetime.stopped_delete".
const CreateFieldLifetimeStoppedDelete = "lifetime.stopped_delete"

// CreateFieldLifetimeStoppedDeleteMode is the canonical path
// "lifetime.stopped_delete.mode".
const CreateFieldLifetimeStoppedDeleteMode = "lifetime.stopped_delete.mode"

// CreateFieldLifetimeStoppedDeleteAfter is the canonical path
// "lifetime.stopped_delete.after".
const CreateFieldLifetimeStoppedDeleteAfter = "lifetime.stopped_delete.after"

// CreateFieldProvisioning is the canonical path "provisioning".
const CreateFieldProvisioning = "provisioning"

// CreateFieldProvisioningTimeout is the canonical path "provisioning.timeout".
const CreateFieldProvisioningTimeout = "provisioning.timeout"

// CreateFieldProvisioningQueueTimeout is the canonical path
// "provisioning.queue_timeout".
const CreateFieldProvisioningQueueTimeout = "provisioning.queue_timeout"

// CreateFieldProvisioningWaitFor is the canonical path
// "provisioning.wait_for".
const CreateFieldProvisioningWaitFor = "provisioning.wait_for"

// CreateFieldReadiness is the canonical path "readiness".
const CreateFieldReadiness = "readiness"

// CreateFieldReadinessTcpPort is the canonical path "readiness.tcp_port".
const CreateFieldReadinessTcpPort = "readiness.tcp_port"

// CreateFieldReadinessCommand is the canonical path "readiness.command".
const CreateFieldReadinessCommand = "readiness.command"

// CreateFieldReadinessCommandArgv is the canonical path
// "readiness.command.argv".
const CreateFieldReadinessCommandArgv = "readiness.command.argv"

// CreateFieldObservability is the canonical path "observability".
const CreateFieldObservability = "observability"

// CreateFieldObservabilityVerbose is the canonical path
// "observability.verbose".
const CreateFieldObservabilityVerbose = "observability.verbose"

// CreateFieldObservabilityTelemetryEndpoint is the canonical path
// "observability.telemetry_endpoint".
const CreateFieldObservabilityTelemetryEndpoint = "observability.telemetry_endpoint"

// CreateFieldProviderOptions is the canonical path "provider_options".
const CreateFieldProviderOptions = "provider_options"

// InfoFieldId is the canonical path "id".
const InfoFieldId = "id"

// InfoFieldProvider is the canonical path "provider".
const InfoFieldProvider = "provider"

// InfoFieldName is the canonical path "name".
const InfoFieldName = "name"

// InfoFieldProviderState is the canonical path "provider_state".
const InfoFieldProviderState = "provider_state"

// InfoFieldRegion is the canonical path "region".
const InfoFieldRegion = "region"

// InfoFieldIsolation is the canonical path "isolation".
const InfoFieldIsolation = "isolation"

// InfoFieldIsolationKind is the canonical path "isolation.kind".
const InfoFieldIsolationKind = "isolation.kind"

// InfoFieldIsolationProviderKind is the canonical path
// "isolation.provider_kind".
const InfoFieldIsolationProviderKind = "isolation.provider_kind"

// InfoFieldIsolationNestedVirtualization is the canonical path
// "isolation.nested_virtualization".
const InfoFieldIsolationNestedVirtualization = "isolation.nested_virtualization"

// InfoFieldResources is the canonical path "resources".
const InfoFieldResources = "resources"

// InfoFieldResourcesCpuCores is the canonical path "resources.cpu_cores".
const InfoFieldResourcesCpuCores = "resources.cpu_cores"

// InfoFieldResourcesCpuLimitCores is the canonical path
// "resources.cpu_limit_cores".
const InfoFieldResourcesCpuLimitCores = "resources.cpu_limit_cores"

// InfoFieldResourcesMemoryMib is the canonical path "resources.memory_mib".
const InfoFieldResourcesMemoryMib = "resources.memory_mib"

// InfoFieldResourcesMemoryLimitMib is the canonical path
// "resources.memory_limit_mib".
const InfoFieldResourcesMemoryLimitMib = "resources.memory_limit_mib"

// InfoFieldResourcesDiskMib is the canonical path "resources.disk_mib".
const InfoFieldResourcesDiskMib = "resources.disk_mib"

// InfoFieldLabels is the canonical path "labels".
const InfoFieldLabels = "labels"

// InfoFieldCreatedAt is the canonical path "created_at".
const InfoFieldCreatedAt = "created_at"

// InfoFieldProviderMetadata is the canonical path "provider_metadata".
const InfoFieldProviderMetadata = "provider_metadata"

// InfoFieldOrigins is the canonical path "origins".
const InfoFieldOrigins = "origins"

// InfoFieldState is the canonical path "state".
const InfoFieldState = "state"

// InfoFieldTerminal is the canonical path "terminal".
const InfoFieldTerminal = "terminal"

// responseFieldPresent reports whether path names a field that is present on
// the sandbox metadata. Presence follows protobuf descriptors rather than a
// handwritten field inventory.
func responseFieldPresent(info *SandboxInfo, path string) bool {
	if info == nil || path == "" {
		return false
	}
	msg := info.ProtoReflect()
	parts := strings.Split(path, ".")
	for i, part := range parts {
		fd := msg.Descriptor().Fields().ByName(protoreflect.Name(part))
		if fd == nil {
			return false
		}
		last := i == len(parts)-1
		switch {
		case fd.IsMap():
			if !last {
				return false
			}
			// ProtoReflect.Has is false for empty maps; Kit treats a non-nil empty
			// map as present (same as clone/origins semantics).
			return msg.Has(fd) || mapFieldNonNil(msg, fd)
		case fd.Message() != nil:
			if !msg.Has(fd) {
				return false
			}
			if last {
				return true
			}
			msg = msg.Get(fd).Message()
		default:
			if !last {
				return false
			}
			if fd.HasPresence() {
				return msg.Has(fd)
			}
			if fd.Kind() == protoreflect.StringKind {
				return msg.Get(fd).String() != ""
			}
			// Proto3 enums and other non-presence scalars are present when the
			// containing message was reached (for example state, isolation.kind).
			return true
		}
	}
	return false
}

// mapFieldNonNil reports whether the Go map backing fd is non-nil (including empty).
func mapFieldNonNil(msg protoreflect.Message, fd protoreflect.FieldDescriptor) bool {
	value := reflect.ValueOf(msg.Interface())
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return false
		}
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return false
	}
	name := string(fd.Name())
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if field.PkgPath != "" {
			continue
		}
		if protobufFieldName(field) != name {
			continue
		}
		mapValue := value.Field(i)
		return mapValue.Kind() == reflect.Map && !mapValue.IsNil()
	}
	return false
}

// validateResponseOrigins checks that every origin names a known field that is
// actually present. A missing origin means the source is unknown.
func validateResponseOrigins(info *SandboxInfo, operation string) error {
	for path, origin := range info.Origins {
		if !validValueOrigin(origin) || !responseFieldPresent(info, path) {
			return NewError(&ErrorInfo{Kind: ErrorKindInvalidResponse, Provider: info.Provider, Operation: operation, Field: "origins." + path, Message: "sandbox-kit: origin must describe a known, present response field"}, nil)
		}
	}
	return nil
}

// protobufFieldName returns the canonical protobuf field name from a generated
// struct tag. Falls back to the Go field name when the tag is missing.
func protobufFieldName(field reflect.StructField) string {
	for _, part := range strings.Split(field.Tag.Get("protobuf"), ",") {
		if name, ok := strings.CutPrefix(part, "name="); ok && name != "" {
			return name
		}
	}
	return field.Name
}
