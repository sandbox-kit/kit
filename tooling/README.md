# SDK generation tooling

This Go module contains the Cobra CLI and retired Kit custom emitters.
Applications using Sandbox Kit do not depend on it. Standard protobuf plugins
still emit contract packages; provider/shared SDK helpers are handwritten.

## Run the CLI

Requires Go 1.26.9 or newer. Generation requires `protoc`; protobuf verification
requires the Buf CLI.

From this directory:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test proto --against '.git#branch=main,subdir=proto'
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate typescript
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Or install the command:

```sh
GOWORK=off go install ./cmd/sandbox-kit
sandbox-kit generate go
sandbox-kit test go
```

Add `GOBIN`, or `$(go env GOPATH)/bin` when unset, to `PATH`.
The command discovers the repository from the current directory.
Use `--repo /absolute/path/to/kit` elsewhere.

| Option         | Purpose                      |
| -------------- | ---------------------------- |
| `--repo`       | Explicit repository root     |
| `--go-binary`  | Go executable                |
| `--npm-binary` | npm executable               |
| `--protoc`     | Protobuf compiler executable |
| `--buf`        | Buf CLI executable           |
| `--help`       | Commands and usage           |

## What commands do

`test proto` checks canonical formatting, lint, and compilation with Buf. Pass
`--against` to add a breaking-change comparison against a Git or Buf input.

`generate go` and `generate typescript` currently find empty Kit emit targets and
leave handwritten SDK sources authoritative. Use `generate contracts-go` /
`contracts-typescript` to refresh protobuf contract packages.

`test go` tests the public SDK, tooling, both providers, and every independent
Go example module registered by the command. Example tests use dummy credentials
and local servers; they create no cloud resources. Creation examples test error
reporting, validation, or request mapping, and lifecycle examples use local
provider doubles. Dedicated error examples test local rejection and error-kind
handling.

`test typescript` builds and tests the shared package and both providers,
discovers every TypeScript example containing a `package.json`, and builds each
example without making provider calls.

Commands forward cancellation to subprocesses and run Go module commands with
`GOWORK=off`. Generation uses an external temporary directory and cleans up
its binaries. Caches stay outside the checkout.

`GOWORK=off` ignores the root workspace for each command so every module is
verified with its own `go.mod` and local replacements. The root `go.work` is for
working across modules in an editor. Run `go version` in the same terminal to
confirm Go 1.26.9 or newer; restart the editor's Go language server after upgrading.

`go work sync` can update dependency versions in individual `go.mod` files.
After using it, run `GOWORK=off go mod tidy` inside each affected module to update
its `go.sum`, then run `sandbox-kit test go` again.

## Compilation model

The retired Kit plugins still compile protobuf descriptors and YAML inventories
into `internal/model/` when invoked. With empty `specs/targets.yaml`, the CLI
does not build or run those plugins. Portable public types come from protobuf
contracts. `specs/templates.yaml` keeps field-path naming for handwritten
constants; helper callables and `specs/behaviors.yaml` signatures are empty.
`specs/generation.yaml` rejects any non-empty emission phase.

## Layout

| Directory                        | Responsibility                                       |
| -------------------------------- | ---------------------------------------------------- |
| `cmd/sandbox-kit/`               | Developer CLI                                        |
| `cmd/protoc-gen-kit-go/`         | Retired Kit Go emit plugin (no-op)                   |
| `cmd/protoc-gen-kit-typescript/` | Retired Kit TypeScript emit plugin (no-op)           |
| `internal/model/`                | Language-neutral descriptor and semantic compilation |
| `internal/go/`                   | Go Kit emit dispatcher (no-op)                       |
| `internal/typescript/`           | TypeScript Kit emit dispatcher (no-op)               |
| `internal/spec/`                 | Strict YAML loader and portable rule types           |
| `internal/gen/`                  | Generated annotation bindings                        |

Inputs live in root `proto/` and `specs/`. Runtime output lives in `sdks/go/`
and `sdks/typescript/`. Protobuf provides descriptors for tooling; it does not
add a Kit network service or public runtime message dependency.

See [API declarations](../docs/api-declarations.md) and
[generation specifications](../docs/code-generation.md) for the vocabulary,
validation semantics, and extension workflow.
