import { absoluteSiteUrl, appName } from "./shared";

export const siteDescription =
  "SDK for cloud sandboxes in TypeScript or Go. Plug in a provider; the client talks to it for create, files, and commands.";

/** Same asset as the site favicon/nav mark (`public/logo.png`). */
export const ogImagePath = "/logo.png";
export const ogImageMime = "image/png";
export const ogImageSize = { width: 500, height: 500 } as const;

export type SocialHeadInput = {
  title: string;
  description: string;
  /** Site path such as `/` or `/docs/create` (BASE_URL applied once). */
  path: string;
  type?: "website" | "article";
};

/** Open Graph + Twitter Card tags with absolute URLs for static GitHub Pages. */
export function socialHead({
  title,
  description,
  path,
  type = "website",
}: SocialHeadInput) {
  const url = absoluteSiteUrl(path);
  const image = absoluteSiteUrl(ogImagePath);
  const fullTitle = title === appName ? title : `${title} | ${appName}`;

  return {
    meta: [
      { title: fullTitle },
      { name: "description", content: description },
      { property: "og:title", content: fullTitle },
      { property: "og:description", content: description },
      { property: "og:url", content: url },
      { property: "og:type", content: type },
      { property: "og:site_name", content: appName },
      { property: "og:locale", content: "en_US" },
      { property: "og:image", content: image },
      { property: "og:image:secure_url", content: image },
      { property: "og:image:type", content: ogImageMime },
      { property: "og:image:width", content: String(ogImageSize.width) },
      { property: "og:image:height", content: String(ogImageSize.height) },
      { property: "og:image:alt", content: appName },
      { name: "twitter:card", content: "summary" },
      { name: "twitter:title", content: fullTitle },
      { name: "twitter:description", content: description },
      { name: "twitter:image", content: image },
      { name: "twitter:image:alt", content: appName },
    ],
    links: [{ rel: "canonical", href: url }],
  };
}

export function siteSocialHead() {
  return socialHead({
    title: appName,
    description: siteDescription,
    path: "/",
  });
}
