# Run a Modal shell command

Run one shell string with `Process.Run` / `sandbox.process.run` using
`shell`. The program loads that sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Modal collects both
streams.

## Edge cases

- Set exactly one of `argv` or `shell`. Setting both, or neither, fails before
  Modal is called. The error field is `argv`.
- A blank `shell` fails with field `shell`. Kit maps `shell` to `bash -c`.
- Provider command timeouts must be whole seconds.

## Run

```sh
cd examples/typescript/modal/process/run-modal-shell-command
cp .env.example .env
# Fill in the values below, then:
npm install && npm run build && npm start
```

Requires Node.js 20/22/24.15+ and npm 11.19.0. Skip copying if `.env` already exists.

### Environment

```dotenv
MODAL_TOKEN_ID=
MODAL_TOKEN_SECRET=
MODAL_APP_NAME=
MODAL_ENVIRONMENT=main
MODAL_SANDBOX_ID=

```

## Warnings

- Needs an existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

## Related

- [Go twin](../../../../go/modal/process/run-modal-shell-command/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
