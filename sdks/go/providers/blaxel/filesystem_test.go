package blaxel

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/sandbox"
)

func TestLocalSourceFileValidation(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.txt")
	if err := localSourceFile(missing, "copy_from_local"); err == nil {
		t.Fatal("expected missing file error")
	} else {
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindNotFound || detail.Field != sandbox.FileFieldLocalPath {
			t.Fatalf("error = %v detail = %#v", err, detail)
		}
	}
	if err := localSourceFile(dir, "copy_from_local"); err == nil {
		t.Fatal("expected directory error")
	} else {
		var detail *sandbox.Error
		if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindInvalidArgument || detail.Field != sandbox.FileFieldLocalPath {
			t.Fatalf("error = %v detail = %#v", err, detail)
		}
	}
	file := filepath.Join(dir, "note.txt")
	if err := os.WriteFile(file, []byte("hi"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := localSourceFile(file, "copy_from_local"); err != nil {
		t.Fatal(err)
	}
}

func TestRejectParentsFalse(t *testing.T) {
	err := rejectParents(&sandbox.MakeDirectoryOptions{Parents: sandbox.Value(false)})
	var detail *sandbox.Error
	if !errors.As(err, &detail) || detail.Kind != sandbox.ErrorKindUnsupported || detail.Field != sandbox.FileFieldParents {
		t.Fatalf("error = %v detail = %#v", err, detail)
	}
}
