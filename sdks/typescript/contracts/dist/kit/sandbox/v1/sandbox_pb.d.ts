import type { GenEnum, GenFile, GenMessage } from "@bufbuild/protobuf/codegenv2";
import type { Duration, Timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";
/**
 * Describes the file kit/sandbox/v1/sandbox.proto.
 */
export declare const file_kit_sandbox_v1_sandbox: GenFile;
/**
 * is the creation request. It does not include lifecycle operations or credentials.
 * Omitted fields use provider defaults. A provider rejects settings it cannot honor.
 *
 * Example:
 *
 * 	options := &CreateOptions{
 * 	    Resources: &Resources{CPUCores: Value(0.5), MemoryMiB: Value(uint64(512))},
 * 	}
 *
 * @generated from message kit.sandbox.v1.CreateOptions
 */
export type CreateOptions = Message<"kit.sandbox.v1.CreateOptions"> & {
    /**
     * is an optional display name. Nil uses the provider default. An explicit empty string asks for an unnamed sandbox.
     *
     * @generated from field: optional string name = 1;
     */
    name?: string;
    /**
     * selects an image, snapshot, warm pool, or the provider default. Nil means the provider default.
     *
     * @generated from field: kit.sandbox.v1.SandboxSource source = 2;
     */
    source?: SandboxSource;
    /**
     * sets language, user, working directory, and entrypoint hints.
     *
     * @generated from field: kit.sandbox.v1.RuntimeConfig runtime = 3;
     */
    runtime?: RuntimeConfig;
    /**
     * selects the isolation class, such as a container or Linux VM.
     *
     * @generated from field: kit.sandbox.v1.IsolationConfig isolation = 4;
     */
    isolation?: IsolationConfig;
    /**
     * sets CPU, memory, and disk. Units are physical cores and MiB.
     *
     * @generated from field: kit.sandbox.v1.Resources resources = 5;
     */
    resources?: Resources;
    /**
     * sets cloud and region preferences. A preference is not a scheduling guarantee.
     *
     * @generated from field: kit.sandbox.v1.Placement placement = 6;
     */
    placement?: Placement;
    /**
     * is the process environment. An empty map is explicit; nil means omitted.
     *
     * @generated from field: map<string, string> environment = 7;
     */
    environment: {
        [key: string]: string;
    };
    /**
     * is caller-supplied labels or tags. An empty map is explicit; nil means omitted.
     *
     * @generated from field: map<string, string> labels = 8;
     */
    labels: {
        [key: string]: string;
    };
    /**
     * sets egress, allowlists, ports, and related network options.
     *
     * @generated from field: kit.sandbox.v1.NetworkConfig network = 9;
     */
    network?: NetworkConfig;
    /**
     * sets public access, workload identity, and secret references.
     *
     * @generated from field: kit.sandbox.v1.SecurityConfig security = 10;
     */
    security?: SecurityConfig;
    /**
     * sets lifetime and idle policies requested at creation. These are not stop or delete calls.
     *
     * @generated from field: kit.sandbox.v1.LifetimePolicy lifetime = 11;
     */
    lifetime?: LifetimePolicy;
    /**
     * sets the creation deadline, queue timeout, and wait condition. JSON name is "creation".
     *
     * @generated from field: kit.sandbox.v1.ProvisioningOptions provisioning = 12 [json_name = "creation"];
     */
    provisioning?: ProvisioningOptions;
    /**
     * @generated from field: kit.sandbox.v1.ReadinessProbe readiness = 14;
     */
    readiness?: ReadinessProbe;
    /**
     * @generated from field: kit.sandbox.v1.ObservabilityConfig observability = 15;
     */
    observability?: ObservabilityConfig;
    /**
     * holds provider-specific declarative options. It is never an SDK object or an auth credential.
     *
     * @generated from field: kit.sandbox.v1.MetadataObject provider_options = 16;
     */
    providerOptions?: MetadataObject;
};
/**
 * Describes the message kit.sandbox.v1.CreateOptions.
 * Use `create(CreateOptionsSchema)` to create a new message.
 */
export declare const CreateOptionsSchema: GenMessage<CreateOptions>;
/**
 * is the sandbox source. Set one alternative. Nil means the provider default, which need not be an image.
 *
 * @generated from message kit.sandbox.v1.SandboxSource
 */
export type SandboxSource = Message<"kit.sandbox.v1.SandboxSource"> & {
    /**
     * is the selected source. Leaving it unset also means the provider default, which need not be an image.
     *
     * @generated from oneof kit.sandbox.v1.SandboxSource.kind
     */
    kind: {
        /**
         * @generated from field: kit.sandbox.v1.ProviderDefaultSource provider_default = 1;
         */
        value: ProviderDefaultSource;
        case: "providerDefault";
    } | {
        /**
         * @generated from field: kit.sandbox.v1.ImageSource image = 2;
         */
        value: ImageSource;
        case: "image";
    } | {
        /**
         * @generated from field: kit.sandbox.v1.SnapshotSource snapshot = 3;
         */
        value: SnapshotSource;
        case: "snapshot";
    } | {
        /**
         * @generated from field: kit.sandbox.v1.WarmPoolSource warm_pool = 4;
         */
        value: WarmPoolSource;
        case: "warmPool";
    } | {
        case: undefined;
        value?: undefined;
    };
};
/**
 * Describes the message kit.sandbox.v1.SandboxSource.
 * Use `create(SandboxSourceSchema)` to create a new message.
 */
export declare const SandboxSourceSchema: GenMessage<SandboxSource>;
/**
 * asks the provider to choose its default source.
 *
 * @generated from message kit.sandbox.v1.ProviderDefaultSource
 */
export type ProviderDefaultSource = Message<"kit.sandbox.v1.ProviderDefaultSource"> & {};
/**
 * Describes the message kit.sandbox.v1.ProviderDefaultSource.
 * Use `create(ProviderDefaultSourceSchema)` to create a new message.
 */
export declare const ProviderDefaultSourceSchema: GenMessage<ProviderDefaultSource>;
/**
 * is a container image reference, such as "alpine:3.21".
 *
 * @generated from message kit.sandbox.v1.ImageSource
 */
export type ImageSource = Message<"kit.sandbox.v1.ImageSource"> & {
    /**
     * is the registry image reference.
     *
     * @generated from field: string reference = 1;
     */
    reference: string;
};
/**
 * Describes the message kit.sandbox.v1.ImageSource.
 * Use `create(ImageSourceSchema)` to create a new message.
 */
export declare const ImageSourceSchema: GenMessage<ImageSource>;
/**
 * is an existing snapshot. Reference identifies it. Kind classifies it.
 *
 * @generated from message kit.sandbox.v1.SnapshotSource
 */
export type SnapshotSource = Message<"kit.sandbox.v1.SnapshotSource"> & {
    /**
     * is the provider snapshot identifier.
     *
     * @generated from field: string reference = 1;
     */
    reference: string;
    /**
     * classifies the snapshot. Default leaves the class to the provider.
     *
     * @generated from field: kit.sandbox.v1.SnapshotKind kind = 2;
     */
    kind: SnapshotKind;
    /**
     * is a provider-specific kind name when Kind is provider-specific.
     *
     * @generated from field: optional string provider_kind = 3;
     */
    providerKind?: string;
};
/**
 * Describes the message kit.sandbox.v1.SnapshotSource.
 * Use `create(SnapshotSourceSchema)` to create a new message.
 */
export declare const SnapshotSourceSchema: GenMessage<SnapshotSource>;
/**
 * is an existing warm pool. Reference identifies it.
 *
 * @generated from message kit.sandbox.v1.WarmPoolSource
 */
export type WarmPoolSource = Message<"kit.sandbox.v1.WarmPoolSource"> & {
    /**
     * is the provider warm-pool identifier.
     *
     * @generated from field: string reference = 1;
     */
    reference: string;
};
/**
 * Describes the message kit.sandbox.v1.WarmPoolSource.
 * Use `create(WarmPoolSourceSchema)` to create a new message.
 */
export declare const WarmPoolSourceSchema: GenMessage<WarmPoolSource>;
/**
 * is process setup applied when the sandbox starts. Language is a hint, not an install promise.
 *
 * @generated from message kit.sandbox.v1.RuntimeConfig
 */
export type RuntimeConfig = Message<"kit.sandbox.v1.RuntimeConfig"> & {
    /**
     * is a language or toolbox hint. It does not promise that a runtime is installed.
     *
     * @generated from field: optional string language = 1;
     */
    language?: string;
    /**
     * @generated from field: optional string version = 2;
     */
    version?: string;
    /**
     * @generated from field: optional string user = 3;
     */
    user?: string;
    /**
     * @generated from field: optional string working_directory = 4;
     */
    workingDirectory?: string;
    /**
     * @generated from field: repeated string entrypoint = 5;
     */
    entrypoint: string[];
    /**
     * @generated from field: optional bool pty = 6;
     */
    pty?: boolean;
};
/**
 * Describes the message kit.sandbox.v1.RuntimeConfig.
 * Use `create(RuntimeConfigSchema)` to create a new message.
 */
export declare const RuntimeConfigSchema: GenMessage<RuntimeConfig>;
/**
 * selects how the sandbox is isolated. Kind default leaves the class to the provider.
 *
 * @generated from message kit.sandbox.v1.IsolationConfig
 */
export type IsolationConfig = Message<"kit.sandbox.v1.IsolationConfig"> & {
    /**
     * @generated from field: kit.sandbox.v1.IsolationKind kind = 1;
     */
    kind: IsolationKind;
    /**
     * @generated from field: optional string provider_kind = 2;
     */
    providerKind?: string;
    /**
     * @generated from field: optional bool nested_virtualization = 3;
     */
    nestedVirtualization?: boolean;
};
/**
 * Describes the message kit.sandbox.v1.IsolationConfig.
 * Use `create(IsolationConfigSchema)` to create a new message.
 */
export declare const IsolationConfigSchema: GenMessage<IsolationConfig>;
/**
 * is the requested compute size. Omitted fields use provider defaults.
 * GPU requests and resizing are not creation settings.
 *
 * @generated from message kit.sandbox.v1.Resources
 */
export type Resources = Message<"kit.sandbox.v1.Resources"> & {
    /**
     * is the requested physical CPU count. Modal allows fractions. Daytona requires whole cores.
     *
     * @generated from field: optional double cpu_cores = 1;
     */
    cpuCores?: number;
    /**
     * is an optional CPU ceiling. A positive limit requires a request and cannot be smaller than it.
     *
     * @generated from field: optional double cpu_limit_cores = 2;
     */
    cpuLimitCores?: number;
    /**
     * is requested memory in MiB. Providers convert this to the native unit.
     *
     * @generated from field: optional uint64 memory_mib = 3;
     */
    memoryMib?: bigint;
    /**
     * is an optional memory ceiling in MiB. A positive limit requires a request.
     *
     * @generated from field: optional uint64 memory_limit_mib = 4;
     */
    memoryLimitMib?: bigint;
    /**
     * is requested disk in MiB. Modal rejects a creation override. Daytona requires whole GiB.
     *
     * @generated from field: optional uint64 disk_mib = 5;
     */
    diskMib?: bigint;
};
/**
 * Describes the message kit.sandbox.v1.Resources.
 * Use `create(ResourcesSchema)` to create a new message.
 */
export declare const ResourcesSchema: GenMessage<Resources>;
/**
 * is a placement preference. The provider may still choose another location.
 *
 * @generated from message kit.sandbox.v1.Placement
 */
export type Placement = Message<"kit.sandbox.v1.Placement"> & {
    /**
     * is a cloud name such as aws, gcp, or oci. Supported names depend on the provider.
     *
     * @generated from field: optional string cloud = 1;
     */
    cloud?: string;
    /**
     * lists region preferences in order. It does not promise multi-region support.
     *
     * @generated from field: repeated string regions = 2;
     */
    regions: string[];
    /**
     * is a provider target. Daytona maps the client region here. Per-request placement is unsupported.
     *
     * @generated from field: optional string target = 3;
     */
    target?: string;
};
/**
 * Describes the message kit.sandbox.v1.Placement.
 * Use `create(PlacementSchema)` to create a new message.
 */
export declare const PlacementSchema: GenMessage<Placement>;
/**
 * is an allowlist. A present empty list means deny all. Nil means the list was omitted.
 *
 * @generated from message kit.sandbox.v1.StringAllowlist
 */
export type StringAllowlist = Message<"kit.sandbox.v1.StringAllowlist"> & {
    /**
     * lists allowed values. A present empty list means deny all. Nil means the list was omitted.
     *
     * @generated from field: repeated string entries = 1;
     */
    entries: string[];
};
/**
 * Describes the message kit.sandbox.v1.StringAllowlist.
 * Use `create(StringAllowlistSchema)` to create a new message.
 */
export declare const StringAllowlistSchema: GenMessage<StringAllowlist>;
/**
 * is creation-time networking. Empty allowlists do not mean the same thing as omitted ones.
 *
 * @generated from message kit.sandbox.v1.NetworkConfig
 */
export type NetworkConfig = Message<"kit.sandbox.v1.NetworkConfig"> & {
    /**
     * selects the outbound policy. Default leaves it to the provider.
     *
     * @generated from field: kit.sandbox.v1.EgressMode egress = 1;
     */
    egress: EgressMode;
    /**
     * limits outbound traffic to these CIDR ranges.
     *
     * @generated from field: kit.sandbox.v1.StringAllowlist outbound_cidrs = 2;
     */
    outboundCidrs?: StringAllowlist;
    /**
     * limits outbound traffic to these domains.
     *
     * @generated from field: kit.sandbox.v1.StringAllowlist outbound_domains = 3;
     */
    outboundDomains?: StringAllowlist;
    /**
     * limits inbound traffic to these CIDR ranges.
     *
     * @generated from field: kit.sandbox.v1.StringAllowlist inbound_cidrs = 4;
     */
    inboundCidrs?: StringAllowlist;
    /**
     * asks for a private network when the provider supports one.
     *
     * @generated from field: optional bool private_network = 5;
     */
    privateNetwork?: boolean;
    /**
     * names another sandbox to link. The provider must already know that sandbox.
     *
     * @generated from field: optional string linked_sandbox = 6;
     */
    linkedSandbox?: string;
    /**
     * names a provider-managed proxy.
     *
     * @generated from field: optional string proxy_reference = 7;
     */
    proxyReference?: string;
    /**
     * is an outbound proxy URL.
     *
     * @generated from field: optional string outbound_proxy_url = 8;
     */
    outboundProxyUrl?: string;
    /**
     * is a custom domain for the sandbox, when the provider supports one.
     *
     * @generated from field: optional string custom_domain = 9;
     */
    customDomain?: string;
    /**
     * publishes these ports.
     *
     * @generated from field: repeated kit.sandbox.v1.PortBinding ports = 10;
     */
    ports: PortBinding[];
};
/**
 * Describes the message kit.sandbox.v1.NetworkConfig.
 * Use `create(NetworkConfigSchema)` to create a new message.
 */
export declare const NetworkConfigSchema: GenMessage<NetworkConfig>;
/**
 * publishes one port. Port is the container port. Transport selects the protocol.
 *
 * @generated from message kit.sandbox.v1.PortBinding
 */
export type PortBinding = Message<"kit.sandbox.v1.PortBinding"> & {
    /**
     * is the port number to publish.
     *
     * @generated from field: uint32 port = 1;
     */
    port: number;
    /**
     * selects TLS, HTTP/2, or plain transport. Default leaves the choice to the provider.
     *
     * @generated from field: kit.sandbox.v1.PortTransport transport = 2;
     */
    transport: PortTransport;
};
/**
 * Describes the message kit.sandbox.v1.PortBinding.
 * Use `create(PortBindingSchema)` to create a new message.
 */
export declare const PortBindingSchema: GenMessage<PortBinding>;
/**
 * sets access and secret references. Secrets are provider-managed names, not credential values.
 *
 * @generated from message kit.sandbox.v1.SecurityConfig
 */
export type SecurityConfig = Message<"kit.sandbox.v1.SecurityConfig"> & {
    /**
     * requests public access when true. Nil leaves the provider default.
     *
     * @generated from field: optional bool public_access = 1;
     */
    publicAccess?: boolean;
    /**
     * requests a workload identity when the provider supports one.
     *
     * @generated from field: optional bool workload_identity = 2;
     */
    workloadIdentity?: boolean;
    /**
     * names provider-managed secrets. These are not credential values.
     *
     * @generated from field: repeated kit.sandbox.v1.SecretReference secrets = 3;
     */
    secrets: SecretReference[];
};
/**
 * Describes the message kit.sandbox.v1.SecurityConfig.
 * Use `create(SecurityConfigSchema)` to create a new message.
 */
export declare const SecurityConfigSchema: GenMessage<SecurityConfig>;
/**
 * names a provider-managed secret and how it is injected.
 *
 * @generated from message kit.sandbox.v1.SecretReference
 */
export type SecretReference = Message<"kit.sandbox.v1.SecretReference"> & {
    /**
     * is the provider secret name.
     *
     * @generated from field: string reference = 1;
     */
    reference: string;
    /**
     * is the environment variable that receives the secret, when injection uses the environment.
     *
     * @generated from field: optional string environment_variable = 2;
     */
    environmentVariable?: string;
    /**
     * selects environment injection or an egress placeholder.
     *
     * @generated from field: kit.sandbox.v1.SecretInjection injection = 3;
     */
    injection: SecretInjection;
};
/**
 * Describes the message kit.sandbox.v1.SecretReference.
 * Use `create(SecretReferenceSchema)` to create a new message.
 */
export declare const SecretReferenceSchema: GenMessage<SecretReference>;
/**
 * is requested at creation. It is not a stop, pause, or delete call.
 *
 * @generated from message kit.sandbox.v1.LifetimePolicy
 */
export type LifetimePolicy = Message<"kit.sandbox.v1.LifetimePolicy"> & {
    /**
     * is the maximum lifetime. Nil uses the provider default. Explicit zero means unlimited when the provider supports that.
     *
     * @generated from field: google.protobuf.Duration maximum_lifetime = 1;
     */
    maximumLifetime?: Duration;
    /**
     * requests deletion when the sandbox stops, where the provider supports it.
     *
     * @generated from field: optional bool ephemeral = 2;
     */
    ephemeral?: boolean;
    /**
     * ends the sandbox after it is idle.
     *
     * @generated from field: kit.sandbox.v1.AutomaticAction idle_terminate = 3;
     */
    idleTerminate?: AutomaticAction;
    /**
     * stops the sandbox after it is idle. Daytona maps a disabled stop to native zero.
     *
     * @generated from field: kit.sandbox.v1.AutomaticAction idle_stop = 4;
     */
    idleStop?: AutomaticAction;
    /**
     * pauses the sandbox after it is idle. A positive delay conflicts with a positive stop delay.
     *
     * @generated from field: kit.sandbox.v1.AutomaticAction idle_pause = 5;
     */
    idlePause?: AutomaticAction;
    /**
     * archives a stopped sandbox after a delay. Daytona rejects disable and immediate archive.
     *
     * @generated from field: kit.sandbox.v1.AutomaticAction stopped_archive = 6;
     */
    stoppedArchive?: AutomaticAction;
    /**
     * deletes a stopped sandbox. Daytona maps immediate delete to native zero.
     *
     * @generated from field: kit.sandbox.v1.AutomaticAction stopped_delete = 7;
     */
    stoppedDelete?: AutomaticAction;
};
/**
 * Describes the message kit.sandbox.v1.LifetimePolicy.
 * Use `create(LifetimePolicySchema)` to create a new message.
 */
export declare const LifetimePolicySchema: GenMessage<LifetimePolicy>;
/**
 * is one idle or stopped policy. Disabled, immediate, and delayed do not share one zero value.
 *
 * @generated from message kit.sandbox.v1.AutomaticAction
 */
export type AutomaticAction = Message<"kit.sandbox.v1.AutomaticAction"> & {
    /**
     * is default, disabled, or after a delay. Default omits the policy.
     *
     * @generated from field: kit.sandbox.v1.PolicyMode mode = 1;
     */
    mode: PolicyMode;
    /**
     * is the delay used when Mode is after. Zero means immediate, not disabled.
     *
     * @generated from field: google.protobuf.Duration after = 2;
     */
    after?: Duration;
};
/**
 * Describes the message kit.sandbox.v1.AutomaticAction.
 * Use `create(AutomaticActionSchema)` to create a new message.
 */
export declare const AutomaticActionSchema: GenMessage<AutomaticAction>;
/**
 * is creation waiting, separate from sandbox lifetime.
 *
 * @generated from message kit.sandbox.v1.ProvisioningOptions
 */
export type ProvisioningOptions = Message<"kit.sandbox.v1.ProvisioningOptions"> & {
    /**
     * is the client operation deadline, distinct from sandbox lifetime and queue time.
     * Nil and explicit zero add no Kit deadline.
     *
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
    /**
     * is how long creation may wait in a provider queue. It is not the sandbox lifetime.
     *
     * @generated from field: google.protobuf.Duration queue_timeout = 2;
     */
    queueTimeout?: Duration;
    /**
     * is the condition creation waits for before it returns.
     *
     * @generated from field: kit.sandbox.v1.WaitCondition wait_for = 3;
     */
    waitFor: WaitCondition;
};
/**
 * Describes the message kit.sandbox.v1.ProvisioningOptions.
 * Use `create(ProvisioningOptionsSchema)` to create a new message.
 */
export declare const ProvisioningOptionsSchema: GenMessage<ProvisioningOptions>;
/**
 * is the check used when waiting until ready. Set either a TCP port or a command.
 *
 * @generated from message kit.sandbox.v1.ReadinessProbe
 */
export type ReadinessProbe = Message<"kit.sandbox.v1.ReadinessProbe"> & {
    /**
     * @generated from oneof kit.sandbox.v1.ReadinessProbe.check
     */
    check: {
        /**
         * @generated from field: uint32 tcp_port = 1;
         */
        value: number;
        case: "tcpPort";
    } | {
        /**
         * @generated from field: kit.sandbox.v1.CommandProbe command = 2;
         */
        value: CommandProbe;
        case: "command";
    } | {
        case: undefined;
        value?: undefined;
    };
};
/**
 * Describes the message kit.sandbox.v1.ReadinessProbe.
 * Use `create(ReadinessProbeSchema)` to create a new message.
 */
export declare const ReadinessProbeSchema: GenMessage<ReadinessProbe>;
/**
 * is a command readiness check. Argv is the command and its arguments.
 *
 * @generated from message kit.sandbox.v1.CommandProbe
 */
export type CommandProbe = Message<"kit.sandbox.v1.CommandProbe"> & {
    /**
     * @generated from field: repeated string argv = 1;
     */
    argv: string[];
};
/**
 * Describes the message kit.sandbox.v1.CommandProbe.
 * Use `create(CommandProbeSchema)` to create a new message.
 */
export declare const CommandProbeSchema: GenMessage<CommandProbe>;
/**
 * sets provider logging and telemetry options.
 *
 * @generated from message kit.sandbox.v1.ObservabilityConfig
 */
export type ObservabilityConfig = Message<"kit.sandbox.v1.ObservabilityConfig"> & {
    /**
     * @generated from field: optional bool verbose = 1;
     */
    verbose?: boolean;
    /**
     * @generated from field: optional string telemetry_endpoint = 2;
     */
    telemetryEndpoint?: string;
};
/**
 * Describes the message kit.sandbox.v1.ObservabilityConfig.
 * Use `create(ObservabilityConfigSchema)` to create a new message.
 */
export declare const ObservabilityConfigSchema: GenMessage<ObservabilityConfig>;
/**
 * is the metadata returned by a provider create call.
 *
 * @generated from message kit.sandbox.v1.CreateResult
 */
export type CreateResult = Message<"kit.sandbox.v1.CreateResult"> & {
    /**
     * is the created sandbox. ID must be present before a public handle is built.
     *
     * @generated from field: kit.sandbox.v1.SandboxInfo sandbox = 1;
     */
    sandbox?: SandboxInfo;
};
/**
 * Describes the message kit.sandbox.v1.CreateResult.
 * Use `create(CreateResultSchema)` to create a new message.
 */
export declare const CreateResultSchema: GenMessage<CreateResult>;
/**
 * configures starting a persistent sandbox. Nil timeout inherits Config.Timeout.
 *
 * @generated from message kit.sandbox.v1.StartOptions
 */
export type StartOptions = Message<"kit.sandbox.v1.StartOptions"> & {
    /**
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
};
/**
 * Describes the message kit.sandbox.v1.StartOptions.
 * Use `create(StartOptionsSchema)` to create a new message.
 */
export declare const StartOptionsSchema: GenMessage<StartOptions>;
/**
 * configures a restartable stop. Force requests an immediate provider-native stop.
 *
 * @generated from message kit.sandbox.v1.StopOptions
 */
export type StopOptions = Message<"kit.sandbox.v1.StopOptions"> & {
    /**
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
    /**
     * @generated from field: optional bool force = 2;
     */
    force?: boolean;
};
/**
 * Describes the message kit.sandbox.v1.StopOptions.
 * Use `create(StopOptionsSchema)` to create a new message.
 */
export declare const StopOptionsSchema: GenMessage<StopOptions>;
/**
 * configures a memory-preserving pause. Nil timeout inherits Config.Timeout.
 *
 * @generated from message kit.sandbox.v1.PauseOptions
 */
export type PauseOptions = Message<"kit.sandbox.v1.PauseOptions"> & {
    /**
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
};
/**
 * Describes the message kit.sandbox.v1.PauseOptions.
 * Use `create(PauseOptionsSchema)` to create a new message.
 */
export declare const PauseOptionsSchema: GenMessage<PauseOptions>;
/**
 * configures permanent completion of the current execution.
 *
 * @generated from message kit.sandbox.v1.TerminateOptions
 */
export type TerminateOptions = Message<"kit.sandbox.v1.TerminateOptions"> & {
    /**
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
    /**
     * @generated from field: optional bool wait = 2;
     */
    wait?: boolean;
};
/**
 * Describes the message kit.sandbox.v1.TerminateOptions.
 * Use `create(TerminateOptionsSchema)` to create a new message.
 */
export declare const TerminateOptionsSchema: GenMessage<TerminateOptions>;
/**
 * configures moving retained sandbox state to provider archive storage.
 *
 * @generated from message kit.sandbox.v1.ArchiveOptions
 */
export type ArchiveOptions = Message<"kit.sandbox.v1.ArchiveOptions"> & {
    /**
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
};
/**
 * Describes the message kit.sandbox.v1.ArchiveOptions.
 * Use `create(ArchiveOptionsSchema)` to create a new message.
 */
export declare const ArchiveOptionsSchema: GenMessage<ArchiveOptions>;
/**
 * configures irreversible sandbox deletion.
 *
 * @generated from message kit.sandbox.v1.DeleteOptions
 */
export type DeleteOptions = Message<"kit.sandbox.v1.DeleteOptions"> & {
    /**
     * @generated from field: google.protobuf.Duration timeout = 1;
     */
    timeout?: Duration;
    /**
     * @generated from field: optional bool wait = 2;
     */
    wait?: boolean;
};
/**
 * Describes the message kit.sandbox.v1.DeleteOptions.
 * Use `create(DeleteOptionsSchema)` to create a new message.
 */
export declare const DeleteOptionsSchema: GenMessage<DeleteOptions>;
/**
 * is the result of permanently completing an execution.
 *
 * @generated from message kit.sandbox.v1.ExitResult
 */
export type ExitResult = Message<"kit.sandbox.v1.ExitResult"> & {
    /**
     * @generated from field: optional int32 exit_code = 1;
     */
    exitCode?: number;
};
/**
 * Describes the message kit.sandbox.v1.ExitResult.
 * Use `create(ExitResultSchema)` to create a new message.
 */
export declare const ExitResultSchema: GenMessage<ExitResult>;
/**
 * is shared metadata. Optional fields stay absent when the provider did not report them.
 *
 * Example:
 *
 * 	info := instance.Info()
 * 	fmt.Println(info.GetID(), info.Origins[InfoFieldName])
 *
 * @generated from message kit.sandbox.v1.SandboxInfo
 */
export type SandboxInfo = Message<"kit.sandbox.v1.SandboxInfo"> & {
    /**
     * is the provider sandbox identifier. Creation fails locally when it is empty.
     *
     * @generated from field: string id = 1;
     */
    id: string;
    /**
     * is the selected provider identifier.
     *
     * @generated from field: string provider = 2;
     */
    provider: string;
    /**
     * is the sandbox name when one is known. Nil means it was not reported.
     *
     * @generated from field: optional string name = 3;
     */
    name?: string;
    /**
     * is the provider's own state string. It is not a shared lifecycle state machine.
     *
     * @generated from field: optional string provider_state = 4;
     */
    providerState?: string;
    /**
     * @generated from field: optional string region = 5;
     */
    region?: string;
    /**
     * @generated from field: kit.sandbox.v1.IsolationConfig isolation = 6;
     */
    isolation?: IsolationConfig;
    /**
     * @generated from field: kit.sandbox.v1.Resources resources = 7;
     */
    resources?: Resources;
    /**
     * @generated from field: map<string, string> labels = 8;
     */
    labels: {
        [key: string]: string;
    };
    /**
     * @generated from field: google.protobuf.Timestamp created_at = 9;
     */
    createdAt?: Timestamp;
    /**
     * @generated from field: kit.sandbox.v1.MetadataObject provider_metadata = 10;
     */
    providerMetadata?: MetadataObject;
    /**
     * records where each present field came from. A missing entry means the origin is unknown.
     *
     * @generated from field: map<string, kit.sandbox.v1.ValueOrigin> origins = 11;
     */
    origins: {
        [key: string]: ValueOrigin;
    };
    /**
     * is the normalized lifecycle state when Kit can map the provider state.
     *
     * @generated from field: kit.sandbox.v1.SandboxState state = 12;
     */
    state: SandboxState;
    /**
     * reports that execution/resource lifecycle is permanently complete.
     *
     * @generated from field: optional bool terminal = 13;
     */
    terminal?: boolean;
};
/**
 * Describes the message kit.sandbox.v1.SandboxInfo.
 * Use `create(SandboxInfoSchema)` to create a new message.
 */
export declare const SandboxInfoSchema: GenMessage<SandboxInfo>;
/**
 * is the portable error detail. The native cause stays on the runtime error and is not serialized.
 * Unknown classifications do not imply that a retry is safe.
 *
 * Example:
 *
 * 	err := NewError(ErrorInfo{
 * 	    Kind:    ErrorKindUnsupported,
 * 	    Field:   CreateFieldResourcesDiskMiB,
 * 	    Message: "disk overrides are unsupported",
 * 	}, cause)
 *
 * @generated from message kit.sandbox.v1.ErrorInfo
 */
export type ErrorInfo = Message<"kit.sandbox.v1.ErrorInfo"> & {
    /**
     * is the shared category, such as invalid argument or unsupported.
     *
     * @generated from field: kit.sandbox.v1.ErrorKind kind = 1;
     */
    kind: ErrorKind;
    /**
     * is the provider identifier, such as "modal" or "daytona". Empty when the failure is local.
     *
     * @generated from field: string provider = 2;
     */
    provider: string;
    /**
     * is the Kit operation, such as "create", "initialize", or "close".
     *
     * @generated from field: string operation = 3;
     */
    operation: string;
    /**
     * is the canonical field path involved, such as "resources.cpu_cores". Empty when no field applies.
     *
     * @generated from field: string field = 4;
     */
    field: string;
    /**
     * is the human-readable explanation.
     *
     * @generated from field: string message = 5;
     */
    message: string;
    /**
     * is the native HTTP or transport status when the provider reported one.
     *
     * @generated from field: optional uint32 status_code = 6;
     */
    statusCode?: number;
    /**
     * is the provider's own error code when one was reported.
     *
     * @generated from field: optional string provider_code = 7;
     */
    providerCode?: string;
    /**
     * is the provider component that reported the error, when available.
     *
     * @generated from field: optional string provider_source = 8;
     */
    providerSource?: string;
    /**
     * is a stable portable validation or policy rule identifier, when applicable.
     *
     * @generated from field: optional string rule = 9;
     */
    rule?: string;
};
/**
 * Describes the message kit.sandbox.v1.ErrorInfo.
 * Use `create(ErrorInfoSchema)` to create a new message.
 */
export declare const ErrorInfoSchema: GenMessage<ErrorInfo>;
/**
 * Exact numeric alternatives avoid protobuf Struct's double-only representation.
 * An empty value represents null.
 *
 * @generated from message kit.sandbox.v1.MetadataValue
 */
export type MetadataValue = Message<"kit.sandbox.v1.MetadataValue"> & {
    /**
     * @generated from oneof kit.sandbox.v1.MetadataValue.value
     */
    value: {
        /**
         * @generated from field: bool boolean = 1;
         */
        value: boolean;
        case: "boolean";
    } | {
        /**
         * @generated from field: string text = 2;
         */
        value: string;
        case: "text";
    } | {
        /**
         * @generated from field: int64 signed_integer = 3;
         */
        value: bigint;
        case: "signedInteger";
    } | {
        /**
         * @generated from field: uint64 unsigned_integer = 4;
         */
        value: bigint;
        case: "unsignedInteger";
    } | {
        /**
         * @generated from field: double number = 5;
         */
        value: number;
        case: "number";
    } | {
        /**
         * @generated from field: bytes binary = 6;
         */
        value: Uint8Array;
        case: "binary";
    } | {
        /**
         * @generated from field: kit.sandbox.v1.MetadataList list = 7;
         */
        value: MetadataList;
        case: "list";
    } | {
        /**
         * @generated from field: kit.sandbox.v1.MetadataObject object = 8;
         */
        value: MetadataObject;
        case: "object";
    } | {
        case: undefined;
        value?: undefined;
    };
};
/**
 * Describes the message kit.sandbox.v1.MetadataValue.
 * Use `create(MetadataValueSchema)` to create a new message.
 */
export declare const MetadataValueSchema: GenMessage<MetadataValue>;
/**
 * is an ordered list of metadata values.
 *
 * @generated from message kit.sandbox.v1.MetadataList
 */
export type MetadataList = Message<"kit.sandbox.v1.MetadataList"> & {
    /**
     * is the list contents. An empty list is explicit.
     *
     * @generated from field: repeated kit.sandbox.v1.MetadataValue values = 1;
     */
    values: MetadataValue[];
};
/**
 * Describes the message kit.sandbox.v1.MetadataList.
 * Use `create(MetadataListSchema)` to create a new message.
 */
export declare const MetadataListSchema: GenMessage<MetadataList>;
/**
 * is a string-keyed metadata object. Go exposes this as map[string]any.
 *
 * @generated from message kit.sandbox.v1.MetadataObject
 */
export type MetadataObject = Message<"kit.sandbox.v1.MetadataObject"> & {
    /**
     * is the object contents. An empty object is explicit.
     *
     * @generated from field: map<string, kit.sandbox.v1.MetadataValue> values = 1;
     */
    values: {
        [key: string]: MetadataValue;
    };
};
/**
 * Describes the message kit.sandbox.v1.MetadataObject.
 * Use `create(MetadataObjectSchema)` to create a new message.
 */
export declare const MetadataObjectSchema: GenMessage<MetadataObject>;
/**
 * classifies a snapshot. Preserve the numeric values.
 *
 * @generated from enum kit.sandbox.v1.SnapshotKind
 */
export declare enum SnapshotKind {
    /**
     * @generated from enum value: SNAPSHOT_KIND_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: SNAPSHOT_KIND_FILESYSTEM = 1;
     */
    FILESYSTEM = 1,
    /**
     * @generated from enum value: SNAPSHOT_KIND_MEMORY = 2;
     */
    MEMORY = 2,
    /**
     * @generated from enum value: SNAPSHOT_KIND_IMAGE = 3;
     */
    IMAGE = 3,
    /**
     * @generated from enum value: SNAPSHOT_KIND_MACHINE = 4;
     */
    MACHINE = 4,
    /**
     * @generated from enum value: SNAPSHOT_KIND_CHECKPOINT = 5;
     */
    CHECKPOINT = 5,
    /**
     * @generated from enum value: SNAPSHOT_KIND_PROVIDER_SPECIFIC = 6;
     */
    PROVIDER_SPECIFIC = 6
}
/**
 * Describes the enum kit.sandbox.v1.SnapshotKind.
 */
export declare const SnapshotKindSchema: GenEnum<SnapshotKind>;
/**
 * is the isolation class. Preserve the numeric values.
 *
 * @generated from enum kit.sandbox.v1.IsolationKind
 */
export declare enum IsolationKind {
    /**
     * @generated from enum value: ISOLATION_KIND_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: ISOLATION_KIND_CONTAINER = 1;
     */
    CONTAINER = 1,
    /**
     * @generated from enum value: ISOLATION_KIND_MICROVM = 2;
     */
    MICROVM = 2,
    /**
     * @generated from enum value: ISOLATION_KIND_LINUX_VM = 3;
     */
    LINUX_VM = 3,
    /**
     * @generated from enum value: ISOLATION_KIND_WORKER = 4;
     */
    WORKER = 4,
    /**
     * @generated from enum value: ISOLATION_KIND_PROVIDER_SPECIFIC = 5;
     */
    PROVIDER_SPECIFIC = 5
}
/**
 * Describes the enum kit.sandbox.v1.IsolationKind.
 */
export declare const IsolationKindSchema: GenEnum<IsolationKind>;
/**
 * controls outbound traffic. Default leaves the policy to the provider.
 *
 * @generated from enum kit.sandbox.v1.EgressMode
 */
export declare enum EgressMode {
    /**
     * @generated from enum value: EGRESS_MODE_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: EGRESS_MODE_ALLOW_ALL = 1;
     */
    ALLOW_ALL = 1,
    /**
     * blocks all networking, including private connections.
     *
     * @generated from enum value: EGRESS_MODE_BLOCK_ALL = 2;
     */
    BLOCK_ALL = 2,
    /**
     * @generated from enum value: EGRESS_MODE_RESTRICTED = 3;
     */
    RESTRICTED = 3
}
/**
 * Describes the enum kit.sandbox.v1.EgressMode.
 */
export declare const EgressModeSchema: GenEnum<EgressMode>;
/**
 * is the protocol used for a published port.
 *
 * @generated from enum kit.sandbox.v1.PortTransport
 */
export declare enum PortTransport {
    /**
     * @generated from enum value: PORT_TRANSPORT_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: PORT_TRANSPORT_TLS = 1;
     */
    TLS = 1,
    /**
     * @generated from enum value: PORT_TRANSPORT_HTTP2_TLS = 2;
     */
    HTTP2_TLS = 2,
    /**
     * @generated from enum value: PORT_TRANSPORT_PLAIN = 3;
     */
    PLAIN = 3
}
/**
 * Describes the enum kit.sandbox.v1.PortTransport.
 */
export declare const PortTransportSchema: GenEnum<PortTransport>;
/**
 * selects how a secret is exposed to the sandbox.
 *
 * @generated from enum kit.sandbox.v1.SecretInjection
 */
export declare enum SecretInjection {
    /**
     * @generated from enum value: SECRET_INJECTION_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: SECRET_INJECTION_ENVIRONMENT = 1;
     */
    ENVIRONMENT = 1,
    /**
     * @generated from enum value: SECRET_INJECTION_EGRESS_PLACEHOLDER = 2;
     */
    EGRESS_PLACEHOLDER = 2
}
/**
 * Describes the enum kit.sandbox.v1.SecretInjection.
 */
export declare const SecretInjectionSchema: GenEnum<SecretInjection>;
/**
 * selects whether an automatic action is omitted, disabled, or delayed.
 *
 * @generated from enum kit.sandbox.v1.PolicyMode
 */
export declare enum PolicyMode {
    /**
     * @generated from enum value: POLICY_MODE_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: POLICY_MODE_DISABLED = 1;
     */
    DISABLED = 1,
    /**
     * @generated from enum value: POLICY_MODE_AFTER = 2;
     */
    AFTER = 2
}
/**
 * Describes the enum kit.sandbox.v1.PolicyMode.
 */
export declare const PolicyModeSchema: GenEnum<PolicyMode>;
/**
 * is how far creation waits before returning. Ready can require a probe and a positive timeout.
 *
 * @generated from enum kit.sandbox.v1.WaitCondition
 */
export declare enum WaitCondition {
    /**
     * @generated from enum value: WAIT_CONDITION_DEFAULT = 0;
     */
    DEFAULT = 0,
    /**
     * @generated from enum value: WAIT_CONDITION_SUBMITTED = 1;
     */
    SUBMITTED = 1,
    /**
     * @generated from enum value: WAIT_CONDITION_SCHEDULED = 2;
     */
    SCHEDULED = 2,
    /**
     * @generated from enum value: WAIT_CONDITION_STARTED = 3;
     */
    STARTED = 3,
    /**
     * @generated from enum value: WAIT_CONDITION_READY = 4;
     */
    READY = 4
}
/**
 * Describes the enum kit.sandbox.v1.WaitCondition.
 */
export declare const WaitConditionSchema: GenEnum<WaitCondition>;
/**
 * is a portable lifecycle state. ProviderState retains the native value.
 *
 * @generated from enum kit.sandbox.v1.SandboxState
 */
export declare enum SandboxState {
    /**
     * @generated from enum value: SANDBOX_STATE_UNKNOWN = 0;
     */
    UNKNOWN = 0,
    /**
     * @generated from enum value: SANDBOX_STATE_QUEUED = 1;
     */
    QUEUED = 1,
    /**
     * @generated from enum value: SANDBOX_STATE_STARTING = 2;
     */
    STARTING = 2,
    /**
     * @generated from enum value: SANDBOX_STATE_RUNNING = 3;
     */
    RUNNING = 3,
    /**
     * @generated from enum value: SANDBOX_STATE_STOPPING = 4;
     */
    STOPPING = 4,
    /**
     * @generated from enum value: SANDBOX_STATE_STOPPED = 5;
     */
    STOPPED = 5,
    /**
     * @generated from enum value: SANDBOX_STATE_PAUSING = 6;
     */
    PAUSING = 6,
    /**
     * @generated from enum value: SANDBOX_STATE_PAUSED = 7;
     */
    PAUSED = 7,
    /**
     * @generated from enum value: SANDBOX_STATE_ARCHIVING = 8;
     */
    ARCHIVING = 8,
    /**
     * @generated from enum value: SANDBOX_STATE_ARCHIVED = 9;
     */
    ARCHIVED = 9,
    /**
     * @generated from enum value: SANDBOX_STATE_TERMINATING = 10;
     */
    TERMINATING = 10,
    /**
     * @generated from enum value: SANDBOX_STATE_TERMINATED = 11;
     */
    TERMINATED = 11,
    /**
     * @generated from enum value: SANDBOX_STATE_DELETING = 12;
     */
    DELETING = 12,
    /**
     * @generated from enum value: SANDBOX_STATE_DELETED = 13;
     */
    DELETED = 13,
    /**
     * @generated from enum value: SANDBOX_STATE_ERROR = 14;
     */
    ERROR = 14
}
/**
 * Describes the enum kit.sandbox.v1.SandboxState.
 */
export declare const SandboxStateSchema: GenEnum<SandboxState>;
/**
 * identifies an operation supported by a particular resolved sandbox handle.
 *
 * @generated from enum kit.sandbox.v1.LifecycleCapability
 */
export declare enum LifecycleCapability {
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_UNKNOWN = 0;
     */
    UNKNOWN = 0,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_REFRESH = 1;
     */
    REFRESH = 1,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_START = 2;
     */
    START = 2,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_STOP = 3;
     */
    STOP = 3,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_PAUSE = 4;
     */
    PAUSE = 4,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_TERMINATE = 5;
     */
    TERMINATE = 5,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_ARCHIVE = 6;
     */
    ARCHIVE = 6,
    /**
     * @generated from enum value: LIFECYCLE_CAPABILITY_DELETE = 7;
     */
    DELETE = 7
}
/**
 * Describes the enum kit.sandbox.v1.LifecycleCapability.
 */
export declare const LifecycleCapabilitySchema: GenEnum<LifecycleCapability>;
/**
 * records whether a metadata field came from the provider or from the request.
 *
 * @generated from enum kit.sandbox.v1.ValueOrigin
 */
export declare enum ValueOrigin {
    /**
     * @generated from enum value: VALUE_ORIGIN_UNKNOWN = 0;
     */
    UNKNOWN = 0,
    /**
     * @generated from enum value: VALUE_ORIGIN_PROVIDER = 1;
     */
    PROVIDER = 1,
    /**
     * @generated from enum value: VALUE_ORIGIN_REQUEST = 2;
     */
    REQUEST = 2
}
/**
 * Describes the enum kit.sandbox.v1.ValueOrigin.
 */
export declare const ValueOriginSchema: GenEnum<ValueOrigin>;
/**
 * is the shared error category. It does not imply that a retry is safe.
 *
 * @generated from enum kit.sandbox.v1.ErrorKind
 */
export declare enum ErrorKind {
    /**
     * @generated from enum value: ERROR_KIND_UNKNOWN = 0;
     */
    UNKNOWN = 0,
    /**
     * @generated from enum value: ERROR_KIND_INVALID_ARGUMENT = 1;
     */
    INVALID_ARGUMENT = 1,
    /**
     * @generated from enum value: ERROR_KIND_UNSUPPORTED = 2;
     */
    UNSUPPORTED = 2,
    /**
     * @generated from enum value: ERROR_KIND_AUTHENTICATION = 3;
     */
    AUTHENTICATION = 3,
    /**
     * @generated from enum value: ERROR_KIND_PERMISSION_DENIED = 4;
     */
    PERMISSION_DENIED = 4,
    /**
     * @generated from enum value: ERROR_KIND_NOT_FOUND = 5;
     */
    NOT_FOUND = 5,
    /**
     * @generated from enum value: ERROR_KIND_ALREADY_EXISTS = 6;
     */
    ALREADY_EXISTS = 6,
    /**
     * @generated from enum value: ERROR_KIND_CONFLICT = 7;
     */
    CONFLICT = 7,
    /**
     * @generated from enum value: ERROR_KIND_RATE_LIMITED = 8;
     */
    RATE_LIMITED = 8,
    /**
     * @generated from enum value: ERROR_KIND_RESOURCE_EXHAUSTED = 9;
     */
    RESOURCE_EXHAUSTED = 9,
    /**
     * @generated from enum value: ERROR_KIND_TIMEOUT = 10;
     */
    TIMEOUT = 10,
    /**
     * @generated from enum value: ERROR_KIND_CANCELED = 11;
     */
    CANCELED = 11,
    /**
     * @generated from enum value: ERROR_KIND_UNAVAILABLE = 12;
     */
    UNAVAILABLE = 12,
    /**
     * @generated from enum value: ERROR_KIND_INVALID_RESPONSE = 13;
     */
    INVALID_RESPONSE = 13,
    /**
     * @generated from enum value: ERROR_KIND_INTERNAL = 14;
     */
    INTERNAL = 14,
    /**
     * @generated from enum value: ERROR_KIND_FAILED_PRECONDITION = 15;
     */
    FAILED_PRECONDITION = 15
}
/**
 * Describes the enum kit.sandbox.v1.ErrorKind.
 */
export declare const ErrorKindSchema: GenEnum<ErrorKind>;
