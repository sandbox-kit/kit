import { createGetUrl } from "fumadocs-core/source";

export const appName = "Sandbox Kit";
export const docsRoute = "/docs";
/** Production origin for absolute OG/canonical URLs on GitHub Pages. */
export const siteOrigin = "https://sandbox-kit.github.io";

export const gitConfig = {
  user: "sandbox-kit",
  repo: "kit",
  branch: "main",
};

const getDocsUrl = createGetUrl(docsRoute);

export function getPageMarkdownUrl(page: { slugs: string[]; locale?: string }) {
  const segments = [...page.slugs];
  if (segments.length === 0) {
    segments.push("index.md");
  } else {
    segments[segments.length - 1] += ".md";
  }

  return { segments, url: getDocsUrl(segments, page.locale) };
}

/** @returns page slugs */
export function decodeMarkdownUrl(segments: string[]) {
  if (segments.length === 0) return [];

  const out = [...segments];
  out[out.length - 1] = out[out.length - 1].replace(/\.md$/, "");
  if (out.length === 1 && out[0] === "index") out.pop();
  return out;
}

/** Public file URL. Honors the GitHub Pages base (`/kit/`) when the site is built for it. */
export function assetUrl(path: string) {
  const base = import.meta.env.BASE_URL ?? "/";
  return `${base}${path.replace(/^\//, "")}`;
}

/** Absolute URL for a root-relative site path (applies `BASE_URL` once). */
export function absoluteSiteUrl(path: string) {
  const base = (import.meta.env.BASE_URL ?? "/").replace(/\/$/, "");
  const normalized = path.startsWith("/") ? path : `/${path}`;
  const withBase =
    !base || normalized === base || normalized.startsWith(`${base}/`)
      ? normalized
      : `${base}${normalized}`;
  if (typeof window !== "undefined") {
    return `${window.location.origin}${withBase}`;
  }
  // Static prerender for GitHub Pages project site.
  if (base === "/kit") {
    return `${siteOrigin}${withBase}`;
  }
  return withBase;
}
