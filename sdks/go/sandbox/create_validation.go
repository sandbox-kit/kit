package sandbox

import (
	"math"
	"regexp"
	"strings"

	"google.golang.org/protobuf/types/known/durationpb"
)

// imageReferencePattern matches portable container image references.
var imageReferencePattern = regexp.MustCompile(`^(?:[A-Za-z0-9](?:[A-Za-z0-9.-]*[A-Za-z0-9])?(?::[0-9]+)?/)?[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*(?:/[a-z0-9]+(?:(?:[._]|__|-+)[a-z0-9]+)*)*(?::[A-Za-z0-9_][A-Za-z0-9_.-]{0,127})?(?:@sha256:[a-fA-F0-9]{64})?$`)

// ValidateCreateOptions rejects CreateOptions values the shared contract cannot
// represent. A nil request is treated as omitted configuration. Account quotas
// stay with the provider.
//
// Portable rules are handwritten here and inventoried in specs/validation.yaml.
// They can move onto protobuf (Protovalidate) later if licensing allows.
func ValidateCreateOptions(request *CreateOptions) error {
	if request == nil {
		return nil
	}
	if err := validateSandboxSource(request.Source, "source"); err != nil {
		return err
	}
	if err := validateRuntimeConfig(request.Runtime, "runtime"); err != nil {
		return err
	}
	if err := validateIsolationConfig(request.Isolation, "isolation"); err != nil {
		return err
	}
	if err := validateResources(request.Resources, "resources"); err != nil {
		return err
	}
	if err := validatePlacement(request.Placement, "placement"); err != nil {
		return err
	}
	if err := validateNetworkConfig(request.Network, "network"); err != nil {
		return err
	}
	if err := validateSecurityConfig(request.Security, "security"); err != nil {
		return err
	}
	if err := validateLifetimePolicy(request.Lifetime, "lifetime"); err != nil {
		return err
	}
	if err := validateProvisioningOptions(request.Provisioning, "provisioning"); err != nil {
		return err
	}
	if err := validateReadinessProbe(request.Readiness, "readiness"); err != nil {
		return err
	}
	return validateObservabilityConfig(request.Observability, "observability")
}

func validateSandboxSource(source *SandboxSource, path string) error {
	if source == nil {
		return nil
	}
	switch kind := source.Kind.(type) {
	case nil, *SandboxSource_ProviderDefault:
		return nil
	case *SandboxSource_Image:
		return validateImageSource(kind.Image, path+".image")
	case *SandboxSource_Snapshot:
		return validateSnapshotSource(kind.Snapshot, path+".snapshot")
	case *SandboxSource_WarmPool:
		return validateWarmPoolSource(kind.WarmPool, path+".warm_pool")
	default:
		return invalidCreate(path, "known_enum")
	}
}

func validateImageSource(source *ImageSource, path string) error {
	if source == nil {
		return invalidCreate(path, "required")
	}
	ref := source.GetReference()
	if strings.TrimSpace(ref) == "" {
		return invalidCreate(path+".reference", "nonblank")
	}
	if len(ref) > 512 || !imageReferencePattern.MatchString(ref) {
		return invalidCreate(path+".reference", "image_reference")
	}
	return nil
}

func validateSnapshotSource(source *SnapshotSource, path string) error {
	if source == nil {
		return invalidCreate(path, "required")
	}
	if strings.TrimSpace(source.GetReference()) == "" {
		return invalidCreate(path+".reference", "nonblank")
	}
	return knownEnum(path+".kind", int32(source.GetKind()), 0, 6)
}

func validateWarmPoolSource(source *WarmPoolSource, path string) error {
	if source == nil {
		return invalidCreate(path, "required")
	}
	if strings.TrimSpace(source.GetReference()) == "" {
		return invalidCreate(path+".reference", "nonblank")
	}
	return nil
}

func validateRuntimeConfig(config *RuntimeConfig, path string) error {
	if config == nil {
		return nil
	}
	if err := nonblankOptional(config.Language, path+".language"); err != nil {
		return err
	}
	if err := nonblankOptional(config.Version, path+".version"); err != nil {
		return err
	}
	if err := nonblankOptional(config.User, path+".user"); err != nil {
		return err
	}
	if err := nonblankOptional(config.WorkingDirectory, path+".working_directory"); err != nil {
		return err
	}
	return nonblankEach(config.Entrypoint, path+".entrypoint")
}

