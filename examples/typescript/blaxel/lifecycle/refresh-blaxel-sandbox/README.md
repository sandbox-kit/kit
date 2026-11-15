# Refresh a Blaxel sandbox

Refresh metadata for an existing Blaxel sandbox by name.

## What this means

`sandbox.refresh` updates `info()` on this handle (state, labels, and related metadata).

## Edge cases

- If this handle does not support refresh, you get `Unsupported`.

## Run

```sh
cd examples/typescript/blaxel/lifecycle/refresh-blaxel-sandbox
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
# Sandbox name (Blaxel id is the name).
BL_SANDBOX_NAME=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/blaxel/lifecycle/refresh-blaxel-sandbox/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
