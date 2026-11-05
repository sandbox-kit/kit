# Docs site

Fumadocs + TanStack Start site for the Sandbox Kit SDK user guide. Pages live in `content/docs`.

```sh
npm install
npm run dev
```

Open http://localhost:3000.

```sh
npm run build:pages
```

GitHub Actions deploys `main` to https://sandbox-kit.github.io/kit/. Set Pages to deploy from GitHub Actions.

Language and provider tabs persist in the browser. Runnable samples are linked from the Examples page into the repo `examples/` tree.
