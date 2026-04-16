# Conflict Resolution

> Status: planned for v0.3. Design is settled; UX details will be decided with collaborators before implementation.

---

## The problem

When `dot add module` modifies a file the user has already changed, dot cannot safely overwrite it.

Example: you ran `dot init` with a `go-rest-api` generator. It created `main.go`. You then edited `main.go` to add your own startup logic. Three weeks later you run `dot add module postgres` — the Postgres generator wants to add a database connection to `main.go`. But your version of `main.go` is different from what dot generated.

dot must surface this conflict and let you resolve it. It must never silently overwrite your changes.

---

## The approach: git-style conflict markers

dot writes conflict markers directly into the file, in the same format as git merge conflicts. You resolve them in your editor exactly like a git conflict.

```
<<<<<<< dot (go-postgres generator)
db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
if err != nil {
    log.Fatal(err)
}
=======
// user's existing code that dot would have overwritten
>>>>>>> current
```

This is familiar, tooling-friendly (most editors highlight git conflicts), and explicit. You see exactly what dot wanted to add and exactly what was already there.

---

## How dot detects conflicts

1. `.dot/manifest.json` stores the SHA-256 hash of every generated file at the time `dot init` ran.

2. On `dot add module`, for each FileOp that targets an existing file:
   - Compute the current hash of the file on disk.
   - If the hash **matches** the manifest → file is unmodified → apply the op directly.
   - If the hash **differs** → user modified the file → write conflict markers instead.

3. `dot status` lists all files with unresolved conflict markers.

4. You resolve the conflicts in your editor (same as `git merge` conflicts).

5. `dot resolve` marks the conflicts as done and updates `manifest.json` with the new hashes.

---

## What is not decided yet

The exact conflict marker format (see [open-decisions.md](open-decisions.md) #2).

`dot resolve` UX details — how the developer signals that a conflict is resolved.

How to handle binary files (images, compiled artifacts) — cannot use text markers.

How to handle deleted files — the user deleted a file that dot wants to modify.

Renamed files — dot's manifest records the original path; renames break hash lookup.

These details will be worked out with collaborators before v0.3 implementation starts.
