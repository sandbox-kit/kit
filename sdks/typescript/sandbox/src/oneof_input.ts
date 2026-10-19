import type { StripProto } from "./strip_proto.js";

type Simplify<T> = { [K in keyof T]: T[K] } & {};

type ProtoOneof = { case: string | undefined; value?: unknown };

/**
 * Derive nested exclusive input fields from a protobuf-es `{ case, value }` oneof.
 * Variant payloads are StripProto'd (no `$typeName` / `$unknown`).
 */
export type FlattenOneof<T extends ProtoOneof> = Simplify<{
  [K in Exclude<T["case"], undefined>]?: StripProto<
    Exclude<Extract<T, { case: K }>["value"], undefined>
  >;
}>;

/**
 * Strip protobuf-es internals, then replace listed fields with flattened oneof
 * bags. Use when a wire field is a carrier message whose oneof should surface
 * as exclusive keys on the public input (e.g. `source.kind` → `source.image`).
 *
 * `Oneofs` maps each remapped field to the oneof union to flatten.
 * Extra remaps (plain metadata maps, ms durations) stay at the call site.
 */
export type StripAndFlatten<
  T,
  Oneofs extends {
    [K in keyof Oneofs]: K extends keyof T ? ProtoOneof : never;
  },
> = Simplify<
  Omit<StripProto<T>, keyof Oneofs & string> & {
    [K in keyof Oneofs & string]?: FlattenOneof<Oneofs[K]>;
  }
>;
