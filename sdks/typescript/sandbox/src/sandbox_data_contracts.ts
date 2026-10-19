// Portable sandbox data contracts come from standard protobuf generation.
// These aliases keep the SDK entrypoint concise while ensuring one representation.

import type { MessageInitShape } from "@bufbuild/protobuf";
import type { Duration } from "@bufbuild/protobuf/wkt";
import type { DurationFieldsToMs } from "./duration_ms.js";
import {
  ArchiveOptionsSchema,
  AutomaticActionSchema,
  CommandProbeSchema,
  CreateOptionsSchema,
  CreateResultSchema,
  DeleteOptionsSchema,
  EgressMode as ProtoEgressMode,
  ExitResultSchema,
  ImageSourceSchema,
  IsolationConfigSchema,
  IsolationKind as ProtoIsolationKind,
  LifecycleCapability as ProtoLifecycleCapability,
  LifetimePolicySchema,
  MetadataListSchema,
  MetadataObjectSchema,
  MetadataValueSchema,
  NetworkConfigSchema,
  ObservabilityConfigSchema,
  PauseOptionsSchema,
  PlacementSchema,
  PolicyMode as ProtoPolicyMode,
  PortBindingSchema,
  PortTransport as ProtoPortTransport,
  ProviderDefaultSourceSchema,
  ProvisioningOptionsSchema,
  ReadinessProbeSchema,
  ResourcesSchema,
  RuntimeConfigSchema,
  SandboxInfoSchema,
  SandboxSourceSchema,
  SandboxState as ProtoSandboxState,
  SecretInjection as ProtoSecretInjection,
  SecretReferenceSchema,
  SecurityConfigSchema,
  SnapshotKind as ProtoSnapshotKind,
  SnapshotSourceSchema,
  StartOptionsSchema,
  StopOptionsSchema,
  StringAllowlistSchema,
  TerminateOptionsSchema,
  ValueOrigin as ProtoValueOrigin,
  WaitCondition as ProtoWaitCondition,
  WarmPoolSourceSchema,
  type ArchiveOptions as ProtoArchiveOptions,
  type AutomaticAction as ProtoAutomaticAction,
  type CommandProbe as ProtoCommandProbe,
  type CreateOptions as ProtoCreateOptions,
  type CreateResult as ProtoCreateResult,
  type DeleteOptions as ProtoDeleteOptions,
  type ErrorInfo as ProtoErrorInfo,
  type ExitResult as ProtoExitResult,
  type ImageSource as ProtoImageSource,
  type IsolationConfig as ProtoIsolationConfig,
  type LifetimePolicy as ProtoLifetimePolicy,
  type MetadataList as ProtoMetadataList,
  type MetadataObject as ProtoMetadataObject,
  type MetadataValue as ProtoMetadataValue,
  type NetworkConfig as ProtoNetworkConfig,
  type ObservabilityConfig as ProtoObservabilityConfig,
  type PauseOptions as ProtoPauseOptions,
  type Placement as ProtoPlacement,
  type PortBinding as ProtoPortBinding,
  type ProviderDefaultSource as ProtoProviderDefaultSource,
  type ProvisioningOptions as ProtoProvisioningOptions,
  type ReadinessProbe as ProtoReadinessProbe,
  type Resources as ProtoResources,
  type RuntimeConfig as ProtoRuntimeConfig,
  type SandboxInfo as ProtoSandboxInfo,
  type SandboxSource as ProtoSandboxSource,
  type SecretReference as ProtoSecretReference,
  type SecurityConfig as ProtoSecurityConfig,
  type SnapshotSource as ProtoSnapshotSource,
  type StartOptions as ProtoStartOptions,
  type StopOptions as ProtoStopOptions,
  type StringAllowlist as ProtoStringAllowlist,
  type TerminateOptions as ProtoTerminateOptions,
  type WarmPoolSource as ProtoWarmPoolSource,
} from "@sandbox-kit/contracts/kit/sandbox/v1/sandbox_pb.js";
import type { FlattenOneof } from "./oneof_input.js";
import type { MetadataPlain } from "./metadata_plain.js";

export type { Duration };

