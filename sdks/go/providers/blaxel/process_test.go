package blaxel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	blaxelsdk "github.com/blaxel-ai/sdk-go"
	"github.com/blaxel-ai/sdk-go/option"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestBlaxelCommandArgvAndShell(t *testing.T) {
	command, err := blaxelCommand(&sandbox.ExecOptions{Argv: []string{"echo", "hello world"}})
	if err != nil || command != "'echo' 'hello world'" {
		t.Fatalf("argv command = %q err=%v", command, err)
	}
	command, err = blaxelCommand(&sandbox.ExecOptions{Shell: sandbox.Value("echo hi")})
	if err != nil || command != "echo hi" {
		t.Fatalf("shell command = %q err=%v", command, err)
	}
}

func TestCommandTimeoutWholeSeconds(t *testing.T) {
	timeout, err := commandTimeout("run", sandbox.Duration(2*time.Second))
	if err != nil || timeout != 2*time.Second {
		t.Fatalf("timeout = %v err=%v", timeout, err)
	}
	_, err = commandTimeout("run", sandbox.Duration(1500*time.Millisecond))
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Field != sandbox.ExecFieldTimeout {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}

func TestRunCommandStartsWithoutCompletionHold(t *testing.T) {
	var posts []map[string]any
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/sbx":
			writeJSON(t, w, map[string]any{"metadata": map[string]string{"name": "sbx"}, "spec": map[string]any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/process":
			body := readJSON(t, r)
			posts = append(posts, body)
			writeProcess(t, w, "completed", 0, "ok", "")
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	for _, timeout := range []*sandbox.ExecOptions{
		{Argv: []string{"echo", "hi"}, Timeout: sandbox.Duration(5 * time.Second), WorkingDirectory: sandbox.Value("/workspace"), Environment: map[string]string{"MODE": "test"}},
		{Argv: []string{"echo", "hi"}, Timeout: sandbox.Duration(60 * time.Second)},
	} {
		result, err := backend.RunCommand(context.Background(), "sbx", timeout)
		if err != nil {
			t.Fatal(err)
		}
		if result.GetExitCode() != 0 || string(result.GetStdout()) != "ok" || string(result.GetStderr()) != "" {
			t.Fatalf("result = %#v", result)
		}
	}
	if len(posts) != 2 {
		t.Fatalf("posts = %#v", posts)
	}
	if posts[0]["waitForCompletion"] != false || posts[0]["timeout"] != float64(5) || posts[0]["command"] != "'echo' 'hi'" || posts[0]["workingDir"] != "/workspace" {
		t.Fatalf("short post = %#v", posts[0])
	}
	env, _ := posts[0]["env"].(map[string]any)
	if env["MODE"] != "test" {
		t.Fatalf("env = %#v", posts[0]["env"])
	}
	if posts[1]["waitForCompletion"] != false || posts[1]["timeout"] != float64(60) {
		t.Fatalf("60s post = %#v", posts[1])
	}
}

func TestRunCommandWaitsWhenStartReturnsRunning(t *testing.T) {
	var gets atomic.Int32
	var posts []map[string]any
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/sbx":
			writeJSON(t, w, map[string]any{"metadata": map[string]string{"name": "sbx"}, "spec": map[string]any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/process":
			posts = append(posts, readJSON(t, r))
			writeProcess(t, w, "running", 0, "", "")
		case r.Method == http.MethodGet && r.URL.Path == "/process/123":
			if gets.Add(1) == 1 {
				http.Error(w, "gateway timeout", http.StatusGatewayTimeout)
				return
			}
			writeProcess(t, w, "failed", 7, "out", "err")
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	cases := []struct {
		name    string
		timeout *sandbox.ExecOptions
		seconds float64
	}{
		{"omit", &sandbox.ExecOptions{Shell: sandbox.Value("echo hi")}, 0},
		{"zero", &sandbox.ExecOptions{Argv: []string{"echo", "hi"}, Timeout: sandbox.Duration(0)}, 0},
		{"long", &sandbox.ExecOptions{Argv: []string{"echo", "hi"}, Timeout: sandbox.Duration(61 * time.Second)}, 61},
		{"short", &sandbox.ExecOptions{Argv: []string{"echo", "hi"}, Timeout: sandbox.Duration(5 * time.Second)}, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gets.Store(0)
			result, err := backend.RunCommand(context.Background(), "sbx", tc.timeout)
			if err != nil {
				t.Fatal(err)
			}
			if result.GetExitCode() != 7 || string(result.GetStdout()) != "out" || string(result.GetStderr()) != "err" {
				t.Fatalf("result = %#v", result)
			}
			if gets.Load() != 2 {
				t.Fatalf("gets = %d", gets.Load())
			}
		})
	}
	if len(posts) != len(cases) {
		t.Fatalf("posts = %d", len(posts))
	}
	for i, tc := range cases {
		if posts[i]["waitForCompletion"] != false || posts[i]["timeout"] != tc.seconds {
			t.Fatalf("%s post = %#v", tc.name, posts[i])
		}
	}
}

func TestRunCommandReturnsAFinishedStartWithoutPolling(t *testing.T) {
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/sbx":
			writeJSON(t, w, map[string]any{"metadata": map[string]string{"name": "sbx"}, "spec": map[string]any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/process":
			writeProcess(t, w, "killed", 9, "partial", "stopped")
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	result, err := backend.RunCommand(context.Background(), "sbx", &sandbox.ExecOptions{Argv: []string{"sleep", "1"}})
	if err != nil {
		t.Fatal(err)
	}
	if result.GetExitCode() != 9 || string(result.GetStdout()) != "partial" || string(result.GetStderr()) != "stopped" {
		t.Fatalf("result = %#v", result)
	}
}

func TestRunCommandRequiresAProcessIdentifier(t *testing.T) {
	backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/sbx":
			writeJSON(t, w, map[string]any{"metadata": map[string]string{"name": "sbx"}, "spec": map[string]any{}})
		case r.Method == http.MethodPost && r.URL.Path == "/process":
			writeJSON(t, w, map[string]any{"status": "running", "exitCode": 0, "stdout": "", "stderr": ""})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	_, err := backend.RunCommand(context.Background(), "sbx", &sandbox.ExecOptions{Argv: []string{"echo", "hi"}})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidResponse || detail.Operation != "run" {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}

func TestRunCommandCancelAndGatewayTimeout(t *testing.T) {
	t.Run("cancel", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/sbx":
				writeJSON(t, w, map[string]any{"metadata": map[string]string{"name": "sbx"}, "spec": map[string]any{}})
			case r.Method == http.MethodPost && r.URL.Path == "/process":
				writeProcess(t, w, "running", 0, "", "")
			case r.Method == http.MethodGet && r.URL.Path == "/process/123":
				cancel()
				writeProcess(t, w, "running", 0, "", "")
			default:
				http.NotFound(w, r)
			}
		})
		_, err := backend.RunCommand(ctx, "sbx", &sandbox.ExecOptions{Argv: []string{"sleep", "30"}, Timeout: sandbox.Duration(0)})
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindCanceled {
			t.Fatalf("error = %v detail = %#v", err, detail)
		}
	})
	t.Run("gateway timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()
		backend := processBackend(t, func(w http.ResponseWriter, r *http.Request) {
			switch {
			case r.Method == http.MethodGet && r.URL.Path == "/sandboxes/sbx":
				writeJSON(t, w, map[string]any{"metadata": map[string]string{"name": "sbx"}, "spec": map[string]any{}})
			case r.Method == http.MethodPost && r.URL.Path == "/process":
				writeProcess(t, w, "running", 0, "", "")
			case r.Method == http.MethodGet && r.URL.Path == "/process/123":
				http.Error(w, "gateway timeout", http.StatusGatewayTimeout)
			default:
				http.NotFound(w, r)
			}
		})
		_, err := backend.RunCommand(ctx, "sbx", &sandbox.ExecOptions{Argv: []string{"sleep", "90"}, Timeout: sandbox.Duration(90 * time.Second)})
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindTimeout {
			t.Fatalf("error = %v detail = %#v", err, detail)
		}
	})
}

func processBackend(t *testing.T, handler http.HandlerFunc) *backend {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := blaxelsdk.NewClient(
		option.WithBaseURL(server.URL),
		option.WithAPIKey("test"),
		option.WithMaxRetries(0),
	)
	return &backend{client: &client}
}

func writeProcess(t *testing.T, w http.ResponseWriter, status string, exitCode int, stdout, stderr string) {
	t.Helper()
	writeJSON(t, w, map[string]any{
		"command":     "echo",
		"completedAt": "",
		"exitCode":    exitCode,
		"logs":        stdout + stderr,
		"name":        "proc",
		"pid":         "123",
		"startedAt":   "",
		"status":      status,
		"stderr":      stderr,
		"stdout":      stdout,
		"workingDir":  "/",
	})
}

func writeJSON(t *testing.T, w http.ResponseWriter, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("encode: %v", err)
	}
}

func readJSON(t *testing.T, r *http.Request) map[string]any {
	t.Helper()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("decode %s: %v", data, err)
	}
	return body
}

func TestMapProviderErrorWaitDeadlineKeepsGatewayTimeout(t *testing.T) {
	parsed, err := http.NewRequest(http.MethodGet, "https://example.invalid/process/123", nil)
	if err != nil {
		t.Fatal(err)
	}
	apiErr := &blaxelsdk.Error{
		StatusCode: http.StatusGatewayTimeout,
		Request:    parsed,
		Response:   &http.Response{StatusCode: http.StatusGatewayTimeout},
	}
	mapped := mapProviderError(fmt.Errorf(
		"process 123 did not finish in time; it may still be running: %w",
		errors.Join(context.DeadlineExceeded, apiErr),
	), "run")
	var detail *sandbox.Error
	if !errors.As(mapped, &detail) || detail.Kind != sandbox.ErrorKindTimeout || detail.GetStatusCode() != http.StatusGatewayTimeout {
		t.Fatalf("error = %v detail = %#v", mapped, detail)
	}
}
