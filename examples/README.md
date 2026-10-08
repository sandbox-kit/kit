# Examples

## Go

| Project | What it demonstrates | Configuration |
| --- | --- | --- |
| [Create a Modal sandbox](go/create-modal-sandbox/README.md) | Token-pair auth, app/environment scope, registry image creation | Its own `.env` |
| [Create a Daytona sandbox](go/create-daytona-sandbox/README.md) | API-key auth, default-snapshot creation | Its own `.env` |

These are independent Go projects. Choose one, copy its blank `.env.example` to
`.env`, fill in the required values, and run `go run .` from that directory.
The programs use Sandbox Kit's common client API and only their selected provider.

They create real sandboxes. They do not execute commands or stop/delete sandboxes;
SDK cleanup is distinct from sandbox cleanup. Read the project's guide first.

Repository verification compiles both projects without executing them:

```sh
cd tooling
GOWORK=off go run ./cmd/sandbox-kit test go
```

Only Go examples are implemented today.
