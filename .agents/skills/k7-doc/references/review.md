# Documentation review

## Summary

Review by whether a contributor can complete the task the page promises, not by whether every section template exists. Verify prose against code, tests, and the file actually on disk; preserve exact contracts; never soften `must`, `may`, or `never`; run the real lint, spell, and swagger commands before approving.

## Audience test

A contributor with no Gitea context should answer these after reading the page plus the cross-links in its first paragraph: what the subject covers, how to start, where state lives, which file owns it, how it fails, and where to change it. If they have to read source to discover the public flow, restore the missing explanation. If they have to absorb unrelated internals, move those details to a sibling or a code comment.

## Evidence review

For every durable claim, find the strongest owner and cite it:

- Configuration defaults → `modules/setting/*.go` and the matching section of `custom/conf/app.example.ini`.
- API behavior → go-swagger comment on the route, the struct in `modules/structs/`, and the committed spec under `templates/swagger/`.
- Test commands and selectors → `Makefile` recipes (`make help` lists them).
- Build tags and LDFLAGS → `docs/build-source.md` and the `build/` package.
- File layout → the current tree, not a remembered description.

For every operational claim (CLI command, config snippet, env-var prefix, default value, error, platform difference) the evidence is running it. Quote only the observed output. If a claim needs a key, network, or database you do not have, name the verification owner instead of asserting behavior.

## Placement checks

- The fact has exactly one owner. If the same fact appears in a `docs/` page and a root file or `AGENTS.md`, the weaker copy is a bug — delete it and link to the owner.
- New `app.ini` options live in `custom/conf/app.example.ini` first, then in the docs cheat sheet at `https://docs.gitea.com/administration/config-cheat-sheet`; the `docs/` page only links to it.
- UI strings live in `options/locale/locale_en-US.json` only. Other locale files sync from Crowdin and must not be hand-edited (`AGENTS.md:12`).
- API changes live in the swagger comment and the matching struct, then propagate via `make generate-swagger`.

## Voice checks

- One actor per sentence; one term per concept; one instruction per sentence.
- Modality (`must`, `may`, `never`), exceptions, numbers, and timing are unchanged from the source.
- No change-history narration (`previously`, `now`, `no longer`, `used to`) outside an explicit migration page such as `docs/migration-repo-restructure.md`.
- No "we plan to" or future-tense spec language in `docs/` prose. Future work belongs in an issue, not in the doc.

## Cross-link checks

- Every page opens with a paragraph that points at the next useful neighbor.
- `CONTRIBUTING.md` topic table lists every `docs/` page that is part of the contributor workflow.
- A rename or move repairs every inbound link in the same change. `git grep` the old path before deleting the file.

## Verification

Run the focused checks while iterating, then the standing documentation checks before approving.

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

CI runs `make swagger-check`, which fails when the committed spec is stale; regenerate and commit the diff before pushing. `make lint-md` runs `markdownlint --config .config/lint/.markdownlint.yaml` against root `*.md`; for pages under `docs/`, run `markdownlint` directly with the same config against the files you changed.

Windows (PowerShell): the `make` invocations are byte-identical once `make` is on `PATH` (install via MSYS2 or Chocolatey). To call `markdownlint` directly:

```powershell
node_modules/.bin/markdownlint --config .config/lint/.markdownlint.yaml docs\<file>.md
```

After the checks pass, re-read the final diff once for factual completeness and once for brevity, navigation, and ownership. Then check that the PR description follows `AGENTS.md`: minimal, what and why only, under 1000 characters, no task or file listings, screenshots for UI changes, references by full URL, Conventional Commits title plus the `enhance` type when applicable, and an `Assisted-by:` trailer.
