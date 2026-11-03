# Run a Modal command with environment variables

Run `printenv KIT_DEMO` with environment variables for that call only. The
program loads that sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exit_code`, `stdout`,
and `stderr`. A non-zero exit code is not an SDK error. Modal collects both
streams.

## Edge cases

- An omitted `environment` map leaves the Modal default. An empty map is
  explicit.
- These variables apply to this command, not to later calls on the same handle.

## Run

```sh
cd examples/typescript/modal/process/run-modal-with-env
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

- [Go twin](../../../../go/modal/process/run-modal-with-env/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
