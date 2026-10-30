# Create a Modal sandbox with a named secret

Create from an image and attach an existing Modal secret by name.

## What this means

Modal resolves the secret by name in the configured app and environment, then
injects it as environment variables.

## Edge cases

- `MODAL_SECRET_NAME` must already exist in that app.
- Renaming individual variables in the request is unsupported. Define that in the Modal secret.

## Run

```sh
cd examples/go/modal/create-modal-with-secrets
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
MODAL_SECRET_NAME=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to Modal.

## Related

- [typescript twin](../../../typescript/modal/create-modal-with-secrets/README.md)

- Catalog: [examples/go/README.md](../../README.md)
