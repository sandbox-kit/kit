package modal

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestMakeDirectoryParams(t *testing.T) {
	for _, options := range []*sandbox.MakeDirectoryOptions{
		nil,
		{},
		{Parents: sandbox.Value(true)},
	} {
		params := makeDirectoryParams(options)
		if params == nil || params.CreateParents == nil || !*params.CreateParents {
			t.Fatalf("parents true: %#v from %#v", params, options)
		}
	}
	params := makeDirectoryParams(&sandbox.MakeDirectoryOptions{Parents: sandbox.Value(false)})
	if params == nil || params.CreateParents == nil || *params.CreateParents {
		t.Fatalf("parents false: %#v", params)
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

func TestMapFileErrorLocalPathOnly(t *testing.T) {
	cause := &os.PathError{Op: "open", Path: "/tmp/note.txt", Err: os.ErrNotExist}
	err := mapFileError(cause, "copy_from_local", sandbox.FileFieldLocalPath)
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindNotFound || detail.Field != sandbox.FileFieldLocalPath || !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("copy path error: %#v", err)
	}
	err = mapFileError(cause, "read_file", "")
	if !errors.As(err, &detail) || detail.Field != "" || detail.Operation != "read_file" {
		t.Fatalf("remote path error: %#v", err)
	}
}
