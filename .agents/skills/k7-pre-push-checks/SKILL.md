---
name: k7-pre-push-checks
description: Use before pushing, marking a PR ready for review, or claiming checks pass on a gitea branch, and after `gh stack sync` publishes rewritten branches, to select the smallest tests and checks that cover the outgoing or just-published diff without reflexively running the full repository suite.
---

# Gitea Pre-Push Checks

Use this skill once before a gitea push to run only the local evidence the outgoing change actually needs. There is no local hook safety net on this checkout (`.git/hooks/` ships only the upstream `.sample` files), so the checks you select are the only local signal before CI. `make` is not on `PATH` on this Windows host by default; install it via MSYS2 or Chocolatey per [docs/build-setup.md](../../../docs/build-setup.md) or invoke it through your shell wrapper. Stack PRs defer to [k7-merging-stacked-prs](../k7-merging-stacked-prs/SKILL.md) for ordering and rebase handling; that skill points back here for evidence selection.

## Inspect the outgoing change

Confirm the checkout, branch, and upstream tracking ref:

```bash
git status --short --branch
git rev-parse --show-toplevel
git rev-parse --abbrev-ref --symbolic-full-name @{u}
```

Verify the live PR base or stack parent and inspect the complete scope against it. Supply the ref already recorded by the PR or stack; never substitute a guess.

```bash
git fetch origin <base>
git diff --stat origin/<base>...HEAD
git diff --name-only origin/<base>...HEAD
```

After a base change (rebase, parent swap), rerun the diff and reassess which behaviors the new combined scope can affect, then rerun only the checks the merge invalidated.

## Select relevant evidence

The Makefile, [docs/testing.md](../../../docs/testing.md), and [AGENTS.md](../../../AGENTS.md) own the available commands. There is no universal local baseline beyond what you run; every behavior change needs the narrowest test or purpose-built check that would fail for its regression. Add broader checks only for surfaces the diff actually reaches. Per AGENTS.md, lint only what changed.

When the change touches a fixture, helper, subprocess, or shared resource, consult [k7-ci-test-reliability](../k7-ci-test-reliability/SKILL.md) first; this skill then selects the commands and skips evidence that already passed.

- **Go packages (`models/`, `modules/`, `services/`, `routers/`, `cmd/`, `tests/integration/`):** run the owning test. Prefer `go test -run '^TestName$' ./modulepath/`; the same target is reachable as `make test-backend#TestName` (Go unit tests) or `make test-integration#TestName` (integration tests under `tests/integration/`). Run the full package once when the change affects shared contracts.
- **Templates (`templates/**/*.tmpl`):** `make lint-templates`. Run the owning integration test when the template renders server-side.
- **Frontend TS/Vue (`frontend/web_src/**`):** run the owning Vitest file with `cd frontend && pnpm exec vitest <path-filter>`. `frontend/package.json` ships no `scripts`; always run through `pnpm exec`.
- **CSS or Tailwind classes in `frontend/web_src/css/` or templates:** `make lint-css`.
- **API routes (`routers/api/**`, swagger annotations):** `make generate-swagger` and `make swagger-validate`. CI runs `make swagger-check`; locally commit the regenerated `templates/swagger/v1-swagger.generated.json` (and `v1-openapi3.generated.json` if it changed) so the diff stays honest.
- **Database-persisted structs under `models/`:** `make test-migration`. Add a new migration under `modelmigration/` when the struct changes.
- **Locale strings:** edit only `options/locale/locale_en-US.json`; other locales sync automatically. No extra check.
- **`go.mod` / `go.sum`:** `make tidy`.
- **`.go` edits:** `make fmt` first, then lint what changed (`make lint-go` accepts golangci-lint args after `--`, e.g. `make lint-go -- ./modulepath/...`).
- **GitHub Actions / Make recipes / shell scripts:** `make lint-actions` and `make lint-shell` for the surfaces touched.
- **UI behavior observable through the browser:** targeted `GITEA_TEST_E2E_FLAGS='tests/e2e/<file>.test.ts' make test-e2e`. Set `GITEA_TEST_E2E_TIMEOUT_FACTOR` only with a reason; the default is 4 on CI, 1 locally.

  ```powershell
  $env:GITEA_TEST_E2E_FLAGS='tests/e2e/repo.test.ts'; make test-e2e
  ```
- **Generated or vendored artifacts under `public/`, `options/bindata.go`:** regenerate via the owning `make` target (`make generate`, `make go-licenses`, `make svg`) and commit the result.

Do not manually repeat a passing check merely because commit or push follows.

## Full local rehearsal

Run `make checks`, `make lint`, and `make test-backend` only when the user asks, while diagnosing a CI failure, or when the change spans the repository so broadly that no narrower set is credible. These are the closest local approximation to what CI runs (`make checks` is `checks-frontend checks-backend`); treat them as the floor for genuinely cross-cutting changes, not as a default before every push.

## Protect history-rewriting pushes

Per AGENTS.md, history rewrites happen only when the user asks. When authorized, fetch the current remote branch and record its exact OID, then publish with `--force-with-lease=<branch>:<observed-oid>` so a concurrent update aborts the push. Raw `--force` is never allowed. Update PRs with new commits and a normal push whenever possible.

After any rewritten push, fetch the live heads again and re-audit unresolved review threads, approvals, mergeability, and checks. Commit hashes and inline-comment anchors from before the rewrite are not current evidence.

### Post-sync validation

For stacked PRs published by `gh stack sync`, see [k7-merging-stacked-prs](../k7-merging-stacked-prs/SKILL.md); the same rebase-then-validate discipline applies, but the sync command bundles fetch + cascade-rebase + push into one operation, so validation happens after the rewrite.

## Handle failures

If a relevant check fails before a push, stop and fix or explain the blocker. Do not push and hope CI differs.

If a failure looks environment-specific, prove it:

- Record the exact command, failing test, and platform-specific mismatch.
- Confirm the relevant non-platform evidence (e.g., the test passes on Linux CI for the same SHA).
- Prefer fixing cross-platform nondeterminism when the check is required.
- Bypass a local check only when the user explicitly asks or agrees, and report exactly what failed and why CI is expected to differ.

## Push procedure

For ordinary and authorized rebase pushes:

1. Run the selected relevant checks once.
2. Inspect any files changed by the formatter or fixer before continuing.
3. Push normally, or use the exact lease for an authorized rewritten branch.
4. Verify the remote ref matches local `HEAD`:

    ```bash
    git rev-parse HEAD origin/$(git branch --show-current)
    ```

    Windows (PowerShell): `git rev-parse HEAD "origin/$(git branch --show-current)"` works as-is.

For GitHub PRs, inspect remote CI after the push:

```bash
gh pr checks
```

Report pending checks as pending. Inspect failures before attributing them to the branch or the environment.

When `gh pr checks` reports "no checks reported" and `/actions/runs?head_sha=<sha>` returns `total_count: 0`, read mergeability before suspecting the push or a dropped GitHub event:

```bash
gh pr view <number> --json mergeable,mergeStateStatus
```

GitHub creates no `pull_request` workflow runs while a PR is `CONFLICTING`/`DIRTY`, so the absent signal is the conflict, not infrastructure. Resolving the conflict is the only fix; empty commits, `--allow-empty` pushes, draft/ready toggles, and revert-and-restore bounces all leave `total_count` at zero and add junk history. Confirm the conflicting paths with `git merge-tree --write-tree HEAD origin/<base>` when the branch cannot be merged locally yet.