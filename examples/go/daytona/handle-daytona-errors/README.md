# Handle Daytona errors

Initialize a Daytona client and submit one invalid CPU request. The application
handles the error by switching on `sandboxError.Kind`. No cloud sandbox is created.

## Run

```sh
cd examples/go/daytona/handle-daytona-errors
cp .env.example .env
# Fill in your provider settings, then:
go run .
```

Skip copying when `.env` already exists. Shell variables take precedence.
The example prints `Check resources.cpu_cores: ...` for the invalid request.

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
See [response and error contracts](../../../../docs/responses-and-errors.md).
