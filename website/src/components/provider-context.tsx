"use client";

import {
  createContext,
  useCallback,
  useContext,
  useLayoutEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react";
import {
  DEFAULT_DOC_PROVIDER,
  DOC_PROVIDER_STORAGE_KEY,
  DOC_PROVIDERS,
  isDocProviderId,
  type DocProviderId,
} from "@/lib/providers";

type ProviderContextValue = {
  provider: DocProviderId;
  setProvider: (id: DocProviderId) => void;
  providers: typeof DOC_PROVIDERS;
};

const ProviderContext = createContext<ProviderContextValue | null>(null);

export function ProviderProvider({ children }: { children: ReactNode }) {
  const [provider, setProviderState] =
    useState<DocProviderId>(DEFAULT_DOC_PROVIDER);

  useLayoutEffect(() => {
    try {
      const stored = localStorage.getItem(DOC_PROVIDER_STORAGE_KEY);
      if (stored && isDocProviderId(stored)) {
        setProviderState(stored);
      }
    } catch {
      // ignore
    }
  }, []);

  const setProvider = useCallback((id: DocProviderId) => {
    setProviderState(id);
    try {
      localStorage.setItem(DOC_PROVIDER_STORAGE_KEY, id);
    } catch {
      // ignore
    }
  }, []);

  const value = useMemo(
    () => ({
      provider,
      setProvider,
      providers: DOC_PROVIDERS,
    }),
    [provider, setProvider],
  );

  return (
    <ProviderContext.Provider value={value}>
      {children}
    </ProviderContext.Provider>
  );
}

export function useDocProvider(): ProviderContextValue {
  const ctx = useContext(ProviderContext);
  if (!ctx) {
    throw new Error("useDocProvider must be used within ProviderProvider");
  }
  return ctx;
}

const providerSelectClassName =
  "cursor-pointer rounded-md border border-fd-border bg-fd-background px-2 py-1 text-sm font-medium text-fd-foreground outline-none focus-visible:ring-2 focus-visible:ring-fd-ring";

export function ProviderSelect({ className }: { className?: string }) {
  const { provider, setProvider, providers } = useDocProvider();

  return (
    <select
      aria-label="Provider"
      className={className ?? providerSelectClassName}
      value={provider}
      onChange={(event) => {
        const next = event.target.value;
        if (isDocProviderId(next)) setProvider(next);
      }}
    >
      {providers.map((item) => (
        <option key={item.id} value={item.id}>
          {item.label}
        </option>
      ))}
    </select>
  );
}

/** Pages that are already scoped to one provider — no page-level selector. */
const PROVIDER_SCOPED_PAGES = new Set(["daytona", "modal"]);

export function DocsProviderSelect({ path }: { path: string }) {
  const slug =
    path
      .replace(/\.mdx?$/, "")
      .split("/")
      .pop() ?? "";
  if (PROVIDER_SCOPED_PAGES.has(slug)) return null;
  return <ProviderSelect className={`ms-auto ${providerSelectClassName}`} />;
}

export function ForProvider({
  id,
  children,
}: {
  id: DocProviderId | DocProviderId[];
  children: ReactNode;
}) {
  const { provider } = useDocProvider();
  const ids = Array.isArray(id) ? id : [id];
  if (!ids.includes(provider)) return null;
  return children;
}
