# SDK generation tooling

This Go module contains the Cobra CLI and the Go and TypeScript SDK generators.
Applications using Sandbox Kit do not depend on it. Language emitters are
implemented in Go over one compiled semantic model.

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

`generate go` rebuilds annotation bindings and emits shared native types,
validation, client/provider contracts, SDK constructors, state capture,
configuration assembly, field mappings, owned copies, common errors, field-path
constants, and response-origin validation.

`test go` tests the public SDK, tooling, both providers, and every independent
Go example module registered by the command. Example tests use dummy credentials
and local servers; they create no cloud resources. Creation examples test error
reporting, validation, or request mapping, and lifecycle examples use local
provider doubles. Dedicated error examples test local rejection and error-kind
handling.

`generate typescript` emits shared contracts, the public client, provider
identities, lifecycle dispatch, and typed Daytona and Modal mappings. Its
lifecycle surface is compiled from annotated `SandboxService` protobuf methods;
reusable source-block emission keeps handwritten SDK orchestration separate from
the generated contract.
`test typescript` builds and tests the shared package and both providers, discovers
every TypeScript example containing a `package.json`, and builds each example
without making provider calls.

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

## Compilation and emission

The plugin compiles the supplied protobuf descriptors and YAML once per plugin invocation into `internal/model/`.
The model resolves schema identities, client methods, field paths, mappings, and
error classifications without Go syntax or SDK dependencies. Provider bindings
remain attached to the specifications for selection by each language emitter.

Portable public types come from protobuf contracts. `specs/templates.yaml` declares reusable
signatures and field-path rules for descriptor-driven helpers.
`specs/behaviors.yaml` defines the signature requirements of implemented
behaviors. `specs/languages/go.yaml` and `specs/languages/typescript.yaml`
control native representation and output conventions. Templates retain structured type
bindings; body emitters resolve explicit symbols instead of renaming text. `specs/generation.yaml` declares
output phases and ordered client instructions. Annotated `SandboxService`
methods are the semantic registry from which handle lifecycle methods and static
capability discovery are emitted.
`internal/go/generate.go` and `internal/typescript/generate.go` dispatch the
compiled program to focused native emitters. Adding another language requires an
emitter for that vocabulary and native SDK bindings; it does not require
duplicating schema resolution or semantic classification.

## Layout

`internal/go/error_gen/` owns shared error-expression emission.
`internal/go/schema/` resolves descriptor-backed paths. The type emitter also
generates public path constants and response provenance validation.

Response/error generation reads `specs/contracts.yaml` and `specs/errors.yaml`.
The latter declares cause classification, context reuse, detail preservation,
and classification precedence. Unsupported policy alternatives fail generation.
Provider error classifiers come from portable protocol rules and native SDK
bindings. Runtime facades and metadata clone helpers stay handwritten.

| Directory                              | Responsibility                                                                          |
| -------------------------------------- | --------------------------------------------------------------------------------------- |
| `cmd/sandbox-kit/`                     | Developer CLI                                                                           |
| `cmd/protoc-gen-kit-go/`               | Protobuf output plugin                                                                  |
| `internal/model/`                      | Language-neutral descriptor and semantic compilation                                    |
| `internal/go/`                         | Go phase dispatcher and native language emitters                                        |
| `internal/spec/`                       | Strict YAML loader and portable rule types                                              |
| `internal/go/declaration_gen/`         | Declaration primitives                                                                  |
| `internal/go/types_gen/`               | Paths, diagnostic field names, and related helpers (not portable public types)          |
| `internal/go/error_gen/`               | Consistent shared error expressions                                                     |
| `internal/go/schema/`                  | Descriptor resolution and canonical field paths                                         |
| `internal/go/validation_gen/`          | Shared validation and enum checks                                                       |
| `internal/go/mapping_gen/`             | Field conversions and policy mappings                                                   |
| `internal/go/provider_gen/`            | Provider checks, config assembly, and error classification                              |
| `internal/go/naming/`                  | Go field and enum naming                                                                |
| `internal/typescript/declaration_gen/` | TypeScript output conventions and declaration primitives                                |
| `internal/typescript/error_gen/`       | Shared TypeScript error construction                                                    |
| `internal/typescript/mapping_gen/`     | Compiled field, unit, policy, and origin mappings                                       |
| `internal/typescript/naming/`          | Protobuf-to-TypeScript identifier conventions                                           |
| `internal/typescript/provider_gen/`    | Provider checks, config assembly, and related helpers                                   |
| `internal/typescript/schema/`          | Descriptor-to-TypeScript type lowering                                                  |
| `internal/typescript/types_gen/`       | Retired portable-type emitter (contracts + handwritten aliases replace it)              |
| `internal/typescript/validation_gen/`  | TypeScript representation and provider-rule validation                                  |
| `internal/gen/`                        | Generated annotation bindings                                                           |

Inputs live in root `proto/` and `specs/`. Runtime output lives in `sdks/go/`
and `sdks/typescript/`, ending in `.gen.go` and `.gen.ts` respectively.
Protobuf provides descriptors for generation; it does not add a Kit network
service or public runtime message dependency.

See [API declarations](../docs/api-declarations.md) and [generation specifications](../docs/code-generation.md) for the vocabulary,
validation semantics, and extension workflow.

## Go declaration emitter structure

`internal/go/declaration_gen/` separates native concerns:

- `emitter.go`: emitter construction and common validation.
- `profile.go`: Go capabilities and output conventions.
- `types.go`: native type representation.
- `signatures.go`: receivers, parameters, results, and generics.
- `declarations.go`: objects, interfaces, aliases, enums, and callbacks.
- `symbols.go`: explicit parameter, result, receiver, field, and local references.
- `templates.go`: backend emission of bound declaration templates.

Shared behavior compatibility is checked in `internal/model/behaviors.go`.
The shared profile loader permits other language extensions and representation
choices; a backend decides which ones it supports.
