# Repository Restructure Migration

This document describes the ongoing repository restructuring that reorganizes
frontend, lint/tooling configs, CI/CD, and Go backend code into a clean
monorepo layout.

## Overview

The restructure isolates frontend and backend code into separate top-level
directories, cleaning up the repository root. Five pillars are involved:

| Pillar | From | To |
|--------|------|----|
| Frontend source | `web_src/` | `frontend/web_src/` |
| Frontend configs | root (`package.json`, `vite.config.ts`, etc.) | `frontend/` |
| Go backend | `cmd/`, `models/`, `modules/`, `routers/`, `services/`, etc. | `backend/` |
| Lint configs | root (`.golangci.yml`, `.markdownlint.yaml`, etc.) | `.config/lint/` |
| Tool configs | root (`.air.toml`, `pyproject.toml`, etc.) | `.config/tools/` |
| CI/CD actions & workflows | `.github/actions/`, `.github/workflows/` | `.gitea/actions/`, `.gitea/workflows/` |

## Target Directory Layout

```
.
├── .agents/                    # Agent skill definitions (new)
├── .config/
│   ├── lint/                   # Lint configs consolidated from root
│   │   ├── .golangci.yml
│   │   ├── .markdownlint.yaml
│   │   ├── .shellcheckrc
│   │   ├── .spectral.yaml
│   │   └── .yamllint.yaml
│   └── tools/                  # Tool configs consolidated from root
│       ├── .air.toml
│       ├── .envrc
│       ├── pyproject.toml
│       └── uv.lock
├── .gitea/
│   ├── actions/                # CI composite actions (moved from .github/actions/)
│   └── workflows/              # CI workflows (moved from .github/workflows/)
├── .github/
│   ├── workflows/              # New GitHub-specific workflows (restructured)
│   ├── ISSUE_TEMPLATE/         # New issue templates (.md format)
│   └── issue-management/       # Issue lifecycle automation (new)
├── backend/                    # All Go backend code
│   ├── build/                  # Build/generate scripts (standalone)
│   ├── cmd/                    # CLI entry points
│   ├── modelmigration/         # DB schema migrations
│   ├── models/                 # Data models
│   ├── modules/                # Shared Go packages (largest: 1,039 files)
│   ├── routers/                # HTTP handlers
│   └── services/               # Business logic
├── frontend/                   # All frontend code and configs
│   ├── eslint.config.ts
│   ├── eslint.json.config.ts
│   ├── package.json
│   ├── pnpm-lock.yaml
│   ├── pnpm-workspace.yaml
│   ├── stylelint.config.ts
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── types.d.ts
│   ├── vite.config.ts
│   ├── vitest.config.ts
│   └── web_src/
│       ├── css/                # Stylesheets
│       ├── js/                 # TypeScript, Vue, features, modules
│       ├── fomantic/           # Vendored Fomantic UI
│       └── svg/                # SVG icon sources
├── assets/                     # Embedded static assets (unchanged)
├── contrib/
│   ├── development/            # Community dev tooling (unchanged)
│   └── packaging/              # Distribution/packaging (moved from root)
│       ├── docker/             # Docker image rootfs (was root docker/)
│       ├── nix/                # Nix flake (was root flake.nix/.lock)
│       └── snap/               # Snapcraft packaging (was root snap/)
├── custom/                     # User configuration (unchanged)
├── docs/                       # Documentation (unchanged)
├── node_modules -> frontend/node_modules   # Symlink created by Makefile
├── options/                    # Locale files (unchanged)
├── public/                     # Built frontend assets (unchanged)
├── templates/                  # Go HTML templates (unchanged)
├── tests/                      # Integration & e2e tests (unchanged)
├── tools/                      # Mixed Go + TS tool scripts (unchanged)
├── main.go                     # Application entry point (unchanged)
├── main_timezones.go           # Windows timezone support (unchanged)
├── go.mod                      # Module definition
├── go.sum
└── Makefile                    # Updated to reference new paths
```

### Packages staying at root

