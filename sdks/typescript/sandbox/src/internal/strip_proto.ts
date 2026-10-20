import type { Duration, Timestamp } from "@bufbuild/protobuf/wkt";

type ProtoKey = "$typeName" | "$unknown";

/**
 * True for proto3 map fields (`{ [key: string]: V }`), including non-string
 * values such as `map<string, ValueOrigin>`. Not true for plain messages.
 */
type IsProtoMap<T> = T extends object
  ? string extends keyof T
    ? T extends readonly unknown[]
      ? false
      : true
    : false
  : false;

/**
 * Required proto3 maps/repeated/enums become optional for construction.
 * protobuf-es enums are numeric, so required `number` fields soften too.
 * Map *entry* keys (index signatures) are not softened — that would yield
 * `Record<string, V | undefined>`.
 */
type SoftenField<V, K> = string extends K
  ? false
  : V extends readonly unknown[]
    ? true
    : IsProtoMap<V> extends true
      ? true
      : V extends number
        ? true
        : false;

type OptionalKeys<T> = Exclude<
  {
    [K in keyof T]-?: undefined extends T[K] ? K : never;
  }[keyof T],
  ProtoKey
>;

type RequiredKeys<T> = Exclude<keyof T, OptionalKeys<T> | ProtoKey>;

type SoftenedKeys<T> = {
  [K in RequiredKeys<T>]: SoftenField<T[K], K> extends true ? K : never;
}[RequiredKeys<T>];

type KeptRequiredKeys<T> = Exclude<RequiredKeys<T>, SoftenedKeys<T>>;

type Simplify<T> = { [K in keyof T]: T[K] } & {};

type StripProtoObject<T> = Simplify<
  {
    [K in OptionalKeys<T> & keyof T]?: StripProto<Exclude<T[K], undefined>>;
  } & {
    [K in SoftenedKeys<T>]?: StripProto<T[K]>;
  } & {
    [K in KeptRequiredKeys<T>]: StripProto<T[K]>;
  }
>;

/**
 * Recursively drop protobuf-es message internals for public SDK types:
 *
 * - removes `$typeName` / `$unknown` at every level (including nested messages)
 * - keeps normal field optionality
 * - softens required proto3 maps/repeated/enums to optional
 * - keeps wire `Duration` opaque for provider helpers
 * - maps `Timestamp` to `Date`
 *
 * Nested messages are stripped in place, so aliases do not need `Omit` + remap
 * just to clean child protobuf fields.
 */
export type StripProto<T> = T extends Duration
  ? Duration
  : T extends Timestamp
    ? Date
    : T extends
          | string
          | number
          | boolean
          | bigint
          | null
          | undefined
          | symbol
          | Date
          | Uint8Array
          | ((...args: never[]) => unknown)
      ? T
      : T extends readonly (infer U)[]
        ? StripProto<U>[]
        : T extends object
          ? StripProtoObject<T>
          : T;
