---
name: k7-ci-test-reliability
description: Use when designing, reviewing, or diagnosing gitea tests (Go unit, `tests/integration`, `tests/e2e`) that can fail nondeterministically under CI concurrency, the shared test database, process-global `setting.*`, spawned git/ssh subprocesses, real clocks, file-locking differences between Linux CI and Windows hosts, or Playwright animation/timing races; use k7-pre-push-checks separately to pick the outgoing command.
---

# Reliable Gitea CI tests

Build gitea tests that remain correct under the repository's real CI topology, not only when run alone on a quiet workstation. CI runs Linux; contributors develop on Windows or macOS, so platform and path handling differences are part of every test's surface. This skill owns isolation and reliability decisions; it does not select every command for a push.

## Read the owning rules

- Use [docs/testing.md](../../../docs/testing.md) for the test-tier policy and the available `make` targets.
- Use [AGENTS.md](../../../AGENTS.md) for style, scope, and the lint-only-what-changed rule.
- Use [docs/guidelines-backend.md](../../../docs/guidelines-backend.md) for Go patterns that affect tests (fixtures, migrations, DB schema).
- Use [docs/guidelines-frontend.md](../../../docs/guidelines-frontend.md) for TS/Vue patterns that affect tests (locators, async helpers).
- Read `tests/test_utils.go` and `tests/integration/git_helper_for_declarative_test.go` when the change touches the integration harness (`onGiteaRun`, `PrepareTestEnv`, `unittest.LoadFixtures`, `unittest.ResetTestDatabase`).
- After the test design is sound, use [k7-pre-push-checks](../k7-pre-push-checks/SKILL.md) to select outgoing validation.

## Model the execution topology

Assume these layers can overlap unless the active configuration proves otherwise:

1. Tests within one `*_test.go` file.
2. Separate Go packages or files in one `go test` process.
3. Independent `go test` invocations or `make` runs on one runner.
4. Multiple Actions jobs on one hosted runner.
5. Different runner hosts in a matrix (Linux primary; Windows when applicable).

Process isolation does not isolate the shared SQLite test database, fixture directories, predictable temp paths, the in-process `onGiteaRun` HTTP listener, the SSH listener started by `setting.SSH`, or inherited child processes. For every acquired resource, identify its owner, allocation mechanism, observable readiness signal, registered cleanup, and quiescent completion signal.

Do not serialize the entire suite because one fixture lacks isolation. Narrow the exclusive scope or change the resource allocation first; `t.Parallel()` markers do not protect a host resource from another package, process, job, or runner.

## Allocate resources atomically

Use the resource owner's allocator instead of checking availability and claiming it later.

- The integration harness binds a `net.Listener` on `setting.AppURL`'s host via `onGiteaRun`; tests do not allocate their own HTTP port for the test server. A test that does need its own listener must use `net.Listen("tcp", "127.0.0.1:0")` and read the assigned address back.
- Create per-test temporary roots with `t.TempDir()`; do not reuse predictable shared paths under `tests/integration/`.
- Use `unittest.LoadFixtures()` between packages, and `unittest.ResetTestDatabase()` for the integration process, to keep the shared DB deterministic. Do not call them inside a `t.Parallel()` body.
- LFS fixtures live under `tests/gitea-lfs-meta/` and `PrepareLFSStorage` reloads them per test; do not assume a test can mutate a shared LFS object and have later tests see the change.
- Git repos for integration scenarios live under `tests/gitea-repositories-meta/`; `PrepareGitRepoDirectory` syncs them. Do not write directly into those paths.
- Keep stable fixture data separate from live resource allocation. A recorded URL can stay stable while the fixture maps its transport to an OS-assigned port.

Literal paths and URLs used only as parser inputs or expected values are not acquired resources. Do not rewrite them merely because they look fixed.

## Contain process-global state

Treat `setting.*`, the loaded fixtures, the shared SQLite connection, `os.Setenv` writes, locale and timezone via `time.Local`, module mocks, `graceful.GetManager()` state, and the registered `testlogger` as exclusive mutable resources. Integration tests mutate `setting.*` directly (see `tests/test_utils.go`); restore it.

When a test must mutate global state:

- capture whether the original value was absent or present;
- restore that exact state, not a default;
- register restoration with `t.Cleanup` immediately, before any other failure can leave it dirty;
- keep the mutation inside the smallest possible scope;
- keep an `afterEach` fallback only when the registered cleanup cannot run.

Do not mutate `setting.*` in tests that run under `t.Parallel()`; the harness reloads state from disk and parallel mutations race.

## Respect platform-owned semantics

CI runs Linux; hosts run Windows or macOS. Values the OS owns do not always come back the way the test wrote them.

