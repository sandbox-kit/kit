# Run a Modal command

Run `echo hello from kit` with `sandbox.process.run` using `argv`. The program
opens that sandbox with `client.attach`.

## What this means

The call waits for the command to finish and returns `exitCode`, `stdout`, and
`stderr`. A non-zero exit code is not an SDK error. Modal pipes both streams
and Kit collects them before returning.

## Edge cases

- Set exactly one of `argv` or `shell`. Setting both, or neither, fails before
  Modal is called. The error field is `argv`.
- A relative working directory is rejected before Modal is called. The error
  field is `working_directory`.
- An omitted or zero command timeout adds no Modal command limit: TypeScript
  omits `timeoutMs` and must not send `timeoutMs: 0`. Provider command
  timeouts must be whole seconds (`2000` is valid; `1500` is rejected). The
  error field is `timeout`.

## Run

```sh
cd examples/typescript/modal/process/run-modal-command
cp .env.example .env
# Fill in the values below, then:
npm install && npm run build && npm start
```

Requires Node.js 20/22/24.15+ and npm 11.19.0. Skip copying if `.env` already
exists.

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

- [go twin](../../../../go/modal/process/run-modal-command/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
