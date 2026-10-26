package sandbox

import "testing"

func TestValidValueOrigin(t *testing.T) {
	t.Parallel()

	for _, origin := range []ValueOrigin{ValueOriginUnknown, ValueOriginProvider, ValueOriginRequest} {
		if !validValueOrigin(origin) {
			t.Fatalf("validValueOrigin(%v) = false, want true", origin)
		}
	}
	if validValueOrigin(ValueOrigin(99)) {
		t.Fatal("validValueOrigin(99) = true, want false")
	}
}

func TestResponseFieldPresent(t *testing.T) {
	t.Parallel()

	info := &SandboxInfo{
		Id:       "sbx-1",
		Provider: "modal",
		Resources: &Resources{
			CpuCores: Value(0.5),
		},
		Isolation: &IsolationConfig{},
		Labels:    map[string]string{},
		State:     SandboxStateRunning,
	}

	cases := []struct {
		path string
		want bool
	}{
		{InfoFieldId, true},
		{InfoFieldProvider, true},
		{InfoFieldResources, true},
		{InfoFieldResourcesCpuCores, true},
		{InfoFieldResourcesMemoryMib, false},
		{InfoFieldIsolation, true},
		{InfoFieldIsolationKind, true},
		{InfoFieldIsolationProviderKind, false},
		{InfoFieldLabels, true},
		{InfoFieldState, true},
		{InfoFieldName, false},
		{"not_a_field", false},
		{"resources.not_a_field", false},
		{"", false},
	}
	for _, tc := range cases {
		if got := responseFieldPresent(info, tc.path); got != tc.want {
			t.Fatalf("responseFieldPresent(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
	if responseFieldPresent(nil, InfoFieldId) {
		t.Fatal("responseFieldPresent(nil) = true, want false")
	}
	emptyID := &SandboxInfo{Provider: "modal"}
	if responseFieldPresent(emptyID, InfoFieldId) {
		t.Fatal("empty id must not count as present")
	}
}

func TestValidateResponseOrigins(t *testing.T) {
	t.Parallel()

	ok := &SandboxInfo{
		Id:       "sbx-1",
		Provider: "modal",
		Resources: &Resources{
			CpuCores: Value(1.0),
		},
		Origins: map[string]ValueOrigin{
			InfoFieldId:                ValueOriginProvider,
			InfoFieldResourcesCpuCores: ValueOriginRequest,
			InfoFieldState:             ValueOriginProvider,
		},
	}
	if err := validateResponseOrigins(ok, "create"); err != nil {
		t.Fatalf("valid origins: %v", err)
	}

	missingField := &SandboxInfo{
		Id:       "sbx-1",
		Provider: "modal",
		Origins: map[string]ValueOrigin{
			InfoFieldResourcesCpuCores: ValueOriginProvider,
		},
	}
	if err := validateResponseOrigins(missingField, "create"); err == nil {
		t.Fatal("expected error for origin on absent field")
	}

	emptyID := &SandboxInfo{
		Provider: "modal",
		Origins: map[string]ValueOrigin{
			InfoFieldId: ValueOriginProvider,
		},
	}
	if err := validateResponseOrigins(emptyID, "create"); err == nil {
		t.Fatal("expected error for origin on empty id")
	}

	badOrigin := &SandboxInfo{
		Id:       "sbx-1",
		Provider: "modal",
		Origins: map[string]ValueOrigin{
			InfoFieldId: ValueOrigin(99),
		},
	}
	if err := validateResponseOrigins(badOrigin, "create"); err == nil {
		t.Fatal("expected error for unknown ValueOrigin")
	}
}
