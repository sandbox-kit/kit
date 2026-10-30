# Create a Daytona sandbox with labels and environment

Create from an image with request `labels` and `environment` maps.

## What this means

Labels become Daytona sandbox labels; environment variables are passed into the sandbox process environment.

## Edge cases

- Environment keys must be valid env names (`[A-Za-z_][A-Za-z0-9_]*`).
- Label values are opaque strings to Kit; Daytona enforces its own size limits.

## Run

```sh
cd examples/typescript/daytona/create-daytona-with-labels
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Uses `await using`
so local client cleanup runs on success or failure.

### Environment

```dotenv
DAYTONA_API_KEY=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
