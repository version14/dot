# Template Dependency Checker

`tools/dep-checker` is a Go binary that detects when package versions hardcoded inside
generator source files have become outdated or deprecated. It is the automated answer to
a gap Dependabot cannot fill: the npm, Go, Rust, and Maven packages that generators write
into scaffolded projects live inside Go map literals, invisible to standard dependency bots.

## Why this exists

Every generator declares its dependencies like this:

```go
// generators/express_server_typescript_deps/generator.go
d.Merge(map[string]interface{}{
    "dependencies": map[string]interface{}{
        "express": "^4.21.0",   // ← hardcoded, never touched by Dependabot
    },
})
```

When a package releases a new major version or gets deprecated, nothing updates these
strings automatically. A user running `dot scaffold` six months later gets a project
pinned to an old baseline. The dep-checker closes this loop.

## Architecture

The tool works in two phases: **scan** and **patch**.

### Scan phase

For each registered generator the scanner:

1. Creates a fresh `state.VirtualProjectState` (no files, empty spec).
2. Calls `generator.Generate(ctx)` with that state.
3. Reads back every file the generator wrote (`package.json`, `go.mod`, `Cargo.toml`, `pom.xml`).
4. Extracts `{package: versionConstraint}` entries from those files.
5. Queries the appropriate registry API for the latest version and deprecation status.

This approach is strictly accurate because it runs the real `Generate()` code path — the
same path that produces real scaffolded output. No source parsing, no second source of
truth.

### Patch phase

Given a generator name, package name, current version, and latest version, the patcher:

1. Resolves the file path: `generators/<name>/generator.go`.
2. Replaces the exact Go string literal pair `"pkg": "^old"` → `"pkg": "^new"`.
3. Preserves the constraint prefix (`^`, `~`, etc.) from the current version.

### Registry support

| Ecosystem | File detected | Registry API |
|---|---|---|
| npm | `package.json` | `registry.npmjs.org/<pkg>/latest` |
| Go modules | `go.mod` | `proxy.golang.org/<module>/@latest` |
| Rust | `Cargo.toml` | `crates.io/api/v1/crates/<name>` |
| Java/Maven | `pom.xml` | `search.maven.org` (groupId:artifactId format) |

To add a new ecosystem, implement `RegistryChecker` in `checkers.go` and add one entry
to `checkerFor()`.

## Running locally

```bash
# Scan all generators and write a report
go run ./tools/dep-checker scan --output=dep-report.json

# Inspect results
cat dep-report.json | jq '.entries[] | select(.outdated or .deprecated)'

# Patch a single generator (done automatically by CI)
go run ./tools/dep-checker patch \
  --generator=express_server_typescript_deps \
  --package=express \
  --current="^4.21.0" \
  --latest=5.1.0
```

## CI workflow

The GitHub Action at `.github/workflows/dep-checker.yml` runs every **Wednesday at 09:00 UTC**
(offset from Dependabot's Monday to avoid PR pile-ups) and on `workflow_dispatch`.

For each dependency group in the report:

```
scan
 └─ for each group:
     ├─ open PR already exists on deps/<ecosystem>/<group>? → skip
     ├─ patch all generators in the group
     ├─ make test-flows  (cache handles unaffected generators)
     │   ├─ PASS → push branch, open PR
     │   └─ FAIL → open issue with test output (no PR)
     └─ package deprecated?
         ├─ issue already open? → reopen if closed, add comment
         └─ no issue → create issue with deprecation notice
```

### PR format

- **Branch:** `deps/<ecosystem>/<group>` (e.g. `deps/npm/express`, `deps/npm/minor-and-patch`)
- **Title:** `chore(deps): bump <package>` (or multiple packages) `in templates`
- **Body:** lists affected generators, confirms test-flows passed, links to dep-checker

### Issue format

Two issue types are created:

**test-flow failure** (label: `dependencies`, `test-failure`)
- Title: `chore(deps): test-flow failure bumping <package> to <latest>`
- Body: full test-flow output so the root cause is immediately visible

**deprecated package** (label: `dependencies`, `deprecated`)
- Title: `deprecated: <package> (<ecosystem>) used in templates`
- Body: deprecation notice, affected generators, recommended action
- Reopened automatically on subsequent CRON runs if still deprecated

## Known limitations

- **Template-based generators** — `plugin_repo_skeleton` embeds `dotModuleVersion` as a
  Go constant rather than via `UpdateGoMod`. This value is invisible to the scanner.
  Update it manually when releasing a new stable version of `dot`.
- **Generators requiring answers** — generators that validate `ctx.Answers` at Generate
  time (e.g. `base_project`) are skipped with a warning. These generators do not declare
  npm/Go deps directly, so this is safe.
- **Go module deprecation** — the Go proxy API does not expose a deprecation field.
  Deprecated Go modules are detected via `go.mod` `Deprecated:` comments, which is not
  yet implemented. A future iteration can use `go list -m -json` for this.
- **Cargo and Maven** — implementations return version info but not deprecation status,
  as neither registry has a stable deprecation API. Treat their results as "outdated only".

## Adding a new generator with npm deps

No action needed. Any generator that calls `ctx.State.UpdateJSON("package.json", ...)` is
automatically picked up on the next scan. The dep-checker uses the same generator registry
as the scaffold pipeline.
