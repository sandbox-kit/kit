# Examples

Run Sandbox Kit through independent Go projects. Each selects one provider,
loads its own `.env`, and creates one sandbox.

## Choose an example

| Project                                                       | Demonstrates                                            |
| ------------------------------------------------------------- | ------------------------------------------------------- |
| [Modal creation](go/create-modal-sandbox/README.md)           | Token-pair auth, app/environment scope, registry image  |
| [Daytona creation](go/create-daytona-sandbox/README.md)       | API-key auth and the default snapshot                   |
| [Modal resources](go/create-modal-with-resources/README.md)   | CPU, memory, lifetime, and validation tests             |
| [Daytona policies](go/create-daytona-with-policies/README.md) | Disabled pause, delayed deletion, and SDK request tests |

Requires Go 1.26.1 or newer and this checkout. In the selected project, copy
`.env.example` to `.env`, fill in credentials, and run `go run .`.
See [Go setup and verification](go/README.md).

Live examples provision real resources. Client cleanup releases SDK resources;
sandbox cleanup uses provider tools or configured lifetime policies.

## Verify locally

From the repository root:

```sh
cd tooling
GOWORK=off go run ./cmd/sandbox-kit test go
```

Verification compiles all four example modules. Resource/policy tests use dummy
credentials and local servers, exercising validation and SDK request mapping
without cloud creation.

Only Go examples are implemented today.
