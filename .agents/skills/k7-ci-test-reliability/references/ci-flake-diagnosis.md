# CI flake diagnosis

Use this workflow only when the task is to investigate an existing probabilistic test or CI failure in gitea. Preserve the requested read/write scope: diagnosis does not authorize a fix, workflow rerun, or CI configuration change.

## Freeze the evidence

Record the repository, workflow, job, commit SHA, runner labels (note whether the runner is `ubuntu-latest` or `windows-latest`), timestamps, exact failing test or command, and the first stable failure signature. Keep infrastructure messages separate from test output.

Compare multiple failing and passing runs. Prefer runs of the same SHA; when that is impossible, verify that the relevant test and CI configuration are identical across the compared commits. One passing rerun does not prove an infrastructure fault, and one timeout does not prove a product race.

Use Actions logs and metadata to establish whether failures overlap on one host or resource namespace. Preserve links to the supporting runs rather than pasting large logs.

## Classify the failure

Classify from recorded evidence, not from the eventual fix:

- **Shared-test-database collision:** two packages or processes write to the same SQLite file or fixture directory between `unittest.ResetTestDatabase()` / `unittest.LoadFixtures()` calls; rows left by an earlier test change the next test's baseline. The integration harness resets fixtures per process, not per test.
- **Incomplete lifecycle:** `t.Cleanup` returns before children, the in-process HTTP server (`s.Shutdown` in `onGiteaRun`), the SSH listener started by `setting.SSH`, or queued log writes reach quiescence; later output or mutations appear in another test.
- **Process-global contamination:** outcome depends on test order or leaked `setting.*`, env vars, locale, timezone (`time.Local`), module mocks, `graceful.GetManager()` state, or fixtures reloaded by a parallel test.
- **Load-sensitive synchronization:** a `time.Sleep`, polling interval, or assumed event-loop turn substitutes for observable readiness or completion. In e2e, this shows up as `waitForTimeout` racing against animations or async network calls.
- **Platform or entry-path mismatch:** the failure consistently follows an OS, shell, filesystem rule, source/build mode, or executable entry. Path separators (`\` vs `/`), line endings, file locking, NTFS mtime precision, env var case, and file handle release timing all differ between Windows and POSIX, so a case passing on macOS says nothing about the Linux CI lane. Git LFS presence on `PATH` is part of the same category when integration fixtures need it.
- **Product concurrency defect:** the test controls its resources, reproduces deterministically with explicit overlap, and exposes a race in shipped behavior.
- **External-provider transience:** the failure is owned by a live API or network boundary and matches its documented retry policy (rare in this repo).
- **Runner infrastructure:** checkout, dependency download, disk, host process, or runner service fails independently of the test command. Require direct runner evidence before assigning this class. Where a self-hosted pool exposes no host metrics, say so and classify from what the logs do carry: one signature repeating across unrelated branches on one pool is evidence of shared-host contention even when the host cannot be inspected.

If evidence supports more than one independent fact, report each one. Do not collapse a timeout, signal, exit code, and assertion into a single inferred outcome.

## Reproduce the smallest relevant topology

Start with the owning test file or focused test name. Increase concurrency only to the first topology that reproduces the signature:

1. one test process (`go test -run '^TestName$' ./modulepath/` or the focused `make test-backend#TestName` / `make test-integration#TestName`);
2. concurrent tests or files in the same package, including under `-count N`;
3. multiple independent `go test` processes or `make` invocations on one host;
4. the owning `make` target with its configured worker count (`make test-backend`, `make test-integration`, `make test-migration`);
5. separate jobs or runner processes sharing the implicated host resource, including the Windows runner when the failure is OS-specific.

Match the active Go build tags, env vars (`GITEA_TEST_DATABASE`, `GITEA_TEST_LOG_SQL`, `GITEA_TEST_E2E_TIMEOUT_FACTOR`), source/build mode, and platform. Do not lower a production timeout or add random load merely to manufacture a different failure.

Where the signature belongs to a platform the available host cannot run, the ladder stops at the last reachable rung. Record that limit rather than substituting a passing run on another platform, then use CI as the reproduction, changing one suspected owner per run so the result stays attributable.

For a suspected race, replace probabilistic timing with a barrier at the contested transition. For a suspected shared-database collision, prove that two goroutines or processes commit conflicting rows in the same window, or prove that `unittest.LoadFixtures()` removes the conflict.

For a suspected platform mismatch, the minimal reproduction is the same test on the matching runner; do not try to manufacture a Windows failure on Linux.

## Fix at the owner

When implementation is authorized, fix the component that allocates, publishes readiness, mutates global state, or owns teardown. Do not hide the failure in a snapshot normalizer, retry wrapper, broader timeout, global serialization setting, or weaker assertion.

Keep stable fixture data separate from live resource allocation. A recorded URL can remain stable while the fixture maps its transport to an OS-assigned port; a stable expected path can remain an assertion without becoming a shared writable directory.

## Close the investigation

The evidence is complete when:

- the original signature has a supported classification;
- the smallest relevant topology reproduces it, or the external evidence is sufficient and the reproduction limit is explicit;
- an authorized fix fails under a negative control or pre-fix state and passes under the same topology afterward;
- any concurrent-process, restoration, or quiescent-teardown proof required by the resource owner passes;
- remaining Actions checks are reported as passing, pending, skipped, or failing from their observed state.

Do not run until a test happens to pass and call that result stable. Stop after the selected evidence establishes the conclusion, or report the missing fact that blocks classification.