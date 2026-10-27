# Get a Daytona sandbox

Load a provider-verified handle for `DAYTONA_SANDBOX_ID` with
`Client.Get` / `client.get` and print identity, metadata, and capabilities.

## What this means

`Get` performs a remote lookup and returns current Daytona metadata (state,
class, labels, …). Use it when you already have an ID and need a fresh handle.
Capabilities list the lifecycle methods Daytona implements.

## Edge cases

- Unknown or deleted IDs fail at the provider (`NotFound` or similar).
- The resolved class can differ from what you assumed at create time.
- Does not change the cloud sandbox.

## Run

```sh
cd examples/typescript/daytona/get-daytona-sandbox
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/daytona/get-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
