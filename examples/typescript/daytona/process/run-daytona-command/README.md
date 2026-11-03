# Run a Daytona command

Run `echo hello from kit` with `sandbox.process.run` using `argv`. The program
loads that sandbox with `client.get`.

## What this means

The call waits for the command to finish and returns `exitCode`, `stdout`, and
`stderr`. A non-zero exit code is not an SDK error. Daytona's basic exec
returns stdout only; `stderr` is empty.

## Edge cases

- Set exactly one of `argv` or `shell`. Setting both, or neither, fails before
  Daytona is called. The error field is `argv`.
- A relative working directory is rejected before Daytona is called. The error
  field is `working_directory`.
- An omitted command timeout passes `0` to Daytona so it adds no command
  limit. Provider command timeouts must be whole seconds (`2000` is valid;
  `1500` is rejected). The error field is `timeout`.

## Run

```sh
cd examples/typescript/daytona/process/run-daytona-command
cp .env.example .env
# Fill in the values below, then:
npm install && npm run build && npm start
```

Requires Node.js 20/22/24.15+ and npm 11.19.0. Skip copying if `.env` already
exists.

### Environment

```dotenv
DAYTONA_API_KEY=
DAYTONA_SANDBOX_ID=

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

## Related

- [go twin](../../../../go/daytona/process/run-daytona-command/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
