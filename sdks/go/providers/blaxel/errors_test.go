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

func TestMapProviderErrorSandboxAlreadyExists(t *testing.T) {
	parsed, err := url.Parse("https://example.invalid/sandboxes/demo")
	if err != nil {
		t.Fatal(err)
	}
	mapped := mapProviderError(&sdkblaxel.Error{
		StatusCode: http.StatusConflict,
		ErrorCode:  "SANDBOX_ALREADY_EXISTS",
		Request:    &http.Request{Method: http.MethodPost, URL: parsed},
		Response:   &http.Response{StatusCode: http.StatusConflict},
	}, "create")
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindAlreadyExists {
		t.Fatalf("error = %v detail = %#v", mapped, detail)
	}
}

func TestMapProviderErrorBareConflictIsAlreadyExists(t *testing.T) {
	parsed, err := url.Parse("https://example.invalid/sandboxes/demo")
	if err != nil {
		t.Fatal(err)
	}
	mapped := mapProviderError(&sdkblaxel.Error{
		StatusCode: http.StatusConflict,
		Request:    &http.Request{Method: http.MethodPost, URL: parsed},
		Response:   &http.Response{StatusCode: http.StatusConflict},
	}, "create")
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindAlreadyExists {
		t.Fatalf("error = %v detail = %#v", mapped, detail)
	}
}

func TestMapProviderErrorDeletionInProgressStaysConflict(t *testing.T) {
	parsed, err := url.Parse("https://example.invalid/sandboxes/demo")
	if err != nil {
		t.Fatal(err)
	}
	mapped := mapProviderError(&sdkblaxel.Error{
		StatusCode: http.StatusConflict,
		ErrorCode:  "SANDBOX_DELETION_IN_PROGRESS",
		Request:    &http.Request{Method: http.MethodDelete, URL: parsed},
		Response:   &http.Response{StatusCode: http.StatusConflict},
	}, "delete")
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindConflict {
		t.Fatalf("error = %v detail = %#v", mapped, detail)
	}
}

func TestMapProviderErrorGatewayTimeout(t *testing.T) {
	parsed, err := url.Parse("https://example.invalid/sandboxes/demo/process")
	if err != nil {
		t.Fatal(err)
	}
	mapped := mapProviderError(&sdkblaxel.Error{
		StatusCode: http.StatusGatewayTimeout,
		Request:    &http.Request{Method: http.MethodPost, URL: parsed},
		Response:   &http.Response{StatusCode: http.StatusGatewayTimeout},
	}, "run")
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindTimeout {
		t.Fatalf("error = %v detail = %#v", mapped, detail)
	}
}