func validateIsolationConfig(config *IsolationConfig, path string) error {
	if config == nil {
		return nil
	}
	return knownEnum(path+".kind", int32(config.GetKind()), 0, 5)
}

func validateResources(resources *Resources, path string) error {
	if resources == nil {
		return nil
	}
	if resources.CpuCores != nil {
		value := resources.GetCpuCores()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return invalidCreate(path+".cpu_cores", "finite")
		}
		if !(value > 0) {
			return invalidCreate(path+".cpu_cores", "exclusive_minimum")
		}
	}
	if resources.CpuLimitCores != nil {
		value := resources.GetCpuLimitCores()
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return invalidCreate(path+".cpu_limit_cores", "finite")
		}
		if value < 0 {
			return invalidCreate(path+".cpu_limit_cores", "minimum")
		}
	}
	if resources.MemoryMib != nil && !(resources.GetMemoryMib() > 0) {
		return invalidCreate(path+".memory_mib", "exclusive_minimum")
	}
	if resources.DiskMib != nil && !(resources.GetDiskMib() > 0) {
		return invalidCreate(path+".disk_mib", "exclusive_minimum")
	}
	if resources.GetCpuLimitCores() > 0 && resources.GetCpuLimitCores() < resources.GetCpuCores() {
		return invalidCreate(path+".cpu_limit_cores", "limit")
	}
	if resources.GetMemoryLimitMib() > 0 && resources.GetMemoryLimitMib() < resources.GetMemoryMib() {
		return invalidCreate(path+".memory_limit_mib", "limit")
	}
	if resources.GetCpuLimitCores() > 0 && resources.GetCpuCores() == 0 {
		return invalidCreate(path+".cpu_limit_cores", "requires_positive")
	}
	if resources.GetMemoryLimitMib() > 0 && resources.GetMemoryMib() == 0 {
		return invalidCreate(path+".memory_limit_mib", "requires_positive")
	}
	return nil
}

func validatePlacement(placement *Placement, path string) error {
	if placement == nil {
		return nil
	}
	return nonblankEach(placement.Regions, path+".regions")
}

func validateNetworkConfig(network *NetworkConfig, path string) error {
	if network == nil {
		return nil
	}
	if err := knownEnum(path+".egress", int32(network.GetEgress()), 0, 3); err != nil {
		return err
	}
	if err := nonblankOptional(network.LinkedSandbox, path+".linked_sandbox"); err != nil {
		return err
	}
	if err := nonblankOptional(network.ProxyReference, path+".proxy_reference"); err != nil {
		return err
	}
	if network.GetEgress() == EgressModeRestricted && network.OutboundCidrs == nil && network.OutboundDomains == nil {
		return invalidCreate(path+".outbound_cidrs,outbound_domains", "requires_any")
	}
	if network.GetEgress() == EgressModeBlockAll {
		if network.OutboundCidrs != nil || network.OutboundDomains != nil || network.InboundCidrs != nil ||
			network.LinkedSandbox != nil || len(network.Ports) > 0 || network.GetPrivateNetwork() {
			return invalidCreate(path+".outbound_cidrs,outbound_domains,private_network,linked_sandbox,inbound_cidrs,ports", "forbids")
		}
	}
	if network.GetEgress() == EgressModeAllowAll && (network.OutboundCidrs != nil || network.OutboundDomains != nil) {
		return invalidCreate(path+".outbound_cidrs,outbound_domains", "forbids")
	}
	for _, port := range network.Ports {
		if port == nil {
			return invalidCreate(path+".ports", "required")
		}
		if err := validatePortBinding(port, path+".ports"); err != nil {
			return err
		}
	}
	return nil
}

func validatePortBinding(port *PortBinding, path string) error {
	if port.Port != 0 {
		if port.GetPort() < 1 {
			return invalidCreate(path+".port", "minimum")
		}
		if port.GetPort() > 65535 {
			return invalidCreate(path+".port", "maximum")
		}
	}
	return knownEnum(path+".transport", int32(port.GetTransport()), 0, 3)
}

func validateSecurityConfig(security *SecurityConfig, path string) error {
	if security == nil {
		return nil
	}
	for _, secret := range security.Secrets {
		if secret == nil {
			return invalidCreate(path+".secrets", "required")
		}
		if err := validateSecretReference(secret, path+".secrets"); err != nil {
			return err
		}
	}
	return nil
}

