# Stat a Blaxel path

Stat one absolute path in a Blaxel sandbox.

## What this means

Returns type and size for the path.

## Edge cases

- Missing paths return not found.

## Run

```sh
cd examples/typescript/blaxel/filesystem/stat-blaxel-file
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
BL_SANDBOX_ID=
# Optional region override (default us-pdx-1).
BL_REGION=
# Absolute path inside the sandbox.
BL_REMOTE_PATH=/tmp/kit-example.txt
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/blaxel/filesystem/stat-blaxel-file/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
