package sdk

import (
	"context"
	"errors"

	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/internal/contracts/kit/sandbox/v1"
)

// Error is the failure type returned by the sandbox client.
// Read it with errors.As. Kind is the category. Field names the setting to change.
type Error struct {
	// Kind is the category, such as ErrorKindUnsupported or ErrorKindAuthentication.
	Kind sandboxv1.ErrorKind
	// Provider is the provider name, such as "daytona" or "modal". It is empty when the failure happened before a provider was selected.
	Provider string
	// Operation is the call that failed, such as "create", "get", or "initialize".
	Operation string
	// Field is the setting to change, such as "resources.cpu_cores" or "endpoint".
	// It is empty when no single setting applies. Compare it with CreateFieldResourcesCpuCores and the other field constants.
	Field string
	// Message explains the failure.
	Message string
	// StatusCode is the HTTP or transport status reported by the provider.
	// It is nil when the provider did not report one. GetStatusCode returns 0 in that case.
	StatusCode *uint32
	// ProviderCode is the provider's own error code.
	// It is nil when the provider did not report one.
	ProviderCode *string
	// ProviderSource is the provider component that reported the error.
	// It is nil when the provider did not name one.
	ProviderSource *string
	// Rule is the validation or policy identifier for this failure.
	// It is nil when the failure has no rule id.
	Rule *string
	// cause is the provider or validator error. Unwrap returns it.
	cause error
}

// NewError builds an Error from info and keeps cause for Unwrap.
// An unrecognized kind becomes ErrorKindUnknown.
// A cause of context.Canceled sets ErrorKindCanceled. A cause of context.DeadlineExceeded sets ErrorKindTimeout.
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

// Error returns Message, with Field appended in parentheses when Field is set.
// A nil error returns an empty string.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Field != "" {
		return e.Message + " (" + e.Field + ")"
	}
	return e.Message
}

// Unwrap returns the provider or validator error so errors.Is and errors.As can match it.
// A nil error returns nil.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// GetStatusCode returns the provider HTTP or transport status, or 0 when StatusCode is nil.
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
