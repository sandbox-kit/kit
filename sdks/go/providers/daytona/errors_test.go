package daytona

import (
	"context"
	"errors"
	"fmt"
	sdkerrors "github.com/daytona/clients/sdk-go/pkg/errors"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
	"testing"
)

func TestGeneratedNativeErrorMappingsPreserveCause(t *testing.T) {
	for _, tc := range []struct {
		status int
		kind   sandbox.ErrorKind
	}{{400, sandbox.ErrorKindInvalidArgument}, {401, sandbox.ErrorKindAuthentication}, {403, sandbox.ErrorKindPermissionDenied}, {404, sandbox.ErrorKindNotFound}, {429, sandbox.ErrorKindRateLimited}, {503, sandbox.ErrorKindUnavailable}, {504, sandbox.ErrorKindTimeout}, {418, sandbox.ErrorKindUnknown}} {
		native := &sdkerrors.DaytonaError{Message: "native failure", StatusCode: tc.status, Code: "provider-code", Source: "DAYTONA_API"}
		err := mapProviderError(fmt.Errorf("wrapped: %w", native), "create")
		var detail *sandbox.Error
		var recovered *sdkerrors.DaytonaError
		if !errors.As(err, &detail) || detail.Kind != tc.kind || detail.Provider != "daytona" || detail.Operation != "create" || detail.GetStatusCode() != uint32(tc.status) || detail.GetProviderCode() != "provider-code" || !errors.As(err, &recovered) || recovered != native || !errors.Is(err, native) {
			t.Fatalf("incorrect mapping for %d: %v", tc.status, err)
		}
		if tc.status == 404 && !errors.Is(err, sdkerrors.ErrNotFound) {
			t.Fatal("native sentinel matching lost")
		}
	}
	var detail *sandbox.Error
	err := mapProviderError(context.DeadlineExceeded, "create")
	if !errors.Is(err, context.DeadlineExceeded) || !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindTimeout {
		t.Fatal("context deadline cause lost")
	}
}
