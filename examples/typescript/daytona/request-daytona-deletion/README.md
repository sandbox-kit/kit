# Request Daytona sandbox deletion without waiting

`Get` a sandbox, check **Delete**, then delete with **`wait: false`**.

## What this means

Waited deletion ([delete-daytona-sandbox](../delete-daytona-sandbox/README.md))
blocks until Daytona finishes. Non-waited deletion returns after the delete is
accepted. The handle stores the metadata Daytona returns with that response.

## Edge cases

- Follow up with `get`/`refresh` if you need to confirm the sandbox is gone.
- Mutates (and eventually destroys) the sandbox named in `.env`.

## Run

```sh
cd examples/typescript/daytona/request-daytona-deletion
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
