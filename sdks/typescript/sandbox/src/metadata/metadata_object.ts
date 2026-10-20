// Handwritten helpers for protobuf metadata objects (mirrors Go sandbox metadata helpers).

import { create } from "@bufbuild/protobuf";
import {
  MetadataListSchema,
  MetadataObjectSchema,
  MetadataValueSchema,
} from "@sandbox-kit/contracts/sandbox";
import type {
  MetadataObject,
  MetadataValue,
} from "#contracts/sandbox_data_contracts.js";
import type { MetadataPlain } from "#metadata/metadata_plain.js";

export type { MetadataPlain } from "#metadata/metadata_plain.js";

/** Converts a declarative map into the protobuf metadata object. */
export function metadataObjectFromMap(
  source: Record<string, MetadataPlain> | undefined,
): MetadataObject | undefined {
  if (source === undefined) return undefined;
  const values: { [key: string]: MetadataValue } = {};
  for (const [key, value] of Object.entries(source)) {
    values[key] = metadataValueFromPlain(value);
  }
  return create(MetadataObjectSchema, { values }) as MetadataObject;
}

/** Converts a protobuf metadata object into a declarative map. */
export function metadataObjectAsMap(
  source: MetadataObject | undefined,
): Record<string, MetadataPlain> | undefined {
  if (source === undefined) return undefined;
  const out: Record<string, MetadataPlain> = {};
  for (const [key, value] of Object.entries(source.values ?? {})) {
    out[key] = metadataValueAsPlain(value);
  }
  return out;
}

function metadataValueFromPlain(value: MetadataPlain): MetadataValue {
  // create() returns a protobuf Message; public aliases are StripProto shapes.
  if (value === null) return create(MetadataValueSchema, {}) as MetadataValue;
  if (typeof value === "boolean")
    return create(MetadataValueSchema, {
      value: { case: "boolean", value },
    }) as MetadataValue;
  if (typeof value === "string")
    return create(MetadataValueSchema, {
      value: { case: "text", value },
    }) as MetadataValue;
  if (typeof value === "bigint") {
    if (value < 0n)
      return create(MetadataValueSchema, {
        value: { case: "signedInteger", value },
      }) as MetadataValue;
    return create(MetadataValueSchema, {
      value: { case: "unsignedInteger", value },
    }) as MetadataValue;
  }
  if (typeof value === "number") {
    if (!Number.isFinite(value)) {
      throw new Error("sandbox-kit: metadata numbers must be finite");
    }
    if (Number.isInteger(value)) {
      const asBigInt = BigInt(value);
      if (value < 0)
        return create(MetadataValueSchema, {
          value: { case: "signedInteger", value: asBigInt },
        }) as MetadataValue;
      return create(MetadataValueSchema, {
        value: { case: "unsignedInteger", value: asBigInt },
      }) as MetadataValue;
    }
    return create(MetadataValueSchema, {
      value: { case: "number", value },
    }) as MetadataValue;
  }
  if (value instanceof Uint8Array)
    return create(MetadataValueSchema, {
      value: { case: "binary", value },
    }) as MetadataValue;
  if (Array.isArray(value)) {
    return create(MetadataValueSchema, {
      value: {
        case: "list",
        value: create(MetadataListSchema, {
          values: value.map(metadataValueFromPlain),
        }),
      },
    }) as MetadataValue;
  }
  return create(MetadataValueSchema, {
    value: {
      case: "object",
      value: metadataObjectFromMap(value)!,
    },
  }) as MetadataValue;
}

function metadataValueAsPlain(value: MetadataValue): MetadataPlain {
  switch (value.value?.case) {
    case undefined:
      return null;
    case "boolean":
    case "text":
    case "number":
    case "binary":
      return value.value.value as MetadataPlain;
    case "signedInteger":
    case "unsignedInteger":
      return value.value.value as MetadataPlain;
    case "list":
      return (value.value.value?.values ?? []).map(metadataValueAsPlain);
    case "object":
      return metadataObjectAsMap(value.value.value) ?? {};
    default:
      return null;
  }
}
