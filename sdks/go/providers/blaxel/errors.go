package blaxel

import (
	"context"
	"errors"

	sdkblaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

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
	if errors.Is(err, context.Canceled) {
		kind = sandbox.ErrorKindCanceled
	} else if errors.Is(err, context.DeadlineExceeded) {
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
		if native.BlaxelSource != "" {
			providerSource = sandbox.Value(native.BlaxelSource)
		} else if native.Origin != "" {
			providerSource = sandbox.Value(native.Origin)
		}
		if kind == sandbox.ErrorKindUnknown {
			switch native.ErrorCode {
			case sdkblaxel.ErrAuthenticationRequired, sdkblaxel.ErrAuthenticationFailed, "TOKEN_REVOKED":
				kind = sandbox.ErrorKindAuthentication
			case sdkblaxel.ErrForbidden, "UNAUTHORIZED", "UNAUTHORIZED_WORKSPACE":
				kind = sandbox.ErrorKindPermissionDenied
			case sdkblaxel.ErrRouteNotFound, sdkblaxel.ErrWorkloadNotFound, sdkblaxel.ErrWorkspaceNotFound,
				"IMAGE_NOT_FOUND", "VOLUME_NOT_FOUND", "ACCOUNT_NOT_FOUND", "SANDBOX_DELETED", "VOLUME_DELETED":
				kind = sandbox.ErrorKindNotFound
			case "SANDBOX_ALREADY_EXISTS", "VOLUME_ALREADY_ATTACHED":
				// Control-plane codes (not all are gateway apierror constants).
				kind = sandbox.ErrorKindAlreadyExists
			case sdkblaxel.ErrBadRequest, sdkblaxel.ErrPolicyViolation,
				"INVALID_INPUT", "INVALID_IMAGE", "INVALID_PORTS", "INVALID_REGION",
				"INVALID_VOLUME", "INVALID_VOLUME_COUNT", "INVALID_VOLUME_REGION",
				"VALIDATION_FAILED", "VALIDATION_ERROR", "WORKSPACE_REQUIRED":
				kind = sandbox.ErrorKindInvalidArgument
			case "CREATION_TIMEOUT", "UPSTREAM_CONNECT_TIMEOUT", "UPSTREAM_TIMEOUT":
				kind = sandbox.ErrorKindTimeout
			case sdkblaxel.ErrUsageLimitExceeded, "QUOTA_EXCEEDED":
				kind = sandbox.ErrorKindResourceExhausted
			case "RATE_LIMIT_EXCEEDED", "WORKSPACE_RATE_LIMIT_EXCEEDED":
				kind = sandbox.ErrorKindRateLimited
			case "SANDBOX_DELETION_IN_PROGRESS":
				kind = sandbox.ErrorKindConflict
			case sdkblaxel.ErrWorkloadUnavailable, "ROUTING_UNAVAILABLE", "SERVICE_UNAVAILABLE",
				"AUTHENTICATION_UNAVAILABLE", "UPSTREAM_CONNECT_FAILED", "UPSTREAM_ERROR",
				"GATEWAY_INTERNAL_ERROR":
				kind = sandbox.ErrorKindUnavailable
			case "WORKLOAD_FAILED", "INTERNAL_ERROR", "HANDLER_ERROR", "DATABASE_ERROR":
				kind = sandbox.ErrorKindInternal
			}
		}
		if kind == sandbox.ErrorKindUnknown {
			switch native.StatusCode {
			case 400, 422:
				kind = sandbox.ErrorKindInvalidArgument
			case 401:
				kind = sandbox.ErrorKindAuthentication
			case 402:
				kind = sandbox.ErrorKindResourceExhausted
			case 403:
				kind = sandbox.ErrorKindPermissionDenied
			case 404, 410:
				kind = sandbox.ErrorKindNotFound
			case 408:
				kind = sandbox.ErrorKindTimeout
			case 409:
				kind = sandbox.ErrorKindConflict
			case 412:
				kind = sandbox.ErrorKindFailedPrecondition
			case 429:
				kind = sandbox.ErrorKindRateLimited
			case 500:
				kind = sandbox.ErrorKindInternal
			case 502, 503:
				kind = sandbox.ErrorKindUnavailable
			case 504:
				kind = sandbox.ErrorKindTimeout
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
