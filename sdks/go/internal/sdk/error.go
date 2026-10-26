package sdk

import (
	"context"
	"errors"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
)

// Error exposes portable details and retains the native cause for errors.Is
// and errors.As. Details are plain fields copied from ErrorInfo. The error
// value is not a protobuf message.
type Error struct {
	// Kind is the shared category. Switch on it.
	Kind sandboxv1.ErrorKind
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

// NewError copies portable details and retains the original cause. An unknown
// kind becomes ErrorKindUnknown. Cancellation and deadline causes override the
// supplied kind.
func NewError(info *sandboxv1.ErrorInfo, cause error) *Error {
	if info == nil {
		info = &sandboxv1.ErrorInfo{}
	}
	kind := info.Kind
	if kind.Descriptor().Values().ByNumber(kind.Number()) == nil {
		kind = sandboxv1.ErrorKind_ERROR_KIND_UNKNOWN
	}
	if errors.Is(cause, context.Canceled) {
		kind = sandboxv1.ErrorKind_ERROR_KIND_CANCELED
	} else if errors.Is(cause, context.DeadlineExceeded) {
		kind = sandboxv1.ErrorKind_ERROR_KIND_TIMEOUT
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
	copied := *value
	return &copied
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
// error that already has them. A nil error returns nil.
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
	return NewError(&sandboxv1.ErrorInfo{Kind: sandboxv1.ErrorKind_ERROR_KIND_UNKNOWN, Provider: provider, Operation: operation, Message: err.Error()}, err)
}
