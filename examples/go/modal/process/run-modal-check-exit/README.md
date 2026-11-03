# Check a non-zero Modal exit code

Run `exit 7` via `shell` and read `exit_code` from the result. The program
loads that sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Modal collects both
streams.

## Edge cases

- A non-zero exit code is not an SDK error. Check `exit_code` / `exitCode` in
  your program.
- Transport, validation, and Modal API failures still return errors.

## Run

```sh
cd examples/go/modal/process/run-modal-check-exit
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

- [typescript twin](../../../../typescript/modal/process/run-modal-check-exit/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
