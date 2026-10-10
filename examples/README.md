# Examples

Run Sandbox Kit through independent Go projects. Each selects one provider,
loads its own `.env`, and demonstrates creation or error handling.

## Choose an example

| Project                                                       | Demonstrates                                             |
| ------------------------------------------------------------- | -------------------------------------------------------- |
| [Modal creation](go/modal/create-modal-sandbox/README.md)           | Token-pair auth, app/environment scope, registry image   |
| [Modal resources](go/modal/create-modal-with-resources/README.md)   | CPU, memory, lifetime, and validation tests              |
| [Modal runtime](go/modal/create-modal-with-runtime/README.md)       | Working directory, entrypoint, and PTY                   |
| [Modal gVisor](go/modal/create-modal-gvisor/README.md)              | Container isolation                                      |
| [Modal Linux VM](go/modal/create-modal-linux-vm/README.md)          | Linux VM isolation                                       |
| [Modal placement](go/modal/create-modal-with-placement/README.md)   | Cloud and optional region preference                     |
| [Modal labels](go/modal/create-modal-with-labels/README.md)         | Environment variables and tags                           |
| [Modal network](go/modal/create-modal-with-network/README.md)       | Blocked outbound access                                  |
| [Modal idle timeout](go/modal/create-modal-idle/README.md)          | Termination after 60 idle seconds                        |
| [Modal readiness](go/modal/create-modal-ready/README.md)            | Wait until a readiness command succeeds                  |
| [Modal verbose](go/modal/create-modal-verbose/README.md)            | Verbose provider logging                                 |
| [Modal identity](go/modal/create-modal-with-identity/README.md)     | Workload identity token                                  |
| [Daytona creation](go/daytona/create-daytona-sandbox/README.md)       | API-key auth and the default snapshot                    |
| [Daytona snapshot](go/daytona/create-daytona-from-snapshot/README.md) | Stock `daytona-small` snapshot                           |
| [Daytona image](go/daytona/create-daytona-from-image/README.md)       | Public image with CPU, memory, and disk                  |
| [Daytona ephemeral](go/daytona/create-daytona-ephemeral/README.md)    | Delete-on-stop from the default snapshot                 |
| [Daytona linked](go/daytona/create-daytona-linked/README.md)          | Parent sandbox and an ephemeral linked child             |
| [Daytona policies](go/daytona/create-daytona-with-policies/README.md) | Disabled pause, delayed deletion, and SDK request tests  |
| [Modal errors](go/modal/handle-modal-errors/README.md)              | Unsupported disk override and common error-kind handling |
| [Daytona errors](go/daytona/handle-daytona-errors/README.md)          | Validation and common error-kind handling                |

Requires Go 1.26.9 or newer and this checkout. In the selected project, copy
`.env.example` to `.env`, fill in credentials, and run `go run .`.
See [Go setup and verification](go/README.md).

Creation examples provision real resources. Error-handling examples fail locally
and create no cloud resources. Client cleanup releases SDK resources;
sandbox cleanup uses provider tools or configured lifetime policies.

## Verify locally

From the repository root:

```sh
cd tooling
GOWORK=off go run ./cmd/sandbox-kit test go
```

Verification compiles every example module. The resources, policies, and
error-handling projects also run local handler and request tests.
Resource/policy tests use dummy
credentials and local servers, exercising validation and SDK request mapping
without cloud creation.

Only Go examples are implemented today.
