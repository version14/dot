# ADR-0001: Template Dependency Checker Architecture

**Status:** Superseded by [ADR-0002](./0002-dependency-catalog.md)
**Date:** 2026-05-28

> **Superseded 2026-07-11.** The "real-state extraction" premise below is false: calling
> `Generate()` with empty `Answers` executes one arbitrary branch, not real state, so
> dependencies behind conditionals were never checked. Worse, the scanner read *values* while
> the patcher rewrote *source literals* — a gap that made some dependencies structurally
> unpatchable, crashed the weekly run, and merged a mislabeled commit (`997e489`) to `main`.
> See [ADR-0002](./0002-dependency-catalog.md).

## Context

Generator files (e.g. `generators/express_server_typescript_deps/generator.go`) hardcode
dependency versions as Go map literals passed to `ctx.State.UpdateJSON("package.json", ...)`.
Dependabot is blind to these because they are embedded inside Go string literals, not in real
manifest files on disk. Over time these pinned versions drift silently: packages become
deprecated or superseded by new major versions with no automated signal.

We need a CI mechanism that:
1. Detects drift for any package ecosystem (npm, Go modules, Rust, Java, …)
2. Opens PRs for outdated packages (one PR per package, gated by test-flows)
3. Raises issues for deprecated packages (blocker signal requiring human decision)
4. Runs on a weekly CRON, offset from Dependabot to avoid PR pile-ups

## Decision

Build `tools/dep-checker/`, a standalone Go binary consistent with `tools/test-flow/`,
structured around two core ideas.

### 1. Real-state extraction

Rather than parsing Go source with regexes or AST rewriting, the tool instantiates each
generator with a real `state.VirtualProjectState` (empty, no files). It calls `Generate()`,
then reads back whatever files were written — `package.json`, `go.mod`, `Cargo.toml`, etc.
This is strictly accurate: it reuses the exact code path that produces real output during
scaffolding and can never drift from it. Adding a new generator automatically makes it
visible to the tool.

### 2. `RegistryChecker` interface

```go
type DepInfo struct {
    Latest     string
    Deprecated bool
    Notice     string
}

type RegistryChecker interface {
    Ecosystem() Ecosystem
    Check(pkg, currentVersion string) (DepInfo, error)
}
```

The scanner detects which file a generator wrote to and dispatches to the matching checker:

| File written | Ecosystem | Registry API |
|---|---|---|
| `package.json` | npm | `registry.npmjs.org/<pkg>/latest` |
| `go.mod` | Go | `proxy.golang.org/<module>/@latest` |
| `Cargo.toml` | Rust | `crates.io/api/v1/crates/<name>` |
| `pom.xml` | Java | `search.maven.org/solrsearch/select` |

Adding a new ecosystem = implementing one `RegistryChecker` and adding one entry to
`checkerFor()` in `checkers.go`.

### PR and issue strategy

- **Outdated (not deprecated):** open one PR per package (branch `deps/<ecosystem>/<package>`),
  gated by `make test-flows`. If test-flows fail, open an issue with the output instead.
- **Deprecated:** open the version-bump PR if a newer version exists under the same name,
  plus a separate tracking issue. If no replacement exists under the same name, issue only.
- **Idempotency:** skip if an open PR already exists on `deps/<ecosystem>/<package>`.
  Reopen closed issues if the package is still deprecated on the next run.
- **Schedule:** weekly on Wednesday, offset from Dependabot (Monday).

## Alternatives considered

| Option | Reason rejected |
|---|---|
| Regex over Go source | Fragile; breaks silently if map literal formatting changes |
| `PackageDependencies()` interface on generators | Requires migrating all 41 generators; creates a second source of truth that can drift |
| Go AST extraction | Heavy machinery for what the real `Generate()` already does for free |
| Separate binary per ecosystem | Duplicates orchestration logic; no shared `RegistryChecker` contract |
| Full project scaffolding + `npm outdated` | Slower; conflates first-party version choices with transitive drift |

## Consequences

- The tool imports generator packages directly, so it lives inside the module.
- Generators that require non-empty answers (e.g. `base_project`, `plugin_repo_skeleton`)
  will return an error when called with an empty context; the scanner logs a warning and
  skips them. This is acceptable since those generators do not declare npm/Go deps via
  `UpdateJSON`/`UpdateGoMod`.
- Template-based generators (those using `render.NewLocalFolderRenderer` with embedded
  files) may embed version constants that the real-state approach cannot see. These are
  tracked as a known limitation; the affected constant (`dotModuleVersion` in
  `plugin_repo_skeleton`) should be updated manually or via a future dedicated check.
- The `dotapi.State` interface remains the canonical write surface; the dep-checker
  depends on this contract being stable.
- Patching `generator.go` files uses targeted `strings.ReplaceAll` on the exact Go
  string literal pair (`"pkg": "^old"` → `"pkg": "^new"`); unique per file by construction.
