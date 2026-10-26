// Portable error behavior is handwritten. ErrorInfo and ErrorKind remain
// protobuf-owned contract data during the schema migration.
// Error carries portable details and the original cause. Use errors.As, then
// switch on Kind.
//
// Example:
//
//	var detail *Error
//	if errors.As(err, &detail) {
//	    fmt.Println(detail.Kind, detail.Field)
//	}
package sandbox

import (
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
)

// ErrorInfo and ErrorKind are the generated portable protobuf contracts.
// Error below only adds native cause chaining.
type ErrorInfo = sandboxv1.ErrorInfo
type ErrorKind = sandboxv1.ErrorKind

const (
	ErrorKindUnknown            = sandboxv1.ErrorKind_ERROR_KIND_UNKNOWN
	ErrorKindInvalidArgument    = sandboxv1.ErrorKind_ERROR_KIND_INVALID_ARGUMENT
	ErrorKindUnsupported        = sandboxv1.ErrorKind_ERROR_KIND_UNSUPPORTED
	ErrorKindAuthentication     = sandboxv1.ErrorKind_ERROR_KIND_AUTHENTICATION
	ErrorKindPermissionDenied   = sandboxv1.ErrorKind_ERROR_KIND_PERMISSION_DENIED
	ErrorKindNotFound           = sandboxv1.ErrorKind_ERROR_KIND_NOT_FOUND
	ErrorKindAlreadyExists      = sandboxv1.ErrorKind_ERROR_KIND_ALREADY_EXISTS
	ErrorKindConflict           = sandboxv1.ErrorKind_ERROR_KIND_CONFLICT
	ErrorKindRateLimited        = sandboxv1.ErrorKind_ERROR_KIND_RATE_LIMITED
	ErrorKindResourceExhausted  = sandboxv1.ErrorKind_ERROR_KIND_RESOURCE_EXHAUSTED
	ErrorKindTimeout            = sandboxv1.ErrorKind_ERROR_KIND_TIMEOUT
	ErrorKindCanceled           = sandboxv1.ErrorKind_ERROR_KIND_CANCELED
	ErrorKindUnavailable        = sandboxv1.ErrorKind_ERROR_KIND_UNAVAILABLE
	ErrorKindInvalidResponse    = sandboxv1.ErrorKind_ERROR_KIND_INVALID_RESPONSE
	ErrorKindInternal           = sandboxv1.ErrorKind_ERROR_KIND_INTERNAL
	ErrorKindFailedPrecondition = sandboxv1.ErrorKind_ERROR_KIND_FAILED_PRECONDITION
)

// Error exposes portable details and retains the native cause for errors.Is
// and errors.As. Details are plain fields copied from ErrorInfo. The error
// value is not a protobuf message.
//
// Example:
//
//	var sandboxErr *Error
//	if errors.As(err, &sandboxErr) {
//	    fmt.Println(sandboxErr.Kind, sandboxErr.Field)
//	}
type Error = sdk.Error

// NewError copies portable details and retains the original cause. An unknown
// kind becomes ErrorKindUnknown.
// Cancellation and deadline causes override the supplied kind.
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
