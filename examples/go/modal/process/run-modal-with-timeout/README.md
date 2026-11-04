# Cap a Modal Process.Run timeout

Run `echo hello with timeout` with `Process.Run` / `sandbox.process.run` and a
whole-second `timeout` of five seconds. The program loads that sandbox with
`client.attach`.

## What this means

A positive `timeout` caps this call. Values must be whole seconds
(`5 * time.Second` / `5000` ms). Omit or zero adds no Modal command limit (Go
sends `Timeout: 0`; TypeScript omits `timeoutMs`).

## Edge cases

- `1500` ms / `1500 * time.Millisecond` is rejected before Modal runs. The
  error field is `timeout`.
- Cancel the Go context (or abort the TypeScript call) to end the wait. The
  remote command may keep running; the handle stays busy until that provider
  call settles.

## Run

```sh
cd examples/go/modal/process/run-modal-with-timeout
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

- [typescript twin](../../../../typescript/modal/process/run-modal-with-timeout/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
