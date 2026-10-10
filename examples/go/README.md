# Go examples

Each example is an independent module with local replacements for the public
SDK and its selected provider. Requires Go 1.26.9 or newer and this checkout.

## Projects

| Project                                                                | Required configuration                               | Local verification                              |
| ---------------------------------------------------------------------- | ---------------------------------------------------- | ----------------------------------------------- |
| [create-modal-sandbox](modal/create-modal-sandbox/README.md)                 | Modal token pair, existing app, optional environment | Error-kind handling tests                       |
| [create-modal-with-resources](modal/create-modal-with-resources/README.md)   | Modal token pair, existing app, optional environment | Validation and scope tests                      |
| [create-modal-with-runtime](modal/create-modal-with-runtime/README.md)       | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-gvisor](modal/create-modal-gvisor/README.md)                   | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-linux-vm](modal/create-modal-linux-vm/README.md)               | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-with-placement](modal/create-modal-with-placement/README.md)   | Modal token pair, app, optional cloud and region     | Compiles                                        |
| [create-modal-with-labels](modal/create-modal-with-labels/README.md)         | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-with-network](modal/create-modal-with-network/README.md)       | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-idle](modal/create-modal-idle/README.md)                       | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-ready](modal/create-modal-ready/README.md)                     | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-verbose](modal/create-modal-verbose/README.md)                 | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-modal-with-identity](modal/create-modal-with-identity/README.md)     | Modal token pair, existing app, optional environment | Compiles                                        |
| [create-daytona-sandbox](daytona/create-daytona-sandbox/README.md)             | Daytona API key, optional endpoint/target            | Error-kind handling tests                       |
| [create-daytona-from-snapshot](daytona/create-daytona-from-snapshot/README.md) | Daytona API key, optional endpoint/target            | Compiles                                        |
| [create-daytona-from-image](daytona/create-daytona-from-image/README.md)       | Daytona API key, optional endpoint/target            | Compiles                                        |
| [create-daytona-ephemeral](daytona/create-daytona-ephemeral/README.md)         | Daytona API key, optional endpoint/target            | Compiles                                        |
| [create-daytona-linked](daytona/create-daytona-linked/README.md)               | Daytona API key, optional endpoint/target            | Compiles                                        |
| [create-daytona-with-policies](daytona/create-daytona-with-policies/README.md) | Daytona API key, optional endpoint/target            | SDK request and policy tests                    |
| [handle-modal-errors](modal/handle-modal-errors/README.md)                   | Modal token pair and app name                        | Unsupported disk override and kind-switch tests |
| [handle-daytona-errors](daytona/handle-daytona-errors/README.md)               | Daytona API key                                      | Local validation and kind-switch tests          |

## Run with your account

From the repository root, choose a project:

```sh
cd examples/go/daytona/create-daytona-sandbox
cp .env.example .env
# Fill in the values, then:
go run .
```

Skip copying if `.env` already exists. Every project has its own blank template.
Modal app setup is in the [Modal creation guide](modal/create-modal-sandbox/README.md).

The programs load `.env` with `godotenv`; shell variables take precedence.
The SDK itself does not load files. Credentials are Git-ignored.

Creation programs create real sandboxes and print their IDs. Each creates one
sandbox except `create-daytona-linked`, which creates a parent and a child.
Error-handling programs exercise local failures without cloud creation. Closing its client
releases SDK resources without stopping/deleting the sandbox. Use provider tools
for cleanup, or the configured lifetime policy where the example supplies one.

## Verify without credentials

From any example directory:

```sh
GOWORK=off go test -v ./...
```

The original creation examples and both error examples include handler tests.
The resources and policies projects also use dummy credentials and local servers:

- Modal tests check numeric rejection, common error details, unsupported configuration, and copied scope.
- Daytona tests run the example through the official SDK and inspect serialized
  policies, common error kinds/fields, explicit zero values, and rejection before creation calls.

Local tests require permission to bind a port. They do not establish real
provider scheduling, account availability, or execution.

The repository-wide `sandbox-kit test go` command includes every example module.
Projects added for each supported create path compile without handler tests.

## Error handling

The original creation examples and both error examples keep provider setup in
`main.go` and a small application error handler in `errors.go`. The handler
extracts `*sandbox.Error` and switches on `Kind`. The other create projects
print the error and exit. All creation programs report client cleanup errors;
a cleanup failure does not replace a creation failure.

The Modal error example handles an unsupported disk override; the Daytona error
example handles an invalid CPU request. Tests cover switch branches and ordinary Go errors.
