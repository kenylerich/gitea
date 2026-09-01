# Markdown style for gitea docs

## Summary

The `docs/` pages are plain GitHub-rendered Markdown. Style choices aim for scannability: short intro paragraph, controlled technical English, disciplined emphasis, fenced code with a language, GitHub alert syntax for non-blocking notes, tables for variable and flag inventories, and line-wrap behavior consistent with the existing pages.

## Short intro

Open every page with one paragraph that names the subject, the reader, and the next useful neighbor. The opening of `docs/development.md:1` is the model. The intro does not need to summarize every section; the table of contents (when present) and the section leads carry the detail.

## Controlled technical English

Follow [../k7-prose-standard/SKILL.md](../k7-prose-standard/SKILL.md) for sentence-level coverage. The non-negotiable rules for `docs/` prose:

- One actor per sentence when ambiguity can change behavior. Prefer active voice.
- One stable term per concept; do not rotate synonyms for variety.
- One instruction per sentence. Use a list for three or more steps or conditions.
- Preserve modality (`must`, `may`, `never`) and exceptions verbatim. Never weaken a contract to shorten a sentence.
- Keep paragraphs on one topic.

## Emphasis

- Reserve **bold** for the clause that changes behavior or for the comparison that matters.
- Use `code` for filenames, paths, commands, env-var names, JSON keys, and Go symbols. Use code, not italics, for `app.ini` section keys like `[database]`.
- Do not use italics for emphasis on contributor docs; they reduce scannability in the rendered GitHub view.

## Fenced code

Always set the language: ```bash for shell, ```powershell for PowerShell variants, ```text for plain output, ```go for Go snippets, ```yaml for `app.ini`-style blocks. The `markdownlint` config in `.config/lint/.markdownlint.yaml` allows the language to be omitted, but the existing `docs/` pages set it consistently; match them.

For shell commands, write POSIX form first and follow with `Windows (PowerShell):` only when env-var prefixes, `&&`, `$(...)`, temp dirs, or path separators differ. Omit the variant when the command is byte-identical. Examples that differ:

- `GITEA_TEST_E2E_FLAGS='<filepath>' make test-e2e` (POSIX) versus `$env:GITEA_TEST_E2E_FLAGS='<filepath>'; make test-e2e` (PowerShell).
- `make build` and `make test-backend` are byte-identical in both shells.
- `find frontend/web_src/js -type f` and `Get-ChildItem -LiteralPath frontend\web_src\js -Recurse -File` are not; show only the form that applies.

## GitHub alerts

Use `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, `> [!WARNING]`, and `> [!CAUTION]` for non-blocking context that the reader can skip. Do not use alerts for primary instructions; the reader who skims alerts will miss them. The existing `docs/` pages use `> [!NOTE]` for asides like `$GOPATH/bin` on PATH (`docs/build-setup.md:18`).

## Tables

Use Markdown tables for:

- Environment variables (see `docs/testing.md:130`).
- Configuration field inventories.
- HTTP status code mappings.
- File-by-file change summaries.

Use a list when order matters more than alignment. Do not nest tables inside tables.

## Line wrap and length

`.config/lint/.markdownlint.yaml` sets `line-length.stern: false` and `line_length: -1`, so markdownlint will not flag long lines. The existing pages wrap around 80–100 columns in prose and let tables and code blocks run wider. Match the wrap of the page you are editing; do not reformat surrounding prose as a side effect.

## Links

- Internal: relative links from the current file (`./testing.md`, `../CONTRIBUTING.md`). Verify the target with `Glob` or by reading the file before writing the link.
- External: full URLs. Verify they resolve before publishing; do not paraphrase a URL from memory.
- Same-page anchor: `[Section name](#section-name)` only when the page has a table of contents that uses those anchors.
- Never invent a leading-slash path (`/docs/foo.md`); it will not resolve on GitHub.

## What to avoid

- `-----` horizontal rules between sections. The existing pages rely on heading hierarchy, not rules.
- `<details>` folds inside contributor docs. Reserve folds for content that would otherwise dominate the page; if you reach for one, consider a new page instead.
- Footnote-style references and external link numbers — use inline links.
- Screenshots in `docs/`. Screenshots live in the docs repository at `https://gitea.com/gitea/docs` and in PR descriptions (`AGENTS.md` requires before/after screenshots only in PR descriptions for UI changes).
