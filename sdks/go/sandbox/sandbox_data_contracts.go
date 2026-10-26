package sandbox

import (
	"time"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"

	"google.golang.org/protobuf/types/known/durationpb"
)

// Portable sandbox data contracts come from standard protobuf generation.
// These aliases keep the SDK entrypoint concise while ensuring one representation.
type (
	CreateOptions = sandboxv1.CreateOptions
	SandboxSource = sandboxv1.SandboxSource
	// ProviderDefault is the provider-default variant of SandboxSource.Kind.
	ProviderDefault = sandboxv1.SandboxSource_ProviderDefault
	// Image is the image variant of SandboxSource.Kind.
	Image = sandboxv1.SandboxSource_Image
	// Snapshot is the snapshot variant of SandboxSource.Kind.
	Snapshot = sandboxv1.SandboxSource_Snapshot
	// WarmPool is the warm-pool variant of SandboxSource.Kind.
	WarmPool              = sandboxv1.SandboxSource_WarmPool
	ProviderDefaultSource = sandboxv1.ProviderDefaultSource
	ImageSource           = sandboxv1.ImageSource
	SnapshotSource        = sandboxv1.SnapshotSource
	WarmPoolSource        = sandboxv1.WarmPoolSource
	RuntimeConfig         = sandboxv1.RuntimeConfig
	IsolationConfig       = sandboxv1.IsolationConfig
	Resources             = sandboxv1.Resources
	Placement             = sandboxv1.Placement
	StringAllowlist       = sandboxv1.StringAllowlist
	NetworkConfig         = sandboxv1.NetworkConfig
	PortBinding           = sandboxv1.PortBinding
	SecurityConfig        = sandboxv1.SecurityConfig
	SecretReference       = sandboxv1.SecretReference
	LifetimePolicy        = sandboxv1.LifetimePolicy
	AutomaticAction       = sandboxv1.AutomaticAction
	ProvisioningOptions   = sandboxv1.ProvisioningOptions
	ReadinessProbe        = sandboxv1.ReadinessProbe
	// TCPPort is the TCP-port variant of ReadinessProbe.Check.
	TCPPort = sandboxv1.ReadinessProbe_TcpPort
	// Command is the command variant of ReadinessProbe.Check.
	Command                 = sandboxv1.ReadinessProbe_Command
	CommandProbe            = sandboxv1.CommandProbe
	ObservabilityConfig     = sandboxv1.ObservabilityConfig
	CreateResult            = sandboxv1.CreateResult
	StartOptions            = sandboxv1.StartOptions
	StopOptions             = sandboxv1.StopOptions
	PauseOptions            = sandboxv1.PauseOptions
	TerminateOptions        = sandboxv1.TerminateOptions
	ArchiveOptions          = sandboxv1.ArchiveOptions
	DeleteOptions           = sandboxv1.DeleteOptions
	ExitResult              = sandboxv1.ExitResult
	SandboxInfo             = sandboxv1.SandboxInfo
	MetadataValue           = sandboxv1.MetadataValue
	MetadataBoolean         = sandboxv1.MetadataValue_Boolean
	MetadataText            = sandboxv1.MetadataValue_Text
	MetadataSignedInteger   = sandboxv1.MetadataValue_SignedInteger
	MetadataUnsignedInteger = sandboxv1.MetadataValue_UnsignedInteger
	MetadataNumber          = sandboxv1.MetadataValue_Number
	MetadataBinary          = sandboxv1.MetadataValue_Binary
	MetadataListValue       = sandboxv1.MetadataValue_List
	MetadataObjectValue     = sandboxv1.MetadataValue_Object
	MetadataList            = sandboxv1.MetadataList
	MetadataObject          = sandboxv1.MetadataObject
	SnapshotKind            = sandboxv1.SnapshotKind
	IsolationKind           = sandboxv1.IsolationKind
	EgressMode              = sandboxv1.EgressMode
	PortTransport           = sandboxv1.PortTransport
	SecretInjection         = sandboxv1.SecretInjection
	PolicyMode              = sandboxv1.PolicyMode
	WaitCondition           = sandboxv1.WaitCondition
	SandboxState            = sandboxv1.SandboxState
	LifecycleCapability     = sandboxv1.LifecycleCapability
	ValueOrigin             = sandboxv1.ValueOrigin
)

