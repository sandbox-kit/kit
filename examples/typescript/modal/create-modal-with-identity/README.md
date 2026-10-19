# Create a Modal sandbox with workload identity

Create with **`security.workloadIdentity = true`**.

## What this means

Enables a Modal OIDC workload identity token inside the sandbox so workloads can
authenticate to identity-aware APIs. Kit client auth still uses your
`MODAL_TOKEN_ID` / `MODAL_TOKEN_SECRET`.

## Edge cases

- Workspace OIDC configuration must allow workload identity or create fails.
- This flag is not a substitute for API token credentials on the Kit client.

## Run

```sh
cd examples/typescript/modal/create-modal-with-identity
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

- [go twin](../../../go/modal/create-modal-with-identity/README.md)

- Catalog: [examples/typescript/README.md](../../README.md)
