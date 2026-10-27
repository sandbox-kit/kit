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
cd examples/typescript/modal/close-modal-handle
cp .env.example .env
# Fill in the values below, then:
npm install
npm run build
npm start
```

Requires Node.js 20/22/24.15+, npm 11.19.0, and this checkout. Uses `await using`
so local client cleanup runs on success or failure.

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

- [go twin](../../../go/modal/close-modal-handle/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
