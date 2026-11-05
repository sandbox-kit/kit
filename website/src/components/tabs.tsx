"use client";

import {
  createContext,
  useContext,
  useEffect,
  useId,
  useMemo,
  type ComponentProps,
  type ReactNode,
} from "react";
import {
  Tabs as RootTabs,
  TabsContent as RootTabsContent,
  TabsList as RootTabsList,
  TabsTrigger as RootTabsTrigger,
} from "fumadocs-ui/components/ui/tabs";
import { cn } from "@/lib/cn";
import { DOC_LANGUAGE_STORAGE_KEY } from "@/lib/languages";
import {
  resolveDocLanguageId,
  useDocLanguageOptional,
} from "@/components/language-context";

type TabsContextValue = {
  items?: string[];
  collection: string[];
};

const TabsContext = createContext<TabsContextValue | null>(null);

function useTabContext() {
  const ctx = useContext(TabsContext);
  if (!ctx) throw new Error("You must wrap your component in <Tabs>");
  return ctx;
}

function escapeValue(v: string) {
  return v.toLowerCase().replace(/\s/, "-");
}

type RootTabsProps = ComponentProps<typeof RootTabs>;

export type TabsProps = Omit<RootTabsProps, "value" | "onValueChange"> & {
  items?: string[];
  defaultIndex?: number;
  label?: ReactNode;
};

export function Tabs({
  ref,
  className,
  items,
  label,
  defaultIndex = 0,
  defaultValue = items ? escapeValue(items[defaultIndex]!) : undefined,
  groupId,
  persist,
  children,
  ...props
}: TabsProps) {
  const collection = useMemo(() => [] as string[], []);
  const languageCtx = useDocLanguageOptional();
  const isLanguageGroup = groupId === DOC_LANGUAGE_STORAGE_KEY;
  const languageValue =
    isLanguageGroup && languageCtx ? languageCtx.language : undefined;

  return (
    <RootTabs
      ref={ref}
      className={cn(
        "my-4 flex flex-col overflow-hidden rounded-xl border bg-fd-secondary",
        className,
      )}
      groupId={groupId}
      persist={persist}
      defaultValue={languageValue ? undefined : defaultValue}
      value={languageValue}
      onValueChange={
        isLanguageGroup && languageCtx
          ? (value) => {
              const id = resolveDocLanguageId(value);
              if (id) languageCtx.setLanguage(id);
            }
          : undefined
      }
      {...props}
    >
      {items && (
        <RootTabsList className="flex w-full items-center gap-3.5 overflow-x-auto px-4 text-fd-secondary-foreground not-prose">
          {label && (
            <span className="my-auto me-auto text-sm font-medium">{label}</span>
          )}
          {items.map((item) => (
            <RootTabsTrigger
              key={item}
              value={escapeValue(item)}
              className="inline-flex items-center gap-2 whitespace-nowrap border-b border-transparent py-2 text-sm font-medium text-fd-muted-foreground transition-colors hover:text-fd-accent-foreground disabled:pointer-events-none disabled:opacity-50 data-[active]:border-fd-primary data-[active]:text-fd-primary"
            >
              {item}
            </RootTabsTrigger>
          ))}
        </RootTabsList>
      )}
      <TabsContext.Provider
        value={useMemo(
          () => ({
            items,
            collection,
          }),
          [collection, items],
        )}
      >
        {children}
      </TabsContext.Provider>
    </RootTabs>
  );
}

export type TabProps = Omit<ComponentProps<typeof RootTabsContent>, "value"> & {
  value?: string;
};

export function Tab({ value, ...props }: TabProps) {
  const { items } = useTabContext();
  const resolved = value ?? items?.at(useCollectionIndex());
  if (!resolved) {
    throw new Error(
      "Failed to resolve tab `value`, please pass a `value` prop to the Tab component.",
    );
  }
  return (
    <RootTabsContent
      value={escapeValue(resolved)}
      className="rounded-xl bg-fd-background p-4 text-[0.9375rem] outline-none prose-no-margin data-[inactive]:hidden [&>figure:only-child]:-m-4 [&>figure:only-child]:border-none"
      {...props}
    />
  );
}

function useCollectionIndex() {
  const key = useId();
  const { collection } = useTabContext();
  useEffect(() => {
    return () => {
      const idx = collection.indexOf(key);
      if (idx !== -1) collection.splice(idx, 1);
    };
  }, [key, collection]);
  if (!collection.includes(key)) collection.push(key);
  return collection.indexOf(key);
}
