---
name: docs-current-state
description: >-
  Update project READMEs and Markdown guides as current-state truth only.
  Use when writing or revising README.md, docs/, examples/**/README.md, or any
  project guide after API, layout, architecture, or dependency changes; when
  removing stale docs; or when the user asks to refine documentation.
---

# Docs: current state only

Write documentation that matches the repository **as it exists now**. Pair with
`crafting-effective-readmes` for README structure, `writing-clearly-and-concisely`
for prose, and `golang-documentation` for Go/doc-comment work. Use
`baoyu-format-markdown` only when the user asks to format or beautify markdown
layout without rewriting meaning.

This skill aligns with `.cursor/rules/docs-current-state.mdc` and `AGENTS.md`.

## Rules

1. **Present tense.** Describe what the project does and how to use it today.
2. **No history in product docs.** Do not write before/after, "previously… now…",
   migration narratives, or changelog-style sections in READMEs/guides unless
   the user explicitly asks for history or a migration guide.
3. **Replace stale content.** Delete outdated sections, examples, commands, and
   claims. Do not keep deprecated paths "for context."
4. **Update with the code change.** When APIs, layout, architecture, examples,
   or dependencies change, edit the relevant Markdown in the same task.
5. **Stay silent in chat.** Do not announce "I also updated the README" unless
   the user asked about docs. PR/commit summaries may mention docs when that is
   part of the deliverable.
6. **Verify.** Check links and commands against the repo. Match examples to the
   implemented API. Prefer one catalog of truth (`examples/README.md`) over
   duplicated incomplete tables.

## Scope

| Include                                       | Exclude                                       |
| --------------------------------------------- | --------------------------------------------- |
| Root/`docs/`/`examples/**` READMEs and guides | CHANGELOG, commit/PR history                  |
| SDK and tooling READMEs that describe usage   | Skill packages under `.agents/skills/`        |
| Out-of-scope notes stated as present facts    | Roadmap or "planned / future work" narratives |

## Workflow

1. Identify which Markdown files describe the changed behavior.
2. Read those files and the matching code/examples.
3. Rewrite so each file is a clean current description; remove contradictions.
4. For README structure or audience fit, follow `crafting-effective-readmes`.
5. For sentence-level clarity, follow `writing-clearly-and-concisely`.
6. Do not prompt about incidental README edits after feature work unless
   something remains ambiguous.
