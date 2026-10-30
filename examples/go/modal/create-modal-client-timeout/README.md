# Create a Modal sandbox using the client timeout

Create with `WaitConditionReady` and a command probe. The example sets the
client `Timeout` to 90 seconds and leaves `Provisioning.Timeout` unset.

## What this means

Kit uses that client timeout as the create deadline, including the ready
wait. The entrypoint is `sleep 300`, so the sandbox stays up after the probe
passes.

## Edge cases

- An explicit `Provisioning.Timeout` replaces the inherited deadline when you set one.
- `WaitConditionReady` requires a readiness probe and a positive timeout.

## Run

```sh
cd examples/go/modal/create-modal-client-timeout
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

- [typescript twin](../../../typescript/modal/create-modal-client-timeout/README.md)

- Catalog: [examples/go/README.md](../../README.md)
