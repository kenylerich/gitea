# Recall batteries

Probes for [the taxonomy](../SKILL.md#taxonomy), tuned during gitea review rounds. Every hit needs semantic judgment — the batteries over-match by design, and they under-match by nature: each review round found cases no battery caught, so pair them with an unpatterned read of the densest prose in scope.

## Invocation rules

- The Go root lives at `modules/`, `models/`, `services/`, `routers/`, `cmd/`, `tests/integration/`, `modelmigration/`; the frontend root is `frontend/web_src/`. `go.mod` module path is `gitea.dev` (in mid-restructure to `gitea.dev/backend`); ignore module-path chatter as change narration unless it documents an in-flight import rewrite that affects behavior.
- Exclusions go last so a later include cannot re-admit them. Gitea's standing exclusions are `vendor/`, `frontend/web_src/js/vendor/`, the generated artifacts under `public/assets/`, and `templates/swagger/v1-swagger.generated.json` (regenerate via `make generate-swagger`, never edit by hand). The skill's own files quote leaked wording as calibration, so exclude this skill's directory and its `references/` subdir as well. The [docs/development.md](../../../docs/development.md) self-hits through its own quoted evidence; judge it as evidence, not usage.
- Natural-language lines carry `-i` so sentence-initial capitals hit ("This PR adds…", "Probably fine…"); the first line, which matches code patterns, stays case-sensitive — `-i` would turn `\bT\d\b` and `\bP-I\b` into noise.
- Bound complete phrases. `\bthis PR\b` must match "this PR adds" without matching "this project", "this process", or "this provider".
- A zero-hit pattern proves nothing until it matches a known positive, and a noisy pattern proves nothing until it rejects a near-miss negative. Calibrate both before trusting a corpus result.
- Target authoring-language probes at the opposite-language surface. A generic ASCII search for English residue in Chinese prose is too noisy around code and identifiers; compare the prose additions against their counterpart instead. Gitea's primary locale is `options/locale/locale_en-US.json`; other locales sync from Crowdin and are not gitea-authored.

## English battery (POSIX)

```sh
rg -n --hidden '\(decision \d|\(audit [A-Z]\d|design §|plan §|design ledger|\(B ruling|\bP-I\b|\bW\d\b|\bT\d\b' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  --glob '!public/assets/**' --glob '!templates/swagger/v1-swagger.generated.json' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
rg -n --hidden -i '\bthis PR\b|\bthis branch\b|\bthis stack\b|\blater PRs?\b|\bprevious commits?\b|\bthis commit\b' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
rg -n --hidden -i '\bused to\b|\bno longer\b|\bpreviously\b|\bthe old\b|\bwas renamed\b|\bwas moved\b' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
rg -n --hidden -i '\bv1\b|this cut|\bcut \d|\btoday\b|\bfor now\b|roadmap' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
rg -n --hidden -i 'rejected in review|review round|reviewer|as of v\d' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
rg -n --hidden -i 'probably |should be enough|should suffice|it simply|is safe —|is safe --' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
rg -n --hidden '§\d' \
  --glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration tests/integration docs frontend/web_src
```

## English battery (Windows PowerShell)

`rg` is byte-identical; only the array construction differs because PowerShell 5.1 does not have a POSIX `&&` chain and arrays cannot be joined into separate `argv` slots with `-join`.

```powershell
$rg = 'rg'
$paths = 'modules','models','services','routers','cmd','modelmigration','tests/integration','docs','frontend/web_src'
$glob  = @(
  '--glob','!vendor/**',
  '--glob','!frontend/web_src/js/vendor/**',
  '--glob','!public/assets/**',
  '--glob','!templates/swagger/v1-swagger.generated.json'
)

& $rg -n --hidden $glob `
  '\(decision \d|\(audit [A-Z]\d|design §|plan §|design ledger|\(B ruling|\bP-I\b|\bW\d\b|\bT\d\b' $paths
& $rg -n --hidden $glob -i '\bthis PR\b|\bthis branch\b|\bthis stack\b|\blater PRs?\b|\bprevious commits?\b|\bthis commit\b' $paths
& $rg -n --hidden $glob -i '\bused to\b|\bno longer\b|\bpreviously\b|\bthe old\b|\bwas renamed\b|\bwas moved\b' $paths
& $rg -n --hidden $glob -i '\bv1\b|this cut|\bcut \d|\btoday\b|\bfor now\b|roadmap' $paths
& $rg -n --hidden $glob -i 'rejected in review|review round|reviewer|as of v\d' $paths
& $rg -n --hidden $glob -i 'probably |should be enough|should suffice|it simply|is safe —|is safe --' $paths
& $rg -n --hidden $glob '§\d' $paths
```

The trailing backtick (`` ` ``) is the PowerShell line-continuation character; do not split a quoted regex across lines without it.

## Markdown battery (POSIX)

`docs/*.md` and `*.md` at the repo root carry their own change-narration profile. Add these to a gitea review pass:

