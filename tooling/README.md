# Developer tooling

Cobra CLI for regenerating protobuf contracts and verifying SDK modules.
Applications using Sandbox Kit do not depend on it. Provider and shared SDK
helpers are handwritten under `sdks/`.

## Usage

Requires Go 1.26.9 or newer. Contract generation and protobuf verification
require the Buf CLI.

From this directory:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate contracts-go
GOWORK=off go run ./cmd/sandbox-kit test proto --against '.git#branch=main,subdir=proto'
GOWORK=off go run ./cmd/sandbox-kit test go
GOWORK=off go run ./cmd/sandbox-kit generate contracts-typescript
GOWORK=off go run ./cmd/sandbox-kit test typescript
```

Or install the command:

```sh
GOWORK=off go install ./cmd/sandbox-kit
sandbox-kit generate contracts-go
sandbox-kit test go
```

Add `GOBIN`, or `$(go env GOPATH)/bin` when unset, to `PATH`.
The command discovers the repository from the current directory.
Use `--repo /absolute/path/to/kit` elsewhere.

| Option         | Purpose                  |
| -------------- | ------------------------ |
| `--repo`       | Explicit repository root |
| `--go-binary`  | Go executable            |
| `--npm-binary` | npm executable           |
| `--buf`        | Buf CLI executable       |
| `--help`       | Commands and usage       |

## What commands do

`test proto` checks canonical formatting, lint, and compilation with Buf. Pass
`--against` to add a breaking-change comparison against a Git or Buf input.

`generate contracts-go` / `contracts-typescript` run `buf generate` using the
root `buf.gen.yaml` (Go + TypeScript plugins) and refresh contract packages
under `sdks/<language>/contracts/`. Protovalidate options travel with the
descriptors; SDKs evaluate create and client-config rules at runtime
(`secure_endpoint` stays handwritten). `contracts-typescript` also runs
`npm run build` in `sdks/typescript/contracts` so `@sandbox-kit/contracts`
exports compiled descriptors from `dist/` via `@sandbox-kit/contracts/sandbox`
and `@sandbox-kit/contracts/client`.

`test go` tests the public SDK, tooling, both providers, and every independent
Go example module registered by the command. Example tests use dummy credentials
and local servers; they create no cloud resources.

`test typescript` builds contracts first, then builds and tests the shared
package and both providers, discovers every TypeScript example containing a
`package.json`, and builds each example without making provider calls. Install
TypeScript workspace deps once from `sdks/typescript/` (`npm install`) before
running it. Manual npm publish steps live in
[`sdks/typescript/README.md`](../sdks/typescript/README.md).

Commands forward cancellation to subprocesses and run Go module commands with
`GOWORK=off`. Caches stay outside the checkout.

`GOWORK=off` ignores the root workspace for each command so every module is
verified with its own `go.mod` and local replacements. The root `go.work` is for
working across modules in an editor. Run `go version` in the same terminal to
confirm Go 1.26.9 or newer; restart the editor's Go language server after upgrading.

`go work sync` can update dependency versions in individual `go.mod` files.
After using it, run `GOWORK=off go mod tidy` inside each affected module to update
its `go.sum`, then run `sandbox-kit test go` again.

## Layout

| Directory            | Responsibility                                    |
| -------------------- | ------------------------------------------------- |
| `cmd/sandbox-kit/`   | Developer CLI                                     |
| `internal/commands/` | CLI command implementations                       |
| `internal/spec/`     | Strict YAML loaders for portable rule inventories |

Inputs live in root `proto/` and `specs/`. Runtime output lives in `sdks/go/`
and `sdks/typescript/`. Protobuf provides descriptors for contract generation;
it does not add a Kit network service or public runtime message dependency.

See [API declarations](../docs/api-declarations.md) and
[generation specifications](../docs/code-generation.md) for the vocabulary,
validation semantics, and extension workflow.
