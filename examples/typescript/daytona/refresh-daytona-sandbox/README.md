# Refresh a Daytona sandbox

`Get` a sandbox, confirm **Refresh** is supported, then refresh cached metadata.

## What this means

`Refresh` re-fetches Daytona metadata onto the handle. It does not start, stop,
or delete the sandbox. Useful after another process may have changed state.

## Edge cases

- Refresh replaces the metadata on the handle. Capabilities stay the Daytona method list.
- Not a substitute for lifecycle operations.
- Requires a live sandbox ID.

## Run

```sh
cd examples/typescript/daytona/refresh-daytona-sandbox
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
DAYTONA_SANDBOX_ID=
# Optional: DAYTONA_API_URL=, DAYTONA_TARGET=
```

## Warnings

- Needs an **existing sandbox ID** in `.env`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/refresh-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
