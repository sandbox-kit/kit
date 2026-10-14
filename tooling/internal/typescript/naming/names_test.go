package naming

import "testing"

func TestNames(t *testing.T) {
	if got, want := Camel("provider_state"), "providerState"; got != want {
		t.Fatalf("Camel() = %q, want %q", got, want)
	}
	if got, want := Pascal("soft-deleted"), "SoftDeleted"; got != want {
		t.Fatalf("Pascal() = %q, want %q", got, want)
	}
	if got, want := EnumPrefix("ErrorKind"), "ERROR_KIND_"; got != want {
		t.Fatalf("EnumPrefix() = %q, want %q", got, want)
	}
}
