# Run a Blaxel command

Run one argv command in a Blaxel sandbox and print exit code plus streams.

## What this means

`Process.Run` is blocking. Set exactly one of `argv` or `shell`. Argv is sent as a POSIX-quoted shell string.

## Edge cases

- Non-zero exit codes are success for the SDK—read `exitCode`.
- Blaxel returns stdout and stderr.

## Run

```sh
cd examples/typescript/blaxel/process/run-blaxel-command
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

- [go twin](../../../../go/blaxel/process/run-blaxel-command/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
