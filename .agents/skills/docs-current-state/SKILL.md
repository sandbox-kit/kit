---
name: docs-current-state
description: >-
  Update project READMEs and Markdown/MDX guides as current-state truth for
  end users and SDK callers (not SDK builders). Use when writing or revising
  README.md, docs/, examples/**/README.md, sdks/**/README.md,
  website/content/docs/**, or any product guide after API, layout,
  architecture, or dependency changes; when removing stale docs; or when the
  user asks to refine documentation.
---

# Docs: current state, SDK-user audience

Write documentation that matches the repository **as it exists now**, for
**people who call the SDK**. Pair with `crafting-effective-readmes` for README
structure, `writing-clearly-and-concisely` for prose, and
`golang-documentation` for Go/doc-comment work. Use `baoyu-format-markdown`
only when the user asks to format or beautify markdown layout without
rewriting meaning.

This skill aligns with `.cursor/rules/docs-current-state.mdc`,
`.cursor/rules/docs-sdk-user-audience.mdc`, and `AGENTS.md`.

## Audience

Primary readers are application developers integrating `@sandbox-kit/*` or the
Go sandbox modules. Write like a technical documentation engineer for that
audience.

| For SDK users                       | Not for SDK builders                                      |
| ----------------------------------- | --------------------------------------------------------- |
| What to call, pass, and expect back | Mapping layers, wire shapes, adapter internals            |
| Provider-specific usage differences | Inventory YAML, proto regen, contributor workflows        |
| Outcomes on the handle (`state`, …) | “Unified metadata”, ignored-option plumbing               |
| Unsupported → say so + alternative  | Negated field catalogs unless they prevent a real mistake |

Omit options the provider does not use from examples. If you must mention an
unused option, say the provider does not use it—do not say it is “accepted and
ignored.” Prefer listing fields the caller receives over listing fields that
stay empty.

## Rules

1. **Present tense.** Describe what the project does and how to use it today.
2. **No history in product docs.** Do not write before/after, "previously… now…",
   migration narratives, or changelog-style sections in READMEs/guides unless
   the user explicitly asks for history or a migration guide.
3. **Replace stale content.** Delete outdated sections, examples, commands, and
   claims. Do not keep deprecated paths "for context."
4. **Update with the code change.** When APIs, layout, architecture, examples,
   or dependencies change, edit the relevant Markdown/MDX in the same task.
5. **Stay silent in chat.** Do not announce "I also updated the README" unless
   the user asked about docs. PR/commit summaries may mention docs when that is
   part of the deliverable.
6. **Verify.** Check links and commands against the repo. Match examples to the
   implemented API. Prefer one catalog of truth (`examples/README.md`) over
   duplicated incomplete tables.
7. **Provider- and language-specific.** When behavior differs, say so under the
   right provider/language section—do not paper over with generic wording.

## Scope

| Include                                       | Exclude                                       |
| --------------------------------------------- | --------------------------------------------- |
| Root/`docs/`/`examples/**` READMEs and guides | CHANGELOG, commit/PR history                  |
| SDK package READMEs that describe usage       | Skill packages under `.agents/skills/`        |
| `website/content/docs/**` product + SDK pages | Contributor-only architecture unless asked    |
| Out-of-scope notes stated as present facts    | Roadmap or "planned / future work" narratives |

## Workflow

1. Identify which Markdown/MDX files describe the changed behavior.
2. Read those files and the matching public API / examples.
3. Rewrite so each file is a clean current description for an SDK caller;
   remove contradictions and builder-only detail.
4. For README structure or audience fit, follow `crafting-effective-readmes`.
5. For sentence-level clarity, follow `writing-clearly-and-concisely`.
6. Do not prompt about incidental README edits after feature work unless
   something remains ambiguous.
