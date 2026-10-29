# Handle Modal unsupported-intent errors

Attempt create with **`diskMib` set**. Modal does not allow a disk override at
creation, so Kit returns **`ErrorKind.Unsupported`** and field
**`resources.disk_mib`** before the API call.

## What this means

Modal documents that Sandbox disk size cannot be increased (VM root disk is a
fixed size). Kit rejects the override before calling Modal. The portable field
is `resources.disk_mib`. See
[sandbox resources → Disk](https://modal.com/docs/guide/sandbox-resources).

## Edge cases

- No sandbox should be created on this intentional failure path.
- Good template for switching on `kind` / `field` / `operation`.
- Kinds are not safe automatic retry signals for other failures.

## Run

```sh
cd examples/typescript/modal/handle-modal-errors
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
```

## Warnings

- Intended to fail **locally** before a provider create call.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/modal/handle-modal-errors/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
