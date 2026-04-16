# FileOp Pipeline

Implementation: `internal/pipeline/pipeline.go`, `internal/pipeline/patch.go`

---

## What the pipeline does

The pipeline is the only thing in dot that writes to disk.

Generators produce descriptions of what to write (`[]FileOp`). They do not call `os.WriteFile` directly. Everything flows through the pipeline, which gives dot its central guarantee: if anything fails, nothing is written.

Three phases: **collect → resolve → write**.

---

## Phase 1: Collect

`dot init` calls `generator.Apply(spec)` on every matched generator. Each returns `[]FileOp`. The pipeline receives all ops from all generators as a flat slice.

Nothing is written yet. This phase is pure computation.

---

## Phase 2: Resolve conflicts

Ops are sorted by `Priority` (descending), then `Generator` name (alphabetical tiebreak).

**Create and Template ops on the same path:**
- Highest priority wins. The lower-priority op is silently skipped.
- Two ops at the same priority on the same path → pipeline aborts with a descriptive error naming both generators.

**Append and Patch ops:**
- All are applied in priority order. No conflict. Multiple generators can append to the same file.

After sorting, the pipeline scans for Create/Template conflicts before any writes begin. A conflict at this stage means something is wrong with the generator registrations, not with user input — it is a programming error, not a user error.

---

## Phase 3: Write atomically

The pipeline builds all file contents in memory (`buildWrites`). Only when every op has been applied without error does it call `flushWrites`, which writes files to disk.

`flushWrites` creates parent directories as needed and writes each file. If writing fails mid-way, the files written before the failure will have landed on disk — there is no rollback at the OS level. But by the time `flushWrites` runs, all validation is complete. The only failures in this phase are OS-level errors (disk full, permissions) which are unrecoverable regardless.

---

## FileOp kinds

| Kind | Behavior | Conflict rule |
|---|---|---|
| `Create` | Write a new file with literal content | Priority wins; same-priority conflict aborts |
| `Template` | Render a Go `text/template` then write | Same as Create |
| `Append` | Add content to the end of a file | No conflict — all ops applied in order |
| `Patch` | Insert content at a named anchor point | No conflict — all ops applied in order |

Full reference: [fileop-reference.md](../generators/fileop-reference.md).
Anchor reference: [patch-strategies.md](../generators/patch-strategies.md).
