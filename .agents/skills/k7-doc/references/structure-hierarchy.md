# Structure and hierarchy

## Summary

Gitea's developer documentation is a flat `docs/` tree plus a small set of root files (`README.md`, `CONTRIBUTING.md`, `AGENTS.md`, `CHANGELOG.md`, `CODE_OF_CONDUCT.md`, `SECURITY.md`, `CLAUDE.md`). User and admin guides live in a separate repository. A new page is warranted when a topic has its own audience, its own cross-link graph, or its own change cadence; otherwise add a section to an existing page.

## When to add a page versus a section

Add a new `docs/<topic>.md` when at least one of these is true:

- The topic is referenced from `CONTRIBUTING.md`'s topic table (see `CONTRIBUTING.md:5`) and would otherwise overload a sibling.
- The audience differs from every sibling (contributors, release managers, refactorers, governance readers).
- The page would exceed a few hundred lines if forced into a sibling and would dominate its table of contents.
- The page will outlive the change that introduces it (governance, release process, layout reference).

Otherwise add a section to the closest sibling. The first paragraph of every page points at the next useful neighbor; adding a sibling changes that neighbor map, so update the inbound links in the same change.

## Required cross-links

The contributor docs form a small directed graph. `CONTRIBUTING.md` is the entry point; `build-setup.md`, `development.md`, and `testing.md` are the working triangle; the three `guidelines-*.md` files hang off it.

- `docs/build-setup.md` → `development.md`, `testing.md`, `CONTRIBUTING.md`
- `docs/development.md` → `build-setup.md`, `testing.md`, `CONTRIBUTING.md`, each `guidelines-*.md`
- `docs/testing.md` → `build-setup.md`, `development.md`
- `docs/guidelines-backend.md`, `guidelines-frontend.md`, `guidelines-refactoring.md` → `CONTRIBUTING.md`, `development.md`, `testing.md`
- `docs/release-management.md`, `community-governance.md` → `CONTRIBUTING.md`
- `docs/migration-repo-restructure.md` → the current layout it explains

When you add or move a page, repair these cross-links atomically. A link in a sibling that points at a renamed or moved page is a bug.

## Heading structure

`docs/*.md` pages use a flat structure:

1. H1 with the page title (`# Backend development guidelines`, `# Testing`).
2. A short orienting paragraph that names the reader, the topic, and the next useful neighbor.
3. Optional `## Background` or `## Requirements` when prerequisites must precede dependent concepts.
4. `##` sections ordered from basic use to advanced use to developer detail.
5. `##` sections for tooling, validation, and submission (in that order) when the page covers workflow.

Open every substantive `##` with one orienting sentence before tables, lists, or `###` subsections. Do not skip heading levels. Reserve `> [!NOTE]` (and other GitHub alerts) for non-blocking context the reader can skip.

## When not to add to `docs/`

- A user or admin task belongs in the docs repository at `https://gitea.com/gitea/docs`. Add a one-line link from the relevant page in this repo instead of a new `docs/` page.
- A package or directory contract belongs in a sibling `README.md` next to the code. Link from `docs/` instead of duplicating.
- A one-off rationale that code already explains belongs in a same-line code comment, not in `docs/`.

## Verifying placement

Before creating a page, run the focused checks against the current tree:

```bash
Get-ChildItem -LiteralPath docs -Filter '*.md' | Select-Object -ExpandProperty Name
```

Windows (PowerShell): same command.

Open the closest sibling and the `CONTRIBUTING.md` topic table. If a sibling already owns the topic, add a section there. If neither owns it and the topic fits the audience and cadence test above, create the page and add one row to the topic table.
