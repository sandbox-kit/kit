package modal

import (
	"context"
	"time"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

type createPlan struct {
	app          string
	environment  string
	image        string
	params       sdk.SandboxCreateParams
	ready        bool
	readyTimeout time.Duration
	secrets      []*sandbox.SecretReference
}

func (a *backend) create(ctx context.Context, request *sandbox.CreateOptions) (*sandbox.CreateResult, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	region := ""
	if a.region != nil {
		region = *a.region
	}
	plan, err := planCreate(request, a.scope, region)
	if err != nil {
		return nil, err
	}
	app, err := a.client.Apps.FromName(ctx, plan.app, &sdk.AppFromNameParams{Environment: plan.environment})
	if err != nil {
		return nil, err
	}
	for _, ref := range plan.secrets {
		secret, err := a.client.Secrets.FromName(ctx, ref.Reference, &sdk.SecretFromNameParams{Environment: plan.environment})
		if err != nil {
			return nil, err
		}
		plan.params.Secrets = append(plan.params.Secrets, secret)
	}
	image := a.client.Images.FromRegistry(plan.image, nil)
	instance, err := a.client.Sandboxes.Create(ctx, app, image, &plan.params)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidResponse,
			Provider:  "modal",
			Operation: "create",
			Field:     sandbox.InfoFieldID,
			Message:   "sandbox-kit modal: SDK returned no sandbox",
		}, nil)
	}
	if plan.ready {
		if err := instance.WaitUntilReady(ctx, plan.readyTimeout, nil); err != nil {
			return nil, err
		}
	}
	info, err := mapSandboxIdentity(instance)
	if err != nil {
		return nil, err
	}
	info.Name = request.Name
	info.Labels = plan.params.Tags
	if info.Name != nil {
		info.Origins["name"] = sandbox.ValueOriginRequest
	}
	if info.Labels != nil {
		info.Origins["labels"] = sandbox.ValueOriginRequest
	}
	return &sandbox.CreateResult{Sandbox: &info}, nil
}

