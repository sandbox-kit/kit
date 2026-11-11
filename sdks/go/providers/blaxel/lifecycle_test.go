package blaxel

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestDeleteSandboxReportsDeletedWithoutWait(t *testing.T) {
	var gets atomic.Int32
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/sandboxes/doomed":
			writeJSON(t, w, map[string]any{
				"metadata": map[string]string{"name": "doomed"},
				"status":   "DELETING",
				"spec":     map[string]any{},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/doomed":
			gets.Add(1)
			http.NotFound(w, r)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	info, err := backend.DeleteSandbox(context.Background(), "doomed", &sandbox.DeleteOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if info.State != sandbox.SandboxStateDeleted || info.GetTerminal() != true || info.GetProviderState() != "TERMINATED" {
		t.Fatalf("info = %#v", info)
	}
	if gets.Load() != 0 {
		t.Fatalf("gets = %d", gets.Load())
	}
}

func TestDeleteSandboxWaitsUntilNotFound(t *testing.T) {
	var gets atomic.Int32
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/sandboxes/doomed":
			writeJSON(t, w, map[string]any{
				"metadata": map[string]string{"name": "doomed"},
				"status":   "DELETING",
				"spec":     map[string]any{},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/doomed":
			n := gets.Add(1)
			if n < 2 {
				writeJSON(t, w, map[string]any{
					"metadata": map[string]string{"name": "doomed"},
					"status":   "DELETING",
					"spec":     map[string]any{},
				})
				return
			}
			http.NotFound(w, r)
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	info, err := backend.DeleteSandbox(ctx, "doomed", &sandbox.DeleteOptions{Wait: sandbox.Value(true)})
	if err != nil {
		t.Fatal(err)
	}
	if info.State != sandbox.SandboxStateDeleted || info.GetTerminal() != true {
		t.Fatalf("info = %#v", info)
	}
	if gets.Load() < 2 {
		t.Fatalf("gets = %d", gets.Load())
	}
}

func TestDeleteSandboxWaitHonorsCancel(t *testing.T) {
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodDelete && r.URL.Path == "/sandboxes/doomed":
			writeJSON(t, w, map[string]any{
				"metadata": map[string]string{"name": "doomed"},
				"status":   "DELETING",
				"spec":     map[string]any{},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/doomed":
			writeJSON(t, w, map[string]any{
				"metadata": map[string]string{"name": "doomed"},
				"status":   "DELETING",
				"spec":     map[string]any{},
			})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel()
	}()
	_, err := backend.DeleteSandbox(ctx, "doomed", &sandbox.DeleteOptions{Wait: sandbox.Value(true)})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindCanceled {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}
