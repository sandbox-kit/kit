# Create a Modal sandbox and wait until ready

Create with `readiness: { command: { argv: ["true"] } }`,
`provisioning: { waitFor: WaitCondition.Ready, timeout: 120_000 }`.

## What this means

Modal readiness probes (TCP or exec) run inside the sandbox until they succeed.
Kit's ready wait maps to waiting until the probe passes before returning the
handle. See
[Sandboxes → Readiness Probes](https://modal.com/docs/guide/sandboxes).

Ready wait requires a probe and a positive `provisioning.timeout` (ms).

## Edge cases

- Missing probe or non-positive timeout → local rejection.
- A failing/slow probe can time out create; check the Modal dashboard if unsure
  whether a sandbox remained.
- Entrypoint `sleep 300` keeps the process alive after ready.

## Run

```sh
cd examples/typescript/modal/create-modal-ready
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

- [go twin](../../../go/modal/create-modal-ready/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
