# Open Decisions

Unresolved questions that block specific features. Each decision must be made with collaborators before the affected feature is implemented.

---

## 1. Community generator loading mechanism

**Blocks:** v0.2 local custom generators

Go has no simple cross-platform plugin system. Three options:

| Option | How it works | Pros | Cons |
|--------|-------------|------|------|
| **In-process** | Community generators are Go modules imported at compile time. Users compile a custom dot binary. | Simple, no IPC | Requires recompiling dot to add generators |
| **Subprocess / RPC** | Community generators are separate binaries. dot spawns them, communicates via stdin/stdout JSON. | Flexible, any language | More complex, FileOpKind must be string-typed (already done) |
| **Embedded registry** | dot fetches generator binaries from a registry URL, runs as subprocesses. Like Terraform providers. | Best UX | Most infrastructure to build |

**Note:** The Generator interface is compatible with all three options. The choice affects distribution and loading, not the interface contract. `FileOpKind` is already string-typed in anticipation of the subprocess option.

**Decision needed before:** v0.2 local custom generator work begins.

---

## 2. dot resolve UX and conflict marker format

**Blocks:** v0.3 `dot add module` and conflict resolution

The conflict strategy is settled (git-style markers). What is not decided:

- Exact format of the conflict markers (which metadata to include in the header lines)
- What `dot resolve` does step by step
- How to handle binary files (images, compiled assets) — cannot use text markers
- How to handle files deleted by the user after `dot init`
- How to handle renamed files — the manifest stores the original path

**Decision needed before:** v0.3 conflict resolution implementation begins.

---

## 3. dot plan diff algorithm

**Blocks:** v0.5 `dot plan` command

The output format is clear ("will add: user-service/redis") but the algorithm is not.
A complete diff must handle:

- Added apps (new key in `dot.yaml` apps)
- Removed apps (key in `.dot/config.json` but not in `dot.yaml`)
- Added modules per app
- Removed modules per app
- Changed `CoreConfig` fields
- Changed `Extensions`

**Decision needed before:** v0.5 `dot plan` implementation begins.
