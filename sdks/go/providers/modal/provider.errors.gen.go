// Code generated from provider error bindings. DO NOT EDIT.
package modal

import (
	errors "errors"
	_go "github.com/modal-labs/modal-client/go"
	sandbox "github.com/sandbox-kit/kit/sdks/go/sandbox"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func mapProviderError(err error, operation string) error {
	if err == nil {
		return nil
	}
	var local *sandbox.Error
	if errors.As(err, &local) {
		return sandbox.WithErrorContext(err, "modal", operation)
	}
	kind := sandbox.ErrorKindUnknown
	var statusCode *uint32
	var providerCode, providerSource *string
	if kind == sandbox.ErrorKindUnknown {
		var value0 _go.InvalidError
		var pointer0 *_go.InvalidError
		if errors.As(err, &value0) || errors.As(err, &pointer0) {
			kind = sandbox.ErrorKindInvalidArgument
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value1 _go.NotFoundError
		var pointer1 *_go.NotFoundError
		if errors.As(err, &value1) || errors.As(err, &pointer1) {
			kind = sandbox.ErrorKindNotFound
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value2 _go.AlreadyExistsError
		var pointer2 *_go.AlreadyExistsError
		if errors.As(err, &value2) || errors.As(err, &pointer2) {
			kind = sandbox.ErrorKindAlreadyExists
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value3 _go.ConflictError
		var pointer3 *_go.ConflictError
		if errors.As(err, &value3) || errors.As(err, &pointer3) {
			kind = sandbox.ErrorKindConflict
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value4 _go.ResourceExhaustedError
		var pointer4 *_go.ResourceExhaustedError
		if errors.As(err, &value4) || errors.As(err, &pointer4) {
			kind = sandbox.ErrorKindResourceExhausted
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value5 _go.TimeoutError
		var pointer5 *_go.TimeoutError
		if errors.As(err, &value5) || errors.As(err, &pointer5) {
			kind = sandbox.ErrorKindTimeout
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value6 _go.SandboxTimeoutError
		var pointer6 *_go.SandboxTimeoutError
		if errors.As(err, &value6) || errors.As(err, &pointer6) {
			kind = sandbox.ErrorKindTimeout
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value7 _go.ClientClosedError
		var pointer7 *_go.ClientClosedError
		if errors.As(err, &value7) || errors.As(err, &pointer7) {
			kind = sandbox.ErrorKindFailedPrecondition
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		switch status.Code(err) {
		case codes.Code(1):
			kind = sandbox.ErrorKindCanceled
		case codes.Code(3):
			kind = sandbox.ErrorKindInvalidArgument
		case codes.Code(4):
			kind = sandbox.ErrorKindTimeout
		case codes.Code(5):
			kind = sandbox.ErrorKindNotFound
		case codes.Code(6):
			kind = sandbox.ErrorKindAlreadyExists
		case codes.Code(7):
			kind = sandbox.ErrorKindPermissionDenied
		case codes.Code(8):
			kind = sandbox.ErrorKindResourceExhausted
		case codes.Code(9):
			kind = sandbox.ErrorKindFailedPrecondition
		case codes.Code(10):
			kind = sandbox.ErrorKindConflict
		case codes.Code(11):
			kind = sandbox.ErrorKindInvalidArgument
		case codes.Code(12):
			kind = sandbox.ErrorKindUnsupported
		case codes.Code(13):
			kind = sandbox.ErrorKindInternal
		case codes.Code(14):
			kind = sandbox.ErrorKindUnavailable
		case codes.Code(16):
			kind = sandbox.ErrorKindAuthentication
		}
	}
	return sandbox.NewError(sandbox.ErrorInfo{Kind: kind, Provider: "modal", Operation: operation, Message: err.Error(), StatusCode: statusCode, ProviderCode: providerCode, ProviderSource: providerSource}, err)
}
