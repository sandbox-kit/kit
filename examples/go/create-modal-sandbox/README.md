# Create a Modal sandbox

```sh
cd examples/go/create-modal-sandbox
cp .env.example .env   # Skip if .env already exists.
# Fill in .env with your dashboard credentials and existing Modal app name.
go run .
```

`main.go` authenticates through the unified SDK, creates one sandbox, and prints
its ID. There are no flags or provider selectors. This project depends only on
Sandbox Kit and its Modal provider.

Set `MODAL_ENVIRONMENT=main` in `.env`, or use the environment containing your app.
The example defaults to `main` when the value is missing or empty and passes it
through `sandbox.Config.Scope.Environment`.

Uses `alpine:3.21`, following the creation portion of [Modal’s Go example](https://modal.com/docs/guide/sdk-javascript-go). The app must already exist.

`.env` is Git-ignored. Running the example creates a real sandbox.
`client.Close` releases SDK resources; it does not delete the sandbox.
No command execution or sandbox lifecycle operations are included.

Compile/check without provisioning: `GOWORK=off go test ./...`.
