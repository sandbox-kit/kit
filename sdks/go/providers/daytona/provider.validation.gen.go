// Code generated from provider validation specs. DO NOT EDIT.
package daytona

import (
	fmt "fmt"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
	netip "net/netip"
	url "net/url"
	regexp "regexp"
	strings "strings"
)

var providerEnvNamePattern = regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_]*$")
var providerDomainPattern = regexp.MustCompile("(?i)^(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)*[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.?$")

func validateProviderCreate(request *sandbox.CreateOptions) error {
	if request == nil {
		return nil
	}
	if request.Runtime != nil {
		if request.Runtime.Language != nil {
			if *request.Runtime.Language != "python" && *request.Runtime.Language != "javascript" && *request.Runtime.Language != "typescript" {
				return fmt.Errorf("sandbox-kit daytona: runtime.language is outside supported values")
			}
		}
	}
	if request.Network != nil {
		if request.Network.OutboundProxyURL != nil {
			parsed, err := url.Parse(*request.Network.OutboundProxyURL)
			if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
				return fmt.Errorf("sandbox-kit daytona: network.outbound_proxy_url is outside supported values")
			}
		}
	}
	if request.Network != nil {
		if request.Network.OutboundCIDRs != nil {
			for _, value := range request.Network.OutboundCIDRs.Entries {
				_, err := netip.ParsePrefix(value)
				if err != nil {
					return fmt.Errorf("sandbox-kit daytona: network.outbound_cidrs.entries is outside supported values")
				}
			}
		}
	}
	if request.Network != nil {
		if request.Network.OutboundDomains != nil {
			for _, value := range request.Network.OutboundDomains.Entries {
				domain := strings.TrimPrefix(value, "*.")
				matched := providerDomainPattern.MatchString(domain)
				if !matched || len(domain) > 253 {
					return fmt.Errorf("sandbox-kit daytona: network.outbound_domains.entries is outside supported values")
				}
			}
		}
	}
	if request.Security != nil {
		for _, item := range request.Security.Secrets {
			if item == nil {
				return fmt.Errorf("sandbox-kit: nil configuration entry")
			}
			if item.EnvironmentVariable != nil {
				matched := providerEnvNamePattern.MatchString(*item.EnvironmentVariable)
				if !matched {
					return fmt.Errorf("sandbox-kit daytona: security.secrets.environment_variable is outside supported values")
				}
			}
		}
	}
	if lifetime := request.Lifetime; lifetime != nil {
		if (lifetime.IdleStop != nil && lifetime.IdleStop.Mode == 2 && lifetime.IdleStop.After != nil && *lifetime.IdleStop.After > 0) && (lifetime.IdlePause != nil && lifetime.IdlePause.Mode == 2 && lifetime.IdlePause.After != nil && *lifetime.IdlePause.After > 0) {
			return fmt.Errorf("sandbox-kit: incompatible lifetime policies")
		}
	}
	if lifetime := request.Lifetime; lifetime != nil {
		if (lifetime.IdlePause != nil && lifetime.IdlePause.Mode == 2 && lifetime.IdlePause.After != nil && *lifetime.IdlePause.After > 0) && (lifetime.GetEphemeral()) {
			return fmt.Errorf("sandbox-kit: incompatible lifetime policies")
		}
	}
	if lifetime := request.Lifetime; lifetime != nil {
		if (lifetime.IdlePause != nil && lifetime.IdlePause.Mode == 2 && lifetime.IdlePause.After != nil && *lifetime.IdlePause.After > 0) && (lifetime.StoppedDelete != nil && lifetime.StoppedDelete.Mode == 2 && lifetime.StoppedDelete.After != nil && *lifetime.StoppedDelete.After == 0) {
			return fmt.Errorf("sandbox-kit: incompatible lifetime policies")
		}
	}
	if lifetime := request.Lifetime; lifetime != nil {
		if (lifetime.StoppedDelete != nil && lifetime.StoppedDelete.Mode == 2 && lifetime.StoppedDelete.After != nil && *lifetime.StoppedDelete.After > 0) && (lifetime.GetEphemeral()) {
			return fmt.Errorf("sandbox-kit: incompatible lifetime policies")
		}
	}
	if lifetime := request.Lifetime; lifetime != nil {
		if (lifetime.StoppedArchive != nil && lifetime.StoppedArchive.Mode == 2 && lifetime.StoppedArchive.After != nil && *lifetime.StoppedArchive.After > 0) && (lifetime.GetEphemeral()) {
			return fmt.Errorf("sandbox-kit: incompatible lifetime policies")
		}
	}
	return nil
}
