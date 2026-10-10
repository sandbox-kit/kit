# SDK generation tooling

This Go module contains the Cobra CLI and all SDK generators. Applications using
Sandbox Kit do not depend on it. Future language emitters will also be written
in Go; current output targets Go.

## Run the CLI

Requires Go 1.26.1 or newer. Generation also requires `protoc`.

From this directory:

```sh
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test go
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

| Option        | Purpose                      |
| ------------- | ---------------------------- |
| `--repo`      | Explicit repository root     |
| `--go-binary` | Go executable                |
| `--protoc`    | Protobuf compiler executable |
| `--help`      | Commands and usage           |

## What commands do

`generate go` rebuilds annotation bindings and emits shared native types,
validation, client/provider contracts, SDK constructors, state capture,
configuration assembly, field mappings, owned copies, common errors, field-path
constants, and response-origin validation.

`test go` tests the public SDK, tooling, both providers, and six independent
example modules. Example tests use dummy credentials and local servers; they
create no cloud resources. The minimal creation examples compile and run reporting tests. Dedicated error
examples also run validation and kind-switch paths without cloud creation.

Commands forward cancellation to subprocesses and run Go module commands with
`GOWORK=off`. Generation uses an external temporary directory and cleans up
its binaries. Caches stay outside the checkout.

## Layout

`internal/go/error_gen/` owns shared error-expression emission.
`internal/go/schema/` resolves descriptor-backed paths. The type emitter also
generates public path constants and response provenance validation.

Response/error generation reads `specs/contracts.yaml` and `specs/errors.yaml`.
The latter declares cause classification, context reuse, detail preservation,
and classification precedence. Unsupported policy alternatives fail generation.
Native runtime
helpers use `g.P(...)` in `internal/go/types_gen/runtime.go`; provider error
classifiers are emitted from portable protocol rules and native SDK bindings.

| Directory                     | Responsibility                                                    |
| ----------------------------- | ----------------------------------------------------------------- |
| `cmd/sandbox-kit/`            | Developer CLI                                                     |
| `cmd/protoc-gen-kit-go/`      | Protobuf output plugin                                            |
| `internal/spec/`              | Strict YAML loader and portable rule types                        |
| `internal/go/types_gen/`      | Native types, copying, error runtime, paths, and origin validation |
| `internal/go/error_gen/`      | Consistent shared error expressions                               |
| `internal/go/schema/`         | Descriptor resolution and canonical field paths                   |
| `internal/go/validation_gen/` | Shared validation and enum checks                                 |
| `internal/go/mapping_gen/`    | Field conversions and policy mappings                             |
| `internal/go/client_gen/`     | Client contracts, deadlines, and delegation                       |
| `internal/go/provider_gen/`   | SDK construction, client assembly, state capture, provider checks, error classification |
| `internal/go/naming/`         | Go field and enum naming                                          |
| `internal/gen/`               | Generated annotation bindings                                     |

Inputs live in root `proto/` and `specs/`. Runtime output lives in `sdks/go/`
and ends in `.gen.go`. Protobuf provides descriptors for generation; it does
not add a Kit network service or public runtime message dependency.

See [generation specifications](../docs/code-generation.md) for the vocabulary,
validation semantics, and extension workflow.