| Directory | Reason |
|-----------|--------|
| `main.go` | Application entry point; must stay at module root |
| `go.mod` / `go.sum` | Module definition |
| `templates/` | Go HTML templates loaded by both backend and frontend tooling |
| `options/` | Locale files embedded via `go:embed` |
| `public/` | Built frontend assets served by the backend |
| `assets/` | Embedded static assets (go-bindata) |
| `custom/` | User configuration overrides |
| `tests/` | Integration tests (self-contained, not imported by production code) |
| `tools/` | Mixed Go + TS tool scripts (standalone, not imported) |
| `docs/` | Documentation |
| `contrib/` | Community contributions and extras |
| `contrib/packaging/` | Distribution assets: `docker/` (image rootfs), `nix/` (flake), `snap/` (snapcraft) |

## Go Backend Changes

### Phase 1: Path-reference updates (completed)

The first phase updated Go files that referenced frontend paths to accommodate
the `web_src/` → `frontend/web_src/` move.

| File | Change |
|------|--------|
| `modules/public/vitedev.go` | Added `web_src/` → `frontend/web_src/` disk-path mapping in `viteDevModuleID()` |
| `modules/util/color.go` | Comment: `web_src/` → `frontend/web_src/` |
| `services/websocket/events.go` | Comment: `web_src/` → `frontend/web_src/` |
| `services/webtheme/webtheme.go` | Vite dev-mode path: `"web_src/css/themes"` → `"frontend/web_src/css/themes"` |
| `tools/lint-go-all.go` | Added `--config=.config/lint/.golangci.yml` to golangci-lint invocations |
| `go.mod` | `go generate` ignore: `./web_src` → `./frontend/web_src` |

### Phase 2: Backend directory consolidation (proposed)

Move all Go backend packages into a `backend/` directory to create a clean
separation between frontend, backend, and shared resources.

#### Packages to move into `backend/`

| Package | Files | Subdirs | Description |
|---------|-------|---------|-------------|
| `cmd/` | 56 | `cmdtest/` | CLI entry points and subcommands |
| `modelmigration/` | 341 | 26 version dirs (`v1_6/`-`v1_27/`, `v28/`), `base/`, `migrationtest/`, `fixtures/` | DB schema migrations |
| `models/` | 380 | 28 subdirs | Data models and ORM definitions |
| `modules/` | 1,039 | 83 subdirs | Shared Go packages (largest package) |
| `routers/` | 464 | `api/`, `web/`, `private/`, `install/`, `common/`, `utils/` | HTTP route handlers |
| `services/` | 531 | 42 subdirs | Business logic layer |
| `build/` | 9 | `openapi3gen/` | Build/generate scripts (standalone) |

**Total: ~2,820 Go files** to be relocated.

#### Import path changes

The Go `module` path stays `gitea.dev` (defined by the root `go.mod`). Because
the packages move under the `backend/` subdirectory, every import path gains the
`backend/` segment, since Go import paths are `module + relative-path-from-go.mod`:

```
# Before
import "gitea.dev/modules/setting"
import "gitea.dev/services/pull"
import "gitea.dev/routers/api/v1/utils"

# After
import "gitea.dev/backend/modules/setting"
import "gitea.dev/backend/services/pull"
import "gitea.dev/backend/routers/api/v1/utils"
```

Packages left at the module root (`options/`, `tests/`, `tools/`, `main.go`) keep
their `gitea.dev/...` import paths unchanged (no `backend/` segment).

#### `main.go` changes

The root entry point imports the backend's `cmd` package:

```go
// main.go
import (
    "gitea.dev/backend/cmd"
    // ... other imports
)
```

#### Estimated impact

| Metric | Count |
|--------|-------|
| Go files requiring import path updates | ~2,820 |
| Unique import statements to rewrite | ~8,500+ |
| `go.mod` module path | unchanged (`gitea.dev`) |

#### Migration strategy

1. **Create `backend/` directory** with target subdirectories
2. **Move packages** one at a time, updating imports after each move
3. **Use `sed` or `goimports`** to batch-update import paths
4. **Run `make fmt`** after each batch to ensure formatting
5. **Run `make lint-go`** to catch broken imports
6. **Update Makefile** targets to reference `backend/` paths
7. **Keep `go.mod` module path** as `gitea.dev` (imports gain `backend/` segment)
8. **Verify build** with `make build`
9. **Run tests** with `make test-backend`

