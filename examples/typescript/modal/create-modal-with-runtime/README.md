# Create a Modal sandbox with runtime settings

Create with **working directory**, **entrypoint**, and **PTY** enabled.

## What this means

These fields configure the sandbox process: where it starts (`workdir`), what
command runs (`entrypoint`), and whether a PTY is attached. See Modal
[Sandboxes](https://modal.com/docs/guide/sandboxes) for workdir/entrypoint usage.

## Edge cases

- Working directory must be an absolute path.
- Entrypoint replaces the image default command for this sandbox.
- PTY affects interactive IO for exec sessions.

## Run

```sh
cd examples/typescript/modal/create-modal-with-runtime
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

- Creates a **real** cloud sandbox when the request succeeds.
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- [go twin](../../../go/modal/create-modal-with-runtime/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
