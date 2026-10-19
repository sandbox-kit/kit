import { describe, expect, expectTypeOf, it } from "vitest";
import { durationMs } from "@bufbuild/protobuf/wkt";
import {
  materializeClientConfig,
  materializeCreateOptions,
  type ConfigInput,
} from "./client.js";
import {
  EgressMode,
  type CreateOptionsInput,
  type NetworkConfigInput,
  type SandboxSourceInput,
} from "./sandbox_data_contracts.js";

describe("CreateOptionsInput remaps", () => {
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
    expectTypeOf<
      NonNullable<CreateOptionsInput["environment"]>
    >().toEqualTypeOf<Record<string, string>>();
    expectTypeOf<NonNullable<CreateOptionsInput["labels"]>>().toEqualTypeOf<
      Record<string, string>
    >();
    expectTypeOf<NonNullable<ConfigInput["timeout"]>>().toEqualTypeOf<number>();

    // Wire allowlist / empty-message shapes must not remain public.
    expectTypeOf<{ entries: string[] }>().not.toMatchTypeOf<
      NonNullable<NetworkConfigInput["outboundCidrs"]>
    >();
    expectTypeOf<{}>().not.toMatchTypeOf<
      NonNullable<SandboxSourceInput["providerDefault"]>
    >();
  });

  it("accepts the public shapes and materializes wire allowlists / markers", () => {
    const network = {
      egress: EgressMode.Restricted,
      outboundCidrs: ["192.0.2.0/24"],
      outboundDomains: ["example.com"],
      inboundCidrs: [] as string[],
      privateNetwork: true,
    } satisfies NetworkConfigInput;

    const request = {
      source: { providerDefault: true } satisfies SandboxSourceInput,
      network,
      environment: { A: "1" },
      labels: { k: "v" },
      providerOptions: { flag: true },
      provisioning: { timeout: 1_000 },
    } satisfies CreateOptionsInput;

    const wire = materializeCreateOptions(request);
    expect(wire.source?.kind).toMatchObject({
      case: "providerDefault",
      value: {},
    });
    expect(wire.network?.outboundCidrs?.entries).toEqual(["192.0.2.0/24"]);
    expect(wire.network?.outboundDomains?.entries).toEqual(["example.com"]);
    expect(wire.network?.inboundCidrs?.entries).toEqual([]);
    expect(wire.environment).toEqual({ A: "1" });
    expect(wire.labels).toEqual({ k: "v" });

    const config = {
      auth: { apiKey: { key: "k" } },
      timeout: 30_000,
    } satisfies ConfigInput;
    const wireConfig = materializeClientConfig(config);
    expect(wireConfig.auth?.credentials?.case).toBe("apiKey");
    expect(durationMs(wireConfig.timeout!)).toBe(30_000);
  });
});
