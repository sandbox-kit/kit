type Simplify<T> = { [K in keyof T]: T[K] } & {};

type PreserveOptional<Orig, New> = undefined extends Orig
  ? New | undefined
  : New;

/**
 * Replace selected fields with R[K], keeping each field's optionality from T.
 * Prefer this over call-site `Omit<T, K> & { K?: New }`. Keys in R that are
 * absent from T are ignored; only overlapping keys are rewritten.
 */
export type RemapFields<T, R> = Simplify<{
  [K in keyof T]: K extends keyof R ? PreserveOptional<T[K], R[K]> : T[K];
}>;

type AllowlistField = "outboundCidrs" | "outboundDomains" | "inboundCidrs";

/**
 * Recursively unwrap StringAllowlist-shaped fields to plain `string[]`.
 * Present empty arrays still mean deny-all after materialize.
 */
export type FlattenAllowlistFields<T> = T extends readonly (infer U)[]
  ? FlattenAllowlistFields<U>[]
  : T extends object
    ? {
        [K in keyof T]: K extends AllowlistField
          ? PreserveOptional<T[K], string[]>
          : FlattenAllowlistFields<T[K]>;
      }
    : T;

type EmptyMarkerField = "providerDefault";

/**
 * Recursively turn empty wire marker messages into `true` presence flags
 * (e.g. `source.providerDefault?: true`).
 */
export type EmptyMarkersAsTrue<T> = T extends readonly (infer U)[]
  ? EmptyMarkersAsTrue<U>[]
  : T extends object
    ? {
        [K in keyof T]: K extends EmptyMarkerField
          ? PreserveOptional<T[K], true>
          : EmptyMarkersAsTrue<T[K]>;
      }
    : T;

type StringMapField = "environment" | "labels";

/**
 * Recursively replace noisy proto string maps with `Record<string, string>`.
 */
export type PlainStringMapFields<T> = T extends readonly (infer U)[]
  ? PlainStringMapFields<U>[]
  : T extends object
    ? {
        [K in keyof T]: K extends StringMapField
          ? PreserveOptional<T[K], Record<string, string>>
          : PlainStringMapFields<T[K]>;
      }
    : T;
