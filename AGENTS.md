# Repository instructions

For every implemented public feature, add or update a runnable example under
`examples/`, document how users run it, and update `examples/README.md`.
Maintain one runnable program at the root of `examples/`; extend it instead of
adding parallel tutorials or example modules.
Keep examples synchronized with the actual implemented API. Clearly label
planned APIs and demonstration providers; do not present them as working
production integrations. Verify changed examples before handing off the feature.
Update relevant project Markdown whenever the API, architecture, dependencies,
or file layout changes. Check links and commands, and distinguish implemented
behavior from planned contracts and integrations.

Sandbox Kit is a unified interface over existing provider SDKs. Core operations
must expose shared request and response contracts; adapters map those contracts
to their SDKs. Keep borrowed SDK clients internal, and do not expose a generic
provider-specific client type or raw SDK accessor as the shared public API.
Core, sandbox adapters, and harness adapters remain separate modules. Applications
initialize/authenticate their SDK clients and retain ownership of their lifetime.
Reusable adapter modules must not import or require official provider SDKs.
Future SDK-specific mappings belong in application-owned bindings, and their
callbacks must return shared contracts before crossing into core. Add that
machinery when operation contracts exist, rather than speculative scaffolding.
