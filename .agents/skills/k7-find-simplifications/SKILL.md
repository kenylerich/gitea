---
name: k7-find-simplifications
description: "Use when working in the gitea repo to find non-obvious simplification candidates, remove redundant comments or implementation-heavy documentation, write inline TODO/FIXME/XXX notes or a tracking issue referenced by full URL, audit superseded design docs, or fold worthwhile simplification ideas from another PR; especially for dead, duplicated, speculative, over-built, added-then-removed, or hand-rolled-where-a-dependency-exists surfaces in modules/, services/, models/, routers/, cmd/, templates/, frontend/web_src/."
---

# Finding Gitea simplifications

This skill helps turn a broad "find things to simplify" request into evidence-backed changes that remove or collapse existing gitea surface area. It is guidance, not a checklist: follow the code, keep judgment active, and prefer a few well-proven candidates over a pile of thin guesses.

## Start with repo context

- Read [AGENTS.md](../../../AGENTS.md), especially the comment and authorship rules, plus [docs/development.md](../../../docs/development.md) and [docs/testing.md](../../../docs/testing.md).
- Read [docs/guidelines-backend.md](../../../docs/guidelines-backend.md) and [docs/guidelines-frontend.md](../../../docs/guidelines-frontend.md) before judging anything under `modules/`, `services/`, `models/`, `routers/`, `cmd/`, `templates/`, or `frontend/web_src/`. Simplifications that fight the package layout, the `db.WithTx` pattern, the `cmd → routers → services → models → modules` dependency direction, the `tw-*` utility convention, or the API/swagger rules need extra evidence.
- Read [docs/guidelines-refactoring.md](../../../docs/guidelines-refactoring.md) before proposing a refactor: a refactoring PR must be forward-looking, keep behavior, include tests, and stay within the agreed scope.
- Read [docs/migration-repo-restructure.md](../../../docs/migration-repo-restructure.md) when the candidate lives in a directory affected by the in-flight move to `backend/`. A path that exists at the root *and* under `backend/` is in transition; cite both paths in the change description.
- Read [docs/build-setup.md](../../../docs/build-setup.md) and the `Makefile` (run `make help`) for the available targets. `make` is not on `PATH` on this Windows host by default; install it via MSYS2 or Chocolatey per `docs/build-setup.md`, or invoke tools through `node_modules/.bin/` after `make deps` has run.

## What counts as a strong candidate

A strong simplification removes, folds, or demotes something real and has clear evidence that the current design costs more than it buys:

- An exported function, type, configuration knob, swagger parameter, CLI subcommand, helper package, route, or template partial has no production consumer.
- Tests or docs are the only consumers, and the behavior they pin is not load-bearing.
- Two implementations mirror the same fact — duplicated helpers between `modules/` and `services/`, redundant template partials under `templates/shared/` and `templates/repo/`, an exported wrapper that only re-exports its inner type.
- A seam has methods every implementation must support but no consumer uses.
- A separate package exists only for test/demo/support code and adds publish or dependency overhead.
- A feature implements speculative product generality (multi-session support, dual LLM adapters, alternate storage backends, alternate markup renderers) with no consumer in the configuration or the docs.
- An invariant, rollback path, set of expected outputs, or special-case test exists only to protect an unused API.
- Hand-rolled code reimplements what an existing Go standard-library helper, a maintained dependency in `go.mod`, or an existing `modules/util` / `modules/htmlutil` / `modules/typesniffer` / `modules/validation` helper already provides — and the swap would delete the implementation plus its dedicated tests.
- The simplified behavior may differ slightly, but the new behavior is still reasonable and easier to explain.

Thin candidates are not enough on their own: deleting one typo, removing an intentionally documented dual adapter or backed-by-API-v1 route, or flagging "this looks complex" without call-site proof. For thin candidates, prefer a `TODO(name):` with a one-line reason at the call site over a sweeping design doc.

