# Go examples

Two independent projects demonstrate authentication and sandbox creation:

- [Create a Modal sandbox](create-modal-sandbox/README.md)
- [Create a Daytona sandbox](create-daytona-sandbox/README.md)

Each project contains its own `main.go`, `go.mod`, `.env.example`, and local `.env`.
Fill in the local `.env`, then run `go run .` from that project directory.
No flags, provider selection, command execution, or lifecycle operations are included.
Each project installs only its own optional provider integration.

The previous combined example and its `.env` have been removed. Relevant local
settings were migrated to the corresponding project before deletion.
