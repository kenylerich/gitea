---
name: k7-prose-standard
description: "Use when writing, reviewing, restoring, trimming, or auditing prose in the gitea repo, including deciding where documentation or comments are required across Go doc comments, Go error and log messages, same-line code comments, Go template and Vue SFC comments, CSS comments, options/locale/locale_en-US.json UI strings, go-swagger annotations, CLI help text under cmd/, commit messages, PR descriptions, and docs/*.md."
---

# Gitea prose standard

Write enough to preserve the contract, then remove reasoning transcripts, repetition, and decoration. A contract is an obligation, invariant, precondition, postcondition, or compatibility promise that a caller, callee, implementer, producer, or consumer relies on. This skill owns editorial judgment and required prose coverage; use [k7-doc](../k7-doc/SKILL.md) for placement and one-owner-per-fact, and [k7-trim-cot-leakage](../k7-trim-cot-leakage/SKILL.md) for hunting and fixing reasoning-transcript leakage. It is guidance, not a script.

`AGENTS.md` is the authority for comment discipline and authorship metadata; this skill restates the rule, does not override it. When a rule conflicts, `AGENTS.md` wins.

## Inputs and exclusions

Require an explicit `scope`. If it is missing, report the required input and stop; do not infer a repository-wide scope or begin an interview.

Accept `mode: automatic | interactive`; default to `automatic`. Enter interactive mode only when the user explicitly requests questions or calibration.

`mode` controls questions, not write authority. Review and audit tasks report findings without editing; explicitly requested write, fix, or trim tasks apply clear changes.

Always exclude `vendor/` and `frontend/web_src/js/vendor/` from discovery, review, and edits, even when the requested scope is the whole repository. Do not follow a symlink into them. Put exclusions after inclusion globs so a later include cannot re-admit them: for example, end ripgrep commands with `--glob '!vendor/**' --glob '!frontend/web_src/js/vendor/**'`, and give Git commands an explicit `:(exclude)vendor/**` pathspec. If the requested scope contains only excluded paths, report that no eligible files remain.

Treat generated artifacts as derivative. The Swagger spec (`templates/swagger/v1-swagger.generated.json` and `templates/swagger/v1-openapi3.generated.json`), compiled CSS in `public/assets/`, and the embedded bindata (`options/bindata.go`) trace back to a source: a go-swagger comment on the route, a CSS source under `frontend/web_src/css/`, and a `make` target. Edit the owning source first, then regenerate the derivative with `make generate-swagger`, `make svg`, `make go-licenses`, or the appropriate `make generate-*` target. When a generator extracts a summary from owner prose, make the extracted sentence complete for that surface.

## Preserve the complete proposition

Before editing, identify every proposition in the passage. Preserve each relevant:

- actor and action;
- condition, timing, and ordering;
- modality such as must, may, or never;
- negative guarantee and exception;
- ownership, side effect, failure mode, and consequence.

Remove adjectives, repetition, and narration only when every factual clause survives and the result is clearer. A smaller word count alone is not an improvement.

Keep a complete local contract at the point of use: behavior, failure, ownership, and consequence that a caller or maintainer needs there. Aggressively link to the owning document for architecture, rationale, algorithms, history, or extended examples. One explanation has one home; essential contract facts may repeat locally.

Keep non-obvious rationale when omitting it could plausibly cause misuse or an incorrect simplification. Otherwise state the consequence and link the rationale home.

## Required coverage by prose location

This is not a one-way shortening pass. Add or restore prose when code, types, and structure do not communicate a required contract below. Do not add a comment when those facts are already obvious locally.