- File paths: use `filepath.Join`, not hard-coded `/`. Path separators differ; the case-insensitive Windows filesystem collides distinct names that stay distinct on Linux. Accept either case via `strings.EqualFold` only when a recorded ID must survive a Windows round trip.
- Line endings: `git config core.autocrlf` and `core.eol` shape what `git` writes. A fixture that depends on a specific byte sequence (`\n`, `\r\n`) is a fixture that needs explicit normalization at load, not at write.
- File handles: Windows releases them asynchronously. A rename or delete that completes instantly on Linux needs a bounded retry sized to observed contention. `os.RemoveAll` followed by an immediate `os.MkdirAll` on the same path can fail on Windows; space the operations or use `t.TempDir()`.
- File locking: SQLite file locking differs between Linux (POSIX advisory) and Windows (mandatory ranges). A test that opens the SQLite file directly (rare) needs a path that is not the shared integration DB.
- Time: NTFS `mtime` is 100-nanosecond ticks, not fractional milliseconds. A test that round-trips a timestamp through the filesystem must take the expected value from a fresh read rather than from a remembered one.
- Env vars: Windows matches variable names case-insensitively. A test that seeds `HTTP_PROXY` and `http_proxy` separately holds one entry.
- Permissions and signals: Windows has no POSIX permission or signal semantics. A test that depends on them takes an explicit platform skip naming the reason, rather than weakening the assertion everywhere.
- Git LFS: integration tests assume `git-lfs` is installed and on `PATH` (see [docs/testing.md](../../../docs/testing.md)). A test that uses LFS objects must not assume CI seeds the binary; declare the dependency in the Makefile target, not in the test.

Prefer an observation that holds on every platform. When a case genuinely cannot, exclude it on that platform explicitly.

## Budget timeouts against the lane

The integration runner uses `GOTEST_FLAGS=-timeout 40m` per `make test-integration`. Per-test `time.AfterFunc` and `-timeout` overrides do not yield to that budget, so a value below the lane's budget lowers what CI already granted.

The e2e runner uses Playwright timeouts multiplied by `GITEA_TEST_E2E_TIMEOUT_FACTOR` (4 on CI, 1 locally). Setting the factor to mask a flake rather than removing the flake is masking, not restoring.

Where a timeout is the subject of the test, keep the outer wait far larger than the timeout under test. A case proving that a 20 ms deadline fires must not race the harness's own wait.

## Synchronize on state

A fixed sleep is not evidence that setup completed or cleanup settled.

- Wait for the explicit readiness signal the owner publishes: `tests.PrepareTestEnv` returns a defer function, `onGiteaRun` returns once the test server is listening, and `unittest.LoadFixtures` is synchronous.
- Use a barrier (channel, `sync.WaitGroup`, or a `testing.TB`-safe channel) when two operations must overlap to prove a race; repeating the test serially does not prove it.
- Use a timeout only to bound a wait, never as the condition that makes the assertion correct.
- When time itself is the subject, inject a fake clock and always restore real time.
- In e2e tests, prefer Playwright's locator-based auto-waiting assertions (`toBeVisible`, `toHaveText`) and `page.waitForResponse` / `page.waitForFunction` over `waitForTimeout`.

## Dispose to quiescence

Register cleanup with `t.Cleanup` immediately after acquisition so assertion failures also release the resource. Cleanup stops new callbacks, drains goroutines, closes spawned listeners (`s.Shutdown(ctx)` in `onGiteaRun` is the model), terminates the SSH listener started by `setting.SSH`, and awaits child processes.

`exec.Command(...).Start()` without `Wait()` is incomplete teardown. When late completion is possible, prove that disposal prevents it from mutating another test (e.g., a queued HTTP request or a buffered log write).

`graceful.GetManager().HammerContext()` provides a context that signals hard shutdown; use it for cleanup that must finish regardless of test timeout.

## Prove the intended regression

- Observe an ordinary regression fail before the fix when practical. For a new guard, temporarily introduce the rejected case and observe the intended failure.
- For a race, use barriers to prove overlap; repeated execution alone is not a race test. Combine the focused run with `go test -run '^TestName$' -count 20 ./tests/integration/` (or the package owning the test) to expose order/load-sensitive bugs.
- For ports, sockets, shared paths, subprocesses, or other host resources, run independent test processes concurrently when cross-process isolation is part of the fix.
- Where a fixture spawns with its own deadline, assert that no signal or timeout ended the child before asserting its exit status, so a killed child reports as a timeout instead of a status mismatch.
- Verify external state (events, files, logs, exits, disposal) instead of trusting the component's self-report.

Stress runs supplement a deterministic regression; they do not replace one.

## Reject flake-masking fixes

Do not present these as root-cause fixes for deterministic local tests:

- increasing a timeout without identifying the awaited state;
- adding retries or `t.Skip` after the failure;
- making all tests serial via global `t.Parallel()` removal;
- swallowing an error or unhandled rejection;
- weakening an assertion (`.Contains` instead of `.Equal`, `len(...) > 0` instead of `==`);
- normalizing away unstable behavior;
- adding a `sleep` before cleanup or assertion.

Retries remain valid only for documented external-provider tests under the live-API policy (rare in this repo). Keep that exception at the external boundary.

Restoring a budget is not masking. Raising a suite to the lane budget it already had, or sizing a bounded retry to the contention actually measured on the runner, names the awaited work and returns what the lane granted; neither invents headroom around an unexamined wait.

## Diagnose existing flakes

For an existing probabilistic CI failure, read [references/ci-flake-diagnosis.md](references/ci-flake-diagnosis.md). A diagnosis-only request remains read-only: report the cause and evidence unless the user also asks for a fix.

## Validate and report

Run the smallest focused regression for the affected behavior. Add topology-specific evidence only when the change owns that risk:

- global mutation needs restoration evidence;
- lifecycle or subprocess work needs quiescent-teardown evidence;
- ports, sockets, or shared paths need concurrent independent-process evidence;
- a new guard needs a negative control.

Before a push, use [k7-pre-push-checks](../k7-pre-push-checks/SKILL.md). Report exact commands and observed results; do not describe retries, skipped tests, or pending CI as passing.