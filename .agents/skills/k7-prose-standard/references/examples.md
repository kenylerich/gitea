# Distilled prose examples

Use these examples to identify the governing principle, not as text templates. "Balanced" preserves every load-bearing proposition with the least explanation needed at that location. Every example targets a gitea surface (Go, Go template, Vue, TS, swagger, locale JSON, Markdown).

## Preserve every factual clause

**Original:** "The migration coordinator carefully serializes writes per session, flushes buffered events before disposal resolves, and reports backend failures to the caller."

**Over-trimmed:** "The coordinator serializes persistence."

**Balanced:** "The coordinator serializes writes per session, flushes buffered events before disposal resolves, and reports backend failures to the caller."

Remove decoration and repetition, not propositions. Actor, per-session scope, disposal ordering, and failure visibility are separate facts.

## Explicit skill scope is functional

**Over-trimmed:** "Read the sources and use judgment."

**Balanced:** "This skill is guidance, not a complete checklist. Use judgment beyond the named checks; documented requirements still apply."

**Over-detailed:** Several paragraphs defending why lists cannot replace independent reasoning.

Keep the explicit limitation because it changes how an agent applies the workflow. Trim repeated persuasion, not the guardrail.

## Go doc comments preserve the returned resource contract

**Over-trimmed:**

```go
// CopyFile copies file from source to target path.
func CopyFile(src, dest string) error {
```

**Balanced:**

```go
// CopyFile copies src to dest, preserving the source's modification time.
// It returns an error wrapping os.PathError if either path is missing or
// the destination cannot be created; the caller owns dest on a non-nil error.
func CopyFile(src, dest string) error {
```

**Over-detailed:** A walkthrough of every `os.Open`/`defer`/`os.Chtimes` branch already visible in the body.

Behavior, side effect on `dest` on error, and wrapping convention are caller-visible contract facts.

## Same-line Go comment explains a why

**Over-trimmed:** No comment on a `defer rows.Close()` call.

**Balanced:** Same-line comment that names the non-obvious resource being released or the invariant being upheld:

```go
defer rows.Close() // driver holds the connection until Rows is closed
```

**Over-detailed:** A paragraph above the function explaining Go's database/sql contract.

`AGENTS.md` says: write almost none, short and preferably same-line, explaining why. If the code already shows the why (a literal `return nil, err` or an obvious defer), omit the comment.

## Go error string names the failing subject

**Over-trimmed:** `return errors.New("invalid input")`

**Balanced:** `return fmt.Errorf("invalid time string %q: expected RFC3339", timeStr)`

**Over-detailed:** A multi-line error chain narrating each parser decision.

The subject (`timeStr`), the format the caller should use, and the wrapped `error` are separate facts. The wrapped error is the diagnostic; the prefix is the contract.

## Go log message carries a scope tag

**Over-trimmed:** `log.Error("Error reading file: %v", err)`

**Balanced:** `log.Error("Error reading file for %s: %v", envKey, err)`

**Over-detailed:** `log.Error("an unexpected error occurred while attempting to read the file referenced by env key %s; the underlying error was: %v", envKey, err)`

The `[scope]` tag and the failing key are what the operator greps for; the wrapped error carries the diagnostic. Compare the house style at `backend/modules/setting/config_env.go:132`.

## Go template attribute block documents required inputs

**Over-trimmed:**

```tmpl
{{define "form"}}
  <input name="{{.Name}}">
{{end}}
```

**Balanced:**

```tmpl
{{/*
Template Attributes:
* Name: name attribute for the input (required)
* Label: visible label; falls back to Name when empty
* Value: initial input value
*/}}
{{define "form"}}
  <label>{{or .Label .Name}}</label>
  <input name="{{.Name}}" value="{{.Value}}">
{{end}}
```

**Over-detailed:** A comment per attribute repeating what the next line already shows.

Document the contract: which attributes are required, what the empty-value behavior is, and where the value is rendered. See the real shape at `templates/shared/combomarkdowneditor.tmpl:1`.

## Vue SFC prop comment names a non-obvious contract

**Over-trimmed:** No comment above `defineProps<{...}>()`.

**Balanced:**

```vue
<script setup lang="ts">
// `tooltipUnit` is interpolated into the cell's aria-label; never pass user-controlled text.
const props = defineProps<{values: HeatmapValue[]; tooltipUnit: string}>();
</script>
```

**Over-detailed:** A `<!-- ... -->` block above every prop explaining its TypeScript type.

Comments state the *why* of a contract — what would happen if a caller broke it — not the type the type system already enforces.

## TS single-line JSDoc-style comment on an exported helper

**Over-trimmed:** No comment above an exported `globCompile` helper.

**Balanced:**

```ts
/** Match paths against the given pattern; returns null when the pattern is invalid. */
export function globCompile(pattern: string, sep = '/'): CompiledGlob | null {
```

**Over-detailed:** A multi-line JSDoc with `@param`, `@returns`, and an `@example` block.

Gitea uses single-line `/** ... */` comments for TS (see `frontend/web_src/js/utils/url.ts:1`). One sentence on the caller-visible contract is enough.

## CSS comment flags a framework conflict

**Over-trimmed:** No comment on a custom utility class.

**Balanced:**

```css
/* matches Fomantic-UI's `.interact-fg` so third-party themes keep both */
.g-fg { color: var(--color-primary); }
```

