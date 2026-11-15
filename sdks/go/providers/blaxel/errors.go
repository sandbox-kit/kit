package blaxel

import (
	"context"
	"errors"
	"net/http"

	sdkblaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

var providerCodeKinds = map[string]sandbox.ErrorKind{
	sdkblaxel.ErrAuthenticationRequired: sandbox.ErrorKindAuthentication,
	sdkblaxel.ErrAuthenticationFailed:   sandbox.ErrorKindAuthentication,
	"TOKEN_REVOKED":                     sandbox.ErrorKindAuthentication,

	sdkblaxel.ErrForbidden:   sandbox.ErrorKindPermissionDenied,
	"UNAUTHORIZED":           sandbox.ErrorKindPermissionDenied,
	"UNAUTHORIZED_WORKSPACE": sandbox.ErrorKindPermissionDenied,

	sdkblaxel.ErrRouteNotFound:     sandbox.ErrorKindNotFound,
	sdkblaxel.ErrWorkloadNotFound:  sandbox.ErrorKindNotFound,
	sdkblaxel.ErrWorkspaceNotFound: sandbox.ErrorKindNotFound,
	"IMAGE_NOT_FOUND":              sandbox.ErrorKindNotFound,
	"VOLUME_NOT_FOUND":             sandbox.ErrorKindNotFound,
	"ACCOUNT_NOT_FOUND":            sandbox.ErrorKindNotFound,
	"SANDBOX_DELETED":              sandbox.ErrorKindNotFound,
	"VOLUME_DELETED":               sandbox.ErrorKindNotFound,

	"SANDBOX_ALREADY_EXISTS":  sandbox.ErrorKindAlreadyExists,
	"VOLUME_ALREADY_ATTACHED": sandbox.ErrorKindAlreadyExists,

	sdkblaxel.ErrBadRequest:      sandbox.ErrorKindInvalidArgument,
	sdkblaxel.ErrPolicyViolation: sandbox.ErrorKindInvalidArgument,
	"INVALID_INPUT":              sandbox.ErrorKindInvalidArgument,
	"INVALID_IMAGE":              sandbox.ErrorKindInvalidArgument,
	"INVALID_PORTS":              sandbox.ErrorKindInvalidArgument,
	"INVALID_REGION":             sandbox.ErrorKindInvalidArgument,
	"INVALID_VOLUME":             sandbox.ErrorKindInvalidArgument,
	"INVALID_VOLUME_COUNT":       sandbox.ErrorKindInvalidArgument,
	"INVALID_VOLUME_REGION":      sandbox.ErrorKindInvalidArgument,
	"VALIDATION_FAILED":          sandbox.ErrorKindInvalidArgument,
	"VALIDATION_ERROR":           sandbox.ErrorKindInvalidArgument,
	"WORKSPACE_REQUIRED":         sandbox.ErrorKindInvalidArgument,
	"PARSE_FAILED":               sandbox.ErrorKindInvalidArgument,
	"REGION_MISMATCH":            sandbox.ErrorKindInvalidArgument,
	"VOLUME_PROVIDER_MISMATCH":   sandbox.ErrorKindInvalidArgument,
	"KERNEL_GENERATION_MISMATCH": sandbox.ErrorKindInvalidArgument,
	"VALIDATED_DATA_ERROR":       sandbox.ErrorKindInvalidArgument,
	"WORKSPACE_DATA_MISSING":     sandbox.ErrorKindInvalidArgument,

	"UPLOAD_NOT_SUPPORTED":   sandbox.ErrorKindUnsupported,
	"UNSUPPORTED_GENERATION": sandbox.ErrorKindUnsupported,

	"CREATION_TIMEOUT":         sandbox.ErrorKindTimeout,
	"UPSTREAM_CONNECT_TIMEOUT": sandbox.ErrorKindTimeout,
	"UPSTREAM_TIMEOUT":         sandbox.ErrorKindTimeout,

	sdkblaxel.ErrUsageLimitExceeded: sandbox.ErrorKindResourceExhausted,
	"QUOTA_EXCEEDED":                sandbox.ErrorKindResourceExhausted,

	"RATE_LIMIT_EXCEEDED":           sandbox.ErrorKindRateLimited,
	"WORKSPACE_RATE_LIMIT_EXCEEDED": sandbox.ErrorKindRateLimited,

	"SANDBOX_DELETION_IN_PROGRESS": sandbox.ErrorKindConflict,

	sdkblaxel.ErrWorkloadUnavailable: sandbox.ErrorKindUnavailable,
	"ROUTING_UNAVAILABLE":            sandbox.ErrorKindUnavailable,
	"SERVICE_UNAVAILABLE":            sandbox.ErrorKindUnavailable,
	"AUTHENTICATION_UNAVAILABLE":     sandbox.ErrorKindUnavailable,
	"UPSTREAM_CONNECT_FAILED":        sandbox.ErrorKindUnavailable,
	"UPSTREAM_ERROR":                 sandbox.ErrorKindUnavailable,
	"GATEWAY_INTERNAL_ERROR":         sandbox.ErrorKindUnavailable,
	"UPLOAD_UNAVAILABLE":             sandbox.ErrorKindUnavailable,
	"LOCK_ACQUISITION_FAILED":        sandbox.ErrorKindUnavailable,
	"VOLUME_LOOKUP_FAILED":           sandbox.ErrorKindUnavailable,
	"CLUSTER_GATEWAY_ERROR":          sandbox.ErrorKindUnavailable,

	"WORKLOAD_FAILED":      sandbox.ErrorKindInternal,
	"INTERNAL_ERROR":       sandbox.ErrorKindInternal,
	"HANDLER_ERROR":        sandbox.ErrorKindInternal,
	"DATABASE_ERROR":       sandbox.ErrorKindInternal,
	"READ_FAILED":          sandbox.ErrorKindInternal,
	"BUILD_TRIGGER_ERROR":  sandbox.ErrorKindInternal,
	"PRESIGNED_URL_ERROR":  sandbox.ErrorKindInternal,
	"DEPLOYMENT_FAILED":    sandbox.ErrorKindInternal,
	"PROXY_CONFIG_ERROR":   sandbox.ErrorKindInternal,
	"CALLBACK_SEND_FAILED": sandbox.ErrorKindInternal,
	"FORK_FAILED":          sandbox.ErrorKindInternal,
	"LOCKDOWN_ERROR":       sandbox.ErrorKindInternal,
	"UNKNOWN_ERROR":        sandbox.ErrorKindInternal,
}

var providerStatusKinds = map[int]sandbox.ErrorKind{
	400: sandbox.ErrorKindInvalidArgument,
	422: sandbox.ErrorKindInvalidArgument,
	401: sandbox.ErrorKindAuthentication,
	402: sandbox.ErrorKindResourceExhausted,
	403: sandbox.ErrorKindPermissionDenied,
	404: sandbox.ErrorKindNotFound,
	410: sandbox.ErrorKindNotFound,
	408: sandbox.ErrorKindTimeout,
	// Bare 409 is Conflict for non-create ops. Create maps it to AlreadyExists
	// in statusKindFor (name collision). Coded conflicts still win above.
	409: sandbox.ErrorKindConflict,
	412: sandbox.ErrorKindFailedPrecondition,
	429: sandbox.ErrorKindRateLimited,
	500: sandbox.ErrorKindInternal,
	502: sandbox.ErrorKindUnavailable,
	503: sandbox.ErrorKindUnavailable,
	504: sandbox.ErrorKindTimeout,
}

// statusKindFor maps an HTTP status. Bare 409 on create is treated as a name
// collision (AlreadyExists); other operations keep Conflict so portable code
// that branches on kind matches Daytona's conflict semantics outside create.
func statusKindFor(status int, operation string) (sandbox.ErrorKind, bool) {
	if status == http.StatusConflict && operation == "create" {
		return sandbox.ErrorKindAlreadyExists, true
	}
	kind, ok := providerStatusKinds[status]
	return kind, ok
}

func mapProviderError(err error, operation string) error {
	if err == nil {
		return nil
	}
	var local *sandbox.Error
	if errors.As(err, &local) {
		return sdk.WithErrorContext(err, "blaxel", operation)
	}

	kind := sandbox.ErrorKindUnknown
	var statusCode *uint32
	var providerCode, providerSource *string
	switch {
	case errors.Is(err, context.Canceled):
		kind = sandbox.ErrorKindCanceled
	case errors.Is(err, context.DeadlineExceeded):
		kind = sandbox.ErrorKindTimeout
	}

	var native *sdkblaxel.Error
	if errors.As(err, &native) {
		if native.StatusCode > 0 {
			statusCode = sandbox.Value(uint32(native.StatusCode))
		}
		if native.ErrorCode != "" {
			providerCode = sandbox.Value(native.ErrorCode)
		}
		switch {
		case native.BlaxelSource != "":
			providerSource = sandbox.Value(native.BlaxelSource)
		case native.Origin != "":
			providerSource = sandbox.Value(native.Origin)
		}
		if kind == sandbox.ErrorKindUnknown {
			if mapped, ok := providerCodeKinds[native.ErrorCode]; ok {
				kind = mapped
			} else if mapped, ok := statusKindFor(native.StatusCode, operation); ok {
				kind = mapped
			}
		}
	}

	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:           kind,
		Provider:       "blaxel",
		Operation:      operation,
		Message:        err.Error(),
		StatusCode:     statusCode,
		ProviderCode:   providerCode,
		ProviderSource: providerSource,
	}, err)
}

// kitFault is the internal shape for Blaxel-scoped validation / response errors.
type kitFault struct {
	Kind      sandbox.ErrorKind
	Operation string
	Field     string
	Message   string
}

func invalidResponse(operation, field, message string) error {
	return kitError(kitFault{
		Kind:      sandbox.ErrorKindInvalidResponse,
		Operation: operation,
		Field:     field,
		Message:   message,
	})
}

func kitError(fault kitFault) error {
	return sandbox.NewError(&sandbox.ErrorInfo{
		Kind:      fault.Kind,
		Provider:  "blaxel",
		Operation: fault.Operation,
		Field:     fault.Field,
		Message:   "sandbox-kit blaxel: " + fault.Message,
	}, nil)
}
