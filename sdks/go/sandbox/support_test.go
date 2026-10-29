package sandbox

import (
	"errors"
	"testing"

	"github.com/sandbox-kit/kit/sdks/go/internal/sdk"
)

func TestRejectUnmappedUsesProtobufFieldName(t *testing.T) {
	t.Parallel()

	err := sdk.RejectUnmapped("test", "", &Resources{CpuCores: Value(2.0)})
	if err == nil {
		t.Fatal("expected unsupported error")
	}
	var kitErr *Error
	if !errors.As(err, &kitErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if kitErr.Field != "cpu_cores" {
		t.Fatalf("Field = %q, want cpu_cores", kitErr.Field)
	}
	nested := sdk.RejectUnmapped("test", "resources", &Resources{CpuCores: Value(2.0)})
	if !errors.As(nested, &kitErr) || kitErr.Field != "resources.cpu_cores" {
		t.Fatalf("nested Field = %#v", kitErr)
	}
	if kitErr.Kind != ErrorKindUnsupported {
		t.Fatalf("Kind = %v, want Unsupported", kitErr.Kind)
	}
}

func TestRejectUnmappedUsesProtobufNameOverJSONAlias(t *testing.T) {
	t.Parallel()

	// provisioning serializes as json=creation but proto name stays provisioning.
	err := sdk.RejectUnmapped("test", "", &CreateOptions{
		Provisioning: &ProvisioningOptions{},
	})
	if err == nil {
		t.Fatal("expected unsupported error")
	}
	var kitErr *Error
	if !errors.As(err, &kitErr) {
		t.Fatalf("expected *Error, got %T", err)
	}
	if kitErr.Field != "provisioning" {
		t.Fatalf("Field = %q, want provisioning", kitErr.Field)
	}
}

func TestRejectUnmappedNilRemainder(t *testing.T) {
	t.Parallel()

	if err := sdk.RejectUnmapped("test", "", (*Resources)(nil)); err != nil {
		t.Fatalf("nil remainder: %v", err)
	}
	if err := sdk.RejectUnmapped("test", "resources", &Resources{}); err != nil {
		t.Fatalf("empty remainder: %v", err)
	}
}
