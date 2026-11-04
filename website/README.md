# Sandbox Kit docs site

User guide for the Sandbox Kit SDK—a lightweight, provider-pluggable client for cloud sandboxes. Fumadocs + TanStack Start. Source pages: `content/docs`.

```sh
npm install
npm run dev
```

Open http://localhost:3000.

```sh
npm run build:pages
```

GitHub Actions publishes from `main` to https://sandbox-kit.github.io/kit/. Set Pages to deploy from GitHub Actions.

## Structure

| Section    | Pages                                                                     |
| ---------- | ------------------------------------------------------------------------- |
| Quickstart | `index`                                                                   |
| Concepts   | `concepts`                                                                |
| Providers  | `daytona`, `modal`                                                        |
| Guides     | `configuration`, `create`, `lifecycle`, `filesystem`, `process`, `errors` |
| Examples   | `examples`                                                                |

Language tabs (TypeScript / Go) persist in the browser. Behavior reference for the same API lives in the repository `docs/` directory. Runnable projects live in `examples/`.
