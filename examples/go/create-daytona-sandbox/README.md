# Create a Daytona sandbox

```sh
cd examples/go/create-daytona-sandbox
cp .env.example .env   # Skip if .env already exists.
# Fill in .env with your dashboard credentials.
go run .
```

`main.go` authenticates through the unified SDK, creates one sandbox, and prints
its ID. There are no flags or provider selectors. This project depends only on
Sandbox Kit and its Daytona provider.

Uses the default snapshot, following [Daytona’s Go quickstart](https://www.daytona.io/docs/en/go-sdk/).

`.env` is Git-ignored. Running the example creates a real sandbox.
`client.Close` releases SDK resources; it does not delete the sandbox.
No command execution or sandbox lifecycle operations are included.

Compile/check without provisioning: `GOWORK=off go test ./...`.
