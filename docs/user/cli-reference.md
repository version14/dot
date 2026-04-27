# CLI Reference

Complete reference for every `dot` command, flag, and exit code.

---

## Global synopsis

```
dot <command> [flags] [arguments]
```

Flags are command-specific. There are no global flags — every flag is scoped to the subcommand that defines it.

---

## Commands

### `dot scaffold`

Run an interactive scaffolding flow, generate a project, and execute post-generation commands.

```
dot scaffold [flow-id] [-out DIR] [-skip-post]
```

**Arguments**

| Argument | Description |
|----------|-------------|
| `flow-id` | Optional. ID of the flow to run (e.g. `monorepo`, `fullstack`). When omitted and more than one flow is registered, DOT lists available flows and exits. |

**Flags**

| Flag | Default | Description |
|------|---------|-------------|
| `-out DIR` | `.` (current directory) | Parent directory where the project will be created. The project name (from the `project_name` answer) is appended as a subdirectory. |
| `-skip-post` | `false` | Skip all `PostGenerationCommands`. Useful for offline runs or CI environments that handle installation separately. |

**Pipeline**

1. The selected flow's question graph is presented as an interactive TUI form.
2. Answers are used to build a `ProjectSpec`.
3. The flow's `Generators` resolver produces an invocation list.
4. Active plugins may append additional invocations via `ResolveExtras`.
5. Invocations are topologically sorted (dependency order) and executed.
6. The virtual filesystem is persisted to `<out>/<project_name>/`.
7. `.dot/spec.json` and `.dot/manifest.json` are written.
8. `PostGenerationCommands` from all manifests run in order (unless `-skip-post`).

**Output**

Each generator prints a status line. Post-gen commands show a live spinner with elapsed time. Only failure output is printed.

**Exit codes**

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Scaffold error, flow aborted, or post-gen command failed |
| `2` | Usage error (unknown flag, unknown flow) |

**Examples**

```bash
dot scaffold                          # pick flow interactively
dot scaffold monorepo                 # use the monorepo flow
dot scaffold fullstack -out /tmp      # write to /tmp/<project_name>/
dot scaffold monorepo --skip-post     # generate files only, no pnpm install
```

---

### `dot update`

Re-run generators against an existing DOT project without re-prompting. Reads answers from `.dot/spec.json`.

```
dot update [PATH]
```

**Arguments**

| Argument | Default | Description |
|----------|---------|-------------|
| `PATH` | `.` | Path to the project root (the directory containing `.dot/`). |

**What it does**

1. Loads `.dot/spec.json` to recover the original answers and flow ID.
2. Re-runs the generator pipeline with those answers.
3. Persists changed files back to `PATH`.

Use this after upgrading DOT (new generator versions), after a generator adds new files, or to re-apply templates after editing them locally.

**Exit codes**

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | Error reading spec, generator failure, or persist error |
| `2` | Usage error |

**Example**

```bash
dot update ./my-project
dot update .
```

---

### `dot doctor`

Diagnose drift between a project's `.dot/spec.json` and the currently installed generators. Does not write any files.

```
dot doctor [PATH]
```

**Arguments**

| Argument | Default | Description |
|----------|---------|-------------|
| `PATH` | `.` | Path to the project root. |

**Checks performed**

1. **Flow exists** — the `flow_id` in `spec.json` is still registered.
2. **Generators found** — every generator recorded in `manifest.json` exists in the current registry.
3. **Version drift** — compares recorded versions against installed versions using semver constraints.
4. **Validators** — runs every `Check` in each manifest's `Validators` against the on-disk files.

**Output**

```
✓ flow: monorepo
✓ generators: all found
⚠ version drift: base_project spec=0.1.0 installed=0.2.0
✓ validators: 12 passed
```

**Exit codes**

| Code | Meaning |
|------|---------|
| `0` | No issues detected |
| `1` | One or more checks failed |
| `2` | Usage error or spec not found |

---

### `dot flows`

List all registered flows with their ID, title, and description.

```
dot flows
```

No flags. No arguments.

**Example output**

