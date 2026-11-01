# Copy a file from Modal

Copy `MODAL_REMOTE_PATH` from the sandbox to `MODAL_LOCAL_PATH` on this machine. The program opens that sandbox with `client.attach`. Attach does not look the id up.

## What this means

The call copies `MODAL_REMOTE_PATH` from the sandbox to `MODAL_LOCAL_PATH` on this machine.

## Edge cases

- The sandbox path must be absolute. A blank local path is rejected before Modal is called. The error field is `local_path`.
- A missing remote file comes back from Modal as not found.
- A missing local parent directory is created. A local directory is an invalid argument, and the field is `local_path`.
- An existing local file is replaced.

## Run

```sh
cd examples/typescript/modal/copy-file-from-modal
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
MODAL_LOCAL_PATH=./note.txt

```

## Warnings

- Overwrites the local file when it already exists.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/modal/copy-file-from-modal/README.md)
- Catalog: [examples/typescript/README.md](../../README.md)
