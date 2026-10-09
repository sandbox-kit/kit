# Create a Modal sandbox with resources

This standalone project creates one sandbox using 0.5 CPU cores, 512 MiB
memory, and a five-minute lifetime. It uses the common Sandbox Kit configuration
and creation types and installs only the Modal provider. Requires Go 1.26.1+.

## Run with your account

From the repository root:

```sh
cd examples/go/create-modal-with-resources
cp .env.example .env
# Fill in MODAL_TOKEN_ID, MODAL_TOKEN_SECRET, MODAL_APP_NAME.
go run .
```

`MODAL_ENVIRONMENT` defaults to `main`. The app must already exist in that
environment. See [token and app setup](../create-modal-sandbox/README.md).
Shell environment variables take precedence over `.env`.

The program creates one real sandbox and prints its ID and requested settings.
Those settings describe the request, not a claim about observed allocation.
Closing the client releases SDK resources; it does not terminate the sandbox.
The requested lifetime bounds its execution. `.env` is ignored by Git.

## Verify without credentials

```sh
GOWORK=off go test -v ./...
```

Tests initialize the actual public client with dummy credentials and direct the
SDK to a local test server. No dashboard keys or `.env` are read by these tests.
They verify that:

- Below-minimum CPU and CPU overflow are rejected before remote requests.
- Fractional-second and overflowing lifetimes are rejected before remote requests.
- Changing the supplied app name after client initialization does not change
  captured scope: the numeric validation errors still occur, not missing-app errors.
- Unsupported Modal endpoint configuration is rejected.
- The positive example options pass common validation.

These tests require permission to bind a local port. Positive Modal provisioning
is verified only when you run `go run .` with real credentials; it is not simulated
as a successful cloud operation. The provider module also tests generated numeric
mapping boundaries directly. See [main.go](main.go), [tests](main_test.go), and
[creation semantics](../../../docs/sandbox-creation.md).
