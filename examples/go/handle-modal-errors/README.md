# Handle Modal errors

Initialize a Modal client and request a 1024 MiB disk override, which Modal does not support. The application
handles the error by switching on `sandboxError.Kind`. No cloud sandbox is created.

## Run

```sh
cd examples/go/handle-modal-errors
cp .env.example .env
# Fill in your provider settings, then:
go run .
```

Skip copying when `.env` already exists. Shell variables take precedence.
The example prints `Check disk_mib: ...` for the unsupported request
(`sandbox.ErrorKindUnsupported`).

## Files

- `main.go`: provider configuration and one creation attempt.
- `errors.go`: application handling with `errors.As` and a switch on common kinds.
- Tests: local validation and each switch branch, using dummy credentials.

The same error kinds work for either provider. Ordinary errors, such as a missing
`.env`, are printed normally. Client cleanup errors use the same handler.

## Verify

```sh
GOWORK=off go test ./...
```

Tests create no cloud resources and do not verify live authentication.
See [response and error contracts](../../../docs/responses-and-errors.md).
