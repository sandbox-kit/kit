// This file turns a shared creation request into a Daytona SDK create call.
// Generated mappings cover field and policy conversion. This file chooses the
// snapshot or image source and applies the mapped options.
//
// Example:
//
//	client, err := sandbox.NewClient(sandbox.Config{Provider: daytona.New()})
//	instance, err := client.Create(ctx, nil)

package daytona

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/daytona/clients/sdk-go/pkg/options"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// createPlan is the Daytona source and option list for one creation request.
type createPlan struct {
	// params is the snapshot or image argument accepted by the SDK create method.
	params any
	// options are the mapped creation settings, including policies and resources.
	options []func(*options.CreateSandbox)
}

// create calls the Daytona SDK and returns shared metadata.
// A nil request uses the default snapshot and provider defaults.
func (a *backend) create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	plan, err := planCreate(request)
	if err != nil {
		return nil, err
	}
	instance, err := a.client.Create(ctx, plan.params, plan.options...)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "daytona",
			Operation: "create",
			Field:     sandbox.InfoFieldID,
			Message:   "sandbox-kit daytona: SDK returned no sandbox",
		}, nil)
	}
	info, err := mapSandboxInfo(instance)
	if err != nil {
		return nil, err
	}
	resources, err := mapAllocatedResources(instance)
	if err != nil {
		return nil, err
	}
	if resources.CPUCores != nil {
		info.Resources = &resources
		info.Origins["resources"] = sandbox.ValueOriginProvider
	}
	return &sandbox.CreateResult{Sandbox: &info}, nil
}

