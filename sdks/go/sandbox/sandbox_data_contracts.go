package sandbox

import (
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
)

type (
	// CreateOptions is the request passed to Client.Create.
	// A nil request uses the provider default. Modal requires an image in Source.
	CreateOptions = sandboxv1.CreateOptions

	// SandboxSource chooses how the sandbox is created.
	// Set Kind to ProviderDefault, Image, Snapshot, or WarmPool.
	SandboxSource = sandboxv1.SandboxSource

	// ProviderDefault asks the provider to choose the sandbox source.
	// Daytona Create with a nil request uses the default snapshot.
	ProviderDefault = sandboxv1.SandboxSource_ProviderDefault

	// Image starts the sandbox from a container image reference.
	Image = sandboxv1.SandboxSource_Image

	// Snapshot starts the sandbox from a saved snapshot.
	Snapshot = sandboxv1.SandboxSource_Snapshot

	// WarmPool starts the sandbox from a warm pool reference.
	WarmPool = sandboxv1.SandboxSource_WarmPool

	// ProviderDefaultSource is the empty source used with ProviderDefault.
	ProviderDefaultSource = sandboxv1.ProviderDefaultSource

	// ImageSource names the image reference, such as "python:3.11" or "alpine:3.21".
	ImageSource = sandboxv1.ImageSource

	// SnapshotSource names a snapshot and its SnapshotKind.
	SnapshotSource = sandboxv1.SnapshotSource

	// WarmPoolSource names a warm pool reference.
	WarmPoolSource = sandboxv1.WarmPoolSource

	// RuntimeConfig sets the language, user, working directory, and entrypoint.
	RuntimeConfig = sandboxv1.RuntimeConfig

	// IsolationConfig selects the sandbox class, such as a container or a microVM.
	IsolationConfig = sandboxv1.IsolationConfig

	// Resources sets CPU, memory, and disk for a create request.
	// CpuCores is a core count. MemoryMib and DiskMib are mebibytes.
	// Daytona accepts whole cores, and memory and disk only in whole gibibytes. Setting resources requires an image.
	// Modal accepts at least 0.125 cores and 128 MiB, and rejects DiskMib.
	Resources = sandboxv1.Resources

	// Placement selects the cloud, regions, or a provider target.
	Placement = sandboxv1.Placement

	// StringAllowlist is an explicit list of allowed strings, such as domains or CIDRs.
	StringAllowlist = sandboxv1.StringAllowlist

	// NetworkConfig sets egress and published ports.
	NetworkConfig = sandboxv1.NetworkConfig

	// PortBinding publishes one port and names its transport.
	PortBinding = sandboxv1.PortBinding

	// SecurityConfig sets secrets and other security options on a create request.
	SecurityConfig = sandboxv1.SecurityConfig

	// SecretReference names a secret and how the provider injects it.
	SecretReference = sandboxv1.SecretReference

	// LifetimePolicy limits how long the sandbox may run and what happens afterward.
	LifetimePolicy = sandboxv1.LifetimePolicy

	// AutomaticAction is one action in a lifetime policy, such as stop or delete after a delay.
	AutomaticAction = sandboxv1.AutomaticAction

	// ProvisioningOptions sets the create timeout and the condition to wait for.
	ProvisioningOptions = sandboxv1.ProvisioningOptions

	// ReadinessProbe checks that a sandbox is ready on a TCP port or by running a command.
	ReadinessProbe = sandboxv1.ReadinessProbe

	// TCPPort checks readiness by connecting to a TCP port.
	TCPPort = sandboxv1.ReadinessProbe_TcpPort

	// Command checks readiness by running a command.
	Command = sandboxv1.ReadinessProbe_Command

	// CommandProbe is the argv used by a Command readiness check.
	CommandProbe = sandboxv1.CommandProbe

	// ObservabilityConfig sets verbose logging and a telemetry endpoint.
	ObservabilityConfig = sandboxv1.ObservabilityConfig

	// CreateResult is the provider response from a successful create.
	// Sandbox holds the metadata. Client.Create returns a handle built from it.
	CreateResult = sandboxv1.CreateResult

	// StartOptions sets the timeout for Sandbox.Start.
	// A nil options value uses the client timeout.
	StartOptions = sandboxv1.StartOptions

	// StopOptions sets the timeout and Force flag for Sandbox.Stop.
	// Set Force with Value(true) to request an immediate stop.
	StopOptions = sandboxv1.StopOptions

	// PauseOptions sets the timeout for Sandbox.Pause.
	PauseOptions = sandboxv1.PauseOptions

	// TerminateOptions sets the timeout and Wait flag for Sandbox.Terminate.
	// Set Wait with Value(true) to wait for the exit code.
	TerminateOptions = sandboxv1.TerminateOptions

	// ArchiveOptions sets the timeout for Sandbox.Archive.
	ArchiveOptions = sandboxv1.ArchiveOptions

	// DeleteOptions sets the timeout for Sandbox.Delete.
	DeleteOptions = sandboxv1.DeleteOptions

	// ExitResult is the outcome of Sandbox.Terminate, including the exit code when Wait is set.
	ExitResult = sandboxv1.ExitResult

	// SandboxInfo is the metadata for one sandbox.
	// Sandbox.Info returns a copy. State is a SandboxState. Origins records who set each field.
	SandboxInfo = sandboxv1.SandboxInfo

	// MetadataValue is one metadata value: a boolean, text, number, bytes, list, or object.
	MetadataValue = sandboxv1.MetadataValue

	// MetadataBoolean is the boolean form of MetadataValue.Kind.
	MetadataBoolean = sandboxv1.MetadataValue_Boolean

	// MetadataText is the string form of MetadataValue.Kind.
	MetadataText = sandboxv1.MetadataValue_Text

	// MetadataSignedInteger is the signed-integer form of MetadataValue.Kind.
	MetadataSignedInteger = sandboxv1.MetadataValue_SignedInteger

	// MetadataUnsignedInteger is the unsigned-integer form of MetadataValue.Kind.
	MetadataUnsignedInteger = sandboxv1.MetadataValue_UnsignedInteger

	// MetadataNumber is the floating-point form of MetadataValue.Kind.
	MetadataNumber = sandboxv1.MetadataValue_Number

	// MetadataBinary is the byte-slice form of MetadataValue.Kind.
	MetadataBinary = sandboxv1.MetadataValue_Binary

	// MetadataListValue is the list form of MetadataValue.Kind.
	MetadataListValue = sandboxv1.MetadataValue_List

	// MetadataObjectValue is the nested-object form of MetadataValue.Kind.
	MetadataObjectValue = sandboxv1.MetadataValue_Object

	// MetadataList is an ordered list of metadata values.
	MetadataList = sandboxv1.MetadataList

	// MetadataObject is a string-keyed metadata map.
	// Build one with MetadataObjectFromMap and read one with MetadataObjectAsMap.
	MetadataObject = sandboxv1.MetadataObject

	// SnapshotKind classifies a SnapshotSource.
	SnapshotKind = sandboxv1.SnapshotKind

	// IsolationKind classifies IsolationConfig, such as a container or a microVM.
	IsolationKind = sandboxv1.IsolationKind

	// EgressMode controls outbound network access on NetworkConfig.
	EgressMode = sandboxv1.EgressMode

	// PortTransport selects how a published port is served.
	PortTransport = sandboxv1.PortTransport

	// SecretInjection selects how a SecretReference is delivered into the sandbox.
	SecretInjection = sandboxv1.SecretInjection

	// PolicyMode selects whether an AutomaticAction is disabled or runs after a delay.
	PolicyMode = sandboxv1.PolicyMode

	// WaitCondition is the create condition named by ProvisioningOptions.
	WaitCondition = sandboxv1.WaitCondition

	// SandboxState is the lifecycle state reported on SandboxInfo.
	SandboxState = sandboxv1.SandboxState

	// LifecycleCapability names one operation Sandbox.Supports can check.
	// Check Pause and Archive before calling them. A terminal handle supports nothing further.
	LifecycleCapability = sandboxv1.LifecycleCapability

	// ValueOrigin records whether a metadata field came from the request or the provider.
	// SandboxInfo.Origins maps a field path, such as InfoFieldId, to a ValueOrigin.
	ValueOrigin = sandboxv1.ValueOrigin
)

