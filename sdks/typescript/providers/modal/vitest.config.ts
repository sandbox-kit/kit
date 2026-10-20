import path from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig } from "vitest/config";

const root = path.dirname(fileURLToPath(import.meta.url));

export default defineConfig({
  resolve: {
    alias: [
      {
        find: /^#(mappings|checks|config)\/(.+)\.js$/,
        replacement: path.join(root, "src/$1/$2.ts"),
      },
    ],
  },
  test: {
    environment: "node",
  },
});
