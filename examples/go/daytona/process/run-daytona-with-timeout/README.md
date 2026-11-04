# Cap a Daytona Process.Run timeout

Run `echo hello with timeout` with `Process.Run` / `sandbox.process.run` and a
whole-second `timeout` of five seconds. The program loads that sandbox with
`client.get`.

## What this means

A positive `timeout` caps this call. Values must be whole seconds
(`5 * time.Second` / `5000` ms). Omit or zero adds no extra client deadline and
no separate Daytona command limit (Daytona would otherwise default to 10
seconds).

## Edge cases

- `1500` ms / `1500 * time.Millisecond` is rejected before Daytona runs. The
  error field is `timeout`.
- Cancel the Go context (or abort the TypeScript call) to end the wait. The
  remote command may keep running; the handle stays busy until that provider
  call settles.

## Run

```sh
cd examples/go/daytona/process/run-daytona-with-timeout
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
DAYTONA_API_KEY=
DAYTONA_SANDBOX_ID=

# Optional Daytona endpoint and target overrides.
# DAYTONA_API_URL=https://app.daytona.io/api
# DAYTONA_TARGET=

```

## Warnings

- Needs an existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes
  your connection to the provider.

## Related

- [typescript twin](../../../../typescript/daytona/process/run-daytona-with-timeout/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
