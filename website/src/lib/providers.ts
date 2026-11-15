export const DOC_PROVIDERS = [
  { id: "daytona", label: "Daytona" },
  { id: "modal", label: "Modal" },
  { id: "blaxel", label: "Blaxel" },
] as const;

export type DocProviderId = (typeof DOC_PROVIDERS)[number]["id"];

export const DEFAULT_DOC_PROVIDER: DocProviderId = "daytona";

export const DOC_PROVIDER_STORAGE_KEY = "sandbox-kit.docs.provider";

export function isDocProviderId(value: string): value is DocProviderId {
  return DOC_PROVIDERS.some((provider) => provider.id === value);
}

/**
 * TOC entries that only exist for some providers. Keyed by docs page slug
 * (filename without extension), then by heading id. Keep in sync with
 * headings wrapped in `<ForProvider>` in `content/docs`.
 */
export const PROVIDER_SCOPED_HEADINGS: Record<
  string,
  Record<string, readonly DocProviderId[]>
> = {
  lifecycle: {
    "stop-start-pause-and-archive": ["daytona"],
    delete: ["daytona", "blaxel"],
    terminate: ["modal"],
  },
  client: {
    get: ["daytona", "blaxel"],
    attach: ["modal"],
  },
  sandbox: {
    refresh: ["daytona", "blaxel"],
    start: ["daytona"],
    stop: ["daytona"],
    pause: ["daytona"],
    archive: ["daytona"],
    delete: ["daytona", "blaxel"],
    terminate: ["modal"],
  },
};

export function filterTocForProvider<T extends { url: string }>(
  pageSlug: string,
  toc: T[],
  provider: DocProviderId,
): T[] {
  const scoped = PROVIDER_SCOPED_HEADINGS[pageSlug];
  if (!scoped) return toc;

  return toc.filter((item) => {
    const id = item.url.replace(/^#/, "");
    const providers = scoped[id];
    if (!providers) return true;
    return providers.includes(provider);
  });
}
