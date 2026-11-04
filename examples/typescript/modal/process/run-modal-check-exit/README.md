# Check a non-zero Modal exit code

Run `exit 7` via `shell` and read `exit_code` from the result. The program
loads that sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Modal returns both stdout and stderr.

## Edge cases

- A non-zero exit code is not an SDK error. Check `exit_code` / `exitCode` in
  your program.
- Transport, validation, and Modal API failures still return errors.

## Run

```sh
cd examples/typescript/modal/process/run-modal-check-exit
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

- [Go twin](../../../../go/modal/process/run-modal-check-exit/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
