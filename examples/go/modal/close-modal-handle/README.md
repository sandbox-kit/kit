# Close a Modal sandbox handle

Attach by ID, then call **`Sandbox.Close` / `sandbox.close()`**.

## What this means

Closing the handle does not terminate the sandbox. The sandbox keeps running
until idle timeout, max lifetime, entrypoint exit, or an explicit terminate.
Closing the client closes your connection to Modal. You can keep using the
handle.

## Edge cases

- Safe against sandboxes you intend to keep.
- Closing the client afterward also does not terminate.
- A later call on the handle goes to Modal.

## Run

```sh
cd examples/go/modal/close-modal-handle
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
- Handle close does not terminate the cloud sandbox. Client close releases the SDK connection.

## Related

- [typescript twin](../../../typescript/modal/close-modal-handle/README.md)

- Catalog: [examples/go/README.md](../../README.md)
