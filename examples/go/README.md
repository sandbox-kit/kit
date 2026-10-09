# Go examples

Each example is an independent module with local replacements for the public
SDK and its selected provider. Requires Go 1.26.1 or newer and this checkout.

## Projects

| Project                                                                | Required configuration                               | Local verification           |
| ---------------------------------------------------------------------- | ---------------------------------------------------- | ---------------------------- |
| [create-modal-sandbox](create-modal-sandbox/README.md)                 | Modal token pair, existing app, optional environment | Compile                      |
| [create-daytona-sandbox](create-daytona-sandbox/README.md)             | Daytona API key, optional endpoint/target            | Compile                      |
| [create-modal-with-resources](create-modal-with-resources/README.md)   | Modal token pair, existing app, optional environment | Validation and scope tests   |
| [create-daytona-with-policies](create-daytona-with-policies/README.md) | Daytona API key, optional endpoint/target            | SDK request and policy tests |

## Run with your account

From the repository root, choose a project:

```sh
cd examples/go/create-daytona-sandbox
cp .env.example .env
# Fill in the values, then:
go run .
```

Skip copying if `.env` already exists. Every project has its own blank template.
Modal app setup is in the [Modal creation guide](create-modal-sandbox/README.md).

The programs load `.env` with `godotenv`; shell variables take precedence.
The SDK itself does not load files. Credentials are Git-ignored.

Each program creates one real sandbox and prints its ID. Closing its client
releases SDK resources without stopping/deleting the sandbox. Use provider tools
for cleanup, or the configured lifetime policy where the example supplies one.

## Verify without credentials

From any example directory:

```sh
GOWORK=off go test -v ./...
```

Minimal creation examples compile without being run. The resources/policies
projects also run tests with dummy credentials and local servers:

* Modal tests check numeric rejection, unsupported configuration, and copied scope.
* Daytona tests run the example through the official SDK and inspect serialized
  policies, including explicit zero values and rejection before creation calls.

Local tests require permission to bind a port. They do not establish real
provider scheduling, account availability, or execution.

The repository-wide `sandbox-kit test go` command includes all four modules.
