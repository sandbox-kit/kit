# Repository agent skills

This repository keeps shared agent skills in `.agents/skills/` and records their
sources and content hashes in `skills-lock.json`. Install project skills from the
repository root with `npx skills add <source> --skill <name> --agent codex --copy`.
Using `--copy` preserves the repository's existing checked-in layout.

## Selected skills

### `typescript-advanced-types`

- **Source:** `wshobson/agents`
- **Why selected:** Focused guidance for generics, conditional and mapped types,
  discriminated unions, type-safe builders, and compile-time tests. These are
  directly relevant to a generated TypeScript SDK that must preserve protobuf
  presence, provider variants, and configuration constraints.
- **Scope:** TypeScript type and public API design. It complements the existing
  Go-specific skill instead of competing with it.

### `coding-standards`

- **Source:** `iurysza/agent-skills`
- **Why selected:** Language-neutral guidance for parsing at boundaries, typed
  failures, adapter ownership, explicit dependencies, and behavior-focused tests.
  Those rules reinforce this repository's provider isolation and cross-language
  contract model.
- **Precedence:** `AGENTS.md`, specifications, and established repository
  conventions remain authoritative when a generic recommendation differs.

## Existing complementary skills

- `golang-design-patterns` covers idiomatic Go API and lifecycle choices.
- `golang-documentation` covers Go and project documentation.
- `crafting-effective-readmes`, `baoyu-format-markdown`, and
  `writing-clearly-and-concisely` cover documentation structure and prose rather
  than SDK architecture.

## Evaluated but not installed

### `typescript-patterns` from `heldinhow`

The previously suggested source could not be verified in the current public
skills registry or its GitHub source. The selected advanced-types skill covers
the useful TypeScript-specific portion without relying on an unverified package.

### `design-patterns` from `FlorianBruniaux/claude-code-ultimate-guide`

This is a broad catalog of GoF and framework patterns. It overlaps with both
`golang-design-patterns` and `coding-standards`, while its React, Angular,
NestJS, Vue, Redux, and ORM material is unrelated to this SDK. A generic pattern
catalog also risks encouraging abstractions before the repository's specs and
provider seams require them.

### `typescript-pro` and other broad TypeScript bundles

These substantially overlap with `typescript-advanced-types` and add unrelated
application and framework guidance. One focused TypeScript skill keeps activation
predictable and avoids contradictory recommendations.

### `codebase-design` from `mattpocock/skills`

Its deep-module and seam guidance is strong, but the selected `coding-standards`
skill already includes cohesive-module, adapter, and boundary guidance. Installing
both would duplicate architecture advice and vocabulary.

### Microsoft `skill-creator`

This helps authors create agent skills; it is not an SDK design or TypeScript/Go
interoperability skill. It does not contribute to implementation work in this
repository.

### SDK generators and framework-specific skills

Low-adoption SDK generator skills and React, Node backend, OpenAPI, NestJS, and
frontend-focused skills were excluded. This repository has its own protobuf and
YAML-driven generator architecture, and generic generators could conflict with
the rules in `AGENTS.md` and `specs/`.
