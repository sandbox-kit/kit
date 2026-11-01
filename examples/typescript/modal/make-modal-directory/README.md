# Create a Modal directory

Create the directory at `MODAL_REMOTE_PATH`, including any missing parents. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call creates `MODAL_REMOTE_PATH` and any missing parents.

## Edge cases

- A relative path is rejected before Modal is called.
- Modal accepts a directory that already exists.

## Run

```sh
cd examples/typescript/modal/make-modal-directory
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
MODAL_REMOTE_PATH=/tmp/kit-files

```

## Warnings

- Creates the directory inside the existing sandbox.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/modal/make-modal-directory/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
