# Generator: `monorepo_ts_workspaces`

pnpm workspaces root: `package.json` + `pnpm-workspace.yaml` for TypeScript monorepos.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `monorepo_ts_workspaces` |
| Version | `0.1.0` |
| Package | `generators/monorepo_ts_workspaces` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `base_project` | Runs against the project root the base generator establishes |

---

## Answers consumed

| Answer | Type | Used for |
|--------|------|----------|
| `project_name` | string | The workspace root's `name` field (defaults to `monorepo`) |

---

## Files written

| Path | Description |
|------|-------------|
| `pnpm-workspace.yaml` | Declares `apps/*` as the workspace package globs |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `name`, `version`, `private`, `workspaces` |

Declares no third-party dependencies — the workspace root is a container, and each
app under `apps/*` brings its own. The `version` field here is the *scaffolded
project's* own version, not a dependency Pin.

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `pnpm-workspace.yaml` | `file_exists` | File is present after generation |

---

## Post-generation commands

None.

## Test commands

None.

---

## Conflicts

Owns the root `package.json` identity keys (`name`, `version`, `private`, `workspaces`).
Generators that scaffold an individual app write under `apps/<name>/` and do not
collide with it.