#### Potential issues

| Issue | Mitigation |
|-------|------------|
| Cross-layer dependency: `services/repository/files/` imports `routers/api/v1/utils` | Refactor to remove reverse dependency before or during move |
| `modelmigration/` imports from `models/` and `modules/` | Move together; internal imports are self-contained |
| `build/` scripts use `go:generate` with current paths | Update generate directives after move |
| Integration tests in `tests/` import backend packages | Update test imports; `tests/` stays at root |
| Template embedding (`go:embed`) in `modules/templates` | Verify embed paths resolve correctly after move |

### URL/manifest paths vs disk paths (critical distinction)

The Vite manifest keys and URL paths **intentionally remain `web_src/...`** even
though the files moved to `frontend/web_src/`. This means:

**Must NOT change** — these are manifest/URL keys:
```go
public.AssetURI("web_src/js/index.ts")           // → /assets/js/index.C6Z2MRVQ.js
public.AssetURI("web_src/css/themes/theme-...")   // → /assets/css/theme-....CyAaQnn5.css
AssetCSSLinks "web_src/js/index.ts" "web_src/css/index.css"  // in templates
ctx.ScriptImport "web_src/js/index.ts"             // in templates
```

**Must update** — these are disk-path references:
```go
os.DirFS(setting.StaticRootPath), "frontend/web_src/css/themes"   // disk path
// Keep this in sync with frontend/web_src/js/utils/color.ts      // comment
```

The bridge between the two is `viteDevModuleID()` in `modules/public/vitedev.go`:
```go
func viteDevModuleID(srcPath string) string {
    // Translate manifest key to disk path
    if fsPath, ok := strings.CutPrefix(srcPath, "web_src/"); ok {
        srcPath = "frontend/web_src/" + fsPath
    }
    return filepath.ToSlash(util.FilePathJoinAbs(setting.StaticRootPath, srcPath))
}
```

### Cosmetic-only stale references (low priority)

Two comments still reference `vite.config.ts` without the `frontend/` prefix:

- `modules/public/vitedev.go:170` — comment "see vite.config.ts: fs.allow"
- `build/generate-go-licenses.go:22` — comment "also defined in vite.config.ts"

These are harmless but could be updated for accuracy.

## Makefile Changes

All Makefile targets have been updated to reference new paths:

| Target | Key changes |
|--------|------------|
| `node_modules` | Runs `cd frontend && pnpm install --frozen-lockfile`, creates root symlink |
| `vite` | `vite build --config frontend/vite.config.ts` |
| `watch-frontend` | `vite --config frontend/vite.config.ts` |
| `watch-backend` | `air -c .config/tools/.air.toml` |
| `lint-js` | `eslint --config frontend/eslint.config.ts`, `vue-tsc -p frontend/tsconfig.json` |
| `lint-css` | `stylelint --config frontend/stylelint.config.ts` |
| `lint-json` | `eslint -c frontend/eslint.json.config.ts` |
| `lint-md` | `markdownlint --config .config/lint/.markdownlint.yaml` |
| `lint-swagger` | `spectral lint --ruleset .config/lint/.spectral.yaml` |
| `lint-yaml` | `yamllint -c .config/lint/.yamllint.yaml` |
| `lint-actions` | `uv run --project .config/tools` |
| `lint-shell` | `shellcheck --rcfile=.config/lint/.shellcheckrc` |
| `lint-templates` | `uv run --project .config/tools` |
| `test-frontend` | `vitest --config frontend/vitest.config.ts` |
| `lockfile-check` | Runs in `frontend/` directory |
| `clean-all` | Removes both `node_modules` and `frontend/node_modules` |

### Pending Makefile updates (backend consolidation)

After moving Go packages into `backend/`, these targets need updating:

| Target | Required change |
|--------|----------------|
| `build` | Update `go build` paths to `backend/` |
| `fmt` | Update `gofmt` paths to `backend/**/*.go` |
| `lint-go` | Update `golangci-lint` paths to `backend/` |
| `test-backend` | Update `go test` paths to `backend/...` |
| `generate-swagger` | Update `go generate` paths to `backend/` |
| `watch-backend` | Update `air` config to watch `backend/` |

