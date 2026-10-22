// This file turns a shared creation request into Modal app, image, and secret lookups,
// then calls the official SDK. Generated mappings cover field conversion. This file
// covers the call order those mappings cannot express.
//
// Example:
//
//	client, err := sandbox.NewClient(modal.New(), &sandbox.Config{})
//	instance, err := client.Create(ctx, &sandbox.CreateOptions{
//	    Source: &sandbox.SandboxSource{
//	        Image: &sandbox.ImageSource{Reference: "alpine:3.21"},
//	    },
//	})

package modal

import (
	"context"
	"math"
	"time"

	sdk "github.com/modal-labs/modal-client/go"
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/protobuf/proto"
)

var _ sandbox.SandboxAttacher = (*backend)(nil)
var _ sandbox.SandboxTerminator = (*backend)(nil)

// createPlan is the Modal call sequence for one creation request.
type createPlan struct {
	// app is the existing Modal app name.
	app string
	// environment is the Modal environment used for app and secret lookup.
	environment string
	// image is the registry image reference.
	image string
	// params is the native creation payload after shared fields are mapped.
	params sdk.SandboxCreateParams
	// ready reports whether creation should wait for a readiness probe.
	ready bool
	// readyTimeout is the wait limit used when ready is true.
	readyTimeout time.Duration
	// secrets are provider-managed secret names to resolve before creation.
	secrets []*sandbox.SecretReference
	// proxy is the optional provider-managed proxy name to resolve before creation.
	proxy string
}

// create looks up the app, secrets, and image, then creates the sandbox.
// A nil request uses provider defaults. The returned metadata is shared,
// not a Modal object.
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
	if plan.proxy != "" {
		proxy, err := a.client.Proxies.FromName(ctx, plan.proxy, &sdk.ProxyFromNameParams{Environment: plan.environment})
		if err != nil {
			return nil, err
		}
		plan.params.Proxy = proxy
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
			Field:     sandbox.InfoFieldId,
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
	info.State = sandbox.SandboxStateRunning
	info.Labels = plan.params.Tags
	if request.Isolation != nil {
		info.Isolation = proto.Clone(request.Isolation).(*sandbox.IsolationConfig)
		info.Origins[sandbox.InfoFieldIsolation] = sandbox.ValueOriginRequest
	}
	if info.Name != nil {
		info.Origins["name"] = sandbox.ValueOriginRequest
	}
	if info.Labels != nil {
		info.Origins["labels"] = sandbox.ValueOriginRequest
	}
	return &sandbox.CreateResult{Sandbox: &info}, nil
}

// AttachSandbox constructs a Modal handle from an ID. Modal FromID performs no
// provider request, so this deliberately does not implement SandboxRefresher.
func (a *backend) AttachSandbox(ctx context.Context, id string) (*sandbox.SandboxInfo, error) {
	instance, err := a.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return nil, mapProviderError(err, "attach")
	}
	info, err := mapSandboxIdentity(instance)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// TerminateSandbox permanently finishes a Modal sandbox execution.
func (a *backend) TerminateSandbox(ctx context.Context, id string, options *sandbox.TerminateOptions) (*sandbox.ExitResult, error) {
	instance, err := a.client.Sandboxes.FromID(ctx, id, nil)
	if err != nil {
		return nil, mapProviderError(err, "terminate")
	}
	wait := options != nil && options.GetWait()
	code, err := instance.Terminate(ctx, &sdk.SandboxTerminateParams{Wait: wait})
	if err != nil {
		return nil, mapProviderError(err, "terminate")
	}
	result := &sandbox.ExitResult{}
	if wait {
		if code < math.MinInt32 || code > math.MaxInt32 {
			return nil, sandbox.NewError(sandbox.ErrorInfo{Kind: sandbox.ErrorKindInvalidResponse, Provider: "modal", Operation: "terminate", Field: "exit_code", Message: "sandbox-kit modal: exit code exceeds the shared int32 range"}, nil)
		}
		result.ExitCode = sandbox.Value(int32(code))
	}
	return result, nil
}

