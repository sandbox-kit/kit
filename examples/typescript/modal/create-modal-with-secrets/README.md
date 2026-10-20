# Create a Modal sandbox with a named secret

Create from an image and attach an existing Modal secret by name.

## What this means

Modal secrets are resolved by name in the configured app/environment and injected as environment variables by Modal.

## Edge cases

- `MODAL_SECRET_NAME` must already exist in the Modal app.
- Per-variable renaming is unsupported; use Modal’s secret definitions.

## Run

```sh
cd examples/typescript/modal/create-modal-with-secrets
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
MODAL_SECRET_NAME=
```

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
