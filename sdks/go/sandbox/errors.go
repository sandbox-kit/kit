package sandbox

import (
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
)

// ErrorInfo is the detail carried by Error.
// Kind is the category. Field is the setting to change. Message explains the failure.
type ErrorInfo = sandboxv1.ErrorInfo

// ErrorKind is the category of an Error, such as an authentication failure or an unsupported setting.
type ErrorKind = sandboxv1.ErrorKind

const (
	// ErrorKindUnknown means the failure does not match a known category.
	ErrorKindUnknown = sandboxv1.ErrorKind_ERROR_KIND_UNKNOWN
	// ErrorKindInvalidArgument means a setting was rejected before or during the call.
	// Field names the setting, such as "resources.cpu_cores".
	ErrorKindInvalidArgument = sandboxv1.ErrorKind_ERROR_KIND_INVALID_ARGUMENT
	// ErrorKindUnsupported means this provider does not accept the setting or operation.
	ErrorKindUnsupported = sandboxv1.ErrorKind_ERROR_KIND_UNSUPPORTED
	// ErrorKindAuthentication means the credentials were rejected.
	ErrorKindAuthentication = sandboxv1.ErrorKind_ERROR_KIND_AUTHENTICATION
	// ErrorKindPermissionDenied means the account may not perform the operation.
	ErrorKindPermissionDenied = sandboxv1.ErrorKind_ERROR_KIND_PERMISSION_DENIED
	// ErrorKindNotFound means the sandbox or another named resource does not exist.
	ErrorKindNotFound = sandboxv1.ErrorKind_ERROR_KIND_NOT_FOUND
	// ErrorKindAlreadyExists means the named sandbox or resource already exists.
	ErrorKindAlreadyExists = sandboxv1.ErrorKind_ERROR_KIND_ALREADY_EXISTS
	// ErrorKindConflict means the sandbox is in a state that rejects this operation.
	ErrorKindConflict = sandboxv1.ErrorKind_ERROR_KIND_CONFLICT
	// ErrorKindRateLimited means the provider asked the caller to slow down.
	ErrorKindRateLimited = sandboxv1.ErrorKind_ERROR_KIND_RATE_LIMITED
	// ErrorKindResourceExhausted means the account or provider is out of capacity.
	ErrorKindResourceExhausted = sandboxv1.ErrorKind_ERROR_KIND_RESOURCE_EXHAUSTED
	// ErrorKindTimeout means the operation exceeded its deadline.
	ErrorKindTimeout = sandboxv1.ErrorKind_ERROR_KIND_TIMEOUT
	// ErrorKindCanceled means the context was canceled.
	ErrorKindCanceled = sandboxv1.ErrorKind_ERROR_KIND_CANCELED
	// ErrorKindUnavailable means the provider could not be reached.
	ErrorKindUnavailable = sandboxv1.ErrorKind_ERROR_KIND_UNAVAILABLE
	// ErrorKindInvalidResponse means the provider returned a response this client cannot use.
	ErrorKindInvalidResponse = sandboxv1.ErrorKind_ERROR_KIND_INVALID_RESPONSE
	// ErrorKindInternal means the provider or client failed for a reason it did not classify.
	ErrorKindInternal = sandboxv1.ErrorKind_ERROR_KIND_INTERNAL
	// ErrorKindFailedPrecondition means the client or sandbox is not in a state that allows the call.
	ErrorKindFailedPrecondition = sandboxv1.ErrorKind_ERROR_KIND_FAILED_PRECONDITION
)

// Error is the failure type returned by the client and sandbox handle.
// Use errors.As to read Kind, Field, and Message. Unwrap returns the provider error when there is one.
//
// Example:
//
//	var sandboxErr *Error
//	if errors.As(err, &sandboxErr) {
//	    fmt.Println(sandboxErr.Kind, sandboxErr.Field)
//	}
type Error = sdk.Error

// NewError builds an error a caller can read with errors.As.
// Kind is the category, Field is the setting to change, and cause is the
// original error returned by Unwrap. A canceled or timed-out cause sets Kind.
//
// Example:
//
//	err := NewError(&ErrorInfo{
//	    Kind:    ErrorKindUnsupported,
//	    Field:   CreateFieldResourcesDiskMib,
//	    Message: "disk overrides are unsupported",
//	}, cause)
func NewError(info *ErrorInfo, cause error) *Error {
	return sdk.NewError(info, cause)
}
