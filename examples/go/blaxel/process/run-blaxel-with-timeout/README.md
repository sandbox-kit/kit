# Run a Blaxel command with timeout

Cap a Blaxel Process.Run with a whole-second timeout.

## What this means

Command timeouts must be whole seconds. Omit or zero for no Blaxel process
auto-kill; the blocking wait can still end at Blaxel's edge hold.

## Edge cases

- `1500` ms is rejected; `2000` is two seconds.
- Canceling the local wait does not guarantee the remote command stops.

## Run

```sh
cd examples/go/blaxel/process/run-blaxel-with-timeout
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.0+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

### Environment

```dotenv
BL_API_KEY=
BL_WORKSPACE=
# Sandbox name (Blaxel id is the name).
BL_SANDBOX_ID=
# Optional region override (default us-pdx-1).
BL_REGION=
```

## Warnings

- Acts on **real** Blaxel resources when credentials succeed.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../../typescript/blaxel/process/run-blaxel-with-timeout/README.md)
- Catalog: [examples/go/README.md](../../../README.md)
