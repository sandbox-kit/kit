package blaxel

import (
	"time"

	blaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// createPlan is the Blaxel SDK argument for one creation request.
type createPlan struct {
	params blaxel.SandboxNewParams
}

// creation translates one CreateOptions value into native Blaxel arguments.
type creation struct {
	request *sandbox.CreateOptions
	rest    *sandbox.CreateOptions
	params  blaxel.SandboxNewParams
	region  string
}

var reservedBlaxelPorts = map[uint32]struct{}{
	80:   {},
	443:  {},
	8080: {},
}

// planCreate checks the request and builds SDK arguments without contacting Blaxel.
func planCreate(request *sandbox.CreateOptions, region string) (createPlan, error) {
	if request == nil {
		request = &sandbox.CreateOptions{}
	}
	if err := validateProviderCreate(request); err != nil {
		return createPlan{}, err
	}
	if region == "" {
		return createPlan{}, providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.ConfigFieldRegion, "Config.Region is required")
	}
	params, rest, err := mapCommonCreateFields(request)
	if err != nil {
		return createPlan{}, err
	}
	planned := &creation{request: request, rest: rest, params: params, region: region}
	planned.params.Sandbox.Spec.Region = blaxel.String(region)
	for _, apply := range []func() error{
		planned.source,
		planned.resources,
		planned.lifetime,
		planned.runtime,
		planned.isolation,
		planned.security,
		planned.network,
		planned.provisioning,
		planned.readiness,
		planned.observability,
	} {
		if err := apply(); err != nil {
			return createPlan{}, err
		}
	}
	if err := sdk.RejectUnmapped("blaxel", "", planned.rest); err != nil {
		return createPlan{}, err
	}
	return createPlan{params: planned.params}, nil
}

func (c *creation) source() error {
	image := c.request.GetSource().GetImage()
	c.params.Sandbox.Spec.Runtime.Image = blaxel.String(image.Reference)
	c.rest.Source = nil
	return nil
}

func (c *creation) resources() error {
	if c.request.Resources == nil {
		return nil
	}
	mapped, remaining, err := mapResources(c.request.Resources)
	if err != nil {
		return err
	}
	if hasUnmappedResourceFields(remaining) {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldResources, "one or more Blaxel resource fields are not mapped")
	}
	if err := sdk.RejectUnmapped("blaxel", sandbox.CreateFieldResources, remaining); err != nil {
		return err
	}
	if mapped.Memory.Valid() {
		c.params.Sandbox.Spec.Runtime.Memory = mapped.Memory
	}
	c.rest.Resources = nil
	return nil
}

func (c *creation) lifetime() error {
	if c.request.Lifetime == nil {
		return nil
	}
	config := cloneMessage(c.request.Lifetime)
	if config.MaximumLifetime != nil {
		formatted, err := formatBlaxelDuration(config.MaximumLifetime.AsDuration(), 5*time.Minute)
		if err != nil {
			return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldLifetimeMaximumLifetime, err.Error())
		}
		c.params.Sandbox.Spec.Runtime.Ttl = blaxel.String(formatted)
	}
	config.MaximumLifetime = nil
	if config.IdleTerminate != nil {
		policy, err := mapIdleTerminate(config.IdleTerminate)
		if err != nil {
			return err
		}
		if policy != nil {
			c.params.Sandbox.Spec.Lifecycle.ExpirationPolicies = append(
				c.params.Sandbox.Spec.Lifecycle.ExpirationPolicies,
				*policy,
			)
		}
	}
	config.IdleTerminate = nil
	if hasUnsupportedLifetimeFields(config) {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldLifetime, "one or more Blaxel lifetime fields are not mapped")
	}
	if err := sdk.RejectUnmapped("blaxel", sandbox.CreateFieldLifetime, config); err != nil {
		return err
	}
	c.rest.Lifetime = nil
	return nil
}

func mapIdleTerminate(action *sandbox.AutomaticAction) (*blaxel.ExpirationPolicyParam, error) {
	if action == nil {
		return nil, nil
	}
	switch action.Mode {
	case sandbox.PolicyModeDefault:
		return nil, nil
	case sandbox.PolicyModeDisabled:
		return nil, providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldLifetimeIdleTerminate, "disabling idle terminate is not representable")
	case sandbox.PolicyModeAfter:
		if action.After == nil {
			return nil, providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldLifetimeIdleTerminateAfter, "duration is required")
		}
		formatted, err := formatBlaxelDuration(action.After.AsDuration(), time.Minute)
		if err != nil {
			return nil, providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldLifetimeIdleTerminateAfter, err.Error())
		}
		return &blaxel.ExpirationPolicyParam{
			Action: blaxel.ExpirationPolicyActionDelete,
			Type:   blaxel.ExpirationPolicyTypeTtlIdle,
			Value:  blaxel.String(formatted),
		}, nil
	default:
		return nil, providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldLifetimeIdleTerminateMode, "unknown policy mode")
	}
}

func (c *creation) runtime() error {
	if c.request.Runtime == nil {
		return nil
	}
	config := cloneMessage(c.request.Runtime)
	if config.WorkingDirectory != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldRuntimeWorkingDirectory, "Blaxel create-time working directory is not mapped")
	}
	if hasUnsupportedRuntimeFields(config) {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldRuntime, "runtime user, language, version, entrypoint, and PTY are not mapped by Blaxel")
	}
	if err := sdk.RejectUnmapped("blaxel", sandbox.CreateFieldRuntime, config); err != nil {
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
	if hasUnsupportedIsolation(config) {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldIsolation, "Blaxel isolation class is not mapped")
	}
	if err := sdk.RejectUnmapped("blaxel", sandbox.CreateFieldIsolation, config); err != nil {
		return err
	}
	c.rest.Isolation = nil
	return nil
}

