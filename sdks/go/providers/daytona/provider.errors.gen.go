// Code generated from provider error bindings. DO NOT EDIT.
package daytona

import (
	errors "errors"
	errors1 "github.com/daytona/clients/sdk-go/pkg/errors"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func mapProviderError(err error, operation string) error {
	if err == nil {
		return nil
	}
	var local *sandbox.Error
	if errors.As(err, &local) {
		return sandbox.WithErrorContext(err, "daytona", operation)
	}
	kind := sandbox.ErrorKindUnknown
	var statusCode *uint32
	var providerCode, providerSource *string
	{
		var native *errors1.DaytonaError
		if errors.As(err, &native) {
			if native.StatusCode > 0 {
				statusCode = sandbox.Value(uint32(native.StatusCode))
			}
			if native.Code != "" {
				providerCode = sandbox.Value(native.Code)
			}
			if native.Source != "" {
				providerSource = sandbox.Value(native.Source)
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
				case 412:
					kind = sandbox.ErrorKindFailedPrecondition
				case 422:
					kind = sandbox.ErrorKindInvalidArgument
				case 429:
					kind = sandbox.ErrorKindRateLimited
				case 500:
					kind = sandbox.ErrorKindInternal
				case 502:
					kind = sandbox.ErrorKindUnavailable
				case 503:
					kind = sandbox.ErrorKindUnavailable
				case 504:
					kind = sandbox.ErrorKindTimeout
				}
			}
		}
	}
	return sandbox.NewError(sandbox.ErrorInfo{Kind: kind, Provider: "daytona", Operation: operation, Message: err.Error(), StatusCode: statusCode, ProviderCode: providerCode, ProviderSource: providerSource}, err)
}
