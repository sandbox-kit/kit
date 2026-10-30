# Attach to a Modal sandbox by ID

Build a handle from `MODAL_SANDBOX_ID` with `Client.Attach` / `client.attach`.

## What this means

Modal's SDKs expose ID-based sandbox handles (for example `Sandbox.from_id` /
equivalent). Kit's `Attach` matches that: construct a handle from an ID without
claiming a full remote metadata refresh. You can then terminate or close the
handle.

## Edge cases

- A wrong ID may only fail when you invoke an operation, not at attach time.
- Attach is not proof the sandbox is still running.
- Does not terminate the sandbox by itself.

## Run

```sh
cd examples/typescript/modal/attach-modal-sandbox
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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/modal/attach-modal-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
