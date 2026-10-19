import { describe, expect, it } from "vitest";
import {
  EgressMode as ProtoEgressMode,
  PortTransport as ProtoPortTransport,
} from "@sandbox-kit/contracts/kit/sandbox/v1/sandbox_pb.js";
import { EgressMode, PortTransport } from "./sandbox_data_contracts.js";
import { protoEnum } from "./proto_enum.js";

describe("protoEnum", () => {
  it("maps SCREAMING_SNAKE members to PascalCase with the same values", () => {
    expect(EgressMode).toEqual({
      Default: ProtoEgressMode.DEFAULT,
      AllowAll: ProtoEgressMode.ALLOW_ALL,
      BlockAll: ProtoEgressMode.BLOCK_ALL,
      Restricted: ProtoEgressMode.RESTRICTED,
    });
    expect(PortTransport.Http2Tls).toBe(ProtoPortTransport.HTTP2_TLS);
    expect(PortTransport.Tls).toBe(ProtoPortTransport.TLS);
  });

  it("ignores reverse numeric enum keys", () => {
    const mapped = protoEnum(ProtoEgressMode);
    expect(Object.keys(mapped).every((key) => !/^\d+$/.test(key))).toBe(true);
    expect(Object.keys(mapped)).toHaveLength(4);
  });
});
