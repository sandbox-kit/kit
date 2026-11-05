import {
  createRootRoute,
  HeadContent,
  Outlet,
  Scripts,
} from "@tanstack/react-router";
import * as React from "react";
import appCss from "@/styles/app.css?url";
import { RootProvider } from "fumadocs-ui/provider/tanstack";
import { ProviderProvider } from "@/components/provider-context";
import { assetUrl } from "@/lib/shared";
import { siteDescription } from "@/lib/seo";

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      {
        name: "viewport",
        content: "width=device-width, initial-scale=1",
      },
      { title: "Sandbox Kit" },
      { name: "description", content: siteDescription },
    ],
    links: [
      { rel: "stylesheet", href: appCss },
      { rel: "icon", href: assetUrl("favicon.png"), type: "image/png" },
    ],
  }),
  component: RootComponent,
});

function RootComponent() {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <HeadContent />
      </head>
      <body className="flex min-h-screen flex-col">
        <RootProvider
          theme={{ defaultTheme: "light" }}
          search={{ options: { type: "static" } }}
        >
          <ProviderProvider>
            <Outlet />
          </ProviderProvider>
        </RootProvider>
        <Scripts />
      </body>
    </html>
  );
}
