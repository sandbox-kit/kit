export const DOC_PROVIDERS = [
  { id: "daytona", label: "Daytona" },
  { id: "modal", label: "Modal" },
] as const;

export type DocProviderId = (typeof DOC_PROVIDERS)[number]["id"];

export const DEFAULT_DOC_PROVIDER: DocProviderId = "daytona";

export const DOC_PROVIDER_STORAGE_KEY = "sandbox-kit.docs.provider";

export function isDocProviderId(value: string): value is DocProviderId {
  return DOC_PROVIDERS.some((provider) => provider.id === value);
}
