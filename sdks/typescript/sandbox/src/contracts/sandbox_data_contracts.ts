// Portable sandbox data contracts come from standard protobuf generation.
// These aliases keep the SDK entrypoint concise while ensuring one representation.

import type { DurationFieldsToMs } from "#internal/duration_ms.js";
import type {
  EmptyMarkersAsTrue,
  FlattenAllowlistFields,
  PlainStringMapFields,
} from "#internal/input_remap.js";
import type { StripProto } from "#internal/strip_proto.js";
import {
  EgressMode as ProtoEgressMode,
  IsolationKind as ProtoIsolationKind,
  LifecycleCapability as ProtoLifecycleCapability,
  PolicyMode as ProtoPolicyMode,
  PortTransport as ProtoPortTransport,
  SandboxState as ProtoSandboxState,
  SecretInjection as ProtoSecretInjection,
  SnapshotKind as ProtoSnapshotKind,
  WaitCondition as ProtoWaitCondition,
  type ArchiveOptions as ProtoArchiveOptions,
  type CreateOptions as ProtoCreateOptions,
  type CreateResult as ProtoCreateResult,
  type DeleteOptions as ProtoDeleteOptions,
  type ErrorInfo as ProtoErrorInfo,
  type ExitResult as ProtoExitResult,
  type IsolationConfig as ProtoIsolationConfig,
  type MetadataList as ProtoMetadataList,
  type MetadataObject as ProtoMetadataObject,
  type MetadataValue as ProtoMetadataValue,
  type ObservabilityConfig as ProtoObservabilityConfig,
  type PauseOptions as ProtoPauseOptions,
  type Placement as ProtoPlacement,
  type PortBinding as ProtoPortBinding,
  type Resources as ProtoResources,
  type RuntimeConfig as ProtoRuntimeConfig,
  type SandboxInfo as ProtoSandboxInfo,
  type SandboxSource as ProtoSandboxSource,
  type SecretReference as ProtoSecretReference,
  type SecurityConfig as ProtoSecurityConfig,
  type StartOptions as ProtoStartOptions,
  type StopOptions as ProtoStopOptions,
  type TerminateOptions as ProtoTerminateOptions,
  type ReadinessProbe as ProtoReadinessProbe,
} from "@sandbox-kit/contracts/sandbox";
import type { StripAndFlatten } from "#internal/oneof_input.js";
import { protoEnum } from "#internal/proto_enum.js";

/** Public error payload; identity fields stay required after StripProto softening. */
export type ErrorInfo = StripProto<ProtoErrorInfo> & {
  kind: NonNullable<StripProto<ProtoErrorInfo>["kind"]>;
  provider: string;
  operation: string;
  field: string;
  message: string;
};
export type RuntimeConfig = StripProto<ProtoRuntimeConfig>;
export type IsolationConfig = StripProto<ProtoIsolationConfig>;
export type Resources = StripProto<ProtoResources>;
export type Placement = StripProto<ProtoPlacement>;
export type PortBinding = StripProto<ProtoPortBinding>;
export type SecurityConfig = StripProto<ProtoSecurityConfig>;
export type SecretReference = StripProto<ProtoSecretReference>;
export type ObservabilityConfig = StripProto<ProtoObservabilityConfig>;
export type SandboxInfo = StripProto<ProtoSandboxInfo>;
export type CreateResult = StripProto<ProtoCreateResult>;

/**
 * Public lifecycle options for `Backend` and provider adapters (ms timeouts).
 * Same duration remapping as {@link CreateOptions}; wire `Duration` stays private.
 */
export type StartOptions = DurationFieldsToMs<StripProto<ProtoStartOptions>>;
export type StopOptions = DurationFieldsToMs<StripProto<ProtoStopOptions>>;
export type PauseOptions = DurationFieldsToMs<StripProto<ProtoPauseOptions>>;
export type TerminateOptions = DurationFieldsToMs<
  StripProto<ProtoTerminateOptions>
>;
export type ArchiveOptions = DurationFieldsToMs<
  StripProto<ProtoArchiveOptions>
>;
export type DeleteOptions = DurationFieldsToMs<StripProto<ProtoDeleteOptions>>;
export type ExitResult = StripProto<ProtoExitResult>;
export type MetadataValue = StripProto<ProtoMetadataValue>;
export type MetadataList = StripProto<ProtoMetadataList>;
export type MetadataObject = StripProto<ProtoMetadataObject>;

/**
 * Public create request for `Client.create`.
 * Composed from reusable input transforms: strip + flatten oneofs → plain maps /
 * allowlists / empty markers → ms durations.
 */
export type CreateOptions = DurationFieldsToMs<
  EmptyMarkersAsTrue<
    FlattenAllowlistFields<
      PlainStringMapFields<
        StripAndFlatten<
          ProtoCreateOptions,
          {
            source: NonNullable<ProtoSandboxSource["kind"]>;
            readiness: NonNullable<ProtoReadinessProbe["check"]>;
          }
        >
      >
    >
  >
>;

/**
 * Public network create input (allowlists are plain `string[]`).
 * Derived from {@link CreateOptions}; present empty arrays mean deny-all.
 */
export type NetworkConfigInput = NonNullable<CreateOptions["network"]>;

/** Public create source: set exactly one variant (no protobuf `{ case, value }`). */
export type SandboxSourceInput = NonNullable<CreateOptions["source"]>;
/** Public readiness: set tcpPort or command (no protobuf `{ case, value }`). */
export type ReadinessProbeInput = NonNullable<CreateOptions["readiness"]>;

/**
 * Public lifetime create input (ms durations).
 * Derived from {@link CreateOptions}; wire `Duration` fields stay private.
 */
export type LifetimePolicyInput = NonNullable<CreateOptions["lifetime"]>;
/**
 * Public idle/stopped policy leaf under {@link LifetimePolicyInput} (`after` in ms).
 */
export type AutomaticActionInput = NonNullable<
  LifetimePolicyInput["idleTerminate"]
>;
/**
 * Public provisioning create input (ms timeouts).
 * Derived from {@link CreateOptions}; wire `Duration` fields stay private.
 */
export type ProvisioningOptionsInput = NonNullable<
  CreateOptions["provisioning"]
>;

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
