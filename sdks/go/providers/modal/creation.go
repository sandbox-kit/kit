package modal

import (
	"time"

	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

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

// creation translates one CreateOptions value into a Modal call sequence.
// Methods own the request, the unmapped remainder, and the native plan.
type creation struct {
	request *sandbox.CreateOptions
	rest    *sandbox.CreateOptions
	plan    createPlan
}

// planCreate checks the request and builds the Modal call sequence
// without contacting Modal.
func planCreate(request *sandbox.CreateOptions, scope *sandbox.Scope, region string) (createPlan, error) {
	planned := &creation{request: request}
	if err := planned.assemble(scope, region); err != nil {
		return createPlan{}, err
	}
	return planned.plan, nil
}

func (c *creation) assemble(scope *sandbox.Scope, region string) error {
	if err := validateProviderCreate(c.request); err != nil {
		return err
	}
	params, rest, err := mapCommonCreateFields(c.request)
	if err != nil {
		return err
	}
	c.rest = rest
	c.plan.params = params
	for _, apply := range []func() error{
		func() error { return c.app(scope, region) },
		c.runtime,
		c.isolation,
		c.resources,
		c.placement,
		c.network,
		c.security,
		c.lifetime,
		c.readiness,
		c.provisioning,
		c.observability,
	} {
		if err := apply(); err != nil {
			return err
		}
	}
	return sandbox.RejectUnmapped("modal", c.rest)
}

func (c *creation) app(scope *sandbox.Scope, region string) error {
	// Modal's API requires app and image; no image or app defaults are invented.
	if scope.GetAppName() == "" {
		return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.ConfigFieldScopeAppName, "Config.Scope.AppName is required")
	}
	c.plan.app = scope.GetAppName()
	c.plan.environment = scope.GetEnvironment()
	if region != "" {
		c.plan.params.Regions = []string{region}
	}
	image := c.request.GetSource().GetImage()
	if image == nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSource, "an explicit image source is required")
	}
	c.plan.image = image.Reference
	c.rest.Source = nil
	return nil
}

func (c *creation) runtime() error {
	if c.request.Runtime == nil {
		return nil
	}
	config := cloneMessage(c.request.Runtime)
	c.plan.params.Workdir = config.GetWorkingDirectory()
	c.plan.params.Command = config.Entrypoint
	c.plan.params.PTY = config.GetPty()
	config.WorkingDirectory, config.Entrypoint, config.Pty = nil, nil, nil
	if err := sandbox.RejectUnmapped("modal runtime (language must come from image)", config); err != nil {
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
	switch config.Kind {
	case sandbox.IsolationKindDefault:
	case sandbox.IsolationKindLinuxVM:
		c.plan.params.Runtime = sdk.SandboxRuntimeVM
	case sandbox.IsolationKindContainer:
		c.plan.params.Runtime = sdk.SandboxRuntimeGVisor
	default:
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldIsolationKind, "requested isolation mapping is unsupported")
	}
	config.Kind = 0
	if err := sandbox.RejectUnmapped("modal isolation", config); err != nil {
		return err
	}
	c.rest.Isolation = nil
	return nil
}

func (c *creation) resources() error {
	if c.request.Resources == nil {
		return nil
	}
	resources, remaining, err := mapResources(c.request.Resources)
	if err != nil {
		return err
	}
	if err := sandbox.RejectUnmapped("modal resources", remaining); err != nil {
		return err
	}
	c.plan.params.CPU, c.plan.params.CPULimit = resources.CPU, resources.CPULimit
	c.plan.params.MemoryMiB, c.plan.params.MemoryLimitMiB = resources.MemoryMiB, resources.MemoryLimitMiB
	c.rest.Resources = nil
	return nil
}

func (c *creation) placement() error {
	if c.request.Placement == nil {
		return nil
	}
	config := cloneMessage(c.request.Placement)
	c.plan.params.Cloud = config.GetCloud()
	if config.Regions != nil {
		c.plan.params.Regions = config.Regions
	}
	config.Cloud, config.Regions = nil, nil
	if err := sandbox.RejectUnmapped("modal placement", config); err != nil {
		return err
	}
	c.rest.Placement = nil
	return nil
}

func (c *creation) security() error {
	if c.request.Security == nil {
		return nil
	}
	config := cloneMessage(c.request.Security)
	c.plan.params.IncludeOidcIdentityToken = config.GetWorkloadIdentity()
	config.WorkloadIdentity = nil
	for _, secret := range config.Secrets {
		if secret.EnvironmentVariable != nil {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSecuritySecretsEnvironmentVariable, "per-variable secret renaming is unsupported")
		}
		if secret.Injection != sandbox.SecretInjectionDefault && secret.Injection != sandbox.SecretInjectionEnvironment {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldSecuritySecretsInjection, "egress placeholder injection mapping is unsupported")
		}
	}
	c.plan.secrets = config.Secrets
	config.Secrets = nil
	if err := sandbox.RejectUnmapped("modal security", config); err != nil {
		return err
	}
	c.rest.Security = nil
	return nil
}

