package main

import (
	"errors"
	"fmt"
	"io"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

// handleError uses the same common kinds for either provider.
func handleError(w io.Writer, err error) {
	if err == nil {
		return
	}
	var sandboxError *sandbox.Error
	if !errors.As(err, &sandboxError) {
		fmt.Fprintln(w, err)
		return
	}
	switch sandboxError.Kind {
	case sandbox.ErrorKindInvalidArgument, sandbox.ErrorKindUnsupported:
		fmt.Fprintf(w, "Check %s: %s\n", sandboxError.Field, sandboxError.Message)
	case sandbox.ErrorKindAuthentication:
		fmt.Fprintln(w, "Check your provider credentials.")
	case sandbox.ErrorKindPermissionDenied:
		fmt.Fprintln(w, "Your account does not have permission for this operation.")
	case sandbox.ErrorKindNotFound:
		fmt.Fprintln(w, "The requested provider resource was not found.")
	case sandbox.ErrorKindRateLimited, sandbox.ErrorKindResourceExhausted:
		fmt.Fprintln(w, "The provider's rate or resource limit was reached.")
	case sandbox.ErrorKindCanceled:
		fmt.Fprintln(w, "The operation was canceled.")
	case sandbox.ErrorKindTimeout:
		fmt.Fprintln(w, "The operation timed out.")
	default:
		fmt.Fprintln(w, err)
	}
}
