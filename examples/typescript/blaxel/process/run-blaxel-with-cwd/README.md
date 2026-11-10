# Run a Blaxel command with cwd

Run a command with an explicit working directory.

## What this means

Sets `workingDirectory` / `WorkingDirectory` for the process.

## Edge cases

- Paths in the sandbox are absolute.

## Run

```sh
cd examples/typescript/blaxel/process/run-blaxel-with-cwd
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
# Working directory inside the sandbox.
BL_WORKING_DIRECTORY=/tmp
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/blaxel/process/run-blaxel-with-cwd/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
