# Create linked Daytona sandboxes

Create a parent sandbox, then an ephemeral child with
`network.linkedSandbox = parent.id`. Prints both IDs.

## What this means

Linked sandboxes attach ephemeral children to a parent. Daytona schedules each
child on the same runner and joins them on a shared link network so the group can
reach each other by sandbox ID or name (DNS aliases). Children are always
ephemeral (auto-delete `0`). Deleting the parent cascades to linked children.
See [sandboxes → linked sandboxes](https://www.daytona.io/docs/en/sandboxes/).

This example creates **two** cloud sandboxes.

## Edge cases

- Kit requires `lifetime.ephemeral = true` on the child or rejects before the API.
- If parent create succeeds and child create fails, the parent still exists.
- Clean up both IDs when you are done (or delete the parent to cascade).

## Run

```sh
cd examples/typescript/daytona/create-daytona-linked
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

- [go twin](../../../go/daytona/create-daytona-linked/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
