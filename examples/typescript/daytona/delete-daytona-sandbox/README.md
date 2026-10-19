# Delete a Daytona sandbox

`Get` a sandbox, check **Delete**, then delete with **`wait: true`** until Daytona
finishes removal.

## What this means

Delete irreversibly removes the sandbox resource. After success, Kit marks the
handle terminal and rejects further lifecycle calls.

## Edge cases

- **Destructive.** The ID in `.env` will be gone.
- `wait: false` returns earlier without proving completion; this example waits.
- Deleting a parent also deletes linked children (cascade).

## Run

```sh
cd examples/typescript/daytona/delete-daytona-sandbox
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

- [go twin](../../../go/daytona/delete-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
