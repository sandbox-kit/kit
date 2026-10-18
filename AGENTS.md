# Repository instructions

## Contracts and modules

- Keep shared contracts in `proto/`, portable rule inventories in versioned YAML
  under `specs/`, and runtime SDKs in `sdks/<language>/`.
- Keep the developer CLI and YAML inventory loaders in the Go module under
  `tooling/`. Kit custom emitters are retired; SDK helpers are handwritten.
- Keep contracts independent of language runtime dependencies. Each provider
  owns native constructor conventions and names.
- Keep the public SDK, optional providers, and future harness integrations in
  separate modules. Each provider depends on its official SDK.

## Public API and provider ownership

Applications construct `sandbox.NewClient` with `sandbox.Config` and select a
provider through `Config.Provider`. The provider initializes its native SDK; the
public client owns cleanup through `Close`.

- Keep initialized native SDK clients private.
- Implement a unified interface over existing SDKs; do not expose a generic raw
  SDK accessor or require application callbacks for standard mappings.
- Use typed integration code. Do not invoke SDK methods through reflection.
- Expose shared request/response types as native language types from protobuf
  contracts plus handwritten aliases. Protobuf remains a schema input.
- Creation returns a sandbox handle. Implement lifecycle methods on that handle
  separately from creation configuration.
- Preserve field presence, units, source kinds, and provider defaults.
- Reject unsupported intent before provider calls. Account quotas, permissions,
  and capacity remain provider-enforced.
- Closing an SDK client must not imply stopping or deleting cloud sandboxes.

## Generation and validation

Kit custom emitters are retired. Shared create validation
(`ValidateCreateOptions`), provider checks (`validateProviderCreate`), field
mappings, configuration assembly, path helpers, and Client/Sandbox choreography
are handwritten in `sdks/`. `specs/validation.yaml`, `specs/client.yaml`,
`specs/contracts.yaml`, and each provider's `checks:` / `groups:` remain
inventories to keep in sync with handwritten code.

Shared data contracts live in `proto/` and are emitted by standard protobuf
plugins into language contract packages; SDK runtimes alias those types by hand
rather than regenerating parallel structs. Portable error payloads (`ErrorInfo` /
`ErrorKind`) live in protobuf; error wrappers, cause chaining, and provider
classification are handwritten per language.

- Keep portable rules in YAML as inventories; implement behavior in SDK code.
- Document SDK orchestration that remains handwritten.
- Compare references with the pinned SDK version. Resolve conflicting docs
  explicitly rather than guessing semantics.
- Edit `proto/` and regenerate contract packages with
  `sandbox-kit generate contracts-go` / `contracts-typescript`; never hand-edit
  protobuf-plugin output.
- Follow [naming conventions](docs/naming.md), preserve numeric enum values,
  and preserve explicit serialized names.

## Examples and documentation

For every public feature, update its language's runnable example and guide.
Use small, independent provider-specific projects under `examples/<language>/`.
Each example demonstrates one operation with its own configuration and no provider
selector or CLI flags.

Keep examples aligned with the implemented API. Label test doubles, planned
features, and local verification limits clearly. Verify changed examples before
handoff and update `examples/README.md` when adding a language.

Update project Markdown whenever APIs, layout, architecture, or dependencies
change. Verify links and commands and distinguish current behavior from plans.

## Development tooling

Use reusable Cobra commands for workflows. Constructors must avoid global flag
state and side effects; subprocesses must receive the command context.

Keep temporary binaries, caches, and scratch files outside the checkout.
Never create a repository-local `work/` directory. Generation must clean up
temporary binaries and write SDK output only to its declared directories.
