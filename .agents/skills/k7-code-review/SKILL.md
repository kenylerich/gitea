---
name: k7-code-review
description: Use when reviewing a pull request in the gitea repo — orients the reviewer to AGENTS.md conventions, Go and TS/Vue quality bars, and the review-specific checks that code alone can't show
---

# Reviewing a Gitea PR

**This skill is guidance, not a checklist.** Verify and fetch the PR's live base and exact head, then inspect the diff against that base before reading enough surrounding code to understand the design. Re-establish the base and re-run the inspection after a retarget, a force-push, or a merge. Prioritise correctness, lifecycle, security, and broken required behaviour over style; a short review with one substantiated blocker is better than a list of nits. Review all commits in the PR, not only the latest.

## Sources of truth

- [AGENTS.md](../../../AGENTS.md) and [CONTRIBUTING.md](../../../CONTRIBUTING.md): standing authoring and contribution rules.
- [docs/development.md](../../../docs/development.md), [docs/build-setup.md](../../../docs/build-setup.md), [docs/testing.md](../../../docs/testing.md): build, host, and test workflow.
- [docs/guidelines-backend.md](../../../docs/guidelines-backend.md): package layout, `db.WithTx`, XORM, API and swagger rules.
- [docs/guidelines-frontend.md](../../../docs/guidelines-frontend.md): TS, Vue, Tailwind `tw-*`, DOM and data-fetch helpers.
- [docs/guidelines-refactoring.md](../../../docs/guidelines-refactoring.md) and [docs/community-governance.md](../../../docs/community-governance.md): how refactors and reviews land.
- [k7-prose-standard](../k7-prose-standard/SKILL.md): coverage and editorial judgment for comments, docs, prompts, and visible strings.
- [k7-pre-push-checks](../k7-pre-push-checks/SKILL.md): how to pick the smallest local checks that cover the outgoing diff.
- [k7-ci-test-reliability](../k7-ci-test-reliability/SKILL.md): isolation, flake risk, and deterministic waits.

## Establish the live base and head

Fetch the PR's base and head before reading the diff. Trust nothing earlier in the thread.

```bash
gh pr view <pr> --json number,baseRefName,headRefOid,headRefName,isCrossRepository,state
git fetch origin <base> <head>
git diff --stat origin/<base>...origin/<head>
git diff --name-only origin/<base>...origin/<head>
```

Re-run after a retarget, a force-push, or a merge. Review all commits in the PR, not only the latest.

## Blocking requirements

1. **New prose receives semantic review.** Use [k7-prose-standard](../k7-prose-standard/SKILL.md) on every added or changed Markdown, Go/template comment, prompt, description, diagnostic, and visible string; lint does not establish coverage or accuracy.
2. **Docs, config, and locale match the code.** Behaviour, defaults, errors, wire fields, and visible strings update the relevant `docs/` page, `custom/conf/app.example.ini`, swagger comments, and `options/locale/locale_en-US.json` in the same diff. Only edit `locale_en-US.json`; other locales sync automatically. Comments state non-obvious contracts in one short same-line note; flag implementation narration, review history, and duplicated rationale.
3. **Persisted model changes ship a migration.** A field added to a `models/` struct that is written to the database lands with a matching migration under `modelmigration/` and a passing `make test-migration` run (see [docs/testing.md](../../../docs/testing.md)).
4. **API changes regenerate swagger.** Every change under `routers/api/v1/` or `modules/structs/` runs `make generate-swagger`; the PR includes the regenerated `templates/swagger/v1_json.json` and the matching swagger option registration.
5. **Required evidence exists.** The author ran the checks selected by [k7-pre-push-checks](../k7-pre-push-checks/SKILL.md) for the diff and CI covers the rest. Review the semantic gaps neither can detect.
6. **Authorship and metadata.** New `.go` files carry the current year in the header. Commits use Conventional Commits plus the `enhance` type for user-facing enhancements, with one `Assisted-by: AGENT_NAME:MODEL_VERSION` trailer and never `Co-Authored-By` or `Signed-off-by`. PR descriptions stay under 1000 characters: only what and why, with before/after screenshots for UI changes. Reference issues and PRs by full URL, not by number. Attribute agent authorship on one trailing line in comments, never as a PR description section.

## Manual checks

- **Intent and interface contracts.** Trace both sides of every changed interface: handler/service/model, route/middleware, Go API/TS client, Go template/frontend helper. Confirm implementation matches the PR and surrounding callers, including errors, context, ownership, and disposal.
- **Go correctness and lifecycle.** Trace `context.Context` propagation, every `db.WithTx` boundary, error wrapping (`fmt.Errorf("...: %w", err)`), goroutine lifetimes, `defer` ordering, and resource close/cleanup. Flag missing cancellation, lost errors, and shared mutable state without synchronisation. Permission and access checks must run on every web and API route.
- **Security surfaces.** SQL: every dynamic parameter is bound, never concatenated. XSS: user-controlled values rendered by Go templates are auto-escaped or routed through a helper that is; nothing reaches `innerHTML` in TS or `v-html` in Vue. CSRF, auth scope, and rate limiting on new routes match the surrounding pattern.
- **Frontend conventions.** New code is TS or Vue SFC under `frontend/web_src/`. Templates prefer Tailwind `tw-*` utilities (with `flex-*` helpers and `gt-*`/`g-*` only where `tw-*` does not exist) over inline `style`. TS uses `!` instead of `?.`/`??` when a value always exists; `import type` for type-only imports. Vue event listeners are not `async` unless they call `e.preventDefault()` before the first `await`.
- **i18n ownership.** New product text is a stable locale key, not inline JSX, template literal, accessibility attribute, or primitive default. Missing keys are a blocker.
- **Test strength.** Assertions fail for the intended regression and verify external state, logs, events, or disposal rather than restating the implementation. Coverage is necessary but not evidence the scenario is correct. Prefer unit tests where logic is testable in isolation; integration tests stay under 2s and e2e tests under 4s. Wait on deterministic conditions, never `sleep`; prefer semantic locators in e2e.
- **Test reliability.** For resource-owning, asynchronous, or platform-sensitive tests, apply [k7-ci-test-reliability](../k7-ci-test-reliability/SKILL.md) to worker/job topology, database isolation, global-state restoration, synchronisation, timeout budget, and quiescent teardown.
- **Configuration choices.** Each new option, default, or public flag has current-consumer evidence or prior art. New options land in `custom/conf/app.example.ini` with a comment describing the trade-off.
- **Enforcement paths.** Follow every denial path to the operation that executes it; exercise direct and alternate callers that can bypass middleware, schema validation, or route-level permission checks.
- **Snapshot or output changes.** Visible user-facing changes (HTML, rendered Markdown, e2e screenshots) ship updated snapshots or explain why no snapshot applies. Review expected-output diffs as behaviour changes, not formatting noise.

## Reporting findings

State the defect, location, impact, and evidence. Place a localised defect inline on the tightest relevant diff range; use a PR-level comment for cross-cutting architecture, scope, or review-wide synthesis. Separate blockers from suggestions and omit issues already enforced by a green gate. Use the existing GitHub review thread for replies. When receiving review, verify each claim and fix or rebut it on technical grounds without performative agreement.
