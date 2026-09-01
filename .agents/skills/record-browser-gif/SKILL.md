---
name: record-browser-gif
description: 'Use when recording a browser or Web UI interaction demo as an optimized GIF for a Gitea pull request that changes product-user-visible GUI behavior (which AGENTS.md requires to ship a screenshot or recording). Drives a real Gitea server booted from the pull request tree, captures state-based frames through the existing Playwright dev dependency, encodes them with a deterministic ffmpeg pipeline, and publishes the artifact to a dedicated assets branch only when the task includes attaching the GIF. Never substitutes fixtures, mocks, or test-only hooks for the real configuration path.'
---

# Record Browser GIF

Produce a short, truthful UI demonstration as a local GIF, and — only when the task includes attaching it to a pull request — publish it through the assets-branch workflow at the end of this skill. Use the bundled encoder for repeatable timing, dimensions, and size.

## Every GUI pull request includes a GIF

A pull request that changes product-user-visible GUI behavior MUST include a demonstration GIF recorded with this skill and embedded in the pull request body via [the assets-branch workflow](#publish-to-an-assets-branch).

The recording itself is part of the evidence: use a real Gitea server booted from that pull request's branch tree, a real `app.ini` with a scratch SQLite database, and a real session through the browser. Never substitute fixture data, mock transports, synthetic event injection, or test-only hooks for a working server unless the user explicitly asked for a fixture recording. Next to the embed, state the exact demonstrated commit SHA, the tree and origin that served it, and any mode flags or browser-state exceptions, so reviewers know exactly what the recording proves.

## Keep recording separate from publication

- Recording produces frame images and one local `.gif` artifact only; it never mutates remote state.
- Publication — pushing the GIF to an assets branch and embedding it in a pull request body — is the separate final step, performed only when the task includes attaching the GIF to a pull request. It never touches the pull request's own branch.
- Preserve the requested recording conditions. A real-server demo must not use fixture data, mock transports, synthetic event injection, or test-only hooks. If the server or its dependencies are unavailable, report that limitation instead of substituting a fixture.
- Never read or expose credential values. Use Gitea's normal configuration path (`GITEA_WORK_DIR` plus a fresh `custom/conf/app.ini`) and a benign demonstration session.

## Stage the application

A GIF for a specific pull request demonstrates that pull request's tree, so stage per pull request:

1. Require a clean worktree, record its exact commit with `git rev-parse HEAD`, then build that recorded tree with `make build` (which runs the `frontend` and `backend` sub-targets; see [docs/development.md](../../../docs/development.md)). There is no `pnpm run build` or `pnpm run build:web` to substitute — the pnpm project root is `frontend/`, and the Makefile drives pnpm for you. `make` is not on `PATH` on this Windows host by default; install it via MSYS2 or Chocolatey per [docs/build-setup.md](../../../docs/build-setup.md) or invoke it through your shell wrapper. For active development, `make watch`, `make watch-frontend`, or `make watch-backend` rebuilds in place — pick the target that matches what your GIF demonstrates. A GIF recorded against another commit's build misattributes the evidence.
2. Boot one server per port from that tree with fresh scratch state. Set `GITEA_WORK_DIR` to a fresh temporary directory, write a fresh `custom/conf/app.ini` there with `DB_TYPE = sqlite3` and `PATH = <GITEA_WORK_DIR>/data/gitea.db`, then run the built binary as `./gitea web` (POSIX) or `.\gitea.exe web` (Windows) in the background. `tools/test-e2e.sh` writes a known-good shape of this config during the existing e2e harness — mirror it. Give the browser a fresh isolated context or profile as well; if the browser workflow cannot create one, clear that origin's cookies and site storage before navigation so persisted client state cannot affect the evidence.
3. Treat one storyboard as one evidence run: every published frame comes from that server, that workdir, and that session. If capture automation fails, discard its frames and rerun from fresh state; never splice frames from separate runs.
4. When switching between pull requests, stop the old server by PID or an exact match on its command line. A broad `pkill -f` pattern can match and kill the shell that launched it — including your own.

Pick a scratch location the way your shell prefers:

```bash
GITEA_WORK_DIR=$(mktemp -d)
```

```powershell
$dir = New-Item -ItemType Directory -Path (Join-Path $env:TEMP ('gitea-' + [System.Guid]::NewGuid().Guid))
$env:GITEA_WORK_DIR = $dir.FullName
```

## Record the flow

1. Drive the browser through `@playwright/test` (already a dev dependency in `frontend/package.json`) or another browser tool already set up locally. Use the user's existing Chrome state only when requested or required; state that exception in the provenance and do not claim fresh client state. If browser automation is unavailable, prefer the repository-declared Playwright dependency in an isolated headless browser; do not install another driver or launch the user's browser without authorization. State that fallback in the provenance.
2. Before recording, identify the exact origin, whether the app is built or in development, and any local configuration overrides. Record only claims that the observed setup supports.
3. When a production default opens a native operating-system surface that headless automation cannot drive, select a documented Gitea configuration override through the app's normal `app.ini`. State the override in the provenance; a fixture or test-only hook is not an acceptable substitute.
4. Choose three to six states that tell one story, such as typed, running, settled, and detail. Prefer semantic state changes over continuous capture; omit loading churn that does not help the viewer.
5. Keep one viewport and crop for every frame, and name frames lexically: `00-initial.png`, `01-typed.png`, and so on.
6. Store frames under the repository's gitignored output location — `tests/e2e-output/` is the existing Playwright output directory (already in `.gitignore`), so write frames into a per-storyboard subdirectory there. Create the frame subdirectory first; writing into a missing directory fails at capture time:

    ```bash
    mkdir -p tests/e2e-output/gif-frames-<label>
    ```

    Windows (PowerShell):

    ```powershell
    New-Item -ItemType Directory -Force -Path "tests\e2e-output\gif-frames-<label>" | Out-Null
    ```
7. Before each screenshot, wait for a concrete UI condition such as a unique label, enabled control, changed document title, or completed response. Require the locator to resolve exactly one element; for Playwright accessible-name locators, use `exact: true` when equality is intended because descendant text can otherwise create a false match. Do not use a fixed delay as proof that the application reached the state.
8. Make completion predicates match an exact-text element — for example, an element whose trimmed text equals the expected reply — never a substring check such as `body.textContent.includes(...)`, which an unrelated echo also satisfies.
9. When the claim involves a server-side failure, a recovery path, or a configuration override, include a detail or follow-up frame that shows the affected page or error text. A successful landing page alone does not prove the recovery path behaved that way.
10. Capture a transient state (spinner, running row) by driving a slow foreground operation — for example, a long-running `make` task in another shell — and polling a concrete DOM marker (a `data-*` attribute) inside one browser-script call that also takes the screenshot. State polled across separate tool calls is lost, because the turn settles between calls.
11. Engineer the input so the state you need actually occurs: when an operation would otherwise return immediately, drive a longer foreground path and use a settle sentinel to anchor the completion predicate.
12. Capture no secrets, personal data, unrelated tabs, or transient notifications. Stop any unnecessarily long real run after the demonstrated state is visible.

Use the browser's own screenshot API. When it returns image bytes, save those bytes directly; the encoder detects image content independently of the filename extension.

## Encode the GIF

Probe for a working Python first; the binary name differs between shells and OS images, and `python` on `PATH` may be a Microsoft Store stub on Windows:

```bash
python3 --version
```

```powershell
python --version
py -3 --version
```

If none reports a working interpreter, install one before continuing. Do not rewrite the encoder to depend on another runtime.

Probe for the media binaries; on Windows, `Get-Command` is the equivalent of `which`:

```bash
command -v ffmpeg && command -v ffprobe
```

```powershell
Get-Command ffmpeg; Get-Command ffprobe
```

If either media binary is missing, report the dependency instead of installing software without authorization.

Set `GIF_SKILL_DIR` to this skill's absolute directory on its own line before the python command — an inline `GIF_SKILL_DIR=... python ...` assignment fails, because the argument expands before the assignment takes effect:

```bash
export GIF_SKILL_DIR=/absolute/path/to/this/skill
python3 "$GIF_SKILL_DIR/scripts/encode_gif.py" \
  /absolute/path/to/frames \
  /absolute/path/to/demo.gif \
  --durations 1.5,1.5,1.5,3.5 \
  --fps 10 \
  --max-width 1200 \
  --colors 128
```

```powershell
$env:GIF_SKILL_DIR = 'C:\absolute\path\to\this\skill'
python "$env:GIF_SKILL_DIR\scripts\encode_gif.py" `
  C:\absolute\path\to\frames `
  C:\absolute\path\to\demo.gif `
  --durations 1.5,1.5,1.5,3.5 `
  --fps 10 `
  --max-width 1200 `
  --colors 128
```

One duration applies to every frame; otherwise provide one comma-separated positive duration per frame, holding the final settled state longest. The encoder rejects fewer than two frames, mismatched dimensions or durations, invalid limits, accidental overwrite, unexpected duration, and output above `--max-bytes`. For a large artifact, reduce `--max-width` first, then `--colors` or `--fps`; retain readable text and the final state long enough to inspect. Use `--force` only after resolving the exact output path.

Verify the encoder still parses after local edits:

```bash
python3 -m py_compile "$GIF_SKILL_DIR/scripts/encode_gif.py"
```

```powershell
python -m py_compile "$env:GIF_SKILL_DIR\scripts\encode_gif.py"
```

## Verify the artifact

1. Read the encoder's JSON summary and confirm the output path, source and encoded frame counts, dimensions, duration, and byte size.
2. Visually read the encoded GIF itself, not only the source frames. Confirm that the transition is legible, the last state is held long enough, and no sensitive content appears. If the viewer renders only the first frame, decode representative frames from the encoded GIF with `ffmpeg` and inspect those; the pre-encode screenshots do not prove the encoded order, palette, or final hold.
3. Run `git status --short` and confirm frames and the artifact landed only under ignored paths.
4. Return the absolute GIF path, render it when the client supports local media, and state whether the recording used a real server, fixture data, or another transport. When the task does not include attaching the GIF to a pull request, stop here.

## Publish to an assets branch

Perform this step only when the task includes attaching the GIF to a pull request.

Never commit a GIF to the pull request's own branch or any branch that merges into a long-lived branch: binary media committed there bloats the repository history for every future clone. GIFs live on a dedicated orphan assets branch — a branch with no parent commit and nothing but media — and one assets branch serves a whole pull request series (named `<series>-assets`; list existing ones with `git ls-remote --heads origin '*assets*'`).

Before either workflow below pushes, verify that the assets branch contains media only and that the staged GIF's checksum matches the verified local artifact.

Work in a shallow single-branch scratch clone so the publication cannot touch your working tree. Pick a scratch location the way your shell prefers:

```bash
ASSETS_CHECKOUT=$(mktemp -d)
git clone --branch <assets-branch> --single-branch --depth 1 <repo-url> "$ASSETS_CHECKOUT"
cp /absolute/path/to/demo.gif "$ASSETS_CHECKOUT/<name>.gif"
git -C "$ASSETS_CHECKOUT" add <name>.gif
git -C "$ASSETS_CHECKOUT" commit -m "assets: <what it shows> gif (#<pr>)"
git -C "$ASSETS_CHECKOUT" push origin <assets-branch>
```

For a new series, use the same shallow clone pattern, then create the orphan branch with `git -C "$ASSETS_CHECKOUT" switch --orphan <assets-branch>`, then add, commit, and push the same way. Windows (PowerShell): the `git` invocation accepts the same arguments once a scratch directory exists, so the only platform change is picking the scratch path through PowerShell-native means (e.g. `Join-Path $env:TEMP ...`).

After pushing, use authenticated GitHub API or raw requests to confirm the remote path, byte size, checksum, `200` response, and `image/gif` content type. An anonymous `404` does not disprove a private-repository asset; authenticate the verification instead. This proves the repository-member review path, not public availability.

Immediately before editing the pull-request body, re-read its live head and compare it with the commit recorded next to the GIF. Stop and re-record when it moved. After the edit, re-read the live head and require it to remain at that recorded commit. Separately, render the body through GitHub's Markdown API and confirm that the expected `<img>` is present.

Embed the GIF in the pull request body with the raw blob URL; the `?raw=true` suffix is required, because the plain blob URL renders GitHub's file page instead of the image:

```markdown
![<alt text>](https://github.com/<owner>/<repo>/blob/<assets-branch>/<name>.gif?raw=true)
```

Never delete or rewrite an assets branch, and never force-push it: merged pull request bodies reference its URLs forever. Append new commits only.