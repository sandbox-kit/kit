# Run a Daytona command with a working directory

Run `pwd` with an absolute `working_directory`. The program loads that
sandbox with `client.get`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Daytona's basic exec
returns stdout only; `stderr` is empty.

## Edge cases

- A relative working directory is rejected before Daytona is called. The error
  field is `working_directory`.
- Omit `working_directory` to use the Daytona sandbox default.

## Run

```sh
cd examples/typescript/daytona/process/run-daytona-with-cwd
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

# Absolute working directory for this command.
DAYTONA_WORKING_DIRECTORY=/tmp

```

## Warnings

- Needs an existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

## Related

- [Go twin](../../../../go/daytona/process/run-daytona-with-cwd/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
