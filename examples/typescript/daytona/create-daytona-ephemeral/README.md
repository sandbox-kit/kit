# Create an ephemeral Daytona sandbox

Create from the default snapshot with **`lifetime.ephemeral = true`**.

## What this means

Daytona sandboxes persist by default. An **ephemeral** sandbox is deleted as soon
as it stops, discarding state and avoiding stopped-state disk billing. Daytona
sets auto-delete to immediate (`0`) for ephemeral creates. See
[persistence](https://www.daytona.io/docs/en/persistence/).

## Edge cases

- Stopping an ephemeral sandbox deletes it; do not expect to start it again.
- Ephemeral conflicts with positive idle-pause, archive, or delayed-delete
  policies—Kit rejects those combinations locally.
- [Linked children](../create-daytona-linked/README.md) must be ephemeral.

## Run

```sh
cd examples/typescript/daytona/create-daytona-ephemeral
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/daytona/create-daytona-ephemeral/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