func (c *creation) readiness() error {
	if c.request.Readiness == nil {
		return nil
	}
	var err error
	if c.request.Readiness.GetTcpPort() != 0 {
		c.plan.params.ReadinessProbe, err = sdk.NewTCPProbe(int(c.request.Readiness.GetTcpPort()), nil)
	} else if command := c.request.Readiness.GetCommand(); command != nil {
		c.plan.params.ReadinessProbe, err = sdk.NewExecProbe(command.Argv, nil)
	}
	if err != nil {
		return err
	}
	c.rest.Readiness = nil
	return nil
}

func (c *creation) provisioning() error {
	if c.request.Provisioning == nil {
		return nil
	}
	config := cloneMessage(c.request.Provisioning)
	config.Timeout = nil // Client.Create applies this context deadline.
	switch config.WaitFor {
	case sandbox.WaitConditionDefault, sandbox.WaitConditionScheduled:
	case sandbox.WaitConditionReady:
		timeout := sandbox.OptionalDuration(c.request.Provisioning.Timeout)
		if timeout == nil || *timeout <= 0 {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldProvisioningTimeout, "readiness waiting requires a positive creation timeout")
		}
		c.plan.readyTimeout = *timeout
		if c.plan.params.ReadinessProbe == nil {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldReadiness, "ready wait requires a readiness probe")
		}
		c.plan.ready = true
	default:
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldProvisioningWaitFor, "requested wait condition is unsupported")
	}
	config.WaitFor = 0
	if err := sandbox.RejectUnmapped("modal creation", config); err != nil {
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
	c.plan.params.Verbose = config.GetVerbose()
	config.Verbose = nil
	if err := sandbox.RejectUnmapped("modal observability", config); err != nil {
		return err
	}
	c.rest.Observability = nil
	return nil
}

func (c *creation) network() error {
	if c.request.Network != nil {
		config := cloneMessage(c.request.Network)
		switch config.Egress {
		case sandbox.EgressModeDefault, sandbox.EgressModeAllowAll:
		case sandbox.EgressModeBlockAll:
			c.plan.params.BlockNetwork = true
		case sandbox.EgressModeRestricted:
		default:
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkEgress, "unknown egress mode")
		}
		config.Egress = 0
		if config.OutboundCidrs != nil {
			c.plan.params.OutboundCIDRAllowlist = &sdk.Allowlist{Entries: config.OutboundCidrs.Entries}
		}
		config.OutboundCidrs = nil
		if config.OutboundDomains != nil {
			c.plan.params.OutboundDomainAllowlist = &sdk.Allowlist{Entries: config.OutboundDomains.Entries}
		}
		config.OutboundDomains = nil
		if config.InboundCidrs != nil {
			if len(config.InboundCidrs.Entries) == 0 {
				return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldNetworkInboundCidrs, "empty inbound deny-all cannot be represented by native allow-all")
			}
			c.plan.params.InboundCIDRAllowlist = config.InboundCidrs.Entries
		}
		config.InboundCidrs = nil
		c.plan.params.I6PN = config.GetPrivateNetwork()
		c.plan.params.CustomDomain = config.GetCustomDomain()
		config.PrivateNetwork, config.CustomDomain = nil, nil
		for _, port := range config.Ports {
			switch port.Transport {
			case sandbox.PortTransportTLS:
				c.plan.params.EncryptedPorts = append(c.plan.params.EncryptedPorts, int(port.Port))
			case sandbox.PortTransportHTTP2TLS:
				c.plan.params.H2Ports = append(c.plan.params.H2Ports, int(port.Port))
			case sandbox.PortTransportPlain:
				c.plan.params.UnencryptedPorts = append(c.plan.params.UnencryptedPorts, int(port.Port))
			default:
				return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkPortsTransport, "explicit port transport is required")
			}
		}
		config.Ports = nil
		c.plan.proxy = config.GetProxyReference()
		config.ProxyReference = nil
		if err := sandbox.RejectUnmapped("modal network", config); err != nil {
			return err
		}
		c.rest.Network = nil
	}
	return nil
}

func (c *creation) lifetime() error {
	if c.request.Lifetime != nil {
		config := cloneMessage(c.request.Lifetime)
		if config.MaximumLifetime != nil {
			mapped, _, err := mapLifetimeDuration(config)
			if err != nil {
				return err
			}
			c.plan.params.Timeout = mapped.Timeout
		}
		config.MaximumLifetime = nil
		if config.IdleTerminate != nil {
			policy := config.IdleTerminate
			switch policy.Mode {
			case sandbox.PolicyModeDefault:
			case sandbox.PolicyModeDisabled:
				c.plan.params.IdleTimeout = 0
			case sandbox.PolicyModeAfter:
				if policy.After == nil {
					return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldLifetimeIdleTerminateAfter, "idle duration missing")
				}
				mapped, _, err := mapIdleDuration(policy)
				if err != nil {
					return err
				}
				c.plan.params.IdleTimeout = mapped.IdleTimeout
			default:
				return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldLifetimeIdleTerminateMode, "unknown idle policy")
			}
		}
		config.IdleTerminate = nil
		if err := sandbox.RejectUnmapped("modal lifetime", config); err != nil {
			return err
		}
		c.rest.Lifetime = nil
	}
	return nil
}
