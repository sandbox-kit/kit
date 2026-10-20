import { describe, expectTypeOf, it } from "vitest";
import { type Config } from "./client.js";
import {
  type CreateOptions,
  type NetworkConfigInput,
  type SandboxSourceInput,
} from "./sandbox_data_contracts.js";

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

    // Wire allowlist / empty-message shapes must not remain public.
    expectTypeOf<{ entries: string[] }>().not.toMatchTypeOf<
      NonNullable<NetworkConfigInput["outboundCidrs"]>
    >();
    expectTypeOf<{}>().not.toMatchTypeOf<
      NonNullable<SandboxSourceInput["providerDefault"]>
    >();
  });
});