func planCreate(request *sandbox.CreateOptions, scope *sandbox.Scope, region string) (createPlan, error) {
	var plan createPlan
	if err := validateProviderCreate(request); err != nil {
		return plan, err
	}
	if err := sandbox.ValidateCreateOptions(request); err != nil {
		return plan, err
	}
	params, rest, err := mapCommonCreateFields(request)
	if err != nil {
		return plan, err
	}
	plan.params = params
	// Modal's API requires app and image; no image or app defaults are invented.
	if scope.GetAppName() == "" {
		return plan, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindInvalidArgument,
			Provider:  "modal",
			Operation: "create",
			Field:     sandbox.ConfigFieldScopeAppName,
			Message:   "sandbox-kit modal: Config.Scope.AppName is required",
		}, nil)
	}
	plan.app = scope.GetAppName()
	plan.environment = scope.GetEnvironment()
	if region != "" {
		plan.params.Regions = []string{region}
	}
	image := request.GetSource().GetImage()
	if image == nil {
		return plan, sandbox.NewError(sandbox.ErrorInfo{
			Kind:      sandbox.ErrorKindUnsupported,
			Provider:  "modal",
			Operation: "create",
			Field:     sandbox.CreateFieldSource,
			Message:   "sandbox-kit modal: explicit image source is required; default/snapshot/pool mappings are not implemented",
		}, nil)
	}
	plan.image = image.Reference
	rest.Source = nil
	if request.Runtime != nil {
		config := *request.Runtime
		plan.params.Workdir = config.GetWorkingDirectory()
		plan.params.Command = config.Entrypoint
		plan.params.PTY = config.GetPTY()
		config.WorkingDirectory, config.Entrypoint, config.PTY = nil, nil, nil
		if err := sandbox.RejectUnmapped("modal runtime (language must come from image)", &config); err != nil {
			return plan, err
		}
		rest.Runtime = nil
	}
	if request.Isolation != nil {
		config := *request.Isolation
		switch config.Kind {
		case sandbox.IsolationKindDefault:
		case sandbox.IsolationKindLinuxVM:
			plan.params.Runtime = sdk.SandboxRuntimeVM
		case sandbox.IsolationKindContainer:
			plan.params.Runtime = sdk.SandboxRuntimeGVisor
		default:
			return plan, sandbox.NewError(sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindUnsupported,
				Provider:  "modal",
				Operation: "create",
				Field:     sandbox.CreateFieldIsolationKind,
				Message:   "sandbox-kit modal: requested isolation mapping is unsupported",
			}, nil)
		}
		config.Kind = 0
		if err := sandbox.RejectUnmapped("modal isolation", &config); err != nil {
			return plan, err
		}
		rest.Isolation = nil
	}
	if request.Resources != nil {
		resources, remaining, err := mapResources(request.Resources)
		if err != nil {
			return plan, err
		}
		if err := sandbox.RejectUnmapped("modal resources", &remaining); err != nil {
			return plan, err
		}
		plan.params.CPU, plan.params.CPULimit = resources.CPU, resources.CPULimit
		plan.params.MemoryMiB, plan.params.MemoryLimitMiB = resources.MemoryMiB, resources.MemoryLimitMiB
		rest.Resources = nil
	}
	if request.Placement != nil {
		config := *request.Placement
		plan.params.Cloud = config.GetCloud()
		if config.Regions != nil {
			plan.params.Regions = config.Regions
		}
		config.Cloud, config.Regions = nil, nil
		if err := sandbox.RejectUnmapped("modal placement", &config); err != nil {
			return plan, err
		}
		rest.Placement = nil
	}
	if request.Network != nil {
		config := *request.Network
		switch config.Egress {
		case sandbox.EgressModeDefault, sandbox.EgressModeAllowAll:
		case sandbox.EgressModeBlockAll:
			plan.params.BlockNetwork = true
		case sandbox.EgressModeRestricted:
		default:
			return plan, sandbox.NewError(sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindInvalidArgument,
				Provider:  "modal",
				Operation: "create",
				Field:     sandbox.CreateFieldNetworkEgress,
				Message:   "sandbox-kit modal: unknown egress mode",
			}, nil)
		}
		config.Egress = 0
		if config.OutboundCIDRs != nil {
			plan.params.OutboundCIDRAllowlist = &sdk.Allowlist{Entries: config.OutboundCIDRs.Entries}
		}
		config.OutboundCIDRs = nil
		if config.OutboundDomains != nil {
			plan.params.OutboundDomainAllowlist = &sdk.Allowlist{Entries: config.OutboundDomains.Entries}
		}
		config.OutboundDomains = nil
		if config.InboundCIDRs != nil {
			if len(config.InboundCIDRs.Entries) == 0 {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldNetworkInboundCIDRs,
					Message:   "sandbox-kit modal: empty inbound deny-all cannot be represented by native allow-all",
				}, nil)
			}
			plan.params.InboundCIDRAllowlist = config.InboundCIDRs.Entries
		}
		config.InboundCIDRs = nil
		plan.params.I6PN = config.GetPrivateNetwork()
		plan.params.CustomDomain = config.GetCustomDomain()
		config.PrivateNetwork, config.CustomDomain = nil, nil
		for _, port := range config.Ports {
			switch port.Transport {
			case sandbox.PortTransportTLS:
				plan.params.EncryptedPorts = append(plan.params.EncryptedPorts, int(port.Port))
			case sandbox.PortTransportHTTP2TLS:
				plan.params.H2Ports = append(plan.params.H2Ports, int(port.Port))
			case sandbox.PortTransportPlain:
				plan.params.UnencryptedPorts = append(plan.params.UnencryptedPorts, int(port.Port))
			default:
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldNetworkPortsTransport,
					Message:   "sandbox-kit modal: explicit port transport is required",
				}, nil)
			}
		}
		config.Ports = nil
		if err := sandbox.RejectUnmapped("modal network", &config); err != nil {
			return plan, err
		}
		rest.Network = nil
	}
	if request.Security != nil {
		config := *request.Security
		plan.params.IncludeOidcIdentityToken = config.GetWorkloadIdentity()
		config.WorkloadIdentity = nil
		for _, secret := range config.Secrets {
			if secret.EnvironmentVariable != nil {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldSecuritySecretsEnvironmentVariable,
					Message:   "sandbox-kit modal: per-variable secret renaming is unsupported",
				}, nil)
			}
			if secret.Injection != sandbox.SecretInjectionDefault && secret.Injection != sandbox.SecretInjectionEnvironment {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldSecuritySecretsInjection,
					Message:   "sandbox-kit modal: egress placeholder injection mapping is unsupported",
				}, nil)
			}
		}
		plan.secrets = config.Secrets
		config.Secrets = nil
		if err := sandbox.RejectUnmapped("modal security", &config); err != nil {
			return plan, err
		}
		rest.Security = nil
	}
	if request.Lifetime != nil {
		config := *request.Lifetime
		if config.MaximumLifetime != nil {
			mapped, _, err := mapLifetimeDuration(&config)
			if err != nil {
				return plan, err
			}
			plan.params.Timeout = mapped.Timeout
		}
		config.MaximumLifetime = nil
		if config.IdleTerminate != nil {
			policy := config.IdleTerminate
			switch policy.Mode {
			case sandbox.PolicyModeDefault:
			case sandbox.PolicyModeDisabled:
				plan.params.IdleTimeout = 0
			case sandbox.PolicyModeAfter:
				if policy.After == nil {
					return plan, sandbox.NewError(sandbox.ErrorInfo{
						Kind:      sandbox.ErrorKindInvalidArgument,
						Provider:  "modal",
						Operation: "create",
						Field:     sandbox.CreateFieldLifetimeIdleTerminateAfter,
						Message:   "sandbox-kit modal: idle duration missing",
					}, nil)
				}
				mapped, _, err := mapIdleDuration(policy)
				if err != nil {
					return plan, err
				}
				plan.params.IdleTimeout = mapped.IdleTimeout
			default:
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldLifetimeIdleTerminateMode,
					Message:   "sandbox-kit modal: unknown idle policy",
				}, nil)
			}
		}
		config.IdleTerminate = nil
		if err := sandbox.RejectUnmapped("modal lifetime", &config); err != nil {
			return plan, err
		}
		rest.Lifetime = nil
	}
	if request.Readiness != nil {
		var err error
		if request.Readiness.TCPPort != nil {
			plan.params.ReadinessProbe, err = sdk.NewTCPProbe(int(*request.Readiness.TCPPort), nil)
		} else {
			plan.params.ReadinessProbe, err = sdk.NewExecProbe(request.Readiness.Command.Argv, nil)
		}
		if err != nil {
			return plan, err
		}
		rest.Readiness = nil
	}
	if request.Provisioning != nil {
		config := *request.Provisioning
		config.Timeout = nil // Client.Create applies this context deadline.
		switch config.WaitFor {
		case sandbox.WaitConditionDefault, sandbox.WaitConditionScheduled:
		case sandbox.WaitConditionReady:
			if request.Provisioning.Timeout == nil || *request.Provisioning.Timeout <= 0 {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldProvisioningTimeout,
					Message:   "sandbox-kit modal: readiness waiting requires a positive creation timeout",
				}, nil)
			}
			plan.readyTimeout = *request.Provisioning.Timeout
			if plan.params.ReadinessProbe == nil {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldReadiness,
					Message:   "sandbox-kit modal: ready wait requires a readiness probe",
				}, nil)
			}
			plan.ready = true
		default:
			return plan, sandbox.NewError(sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindUnsupported,
				Provider:  "modal",
				Operation: "create",
				Field:     sandbox.CreateFieldProvisioningWaitFor,
				Message:   "sandbox-kit modal: requested wait condition is unsupported",
			}, nil)
		}
		config.WaitFor = 0
		if err := sandbox.RejectUnmapped("modal creation", &config); err != nil {
			return plan, err
		}
		rest.Provisioning = nil
	}
	if request.Observability != nil {
		config := *request.Observability
		plan.params.Verbose = config.GetVerbose()
		config.Verbose = nil
		if err := sandbox.RejectUnmapped("modal observability", &config); err != nil {
			return plan, err
		}
		rest.Observability = nil
	}
	if err := sandbox.RejectUnmapped("modal", &rest); err != nil {
		return plan, err
	}
	return plan, nil
}
