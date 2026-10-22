package sandbox

import (
	"errors"
	"testing"
)

func TestRejectUnmappedUsesProtobufFieldName(t *testing.T) {
	t.Parallel()

	err := RejectUnmapped("test", &Resources{CpuCores: Value(2.0)})
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
	if kitErr.Kind != ErrorKindUnsupported {
		t.Fatalf("Kind = %v, want Unsupported", kitErr.Kind)
	}
}

func TestRejectUnmappedUsesProtobufNameOverJSONAlias(t *testing.T) {
	t.Parallel()

	// provisioning serializes as json=creation but proto name stays provisioning.
	err := RejectUnmapped("test", &CreateOptions{
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

	if err := RejectUnmapped("test", (*Resources)(nil)); err != nil {
		t.Fatalf("nil remainder: %v", err)
	}
	if err := RejectUnmapped("test", &Resources{}); err != nil {
		t.Fatalf("empty remainder: %v", err)
	}
}
