# Repository agent skills

This repository keeps shared agent skills in `.agents/skills/` and records their
sources and content hashes in `skills-lock.json`. Install project skills from the
repository root with
`npx skills add <source> --skill <name> --agent codex --agent cursor --copy`.
Using `--copy` preserves the repository's existing checked-in layout.

Cursor loads project skills reliably from `.cursor/skills/`. This repo keeps a
symlink at `.cursor/skills` → `.agents/skills` so Codex and Cursor share one
checked-in copy. Without that path, Cursor may treat `.agents/skills` entries as
rules instead of invocable skills (`/skill-name`).

## Selected skills

### `typescript-advanced-types`

- **Source:** `wshobson/agents`
- **Why selected:** Focused guidance for generics, conditional and mapped types,
  discriminated unions, type-safe builders, and compile-time tests. These matter
  for the handwritten TypeScript public facades (`CreateOptions` / `Config`),
  nested exclusive fields, and provider variants over protobuf contracts.
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

### `golang-testing`

- **Source:** `samber/cc-skills-golang`
- **Why selected:** Adds idiomatic table-driven tests, helpers, integration tests,
  HTTP tests, examples, coverage, benchmarks, and mocking guidance for the Go SDK
  and developer CLI under `tooling/`.

### `golang-error-handling`

- **Source:** `samber/cc-skills-golang`
- **Why selected:** Covers error creation, wrapping, inspection, and propagation.
  This supports stable public SDK errors and translation of provider-specific
  failures at adapter boundaries.

### `javascript-testing-patterns`

- **Source:** `wshobson/agents`
- **Why selected:** Complements `typescript-advanced-types` with focused unit,
  integration, asynchronous, and test-double practices for the TypeScript SDK.

### `code-review`

- **Source:** `mattpocock/skills`
- **Why selected:** Provides a review workflow that checks changes separately
  against repository standards and the originating specification. This is useful
  for a spec-driven SDK without duplicating implementation guidance.
- **Operational note:** Its registry assessment reported safe generation and no
  dependency alerts, but a medium Snyk rating. The installed package contains
  instructions and agent metadata only; review its workflow again when updating.

### `security-best-practices`

- **Source:** `openai/skills`
- **Why selected:** Adds trusted language-specific security review guidance for
  Go and TypeScript, especially around untrusted input, credentials, command
  execution, networking, and provider configuration.

### `python-design-patterns`

- **Source:** `wshobson/agents`
- **Why selected:** Provides Python-specific guidance for choosing composition,
  protocols, factories, strategies, adapters, repositories, and other patterns
  without importing Java-style structure into a Python SDK.

### `python-project-structure`

- **Source:** `wshobson/agents`
- **Why selected:** Covers package boundaries, public exports, `pyproject.toml`,
  source layouts, and dependency organization for a Python SDK under
  `sdks/python/`.

### `python-testing-patterns`

- **Source:** `wshobson/agents`
- **Why selected:** Adds focused pytest fixtures, parametrization, mocking,
  integration tests, and asynchronous test patterns for a Python SDK.

### `python-error-handling`

- **Source:** `wshobson/agents`
- **Why selected:** Covers exception hierarchies, contextual failures, cleanup,
  retries, and error propagation needed for a stable provider-neutral SDK API.

### `async-python-patterns`

- **Source:** `wshobson/agents`
- **Why selected:** Covers structured concurrency, cancellation, timeouts,
  resource cleanup, and async context managers for a Python SDK that wraps
  networked sandbox providers.

### `docs-current-state`

- **Source:** local (`.agents/skills/docs-current-state`)
- **Why selected:** Enforces this repository's documentation policy: write
  present-tense current-state docs, delete outdated sections, update Markdown
  with code changes, and avoid before/after or roadmap narration in product
  guides. Complements the always-on Cursor rule
  `.cursor/rules/docs-current-state.mdc`.

### `crafting-effective-readmes`

- **Source:** `softaworks/agent-toolkit`
- **Why selected:** README structure, audience fit, and templates for root,
  SDK, tooling, and example project READMEs. Prefer this when creating or
  reshaping a README; pair with `docs-current-state` so content stays current.

### `writing-clearly-and-concisely`

- **Source:** `softaworks/agent-toolkit`
- **Why selected:** Prose guidance for guides, READMEs, error messages, and
  other human-readable text. Cuts filler and AI-pattern wording.

### `baoyu-format-markdown`

- **Source:** `JimLiu/baoyu-skills` (local copy)
- **Why selected:** Formatting-only markdown passes (headings, lists, layout)
  when the user asks to beautify or reformat without changing meaning.

### `golang-documentation`

- **Source:** `samber/cc-skills-golang`
- **Why selected:** Go doc comments, example tests, and project documentation
  patterns for the Go SDK and `tooling/` CLI.

## Existing complementary skills

- `golang-design-patterns` covers idiomatic Go API and lifecycle choices.