export type ErrorInfo = Omit<ProtoErrorInfo, "$typeName">;
export type ProviderDefaultSource = Omit<
  ProtoProviderDefaultSource,
  "$typeName"
>;
export type ImageSource = Omit<ProtoImageSource, "$typeName">;
export type SnapshotSource = Omit<ProtoSnapshotSource, "$typeName">;
export type WarmPoolSource = Omit<ProtoWarmPoolSource, "$typeName">;
/** Public create source: set exactly one variant (no protobuf `{ case, value }`). */
export type SandboxSourceInput = FlattenOneof<ProtoSandboxSource["kind"]>;
/** Public readiness: set tcpPort or command (no protobuf `{ case, value }`). */
export type ReadinessProbeInput = FlattenOneof<ProtoReadinessProbe["check"]>;
/** Public create request: nested oneofs + millisecond timeouts (not Duration). */
export type CreateOptionsInput = DurationFieldsToMs<
  Omit<
    MessageInitShape<typeof CreateOptionsSchema>,
    "source" | "readiness" | "providerOptions"
  >
> & {
  source?: SandboxSourceInput;
  readiness?: ReadinessProbeInput;
  providerOptions?: Record<string, MetadataPlain>;
};
// Proto3 maps/repeated/enums are required on Message; keep them optional on
// aliases so request construction matches MessageInitShape and Go nil maps.
export type RuntimeConfig = Omit<
  ProtoRuntimeConfig,
  "$typeName" | "entrypoint"
> & {
  entrypoint?: string[];
};
export type IsolationConfig = Omit<ProtoIsolationConfig, "$typeName">;
export type Resources = Omit<ProtoResources, "$typeName">;
export type Placement = Omit<ProtoPlacement, "$typeName" | "regions"> & {
  regions?: string[];
};
export type StringAllowlist = Omit<
  ProtoStringAllowlist,
  "$typeName" | "entries"
> & {
  entries?: string[];
};
export type PortBinding = Omit<ProtoPortBinding, "$typeName">;
export type SecretReference = Omit<ProtoSecretReference, "$typeName">;
export type AutomaticAction = Omit<ProtoAutomaticAction, "$typeName">;
export type CommandProbe = Omit<ProtoCommandProbe, "$typeName" | "argv"> & {
  argv?: string[];
};
export type ObservabilityConfig = Omit<ProtoObservabilityConfig, "$typeName">;
/** Wire lifecycle options (Duration timeouts) for providers / Backend. */
export type StartOptions = Omit<ProtoStartOptions, "$typeName">;
export type StopOptions = Omit<ProtoStopOptions, "$typeName">;
export type PauseOptions = Omit<ProtoPauseOptions, "$typeName">;
export type TerminateOptions = Omit<ProtoTerminateOptions, "$typeName">;
export type ArchiveOptions = Omit<ProtoArchiveOptions, "$typeName">;
export type DeleteOptions = Omit<ProtoDeleteOptions, "$typeName">;
/** Public lifecycle options: timeouts are milliseconds. */
export type StartOptionsInput = DurationFieldsToMs<StartOptions>;
export type StopOptionsInput = DurationFieldsToMs<StopOptions>;
export type PauseOptionsInput = DurationFieldsToMs<PauseOptions>;
export type TerminateOptionsInput = DurationFieldsToMs<TerminateOptions>;
export type ArchiveOptionsInput = DurationFieldsToMs<ArchiveOptions>;
export type DeleteOptionsInput = DurationFieldsToMs<DeleteOptions>;
export type ExitResult = Omit<ProtoExitResult, "$typeName">;
export type MetadataValue = Omit<ProtoMetadataValue, "$typeName">;
export type MetadataList = Omit<ProtoMetadataList, "$typeName">;
export type MetadataObject = Omit<ProtoMetadataObject, "$typeName">;

export type SandboxSource = Omit<ProtoSandboxSource, "$typeName" | "kind"> & {
  kind:
    | { case: "providerDefault"; value: ProviderDefaultSource }
    | { case: "image"; value: ImageSource }
    | { case: "snapshot"; value: SnapshotSource }
    | { case: "warmPool"; value: WarmPoolSource }
    | { case: undefined; value?: undefined };
};
export type NetworkConfig = Omit<
  ProtoNetworkConfig,
  | "$typeName"
  | "egress"
  | "outboundCidrs"
  | "outboundDomains"
  | "inboundCidrs"
  | "ports"