// planCreate checks the request and builds the Modal call sequence
// without contacting Modal.
func planCreate(request *sandbox.CreateOptions, scope *sandboxv1.Scope, region string) (createPlan, error) {
	var plan createPlan
	if err := validateProviderCreate(request); err != nil {
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
		plan.params.PTY = config.GetPty()
		config.WorkingDirectory, config.Entrypoint, config.Pty = nil, nil, nil
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
	if err := planCreateNetwork(request, &plan, &rest); err != nil {
		return plan, err
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
	if err := planCreateLifetime(request, &plan, &rest); err != nil {
		return plan, err
	}
	if request.Readiness != nil {
		var err error
		if request.Readiness.GetTcpPort() != 0 {
			plan.params.ReadinessProbe, err = sdk.NewTCPProbe(int(request.Readiness.GetTcpPort()), nil)
		} else if command := request.Readiness.GetCommand(); command != nil {
			plan.params.ReadinessProbe, err = sdk.NewExecProbe(command.Argv, nil)
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
			timeout := sandbox.OptionalDuration(request.Provisioning.Timeout)
			if timeout == nil || *timeout <= 0 {
				return plan, sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldProvisioningTimeout,
					Message:   "sandbox-kit modal: readiness waiting requires a positive creation timeout",
				}, nil)
			}
			plan.readyTimeout = *timeout
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

func planCreateNetwork(request *sandbox.CreateOptions, plan *createPlan, rest *sandbox.CreateOptions) error {
	if request.Network != nil {
		config := *request.Network
		switch config.Egress {
		case sandbox.EgressModeDefault, sandbox.EgressModeAllowAll:
		case sandbox.EgressModeBlockAll:
			plan.params.BlockNetwork = true
		case sandbox.EgressModeRestricted:
		default:
			return sandbox.NewError(sandbox.ErrorInfo{
				Kind:      sandbox.ErrorKindInvalidArgument,
				Provider:  "modal",
				Operation: "create",
				Field:     sandbox.CreateFieldNetworkEgress,
				Message:   "sandbox-kit modal: unknown egress mode",
			}, nil)
		}
		config.Egress = 0
		if config.OutboundCidrs != nil {
			plan.params.OutboundCIDRAllowlist = &sdk.Allowlist{Entries: config.OutboundCidrs.Entries}
		}
		config.OutboundCidrs = nil
		if config.OutboundDomains != nil {
			plan.params.OutboundDomainAllowlist = &sdk.Allowlist{Entries: config.OutboundDomains.Entries}
		}
		config.OutboundDomains = nil
		if config.InboundCidrs != nil {
			if len(config.InboundCidrs.Entries) == 0 {
				return sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindUnsupported,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldNetworkInboundCidrs,
					Message:   "sandbox-kit modal: empty inbound deny-all cannot be represented by native allow-all",
				}, nil)
			}
			plan.params.InboundCIDRAllowlist = config.InboundCidrs.Entries
		}
		config.InboundCidrs = nil
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
				return sandbox.NewError(sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "modal",
					Operation: "create",
					Field:     sandbox.CreateFieldNetworkPortsTransport,
					Message:   "sandbox-kit modal: explicit port transport is required",
				}, nil)
			}
		}
		config.Ports = nil
		plan.proxy = config.GetProxyReference()
		config.ProxyReference = nil
		if err := sandbox.RejectUnmapped("modal network", &config); err != nil {
			return err
		}
		rest.Network = nil
	}
	return nil
}

func planCreateLifetime(request *sandbox.CreateOptions, plan *createPlan, rest *sandbox.CreateOptions) error {
	if request.Lifetime != nil {
		config := *request.Lifetime
		if config.MaximumLifetime != nil {
			mapped, _, err := mapLifetimeDuration(&config)
			if err != nil {
				return err
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
					return sandbox.NewError(sandbox.ErrorInfo{
						Kind:      sandbox.ErrorKindInvalidArgument,
						Provider:  "modal",
						Operation: "create",
						Field:     sandbox.CreateFieldLifetimeIdleTerminateAfter,
						Message:   "sandbox-kit modal: idle duration missing",
					}, nil)
				}
				mapped, _, err := mapIdleDuration(policy)
				if err != nil {
					return err
				}
				plan.params.IdleTimeout = mapped.IdleTimeout
			default:
				return sandbox.NewError(sandbox.ErrorInfo{
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
			return err
		}
		rest.Lifetime = nil
	}
	return nil
}

