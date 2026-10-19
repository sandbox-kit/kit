# Start a Daytona sandbox with an operation timeout

`Get` a sandbox, check **Start**, then start with **`timeout: 60_000`** (ms).

## What this means

Lifecycle methods take one options bag. `timeout` is both Kit’s abort deadline
and the native start timeout passed to Daytona (converted to seconds).

## Edge cases

- Omitting `timeout` falls back to the client’s default timeout.
- Ephemeral sandboxes that already stopped cannot be started again.
- Mutates the sandbox named in `.env`.

## Run

```sh
cd examples/typescript/daytona/start-daytona-with-timeout
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
- **Changes or destroys** that sandbox.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
