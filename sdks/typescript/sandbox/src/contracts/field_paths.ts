// Field-path constants identify create/info settings in errors.
// Values are canonical protobuf paths (snake_case), matching Go.
//
// Example:
//
//   console.log(CreateFieldResourcesCpuCores); // resources.cpu_cores

/** Canonical path "name". */
export const CreateFieldName = "name" as const;

/** Canonical path "source". */
export const CreateFieldSource = "source" as const;

/** Canonical path "source.provider_default". */
export const CreateFieldSourceProviderDefault = "source.provider_default" as const;

/** Canonical path "source.image". */
export const CreateFieldSourceImage = "source.image" as const;

/** Canonical path "source.image.reference". */
export const CreateFieldSourceImageReference = "source.image.reference" as const;

/** Canonical path "source.snapshot". */
export const CreateFieldSourceSnapshot = "source.snapshot" as const;

/** Canonical path "source.snapshot.reference". */
export const CreateFieldSourceSnapshotReference = "source.snapshot.reference" as const;

/** Canonical path "source.snapshot.kind". */
export const CreateFieldSourceSnapshotKind = "source.snapshot.kind" as const;

/** Canonical path "source.snapshot.provider_kind". */
export const CreateFieldSourceSnapshotProviderKind = "source.snapshot.provider_kind" as const;

/** Canonical path "source.warm_pool". */
export const CreateFieldSourceWarmPool = "source.warm_pool" as const;

/** Canonical path "source.warm_pool.reference". */
export const CreateFieldSourceWarmPoolReference = "source.warm_pool.reference" as const;

/** Canonical path "runtime". */
export const CreateFieldRuntime = "runtime" as const;

/** Canonical path "runtime.language". */
export const CreateFieldRuntimeLanguage = "runtime.language" as const;

/** Canonical path "runtime.version". */
export const CreateFieldRuntimeVersion = "runtime.version" as const;

/** Canonical path "runtime.user". */
export const CreateFieldRuntimeUser = "runtime.user" as const;

/** Canonical path "runtime.working_directory". */
export const CreateFieldRuntimeWorkingDirectory = "runtime.working_directory" as const;

/** Canonical path "runtime.entrypoint". */
export const CreateFieldRuntimeEntrypoint = "runtime.entrypoint" as const;

/** Canonical path "runtime.pty". */
export const CreateFieldRuntimePty = "runtime.pty" as const;

/** Canonical path "isolation". */
export const CreateFieldIsolation = "isolation" as const;

/** Canonical path "isolation.kind". */
export const CreateFieldIsolationKind = "isolation.kind" as const;

/** Canonical path "isolation.provider_kind". */
export const CreateFieldIsolationProviderKind = "isolation.provider_kind" as const;

/** Canonical path "isolation.nested_virtualization". */
export const CreateFieldIsolationNestedVirtualization = "isolation.nested_virtualization" as const;

/** Canonical path "resources". */
export const CreateFieldResources = "resources" as const;

/** Canonical path "resources.cpu_cores". */
export const CreateFieldResourcesCpuCores = "resources.cpu_cores" as const;

/** Canonical path "resources.cpu_limit_cores". */
export const CreateFieldResourcesCpuLimitCores = "resources.cpu_limit_cores" as const;

/** Canonical path "resources.memory_mib". */
export const CreateFieldResourcesMemoryMib = "resources.memory_mib" as const;

/** Canonical path "resources.memory_limit_mib". */
export const CreateFieldResourcesMemoryLimitMib = "resources.memory_limit_mib" as const;

/** Canonical path "resources.disk_mib". */
export const CreateFieldResourcesDiskMib = "resources.disk_mib" as const;

/** Canonical path "placement". */
export const CreateFieldPlacement = "placement" as const;

/** Canonical path "placement.cloud". */
export const CreateFieldPlacementCloud = "placement.cloud" as const;

/** Canonical path "placement.regions". */
export const CreateFieldPlacementRegions = "placement.regions" as const;

/** Canonical path "placement.target". */
export const CreateFieldPlacementTarget = "placement.target" as const;

/** Canonical path "environment". */
export const CreateFieldEnvironment = "environment" as const;

/** Canonical path "labels". */
export const CreateFieldLabels = "labels" as const;

/** Canonical path "network". */
export const CreateFieldNetwork = "network" as const;

/** Canonical path "network.egress". */
export const CreateFieldNetworkEgress = "network.egress" as const;

/** Canonical path "network.outbound_cidrs". */
export const CreateFieldNetworkOutboundCidrs = "network.outbound_cidrs" as const;

