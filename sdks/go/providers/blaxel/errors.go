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
			case sdkblaxel.ErrAuthenticationRequired, sdkblaxel.ErrAuthenticationFailed:
				kind = sandbox.ErrorKindAuthentication
			case sdkblaxel.ErrForbidden:
				kind = sandbox.ErrorKindPermissionDenied
			case sdkblaxel.ErrRouteNotFound, sdkblaxel.ErrWorkloadNotFound, sdkblaxel.ErrWorkspaceNotFound:
				kind = sandbox.ErrorKindNotFound
			case sdkblaxel.ErrWorkloadUnavailable:
				kind = sandbox.ErrorKindUnavailable
			case sdkblaxel.ErrUsageLimitExceeded:
				kind = sandbox.ErrorKindResourceExhausted
			case sdkblaxel.ErrBadRequest, sdkblaxel.ErrPolicyViolation:
				kind = sandbox.ErrorKindInvalidArgument
			case "SANDBOX_ALREADY_EXISTS", "VOLUME_ALREADY_ATTACHED":
				// Control-plane codes (not in gateway apierror constants).
				kind = sandbox.ErrorKindAlreadyExists
			}
		}
		if kind == sandbox.ErrorKindUnknown {
			switch native.StatusCode {
			case 400:
				kind = sandbox.ErrorKindInvalidArgument
			case 401:
				kind = sandbox.ErrorKindAuthentication
			case 403:
				kind = sandbox.ErrorKindPermissionDenied
			case 404:
				kind = sandbox.ErrorKindNotFound
			case 408:
				kind = sandbox.ErrorKindTimeout
			case 409:
				kind = sandbox.ErrorKindConflict
			case 410:
				kind = sandbox.ErrorKindNotFound
			case 412:
				kind = sandbox.ErrorKindFailedPrecondition
			case 422:
				kind = sandbox.ErrorKindInvalidArgument
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
