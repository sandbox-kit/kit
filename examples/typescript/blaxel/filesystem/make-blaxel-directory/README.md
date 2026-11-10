# Make a Blaxel directory

Create a directory at an absolute path in a Blaxel sandbox.

## What this means

`parents` defaults to true.

## Edge cases

- Absolute paths only.

## Run

```sh
cd examples/typescript/blaxel/filesystem/make-blaxel-directory
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
BL_REMOTE_PATH=/tmp/kit-example-dir
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/blaxel/filesystem/make-blaxel-directory/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