- **Go doc comments on exported identifiers:** state the behavior callers rely on. For `func` and `method`, the first sentence is the summary; subsequent sentences document returned values, error sentinels, ownership of returned resources, ordering, and side effects. For `type` and `const` blocks, document the value the identifier represents, not its declaration shape. Follow [Effective Go's "Comment Sentences"](https://go.dev/doc/effective_go#comment-sentences); run `make lint-go` to catch godoc and golangci-lint issues.
- **Go error and log messages:** name the failing subject (file path, repo, user, key, request) and the violated condition. `fmt.Errorf("...: %w", err)` wraps; the prefix states the caller's failed step. `log.Error("[scope] ...", key, err)` uses a `[scope]` tag the operator can grep. Avoid "should" and "failed to"; state the rule that was violated and the consequence.
- **Same-line code comments:** the default form. One short clause on the same line as the statement, explaining a non-obvious why (a workaround, an ordering invariant, a security boundary, an upstream bug). Do not narrate the line below.
- **Go template comments (`{{/* ... */}}` at the top of a `*.tmpl`):** document the `Template Attributes` block: what each attribute means, whether it is required, and what happens when it is empty or absent. The current attribute list at `templates/shared/combomarkdowneditor.tmpl:1` is the shape to follow.
- **Vue SFC comments:** one short same-line or adjacent comment that names a non-obvious prop contract, an unmounted-state invariant, or a key accessibility rule. Avoid `// TODO`-style placeholders in `.vue` files; either fix the smell or leave a `TODO(name):` next to the call site with a reason.
- **TS single-line `/** ... */` comments:** state the caller-visible contract of the exported function or constant in one sentence. Gitea uses single-line JSDoc-style comments for TS (see `frontend/web_src/js/utils/url.ts:1`); do not introduce multi-line JSDoc blocks unless the function needs a `@param`/`@returns` contract.
- **CSS comments:** one short note explaining a non-obvious selector, a fallback chain, or a known framework conflict (Fomantic-UI and Tailwind both ship utilities). The header in `frontend/web_src/css/helpers.css:1` shows the house style.
- **`options/locale/locale_en-US.json` UI strings:** wording that reaches a user is behavior. Per `AGENTS.md`, edit only `locale_en-US.json`; other locales sync from Crowdin and must not be hand-edited. Treat the JSON value as a stable locale key whose text a translator may paraphrase, but whose interpolation placeholders (`%s`, `%d`) and HTML are not negotiable. Update the key when the user-visible concept changes; do not rename a key just to shorten it.
- **go-swagger annotations:** every API route under `routers/api/v1/` ships a `// swagger:operation ...` block (see `backend/routers/api/v1/repo/issue_pin.go:17`). Document the `summary`, every `parameter` (`name`, `in`, `type`, `required`), every response (`status`, `description`), and every failure case the caller can hit. The generated `templates/swagger/v1-swagger.generated.json` (and `templates/swagger/v1-openapi3.generated.json`) is what consumers read; missing fields are silent failures. After editing, run `make generate-swagger` and `make swagger-validate`.
- **CLI help text under `cmd/` (currently being moved to `backend/cmd/` per `docs/migration-repo-restructure.md`):** `urfave/cli/v3` `Usage:` and `Description:` strings. Name the action, the inputs it requires, the side effects, and the exit-on-error conditions. See `backend/cmd/admin.go:23` for the house shape.
- **`custom/conf/app.example.ini`:** every new `app.ini` option lands here with a comment that names the default, the trade-off, and the version that introduced it.
- **Commit messages and PR descriptions:** Conventional Commits (see `AGENTS.md`) plus Gitea's `enhance` type for user-facing enhancements. PR description is what and why only, under 1000 characters, with before/after screenshots for UI changes. Reference issues and PRs by full URL, not by number. The author adds one `Assisted-by: AGENT_NAME:MODEL_VERSION` trailer and never `Co-Authored-By` or `Signed-off-by`. Agent authorship on PRs and issues is one trailing line in a comment, not a PR description section.
- **`docs/*.md`:** match the house style documented at [k7-doc](../k7-doc/SKILL.md). The first paragraph names the subject, the reader, and the next useful neighbor.
- **Tests:** explain only non-obvious test design — why a fixture, assertion, real entry path, or platform accommodation is necessary. Delete walkthroughs and step-by-step narration.

Preserve searchable mechanism names and meaningful modal, temporal, or negative emphasis. Normalize decorative emphasis only.

## Workflow

1. Confirm the scope, mode, current branch or PR base, and the rules in `AGENTS.md`. Do not inspect unrelated branches.
2. Read [k7-doc](../k7-doc/SKILL.md) and the owning code or document before judging a passage. For calibration or unfamiliar cases, read [the distilled examples](references/examples.md).
3. Inspect the requested scope, not only the largest files. Use searches and word counts to find candidates, then judge passages semantically.
4. Classify each candidate as keep, add, trim, restore, restructure, or defer. Apply clear changes only when the task authorizes edits; do not manufacture edits to satisfy a deletion target.
5. Update the owner before derivative artifacts. Re-check analogous passages after learning a new rule.
6. Run the narrow relevant checks, `git diff --check`, and behavior tests for visible strings. After touching swagger annotations: `make generate-swagger` and `make swagger-validate`. After touching `.go`: `make fmt` then `make lint-go` for what changed. After touching `docs/*.md`: `make lint-md` and `make lint-spell`. After touching `options/locale/locale_en-US.json`: no extra check (other locales sync from Crowdin). Verify the final diff contains no `vendor/` or `frontend/web_src/js/vendor/` path and report any accidental match rather than claiming a clean exclusion.
7. Report the inspected scope, clear changes, deliberate keeps, deferred cases, and checks actually run.

## Borderline decisions

A case is borderline only when at least two versions satisfy the complete-proposition rule but trade accepted principles, and this skill does not already resolve the tradeoff. A rewrite with one proposition-preserving answer is not borderline.

In automatic mode, apply clear edits when authorized and report genuine borderline cases without asking questions. Do not weaken a proposition to make progress.

In interactive mode, group analogous passages under the governing principle. Present two or three viable versions, recommend one, and state the factual or structural difference. Do not offer inferior distractors. Use the user's requested channel; when calibrating a PR through inline comments, place the recommended provisional version in the diff and attach the alternatives to that exact line.

After the user decides, distill the principle and versions into [the examples](references/examples.md), without PR history or reviewer narration, and apply the learned rule to every analogous passage in scope.