Key variables updated:

```makefile
FRONTEND_SOURCES := $(shell find frontend/web_src/js frontend/web_src/css -type f)
FRONTEND_CONFIGS := frontend/vite.config.ts frontend/tailwind.config.ts
WEB_DIRS := frontend/web_src/js frontend/web_src/css
ESLINT_FILES := frontend/web_src/js tools frontend/*.ts tests/e2e
STYLELINT_FILES := frontend/web_src/css frontend/web_src/js/components/*.vue
```

## Frontend Tool Execution Pattern

Tools are now invoked via `node_modules/.bin/<tool>` (the root symlink resolves
to `frontend/node_modules/.bin/`) with explicit `--config` flags:

```bash
# Before (pnpm workspace-based)
pnpm exec eslint ...
pnpm exec vite build

# After (direct binary via symlink)
node_modules/.bin/eslint --config frontend/eslint.config.ts ...
node_modules/.bin/vite build --config frontend/vite.config.ts
```

Python/uv tools use `--project .config/tools` to find their configs:

```bash
# Before
uv run --frozen yamllint -s .

# After
uv run --project .config/tools --frozen yamllint -s -c .config/lint/.yamllint.yaml .
```

## Template References

Go HTML templates reference frontend assets using **manifest keys** (not disk
paths). All `{{AssetURI "web_src/..."}}` calls in templates are correct and
must not be changed:

```html
<!-- templates/base/head_style.tmpl -->
{{AssetCSSLinks "web_src/js/index.ts" "web_src/css/index.css"}}

<!-- templates/base/head_script.tmpl -->
{{ctx.ScriptImport "web_src/js/iife.ts"}}

<!-- templates/swagger/openapi-viewer.tmpl -->
{{AssetCSSLinks "web_src/js/swagger.ts" "web_src/css/swagger-standalone.css"}}
{{ctx.ScriptImport "web_src/js/swagger.ts" "module"}}
```

## Architecture Benefits

### Before restructure

```
.
├── cmd/                    # Go backend (mixed with frontend at root)
├── models/                 #
├── modules/                #
├── routers/                #
├── services/               #
├── web_src/                # Frontend source
├── package.json            # Frontend config
├── vite.config.ts          #
├── .golangci.yml           # Lint config
├── .github/workflows/      # CI/CD
└── ...50+ root files
```

Root directory is cluttered with ~50+ files mixing frontend, backend, and
tooling concerns.

### After restructure

```
.
├── backend/                # All Go backend code
├── frontend/               # All frontend code and configs
├── .config/                # Lint and tool configs
├── .gitea/                 # CI/CD (Gitea-native)
├── docs/                   # Documentation
├── templates/              # Go HTML templates (shared)
├── options/                # Locale files (shared)
├── main.go                 # Entry point
├── go.mod                  # Module definition
└── Makefile                # Build system
```

Clear separation of concerns:
- `backend/` — Go code only, self-contained
- `frontend/` — TypeScript/Vue/CSS, self-contained
- Root — minimal glue (entry point, module definition, shared resources)

### Dependency analysis

#### Import dependency matrix (current, pre-move)

| Importing → Imported | `cmd/` | `models/` | `modules/` | `routers/` | `services/` |
|---------------------|--------|-----------|------------|------------|-------------|
| `cmd/`              | —      | 39        | 39         | 1          | 17          |
| `models/`           | —      | 319       | 583        | —          | —           |
| `modules/`          | —      | 52        | 583        | 2          | 2           |
| `routers/`          | —      | 397       | 397        | 156        | 387         |
| `services/`         | —      | 418       | 418        | 3          | 157         |
| `tests/`            | —      | 226       | 226        | 23         | 70          |

**Key observations:**
- `modules/` is the most-imported package (shared library)
- `models/` and `modules/` have no dependency on `routers/` or `services/`
- `routers/` depends on `services/` (correct layering)
- `services/` has one reverse dependency on `routers/` (needs refactoring)

#### Dependency direction (target)

```
main.go → cmd → routers → services → models
                     ↑        ↑
                     └── modules (shared)
```

