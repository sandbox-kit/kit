# Create a Modal sandbox with resources

Create one sandbox with 0.5 physical CPU cores, 512 MiB memory, and a five-minute
lifetime. This standalone module uses Sandbox Kit and only the Modal provider.

## Run with your account

Requires Go 1.26.9 or newer and this checkout. From the repository root:

```sh
cd examples/go/modal/create-modal-with-resources
cp .env.example .env
# Fill in Modal credentials and the existing app name.
go run .
```

Skip copying if `.env` already exists. Set `MODAL_TOKEN_ID`,
`MODAL_TOKEN_SECRET`, and `MODAL_APP_NAME`.
`MODAL_ENVIRONMENT` defaults to `main`. The app must exist in that environment;
see [token and app setup](../create-modal-sandbox/README.md).

The program prints the sandbox ID and requested settings. It does not report
observed resource allocation. Shell variables take precedence over `.env`,
and credentials stay Git-ignored.

Running creates a real sandbox. Closing the client releases SDK resources;
the requested maximum lifetime bounds execution. Use Modal's tools for any
additional sandbox cleanup.

## Verify without credentials

The example prints the identity origin and whether allocated resources were
reported. It also demonstrates `errors.As` with `*sandbox.Error`; local numeric
failures expose provider, operation, kind, and field. Requested resources remain
separate from reported allocation. See [contracts](../../../../docs/responses-and-errors.md).

```sh
GOWORK=off go test -v ./...
```

Tests use the public client, dummy credentials, and a local endpoint. They verify:

| Case                                        | Expected behavior                  |
| ------------------------------------------- | ---------------------------------- |
| CPU below minimum or beyond native range    | Rejected before remote requests    |
| Fractional-second or overflowing lifetime   | Rejected before remote requests    |
| App-name mutation after client construction | Captured scope remains independent |
| Explicit endpoint override                  | Rejected by the Modal integration  |
| Example options                             | Pass common validation             |

Tests do not read your `.env` and require permission to bind a local port.
Positive Modal provisioning requires running `go run .` with real credentials;
local tests do not simulate it as a successful cloud operation.

See [main.go](main.go), [tests](main_test.go), and
[creation semantics](../../../../docs/sandbox-creation.md).

## Error handling

`main.go` configures the provider and creates the sandbox. `errors.go` handles
application errors by extracting `*sandbox.Error` and switching on its `Kind`.
Ordinary `.env` errors are printed normally. Cleanup errors use the same handler
and do not replace a creation failure.

See the [provider error example](../handle-modal-errors/README.md).
