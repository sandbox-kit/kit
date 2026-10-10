package modal

import (
	"errors"
	sdk "github.com/modal-labs/modal-client/go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestGeneratedTypedAndTransportErrors(t *testing.T) {
	for _, tc := range []struct {
		native error
		kind   sandbox.ErrorKind
	}{
		{sdk.NotFoundError{Exception: "missing"}, sandbox.ErrorKindNotFound},
		{&sdk.NotFoundError{Exception: "missing"}, sandbox.ErrorKindNotFound},
		{sdk.InvalidError{Exception: "invalid"}, sandbox.ErrorKindInvalidArgument},
		{sdk.ResourceExhaustedError{Exception: "quota"}, sandbox.ErrorKindResourceExhausted},
		{status.Error(codes.Unauthenticated, "credentials"), sandbox.ErrorKindAuthentication},
		{status.Error(codes.PermissionDenied, "permissions"), sandbox.ErrorKindPermissionDenied},
		{errors.New("unclassified"), sandbox.ErrorKindUnknown},
	} {
		err := mapProviderError(tc.native, "create")
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != tc.kind || detail.Provider != "modal" || !errors.Is(err, tc.native) {
			t.Fatalf("incorrect native classification: %v", err)
		}
	}
}
