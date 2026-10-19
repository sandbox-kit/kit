# Create a Modal sandbox with a network proxy

Create with **`network.proxyReference`** set to an existing Modal Proxy name
(`MODAL_PROXY_NAME`).

## What this means

Attaches the sandbox to a named Modal Proxy in the selected environment so
outbound traffic can be routed through that proxy. The proxy must already exist.

## Edge cases

- Missing proxy → provider error at create.
- Requires `MODAL_PROXY_NAME` in addition to token/app vars.
- Empty egress allowlists mean deny-all when you use allowlist fields; this
  example only sets the proxy reference.

## Run

```sh
cd examples/typescript/modal/create-modal-with-network
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

Required: `MODAL_PROXY_NAME` (existing proxy in the environment).

## Warnings

- Creates a **real** cloud sandbox when the request succeeds.
- Client/handle close releases **local** resources only unless a lifecycle method above says otherwise.

## Related

- [go twin](../../../go/modal/create-modal-with-network/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
