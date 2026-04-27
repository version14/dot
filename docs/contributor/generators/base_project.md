# Generator: `base_project`

Universal project scaffolding. Every flow runs this generator first; it writes the files that every project regardless of language or framework must have.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `base_project` |
| Version | `0.1.0` |
| Package | `generators/base_project` |

---

## Dependencies

None. `base_project` is the root of the dependency graph — every other generator depends on it, directly or transitively.

---

## Answers consumed

| Key | Type | Notes |
|-----|------|-------|
| `project_name` | string | Used as the heading in `README.md` and in the LICENSE copyright line. Falls back to `spec.Metadata.ProjectName`, then `"my-project"`. |

---

## Files written

| Path | Description |
|------|-------------|
| `README.md` | Minimal project README with project name |
| `.gitignore` | Common ignores: `.DS_Store`, `node_modules/`, `dist/`, `vendor/`, `*.log` |
| `LICENSE` | MIT License text using `project_name` |

Later generators (e.g. `typescript_base`) may append lines to `.gitignore` or rewrite `README.md` — this is safe because all writes are cooperative via the virtual filesystem.

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `README.md` exists | `file_exists` | File is present in the virtual state |
| `.gitignore` exists | `file_exists` | File is present |
| `LICENSE` exists | `file_exists` | File is present |

---

## Commands

No `PostGenerationCommands`. No `TestCommands`.

---

## Conflicts

None.
