# Create a Blaxel sandbox with memory

Create one Blaxel sandbox with an explicit memory size in MiB. Blaxel derives CPU from memory.

## What this means

You can set `resources.memoryMib` only. CPU cores, limits, and disk overrides are rejected. Omit memory to use Blaxel's default.

## Edge cases

- Memory values that Blaxel rejects are returned as provider errors.

## Run

```sh
cd examples/typescript/blaxel/configuration/create-blaxel-with-resources
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

- [go twin](../../../../go/blaxel/configuration/create-blaxel-with-resources/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
