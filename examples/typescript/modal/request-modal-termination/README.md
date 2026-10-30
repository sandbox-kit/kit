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
cd examples/typescript/modal/request-modal-termination
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
- **Changes or destroys** that sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/modal/request-modal-termination/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