// planCreate checks the request and builds SDK arguments without contacting Daytona.
func planCreate(request *sandbox.CreateOptions) (createPlan, error) {
	var plan createPlan
	if err := validateProviderCreate(request); err != nil {
		return plan, err
	}
	if err := sandbox.ValidateCreateOptions(request); err != nil {
		return plan, err
	}
	base, rest, err := mapCommonCreateFields(request)
	if err != nil {
		return plan, err
	}
	if request.Runtime != nil {
		config := *request.Runtime
		base.User = config.GetUser()
		base.Language = types.CodeLanguage(config.GetLanguage())
		config.User, config.Language = nil, nil
		if err := sandbox.RejectUnmapped("daytona runtime", &config); err != nil {
			return plan, err
		}
		rest.Runtime = nil
	}
	if request.Isolation != nil {
		config := *request.Isolation
		base.Kvm = config.GetNestedVirtualization()
		config.NestedVirtualization = nil
		if err := sandbox.RejectUnmapped("daytona isolation (class comes from snapshot)", &config); err != nil {
			return plan, err
		}
		rest.Isolation = nil
	}
	if request.Security != nil {
		config := *request.Security
		base.Public = config.GetPublicAccess()
		config.PublicAccess = nil
		if len(config.Secrets) > 0 {
			base.Secrets = map[string]string{}
		}
		for _, secret := range config.Secrets {
			if secret.GetEnvironmentVariable() == "" {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldSecuritySecretsEnvironmentVariable,
					Message:   "sandbox-kit daytona: secret environment variable is required",
				}, nil)
			}
			if secret.Injection != sandbox.SecretInjectionDefault && secret.Injection != sandbox.SecretInjectionEgressPlaceholder {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldSecuritySecretsInjection,
					Message:   "sandbox-kit daytona: secret injection must use egress placeholders",
				}, nil)
			}
			if _, exists := base.Secrets[secret.GetEnvironmentVariable()]; exists {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldSecuritySecretsEnvironmentVariable,
					Message:   "sandbox-kit daytona: duplicate secret environment variable",
				}, nil)
			}
			base.Secrets[secret.GetEnvironmentVariable()] = secret.Reference
		}
		config.Secrets = nil
		if err := sandbox.RejectUnmapped("daytona security", &config); err != nil {
			return plan, err
		}
		rest.Security = nil
	}
	if request.Network != nil {
		config := *request.Network
		switch config.Egress {
		case sandbox.EgressModeDefault, sandbox.EgressModeAllowAll:
		case sandbox.EgressModeBlockAll:
			base.NetworkBlockAll = true
		case sandbox.EgressModeRestricted:
		default:
			return plan, sandbox.NewError(sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindInvalidArgument,
				Provider:  "daytona",
				Operation: "create",
				Field:     sandbox.CreateFieldNetworkEgress,
				Message:   "sandbox-kit daytona: unknown egress mode",
			}, nil)
		}
		config.Egress = 0
		if config.OutboundCIDRs != nil {
			if len(config.OutboundCIDRs.Entries) == 0 {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldNetworkOutboundCIDRs,
					Message:   "sandbox-kit daytona: empty CIDR allowlist mapping is not verified",
				}, nil)
			}
			value := strings.Join(config.OutboundCIDRs.Entries, ",")
			base.NetworkAllowList = &value
		}
		config.OutboundCIDRs = nil
		if config.OutboundDomains != nil {
			if len(config.OutboundDomains.Entries) == 0 {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldNetworkOutboundDomains,
					Message:   "sandbox-kit daytona: empty domain allowlist mapping is not verified",
				}, nil)
			}
			value := strings.Join(config.OutboundDomains.Entries, ",")
			base.DomainAllowList = &value
		}
		config.OutboundDomains = nil
		base.LinkedSandbox = config.GetLinkedSandbox()
		base.OutboundProxyUrl = config.OutboundProxyURL
		config.LinkedSandbox, config.OutboundProxyURL = nil, nil
		if err := sandbox.RejectUnmapped("daytona network", &config); err != nil {
			return plan, err
		}
		rest.Network = nil
	}
	if request.Lifetime != nil {
		config := *request.Lifetime
		base.Ephemeral = config.GetEphemeral()
		config.Ephemeral = nil
		if config.MaximumLifetime != nil {
			minutes, err := durationMinutes(*config.MaximumLifetime, true)
			if err != nil {
				return plan, err
			}
			base.TtlMinutes = &minutes
		}
		config.MaximumLifetime = nil
		mapped, remaining, err := mapLifetimePolicies(&config)
		if err != nil {
			return plan, err
		}
		base.AutoStopInterval, base.AutoPauseInterval = mapped.AutoStopInterval, mapped.AutoPauseInterval
		base.AutoArchiveInterval, base.AutoDeleteInterval = mapped.AutoArchiveInterval, mapped.AutoDeleteInterval
		config = remaining
		if err := sandbox.RejectUnmapped("daytona lifetime", &config); err != nil {
			return plan, err
		}
		rest.Lifetime = nil
	}
	if request.Provisioning != nil {
		config := *request.Provisioning
		if config.Timeout != nil && *config.Timeout > 0 {
			plan.options = append(plan.options, options.WithTimeout(*config.Timeout))
		}
		config.Timeout = nil
		if config.QueueTimeout != nil {
			minutes, err := durationMinutes(*config.QueueTimeout, false)
			if err != nil {
				return plan, err
			}
			base.QueueTimeout = &minutes
		}
		config.QueueTimeout = nil
		switch config.WaitFor {
		case sandbox.WaitConditionDefault:
		case sandbox.WaitConditionSubmitted:
			plan.options = append(plan.options, options.WithWaitForStart(false))
		case sandbox.WaitConditionStarted:
			plan.options = append(plan.options, options.WithWaitForStart(true))
		default:
			return plan, sandbox.NewError(sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindUnsupported,
				Provider:  "daytona",
				Operation: "create",
				Field:     sandbox.CreateFieldProvisioningWaitFor,
				Message:   "sandbox-kit daytona: requested wait condition is unsupported",
			}, nil)
		}
		config.WaitFor = 0
		if err := sandbox.RejectUnmapped("daytona creation", &config); err != nil {
			return plan, err
		}
		rest.Provisioning = nil
	}
	if request.Observability != nil {
		config := *request.Observability
		base.OtelEndpointOverride = config.TelemetryEndpoint
		config.TelemetryEndpoint = nil
		if err := sandbox.RejectUnmapped("daytona observability", &config); err != nil {
			return plan, err
		}
		rest.Observability = nil
	}
	source := request.GetSource()
	if source.GetWarmPool() != nil {
		return plan, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindUnsupported,
			Provider:  "daytona",
			Operation: "create",
			Field:     sandbox.CreateFieldSourceWarmPool,
			Message:   "sandbox-kit daytona: explicit warm-pool creation mapping is unavailable",
		}, nil)
	}
	if source.GetImage() != nil {
		params := types.ImageParams{SandboxBaseParams: base, Image: source.Image.Reference}
		if request.Resources != nil {
			resources, remaining, err := mapResources(request.Resources)
			if err != nil {
				return plan, err
			}
			if err := sandbox.RejectUnmapped("daytona resources", &remaining); err != nil {
				return plan, err
			}
			params.Resources = &resources
			rest.Resources = nil
		}
		plan.params = params
	} else {
		params := types.SnapshotParams{SandboxBaseParams: base}
		if snapshot := source.GetSnapshot(); snapshot != nil {
			if snapshot.Kind != sandbox.SnapshotKindDefault {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldSourceSnapshotKind,
					Message:   "sandbox-kit daytona: explicit snapshot-kind restore mapping is unavailable",
				}, nil)
			}
			if snapshot.ProviderKind != nil {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldSourceSnapshotProviderKind,
					Message:   "sandbox-kit daytona: provider snapshot kind is unsupported",
				}, nil)
			}
			params.Snapshot = snapshot.Reference
		}
		plan.params = params
	}
	rest.Source = nil
	if err := sandbox.RejectUnmapped("daytona", &rest); err != nil {
		return plan, err
	}
	return plan, nil
}

// durationMinutes converts a policy delay to whole minutes.
// zeroAllowed permits an explicit zero. Daytona uses that for
// immediate delete, or for disabled stop and pause.
func durationMinutes(value time.Duration, zeroAllowed bool) (int, error) {
	if value < 0 || value%time.Minute != 0 || value/time.Minute > math.MaxInt32 || (!zeroAllowed && value == 0) {
		return 0, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "daytona",
			Operation: "create",
			Field:     sandbox.CreateFieldLifetime,
			Message:   "sandbox-kit daytona: policy/queue durations require representable whole minutes",
		}, nil)
	}
	return int(value / time.Minute), nil
}
