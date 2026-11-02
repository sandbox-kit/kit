package daytona

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	sdkerrors "github.com/daytona/clients/sdk-go/pkg/errors"
	"github.com/daytona/clients/toolbox-api-client-go"
	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestRejectParents(t *testing.T) {
	if err := rejectParents(nil); err != nil {
		t.Fatal(err)
	}
	if err := rejectParents(&sandbox.MakeDirectoryOptions{}); err != nil {
		t.Fatal(err)
	}
	if err := rejectParents(&sandbox.MakeDirectoryOptions{Parents: sandbox.Value(true)}); err != nil {
		t.Fatal(err)
	}
	err := rejectParents(&sandbox.MakeDirectoryOptions{Parents: sandbox.Value(false)})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != sandbox.FileFieldParents || detail.Operation != "make_directory" {
		t.Fatalf("parents false: %#v", err)
	}
}

func TestLocalSourceFile(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.txt")
	err := localSourceFile(missing, "copy_from_local")
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindNotFound || detail.Field != sandbox.FileFieldLocalPath || detail.Operation != "copy_from_local" || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("missing file: %#v", err)
	}
	err = localSourceFile(dir, "copy_from_local")
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldLocalPath {
		t.Fatalf("directory: %#v", err)
	}
	path := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(path, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := localSourceFile(path, "copy_from_local"); err != nil {
		t.Fatal(err)
	}
}

func TestLocalDestination(t *testing.T) {
	dir := t.TempDir()
	if err := localDestination(filepath.Join(dir, "out.txt"), "copy_to_local"); err != nil {
		t.Fatal(err)
	}
	if err := localDestination(filepath.Join(dir, "missing", "out.txt"), "copy_to_local"); err != nil {
		t.Fatal(err)
	}
	err := localDestination(dir, "copy_to_local")
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldLocalPath || detail.Operation != "copy_to_local" {
		t.Fatalf("directory: %#v", err)
	}
	path := filepath.Join(dir, "out.txt")
	if err := os.WriteFile(path, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := localDestination(path, "copy_to_local"); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureLocalParent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "out.txt")
	if err := ensureLocalParent(path, "copy_to_local"); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Dir(path))
	if err != nil || !info.IsDir() {
		t.Fatalf("parent: %#v %#v", info, err)
	}
	if err := ensureLocalParent(filepath.Join(dir, "out.txt"), "copy_to_local"); err != nil {
		t.Fatal(err)
	}
}

func TestMapCallerPathError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.txt")
	cause := &os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}
	err := mapCallerPathError(sdkerrors.WrapDaytonaError(cause, "Failed to write file", 0, nil), "copy_to_local", path)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindNotFound || detail.Field != sandbox.FileFieldLocalPath || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("local write: %#v", err)
	}
	temp := &os.PathError{Op: "open", Path: filepath.Join(t.TempDir(), "daytona-upload"), Err: os.ErrNotExist}
	err = mapCallerPathError(sdkerrors.WrapDaytonaError(temp, "Failed to read file", 0, nil), "copy_from_local", path)
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnavailable || detail.Field != "" {
		t.Fatalf("temp file: %#v", err)
	}
	denied := &os.PathError{Op: "open", Path: path, Err: os.ErrPermission}
	err = mapCallerPathError(sdkerrors.WrapDaytonaError(denied, "Failed to write file", 0, nil), "copy_to_local", path)
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldLocalPath {
		t.Fatalf("permission: %#v", err)
	}
}

func TestDaytonaFileInfoKeepsToolboxFields(t *testing.T) {
	modified := "2020-01-02T15:04:05.123456789Z"
	info, err := daytonaFileInfo(&toolbox.FileInfo{
		Name:        "note.txt",
		Path:        sandbox.Value("/tmp/note.txt"),
		Size:        12,
		Mode:        "0644",
		Permissions: "rw-r--r--",
		Owner:       "daytona",
		Group:       "staff",
		ModifiedAt:  modified,
		ModTime:     "not-used",
		IsDir:       false,
	}, "stat")
	if err != nil {
		t.Fatal(err)
	}
	if info.GetPath() != "/tmp/note.txt" || info.GetName() != "note.txt" || info.GetSize() != 12 || info.GetType() != sandbox.FileTypeFile {
		t.Fatalf("identity: %#v", info)
	}
	if info.GetMode() != "0644" || info.GetPermissions() != "rw-r--r--" || info.GetOwner() != "daytona" || info.GetGroup() != "staff" {
		t.Fatalf("metadata: %#v", info)
	}
	want := time.Date(2020, 1, 2, 15, 4, 5, 123456789, time.UTC)
	if info.ModifiedAt == nil || !info.ModifiedAt.AsTime().Equal(want) {
		t.Fatalf("modified: %#v", info.ModifiedAt)
	}
}

