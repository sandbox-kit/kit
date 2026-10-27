import { durationFromMs, type Duration } from "@bufbuild/protobuf/wkt";

/** True when a value looks like a protobuf Duration message or init. */
export function isDurationValue(value: unknown): value is Duration {
  return (
    typeof value === "object" &&
    value !== null &&
    ("seconds" in value || "nanos" in value) &&
    !Array.isArray(value)
  );
}

type DurationFieldName =
  | "timeout"
  | "queueTimeout"
  | "maximumLifetime"
  | "after";

/**
 * Recursively replace known protobuf Duration fields with millisecond numbers
 * for public request construction. Field-name based so maps like `labels` are
 * left alone. Public values are always milliseconds.
 */
export type DurationFieldsToMs<T> = T extends readonly (infer U)[]
  ? DurationFieldsToMs<U>[]
  : T extends object
    ? {
        [K in keyof T]: K extends DurationFieldName
          ? undefined extends T[K]
            ? number | undefined
            : number
          : DurationFieldsToMs<T[K]>;
      }
    : T;

function isPlainRecord(value: object): value is Record<string, unknown> {
  const prototype = Object.getPrototypeOf(value);
  return prototype === Object.prototype || prototype === null;
}

const durationFieldNames = new Set<string>([
  "timeout",
  "queueTimeout",
  "maximumLifetime",
  "after",
]);

/**
 * Convert public ms numbers on known duration fields into protobuf Durations.
 * Subtrees with no duration field are returned unchanged.
 */
export function materializeDurationFields<T>(value: T): T {
  return materializeDurationFieldsAt(value, undefined) as T;
}

function materializeDurationFieldsAt(
  value: unknown,
  key: string | undefined,
): unknown {
  if (value === null || value === undefined) return value;
  if (typeof value === "number") {
    if (key !== undefined && durationFieldNames.has(key)) {
      return durationFromMs(value);
    }
    return value;
  }
  if (typeof value !== "object") return value;
  if (isDurationValue(value)) return value;
  if (!isPlainRecord(value) && !Array.isArray(value)) {
    const out: Record<string, unknown> = {};
    for (const [childKey, child] of Object.entries(value)) {
      out[childKey] = materializeDurationFieldsAt(child, childKey);
    }
    return out;
  }
  if (Array.isArray(value)) {
    let copy: unknown[] | undefined;
    for (let index = 0; index < value.length; index++) {
      const item = value[index];
      const next = materializeDurationFieldsAt(item, key);
      if (next !== item) {
        copy ??= value.slice();
        copy[index] = next;
      }
    }
    return copy ?? value;
  }
  const record = value as Record<string, unknown>;
  let copy: Record<string, unknown> | undefined;
  for (const [childKey, child] of Object.entries(record)) {
    const next = materializeDurationFieldsAt(child, childKey);
    if (next !== child) {
      copy ??= { ...record };
      copy[childKey] = next;
    }
  }
  return copy ?? value;
}
