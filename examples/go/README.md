# Go examples

Requires Go 1.26.1 or newer and this repository checkout. Each example is a
standalone module with local replacements for the SDK and its selected provider.

| Project | Required `.env` values |
| --- | --- |
| [create-modal-sandbox](create-modal-sandbox/README.md) | `MODAL_TOKEN_ID`, `MODAL_TOKEN_SECRET`, `MODAL_APP_NAME`; environment defaults to `main` |
| [create-daytona-sandbox](create-daytona-sandbox/README.md) | `DAYTONA_API_KEY` |
| [create-modal-with-resources](create-modal-with-resources/README.md) | Modal credentials, app, and optional environment |
| [create-daytona-with-policies](create-daytona-with-policies/README.md) | Daytona API key and optional endpoint/target |

From the repository root, choose one:

```sh
cd examples/go/create-daytona-sandbox
# Or: cd examples/go/create-modal-sandbox
cp .env.example .env   # Skip if .env already exists.
# Edit .env, then:
go run .
```

Each `main.go` loads `.env`, constructs the common client, creates one sandbox,
prints its ID, and releases SDK resources. There are no CLI flags, provider
selectors, or simulated backends. The SDK itself does not load `.env`; examples
use `godotenv`. Shell environment values take precedence over file values.

`.env` is Git-ignored. `.env.example` contains only blank credentials/defaults.
Modal app setup is documented in its project guide.

Running creates real resources. `client.Close` does not stop/delete the sandbox;
use provider management tools for that cleanup. Execution and lifecycle operations
are not included.

Verify any project without provisioning:

```sh
GOWORK=off go test ./...
```

The resources/policies projects also include tests through the public SDK using
dummy credentials and local servers. `go test -v ./...` runs these without cloud
provisioning. Modal tests verify validation and scope ownership; Daytona tests
also verify the serialized creation request through the official SDK. They do not
claim to verify real provider scheduling or execution.
