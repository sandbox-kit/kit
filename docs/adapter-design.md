# Unified SDK and adapter boundary

Sandbox Kit provides a common interface over existing provider SDKs. It is built
incrementally, starting with client and adapter initialization.

## Language boundaries

Shared declarations live under `proto/`. Go runtime modules are under
`sdks/go/` and its example is under `examples/go/`. Other runtime SDKs and
examples will get sibling directories when implemented.

All generators are written in Go in the single `tooling/` module. Target-specific
emitters live under `tooling/internal/<language>/`; currently only the Go emitter
exists. Future non-Go SDK consumers will not need to install the generator or Go.

Declaration names such as `Client`, `Provider`, and `Adapter` are common metadata.
Constructor conventions belong to each generator: Go emits `NewClient` and
adapter `New`; another language can emit a class constructor or factory.
The shared contracts must not acquire Go-only types or runtime dependencies.
Standard protobuf language options, such as `go_package`, remain output hints.

## Current API

The application imports its SDK, core, and the selected adapter. It initializes
the SDK using its own authentication and configuration, then injects it:

```go
provider, err := modalAdapter.New(existingModalClient)
if err != nil {
    return err
}
client, err := core.NewClient(provider)
if err != nil {
    return err
}
name := client.ProviderName()
```

Daytona uses `daytonaAdapter.New(existingDaytonaClient)` with the same
`core.NewClient`. Both produce one common `*core.Client` type.

Core has no external dependencies. Each adapter requires only core. Its type
parameter retains the supplied SDK pointer privately without importing the SDK.
Constructors validate nil injection; they do not authenticate, modify SDK
configuration, allocate connections, or close caller-owned resources.

Initialization and provider identity are the only implemented shared API.
The opaque pointer is not proof that a client has the expected SDK methods.
Actual operation mappings will need explicit typed bindings. There is no raw
SDK accessor or provider-specific result type in the common public API.

See [the complete runnable Go example](../examples/go/README.md).

## Operation contracts and mappings: planned

For each operation, define a shared protobuf request and response before adding
client methods and provider mappings. Shared responses must not be SDK-specific
sandbox objects or `any` values. Preserve provider differences through explicit
capabilities, defaults, and extensions.

Go cannot call SDK-specific methods on an opaque type parameter. Typed mappings
must live where the application imports the SDK. Add callbacks or application-owned
generation when the operation contracts require them. Those mappings must convert
SDK results to common responses before returning through the reusable adapter.
No application-binding generator is implemented in this increment.

Delegation happens locally in the application process. The provider SDK uses its
own network transport. Sandbox Kit adds no gRPC service or middle service.

Error normalization remains a separate design decision. Current errors describe
only local constructor misuse; no shared provider error taxonomy is implemented.

## Provider details for future mapping work

The example pins Modal v0.11.0 and Daytona v0.222.0. Its tests verify SDK client
injection without authenticating or contacting a provider.

* Modal accepts token ID/secret, OAuth, profiles/environments, logging, and gRPC
  customization. Daytona accepts API key or JWT/organization authentication,
  endpoint configuration, and target selection.
* Modal's SDK exposes `Client.Close()`; the application owns its cleanup.
* Modal lookup uses `Sandboxes.FromID(ctx, id, params)`; Daytona uses
  `Get(ctx, sandboxIDOrName)`.
* Modal creation needs an app, image, and creation parameters. Daytona has its
  own creation parameters and options.
* Mappings must preserve the SDK's cancellation, timeout, retry, streaming,
  pagination, and lifecycle behavior.

These findings describe the pinned SDK versions. Recheck them when designing an
operation or upgrading an SDK; the corresponding mappings do not exist yet.

## Incremental scope

Shared contracts must accommodate provider-specific authentication, isolation
including Linux VMs, resources, environment variables, regions, labels,
lifecycle states/operations, snapshot types including a provider default,
networking/linking, and optional harness support. Availability and semantics
must be verified against each provider's documentation.

Harness adapters will be separately installable. Some providers have native
harness integrations; others need ordinary sandbox execution primitives.
Construction must not install or launch a harness.

Volumes and external storage remain later extensions. GPU, computer use,
browser tools, and macOS/Windows are outside version-one implementation scope,
while the contract design must leave room for them.

The example includes both providers in one program. Applications can choose one.
Local module replacements are development conveniences; real release versions
must replace v0.0.0 requirements before publishing.
