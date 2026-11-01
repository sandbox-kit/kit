import { mkdir, readdir, writeFile } from "node:fs/promises";
import { dirname, join } from "node:path";
import server from "../node_modules/.nitro/vite/services/ssr/index.js";

const root = new URL("../.vercel/output/static/", import.meta.url);
const base = "/kit";

const docFiles = await readdir(new URL("../content/docs/", import.meta.url));
const paths = [
  "/api/search",
  "/llms.txt",
  "/llms-full.txt",
  ...docFiles
    .filter((name) => name.endsWith(".mdx"))
    .flatMap((name) => {
      const slug = name.replace(/\.mdx$/, "");
      const page = slug === "index" ? "" : `/${slug}`;
      return [`/docs${page === "" ? "/index.md" : `${page}.md`}`, `/og/docs${page}/image.webp`];
    }),
];

for (const path of paths) {
  const response = await server.fetch(new Request(`http://localhost${base}${path}`));
  if (!response.ok) {
    throw new Error(`${path} returned ${response.status}`);
  }
  const file = join(root.pathname, path.slice(1));
  await mkdir(dirname(file), { recursive: true });
  await writeFile(file, Buffer.from(await response.arrayBuffer()));
}

await writeFile(join(root.pathname, ".nojekyll"), "");
await writeFile(
  join(root.pathname, "404.html"),
  `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <title>Page not found | Sandbox Kit</title>
    <meta name="viewport" content="width=device-width, initial-scale=1" />
  </head>
  <body>
    <p>That page is not in the guide.</p>
    <p><a href="${base}/">Sandbox Kit</a></p>
  </body>
</html>
`,
);
