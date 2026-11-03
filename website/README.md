# Sandbox Kit docs

User guide for the Sandbox Kit client. The site is a Fumadocs app on TanStack Start. Pages live in `content/docs`.

```sh
npm install
npm run dev
```

Open http://localhost:3000.

```sh
npm run build:pages
```

That build is the GitHub Pages site. GitHub Actions publishes it from `main` to https://sandbox-kit.github.io/kit/. In the repository settings, set Pages to deploy from GitHub Actions.

TypeScript and Go share each page. The language tab is stored in the browser. Each page covers one part of the client: a provider setup, create, lifecycle, filesystem, process, configuration, or errors.

The behavior reference for the same API stays in the repository `docs/` directory. Runnable projects stay in `examples/`.