/** Canonical path "network.outbound_cidrs.entries". */
export const CreateFieldNetworkOutboundCidrsEntries = "network.outbound_cidrs.entries" as const;

/** Canonical path "network.outbound_domains". */
export const CreateFieldNetworkOutboundDomains = "network.outbound_domains" as const;

/** Canonical path "network.outbound_domains.entries". */
export const CreateFieldNetworkOutboundDomainsEntries = "network.outbound_domains.entries" as const;

/** Canonical path "network.inbound_cidrs". */
export const CreateFieldNetworkInboundCidrs = "network.inbound_cidrs" as const;

/** Canonical path "network.inbound_cidrs.entries". */
export const CreateFieldNetworkInboundCidrsEntries = "network.inbound_cidrs.entries" as const;

/** Canonical path "network.private_network". */
export const CreateFieldNetworkPrivateNetwork = "network.private_network" as const;

/** Canonical path "network.linked_sandbox". */
export const CreateFieldNetworkLinkedSandbox = "network.linked_sandbox" as const;

/** Canonical path "network.proxy_reference". */
export const CreateFieldNetworkProxyReference = "network.proxy_reference" as const;

/** Canonical path "network.outbound_proxy_url". */
export const CreateFieldNetworkOutboundProxyUrl = "network.outbound_proxy_url" as const;

/** Canonical path "network.custom_domain". */
export const CreateFieldNetworkCustomDomain = "network.custom_domain" as const;

/** Canonical path "network.ports". */
export const CreateFieldNetworkPorts = "network.ports" as const;

/** Canonical path "network.ports.port". */
export const CreateFieldNetworkPortsPort = "network.ports.port" as const;

/** Canonical path "network.ports.transport". */
export const CreateFieldNetworkPortsTransport = "network.ports.transport" as const;

/** Canonical path "security". */
export const CreateFieldSecurity = "security" as const;

/** Canonical path "security.public_access". */
export const CreateFieldSecurityPublicAccess = "security.public_access" as const;

/** Canonical path "security.workload_identity". */
export const CreateFieldSecurityWorkloadIdentity = "security.workload_identity" as const;

/** Canonical path "security.secrets". */
export const CreateFieldSecuritySecrets = "security.secrets" as const;

/** Canonical path "security.secrets.reference". */
export const CreateFieldSecuritySecretsReference = "security.secrets.reference" as const;

/** Canonical path "security.secrets.environment_variable". */
export const CreateFieldSecuritySecretsEnvironmentVariable = "security.secrets.environment_variable" as const;

/** Canonical path "security.secrets.injection". */
export const CreateFieldSecuritySecretsInjection = "security.secrets.injection" as const;

/** Canonical path "lifetime". */
export const CreateFieldLifetime = "lifetime" as const;

/** Canonical path "lifetime.maximum_lifetime". */
export const CreateFieldLifetimeMaximumLifetime = "lifetime.maximum_lifetime" as const;

/** Canonical path "lifetime.ephemeral". */
export const CreateFieldLifetimeEphemeral = "lifetime.ephemeral" as const;

/** Canonical path "lifetime.idle_terminate". */
export const CreateFieldLifetimeIdleTerminate = "lifetime.idle_terminate" as const;

/** Canonical path "lifetime.idle_terminate.mode". */
export const CreateFieldLifetimeIdleTerminateMode = "lifetime.idle_terminate.mode" as const;

/** Canonical path "lifetime.idle_terminate.after". */
export const CreateFieldLifetimeIdleTerminateAfter = "lifetime.idle_terminate.after" as const;

/** Canonical path "lifetime.idle_stop". */
export const CreateFieldLifetimeIdleStop = "lifetime.idle_stop" as const;

/** Canonical path "lifetime.idle_stop.mode". */
export const CreateFieldLifetimeIdleStopMode = "lifetime.idle_stop.mode" as const;

/** Canonical path "lifetime.idle_stop.after". */
export const CreateFieldLifetimeIdleStopAfter = "lifetime.idle_stop.after" as const;

/** Canonical path "lifetime.idle_pause". */
export const CreateFieldLifetimeIdlePause = "lifetime.idle_pause" as const;

/** Canonical path "lifetime.idle_pause.mode". */
export const CreateFieldLifetimeIdlePauseMode = "lifetime.idle_pause.mode" as const;

/** Canonical path "lifetime.idle_pause.after". */
export const CreateFieldLifetimeIdlePauseAfter = "lifetime.idle_pause.after" as const;

