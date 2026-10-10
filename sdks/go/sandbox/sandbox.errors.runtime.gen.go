// Code generated from contracts.yaml and errors.yaml by the Go emitter. DO NOT EDIT.
package sandbox

import (
	context "context"
	errors "errors"
	v10 "github.com/go-playground/validator/v10"
	strings "strings"
)

// Error exposes portable details and retains the original error for errors.Is/As.
// Error exposes portable details and retains the native cause.
type Error struct {
	ErrorInfo
	cause error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Field != "" {
		return e.Message + " (" + e.Field + ")"
	}
	return e.Message
}

// Unwrap preserves native SDK and context matching through errors.Is/As.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

// NewError owns a copy of portable details and retains the original cause.
// NewError copies portable details and retains the original cause.
func NewError(info ErrorInfo, cause error) *Error {
	owned, _ := info.Clone()
	if !owned.Kind.Valid() {
		owned.Kind = ErrorKindUnknown
	}
	if errors.Is(cause, context.Canceled) {
		owned.Kind = ErrorKindCanceled
	} else if errors.Is(cause, context.DeadlineExceeded) {
		owned.Kind = ErrorKindTimeout
	}
	return &Error{ErrorInfo: *owned, cause: cause}
}

// WithErrorContext adds provider/operation details without mutating an error.
func WithErrorContext(err error, provider string, operation string) error {
	if err == nil {
		return nil
	}
	var existing *Error
	if errors.As(err, &existing) && existing != nil {
		if (provider == "" || existing.Provider == provider) && (operation == "" || existing.Operation == operation) {
			return err
		}
		info, _ := existing.ErrorInfo.Clone()
		if provider != "" {
			info.Provider = provider
		}
		if operation != "" {
			info.Operation = operation
		}
		return &Error{ErrorInfo: *info, cause: err}
	}
	return NewError(ErrorInfo{Kind: ErrorKindUnknown, Provider: provider, Operation: operation, Message: err.Error()}, err)
}

func validationError(err error) error {
	field := ""
	var fields v10.ValidationErrors
	if errors.As(err, &fields) && len(fields) > 0 {
		field = validationPath(fields[0].StructNamespace())
	}
	return NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Field: field, Message: "sandbox-kit: invalid configuration"}, err)
}

func validationPath(namespace string) string {
	parts := strings.Split(namespace, ".")
	if len(parts) > 1 {
		parts = parts[1:]
	}
	for i, part := range parts {
		base, suffix := part, ""
		if index := strings.IndexByte(part, '['); index >= 0 {
			base, suffix = part[:index], part[index:]
		}
		if native, ok := validationFieldNames[base]; ok {
			parts[i] = native + suffix
		}
		if native, ok := configValidationFieldNames[base]; ok {
			parts[i] = native + suffix
		}
	}
	return strings.Join(parts, ".")
}
