# Run a Daytona shell command

Run one shell string with `Process.Run` / `sandbox.process.run` using
`shell`. The program loads that sandbox with `client.get`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Daytona's basic exec
returns stdout only; `stderr` is empty.

## Edge cases

- Set exactly one of `argv` or `shell`. Setting both, or neither, fails before
  Daytona is called. The error field is `argv`.
- A blank `shell` fails with field `shell`.
- An omitted or zero command timeout passes `0` to Daytona so it adds no
  command limit. Provider command timeouts must be whole seconds.

## Run

```sh
cd examples/typescript/daytona/process/run-daytona-shell-command
cp .env.example .env
# Fill in the values below, then:
npm install && npm run build && npm start
```

Requires Node.js 20/22/24.15+ and npm 11.19.0. Skip copying if `.env` already exists.

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

- [Go twin](../../../../go/daytona/process/run-daytona-shell-command/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
