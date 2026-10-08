# Providers and backends

Install the public `sandbox` SDK and the desired optional provider package.
Construct a client with `sandbox.NewClient(sandbox.Config{Provider: modal.New(), ...})`.
The selected provider maps common configuration to its official SDK constructor,
then returns a private backend implementing the SDK's operation contract.

* `sandbox.Provider` identifies an integration and constructs a backend.
* `sandbox.Backend` performs creation and SDK cleanup.
* `sandbox.Client` owns the initialized backend and exposes the public API.

Both contracts are defined where they are consumed, in package `sandbox`.
Configuration is validated before SDK construction. Unsupported intent fails
before the corresponding SDK operation. SDK clients remain private; users supply
no mapping callbacks or raw SDK objects. Provider errors pass through unchanged.
No middle service or generated gRPC transport is added.

The provider abstraction implements the factory pattern while keeping provider
selection as the public concept. Runtime implementations are private. See
[naming decisions](naming.md) for the complete vocabulary and migration table.

Shared contracts live in `proto/`; portable rules live in `specs/`. All emitters
are written in Go in `tooling/internal/<language>/`. Other language SDKs will use
native packages while reusing the same semantic contracts.

Provider modules are optional and bring their own official SDK dependency.
Future harness integrations remain separately installable. Sandbox lifecycle,
GPU, storage, browsers, computer use and macOS/Windows support remain later work.

See [configuration](configuration.md), [sandbox creation](sandbox-creation.md),
[code generation](code-generation.md), and [the Go example](../examples/go/README.md).
