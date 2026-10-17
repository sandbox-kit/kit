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
	v10 "github.com/go-playground/validator/v10"
	sandboxv1 "github.com/sandbox-kit/kit/sdks/go/contracts/kit/sandbox/v1"
	"google.golang.org/protobuf/proto"
	strings "strings"
)

// ErrorInfo and ErrorKind are the generated portable protobuf contracts.
// Error below only adds native cause chaining.
type ErrorInfo = sandboxv1.ErrorInfo
type ErrorKind = sandboxv1.ErrorKind

const (
	ErrorKindUnknown = sandboxv1.ErrorKind_ERROR_KIND_UNKNOWN
	ErrorKindInvalidArgument = sandboxv1.ErrorKind_ERROR_KIND_INVALID_ARGUMENT
	ErrorKindUnsupported = sandboxv1.ErrorKind_ERROR_KIND_UNSUPPORTED
	ErrorKindAuthentication = sandboxv1.ErrorKind_ERROR_KIND_AUTHENTICATION
	ErrorKindPermissionDenied = sandboxv1.ErrorKind_ERROR_KIND_PERMISSION_DENIED
	ErrorKindNotFound = sandboxv1.ErrorKind_ERROR_KIND_NOT_FOUND
	ErrorKindAlreadyExists = sandboxv1.ErrorKind_ERROR_KIND_ALREADY_EXISTS
	ErrorKindConflict = sandboxv1.ErrorKind_ERROR_KIND_CONFLICT
	ErrorKindRateLimited = sandboxv1.ErrorKind_ERROR_KIND_RATE_LIMITED
	ErrorKindResourceExhausted = sandboxv1.ErrorKind_ERROR_KIND_RESOURCE_EXHAUSTED
	ErrorKindTimeout = sandboxv1.ErrorKind_ERROR_KIND_TIMEOUT
	ErrorKindCanceled = sandboxv1.ErrorKind_ERROR_KIND_CANCELED
	ErrorKindUnavailable = sandboxv1.ErrorKind_ERROR_KIND_UNAVAILABLE
	ErrorKindInvalidResponse = sandboxv1.ErrorKind_ERROR_KIND_INVALID_RESPONSE
	ErrorKindInternal = sandboxv1.ErrorKind_ERROR_KIND_INTERNAL
	ErrorKindFailedPrecondition = sandboxv1.ErrorKind_ERROR_KIND_FAILED_PRECONDITION
)

// Error exposes portable details and retains the native cause for errors.Is
// and errors.As.
//
// Example:
//
//	var sandboxErr *Error
//	if errors.As(err, &sandboxErr) {
//	    fmt.Println(sandboxErr.Kind, sandboxErr.Field)
//	}
type Error struct {
	// ErrorInfo holds Kind, Provider, Operation, Field, and Message. Switch on
	// Kind.
	ErrorInfo
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
//	err := NewError(ErrorInfo{
//	    Kind:    ErrorKindUnsupported,
//	    Field:   CreateFieldResourcesDiskMiB,
//	    Message: "disk overrides are unsupported",
//	}, cause)
func NewError(info ErrorInfo, cause error) *Error {
	owned := proto.Clone(&info).(*sandboxv1.ErrorInfo)
	if owned.Kind.Descriptor().Values().ByNumber(owned.Kind.Number()) == nil {
		owned.Kind = ErrorKindUnknown
	}
	if errors.Is(cause, context.Canceled) {
		owned.Kind = ErrorKindCanceled
	} else if errors.Is(cause, context.DeadlineExceeded) {
		owned.Kind = ErrorKindTimeout
	}
	return &Error{ErrorInfo: *owned, cause: cause}
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
		info := proto.Clone(&existing.ErrorInfo).(*sandboxv1.ErrorInfo)
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

// validationError wraps a validator failure as ErrorKindInvalidArgument and
// records the first invalid field path.
func validationError(err error) error {
	field := ""
	rule := ""
	var fields v10.ValidationErrors
	if errors.As(err, &fields) && len(fields) > 0 {
		field = validationPath(fields[0].StructNamespace())
		rule = fields[0].Tag()
		switch rule {
		case "gt":
			rule = "exclusive_minimum"
		case "gte":
			rule = "minimum"
		case "lte":
			rule = "maximum"
		case "min":
			rule = "min_items"
		}
	}
	return NewError(ErrorInfo{Kind: ErrorKindInvalidArgument, Field: field, Rule: Value(rule), Message: "sandbox-kit: invalid configuration"}, err)
}

// validationPath converts a validator namespace such as Config.Auth.APIKey
// into the canonical path auth.api_key.
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
