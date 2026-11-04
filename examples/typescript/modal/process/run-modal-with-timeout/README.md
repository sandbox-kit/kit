# Cap a Modal Process.Run timeout

Run `echo hello with timeout` with `sandbox.process.run` / `Process.Run` and a
whole-second `timeout` of five seconds (`5000` ms). The program loads that
sandbox with `client.attach`.

## What this means

A positive `timeout` caps this call. Values must be whole seconds
(`5000` ms / `5 * time.Second`). Omit or zero adds no Modal command limit
(TypeScript omits `timeoutMs`; Go sends `Timeout: 0`).

## Edge cases

- `1500` ms is rejected before Modal runs. The error field is `timeout`.
- Abort the call (or cancel the Go context) to end the wait. The remote
  command may keep running; the handle stays busy until that provider call
  settles.

## Run

```sh
cd examples/typescript/modal/process/run-modal-with-timeout
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

- [Go twin](../../../../go/modal/process/run-modal-with-timeout/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
