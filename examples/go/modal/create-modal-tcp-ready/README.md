# Create a Modal sandbox that waits for a TCP port

Create with a TCP readiness probe on port 8080, `WaitConditionReady`, and a
two-minute provisioning timeout.

## What this means

The call returns after something accepts TCP connections on that port. The
entrypoint runs `nc -l -p 8080` so the probe can succeed.

## Edge cases

- `WaitConditionReady` requires a probe and a positive provisioning timeout.
- The Alpine image needs busybox `nc`.

## Run

```sh
cd examples/go/modal/create-modal-tcp-ready
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
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Modal.

## Related

- [typescript twin](../../../typescript/modal/create-modal-tcp-ready/README.md)

- Catalog: [examples/go/README.md](../../README.md)
