import { createFileRoute, notFound } from "@tanstack/react-router";
import { DocsLayout } from "fumadocs-ui/layouts/docs";
import { docs, source } from "@/lib/source";
import {
  DocsBody,
  DocsDescription,
  DocsPage,
  DocsTitle,
  MarkdownCopyButton,
  ViewOptionsPopover,
} from "fumadocs-ui/layouts/docs/page";
import { baseOptions } from "@/lib/layout.shared";
import { absoluteSiteUrl, getPageMarkdownUrl, gitConfig } from "@/lib/shared";
import { socialHead } from "@/lib/seo";

import { Suspense, use, useEffect, useMemo, useState } from "react";
import {
  DocsProviderSelect,
  useDocProvider,
} from "@/components/provider-context";
import { useMDXComponents } from "@/components/mdx";
import { filterTocForProvider } from "@/lib/providers";

export const Route = createFileRoute("/docs/$")({
  component: Page,
  loader: async ({ params }) => {
    const slugs = params._splat?.split("/") ?? [];
    const page = source.getPage(slugs);
    if (!page) throw notFound();
    await docs.getPage(page.path)?.preload();

    return {
      path: page.path,
      pageUrl: page.url,
      title: page.data.title,
      description: page.data.description ?? "",
      // Root-relative without BASE_URL — fumadocs page actions add BASE_URL once.
      markdownUrl: getPageMarkdownUrl(page).url,
    };
  },
  head: ({ loaderData }) => {
    if (!loaderData) return {};
    return socialHead({
      title: loaderData.title,
      description: loaderData.description,
      path: loaderData.pageUrl,
      type: "article",
    });
  },
});

function Content({ path, markdownUrl }: { path: string; markdownUrl: string }) {
  const page = docs.getPage(path);
  if (!page) throw new Error(`unknown page: ${path}`);

  const { toc } = use(page.load());
  const { provider } = useDocProvider();
  const pageSlug =
    path
      .replace(/\.mdx?$/, "")
      .split("/")
      .pop() ?? "";
  const visibleToc = useMemo(
    () => filterTocForProvider(pageSlug, toc, provider),
    [pageSlug, toc, provider],
  );
  const MDX = page.body;
  // AI tools fetch this URL; point at the static .md file, not the SPA HTML page.
  const [markdownPageUrl, setMarkdownPageUrl] = useState(markdownUrl);
  useEffect(() => {
    setMarkdownPageUrl(absoluteSiteUrl(markdownUrl));
  }, [markdownUrl]);

  return (
    <DocsPage toc={visibleToc}>
      <DocsTitle>{page.title}</DocsTitle>
      <DocsDescription>{page.description}</DocsDescription>
      <div className="flex flex-row gap-2 items-center border-b -mt-4 pb-6">
        <MarkdownCopyButton markdownUrl={markdownUrl} />
        <ViewOptionsPopover
          markdownUrl={markdownUrl}
          pageUrl={markdownPageUrl}
          githubUrl={`https://github.com/${gitConfig.user}/${gitConfig.repo}/blob/${gitConfig.branch}/website/content/docs/${path}`}
        />
        <DocsProviderSelect path={path} />
      </div>
      <DocsBody>
        <MDX components={useMDXComponents()} />
      </DocsBody>
    </DocsPage>
  );
}

function Page() {
  const { path, markdownUrl } = Route.useLoaderData();
  // Build the tree in the component (SSR + client). Do not put it through
  // serializePageTree: that imports react-dom/server.edge (CJS `require`), which
  // breaks SPA navigation and has no server to call on GitHub Pages.
  const pageTree = source.getPageTree();

  return (
    <DocsLayout {...baseOptions()} tree={pageTree}>
      <Suspense>
        <Content path={path} markdownUrl={markdownUrl} />
      </Suspense>
    </DocsLayout>
  );
}
