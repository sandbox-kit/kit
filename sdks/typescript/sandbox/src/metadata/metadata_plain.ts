/**
 * Plain JSON-like values accepted for public metadata maps.
 * Nested objects use an interface so IDEs do not expand the recursive union.
 */
export type MetadataPlain =
  | null
  | boolean
  | string
  | number
  | bigint
  | Uint8Array
  | MetadataPlain[]
  | MetadataPlainObject;

/** Nested object branch of {@link MetadataPlain}. */
export interface MetadataPlainObject {
  [key: string]: MetadataPlain;
}
