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
	context "context"
	errors "errors"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
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
type Error struct {
	// Kind is the shared category. Switch on it.
	Kind ErrorKind
	// Provider is the provider identifier, such as "daytona" or "modal". Empty
	// when the failure is local.
	Provider string
	// Operation is the Kit operation, such as "create" or "initialize".
	Operation string
	// Field is the canonical field path, such as "resources.cpu_cores". Empty
	// when no field applies.
	Field string
	// Message is the human-readable explanation.
	Message string
	// StatusCode is the native HTTP or transport status when the provider
	// reported one.
	StatusCode *uint32
	// ProviderCode is the provider's own error code when one was reported.
	ProviderCode *string
	// ProviderSource is the provider component that reported the error, when
	// available.
	ProviderSource *string
	// Rule is a stable validation or policy identifier, when applicable.
	Rule *string
	// cause is the original SDK or validator error. Unwrap returns it. It is not
	// serialized.
	cause error
}

// Error returns the portable message. A field path is appended when Field is
// set. A nil error returns an empty string.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Field != "" {
		return e.Message + " (" + e.Field + ")"
	}
	return e.Message
}

// Unwrap returns the native cause so errors.Is and errors.As can match it. A
// nil error returns nil.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

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
	if info == nil {
		info = &ErrorInfo{}
	}
	kind := info.Kind
	if kind.Descriptor().Values().ByNumber(kind.Number()) == nil {
		kind = ErrorKindUnknown
	}
	if errors.Is(cause, context.Canceled) {
		kind = ErrorKindCanceled
	} else if errors.Is(cause, context.DeadlineExceeded) {
		kind = ErrorKindTimeout
	}
	return &Error{
		Kind:           kind,
		Provider:       info.Provider,
		Operation:      info.Operation,
		Field:          info.Field,
		Message:        info.Message,
		StatusCode:     copyOptional(info.StatusCode),
		ProviderCode:   copyOptional(info.ProviderCode),
		ProviderSource: copyOptional(info.ProviderSource),
		Rule:           copyOptional(info.Rule),
		cause:          cause,
	}
}

func (e *Error) clone() *Error {
	if e == nil {
		return &Error{}
	}
	return &Error{
		Kind:           e.Kind,
		Provider:       e.Provider,
		Operation:      e.Operation,
		Field:          e.Field,
		Message:        e.Message,
		StatusCode:     copyOptional(e.StatusCode),
		ProviderCode:   copyOptional(e.ProviderCode),
		ProviderSource: copyOptional(e.ProviderSource),
		Rule:           copyOptional(e.Rule),
		cause:          e.cause,
	}
}

func copyOptional[T any](value *T) *T {
	if value == nil {
		return nil
	}
	return Value(*value)
}

// GetStatusCode returns the native status, or 0 when it was not reported.
func (e *Error) GetStatusCode() uint32 {
	if e == nil || e.StatusCode == nil {
		return 0
	}
	return *e.StatusCode
}

// GetProviderCode returns the provider error code, or "" when it was not reported.
func (e *Error) GetProviderCode() string {
	if e == nil || e.ProviderCode == nil {
		return ""
	}
	return *e.ProviderCode
}

// GetProviderSource returns the provider component, or "" when it was not reported.
func (e *Error) GetProviderSource() string {
	if e == nil || e.ProviderSource == nil {
		return ""
	}
	return *e.ProviderSource
}

// GetRule returns the stable rule identifier, or "" when it was not reported.
func (e *Error) GetRule() string {
	if e == nil || e.Rule == nil {
		return ""
	}
	return *e.Rule
}

// WithErrorContext adds provider and operation details without changing an
// error that already has them.
// A nil error returns nil.
func WithErrorContext(err error, provider string, operation string) error {
	if err == nil {
		return nil
	}
	var existing *Error
	if errors.As(err, &existing) && existing != nil {
		if (provider == "" || existing.Provider == provider) && (operation == "" || existing.Operation == operation) {
			return err
		}
		cloned := existing.clone()
		cloned.cause = err
		if provider != "" {
			cloned.Provider = provider
		}
		if operation != "" {
			cloned.Operation = operation
		}
		return cloned
	}
	return NewError(&ErrorInfo{Kind: ErrorKindUnknown, Provider: provider, Operation: operation, Message: err.Error()}, err)
}