> & {
  egress?: ProtoEgressMode;
  outboundCidrs?: StringAllowlist;
  outboundDomains?: StringAllowlist;
  inboundCidrs?: StringAllowlist;
  ports?: PortBinding[];
};
export type SecurityConfig = Omit<
  ProtoSecurityConfig,
  "$typeName" | "secrets"
> & {
  secrets?: SecretReference[];
};
export type LifetimePolicy = Omit<
  ProtoLifetimePolicy,
  | "$typeName"
  | "idleTerminate"
  | "idleStop"
  | "idlePause"
  | "stoppedArchive"
  | "stoppedDelete"
> & {
  idleTerminate?: AutomaticAction;
  idleStop?: AutomaticAction;
  idlePause?: AutomaticAction;
  stoppedArchive?: AutomaticAction;
  stoppedDelete?: AutomaticAction;
};
export type ProvisioningOptions = Omit<ProtoProvisioningOptions, "$typeName">;
export type ReadinessProbe = Omit<
  ProtoReadinessProbe,
  "$typeName" | "check"
> & {
  check:
    | { case: "tcpPort"; value: number }
    | { case: "command"; value: CommandProbe }
    | { case: undefined; value?: undefined };
};
export type CreateOptions = Omit<
  ProtoCreateOptions,
  | "$typeName"
  | "environment"
  | "labels"
  | "source"
  | "runtime"
  | "isolation"
  | "resources"
  | "placement"
  | "network"
  | "security"
  | "lifetime"
  | "provisioning"
  | "readiness"
  | "observability"
  | "providerOptions"
> & {
  environment?: { [key: string]: string };
  labels?: { [key: string]: string };
  source?: SandboxSource;
  runtime?: RuntimeConfig;
  isolation?: IsolationConfig;
  resources?: Resources;
  placement?: Placement;
  network?: NetworkConfig;
  security?: SecurityConfig;
  lifetime?: LifetimePolicy;
  provisioning?: ProvisioningOptions;
  readiness?: ReadinessProbe;
  observability?: ObservabilityConfig;
  providerOptions?: MetadataObject;
};
// Remap nested message fields so constructed info does not require $typeName.
export type SandboxInfo = Omit<
  ProtoSandboxInfo,
  | "$typeName"
  | "labels"
  | "origins"
  | "isolation"
  | "resources"
  | "providerMetadata"
  | "createdAt"
> & {
  labels?: { [key: string]: string };
  origins?: { [key: string]: ProtoValueOrigin };
  isolation?: IsolationConfig;
  resources?: Resources;
  providerMetadata?: MetadataObject;
  createdAt?: Date;
};
export type CreateResult = Omit<ProtoCreateResult, "$typeName" | "sandbox"> & {
  sandbox?: SandboxInfo;
};

// Message schemas stay package-internal (Protovalidate / materialize). They are
// not part of the public application-facing API.

// Ergonomic names reference protobuf enum values; they do not define a second contract.
export const SnapshotKind = {
  Default: ProtoSnapshotKind.DEFAULT,
  Filesystem: ProtoSnapshotKind.FILESYSTEM,
  Memory: ProtoSnapshotKind.MEMORY,
  Image: ProtoSnapshotKind.IMAGE,
  Machine: ProtoSnapshotKind.MACHINE,
  Checkpoint: ProtoSnapshotKind.CHECKPOINT,
  ProviderSpecific: ProtoSnapshotKind.PROVIDER_SPECIFIC,
} as const;
export type SnapshotKind = (typeof SnapshotKind)[keyof typeof SnapshotKind];
export const SnapshotKindNumber = SnapshotKind;

export const IsolationKind = {
  Default: ProtoIsolationKind.DEFAULT,
  Container: ProtoIsolationKind.CONTAINER,
  Microvm: ProtoIsolationKind.MICROVM,
  LinuxVm: ProtoIsolationKind.LINUX_VM,
  Worker: ProtoIsolationKind.WORKER,
  ProviderSpecific: ProtoIsolationKind.PROVIDER_SPECIFIC,
} as const;
export type IsolationKind = (typeof IsolationKind)[keyof typeof IsolationKind];
export const IsolationKindNumber = IsolationKind;

