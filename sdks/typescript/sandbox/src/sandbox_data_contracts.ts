// Portable sandbox data contracts come from standard protobuf generation.
// These aliases keep the SDK entrypoint concise while ensuring one representation.

import type { DurationFieldsToMs } from "./duration_ms.js";
import type {
  EmptyMarkersAsTrue,
  FlattenAllowlistFields,
  PlainStringMapFields,
  RemapFields,
} from "./input_remap.js";
import type { StripProto } from "./strip_proto.js";
import {
  EgressMode as ProtoEgressMode,
  IsolationKind as ProtoIsolationKind,
  LifecycleCapability as ProtoLifecycleCapability,
  PolicyMode as ProtoPolicyMode,
  PortTransport as ProtoPortTransport,
  SandboxState as ProtoSandboxState,
  SecretInjection as ProtoSecretInjection,
  SnapshotKind as ProtoSnapshotKind,
  ValueOrigin as ProtoValueOrigin,
  WaitCondition as ProtoWaitCondition,
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
import type { StripAndFlatten } from "./oneof_input.js";
import type { MetadataPlain } from "./metadata_plain.js";
import { protoEnum } from "./proto_enum.js";

/** Public error payload; identity fields stay required after StripProto softening. */
export type ErrorInfo = StripProto<ProtoErrorInfo> & {
  kind: NonNullable<StripProto<ProtoErrorInfo>["kind"]>;
  provider: string;
  operation: string;
  field: string;
  message: string;
};
export type ProviderDefaultSource = StripProto<ProtoProviderDefaultSource>;
export type ImageSource = StripProto<ProtoImageSource>;
export type SnapshotSource = StripProto<ProtoSnapshotSource>;
export type WarmPoolSource = StripProto<ProtoWarmPoolSource>;
export type RuntimeConfig = StripProto<ProtoRuntimeConfig>;
export type IsolationConfig = StripProto<ProtoIsolationConfig>;
export type Resources = StripProto<ProtoResources>;
export type Placement = StripProto<ProtoPlacement>;
export type StringAllowlist = StripProto<ProtoStringAllowlist>;
export type PortBinding = StripProto<ProtoPortBinding>;
export type SecretReference = StripProto<ProtoSecretReference>;
export type AutomaticAction = StripProto<ProtoAutomaticAction>;
export type CommandProbe = StripProto<ProtoCommandProbe>;
export type ObservabilityConfig = StripProto<ProtoObservabilityConfig>;
export type NetworkConfig = StripProto<ProtoNetworkConfig>;
export type SecurityConfig = StripProto<ProtoSecurityConfig>;
export type LifetimePolicy = StripProto<ProtoLifetimePolicy>;
export type ProvisioningOptions = StripProto<ProtoProvisioningOptions>;
export type SandboxSource = StripProto<ProtoSandboxSource>;
export type ReadinessProbe = StripProto<ProtoReadinessProbe>;
export type CreateOptions = StripProto<ProtoCreateOptions>;
export type SandboxInfo = StripProto<ProtoSandboxInfo>;
export type CreateResult = StripProto<ProtoCreateResult>;

/** Wire lifecycle options (Duration timeouts) for providers / Backend. */
export type StartOptions = StripProto<ProtoStartOptions>;
export type StopOptions = StripProto<ProtoStopOptions>;
export type PauseOptions = StripProto<ProtoPauseOptions>;
export type TerminateOptions = StripProto<ProtoTerminateOptions>;
export type ArchiveOptions = StripProto<ProtoArchiveOptions>;
export type DeleteOptions = StripProto<ProtoDeleteOptions>;
export type ExitResult = StripProto<ProtoExitResult>;
export type MetadataValue = StripProto<ProtoMetadataValue>;
export type MetadataList = StripProto<ProtoMetadataList>;
export type MetadataObject = StripProto<ProtoMetadataObject>;

/**
 * Public create request, composed from reusable input transforms:
 * strip + flatten oneofs → plain maps / allowlists / empty markers → ms durations.
 */
export type CreateOptionsInput = DurationFieldsToMs<
  EmptyMarkersAsTrue<
    FlattenAllowlistFields<
      PlainStringMapFields<
        RemapFields<
          StripAndFlatten<
            ProtoCreateOptions,
            {
              source: NonNullable<ProtoSandboxSource["kind"]>;
              readiness: NonNullable<ProtoReadinessProbe["check"]>;
            }
          >,
          { providerOptions: Record<string, MetadataPlain> }
        >
      >
    >
  >
>;

/**
 * Public network create input (allowlists are plain `string[]`).
 * Derived from {@link CreateOptionsInput}; present empty arrays mean deny-all.
 */
export type NetworkConfigInput = NonNullable<CreateOptionsInput["network"]>;

/** Public create source: set exactly one variant (no protobuf `{ case, value }`). */
export type SandboxSourceInput = NonNullable<CreateOptionsInput["source"]>;
/** Public readiness: set tcpPort or command (no protobuf `{ case, value }`). */
export type ReadinessProbeInput = NonNullable<CreateOptionsInput["readiness"]>;

// Message schemas stay package-internal (Protovalidate / materialize). They are
// not part of the public application-facing API.

// Ergonomic PascalCase names over protobuf enum values; not a second contract.
export const SnapshotKind = protoEnum(ProtoSnapshotKind);
export type SnapshotKind = (typeof SnapshotKind)[keyof typeof SnapshotKind];
export const SnapshotKindNumber = SnapshotKind;

export const IsolationKind = protoEnum(ProtoIsolationKind);
export type IsolationKind = (typeof IsolationKind)[keyof typeof IsolationKind];
export const IsolationKindNumber = IsolationKind;

export const EgressMode = protoEnum(ProtoEgressMode);
export type EgressMode = (typeof EgressMode)[keyof typeof EgressMode];
export const EgressModeNumber = EgressMode;

export const PortTransport = protoEnum(ProtoPortTransport);
export type PortTransport = (typeof PortTransport)[keyof typeof PortTransport];
export const PortTransportNumber = PortTransport;

export const SecretInjection = protoEnum(ProtoSecretInjection);
export type SecretInjection =
  (typeof SecretInjection)[keyof typeof SecretInjection];
export const SecretInjectionNumber = SecretInjection;

export const PolicyMode = protoEnum(ProtoPolicyMode);
export type PolicyMode = (typeof PolicyMode)[keyof typeof PolicyMode];
export const PolicyModeNumber = PolicyMode;

export const WaitCondition = protoEnum(ProtoWaitCondition);
export type WaitCondition = (typeof WaitCondition)[keyof typeof WaitCondition];
export const WaitConditionNumber = WaitCondition;

export const SandboxState = protoEnum(ProtoSandboxState);
export type SandboxState = (typeof SandboxState)[keyof typeof SandboxState];

export const LifecycleCapability = protoEnum(ProtoLifecycleCapability);
export type LifecycleCapability =
  (typeof LifecycleCapability)[keyof typeof LifecycleCapability];

export const ValueOrigin = protoEnum(ProtoValueOrigin);
export type ValueOrigin = (typeof ValueOrigin)[keyof typeof ValueOrigin];