## Survey broadly

Use parallel subagents when the user asks for breadth or many candidates. Give each agent a domain and require evidence, not guesses. Useful domains for gitea:

- Auth and session: OAuth/LDAP/SMTP sources, webauthn, two-factor, session store, token issuance.
- Repository lifecycle: hooks, mirror sync, LFS transfer, push mirror, releases, tags, branch protection, signing.
- API and routing: `routers/api/v1/` route handlers, swagger annotations, the option registration in `routers/api/v1/swagger/`.
- Markup and rendering: `modules/markup/` adapters, sanitization, code highlighting, math blocks.
- Storage and database: `modules/storage/`, the XORM models under `models/`, migration helpers under `modelmigration/`.
- Frontend: Vue components under `frontend/web_src/js/components/`, fetcher wrappers in `frontend/web_src/js/modules/fetch.ts` and `frontend/web_src/js/modules/fetch-action.ts`, utilities under `frontend/web_src/js/utils/`.

If subagents are unavailable, simulate the same breadth yourself. Do not let the first good candidate stop the survey.

Start with the largest production-code deltas. A broad simplification audit that stops after obvious unused symbols can miss the files where duplicated lifecycle or defensive machinery carries most of the cost.

## Simplify prose with the code

Treat comments and documentation as maintained surface area. Apply [k7-prose-standard](../k7-prose-standard/SKILL.md) when a survey includes prose.

- Delete comments that restate code or explain behavior owned elsewhere; keep required local contracts per `AGENTS.md` (one short same-line note that explains the *why*).
- Keep docs at their owning level; omit implementation details and rare cases unless they change a maintained contract. `docs/*.md` follows [k7-doc](../k7-doc/SKILL.md); design docs and postmortems live in `docs/` only when no other surface owns the content.

## Audit trust and lifecycle boundaries

For every defensive copy, freeze, validator, and callback capture, name where the value came from and who owns it next. Same-process typed service calls ordinarily borrow readonly values; parsers, config loaders, queues, Git command output, durable files, workers, processes, and wire decoders own or validate their data. Tests built around hostile getters, fake typed objects, callback replacement, or mutation after a same-process handoff are evidence of a potentially speculative contract, not automatic justification for keeping it.

For complex asynchronous code, draw the ownership graph and map each sentinel, readiness promise, cancellation path, disposer, and state flag to a distinct owner or transition. When several mechanisms mirror the same liveness or settlement fact, propose one transaction or lifecycle controller instead. Preserve separate machinery where it protects synchronous publication and rollback, callback containment, first-terminal-outcome arbitration, worker/process ownership, or dispose-to-quiescence. `k7-ci-test-reliability` owns the test-side teardown rules; consult it before simplifying lifecycle code.

## Hand-rolled code versus a dependency

Introducing a dependency is a valid simplification move, not a policy exception. When surveying, ask of protocol parsers, framers, retry/backoff loops, glob matchers, diff engines, and similar infrastructure: does a well-maintained Go module, a maintained npm package under `frontend/package.json`, or a stdlib helper at the language floor already do this?

Prove a dependency-swap candidate like any other, plus:

- Read the hand-rolled implementation and name the exact surface the package covers; residual semantics the package does not cover count against the swap.
- Check the package's health honestly (maintenance, adoption, transitive footprint) and prefer stdlib when the language floor has the helper.
- Check whether `go.mod` or `frontend/package.json` already carries the dependency; the swap that re-adds a dropped module is a behavior change, not a simplification.
- Weigh net deletion: implementation plus dedicated tests plus docs, minus the glue that remains. A wrapper that relocates the same complexity is not a win.
- Per [docs/guidelines-backend.md](../../../docs/guidelines-backend.md) and [docs/guidelines-frontend.md](../../../docs/guidelines-frontend.md), every `go.mod` / `frontend/package.json` change must be justified in the PR description and must reference an existing upstream commit; `make tidy` runs after any `go.mod` edit.