const (
	// SnapshotKindDefault leaves the snapshot class to the provider.
	SnapshotKindDefault = sandboxv1.SnapshotKind_SNAPSHOT_KIND_DEFAULT
	// SnapshotKindFilesystem is a filesystem snapshot.
	SnapshotKindFilesystem = sandboxv1.SnapshotKind_SNAPSHOT_KIND_FILESYSTEM
	// SnapshotKindMemory is a memory snapshot.
	SnapshotKindMemory = sandboxv1.SnapshotKind_SNAPSHOT_KIND_MEMORY
	// SnapshotKindImage is an image snapshot.
	SnapshotKindImage = sandboxv1.SnapshotKind_SNAPSHOT_KIND_IMAGE
	// SnapshotKindMachine is a machine snapshot.
	SnapshotKindMachine = sandboxv1.SnapshotKind_SNAPSHOT_KIND_MACHINE
	// SnapshotKindCheckpoint is a checkpoint snapshot.
	SnapshotKindCheckpoint = sandboxv1.SnapshotKind_SNAPSHOT_KIND_CHECKPOINT
	// SnapshotKindProviderSpecific is a provider-defined snapshot class.
	SnapshotKindProviderSpecific = sandboxv1.SnapshotKind_SNAPSHOT_KIND_PROVIDER_SPECIFIC

	// IsolationKindDefault leaves the sandbox class to the provider.
	IsolationKindDefault = sandboxv1.IsolationKind_ISOLATION_KIND_DEFAULT
	// IsolationKindContainer requests a container.
	IsolationKindContainer = sandboxv1.IsolationKind_ISOLATION_KIND_CONTAINER
	// IsolationKindMicroVM requests a microVM.
	IsolationKindMicroVM = sandboxv1.IsolationKind_ISOLATION_KIND_MICROVM
	// IsolationKindLinuxVM requests a Linux VM.
	IsolationKindLinuxVM = sandboxv1.IsolationKind_ISOLATION_KIND_LINUX_VM
	// IsolationKindWorker requests a worker.
	IsolationKindWorker = sandboxv1.IsolationKind_ISOLATION_KIND_WORKER
	// IsolationKindProviderSpecific requests a provider-defined sandbox class.
	IsolationKindProviderSpecific = sandboxv1.IsolationKind_ISOLATION_KIND_PROVIDER_SPECIFIC

	// EgressModeDefault leaves outbound access to the provider.
	EgressModeDefault = sandboxv1.EgressMode_EGRESS_MODE_DEFAULT
	// EgressModeAllowAll allows all outbound traffic.
	EgressModeAllowAll = sandboxv1.EgressMode_EGRESS_MODE_ALLOW_ALL
	// EgressModeBlockAll blocks outbound traffic.
	EgressModeBlockAll = sandboxv1.EgressMode_EGRESS_MODE_BLOCK_ALL
	// EgressModeRestricted allows only the entries in the allowlist.
	EgressModeRestricted = sandboxv1.EgressMode_EGRESS_MODE_RESTRICTED

	// PortTransportDefault leaves the port transport to the provider.
	PortTransportDefault = sandboxv1.PortTransport_PORT_TRANSPORT_DEFAULT
	// PortTransportTLS serves the port with TLS.
	PortTransportTLS = sandboxv1.PortTransport_PORT_TRANSPORT_TLS
	// PortTransportHTTP2TLS serves the port with HTTP/2 over TLS.
	PortTransportHTTP2TLS = sandboxv1.PortTransport_PORT_TRANSPORT_HTTP2_TLS
	// PortTransportPlain serves the port without TLS.
	PortTransportPlain = sandboxv1.PortTransport_PORT_TRANSPORT_PLAIN

	// SecretInjectionDefault leaves secret delivery to the provider.
	SecretInjectionDefault = sandboxv1.SecretInjection_SECRET_INJECTION_DEFAULT
	// SecretInjectionEnvironment delivers the secret as an environment variable.
	SecretInjectionEnvironment = sandboxv1.SecretInjection_SECRET_INJECTION_ENVIRONMENT
	// SecretInjectionEgressPlaceholder delivers the secret as an egress placeholder.
	SecretInjectionEgressPlaceholder = sandboxv1.SecretInjection_SECRET_INJECTION_EGRESS_PLACEHOLDER

	// PolicyModeDefault leaves the automatic action to the provider.
	PolicyModeDefault = sandboxv1.PolicyMode_POLICY_MODE_DEFAULT
	// PolicyModeDisabled turns the automatic action off.
	PolicyModeDisabled = sandboxv1.PolicyMode_POLICY_MODE_DISABLED
	// PolicyModeAfter runs the automatic action after its delay.
	PolicyModeAfter = sandboxv1.PolicyMode_POLICY_MODE_AFTER

	// WaitConditionDefault leaves the create wait condition to the provider.
	WaitConditionDefault = sandboxv1.WaitCondition_WAIT_CONDITION_DEFAULT
	// WaitConditionSubmitted returns when the provider has accepted the create request.
	WaitConditionSubmitted = sandboxv1.WaitCondition_WAIT_CONDITION_SUBMITTED
	// WaitConditionScheduled returns when the sandbox has been scheduled.
	WaitConditionScheduled = sandboxv1.WaitCondition_WAIT_CONDITION_SCHEDULED
	// WaitConditionStarted returns when the sandbox has started.
	WaitConditionStarted = sandboxv1.WaitCondition_WAIT_CONDITION_STARTED
	// WaitConditionReady returns when the readiness probe succeeds.
	WaitConditionReady = sandboxv1.WaitCondition_WAIT_CONDITION_READY

	// SandboxStateUnknown means the provider did not report a state.
	SandboxStateUnknown = sandboxv1.SandboxState_SANDBOX_STATE_UNKNOWN
	// SandboxStateQueued means the sandbox is waiting to start.
	SandboxStateQueued = sandboxv1.SandboxState_SANDBOX_STATE_QUEUED
	// SandboxStateStarting means the sandbox is starting.
	SandboxStateStarting = sandboxv1.SandboxState_SANDBOX_STATE_STARTING
	// SandboxStateRunning means the sandbox is running.
	SandboxStateRunning = sandboxv1.SandboxState_SANDBOX_STATE_RUNNING
	// SandboxStateStopping means the sandbox is stopping.
	SandboxStateStopping = sandboxv1.SandboxState_SANDBOX_STATE_STOPPING
	// SandboxStateStopped means the sandbox is stopped and can be started again on Daytona.
	SandboxStateStopped = sandboxv1.SandboxState_SANDBOX_STATE_STOPPED
	// SandboxStatePausing means the sandbox is pausing.
	SandboxStatePausing = sandboxv1.SandboxState_SANDBOX_STATE_PAUSING
	// SandboxStatePaused means the sandbox is paused.
	SandboxStatePaused = sandboxv1.SandboxState_SANDBOX_STATE_PAUSED
	// SandboxStateArchiving means the sandbox is archiving.
	SandboxStateArchiving = sandboxv1.SandboxState_SANDBOX_STATE_ARCHIVING
	// SandboxStateArchived means the sandbox is archived.
	SandboxStateArchived = sandboxv1.SandboxState_SANDBOX_STATE_ARCHIVED
	// SandboxStateTerminating means the sandbox is terminating.
	SandboxStateTerminating = sandboxv1.SandboxState_SANDBOX_STATE_TERMINATING
	// SandboxStateTerminated means the sandbox has been terminated.
	SandboxStateTerminated = sandboxv1.SandboxState_SANDBOX_STATE_TERMINATED
	// SandboxStateDeleting means the sandbox is being deleted.
	SandboxStateDeleting = sandboxv1.SandboxState_SANDBOX_STATE_DELETING
	// SandboxStateDeleted means the sandbox has been deleted.
	SandboxStateDeleted = sandboxv1.SandboxState_SANDBOX_STATE_DELETED
	// SandboxStateError means the provider reported a failed sandbox.
	SandboxStateError = sandboxv1.SandboxState_SANDBOX_STATE_ERROR

	// LifecycleCapabilityUnknown is not a supported operation.
	LifecycleCapabilityUnknown = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_UNKNOWN
	// LifecycleCapabilityRefresh is the capability checked before Refresh.
	LifecycleCapabilityRefresh = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_REFRESH
	// LifecycleCapabilityStart is the capability checked before Start.
	LifecycleCapabilityStart = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_START
	// LifecycleCapabilityStop is the capability checked before Stop.
	LifecycleCapabilityStop = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_STOP
	// LifecycleCapabilityPause is the capability checked before Pause.
	LifecycleCapabilityPause = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_PAUSE
	// LifecycleCapabilityTerminate is the capability checked before Terminate.
	LifecycleCapabilityTerminate = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_TERMINATE
	// LifecycleCapabilityArchive is the capability checked before Archive.
	LifecycleCapabilityArchive = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_ARCHIVE
	// LifecycleCapabilityDelete is the capability checked before Delete.
	LifecycleCapabilityDelete = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_DELETE

	// ValueOriginUnknown means the metadata does not say who set the field.
	ValueOriginUnknown = sandboxv1.ValueOrigin_VALUE_ORIGIN_UNKNOWN
	// ValueOriginProvider means the provider set the field.
	ValueOriginProvider = sandboxv1.ValueOrigin_VALUE_ORIGIN_PROVIDER
	// ValueOriginRequest means the create request set the field.
	ValueOriginRequest = sandboxv1.ValueOrigin_VALUE_ORIGIN_REQUEST
)

// Value pointers an optional field so the provider can tell a set value from an omitted one.
// Use it for explicit zero as well as nonzero values.
func Value[T any](value T) *T { return &value }

// validValueOrigin reports whether origin is a declared ValueOrigin.
func validValueOrigin(origin ValueOrigin) bool {
	return origin.Descriptor().Values().ByNumber(origin.Number()) != nil
}
