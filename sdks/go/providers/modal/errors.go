package modal

import (
	context "context"
	errors "errors"

	sdk "github.com/modal-labs/modal-client/go"
	sandboxsdk "github.com/sandbox-kit/kit/sdks/go/internal/sdk"
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
		return sandboxsdk.WithErrorContext(err, "modal", operation)
	}
	kind := sandbox.ErrorKindUnknown
	var statusCode *uint32
	var providerCode, providerSource *string
	if errors.Is(err, context.Canceled) {
		kind = sandbox.ErrorKindCanceled
	} else if errors.Is(err, context.DeadlineExceeded) {
		kind = sandbox.ErrorKindTimeout
	}
	if kind == sandbox.ErrorKindUnknown {
		var value0 sdk.InvalidError
		var pointer0 *sdk.InvalidError
		if errors.As(err, &value0) || errors.As(err, &pointer0) {
			kind = sandbox.ErrorKindInvalidArgument
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value1 sdk.NotFoundError
		var pointer1 *sdk.NotFoundError
		if errors.As(err, &value1) || errors.As(err, &pointer1) {
			kind = sandbox.ErrorKindNotFound
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value2 sdk.AlreadyExistsError
		var pointer2 *sdk.AlreadyExistsError
		if errors.As(err, &value2) || errors.As(err, &pointer2) {
			kind = sandbox.ErrorKindAlreadyExists
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value3 sdk.ConflictError
		var pointer3 *sdk.ConflictError
		if errors.As(err, &value3) || errors.As(err, &pointer3) {
			kind = sandbox.ErrorKindConflict
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value4 sdk.ResourceExhaustedError
		var pointer4 *sdk.ResourceExhaustedError
		if errors.As(err, &value4) || errors.As(err, &pointer4) {
			kind = sandbox.ErrorKindResourceExhausted
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value5 sdk.TimeoutError
		var pointer5 *sdk.TimeoutError
		if errors.As(err, &value5) || errors.As(err, &pointer5) {
			kind = sandbox.ErrorKindTimeout
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value6 sdk.SandboxTimeoutError
		var pointer6 *sdk.SandboxTimeoutError
		if errors.As(err, &value6) || errors.As(err, &pointer6) {
			kind = sandbox.ErrorKindTimeout
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var value7 sdk.ClientClosedError
		var pointer7 *sdk.ClientClosedError
		if errors.As(err, &value7) || errors.As(err, &pointer7) {
			kind = sandbox.ErrorKindFailedPrecondition
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var notFound sdk.SandboxFilesystemNotFoundError
		var notFoundPtr *sdk.SandboxFilesystemNotFoundError
		if errors.As(err, &notFound) || errors.As(err, &notFoundPtr) {
			kind = sandbox.ErrorKindNotFound
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var denied sdk.SandboxFilesystemPermissionError
		var deniedPtr *sdk.SandboxFilesystemPermissionError
		if errors.As(err, &denied) || errors.As(err, &deniedPtr) {
			kind = sandbox.ErrorKindPermissionDenied
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var exists sdk.SandboxFilesystemPathAlreadyExistsError
		var existsPtr *sdk.SandboxFilesystemPathAlreadyExistsError
		if errors.As(err, &exists) || errors.As(err, &existsPtr) {
			kind = sandbox.ErrorKindAlreadyExists
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var notEmpty sdk.SandboxFilesystemDirectoryNotEmptyError
		var notEmptyPtr *sdk.SandboxFilesystemDirectoryNotEmptyError
		if errors.As(err, &notEmpty) || errors.As(err, &notEmptyPtr) {
			kind = sandbox.ErrorKindFailedPrecondition
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var tooLarge sdk.SandboxFilesystemFileTooLargeError
		var tooLargePtr *sdk.SandboxFilesystemFileTooLargeError
		if errors.As(err, &tooLarge) || errors.As(err, &tooLargePtr) {
			kind = sandbox.ErrorKindResourceExhausted
		}
	}
	if kind == sandbox.ErrorKindUnknown {
		var notDir sdk.SandboxFilesystemNotADirectoryError
		var notDirPtr *sdk.SandboxFilesystemNotADirectoryError
		var isDir sdk.SandboxFilesystemIsADirectoryError
		var isDirPtr *sdk.SandboxFilesystemIsADirectoryError
		if errors.As(err, &notDir) || errors.As(err, &notDirPtr) || errors.As(err, &isDir) || errors.As(err, &isDirPtr) {
			kind = sandbox.ErrorKindInvalidArgument
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
	return sandbox.NewError(&sandbox.ErrorInfo{Kind: kind, Provider: "modal", Operation: operation, Message: err.Error(), StatusCode: statusCode, ProviderCode: providerCode, ProviderSource: providerSource}, err)
}