The `backend/` directory maintains this layering:
- `cmd/` imports from `routers/`, `modules/`, `services/`
- `routers/` imports from `services/`, `modules/`, `models/`
- `services/` imports from `modules/`, `models/`
- `modules/` imports from `models/` (and self)
- `models/` imports from `modules/` (for utilities)
- `modelmigration/` imports from `models/`, `modules/`

## What Still Needs To Be Done

### Completed
- [x] Frontend files moved to `frontend/web_src/`
- [x] Frontend configs moved to `frontend/`
- [x] Lint configs moved to `.config/lint/`
- [x] Tool configs moved to `.config/tools/`
- [x] CI actions moved to `.gitea/actions/` and `.gitea/workflows/`
- [x] Go backend path references updated (5 files)
- [x] Makefile updated for all new paths
- [x] `go.mod` updated
- [x] Documentation updated (`development.md`, `guidelines-frontend.md`, etc.)
- [x] Template comments updated (action_status, commit_status)
- [x] Shell scripts updated (`playwright.sh`, `test-e2e.sh`)
- [x] Stale comment in `backend/modules/public/vitedev.go:170` (`vite.config.ts` → `frontend/vite.config.ts`)
- [x] Stale comment in `backend/build/generate-go-licenses.go:22` (`vite.config.ts` → `frontend/vite.config.ts`)
- [x] `tools/shared.ts:7` comment references `web_src/js/webcomponents` → `frontend/web_src/js/webcomponents`

### Frontend Toolchain Verification (root-cwd execution model)
- [x] Frontend build stage in `Dockerfile` / `Dockerfile.rootless`: `COPY frontend/package.json frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml /src/frontend/`, then `cd /src/frontend && pnpm install --frozen-lockfile`. Verified by running the install in an isolated temp tree against `pnpm@11.24.0` (`pnpm install --frozen-lockfile` succeeded in 13.5s; `.bin/eslint` produced).
- [x] `.dockerignore` excludes `/frontend/node_modules`, `/.pnpm-store` so the pnpm cache never reaches the image.
- [x] Frontend tools (`eslint`, `vue-tsc`, `stylelint`, `spectral`, `markdownlint`, `vitest`, `vite`) and TS build scripts (`tools/generate-svg.ts`) run from repo root cwd via `node_modules/.bin/<tool>`; the root `node_modules` symlink → `frontend/node_modules` resolves all binaries.
- [x] `eslint.config.ts` parserOptions pinned to `project: ['tsconfig.json']` + `tsconfigRootDir: import.meta.dirname` so `tools/` and `tests/e2e` find the project.
- [x] `eslint.json.config.ts` JSONC block matches `**/tsconfig*.json` (was `tsconfig.json` only).
- [x] Minimal `tools/package.json` and `tests/e2e/package.json` give `eslint-plugin-import-x` `pkgUp` a nearest anchor; `tools/**` and `tests/**` rule block sets `import-x/no-extraneous-dependencies` `packageDir: 'frontend'`.
- [x] Asset imports deepened by `web_src` → `frontend/web_src` move patched in `frontend/web_src/js/svg.ts`, `frontend/web_src/js/webcomponents/overflow-menu.ts`, `frontend/web_src/js/features/emoji.ts`, `frontend/web_src/js/utils/match.ts`, `frontend/web_src/js/modules/codeeditor/main.ts`.

### Pending: Backend Directory Consolidation
- [x] Create `backend/` directory structure
- [x] Move `cmd/` → `backend/cmd/`
- [x] Move `models/` → `backend/models/`
- [x] Move `modules/` → `backend/modules/`
- [x] Move `routers/` → `backend/routers/`
- [x] Move `services/` → `backend/services/`
- [x] Move `modelmigration/` → `backend/modelmigration/`
- [x] Move `build/` → `backend/build/`
- [x] Batch-update all import statements (~2,540 files, `gitea.dev/X` → `gitea.dev/backend/X`)
- [x] Update `main.go` to import `gitea.dev/backend/cmd`
- [x] Update Makefile targets for `backend/` paths
- [x] Refactor cross-layer dependency: `RefCommit` moved to `modules/git`
- [ ] Verify build: `make build`
- [ ] Verify tests: `make test-backend`
- [ ] Verify lint: `make lint-go`
- [ ] Regenerate swagger (``make generate-swagger``) to refresh `templates/swagger/*.json`

