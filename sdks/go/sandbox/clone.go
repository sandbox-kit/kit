package sandbox

import (
	"google.golang.org/protobuf/proto"
)

// cloneCreateOptions returns an owned copy of request that preserves empty maps
// and slices. proto.Clone alone drops empty collections, which would erase
// explicit presence the contract treats as meaningful.
func cloneCreateOptions(request *CreateOptions) *CreateOptions {
	if request == nil {
		return nil
	}
	out := proto.Clone(request).(*CreateOptions)
	preserveCreateOptionsCollections(request, out)
	return out
}

// cloneSandboxInfo returns an owned copy of info that preserves empty maps.
func cloneSandboxInfo(info *SandboxInfo) *SandboxInfo {
	if info == nil {
		return nil
	}
	out := proto.Clone(info).(*SandboxInfo)
	preserveSandboxInfoCollections(info, out)
	return out
}

func preserveCreateOptionsCollections(src, dst *CreateOptions) {
	if src.Environment != nil && dst.Environment == nil {
		dst.Environment = map[string]string{}
	}
	if src.Labels != nil && dst.Labels == nil {
		dst.Labels = map[string]string{}
	}
	if src.Runtime != nil && dst.Runtime != nil {
		if src.Runtime.Entrypoint != nil && dst.Runtime.Entrypoint == nil {
			dst.Runtime.Entrypoint = []string{}
		}
	}
	if src.Placement != nil && dst.Placement != nil {
		if src.Placement.Regions != nil && dst.Placement.Regions == nil {
			dst.Placement.Regions = []string{}
		}
	}
	if src.Network != nil && dst.Network != nil {
		preserveStringAllowlist(src.Network.OutboundCidrs, &dst.Network.OutboundCidrs)
		preserveStringAllowlist(src.Network.OutboundDomains, &dst.Network.OutboundDomains)
		preserveStringAllowlist(src.Network.InboundCidrs, &dst.Network.InboundCidrs)
		if src.Network.Ports != nil && dst.Network.Ports == nil {
			dst.Network.Ports = []*PortBinding{}
		}
	}
	if src.Security != nil && dst.Security != nil {
		if src.Security.Secrets != nil && dst.Security.Secrets == nil {
			dst.Security.Secrets = []*SecretReference{}
		}
	}
	if src.Readiness != nil && dst.Readiness != nil {
		if cmd := src.Readiness.GetCommand(); cmd != nil {
			if dstCmd := dst.Readiness.GetCommand(); dstCmd != nil && cmd.Argv != nil && dstCmd.Argv == nil {
				dstCmd.Argv = []string{}
			}
		}
	}
	if src.ProviderOptions != nil && dst.ProviderOptions != nil {
		preserveMetadataObject(src.ProviderOptions, dst.ProviderOptions)
	}
}

func preserveSandboxInfoCollections(src, dst *SandboxInfo) {
	if src.Labels != nil && dst.Labels == nil {
		dst.Labels = map[string]string{}
	}
	if src.Origins != nil && dst.Origins == nil {
		dst.Origins = map[string]ValueOrigin{}
	}
	if src.ProviderMetadata != nil && dst.ProviderMetadata != nil {
		preserveMetadataObject(src.ProviderMetadata, dst.ProviderMetadata)
	}
}

func preserveStringAllowlist(src *StringAllowlist, dst **StringAllowlist) {
	if src == nil || *dst == nil {
		return
	}
	if src.Entries != nil && (*dst).Entries == nil {
		(*dst).Entries = []string{}
	}
}

func preserveMetadataObject(src, dst *MetadataObject) {
	if src.Values != nil && dst.Values == nil {
		dst.Values = map[string]*MetadataValue{}
	}
}
