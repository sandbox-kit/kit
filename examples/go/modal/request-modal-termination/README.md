# Request Modal termination without waiting

Attach by ID and terminate with **`wait: false`**.

## What this means

Same terminate operation as the waited example, but returns once the request is
accepted instead of blocking for an exit code.

## Edge cases

- Still ends the sandbox; do not assume it is gone the instant the call returns.
- No reliable exit code in this mode.
- Destructive to the ID in `.env`.

## Run

```sh
cd examples/go/modal/request-modal-termination
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
- **Changes or destroys** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [typescript twin](../../../typescript/modal/request-modal-termination/README.md)

- Catalog: [examples/go/README.md](../../README.md)
