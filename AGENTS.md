# Repository instructions

Keep shared contracts under `proto/`. Runtime SDK code belongs in
`sdks/<language>/` and usage examples in `examples/<language>/`.
All generators are implemented in the single Go module under `tooling/`;
target-language emitters belong in `tooling/internal/<language>/`.
Do not put language runtime dependencies in shared contracts.
Keep contract declarations language-neutral; implement constructor conventions
and language-specific APIs in each generator.

For every implemented public feature, update that language's runnable example
and guide. Maintain one example per implemented language instead of parallel
tutorials. Update `examples/README.md` when adding a new language.
Keep examples synchronized with the actual implemented API. Clearly label
planned APIs and demonstration providers; do not present them as working
production integrations. Verify changed examples before handing off the feature.
Update relevant project Markdown whenever the API, architecture, dependencies,
or file layout changes. Check links and commands, and distinguish implemented
behavior from planned contracts and integrations.
Keep temporary generator binaries, verification caches, and scratch files outside
the checkout, in the operating system's temporary/cache directories. Do not
create a repository-local `work/` folder. Generated SDK source remains in its
declared output directories. Generation must clean up its temporary binaries.
Use reusable Cobra commands as developer tooling entry points rather than
duplicating workflows in task-runner files. Command constructors must avoid
global flag state and side effects, and subprocesses must receive the command's
context.

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