```
Flows
  monorepo         Turborepo monorepo (apps + shared packages)
                   Full monorepo with pnpm workspaces, Turborepo, and shared packages.
  fullstack        Full-stack app (Next.js + Go API)
  microservices    Go microservices (multiple services, Docker Compose)
  plugin-template  Plugin Repository Template
                   Scaffold a publishable DOT plugin repo (go.mod + plugin.go + manifest + README + LICENSE).
```

---

### `dot generators`

List all registered generators with their name, version, and description.

```
dot generators
```

No flags. No arguments.

**Example output**

```
Generators
  base_project        v0.1.0    Universal project scaffolding (README, .gitignore, LICENSE)
  typescript_base     v0.1.0    TypeScript base (tsconfig, package.json, tooling)
  react_app           v0.1.0    React application (Vite, React Router, Tailwind)
    depends on: typescript_base
  biome_config        v0.1.0    Biome formatter + linter config
  service_writer      v0.1.0    Go microservice (HTTP server, Dockerfile, healthcheck)
  plugin_repo_skeleton v0.1.0   DOT plugin repository skeleton
```

---

### `dot plugin list`

List built-in plugins (compiled in) and user-installed plugins.

```
dot plugin list
dot plugin ls
```

No flags. No arguments.

**Example output**

```
Built-in plugins
  biome_extras              1 generators · 2 injections

Installed plugins (~/.dot/plugins)
  my-plugin         v0.3.0   Adds extra templates for my stack
```

---

### `dot plugin install`

Install a plugin from a remote git repository or a local path.

```
dot plugin install <source> [-ref REF]
dot plugin install -from PATH
```

**Arguments / flags**

| Flag / Arg | Description |
|------------|-------------|
| `source` | Positional. Remote source (see table below). |
| `-ref REF` | Git ref (tag, branch, commit SHA) to check out after cloning. |
| `-from PATH` | Local directory to copy into the plugin store. Mutually exclusive with `source`. |
| `-id ID` | Override the plugin ID read from `plugin.json`. |
| `-version VER` | Override the version read from `plugin.json`. |

**Source forms**

| Source | Resolved URL |
|--------|-------------|
| `github.com/owner/repo` | `https://github.com/owner/repo.git` |
| `github.com/owner/repo@v1.2.0` | Same URL, checkout `v1.2.0` |
| `https://example.com/path/repo.git` | Used as-is |
| `git@github.com:owner/repo.git` | SSH, checkout with `-ref` |

**Install process**

1. Clone to a temporary staging directory.
2. Read `plugin.json` to extract `id` and `version`.
3. Atomic rename to `~/.dot/plugins/<id>/`.
4. On failure: staging directory is removed, no partial state.

**Exit codes**

| Code | Meaning |
|------|---------|
| `0` | Successfully installed |
| `1` | Clone failed, `plugin.json` missing/invalid, or rename error |
| `2` | Usage error (no source provided, conflicting flags) |

**Examples**

```bash
dot plugin install github.com/version14/dot-plugin-biome-extras
dot plugin install github.com/me/my-plugin -ref v0.2.0
dot plugin install git@github.com:me/my-plugin.git -ref main
dot plugin install -from ./my-local-plugin
```

---

### `dot plugin uninstall`

Remove an installed plugin.

```
dot plugin uninstall <id>
dot plugin remove <id>
dot plugin rm <id>
```

**Arguments**

| Argument | Description |
|----------|-------------|
| `id` | Required. The plugin ID as shown in `dot plugin list`. |

Removes `~/.dot/plugins/<id>/`. Idempotent — no error if the plugin is not installed.

**Exit codes**

| Code | Meaning |
|------|---------|
| `0` | Success (or already not installed) |
| `1` | Filesystem error |
| `2` | No ID provided |

---

### `dot version`

Print the tool version and exit.

```
dot version
dot --version
dot -v
```

**Example output**

```
dot 0.1.0
```

---

### `dot help`

Print the command summary and exit.

```
dot help
dot --help
dot -h
```

---

## Exit code summary

| Code | Meaning |
|------|---------|
| `0` | Command completed successfully |
| `1` | Runtime error (scaffold failed, plugin install failed, etc.) |
| `2` | Usage error (unknown command, missing required argument, unknown flag) |
