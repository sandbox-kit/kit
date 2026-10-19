/** Plain JSON-like values accepted for public metadata / providerOptions maps. */
export type MetadataPlain =
  | null
  | boolean
  | string
  | number
  | bigint
  | Uint8Array
  | MetadataPlain[]
  | { [key: string]: MetadataPlain };
