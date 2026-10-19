# Create a Modal Linux VM sandbox

Create with **`isolation.kind = LinuxVM`**.

## What this means

Requests Modal's Linux VM sandbox runtime. End the sandbox with terminate, idle
timeout, or max lifetime. VM root disk size is fixed per Modal's sandbox resource
docs.

## Edge cases

- Workspace permissions may block VM isolation.
- GPU is not available on the VM runtime path in Modal's docs.

## Run

```sh
cd examples/typescript/modal/create-modal-linux-vm
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

- [go twin](../../../go/modal/create-modal-linux-vm/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
