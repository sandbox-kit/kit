# Start a Daytona sandbox

`Get` a sandbox, check **Start**, then start it from a retained stopped (or
otherwise startable) state.

## What this means

Daytona sandboxes that were stopped (and not deleted) keep their filesystem and
can be started again. See [persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Ephemeral sandboxes that already stopped are gone—`Get` fails.
- GPU start may still be denied by account/features even when the capability appears.
- Mutates the sandbox named in `.env`.

## Run

```sh
cd examples/typescript/daytona/start-daytona-sandbox
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/daytona/start-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