func (c *creation) security() error {
	return rejectIfSet(c.request.Security != nil, sandbox.CreateFieldSecurity, "Blaxel security fields are not mapped")
}

func (c *creation) network() error {
	if c.request.Network == nil {
		return nil
	}
	config := cloneMessage(c.request.Network)
	if hasUnsupportedNetwork(config) {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldNetwork, "requested network settings are unsupported by Blaxel")
	}
	if config.OutboundDomains != nil {
		c.params.Sandbox.Spec.Network.AllowedDomains = append([]string(nil), config.OutboundDomains.Entries...)
		config.OutboundDomains = nil
	}
	if len(config.Ports) > 0 {
		ports, err := mapPorts(config.Ports)
		if err != nil {
			return err
		}
		c.params.Sandbox.Spec.Runtime.Ports = ports
		config.Ports = nil
	}
	config.Egress = 0
	if err := sdk.RejectUnmapped("blaxel", sandbox.CreateFieldNetwork, config); err != nil {
		return err
	}
	c.rest.Network = nil
	return nil
}

func mapPorts(bindings []*sandbox.PortBinding) ([]blaxel.PortParam, error) {
	ports := make([]blaxel.PortParam, 0, len(bindings))
	for _, binding := range bindings {
		if binding == nil {
			continue
		}
		if isReservedPort(binding.Port) {
			return nil, providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkPortsPort, "ports 80, 443, and 8080 are reserved by Blaxel")
		}
		protocol, err := blaxelPortProtocol(binding.Transport)
		if err != nil {
			return nil, err
		}
		ports = append(ports, blaxel.PortParam{
			Target:   int64(binding.Port),
			Protocol: protocol,
		})
	}
	return ports, nil
}

func blaxelPortProtocol(transport sandbox.PortTransport) (blaxel.PortProtocol, error) {
	switch transport {
	case sandbox.PortTransportDefault:
		return "", nil
	case sandbox.PortTransportPlain:
		return blaxel.PortProtocolTcp, nil
	case sandbox.PortTransportTLS:
		return blaxel.PortProtocolTls, nil
	default:
		// Blaxel PortProtocol is HTTP|TCP|UDP|TLS only (Kit Http2Tls included).
		return "", providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldNetworkPortsTransport, "Http2Tls is not a Blaxel port protocol (HTTP|TCP|UDP|TLS)")
	}
}

func (c *creation) provisioning() error {
	if c.request.Provisioning == nil {
		return nil
	}
	config := cloneMessage(c.request.Provisioning)
	// Client.Create already applies Timeout as the context deadline.
	config.Timeout = nil
	if config.QueueTimeout != nil {
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldProvisioningQueueTimeout, "Blaxel queue timeout is unsupported")
	}
	config.QueueTimeout = nil
	switch config.WaitFor {
	case sandbox.WaitConditionDefault, sandbox.WaitConditionScheduled:
	default:
		return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldProvisioningWaitFor, "requested wait condition is unsupported")
	}
	config.WaitFor = 0
	if err := sdk.RejectUnmapped("blaxel", sandbox.CreateFieldProvisioning, config); err != nil {
		return err
	}
	c.rest.Provisioning = nil
	return nil
}

func (c *creation) readiness() error {
	return rejectIfSet(c.request.Readiness != nil, sandbox.CreateFieldReadiness, "Blaxel readiness probes are not mapped")
}

func (c *creation) observability() error {
	return rejectIfSet(c.request.Observability != nil, sandbox.CreateFieldObservability, "Blaxel observability fields are not mapped")
}

func rejectIfSet(set bool, field, message string) error {
	if !set {
		return nil
	}
	return providerCheckError(sandbox.ErrorKindUnsupported, field, message)
}

func isReservedPort(port uint32) bool {
	_, ok := reservedBlaxelPorts[port]
	return ok
}

func hasUnmappedResourceFields(r *sandbox.Resources) bool {
	return r.CpuCores != nil || r.CpuLimitCores != nil || r.MemoryLimitMib != nil || r.DiskMib != nil
}

func hasUnsupportedLifetimeFields(c *sandbox.LifetimePolicy) bool {
	return c.Ephemeral != nil || c.IdleStop != nil || c.IdlePause != nil ||
		c.StoppedArchive != nil || c.StoppedDelete != nil
}

func hasUnsupportedRuntimeFields(c *sandbox.RuntimeConfig) bool {
	return c.Entrypoint != nil || c.Pty != nil || c.User != nil || c.Language != nil || c.Version != nil
}

func hasUnsupportedIsolation(c *sandbox.IsolationConfig) bool {
	return c.Kind != sandbox.IsolationKindDefault || c.ProviderKind != nil || c.NestedVirtualization != nil
}

func hasUnsupportedNetwork(c *sandbox.NetworkConfig) bool {
	return c.Egress != sandbox.EgressModeDefault ||
		c.OutboundCidrs != nil ||
		c.InboundCidrs != nil ||
		c.PrivateNetwork != nil ||
		c.LinkedSandbox != nil ||
		c.ProxyReference != nil ||
		c.OutboundProxyUrl != nil ||
		c.CustomDomain != nil
}
