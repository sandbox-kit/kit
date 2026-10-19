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
 * left alone.
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

const durationFieldNames = new Set<string>([
  "timeout",
  "queueTimeout",
  "maximumLifetime",
  "after",
]);

/** Convert public ms numbers on known duration fields into protobuf Durations. */
export function materializeDurationFields<T>(value: T): T {
  return materializeDurationFieldsAt(value, []) as T;
}

function materializeDurationFieldsAt(value: unknown, path: string[]): unknown {
  if (value === null || value === undefined) return value;
  if (typeof value === "number") {
    const leaf = path[path.length - 1];
    if (leaf !== undefined && durationFieldNames.has(leaf)) {
      return durationFromMs(value);
    }
    return value;
  }
  if (typeof value !== "object") return value;
  if (isDurationValue(value)) return value;
  if (Array.isArray(value)) {
    return value.map((item) => materializeDurationFieldsAt(item, path));
  }
  const out: Record<string, unknown> = {};
  for (const [key, child] of Object.entries(value as Record<string, unknown>)) {
    out[key] = materializeDurationFieldsAt(child, [...path, key]);
  }
  return out;
}
