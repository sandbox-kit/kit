package daytona

import (
	"math"
	"strings"
	"time"

	"github.com/daytona/clients/sdk-go/pkg/options"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// createPlan is the Daytona SDK argument list for one creation request.
type createPlan struct {
	// params is the snapshot or image argument accepted by the SDK create method.
	params any
	// options are the mapped creation settings, including policies and resources.
	options []func(*options.CreateSandbox)
}

// creation translates one CreateOptions value into native Daytona arguments.
// Each method consumes one section and clears it on the remainder.
type creation struct {
	request *sandbox.CreateOptions
	rest    *sandbox.CreateOptions
	base    types.SandboxBaseParams
	plan    createPlan
}

// planCreate checks the request and builds SDK arguments without contacting Daytona.
func planCreate(request *sandbox.CreateOptions) (createPlan, error) {
	if err := validateProviderCreate(request); err != nil {
		return createPlan{}, err
	}
	base, rest, err := mapCommonCreateFields(request)
	if err != nil {
		return createPlan{}, err
	}
	planned := &creation{request: request, rest: rest, base: base}
	for _, apply := range []func() error{
		planned.runtime,
		planned.isolation,
		planned.security,
		planned.network,
		planned.lifetime,
		planned.provisioning,
		planned.observability,
		planned.source,
	} {
		if err := apply(); err != nil {
			return createPlan{}, err
		}
	}
	if err := sdk.RejectUnmapped("daytona", "", planned.rest); err != nil {
		return createPlan{}, err
	}
	return planned.plan, nil
}

func (c *creation) runtime() error {
	if c.request.Runtime == nil {
		return nil
	}
	config := cloneMessage(c.request.Runtime)
	c.base.User = config.GetUser()
	c.base.Language = types.CodeLanguage(config.GetLanguage())
	config.User, config.Language = nil, nil
	if config.Version != nil || config.WorkingDirectory != nil || config.Entrypoint != nil || config.Pty != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldRuntime, "runtime version, working directory, entrypoint, and PTY are not mapped by Daytona")
	}
	if err := sdk.RejectUnmapped("daytona runtime", sandbox.CreateFieldRuntime, config); err != nil {
		return err
	}
	c.rest.Runtime = nil
	return nil
}

func (c *creation) isolation() error {
	if c.request.Isolation == nil {
		return nil
	}
	config := cloneMessage(c.request.Isolation)
	if config.Kind != sandbox.IsolationKindDefault || config.ProviderKind != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldIsolation, "Daytona isolation class comes from the selected source")
	}
	c.base.Kvm = config.GetNestedVirtualization()
	config.Kind, config.ProviderKind, config.NestedVirtualization = 0, nil, nil
	if err := sdk.RejectUnmapped("daytona isolation", sandbox.CreateFieldIsolation, config); err != nil {
		return err
	}
	c.rest.Isolation = nil
	return nil
}

func (c *creation) security() error {
	if c.request.Security == nil {
		return nil
	}
	config := cloneMessage(c.request.Security)
	c.base.Public = config.GetPublicAccess()
	config.PublicAccess = nil
	if len(config.Secrets) > 0 {
		c.base.Secrets = map[string]string{}
	}
	for _, secret := range config.Secrets {
		if secret.GetEnvironmentVariable() == "" {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldSecuritySecretsEnvironmentVariable, "secret environment variable is required")
		}
		if secret.Injection != sandbox.SecretInjectionDefault && secret.Injection != sandbox.SecretInjectionEgressPlaceholder {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSecuritySecretsInjection, "secret injection must use egress placeholders")
		}
		if _, exists := c.base.Secrets[secret.GetEnvironmentVariable()]; exists {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldSecuritySecretsEnvironmentVariable, "duplicate secret environment variable")
		}
		c.base.Secrets[secret.GetEnvironmentVariable()] = secret.Reference
	}
	config.Secrets = nil
	if config.WorkloadIdentity != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSecurityWorkloadIdentity, "Daytona workload identity is not mapped")
	}
	if err := sdk.RejectUnmapped("daytona security", sandbox.CreateFieldSecurity, config); err != nil {
		return err
	}
	c.rest.Security = nil
	return nil
}

