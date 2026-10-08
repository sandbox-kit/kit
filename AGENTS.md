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
to their SDKs. Keep initialized SDK clients private, and do not expose a generic
provider-specific client type or raw SDK accessor as the shared public API.
Sandbox, sandbox adapters, and harness adapters remain separate modules. Applications configure sandbox.NewClient through the generated Config.
Optional provider factories initialize/authenticate SDKs from that common contract;
the sandbox client owns SDK cleanup through Close.
Optional provider integrations import and depend on their own official provider SDK.
Providers own typed SDK method/configuration/response mapping. Applications install
only the adapters they need and select their provider in Config.Provider. Do not use
runtime SDK reflection or require application callbacks for standard operations.
Public SDK configuration and responses use generated native language types;
protobuf remains a schema/generation input and is not required in caller code.
Creation uses the common protobuf request/response and returns a sandbox handle;
lifecycle operations belong to that handle and are implemented separately.
Preserve absent versus explicit configuration values, units, and source kinds.
Bindings must reject unsupported intent before making provider calls.

Validation and portable field mappings belong in versioned YAML under `specs/`.
Keep semantic rules independent of target-language code; SDK member/type bindings
belong under the language's binding section. Go validator tag translation belongs
in the Go emitter. Generate shared validation and expressible mappings rather
than maintaining duplicate handwritten rules. Add vocabulary operations with
semantic tests; do not embed raw language code in YAML. Keep remaining SDK
orchestration explicitly documented until a portable model is implemented.

Use the naming conventions in docs/naming.md. Public Go SDK types belong in
package sandbox; optional provider packages expose New and Provider, with private
backend implementations. Keep initialisms and enum naming in the Go naming emitter.
Generated Go filenames end in .gen.go. Preserve numeric enum values and explicit
serialized field names when changing target-language identifiers.
