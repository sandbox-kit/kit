# Run a Modal command

Run `echo hello from kit` with `Process.Run` / `sandbox.process.run` using
`argv`. The program opens that sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Modal returns both
stdout and stderr.

## Edge cases

- Set exactly one of `argv` or `shell`. Setting both, or neither, fails before
  Modal is called. The error field is `argv`.
- A relative working directory is rejected before Modal is called. The error
  field is `working_directory`.
- Omit or zero `Timeout` for no Modal command limit. Command timeouts must be
  whole seconds. The error field is `timeout`. Cancel the context to stop waiting.

## Run

```sh
cd examples/go/modal/process/run-modal-command
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

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

- [typescript twin](../../../../typescript/modal/process/run-modal-command/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
