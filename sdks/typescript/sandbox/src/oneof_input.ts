/**
 * Derive nested exclusive input fields from a protobuf-es `{ case, value }` oneof.
 * Message values become init-shaped (fields optional, no `$typeName`).
 */
export type FlattenOneof<
  T extends { case: string | undefined; value?: unknown },
> = {
  [K in Exclude<T["case"], undefined>]?: OneofInputValue<
    Extract<T, { case: K }>["value"]
  >;
};

type OneofInputValue<V> = V extends
  | string
  | number
  | boolean
  | bigint
  | null
  | undefined
  | Uint8Array
  | Date
  ? V
  : V extends readonly (infer E)[]
    ? E[]
    : V extends object
      ? {
          [P in keyof V as P extends "$typeName" | "$unknown"
            ? never
            : P]?: OneofInputValue<V[P]>;
        }
      : V;
