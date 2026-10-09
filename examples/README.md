# Examples

## Go

| Project | What it demonstrates | Configuration |
| --- | --- | --- |
| [Create a Modal sandbox](go/create-modal-sandbox/README.md) | Token-pair auth, app/environment scope, registry image creation | Its own `.env` |
| [Create a Daytona sandbox](go/create-daytona-sandbox/README.md) | API-key auth, default-snapshot creation | Its own `.env` |
| [Modal resources](go/create-modal-with-resources/README.md) | Fractional CPU, memory, lifetime, and validation tests | Its own `.env` |
| [Daytona policies](go/create-daytona-with-policies/README.md) | Disabled pause, delayed delete, SDK request and rejection tests | Its own `.env` |

These are independent Go projects. Choose one, copy its blank `.env.example` to
`.env`, fill in the required values, and run `go run .` from that directory.
The programs use Sandbox Kit's common client API and only their selected provider.

They create real sandboxes. They do not execute commands or stop/delete sandboxes;
SDK cleanup is distinct from sandbox cleanup. Read the project's guide first.

Repository verification compiles all four projects and runs the resources/policies
tests against local test servers. It does not run live cloud creation:

```sh
cd tooling
GOWORK=off go run ./cmd/sandbox-kit test go
```

Only Go examples are implemented today.
