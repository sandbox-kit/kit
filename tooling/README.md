# SDK generation tooling

Developer commands and all SDK emitters live in this Go module. They are not
runtime dependencies of applications using Sandbox Kit. Currently only Go output
is implemented; future language emitters will also be written in Go.

## Requirements

* Go 1.26.1 or newer.
* `protoc` for generation; verification does not require it.

## Generate and verify

From the repository root:

```sh
cd tooling
GOWORK=off go run ./cmd/sandbox-kit generate go
GOWORK=off go run ./cmd/sandbox-kit test go
```

Generation rebuilds annotation bindings, then generates native SDK types,
validators, client/provider contracts, constructors, and supported mappings.
Verification tests the public SDK, tooling, both providers, and both standalone
example projects. It does not execute the examples or provision cloud resources.

To install the reusable Cobra CLI:

```sh
GOWORK=off go install ./cmd/sandbox-kit
sandbox-kit generate go
sandbox-kit test go
```

Ensure `GOBIN`, or `$(go env GOPATH)/bin` when unset, is on `PATH`.
The installed command discovers the checkout from the current directory. Use
`--repo /absolute/path/to/kit` from elsewhere. `--go-binary` and `--protoc` select
executables; `sandbox-kit --help` lists commands/options. Subprocesses receive
cancellation and run with `GOWORK=off`.

## Inputs and outputs

| Input | Purpose |
| --- | --- |
| `proto/kit/sandbox/v1/` | Shared client and creation types/method declarations |
| `proto/kit/providers/` | Provider SDK bindings/declarations |
| `specs/client.yaml` | Common config/provider selection and auth validation |
| `specs/validation.yaml` | Shared creation validation rules |
| `specs/providers/` | Portable field conversions and native SDK/runtime bindings |

```text
cmd/sandbox-kit/             Developer CLI
cmd/protoc-gen-kit-go/       Go output plugin
internal/spec/              Strict YAML loader and portable rule model
internal/go/types_gen/       Native structs and validation tags
internal/go/validation_gen/  Field and cross-field validators
internal/go/mapping_gen/     Typed mappings and conversions
internal/go/client_gen/      Client/provider/backend contracts and delegation
internal/go/provider_gen/    Provider constructors and private SDK backends
internal/go/naming/          Go initialisms and enum naming
internal/gen/               Generated annotation bindings used by tooling
```

Generated SDK source lives under `sdks/go/` and ends in `.gen.go`. Edit the
schema/spec/emitter, then regenerate; do not hand-edit generated files. Generator
binaries use a temporary directory outside the checkout and are cleaned up on
exit. The normal external Go cache is used; no repository `work/` folder is needed.

Protobuf supplies build-time descriptors, not public runtime message objects or
gRPC transport. Go runtime validation uses `go-playground/validator/v10`.
Mappings are generated where the vocabulary supports them; SDK construction/call
sequences and remaining provider semantics stay in typed implementation helpers.
See [the spec format and extension workflow](../docs/code-generation.md) and
[naming conventions](../docs/naming.md).
