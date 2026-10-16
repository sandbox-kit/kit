# Examples

Run Sandbox Kit through independent Go and TypeScript projects. Each selects one provider,
loads its own `.env`, and demonstrates one creation, lifecycle, cleanup, or
error-handling operation.

## Choose an example

| Project                                                                            | Demonstrates                                                 |
| ---------------------------------------------------------------------------------- | ------------------------------------------------------------ |
| [All TypeScript examples](typescript/README.md)                                    | Complete Daytona and Modal TypeScript example matrix         |
| [TypeScript Daytona creation](typescript/daytona/create-daytona-sandbox/README.md) | TypeScript client, default snapshot, and explicit cleanup    |
| [TypeScript Modal creation](typescript/modal/create-modal-sandbox/README.md)       | TypeScript client, image, resources, isolation, and policies |
| [Modal creation](go/modal/create-modal-sandbox/README.md)                          | Token-pair auth, app/environment scope, registry image       |
| [Modal resources](go/modal/create-modal-with-resources/README.md)                  | CPU, memory, lifetime, and validation tests                  |
| [Modal runtime](go/modal/create-modal-with-runtime/README.md)                      | Working directory, entrypoint, and PTY                       |
| [Modal gVisor](go/modal/create-modal-gvisor/README.md)                             | Container isolation                                          |
| [Modal Linux VM](go/modal/create-modal-linux-vm/README.md)                         | Linux VM isolation                                           |
| [Modal placement](go/modal/create-modal-with-placement/README.md)                  | Cloud and optional region preference                         |
| [Modal labels](go/modal/create-modal-with-labels/README.md)                        | Environment variables and tags                               |
| [Modal network](go/modal/create-modal-with-network/README.md)                      | Managed outbound proxy                                       |
| [Modal idle timeout](go/modal/create-modal-idle/README.md)                         | Termination after 60 idle seconds                            |
| [Modal readiness](go/modal/create-modal-ready/README.md)                           | Wait until a readiness command succeeds                      |
| [Modal verbose](go/modal/create-modal-verbose/README.md)                           | Verbose provider logging                                     |
| [Modal identity](go/modal/create-modal-with-identity/README.md)                    | Workload identity token                                      |
| [Daytona creation](go/daytona/create-daytona-sandbox/README.md)                    | API-key auth and the default snapshot                        |
| [Daytona snapshot](go/daytona/create-daytona-from-snapshot/README.md)              | Stock `daytona-small` snapshot                               |
| [Daytona image](go/daytona/create-daytona-from-image/README.md)                    | Public image with CPU, memory, and disk                      |
| [Daytona runtime](go/daytona/create-daytona-with-runtime/README.md)                | Toolbox language, version, user, and source-driven isolation |
| [Daytona ephemeral](go/daytona/create-daytona-ephemeral/README.md)                 | Delete-on-stop from the default snapshot                     |
| [Daytona linked](go/daytona/create-daytona-linked/README.md)                       | Parent sandbox and an ephemeral linked child                 |
| [Daytona policies](go/daytona/create-daytona-with-policies/README.md)              | Disabled pause, delayed deletion, and SDK request tests      |
| [Modal errors](go/modal/handle-modal-errors/README.md)                             | Unsupported disk override and common error-kind handling     |
| [Daytona errors](go/daytona/handle-daytona-errors/README.md)                       | Validation and common error-kind handling                    |
| [Daytona stop](go/daytona/stop-daytona-sandbox/README.md)                          | Verified retrieval, capability discovery, and stop           |
| [Modal terminate](go/modal/terminate-modal-sandbox/README.md)                      | ID attachment, capability discovery, and termination         |
| [Daytona get](go/daytona/get-daytona-sandbox/README.md)                            | Current metadata, portable state, and capabilities           |
| [Daytona refresh](go/daytona/refresh-daytona-sandbox/README.md)                    | Refresh cached provider metadata                             |
| [Daytona start](go/daytona/start-daytona-sandbox/README.md)                        | Start an existing retained sandbox                           |
| [Daytona pause](go/daytona/pause-daytona-sandbox/README.md)                        | Memory-preserving VM pause                                   |
| [Daytona archive](go/daytona/archive-daytona-sandbox/README.md)                    | Archive a retained container filesystem                      |
| [Daytona delete](go/daytona/delete-daytona-sandbox/README.md)                      | Wait for irreversible deletion                               |
| [Daytona close handle](go/daytona/close-daytona-handle/README.md)                  | Release a local handle without changing the sandbox          |
| [Modal attach](go/modal/attach-modal-sandbox/README.md)                            | Construct an ID-only handle without remote lookup            |
| [Modal request termination](go/modal/request-modal-termination/README.md)          | Terminate without waiting for an exit code                   |
| [Modal close handle](go/modal/close-modal-handle/README.md)                        | Release a local handle without terminating the sandbox       |

Requires Go 1.26.9 or newer and this checkout. In the selected project, copy
`.env.example` to `.env`, fill in credentials, and run `go run .`.
See [Go setup and verification](go/README.md).

Creation examples provision real resources. Lifecycle examples operate on the
sandbox ID in `.env`; stop, terminate, archive, and delete can change or remove
that resource. Error-handling examples fail locally and create no cloud
resources. `Client.Close` and `Sandbox.Close` release local resources only.

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

TypeScript verification builds and tests the SDK packages, discovers every
TypeScript example package, and builds the complete example matrix through
`go run ./cmd/sandbox-kit test typescript` from `tooling/`.
