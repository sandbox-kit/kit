package daytona

import (
	"errors"
	"testing"

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