**Over-detailed:** A comment above each declaration restating the selector.

Comments name the non-obvious: a framework override, a fallback chain, or a deliberate duplication. See `frontend/web_src/css/helpers.css:1` for the house style.

## UI string in `locale_en-US.json` preserves interpolation

**Over-trimmed:** Renaming the key to fit a new label without checking callers, or rewriting the value to drop `%s`.

**Balanced:** Keep the key, update the value, and keep every `%s` / `%d` placeholder the template interpolates. The Vue template binds both:

```json
"sign_in_with_provider": "Sign in with %s"
```

```vue
<span>{{ locale.signInWithProvider }}</span>
```

If the template is `${name}`, the value must keep `%s`. The interpolation contract belongs to the template; the translator rewrites the surrounding text.

## go-swagger annotation documents every response

**Over-trimmed:** `// summary: Pin an Issue` with no parameters, no responses, and no failure cases.

**Balanced:**

```go
// swagger:operation POST /repos/{owner}/{repo}/issues/{index}/pin issue pinIssue
// ---
// summary: Pin an Issue
// parameters:
// - name: owner
//   in: path
//   type: string
//   required: true
// responses:
//   "204": {}
//   "403": {description: "The authenticated user lacks access to the repo"}
//   "404": {description: "The issue does not exist or is not in this repo"}
```

**Over-detailed:** Repeating the route path and HTTP verb in the description, or pasting the JSON schema into the swagger comment.

`docs/guidelines-backend.md` requires every possible result to be documented. After editing, run `make generate-swagger` and `make swagger-validate`; the committed `templates/swagger/v1-swagger.generated.json` is the contract consumers read.

## CLI help text names the side effect

**Over-trimmed:**

```go
Usage: "Synchronize repository releases with tags",
```

**Balanced:**

```go
Usage:     "Synchronize repository releases with tags",
ArgsUsage: "[owner/repo]",
Description: `Recreates each release whose underlying tag moved since the last sync.
No-op when the tag matches the release's SHA; prints and skips on error.`,
```

The `Usage` line is the short hint `urfave/cli/v3` shows in `--help`; the side effect (recreation, not deletion) and the per-repo filter belong in `Description`. Compare the existing shape at `backend/cmd/admin.go:23`.

## `custom/conf/app.example.ini` comment explains the trade-off

**Over-trimmed:** `[foo]\nBAR = true`

**Balanced:**

```ini
; When enabled, the action runner polls Gitea Actions workflows over SSH
; instead of HTTPS. Reduces certificate-management overhead but requires
; the runner host to trust the Gitea instance's host key.
[actions]
RUNNER_POLL_OVER_SSH = false
```

Defaults, the trade-off, and the version that introduced the option are the three facts a new operator needs.

## Delete reasoning transcripts entirely

**Over-detailed:** "First the loop checks whether the value is absent. If it is absent, the next branch returns early. Otherwise it continues, which is why the final assertion is safe."

**Balanced:** No comment when the code already expresses those branches. If the early return protects a non-obvious invariant, state only that invariant.

Do not compress a reasoning transcript into shorter narration; remove it.

## Configuration comments explain what the tree cannot

**Over-detailed:** "This entry loads the local filesystem provider, followed by the policy plugin, followed by the read, write, and edit tools," when the adjacent entries already show that order.

**Balanced:** "Load policy before the model-facing tools so their write and edit calls pass through the read-before-mutation gate."

Keep the consequence of order, a surprising scope rule, or a security boundary. Let the configuration show its own inventory.

## Do not trim for word count alone

**Current:** "The adapter converts provider errors into the shared error type so callers can handle authentication, rate-limit, and transient failures uniformly."

**Shorter but worse:** "The adapter normalizes provider errors."

**Balanced decision:** Keep the current sentence unless a link or surrounding contract already lists the failure categories. The shorter version loses the consequence and distinctions without improving structure.

## Model-visible text follows ownership

**Over-trimmed:** "The tool returns errors when a call fails."

**Over-detailed:** Copying another package's schema and renderer strings into this backend's README.

**Balanced:** Quote stable prompt, result, and error text owned by this package. Link the generated tool catalog for schemas and the consumer README for text another package owns; state only this package's conditions or deltas locally.

Wording that reaches a user is behavior, but duplication still drifts. Exactness belongs at the owner. Per `AGENTS.md`, edit only `options/locale/locale_en-US.json`; other locales sync automatically.

## PR description keeps "what and why" under 1000 characters

**Over-trimmed:** "Fixed bug."

**Over-detailed:** A bullet list of every file touched and every check that ran.

**Balanced:** Two to four sentences naming the user-visible behavior that changed, the why, and a `Fixes <full URL>` or `Refs <full URL>`. Screenshots before/after for UI changes. Reference issues and PRs by full URL, not by number. The full rule lives in `AGENTS.md`.

## Limitations are contracts, not debt inventories

**Over-trimmed:** Omitting a process-lifetime cache that makes configuration changes require restart.

**Over-detailed:** Listing private helper cleanup and unused test-only accessors with no caller or maintainer consequence.

**Balanced:** "Provider selection is cached for the process lifetime; adding a new provider requires restart." Keep ordinary cleanup in its `TODO` or in a linked issue.

Retain gaps and non-obvious constraints that affect use or safe maintenance. A README is not a backlog dump.