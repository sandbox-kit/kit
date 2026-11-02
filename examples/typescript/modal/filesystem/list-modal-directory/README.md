# List a Modal directory

Print the names in the directory at `MODAL_REMOTE_PATH`. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call prints the names in one directory. It does not walk into subdirectories.

## Edge cases

- A relative path is rejected before Modal is called.
- A missing directory comes back from Modal as not found.
- A file path fails at the provider.

## Run

```sh
cd examples/typescript/modal/filesystem/list-modal-directory
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
MODAL_REMOTE_PATH=/tmp

```

## Warnings

- Needs an existing sandbox and a directory at `MODAL_REMOTE_PATH`.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../../go/modal/filesystem/list-modal-directory/README.md)
- Catalog: [examples/typescript/README.md](../../../README.md)
