# Create a Modal sandbox that waits using the client timeout

Create with `waitFor: WaitCondition.Ready` and a command probe, relying on
`Client` `timeout` (this example omits `provisioning.timeout`).

## What this means

Kit injects the resolved create deadline into `provisioning.timeout` after
validation so Modal ready-wait and the call abort share one deadline.

## Edge cases

- An explicit `provisioning.timeout` is allowed; call-option and client
  timeouts override when resolved.
- `waitFor: WaitCondition.Ready` requires a readiness probe.

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
- Closing the handle leaves the sandbox as it is. Closing the client closes your connection to the provider.

## Related

- Catalog: [examples/typescript/README.md](../../README.md)
