# Create a Blaxel sandbox with labels and env

Create one Blaxel sandbox with labels and environment variables.

## What this means

Set `labels` and create-time environment variables. They appear as sandbox metadata and runtime env.

## Edge cases

- Invalid environment names are rejected before Blaxel is called.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-with-labels
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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-with-labels/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
