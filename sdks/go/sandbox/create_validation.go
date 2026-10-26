package sandbox

import (
	"errors"
	"strings"

	"buf.build/go/protovalidate"
)

// stableRuleIDs maps Protovalidate rule ids onto the portable rule inventory
// names used in ErrorInfo.Rule and specs/validation.yaml.
var stableRuleIDs = map[string]string{
	"enum.defined_only":             "known_enum",
	"duration.gte":                  "minimum",
	"oneof.required":                "exactly_one",
	"auth_credentials_required":     "exactly_one",
	"cpu_limit":                     "limit",
	"memory_limit":                  "limit",
	"cpu_requires_positive":         "requires_positive",
	"memory_requires_positive":      "requires_positive",
	"restricted_requires_allowlist": "requires_any",
	"block_all_forbids_network":     "forbids",
	"allow_all_forbids_allowlists":  "forbids",
	"after_requires_duration":       "requires_any",
	"after_forbids_duration":        "forbids",
}

// ValidateCreateOptions rejects CreateOptions values the shared contract cannot
// represent. A nil request is treated as omitted configuration. Account quotas
// stay with the provider.
//
// Portable rules are declared as Protovalidate options on kit.sandbox.v1 messages
// and evaluated here. Provider-specific checks remain in each provider adapter.
func ValidateCreateOptions(request *CreateOptions) error {
	if request == nil {
		return nil
	}
	if err := rejectNilCreateElements(request); err != nil {
		return err
	}
	if err := protovalidate.Validate(request); err != nil {
		return mapProtovalidateError(err)
	}
	return nil
}

// rejectNilCreateElements catches Go-only nil entries in repeated message fields.
// Wire decoding never produces them; callers constructing options in-process can.
func rejectNilCreateElements(request *CreateOptions) error {
	if network := request.Network; network != nil {
		for _, port := range network.Ports {
			if port == nil {
				return invalidCreate("network.ports", "required")
			}
		}
	}
	if security := request.Security; security != nil {
		for _, secret := range security.Secrets {
			if secret == nil {
				return invalidCreate("security.secrets", "required")
			}
		}
	}
	return nil
}

func mapProtovalidateError(err error) error {
	var violations *protovalidate.ValidationError
	if !errors.As(err, &violations) || len(violations.Violations) == 0 {
		return NewError(&ErrorInfo{
			Kind:    ErrorKindInvalidArgument,
			Message: "sandbox-kit: invalid configuration",
		}, err)
	}
	violation := violations.Violations[0]
	ruleID := violation.Proto.GetRuleId()
	field := protovalidateFieldPath(violation)
	rule := ruleID
	if mapped, ok := stableRuleIDs[ruleID]; ok {
		rule = mapped
	}
	if field == "" {
		field = protovalidateMessageFieldHint(ruleID)
	}
	return invalidCreate(field, rule)
}

func protovalidateFieldPath(violation *protovalidate.Violation) string {
	if violation == nil || violation.Proto == nil {
		return ""
	}
	path := violation.Proto.GetField()
	if path == nil {
		return ""
	}
	parts := make([]string, 0, len(path.GetElements()))
	for _, element := range path.GetElements() {
		name := element.GetFieldName()
		if name == "" {
			continue
		}
		parts = append(parts, name)
	}
	return strings.Join(parts, ".")
}

// protovalidateMessageFieldHint recovers inventory-style field paths for
// message-level CEL rules that do not attach a single field element.
func protovalidateMessageFieldHint(ruleID string) string {
	switch ruleID {
	case "cpu_limit", "cpu_requires_positive":
		return "resources.cpu_limit_cores"
	case "memory_limit", "memory_requires_positive":
		return "resources.memory_limit_mib"
	case "restricted_requires_allowlist":
		return "network.outbound_cidrs,outbound_domains"
	case "block_all_forbids_network":
		return "network.outbound_cidrs,outbound_domains,private_network,linked_sandbox,inbound_cidrs,ports"
	case "allow_all_forbids_allowlists":
		return "network.outbound_cidrs,outbound_domains"
	case "after_requires_duration", "after_forbids_duration":
		return "after"
	case "oneof.required":
		return "tcp_port,command"
	case "auth_credentials_required":
		return "api_key,token_pair,bearer_token,oauth_refresh"
	case "exactly_one":
		return "client_secret,jwt_key"
	default:
		return ""
	}
}

func invalidCreate(field, rule string) error {
	return NewError(&ErrorInfo{Kind: ErrorKindInvalidArgument, Field: field, Rule: Value(rule), Message: "sandbox-kit: invalid configuration"}, nil)
}
