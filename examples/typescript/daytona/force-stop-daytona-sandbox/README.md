# Force-stop a Daytona sandbox

`Get` a sandbox, check **Stop**, then stop it with **`force: true`**.

## What this means

Normal stop asks the sandbox to shut down cleanly. Force stop tells Daytona to
stop even when a graceful stop may not finish. Kit folds `force` into the same
options bag as the call timeout.

## Edge cases

- Prefer plain [stop](../stop-daytona-sandbox/README.md) when a clean shutdown is enough.
- Force stop still requires the Stop capability and a live sandbox ID.
- Mutates the sandbox named in `.env`.

## Run

```sh
cd examples/typescript/daytona/force-stop-daytona-sandbox
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