const (
	SnapshotKindDefault          = sandboxv1.SnapshotKind_SNAPSHOT_KIND_DEFAULT
	SnapshotKindFilesystem       = sandboxv1.SnapshotKind_SNAPSHOT_KIND_FILESYSTEM
	SnapshotKindMemory           = sandboxv1.SnapshotKind_SNAPSHOT_KIND_MEMORY
	SnapshotKindImage            = sandboxv1.SnapshotKind_SNAPSHOT_KIND_IMAGE
	SnapshotKindMachine          = sandboxv1.SnapshotKind_SNAPSHOT_KIND_MACHINE
	SnapshotKindCheckpoint       = sandboxv1.SnapshotKind_SNAPSHOT_KIND_CHECKPOINT
	SnapshotKindProviderSpecific = sandboxv1.SnapshotKind_SNAPSHOT_KIND_PROVIDER_SPECIFIC

	IsolationKindDefault          = sandboxv1.IsolationKind_ISOLATION_KIND_DEFAULT
	IsolationKindContainer        = sandboxv1.IsolationKind_ISOLATION_KIND_CONTAINER
	IsolationKindMicroVM          = sandboxv1.IsolationKind_ISOLATION_KIND_MICROVM
	IsolationKindLinuxVM          = sandboxv1.IsolationKind_ISOLATION_KIND_LINUX_VM
	IsolationKindWorker           = sandboxv1.IsolationKind_ISOLATION_KIND_WORKER
	IsolationKindProviderSpecific = sandboxv1.IsolationKind_ISOLATION_KIND_PROVIDER_SPECIFIC

	EgressModeDefault    = sandboxv1.EgressMode_EGRESS_MODE_DEFAULT
	EgressModeAllowAll   = sandboxv1.EgressMode_EGRESS_MODE_ALLOW_ALL
	EgressModeBlockAll   = sandboxv1.EgressMode_EGRESS_MODE_BLOCK_ALL
	EgressModeRestricted = sandboxv1.EgressMode_EGRESS_MODE_RESTRICTED

	PortTransportDefault  = sandboxv1.PortTransport_PORT_TRANSPORT_DEFAULT
	PortTransportTLS      = sandboxv1.PortTransport_PORT_TRANSPORT_TLS
	PortTransportHTTP2TLS = sandboxv1.PortTransport_PORT_TRANSPORT_HTTP2_TLS
	PortTransportPlain    = sandboxv1.PortTransport_PORT_TRANSPORT_PLAIN

	SecretInjectionDefault           = sandboxv1.SecretInjection_SECRET_INJECTION_DEFAULT
	SecretInjectionEnvironment       = sandboxv1.SecretInjection_SECRET_INJECTION_ENVIRONMENT
	SecretInjectionEgressPlaceholder = sandboxv1.SecretInjection_SECRET_INJECTION_EGRESS_PLACEHOLDER

	PolicyModeDefault  = sandboxv1.PolicyMode_POLICY_MODE_DEFAULT
	PolicyModeDisabled = sandboxv1.PolicyMode_POLICY_MODE_DISABLED
	PolicyModeAfter    = sandboxv1.PolicyMode_POLICY_MODE_AFTER

	WaitConditionDefault   = sandboxv1.WaitCondition_WAIT_CONDITION_DEFAULT
	WaitConditionSubmitted = sandboxv1.WaitCondition_WAIT_CONDITION_SUBMITTED
	WaitConditionScheduled = sandboxv1.WaitCondition_WAIT_CONDITION_SCHEDULED
	WaitConditionStarted   = sandboxv1.WaitCondition_WAIT_CONDITION_STARTED
	WaitConditionReady     = sandboxv1.WaitCondition_WAIT_CONDITION_READY

	SandboxStateUnknown     = sandboxv1.SandboxState_SANDBOX_STATE_UNKNOWN
	SandboxStateQueued      = sandboxv1.SandboxState_SANDBOX_STATE_QUEUED
	SandboxStateStarting    = sandboxv1.SandboxState_SANDBOX_STATE_STARTING
	SandboxStateRunning     = sandboxv1.SandboxState_SANDBOX_STATE_RUNNING
	SandboxStateStopping    = sandboxv1.SandboxState_SANDBOX_STATE_STOPPING
	SandboxStateStopped     = sandboxv1.SandboxState_SANDBOX_STATE_STOPPED
	SandboxStatePausing     = sandboxv1.SandboxState_SANDBOX_STATE_PAUSING
	SandboxStatePaused      = sandboxv1.SandboxState_SANDBOX_STATE_PAUSED
	SandboxStateArchiving   = sandboxv1.SandboxState_SANDBOX_STATE_ARCHIVING
	SandboxStateArchived    = sandboxv1.SandboxState_SANDBOX_STATE_ARCHIVED
	SandboxStateTerminating = sandboxv1.SandboxState_SANDBOX_STATE_TERMINATING
	SandboxStateTerminated  = sandboxv1.SandboxState_SANDBOX_STATE_TERMINATED
	SandboxStateDeleting    = sandboxv1.SandboxState_SANDBOX_STATE_DELETING
	SandboxStateDeleted     = sandboxv1.SandboxState_SANDBOX_STATE_DELETED
	SandboxStateError       = sandboxv1.SandboxState_SANDBOX_STATE_ERROR

	LifecycleCapabilityUnknown   = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_UNKNOWN
	LifecycleCapabilityRefresh   = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_REFRESH
	LifecycleCapabilityStart     = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_START
	LifecycleCapabilityStop      = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_STOP
	LifecycleCapabilityPause     = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_PAUSE
	LifecycleCapabilityTerminate = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_TERMINATE
	LifecycleCapabilityArchive   = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_ARCHIVE
	LifecycleCapabilityDelete    = sandboxv1.LifecycleCapability_LIFECYCLE_CAPABILITY_DELETE

	ValueOriginUnknown  = sandboxv1.ValueOrigin_VALUE_ORIGIN_UNKNOWN
	ValueOriginProvider = sandboxv1.ValueOrigin_VALUE_ORIGIN_PROVIDER
	ValueOriginRequest  = sandboxv1.ValueOrigin_VALUE_ORIGIN_REQUEST
)

// Value returns a pointer to value for optional protobuf scalar fields.
func Value[T any](value T) *T { return &value }

// validValueOrigin reports whether origin is a declared ValueOrigin.
func validValueOrigin(origin ValueOrigin) bool {
	return origin.Descriptor().Values().ByNumber(origin.Number()) != nil
}

// OptionalDuration converts a protobuf duration to a time.Duration pointer.
func OptionalDuration(value *durationpb.Duration) *time.Duration {
	if value == nil {
		return nil
	}
	converted := value.AsDuration()
	return &converted
}
