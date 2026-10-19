# Create a Daytona sandbox

Create one Daytona sandbox from the account's **default snapshot** by calling
`Create` / `create` with no options. Prints the sandbox ID and leaves the sandbox
running.

## What this means

Daytona sandboxes are persistent by default: create boots from a snapshot (a
prebuilt filesystem + toolbox configuration). Omitting `source` means “use the
default snapshot for this target,” not “no image.”

The snapshot chooses the sandbox class (container, Linux VM, Windows, or GPU).
That class decides which lifecycle operations later succeed (`Stop`, `Pause`,
`Archive`, and so on). See Daytona's
[sandbox guide](https://www.daytona.io/docs/en/sandboxes/) and
[snapshots](https://www.daytona.io/docs/en/snapshots/).

Closing the Kit client only drops the local SDK connection. Daytona keeps the
sandbox until you stop/delete it or a lifetime policy fires.

## Edge cases

- Default snapshot resources are inherited. CPU/memory/disk overrides require an
  **image** source ([from-image](../create-daytona-from-image/README.md)).
- Non-default `isolation.kind` values are rejected; class comes from the snapshot.
- Optional `DAYTONA_API_URL` / `DAYTONA_TARGET` select endpoint and region. When
  the URL is omitted, Kit tries `DAYTONA_API_URL`, then `DAYTONA_SERVER_URL`, then
  `https://app.daytona.io/api`. Remote endpoints need HTTPS.
- A create error is not proof that no sandbox was provisioned.

## Run

```sh
cd examples/typescript/daytona/create-daytona-sandbox
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

## Authentication note

This example sets an explicit API key. Omitting `auth` lets the Daytona SDK
resolve its supported environment credentials. Bearer auth maps through
`auth.bearerToken` and requires `scope.organizationId`.


## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/daytona/create-daytona-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
