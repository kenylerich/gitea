---
name: k7-doc
description: "Use when creating, restructuring, reviewing, or auditing Gitea Markdown documentation (docs/, root *.md, in-source READMEs, custom/conf/app.example.ini, go-swagger comments, and options/locale/locale_en-US.json) so each fact has exactly one owner and the page matches the house style of docs/."
---

# Gitea documentation

## Summary

Keep one owner per fact, write for the reader who arrives first, and verify every command and path before publishing. Each piece of Gitea documentation belongs in exactly one of: a code comment, `docs/*.md`, a top-level `README.md` / `CONTRIBUTING.md` / `AGENTS.md`, `custom/conf/app.example.ini`, the external `https://gitea.com/gitea/docs` repository (for user/admin guides), `options/locale/locale_en-US.json` (for UI strings), or a go-swagger comment on the route (for the API). Never duplicate. Cross-link siblings instead of restating them, follow the style of the existing `docs/` pages, then run `make lint-md` (and `make lint-spell` when prose changed) and `git diff --check`.

## Workflow

1. Read `AGENTS.md`, the closest `docs/` sibling, and the file the new fact will live next to before deciding where it goes.
2. Pick exactly one owner from the table in [Placement](#placement). If a fact already has an owner, fix that owner and link to it from anywhere else.
3. Match the house style: short orienting intro, `##` sections, ```bash fences with a language, GitHub `> [!NOTE]` alerts for non-blocking notes, tables for variable or flag inventories, sibling cross-links in the first paragraph.
4. Fact-check every command and path. Run it; quote only what you observed. For paths, read the file or use `Get-ChildItem` / `Glob` to confirm it exists.
5. Keep one term per concept, one actor per sentence, and modality (`must`, `may`, `never`) intact. Do not narrate change history in authored prose.
6. Re-read the complete diff once for correctness and once for brevity, then run the focused checks in [Validation](#validation).

## Placement

| Fact kind | Owner | When to use |
| --- | --- | --- |
| Internal contributor process, build, test, refactor, governance | `docs/*.md` | Audience is contributors working on Gitea itself. |
| Project-level front door, contribution workflow, agent rules | root `README.md`, `CONTRIBUTING.md`, `AGENTS.md` | Reader has just cloned the repo. |
| New `app.ini` option or its semantics | `modules/setting/*.go` + `custom/conf/app.example.ini` + [docs.gitea.com configuration cheat sheet](https://docs.gitea.com/administration/config-cheat-sheet) | Code, shipped example, and admin docs stay in lock-step; the cheat sheet is the user-facing surface. |
| User or admin guide (install, upgrade, reverse proxy, backup, SSH) | `https://gitea.com/gitea/docs` | Not in this repository. |
| UI label, button, error, page title | `options/locale/locale_en-US.json` | Other locales sync from Crowdin and must not be hand-edited. |
| API endpoint, request, response, error | go-swagger comment on the route + `modules/structs/` | Regenerate with `make generate-swagger`; CI runs `make swagger-check`. |
| Package or subsystem README inside the tree | `path/to/<dir>/README.md` | Use when a directory has a distinct contract worth one page; link from `docs/` instead of duplicating. |
| Code rationale that code cannot express | Same-line or short adjacent code comment | `AGENTS.md` says: write almost none, short, same-line, explaining why. |

If a fact appears in two places, delete it from the weaker one and replace the deletion with a link.

## House style

- Open with one short paragraph that names the subject, the reader, and the next useful neighbor (see `docs/development.md:1` or `docs/build-setup.md:1`).
- Use `##` for top-level sections, `###` for sub-sections; do not skip levels.
- Put code in ```bash fences (POSIX form) and add a `Windows (PowerShell):` line with a ```powershell fence only when env-var prefixes, `&&`, `$(...)`, temp dirs, or path separators differ. Omit it when the command is byte-identical.
- Reserve GitHub `> [!NOTE]`, `> [!TIP]`, `> [!IMPORTANT]`, `> [!WARNING]`, `> [!CAUTION]` for non-blocking context the reader can skip; do not use them for primary instructions.
- Use tables for environment variables, configuration fields, and HTTP status codes; use lists for sequential steps.
- Write relative links from the current file (`../CONTRIBUTING.md`, `./testing.md`). Never fabricate a URL; verify with `Glob`/`Read` or by opening it.
- Do not narrate change history, plan future changes, or restate `AGENTS.md`. Do not duplicate a fact owned elsewhere; link it.

## Frontmatter

Exactly two keys for skill files: `name:` and a single-line third-person `description:` that starts with "Use when". Quote the value when it contains a colon. There is no repository-wide frontmatter schema; do not invent `kind`, `audience`, or `tags` for Markdown in this repo.

## References

Open only the reference the current task needs.

- [references/structure-hierarchy.md](references/structure-hierarchy.md): when a `docs/` page is warranted versus a section, required cross-links between `build-setup`, `development`, and `testing`, and heading structure.
- [references/style.md](references/style.md): Markdown conventions that apply to plain GitHub-rendered pages — controlled technical English, emphasis discipline, code fences, alert syntax, table usage.
- [references/review.md](references/review.md): a documentation-review checklist ending with the real lint, spell, and swagger commands.

Use [../k7-prose-standard/SKILL.md](../k7-prose-standard/SKILL.md) for sentence-level contract coverage, [../k7-trim-cot-leakage/SKILL.md](../k7-trim-cot-leakage/SKILL.md) for reasoning-transcript leakage, [../k7-find-simplifications/SKILL.md](../k7-find-simplifications/SKILL.md) when removing prose would change promised behavior, and [../k7-code-review/SKILL.md](../k7-code-review/SKILL.md) for review.

## Validation

Run the smallest focused checks while iterating, then the standing documentation checks.

```bash
make lint-md
make lint-spell
git diff --check
```

When the change touches an API route, struct, or swagger comment:

```bash
make generate-swagger
make swagger-validate
```

`make swagger-check` is the CI gate that fails when the committed spec is stale; run `make generate-swagger` and commit the diff before pushing. `make lint-md` runs `markdownlint --config .config/lint/.markdownlint.yaml` against root `*.md`; to lint pages under `docs/`, run the same `markdownlint` binary directly with the config against the files you changed.

`make` is not on `PATH` by default on Windows; install it via MSYS2 or Chocolatey, or call the tools directly through `node_modules/.bin/` after `make deps` has run.
