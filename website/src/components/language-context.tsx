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
  DEFAULT_DOC_LANGUAGE,
  DOC_LANGUAGE_STORAGE_KEY,
  DOC_LANGUAGES,
  isDocLanguageId,
  type DocLanguageId,
} from "@/lib/languages";

type LanguageContextValue = {
  language: DocLanguageId;
  setLanguage: (id: DocLanguageId) => void;
  languages: typeof DOC_LANGUAGES;
};

const LanguageContext = createContext<LanguageContextValue | null>(null);

function readStoredLanguage(): DocLanguageId {
  if (typeof window === "undefined") return DEFAULT_DOC_LANGUAGE;
  try {
    const stored =
      localStorage.getItem(DOC_LANGUAGE_STORAGE_KEY) ??
      sessionStorage.getItem(DOC_LANGUAGE_STORAGE_KEY);
    if (stored && isDocLanguageId(stored)) return stored;
  } catch {
    // ignore
  }
  return DEFAULT_DOC_LANGUAGE;
}

export function LanguageProvider({ children }: { children: ReactNode }) {
  const [language, setLanguageState] =
    useState<DocLanguageId>(DEFAULT_DOC_LANGUAGE);

  useLayoutEffect(() => {
    setLanguageState(readStoredLanguage());
  }, []);

  const setLanguage = useCallback((id: DocLanguageId) => {
    setLanguageState(id);
    try {
      sessionStorage.setItem(DOC_LANGUAGE_STORAGE_KEY, id);
      localStorage.setItem(DOC_LANGUAGE_STORAGE_KEY, id);
    } catch {
      // ignore
    }
  }, []);

  const value = useMemo(
    () => ({
      language,
      setLanguage,
      languages: DOC_LANGUAGES,
    }),
    [language, setLanguage],
  );

  return (
    <LanguageContext.Provider value={value}>
      {children}
    </LanguageContext.Provider>
  );
}

export function useDocLanguage(): LanguageContextValue {
  const ctx = useContext(LanguageContext);
  if (!ctx) {
    throw new Error("useDocLanguage must be used within LanguageProvider");
  }
  return ctx;
}

export function useDocLanguageOptional(): LanguageContextValue | null {
  return useContext(LanguageContext);
}

export function resolveDocLanguageId(value: string): DocLanguageId | null {
  const normalized = value.toLowerCase().replace(/\s+/g, "-");
  if (isDocLanguageId(normalized)) return normalized;
  return null;
}

/** Show children only for the active docs language. No language switcher. */
export function ForLanguage({
  id,
  children,
}: {
  id: DocLanguageId | DocLanguageId[] | string | string[];
  children: ReactNode;
}) {
  const { language } = useDocLanguage();
  const raw = Array.isArray(id) ? id : [id];
  const ids = raw
    .map((item) => resolveDocLanguageId(String(item)))
    .filter((item): item is DocLanguageId => item !== null);
  if (!ids.includes(language)) return null;
  return children;
}
