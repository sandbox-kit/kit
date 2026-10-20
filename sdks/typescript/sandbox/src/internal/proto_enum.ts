type Simplify<T> = { [K in keyof T]: T[K] } & {};

/** `ALLOW_ALL` → `AllowAll`, `HTTP2_TLS` → `Http2Tls`. */
type SnakeToPascal<S extends string> = S extends `${infer Head}_${infer Tail}`
  ? `${Capitalize<Lowercase<Head>>}${SnakeToPascal<Tail>}`
  : Capitalize<Lowercase<S>>;

/** Forward protobuf-es enum member names (skip reverse numeric keys). */
type EnumMemberKeys<E> = {
  [K in keyof E]: K extends string
    ? E[K] extends number
      ? K extends `${number}`
        ? never
        : K
      : never
    : never;
}[keyof E];

/**
 * Public const object for a protobuf-es enum: PascalCase keys, same numeric
 * values. Does not invent a second contract — only renames member keys.
 */
export type EnumConst<E> = Simplify<{
  readonly [K in EnumMemberKeys<E> as SnakeToPascal<K & string>]: E[K];
}>;

function snakeToPascal(name: string): string {
  return name
    .toLowerCase()
    .split("_")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join("");
}

/**
 * Build an ergonomic enum const from a protobuf-es numeric enum object.
 *
 * @example
 * export const EgressMode = protoEnum(ProtoEgressMode);
 * // { Default: 0, AllowAll: 1, BlockAll: 2, Restricted: 3 }
 */
export function protoEnum<E extends Record<string, string | number>>(
  proto: E,
): EnumConst<E> {
  const out: Record<string, number> = {};
  for (const [key, value] of Object.entries(proto)) {
    if (typeof value !== "number") continue;
    if (/^\d+$/.test(key)) continue;
    out[snakeToPascal(key)] = value;
  }
  return out as EnumConst<E>;
}
