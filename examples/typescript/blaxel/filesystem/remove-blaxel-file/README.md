# Remove a Blaxel path

Remove one absolute path in a Blaxel sandbox.

## What this means

Set `BL_RECURSIVE=true` to remove recursively when supported by options.

## Edge cases

- Removing a missing path returns not found.

## Run

```sh
cd examples/typescript/blaxel/filesystem/remove-blaxel-file
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
# Absolute path inside the sandbox.
BL_REMOTE_PATH=/tmp/kit-example.txt
# Set true for recursive remove.
BL_RECURSIVE=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/blaxel/filesystem/remove-blaxel-file/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
