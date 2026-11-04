import react from "@vitejs/plugin-react";
import { tanstackStart } from "@tanstack/react-start/plugin/vite";
import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import { fumadocsMdx } from "fumadocs-mdx/vite";
import { nitro } from "nitro/vite";

// GitHub project Pages serves this repo at /kit/. Local dev leaves DOCS_BASE unset.
const base = process.env.DOCS_BASE ?? "/";
const basePath = base === "/" ? "" : base.replace(/\/$/, "");

export default defineConfig({
  base,
  // Prefer a project-local cache so root-owned node_modules/.vite does not block local dev.
  cacheDir: process.env.VITE_CACHE_DIR ?? ".vite",
  server: {
    port: 3000,
  },
  plugins: [
    fumadocsMdx(),
    tailwindcss(),
    tanstackStart({
      router: basePath ? { basepath: basePath } : undefined,
      prerender: {
        enabled: true,
        crawlLinks: true,
      },
    }),
    react(),
    nitro({
      preset: process.env.NITRO_PRESET ?? "vercel",
    }),
  ],
  resolve: {
    tsconfigPaths: true,
    alias: {
      tslib: "tslib/tslib.es6.js",
    },
  },
});
