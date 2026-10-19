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
cd examples/go/modal/create-modal-with-identity
cp .env.example .env
# Fill in the values below, then:
go run .
```

Requires Go 1.26.9+ and this checkout. Skip copying if `.env` already exists.
Local check: `GOWORK=off go test ./...`.

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

- [typescript twin](../../../typescript/modal/create-modal-with-identity/README.md)

- Catalog: [examples/go/README.md](../../README.md)