## Prove or reject each candidate

For every symbol or behavior, classify consumers before writing:

- Production corpus: Go under `modules/`, `models/`, `services/`, `routers/`, `cmd/`; Go templates under `templates/`; TS/Vue under `frontend/web_src/`; runtime scripts and CLI entry points.
- Non-production corpus: tests under `tests/integration/`, `tests/e2e/`, package-local `*_test.go` files, READMEs/docs, design docs, generated expected outputs, and comments.
- Ambiguous corpus: examples and scripts that may be product smoke paths. Inspect usage before classifying.

Use `rg` first. Good searches include the exact symbol, event name, package name, config key, method name with both `.name(` and `name(`, swagger `summary:` strings, and any wire strings. Then read the call sites, public interfaces, dynamic event names, tests, docs, and loader/config paths.

Windows (PowerShell) variant of the symbol search:

```powershell
Select-String -Path modules\util\*.go -Pattern '\bCopyFile\b'
```

Reject or downgrade a candidate when:

- A production caller exists and the simplification would be a feature decision rather than a cleanup.
- The API is explicitly justified by an existing design doc, a hard-won defensive pattern, or a documented public contract, and the new evidence does not beat that reason.
- The removal would force unrelated churn without actually reducing the public API or required behavior.
- The idea is correct but tiny. Add a targeted `TODO(name):` instead, using the inline-marker discipline below.

## Coalesce superseded design docs and READMEs

Audit `docs/*.md` and any package-level `README.md` when the user asks to reduce or coalesce documentation, or when the simplification being implemented makes an owning doc obsolete. Do not expand every code-simplification survey into a repository-wide doc audit.

Follow the one-owner-per-fact rule from [k7-doc](../k7-doc/SKILL.md): identify the current owner from shipped code, configuration, generated catalogs, swagger annotations, and newer docs; consolidate the unique rationale into the current owner; repair every inbound link; then delete the obsolete doc. Dates and titles are discovery hints, not proof.

For each candidate chain:

1. Identify the current owner from shipped code, configuration, generated catalogs, package docs, newer design docs, and inbound links.
2. Classify the old doc as fully or partially superseded. Any surviving behavior, current contract, durable format, compatibility obligation, or independently current alternative makes it partial. Rationale that can be transferred to the current owner does not by itself make supersession partial.
3. For full supersession, move every unique rationale, alternative, consequence, shipped verification evidence, and named coverage gap into the current owner. An inventory that only describes deleted implementation mechanics is not one of those decision facts.
4. Repair every inbound link, then delete the obsolete doc.
5. Search exact filenames, symbols, config keys, event names, and wire strings after the edit. Keep partial supersessions cross-linked and current.

An added-then-removed feature is a common full-supersession case. Let the removal note own the history only when the feature is absent from production code, configuration, schemas, durable or wire formats, migration, and compatibility behavior; no current documentation presents it as available; and no test exercises it as supported behavior. Removal rationale and tests that enforce absence may remain. Preserve why the feature originally existed, why that motivation no longer justified it, alternatives to full removal, the capability given up, conditions for reintroduction, and evidence that removal is complete. Old tests and implementation mechanics that verified only the deleted behavior are not current verification evidence.

Reject consolidation when the removal is only one transport, default, implementation, or presentation of a feature; when persisted data or compatibility handling survives; or when the removal note does not yet carry enough rationale to prevent accidental reintroduction. A current negative design decision may legitimately need its own note even though the removed implementation is gone.

## Write the tracking artifact

Gitea does not maintain an internal Agent Note tree; durable simplifications land in one of three places:

- A `TODO(name):` marker at the call site, with a one-line reason that names the smell and the action that would resolve it. `AGENTS.md` calls for almost no comments; a `TODO(name):` is the sanctioned exception when the smell is real but the fix is local. Examples: `TODO(double-default): collapse zero-value defaults once the JSON loader is gone.`, `TODO(unused-export): DeleteBranchProtection was kept for v1 API parity; see <full URL>.`.
- A tracking issue in the gitea repo, referenced from the call site by full URL. Use the full URL form per `AGENTS.md`; cite the URL once at the call site, then re-cite by short name in the PR description and the design doc.
- A design doc under `docs/` for cross-cutting simplifications that need a recorded rationale, alternatives, consequences, and verification evidence. Follow [k7-doc](../k7-doc/SKILL.md) for placement and the [k7-prose-standard](../k7-prose-standard/SKILL.md) rules for sentence-level coverage.

Prefer the smallest artifact the candidate needs. A local rename lives at the call site; a cross-package seam change lives in `docs/`; an in-flight migration may need both.

Prefer this structure for a `docs/` design doc, adjusting when the idea needs it:

- `# <action-oriented title>`
- `## Problem`: name the current API, cite the relevant files, and state the consumer evidence. Separate production callers from tests/docs.
- `## Proposal`: say exactly what to remove, fold, demote, or rehome. Include tests, docs, READMEs, swagger annotations, and template partial cleanup when relevant.
- `## Why not keep it?` or `## What we give up`: make the strongest counterargument legible.
- `## Acceptance criteria`: observable end state and gates.
- `## Risks`: public API changes, behavior changes, future product wants, and why the tradeoff is still reasonable.

Be concrete enough that an implementing PR can follow the trail. Avoid vague "simplify this package" docs. When a proposal overlaps an existing design doc, consolidate the useful details into the existing one rather than creating a duplicate.

## Inline TODO notes

Use inline `TODO(name):` / `FIXME(name):` / `XXX(name):` only for small, local cleanups that are clearly useful but not durable design decisions. Keep them short and actionable:

- Name the smell with a stable tag, e.g. `TODO(double-default)` or `XXX(unused-export)`.
- Explain why it is safe to revisit and what action would simplify it (one short clause, same-line preferred).
- Do not add TODOs for speculative complaints or for behavior that needs a design-doc-level decision.

A `TODO(name):` with a real reason is preferred over a sprawling audit PR that stops at "looks complex".

## When folding another PR or branch

Diff the sibling branch against `origin/main`, not against the current PR branch, so you see its independent contribution. For each item:

- Port non-overlapping `TODO(name):` markers, design docs, and tracking issues that meet the quality bar.
- Consolidate overlapping material into the existing design doc or issue that owns the topic.
- Do not port duplicate or lower-confidence proposals just to preserve the count.
- Update the PR body so reviewers see the true candidate count and scope. Per `AGENTS.md`, PR descriptions stay under 1000 characters; link to the design doc for the rationale.
- Close the duplicate PR only when the user asked you to, or when you clearly own that housekeeping.

## Validation and PR hygiene

Select the smallest focused checks that cover the outgoing diff; [../k7-pre-push-checks/SKILL.md](../k7-pre-push-checks/SKILL.md) owns the selection rules. After the relevant checks pass:

```bash
git diff --check
```

For docs-only work, also run:

```bash
make lint-md
make lint-spell
```

For code-comment work, run the relevant validator when one exists. For API changes, regenerate and validate swagger:

```bash
make generate-swagger
make swagger-validate
```

Do not manually repeat a passing check merely because commit or push follows.

When opening or updating a PR, summarize:

- How many `TODO(name):` markers, design docs, and tracking issues were added, consolidated, retained as partial supersessions, or deleted.
- The main areas surveyed.
- What was intentionally excluded.
- Which checks passed.

For each consolidation group, name the old and current owners, state the evidence for full supersession, and explain why deletion is safe. If an added-then-removed scan finds no qualifying doc, report that result and the representative partial cases retained.

Use a draft PR while the survey is still expanding; mark ready only when the candidate set, review responses, and validation are settled.