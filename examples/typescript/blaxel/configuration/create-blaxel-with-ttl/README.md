# Create a Blaxel sandbox with maximum lifetime

Create one Blaxel sandbox that deletes after a maximum lifetime TTL.

## What this means

Set `lifetime.maximumLifetime` (whole seconds, at least 5 minutes). The sandbox deletes when that max age elapses.

## Edge cases

- Fractional-second durations are rejected.
- Idle cleanup uses `idleTerminate`, not maximum lifetime.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-with-ttl
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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-with-ttl/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
