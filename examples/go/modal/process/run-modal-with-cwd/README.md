# Run a Modal command with a working directory

Run `pwd` with an absolute `working_directory`. The program loads that
sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Modal collects both
streams.

## Edge cases

- A relative working directory is rejected before Modal is called. The error
  field is `working_directory`.
- Omit `working_directory` to use the Modal sandbox default.

## Run

```sh
cd examples/go/modal/process/run-modal-with-cwd
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

# Absolute working directory for this command.
MODAL_WORKING_DIRECTORY=/tmp

```

## Warnings

- Needs an existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

## Related

- [typescript twin](../../../../typescript/modal/process/run-modal-with-cwd/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
