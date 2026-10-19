# Close a Modal sandbox handle

Attach by ID, then close the **local handle** without terminating the sandbox.

## What this means

Handle close only drops local SDK state. The Modal sandbox keeps running until
idle timeout, max lifetime, entrypoint exit, or an explicit terminate.

## Edge cases

- Safe against sandboxes you intend to keep.
- Closing the client afterward also does not terminate.
- Do not use the handle after close.

## Run

```sh
cd examples/go/modal/close-modal-handle
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [typescript twin](../../../typescript/modal/close-modal-handle/README.md)

- Catalog: [examples/go/README.md](../../README.md)
