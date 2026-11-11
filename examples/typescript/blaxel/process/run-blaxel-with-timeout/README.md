# Run a Blaxel command with timeout

Cap a Blaxel `sandbox.process.run` with a whole-second timeout.

## What this means

Command timeouts must be whole seconds. This run uses 5 seconds and caps how
long the command may run. Omit or zero waits until the command ends (or you
cancel), with no Blaxel process time limit.

## Edge cases

- `1500` ms is rejected; `2000` is two seconds.
- Canceling the local wait does not guarantee the remote command stops.

## Run

```sh
cd examples/typescript/blaxel/process/run-blaxel-with-timeout
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

- [go twin](../../../../go/blaxel/process/run-blaxel-with-timeout/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
