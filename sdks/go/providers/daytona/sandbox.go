package daytona

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/daytona/clients/sdk-go/pkg/options"
	"github.com/daytona/clients/sdk-go/pkg/types"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type createPlan struct {
	params  any
	options []func(*options.CreateSandbox)
}

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
		return nil, fmt.Errorf("sandbox-kit daytona: SDK returned no sandbox")
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
	}
	return &sandbox.CreateResult{Sandbox: &info}, nil
}

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
				return plan, fmt.Errorf("sandbox-kit daytona: secret environment variable is required")
			}
			if secret.Injection != sandbox.SecretInjectionDefault && secret.Injection != sandbox.SecretInjectionEgressPlaceholder {
				return plan, fmt.Errorf("sandbox-kit daytona: secret injection must use egress placeholders")
			}
			if _, exists := base.Secrets[secret.GetEnvironmentVariable()]; exists {
				return plan, fmt.Errorf("sandbox-kit daytona: duplicate secret environment variable")
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
			return plan, fmt.Errorf("sandbox-kit daytona: unknown egress mode")
		}
		config.Egress = 0
		if config.OutboundCIDRs != nil {
			if len(config.OutboundCIDRs.Entries) == 0 {
				return plan, fmt.Errorf("sandbox-kit daytona: empty CIDR allowlist mapping is not verified")
			}
			value := strings.Join(config.OutboundCIDRs.Entries, ",")
			base.NetworkAllowList = &value
		}
		config.OutboundCIDRs = nil
		if config.OutboundDomains != nil {
			if len(config.OutboundDomains.Entries) == 0 {
				return plan, fmt.Errorf("sandbox-kit daytona: empty domain allowlist mapping is not verified")
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
			return plan, fmt.Errorf("sandbox-kit daytona: requested wait condition is unsupported")
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
		return plan, fmt.Errorf("sandbox-kit daytona: explicit warm-pool creation mapping is unavailable")
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
				return plan, fmt.Errorf("sandbox-kit daytona: explicit snapshot-kind restore mapping is unavailable")
			}
			if snapshot.ProviderKind != nil {
				return plan, fmt.Errorf("sandbox-kit daytona: provider snapshot kind is unsupported")
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
func durationMinutes(value time.Duration, zeroAllowed bool) (int, error) {
	if value < 0 || value%time.Minute != 0 || value/time.Minute > math.MaxInt32 || (!zeroAllowed && value == 0) {
		return 0, fmt.Errorf("sandbox-kit daytona: policy/queue durations require representable whole minutes")
	}
	return int(value / time.Minute), nil
}
