# Create a Blaxel sandbox with a client timeout

Create one Blaxel sandbox using a shorter client-wide timeout.

## What this means

Sets `Config.Timeout` / `timeout` on the shared client. This bounds local waiting for SDK calls, including create.

## Edge cases

- Client `timeout` is how long create may take.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-client-timeout
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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-client-timeout/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
