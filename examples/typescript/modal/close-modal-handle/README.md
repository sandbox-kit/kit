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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/modal/close-modal-handle/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