/** Canonical path "lifetime.stopped_archive". */
export const CreateFieldLifetimeStoppedArchive = "lifetime.stopped_archive" as const;

/** Canonical path "lifetime.stopped_archive.mode". */
export const CreateFieldLifetimeStoppedArchiveMode = "lifetime.stopped_archive.mode" as const;

/** Canonical path "lifetime.stopped_archive.after". */
export const CreateFieldLifetimeStoppedArchiveAfter = "lifetime.stopped_archive.after" as const;

/** Canonical path "lifetime.stopped_delete". */
export const CreateFieldLifetimeStoppedDelete = "lifetime.stopped_delete" as const;

/** Canonical path "lifetime.stopped_delete.mode". */
export const CreateFieldLifetimeStoppedDeleteMode = "lifetime.stopped_delete.mode" as const;

/** Canonical path "lifetime.stopped_delete.after". */
export const CreateFieldLifetimeStoppedDeleteAfter = "lifetime.stopped_delete.after" as const;

/** Canonical path "provisioning". */
export const CreateFieldProvisioning = "provisioning" as const;

/** Canonical path "provisioning.timeout". */
export const CreateFieldProvisioningTimeout = "provisioning.timeout" as const;

/** Canonical path "provisioning.queue_timeout". */
export const CreateFieldProvisioningQueueTimeout = "provisioning.queue_timeout" as const;

/** Canonical path "provisioning.wait_for". */
export const CreateFieldProvisioningWaitFor = "provisioning.wait_for" as const;

/** Canonical path "readiness". */
export const CreateFieldReadiness = "readiness" as const;

/** Canonical path "readiness.tcp_port". */
export const CreateFieldReadinessTcpPort = "readiness.tcp_port" as const;

/** Canonical path "readiness.command". */
export const CreateFieldReadinessCommand = "readiness.command" as const;

/** Canonical path "readiness.command.argv". */
export const CreateFieldReadinessCommandArgv = "readiness.command.argv" as const;

/** Canonical path "observability". */
export const CreateFieldObservability = "observability" as const;

/** Canonical path "observability.verbose". */
export const CreateFieldObservabilityVerbose = "observability.verbose" as const;

/** Canonical path "observability.telemetry_endpoint". */
export const CreateFieldObservabilityTelemetryEndpoint = "observability.telemetry_endpoint" as const;

/** Canonical path "id". */
export const InfoFieldId = "id" as const;

/** Canonical path "provider". */
export const InfoFieldProvider = "provider" as const;

/** Canonical path "name". */
export const InfoFieldName = "name" as const;

/** Canonical path "provider_state". */
export const InfoFieldProviderState = "provider_state" as const;

/** Canonical path "region". */
export const InfoFieldRegion = "region" as const;

/** Canonical path "isolation". */
export const InfoFieldIsolation = "isolation" as const;

/** Canonical path "isolation.kind". */
export const InfoFieldIsolationKind = "isolation.kind" as const;

/** Canonical path "isolation.provider_kind". */
export const InfoFieldIsolationProviderKind = "isolation.provider_kind" as const;

/** Canonical path "isolation.nested_virtualization". */
export const InfoFieldIsolationNestedVirtualization = "isolation.nested_virtualization" as const;

/** Canonical path "resources". */
export const InfoFieldResources = "resources" as const;

/** Canonical path "resources.cpu_cores". */
export const InfoFieldResourcesCpuCores = "resources.cpu_cores" as const;

/** Canonical path "resources.cpu_limit_cores". */
export const InfoFieldResourcesCpuLimitCores = "resources.cpu_limit_cores" as const;

/** Canonical path "resources.memory_mib". */
export const InfoFieldResourcesMemoryMib = "resources.memory_mib" as const;

/** Canonical path "resources.memory_limit_mib". */
export const InfoFieldResourcesMemoryLimitMib = "resources.memory_limit_mib" as const;

/** Canonical path "resources.disk_mib". */
export const InfoFieldResourcesDiskMib = "resources.disk_mib" as const;

/** Canonical path "labels". */
export const InfoFieldLabels = "labels" as const;

/** Canonical path "created_at". */
export const InfoFieldCreatedAt = "created_at" as const;

/** Canonical path "provider_metadata". */
export const InfoFieldProviderMetadata = "provider_metadata" as const;

/** Canonical path "state". */
export const InfoFieldState = "state" as const;

/** Canonical path "terminal". */
export const InfoFieldTerminal = "terminal" as const;
