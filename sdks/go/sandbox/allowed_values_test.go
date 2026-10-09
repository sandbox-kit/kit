package sandbox

import "testing"

func TestUnknownEnumValuesAndMissingReservations(t *testing.T) {
	for _, request := range []*CreateOptions{
		{Isolation: &IsolationConfig{Kind: IsolationKind(99)}},
		{Network: &NetworkConfig{Egress: EgressMode(99)}},
		{Provisioning: &ProvisioningOptions{WaitFor: WaitCondition(99)}},
		{Source: &SandboxSource{Snapshot: &SnapshotSource{Reference: "snapshot", Kind: SnapshotKind(99)}}},
		{Security: &SecurityConfig{Secrets: []*SecretReference{{Reference: "secret", Injection: SecretInjection(99)}}}},
		{Resources: &Resources{CPULimitCores: Value(2.0)}},
		{Resources: &Resources{MemoryLimitMiB: Value(uint64(1024))}},
	} {
		if err := ValidateCreateOptions(request); err == nil {
			t.Fatalf("invalid request accepted: %#v", request)
		}
	}
	if err := ValidateCreateOptions(&CreateOptions{Resources: &Resources{CPULimitCores: Value(0.0), MemoryLimitMiB: Value(uint64(0))}}); err != nil {
		t.Fatal("explicit no-limit values rejected", err)
	}
}