export const EgressMode = {
  Default: ProtoEgressMode.DEFAULT,
  AllowAll: ProtoEgressMode.ALLOW_ALL,
  BlockAll: ProtoEgressMode.BLOCK_ALL,
  Restricted: ProtoEgressMode.RESTRICTED,
} as const;
export type EgressMode = (typeof EgressMode)[keyof typeof EgressMode];
export const EgressModeNumber = EgressMode;

export const PortTransport = {
  Default: ProtoPortTransport.DEFAULT,
  Tls: ProtoPortTransport.TLS,
  Http2Tls: ProtoPortTransport.HTTP2_TLS,
  Plain: ProtoPortTransport.PLAIN,
} as const;
export type PortTransport = (typeof PortTransport)[keyof typeof PortTransport];
export const PortTransportNumber = PortTransport;

export const SecretInjection = {
  Default: ProtoSecretInjection.DEFAULT,
  Environment: ProtoSecretInjection.ENVIRONMENT,
  EgressPlaceholder: ProtoSecretInjection.EGRESS_PLACEHOLDER,
} as const;
export type SecretInjection =
  (typeof SecretInjection)[keyof typeof SecretInjection];
export const SecretInjectionNumber = SecretInjection;

export const PolicyMode = {
  Default: ProtoPolicyMode.DEFAULT,
  Disabled: ProtoPolicyMode.DISABLED,
  After: ProtoPolicyMode.AFTER,
} as const;
export type PolicyMode = (typeof PolicyMode)[keyof typeof PolicyMode];
export const PolicyModeNumber = PolicyMode;

export const WaitCondition = {
  Default: ProtoWaitCondition.DEFAULT,
  Submitted: ProtoWaitCondition.SUBMITTED,
  Scheduled: ProtoWaitCondition.SCHEDULED,
  Started: ProtoWaitCondition.STARTED,
  Ready: ProtoWaitCondition.READY,
} as const;
export type WaitCondition = (typeof WaitCondition)[keyof typeof WaitCondition];
export const WaitConditionNumber = WaitCondition;

export const SandboxState = {
  Unknown: ProtoSandboxState.UNKNOWN,
  Queued: ProtoSandboxState.QUEUED,
  Starting: ProtoSandboxState.STARTING,
  Running: ProtoSandboxState.RUNNING,
  Stopping: ProtoSandboxState.STOPPING,
  Stopped: ProtoSandboxState.STOPPED,
  Pausing: ProtoSandboxState.PAUSING,
  Paused: ProtoSandboxState.PAUSED,
  Archiving: ProtoSandboxState.ARCHIVING,
  Archived: ProtoSandboxState.ARCHIVED,
  Terminating: ProtoSandboxState.TERMINATING,
  Terminated: ProtoSandboxState.TERMINATED,
  Deleting: ProtoSandboxState.DELETING,
  Deleted: ProtoSandboxState.DELETED,
  Error: ProtoSandboxState.ERROR,
} as const;
export type SandboxState = (typeof SandboxState)[keyof typeof SandboxState];

export const LifecycleCapability = {
  Unknown: ProtoLifecycleCapability.UNKNOWN,
  Refresh: ProtoLifecycleCapability.REFRESH,
  Start: ProtoLifecycleCapability.START,
  Stop: ProtoLifecycleCapability.STOP,
  Pause: ProtoLifecycleCapability.PAUSE,
  Terminate: ProtoLifecycleCapability.TERMINATE,
  Archive: ProtoLifecycleCapability.ARCHIVE,
  Delete: ProtoLifecycleCapability.DELETE,
} as const;
export type LifecycleCapability =
  (typeof LifecycleCapability)[keyof typeof LifecycleCapability];

export const ValueOrigin = {
  Unknown: ProtoValueOrigin.UNKNOWN,
  Provider: ProtoValueOrigin.PROVIDER,
  Request: ProtoValueOrigin.REQUEST,
} as const;
export type ValueOrigin = (typeof ValueOrigin)[keyof typeof ValueOrigin];