func (c *creation) network() error {
	if c.request.Network == nil {
		return nil
	}
	config := cloneMessage(c.request.Network)
	if config.LinkedSandbox != nil && (c.request.Lifetime == nil || !c.request.Lifetime.GetEphemeral()) {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkLinkedSandbox, "linked sandboxes must be ephemeral")
	}
	switch config.Egress {
	case sandbox.EgressModeDefault, sandbox.EgressModeAllowAll:
	case sandbox.EgressModeBlockAll:
		c.base.NetworkBlockAll = true
	case sandbox.EgressModeRestricted:
	default:
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkEgress, "unknown egress mode")
	}
	config.Egress = 0
	if config.OutboundCidrs != nil {
		if len(config.OutboundCidrs.Entries) == 0 {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldNetworkOutboundCidrs, "empty CIDR allowlist mapping is not verified")
		}
		value := strings.Join(config.OutboundCidrs.Entries, ",")
		c.base.NetworkAllowList = &value
	}
	config.OutboundCidrs = nil
	if config.OutboundDomains != nil {
		if len(config.OutboundDomains.Entries) == 0 {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldNetworkOutboundDomains, "empty domain allowlist mapping is not verified")
		}
		value := strings.Join(config.OutboundDomains.Entries, ",")
		c.base.DomainAllowList = &value
	}
	config.OutboundDomains = nil
	c.base.LinkedSandbox = config.GetLinkedSandbox()
	c.base.OutboundProxyUrl = config.OutboundProxyUrl
	config.LinkedSandbox, config.OutboundProxyUrl = nil, nil
	if config.InboundCidrs != nil || config.PrivateNetwork != nil || config.ProxyReference != nil || config.CustomDomain != nil || config.Ports != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldNetwork, "one or more Daytona network fields are not mapped")
	}
	if err := sdk.RejectUnmapped("daytona network", sandbox.CreateFieldNetwork, config); err != nil {
		return err
	}
	c.rest.Network = nil
	return nil
}

func (c *creation) lifetime() error {
	if c.request.Lifetime == nil {
		return nil
	}
	config := cloneMessage(c.request.Lifetime)
	c.base.Ephemeral = config.GetEphemeral()
	config.Ephemeral = nil
	if config.MaximumLifetime != nil {
		minutes, err := durationMinutesAllowZero(config.MaximumLifetime.AsDuration(), sandbox.CreateFieldLifetimeMaximumLifetime)
		if err != nil {
			return err
		}
		c.base.TtlMinutes = &minutes
	}
	config.MaximumLifetime = nil
	mapped, remaining, err := mapLifetimePolicies(config)
	if err != nil {
		return err
	}
	c.base.AutoStopInterval, c.base.AutoPauseInterval = mapped.AutoStopInterval, mapped.AutoPauseInterval
	c.base.AutoArchiveInterval, c.base.AutoDeleteInterval = mapped.AutoArchiveInterval, mapped.AutoDeleteInterval
	config = remaining
	if config.IdleTerminate != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldLifetimeIdleTerminate, "Daytona idle termination is not mapped")
	}
	if err := sdk.RejectUnmapped("daytona lifetime", sandbox.CreateFieldLifetime, config); err != nil {
		return err
	}
	c.rest.Lifetime = nil
	return nil
}