func TestDaytonaFileInfoOmitsEmptyAndUnparsedTimes(t *testing.T) {
	info, err := daytonaFileInfo(&toolbox.FileInfo{Name: "nested", IsDir: true, ModifiedAt: "not-a-timestamp", ModTime: "2020-01-02T15:04:05Z"}, "list_directory")
	if err != nil {
		t.Fatal(err)
	}
	if info.GetType() != sandbox.FileTypeDirectory || info.Path != "" || info.Mode != nil || info.Permissions != nil || info.Owner != nil || info.Group != nil || info.ModifiedAt != nil {
		t.Fatalf("directory: %#v", info)
	}
	info, err = daytonaFileInfo(&toolbox.FileInfo{Name: "note.txt", ModTime: "2020-01-02T15:04:05Z"}, "stat")
	if err != nil {
		t.Fatal(err)
	}
	want := time.Date(2020, 1, 2, 15, 4, 5, 0, time.UTC)
	if info.ModifiedAt == nil || !info.ModifiedAt.AsTime().Equal(want) {
		t.Fatalf("modTime fallback: %#v", info.ModifiedAt)
	}
	if _, err := daytonaFileInfo(nil, "stat"); err == nil {
		t.Fatal("nil metadata accepted")
	}
}

func TestToolboxStatAndListKeepReportedFields(t *testing.T) {
	const body = `{"group":"staff","isDir":false,"modTime":"2020-01-02T15:04:05Z","mode":"0644","modifiedAt":"2020-01-02T15:04:05.123456789Z","name":"note.txt","owner":"daytona","path":"/tmp/note.txt","permissions":"rw-r--r--","size":12}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("path") != "/tmp" && r.URL.Query().Get("path") != "/tmp/note.txt" {
			t.Errorf("path query: %s %s", r.URL.Path, r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/files/info":
			_, _ = w.Write([]byte(body))
		case "/files":
			_, _ = w.Write([]byte("[" + body + "]"))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	cfg := toolbox.NewConfiguration()
	cfg.Servers = toolbox.ServerConfigurations{{URL: server.URL}}
	cfg.HTTPClient = server.Client()
	client := toolbox.NewAPIClient(cfg)

	info, httpResp, err := client.FileSystemAPI.GetFileInfo(context.Background()).Path("/tmp/note.txt").Execute()
	if err != nil {
		t.Fatal(err)
	}
	if httpResp != nil && httpResp.Body != nil {
		httpResp.Body.Close()
	}
	mapped, err := daytonaFileInfo(info, "stat")
	if err != nil {
		t.Fatal(err)
	}
	if mapped.GetPath() != "/tmp/note.txt" || mapped.GetOwner() != "daytona" || mapped.GetPermissions() != "rw-r--r--" || mapped.GetGroup() != "staff" || mapped.ModifiedAt == nil {
		t.Fatalf("stat: %#v", mapped)
	}

	entries, httpResp, err := client.FileSystemAPI.ListFiles(context.Background()).Path("/tmp").Execute()
	if err != nil {
		t.Fatal(err)
	}
	if httpResp != nil && httpResp.Body != nil {
		httpResp.Body.Close()
	}
	if len(entries) != 1 {
		t.Fatalf("entries: %#v", entries)
	}
	mapped, err = daytonaFileInfo(&entries[0], "list_directory")
	if err != nil {
		t.Fatal(err)
	}
	if mapped.GetPath() != "/tmp/note.txt" || mapped.GetOwner() != "daytona" {
		t.Fatalf("list: %#v", mapped)
	}
}
