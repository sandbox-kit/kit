# Create a Blaxel sandbox with idle delete

Create one Blaxel sandbox that deletes after an idle period.

## What this means

Maps `lifetime.idleTerminate` (`PolicyModeAfter`) to a Blaxel `ttl-idle` expiration policy with action delete.

## Edge cases

- Disabling idle terminate is not representable and is rejected.
- Minimum idle duration follows Blaxel rules (1 minute).

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-idle
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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-idle/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