func (c *creation) provisioning() error {
	if c.request.Provisioning == nil {
		return nil
	}
	config := cloneMessage(c.request.Provisioning)
	if timeout := sdk.OptionalDuration(config.Timeout); timeout != nil && *timeout > 0 {
		c.plan.options = append(c.plan.options, options.WithTimeout(*timeout))
	}
	config.Timeout = nil
	if config.QueueTimeout != nil {
		minutes, err := durationMinutesRequirePositive(config.QueueTimeout.AsDuration(), sandbox.CreateFieldProvisioningQueueTimeout)
		if err != nil {
			return err
		}
		c.base.QueueTimeout = &minutes
	}
	config.QueueTimeout = nil
	switch config.WaitFor {
	case sandbox.WaitConditionDefault:
	case sandbox.WaitConditionSubmitted:
		c.plan.options = append(c.plan.options, options.WithWaitForStart(false))
	case sandbox.WaitConditionStarted:
		c.plan.options = append(c.plan.options, options.WithWaitForStart(true))
	default:
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldProvisioningWaitFor, "requested wait condition is unsupported")
	}
	config.WaitFor = 0
	if err := sdk.RejectUnmapped("daytona creation", sandbox.CreateFieldProvisioning, config); err != nil {
		return err
	}
	c.rest.Provisioning = nil
	return nil
}

func (c *creation) observability() error {
	if c.request.Observability == nil {
		return nil
	}
	config := cloneMessage(c.request.Observability)
	c.base.OtelEndpointOverride = config.TelemetryEndpoint
	config.TelemetryEndpoint = nil
	if config.Verbose != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldObservabilityVerbose, "Daytona verbose creation is not mapped")
	}
	if err := sdk.RejectUnmapped("daytona observability", sandbox.CreateFieldObservability, config); err != nil {
		return err
	}
	c.rest.Observability = nil
	return nil
}

func (c *creation) source() error {
	source := c.request.GetSource()
	if source.GetWarmPool() != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSourceWarmPool, "explicit warm-pool creation mapping is unavailable")
	}
	if source.GetImage() != nil {
		params := types.ImageParams{SandboxBaseParams: c.base, Image: source.GetImage().Reference}
		if c.request.Resources != nil {
			resources, remaining, err := mapResources(c.request.Resources)
			if err != nil {
				return err
			}
			if err := sdk.RejectUnmapped("daytona resources", sandbox.CreateFieldResources, remaining); err != nil {
				return err
			}
			params.Resources = &resources
			c.rest.Resources = nil
		}
		c.plan.params = params
	} else {
		params := types.SnapshotParams{SandboxBaseParams: c.base}
		if snapshot := source.GetSnapshot(); snapshot != nil {
			if snapshot.Kind != sandbox.SnapshotKindDefault {
				return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSourceSnapshotKind, "explicit snapshot-kind restore mapping is unavailable")
			}
			if snapshot.ProviderKind != nil {
				return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSourceSnapshotProviderKind, "provider snapshot kind is unsupported")
			}
			params.Snapshot = snapshot.Reference
		}
		c.plan.params = params
	}
	c.rest.Source = nil
	return nil
}

// durationMinutesAllowZero converts a delay to whole minutes; zero is valid
// (Daytona uses it for immediate delete or disabled stop/pause).
func durationMinutesAllowZero(value time.Duration, field string) (int, error) {
	return wholeMinutes(value, zeroMinutesPermitted, field)
}

// durationMinutesRequirePositive converts a delay to whole minutes; zero is rejected.
func durationMinutesRequirePositive(value time.Duration, field string) (int, error) {
	return wholeMinutes(value, zeroMinutesRejected, field)
}

type minuteZeroPolicy int

const (
	zeroMinutesRejected minuteZeroPolicy = iota
	zeroMinutesPermitted
)

func wholeMinutes(value time.Duration, zero minuteZeroPolicy, field string) (int, error) {
	if value < 0 || value%time.Minute != 0 || value/time.Minute > math.MaxInt32 || (zero == zeroMinutesRejected && value == 0) {
		return 0, providerCheckError(sandbox.ErrorKindInvalidArgument, field, "policy/queue durations require representable whole minutes")
	}
	return int(value / time.Minute), nil
}
