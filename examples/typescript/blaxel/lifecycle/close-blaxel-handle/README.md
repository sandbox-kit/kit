# Close a Blaxel sandbox handle

Close one Blaxel sandbox handle without contacting Blaxel or changing the sandbox.

## What this means

`Close` releases local handle resources only. The sandbox keeps running until TTL, idle policy, or delete.

## Edge cases

- After close, further handle calls fail locally.
- Client close also does not delete sandboxes.

## Run

```sh
cd examples/typescript/blaxel/lifecycle/close-blaxel-handle
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

- [go twin](../../../../go/blaxel/lifecycle/close-blaxel-handle/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
