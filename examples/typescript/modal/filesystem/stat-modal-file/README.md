# Stat a Modal path

Print the path, type, and size of `MODAL_REMOTE_PATH`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call prints the path, type, and size. A symlink result describes the symlink. Modal reports `mode` as a four-digit octal string when it has one.

## Edge cases

- A relative path is rejected before Modal is called.
- A missing path comes back from Modal as not found.

## Run

```sh
cd examples/typescript/modal/filesystem/stat-modal-file
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
MODAL_REMOTE_PATH=/tmp/note.txt

```

## Warnings

- Needs an existing sandbox and a path at `MODAL_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/modal/filesystem/stat-modal-file/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
