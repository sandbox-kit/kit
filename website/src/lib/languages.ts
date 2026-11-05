export const DOC_LANGUAGES = [
  { id: "typescript", label: "TypeScript" },
  { id: "go", label: "Go" },
] as const;

export type DocLanguageId = (typeof DOC_LANGUAGES)[number]["id"];

export const DEFAULT_DOC_LANGUAGE: DocLanguageId = "typescript";

/** Matches `<Tabs groupId="language" persist>` storage key in fumadocs tabs. */
export const DOC_LANGUAGE_STORAGE_KEY = "language";

export function isDocLanguageId(value: string): value is DocLanguageId {
  return DOC_LANGUAGES.some((language) => language.id === value);
}