func validateSecretReference(secret *SecretReference, path string) error {
	if strings.TrimSpace(secret.GetReference()) == "" {
		return invalidCreate(path+".reference", "nonblank")
	}
	if err := nonblankOptional(secret.EnvironmentVariable, path+".environment_variable"); err != nil {
		return err
	}
	return knownEnum(path+".injection", int32(secret.GetInjection()), 0, 2)
}

func validateLifetimePolicy(policy *LifetimePolicy, path string) error {
	if policy == nil {
		return nil
	}
	if err := nonNegativeDuration(policy.MaximumLifetime, path+".maximum_lifetime"); err != nil {
		return err
	}
	if err := validateAutomaticAction(policy.IdleTerminate, path+".idle_terminate"); err != nil {
		return err
	}
	if err := validateAutomaticAction(policy.IdleStop, path+".idle_stop"); err != nil {
		return err
	}
	if err := validateAutomaticAction(policy.IdlePause, path+".idle_pause"); err != nil {
		return err
	}
	if err := validateAutomaticAction(policy.StoppedArchive, path+".stopped_archive"); err != nil {
		return err
	}
	return validateAutomaticAction(policy.StoppedDelete, path+".stopped_delete")
}

func validateAutomaticAction(action *AutomaticAction, path string) error {
	if action == nil {
		return nil
	}
	if err := knownEnum(path+".mode", int32(action.GetMode()), 0, 2); err != nil {
		return err
	}
	if err := nonNegativeDuration(action.After, path+".after"); err != nil {
		return err
	}
	if action.GetMode() == PolicyModeAfter && action.After == nil {
		return invalidCreate(path+".after", "requires_any")
	}
	if action.GetMode() != PolicyModeAfter && action.After != nil {
		return invalidCreate(path+".after", "forbids")
	}
	return nil
}

func validateProvisioningOptions(options *ProvisioningOptions, path string) error {
	if options == nil {
		return nil
	}
	if err := nonNegativeDuration(options.Timeout, path+".timeout"); err != nil {
		return err
	}
	if err := nonNegativeDuration(options.QueueTimeout, path+".queue_timeout"); err != nil {
		return err
	}
	return knownEnum(path+".wait_for", int32(options.GetWaitFor()), 0, 4)
}

func validateReadinessProbe(probe *ReadinessProbe, path string) error {
	if probe == nil {
		return nil
	}
	switch check := probe.Check.(type) {
	case *ReadinessProbe_TcpPort:
		if check.TcpPort < 1 {
			return invalidCreate(path+".tcp_port", "minimum")
		}
		if check.TcpPort > 65535 {
			return invalidCreate(path+".tcp_port", "maximum")
		}
		return nil
	case *ReadinessProbe_Command:
		return validateCommandProbe(check.Command, path+".command")
	default:
		return invalidCreate(path+".tcp_port,command", "exactly_one")
	}
}

func validateCommandProbe(probe *CommandProbe, path string) error {
	if probe == nil {
		return invalidCreate(path, "required")
	}
	if len(probe.Argv) < 1 {
		return invalidCreate(path+".argv", "min_items")
	}
	return nil
}

func validateObservabilityConfig(config *ObservabilityConfig, path string) error {
	if config == nil {
		return nil
	}
	return nonblankOptional(config.TelemetryEndpoint, path+".telemetry_endpoint")
}

func invalidCreate(field, rule string) error {
	return NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Field: field, Rule: Value(rule), Message: "sandbox-kit: invalid configuration"}, nil)
}

func nonblankOptional(value *string, field string) error {
	if value != nil && strings.TrimSpace(*value) == "" {
		return invalidCreate(field, "nonblank")
	}
	return nil
}

func nonblankEach(values []string, field string) error {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return invalidCreate(field, "nonblank")
		}
	}
	return nil
}

func nonNegativeDuration(value *durationpb.Duration, field string) error {
	if value != nil && value.AsDuration() < 0 {
		return invalidCreate(field, "minimum")
	}
	return nil
}

func knownEnum(field string, value, min, max int32) error {
	if value < min || value > max {
		return invalidCreate(field, "known_enum")
	}
	return nil
}
