# Terminate a Modal sandbox with a timeout

`Attach` a sandbox, check **Terminate**, then terminate with `Wait` set and a
60-second deadline.

## What this means

The example sets a 60-second context deadline and `TerminateOptions.Timeout`
of one minute. That is how long Kit waits. Modal's terminate call waits until
the sandbox exits or that deadline fires. It does not take a separate duration.

## Edge cases

- `Wait` false returns immediately. See [request-modal-termination](../request-modal-termination/README.md).
- The call destroys the sandbox named in `.env`.

## Run

```sh
cd examples/go/modal/terminate-modal-with-timeout
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

- Needs an **existing sandbox ID** in `.env`.
- **Destroys** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Modal.

## Related

- [typescript twin](../../../typescript/modal/terminate-modal-with-timeout/README.md)

- Catalog: [examples/go/README.md](../../README.md)