### Root Cleanup: packaging consolidation (completed)
- [x] Move `docker/` → `contrib/packaging/docker/`
- [x] Move `snap/` → `contrib/packaging/snap/`
- [x] Move `flake.nix` / `flake.lock` → `contrib/packaging/nix/`
- [x] Update `Dockerfile` / `Dockerfile.rootless` `COPY` paths to `contrib/packaging/docker/root`
- [x] Update `.dockerignore` / `.gitignore` snapcraft path to `contrib/packaging/snap/.snapcraft/`
- [x] Update `files-changed.yml` docker trigger to `contrib/packaging/docker/**`
- [x] Update snapcraft workflow to run `snapcraft pack` inside `contrib/packaging/snap` and output the snap path
- [x] Update `snapcraft.yaml` `source`, `icon`, and override script paths
- [x] Update snap pull/build scripts' root-directory guard to `contrib/packaging/snap`

### Remaining / Low Priority
- `options/fileicon/material-icon-rules.json` still maps `pnpm-lock.yaml` / `pnpm-workspace.yaml` to the pnpm icon by **basename**; that file is `linguist-generated` and editors match by basename, so the entries remain effective for `frontend/pnpm-lock.yaml`. No regeneration needed.

## File Change Statistics

### Current changes (frontend + configs + CI/CD)

| Category | Deleted | Modified | New (Untracked) |
|----------|---------|----------|-----------------|
| Frontend source | 412 | 0 | 440 (in `frontend/`) |
| Frontend configs | 12 | 0 | 12 (in `frontend/`) |
| Lint/tool configs | 9 | 0 | 9 (in `.config/`) |
| CI/CD | 30 | 2 | 44 (in `.gitea/` + `.github/`) |
| Go backend | 0 | 6 | 0 |
| Documentation | 0 | 5 | 6 |
| Makefile | 0 | 1 | 0 |
| Templates | 0 | 3 | 0 |
| Other | 0 | 6 | 0 |
| **Total** | **463** | **23** | **511** |

### Pending changes (backend consolidation)

| Category | Estimated Impact |
|----------|-----------------|
| Go files moved | ~2,540 files into `backend/` |
| Import paths rewritten | ~8,500+ statements |
| `go.mod` module path | unchanged (`gitea.dev`) |
| Makefile targets | ~10 targets updated |
| Total files touched | ~2,550+ |

## Migration Execution Order

### Phase 1: Frontend consolidation (completed)
1. Move `web_src/` → `frontend/web_src/`
2. Move frontend configs → `frontend/`
3. Move lint configs → `.config/lint/`
4. Move tool configs → `.config/tools/`
5. Move CI/CD → `.gitea/`
6. Update Go path references
7. Update Makefile

### Phase 2: Backend consolidation (proposed)
1. **Prepare** — Fix cross-layer dependency (`services/repository/files/` → `routers/api/v1/utils`)
2. **Create structure** — `backend/` with subdirectories
3. **Move packages** (in order of dependency, leaves first):
   a. `modules/` → `backend/modules/` (most-imported, fewest outbound deps)
   b. `models/` → `backend/models/`
   c. `modelmigration/` → `backend/modelmigration/`
   d. `services/` → `backend/services/`
   e. `routers/` → `backend/routers/`
   f. `cmd/` → `backend/cmd/`
   g. `build/` → `backend/build/`
4. **Update imports** — Batch-rewrite all `gitea.dev/X` → `gitea.dev/backend/X`
5. **Keep `go.mod` module path** as `gitea.dev`
6. **Update `main.go`** — Import `gitea.dev/backend/cmd`
7. **Update Makefile** — All backend targets
8. **Verify** — Build, lint, tests

### Why move packages in this order?

Moving leaf packages first minimizes broken imports during the transition:
- `modules/` depends on almost nothing external → safe to move first
- `models/` depends on `modules/` → move after `modules/`
- `services/` depends on `models/` + `modules/` → move after both
- `routers/` depends on `services/` + `modules/` + `models/` → move after `services/`
- `cmd/` depends on everything → move last
- `build/` is standalone → can move anytime
