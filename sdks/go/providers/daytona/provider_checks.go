package daytona

import (
	"net/netip"
	"net/url"
	"regexp"
	"strings"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var providerEnvNamePattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
var providerDomainPattern = regexp.MustCompile(`(?i)^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.?$`)

// validateProviderCreate rejects creation settings this provider cannot honor.
// It runs before the native create call.
//
// Rules are inventoried in specs/providers/daytona.yaml checks and kept in sync
// here by hand (no longer emitted).
func validateProviderCreate(request *sandbox.CreateOptions) error {
	if request == nil {
		return nil
	}
	if request.Runtime != nil && request.Runtime.Language != nil {
		language := *request.Runtime.Language
		if language != "python" && language != "javascript" && language != "typescript" {
			return providerCheckError(sandbox.ErrorKindUnsupported, sandbox.CreateFieldRuntimeLanguage, "runtime.language is outside supported values")
		}
	}
	if network := request.Network; network != nil {
		if network.OutboundProxyUrl != nil {
			parsed, err := url.Parse(*network.OutboundProxyUrl)
			if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
				return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkOutboundProxyUrl, "network.outbound_proxy_url is outside supported values")
			}
		}
		if network.OutboundCidrs != nil {
			for _, value := range network.OutboundCidrs.Entries {
				if _, err := netip.ParsePrefix(value); err != nil {
					return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkOutboundCidrsEntries, "network.outbound_cidrs.entries is outside supported values")
				}
			}
		}
		if network.OutboundDomains != nil {
			for _, value := range network.OutboundDomains.Entries {
				domain := strings.TrimPrefix(value, "*.")
				if !providerDomainPattern.MatchString(domain) || len(domain) > 253 {
					return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldNetworkOutboundDomainsEntries, "network.outbound_domains.entries is outside supported values")
				}
			}
		}
	}
	if request.Security != nil {
		for _, item := range request.Security.Secrets {
			if item == nil {
				return sandbox.NewError(&sandbox.ErrorInfo{
					Kind:      sandbox.ErrorKindInvalidArgument,
					Provider:  "daytona",
					Operation: "create",
					Field:     sandbox.CreateFieldSecuritySecretsEnvironmentVariable,
					Message:   "sandbox-kit: nil configuration entry",
				}, nil)
			}
			if item.EnvironmentVariable != nil && !providerEnvNamePattern.MatchString(*item.EnvironmentVariable) {
				return providerCheckError(sandbox.ErrorKindInvalidArgument, sandbox.CreateFieldSecuritySecretsEnvironmentVariable, "security.secrets.environment_variable is outside supported values")
			}
		}
	}
	if lifetime := request.Lifetime; lifetime != nil {
		if delayedPolicy(lifetime.IdleStop) && delayedPolicy(lifetime.IdlePause) {
			return incompatibleLifetime()
		}
		if delayedPolicy(lifetime.IdlePause) && lifetime.GetEphemeral() {
			return incompatibleLifetime()
		}
		if delayedPolicy(lifetime.IdlePause) && immediatePolicy(lifetime.StoppedDelete) {
			return incompatibleLifetime()
		}
		if delayedPolicy(lifetime.StoppedDelete) && lifetime.GetEphemeral() {
			return incompatibleLifetime()
		}
		if delayedPolicy(lifetime.StoppedArchive) && lifetime.GetEphemeral() {
			return incompatibleLifetime()
		}
	}
	return nil
}

func delayedPolicy(action *sandbox.AutomaticAction) bool {
	return action != nil && action.GetMode() == sandbox.PolicyModeAfter && action.After != nil && action.After.AsDuration() > 0
}

func immediatePolicy(action *sandbox.AutomaticAction) bool {
	return action != nil && action.GetMode() == sandbox.PolicyModeAfter && action.After != nil && action.After.AsDuration() == 0
}

func incompatibleLifetime() error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      sandbox.ErrorKindInvalidArgument,
		Provider:  "daytona",
		Operation: "create",
		Field:     sandbox.CreateFieldLifetime,
		Message:   "sandbox-kit: incompatible lifetime policies",
	}, nil)
}

func providerCheckError(kind sandbox.ErrorKind, field, message string) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      kind,
		Provider:  "daytona",
		Operation: "create",
		Field:     field,
		Message:   "sandbox-kit daytona: " + message,
	}, nil)
}
