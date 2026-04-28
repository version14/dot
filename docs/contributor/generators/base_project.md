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

None.

---

## Files written

This generator fetches the content from the [github-template](https://github.com/mathieusouflis/github-template.git) repository and writes the following files to the project's root:

| Path | Description |
|------|-------------|
| `README.md` | Project README template. |
| `.gitignore` | A standard `.gitignore` file for Go projects. |
| `LICENSE` | MIT License template. |

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
