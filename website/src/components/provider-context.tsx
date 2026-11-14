"use client";

import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import { Check, ChevronDown } from "lucide-react";
import { buttonVariants } from "fumadocs-ui/components/ui/button";
import { cn } from "@/lib/cn";
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

export function ProviderSelect({ className }: { className?: string }) {
  const { provider, setProvider, providers } = useDocProvider();
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  const current =
    providers.find((item) => item.id === provider) ?? providers[0]!;

  useEffect(() => {
    if (!open) return;
    const onPointerDown = (event: PointerEvent) => {
      if (!rootRef.current?.contains(event.target as Node)) setOpen(false);
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setOpen(false);
    };
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  return (
    <div ref={rootRef} className={cn("relative", className)}>
      <button
        type="button"
        aria-haspopup="listbox"
        aria-expanded={open}
        className={cn(
          buttonVariants({
            variant: "secondary",
            size: "sm",
          }),
          "gap-1.5",
          open && "bg-fd-accent text-fd-accent-foreground",
        )}
        onClick={() => setOpen((value) => !value)}
      >
        {current.label}
        <ChevronDown className="size-3.5 text-fd-muted-foreground" />
      </button>
      {open && (
        <div
          role="listbox"
          aria-label="Provider"
          className="absolute end-0 top-[calc(100%+0.35rem)] z-50 min-w-40 rounded-xl border border-fd-border bg-fd-popover p-1 text-fd-popover-foreground shadow-md"
        >
          {providers.map((item) => {
            const selected = item.id === provider;
            return (
              <button
                key={item.id}
                type="button"
                role="option"
                aria-selected={selected}
                className={cn(
                  "flex w-full items-center gap-2 rounded-lg px-2.5 py-1.5 text-left text-sm hover:bg-fd-accent hover:text-fd-accent-foreground",
                  selected && "bg-fd-accent text-fd-accent-foreground",
                )}
                onClick={() => {
                  setProvider(item.id);
                  setOpen(false);
                }}
              >
                <span>{item.label}</span>
                {selected ? (
                  <Check className="ms-auto size-3.5 text-fd-muted-foreground" />
                ) : (
                  <span className="ms-auto size-3.5" aria-hidden />
                )}
              </button>
            );
          })}
        </div>
      )}
    </div>
  );
}

/** Pages that are already scoped to one provider — no page-level selector. */
const PROVIDER_SCOPED_PAGES = new Set(["daytona", "modal", "blaxel"]);

export function DocsProviderSelect({ path }: { path: string }) {
  const slug =
    path
      .replace(/\.mdx?$/, "")
      .split("/")
      .pop() ?? "";
  if (PROVIDER_SCOPED_PAGES.has(slug)) return null;
  return <ProviderSelect className="ms-auto" />;
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
