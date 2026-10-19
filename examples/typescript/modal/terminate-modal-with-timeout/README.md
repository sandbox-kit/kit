# Terminate a Modal sandbox with a wait timeout

`Attach` a sandbox, check **Terminate**, then terminate with **`wait: true`**
and **`timeout: 60_000`** (ms).

## What this means

Like [terminate-modal-sandbox](../terminate-modal-sandbox/README.md), but the
single options bag also sets Kit’s abort deadline for the waited terminate.

## Edge cases

- `wait: false` returns immediately ([request-modal-termination](../request-modal-termination/README.md)).
- Timeout applies to the Kit call; Modal still owns native terminate behavior.
- Mutates (and destroys) the sandbox named in `.env`.

## Run

```sh
cd examples/typescript/modal/terminate-modal-with-timeout
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

- Catalog: [examples/typescript/README.md](../../README.md)
