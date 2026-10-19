# Create a Modal sandbox that waits using the client timeout

Create with `waitFor: Ready` and a command probe, relying on `Client` `timeout` (no `provisioning.timeout` in the request).

## What this means

Kit injects the resolved create deadline into `provisioning.timeout` after validation so Modal ready-wait and the call abort share one deadline.

## Edge cases

- You can still set `provisioning.timeout` explicitly; call-option and client timeouts override when resolved.
- A readiness probe is still required for `waitFor: Ready`.

## Run

```sh
cd examples/typescript/modal/create-modal-client-timeout
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
