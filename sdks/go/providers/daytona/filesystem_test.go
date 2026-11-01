package daytona

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	sdkerrors "github.com/daytona/clients/sdk-go/pkg/errors"
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
	err := localDestination(dir, "copy_to_local")
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldLocalPath || detail.Operation != "copy_to_local" {
		t.Fatalf("directory: %#v", err)
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
