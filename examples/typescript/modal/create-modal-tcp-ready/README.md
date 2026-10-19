# Create a Modal sandbox that waits for a TCP port

Create with a TCP readiness probe on port 8080 and `waitFor: Ready`.

## What this means

Unlike the command-probe ready example, this waits until a TCP listener accepts connections. The entrypoint starts a simple listener so the probe can succeed.

## Edge cases

- `waitFor: Ready` requires a readiness probe and a positive timeout (here `provisioning.timeout`).
- Alpine’s busybox `nc` must be available in the image.

## Run

```sh
cd examples/typescript/modal/create-modal-tcp-ready
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
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
