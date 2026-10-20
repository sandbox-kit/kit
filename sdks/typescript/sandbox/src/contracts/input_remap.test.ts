import { describe, expectTypeOf, it } from "vitest";
import { type Config } from "#client/client.js";
import {
  type AutomaticActionInput,
  type CreateOptions,
  type LifetimePolicyInput,
  type NetworkConfigInput,
  type ObservabilityConfig,
  type ProvisioningOptionsInput,
  type SandboxSourceInput,
  type SecurityConfig,
} from "#contracts/sandbox_data_contracts.js";

describe("CreateOptions remaps", () => {
  it("types allowlists as string[], providerDefault as true, maps as Record", () => {
    expectTypeOf<
      NonNullable<NetworkConfigInput["outboundCidrs"]>
    >().toEqualTypeOf<string[]>();
    expectTypeOf<
      NonNullable<NetworkConfigInput["outboundDomains"]>
    >().toEqualTypeOf<string[]>();
    expectTypeOf<
      NonNullable<NetworkConfigInput["inboundCidrs"]>
    >().toEqualTypeOf<string[]>();
    expectTypeOf<
      NonNullable<SandboxSourceInput["providerDefault"]>
    >().toEqualTypeOf<true>();
    expectTypeOf<NonNullable<CreateOptions["environment"]>>().toEqualTypeOf<
      Record<string, string>
    >();
    expectTypeOf<NonNullable<CreateOptions["labels"]>>().toEqualTypeOf<
      Record<string, string>
    >();
    expectTypeOf<NonNullable<Config["timeout"]>>().toEqualTypeOf<number>();
    expectTypeOf<
      NonNullable<LifetimePolicyInput["maximumLifetime"]>
    >().toEqualTypeOf<number>();
    expectTypeOf<
      NonNullable<AutomaticActionInput["after"]>
    >().toEqualTypeOf<number>();
    expectTypeOf<
      NonNullable<ProvisioningOptionsInput["timeout"]>
    >().toEqualTypeOf<number>();
    expectTypeOf<
      NonNullable<CreateOptions["security"]>
    >().toEqualTypeOf<SecurityConfig>();
    expectTypeOf<
      NonNullable<CreateOptions["observability"]>
    >().toEqualTypeOf<ObservabilityConfig>();

    // Wire allowlist / empty-message shapes must not remain public.
    expectTypeOf<{ entries: string[] }>().not.toMatchTypeOf<
      NonNullable<NetworkConfigInput["outboundCidrs"]>
    >();
    expectTypeOf<{}>().not.toMatchTypeOf<
      NonNullable<SandboxSourceInput["providerDefault"]>
    >();
  });
});
