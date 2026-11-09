package blaxel

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	sdkblaxel "github.com/blaxel-ai/sdk-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestMapProviderErrorClassifiesAPIErrors(t *testing.T) {
	parsed, err := url.Parse("https://example.invalid/sandboxes/demo")
	if err != nil {
		t.Fatal(err)
	}
	mapped := mapProviderError(&sdkblaxel.Error{
		StatusCode: http.StatusNotFound,
		ErrorCode:  sdkblaxel.ErrWorkloadNotFound,
		Request:    &http.Request{Method: http.MethodGet, URL: parsed},
		Response:   &http.Response{StatusCode: http.StatusNotFound},
	}, "refresh")
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindNotFound || detail.Provider != "blaxel" || detail.Operation != "refresh" {
		t.Fatalf("error = %v detail = %#v", mapped, detail)
	}
}

func TestMapProviderErrorPreservesContext(t *testing.T) {
	err := mapProviderError(context.Canceled, "run")
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindCanceled {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}
