import defaultMdxComponents from "fumadocs-ui/mdx";
import type { MDXComponents } from "mdx/types";
import { ForLanguage } from "@/components/language-context";
import { Note } from "@/components/note";
import { ForProvider } from "@/components/provider-context";
import { Tab, Tabs } from "@/components/tabs";

export function getMDXComponents(components?: MDXComponents) {
  return {
    ...defaultMdxComponents,
    Tabs,
    Tab,
    ForProvider,
    ForLanguage,
    Note,
    ...components,
  } satisfies MDXComponents;
}

export const useMDXComponents = getMDXComponents;

declare global {
  type MDXProvidedComponents = ReturnType<typeof getMDXComponents>;
}
