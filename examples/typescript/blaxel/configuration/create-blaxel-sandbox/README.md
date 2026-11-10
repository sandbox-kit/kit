# Create a Blaxel sandbox

Create one Blaxel sandbox from a registry image with an explicit name. Prints the sandbox ID and leaves the sandbox running.

## What this means

Blaxel create requires a nonblank `name`, an image source, and client `region`
(examples read `BL_REGION`, default `us-pdx-1`). Kit calls native create (not
create-if-not-exists). Closing the Kit client only drops the local SDK
connection.

## Edge cases

- Omitting `name` or `source.image` fails before Blaxel is called.
- A create error is not proof that no sandbox was provisioned.
- Reclaim sandboxes with TTL / idle policies or `delete`.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-sandbox
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Skip copying if
`.env` already exists.

### Environment

```dotenv
BL_API_KEY=
BL_WORKSPACE=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/blaxel/configuration/create-blaxel-sandbox/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
