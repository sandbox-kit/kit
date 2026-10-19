# Terminate a Modal sandbox (wait)

Attach by ID, check **Terminate**, then terminate with **`wait: true`** and print
the exit code.

## What this means

`terminate` ends the sandbox execution. Modal also ends sandboxes when max
timeout or idle timeout fires, or when the entrypoint exits. See
[Sandboxes](https://modal.com/docs/guide/sandboxes). Waiting returns an exit
result when available; Kit then marks the handle terminal.

## Edge cases

- **Destructive** to the running sandbox.
- Waiting can still error after termination has begun.
- Use [request-modal-termination](../request-modal-termination/README.md) when you
  do not need the exit code.

## Run

```sh
cd examples/typescript/modal/terminate-modal-sandbox
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/modal/terminate-modal-sandbox/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