```sh
rg -n --hidden -i 'this commit|this change|this PR|this branch|this stack' docs CONTRIBUTING.md AGENTS.md README.md
rg -n --hidden -i 'used to|no longer|previously|the old|was renamed|was moved|now in|v1\b|this cut|today' docs CONTRIBUTING.md AGENTS.md README.md
rg -n --hidden -i 'rejected in review|review round|reviewer|as of v\d|addressed in' docs CONTRIBUTING.md AGENTS.md README.md
rg -n --hidden -i 'probably |should be enough|should suffice|for now\b|roadmap|TODO: ?later' docs CONTRIBUTING.md AGENTS.md README.md
```

## Markdown battery (Windows PowerShell)

```powershell
$md = 'docs','CONTRIBUTING.md','AGENTS.md','README.md'
& $rg -n --hidden -i 'this commit|this change|this PR|this branch|this stack' $md
& $rg -n --hidden -i 'used to|no longer|previously|the old|was renamed|was moved|now in|v1\b|this cut|today' $md
& $rg -n --hidden -i 'rejected in review|review round|reviewer|as of v\d|addressed in' $md
& $rg -n --hidden -i 'probably |should be enough|should suffice|for now\b|roadmap|TODO: ?later' $md
```

## Locale battery

`options/locale/locale_en-US.json` is the only locale file gitea edits; other locales sync from Crowdin and are not prose targets. Inspect the file for hedging language that surfaces in the UI:

```sh
rg -n 'probably|for now|should be enough|TODO' options/locale/locale_en-US.json
```

A hit is a UI string the operator cannot rely on; rewrite as a present-tense contract and re-run `make lint-md` if any other Markdown mentions the same string.

## Swagger annotation battery

The go-swagger comment on every route under `routers/api/v1/` is a prose surface that ships into `templates/swagger/v1-swagger.generated.json`:

```sh
rg -n '// summary:' routers/api/v1 | wc -l     # baseline: should match route count
rg -n --hidden -i 'this PR|this commit|used to|no longer|this cut|roadmap' routers/api/v1
```

A zero-hit `summary:` audit means a route is missing its summary line; the generated spec will reflect that. Re-run `make generate-swagger` after every edit.

## Chinese counterpart battery

Gitea locales other than `en-US` are synced from Crowdin; gitea does not author bilingual Markdown. When a fork or doc branch does carry Chinese counterparts (`*.zh.md`), run the equivalent of:

```sh
# Change or review narration in Chinese counterparts.
rg -n --hidden '评审|上一?轮|旧版|老的|不再|以前|本版|遗留' --glob '*.zh.md' docs

# Chinese authoring-language slips in English Markdown.
rg -n --hidden '设计稿|评审|上一?轮|旧版|老的|不再|以前|本版|遗留|私有|(^|[^a-zA-Z])端([^a-zA-Z]|$)' --glob '*.md' --glob '!*.zh.md' docs

# Chinese authoring-language slips in English code comments and TS comments.
rg -n --hidden '(^[[:space:]]*(//|/\*|\*)|//|/\*)[^\r\n]*(设计稿|评审|上一?轮|旧版|老的|不再|以前|本版|遗留|私有|端)' \
  --glob '*.{ts,tsx,js,jsx,mjs,cjs,go}' --glob '!frontend/web_src/js/vendor/**' \
  modules models services routers cmd modelmigration frontend/web_src
```

A generic ASCII search for English residue in Chinese prose is too noisy around code and identifiers; compare the prose additions against their counterpart instead.

## Known false-positive families

Judged and kept during gitea review rounds; expect them again:

- **Instrumental "used to"** — "the key used to sign requests" is instrumental, not temporal. The temporal form has a subject state before it ("colors used to come from…").
- **Runtime old/new** — "the old connection drains before the new one accepts" names live objects during handover, not repo states.
- **"This PR" in process docs** — documentation *about* PR workflow ("the PR body should…", templates, `CONTRIBUTING.md`) legitimately says "PR"; the ban is on a doc adopting one PR's vantage about the code.
- **`v1` as API path or module segment** — `/v1/...` endpoints, the `gitea.dev` module path, and wire-format names are identifiers, not version stamps.
- **`§N` with a committed owner** — external standards (RFC 9110 §10.1.5) and committed docs that own their §-numbering stay citable by section.
- **Contrastive "actually" and noun "wait"** — ordinary English, not hedging; no committed line probes them, so they surface only when you extend the battery with broader hedging patterns.
- **Runtime "today" and recorded timestamps** — tests that ask for the current date use natural time, not a repository version stamp; recorded CLI output keeps its voice. Wording that reaches a user still follows the behavior-evidence rule before any edit.
- **`enhance` and other Conventional Commit types** — Conventional Commits types are not "v1 stamps"; they are required commit vocabulary. `git log --oneline --grep='^enhance'` is the audit for enhancements, not a leakage probe.
- **`TODO(name):` markers** — required discipline per `AGENTS.md` and the [k7-find-simplifications](../k7-find-simplifications/SKILL.md) skill. A `TODO` whose reason is missing or stale is a different smell (the *reason* clause is missing), not change narration.
- **Alternatives-considered sections** — "rejected" inside a design doc's genre slot is the sanctioned home, not review choreography.