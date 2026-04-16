# Architecture Overview

---

## The big picture

dot has three layers. Every user action flows through all three.

```
┌─────────────────────────────────────────┐
│           Input layers                  │
│  CLI TUI  │  dot.yaml  │  MCP (future)  │
└─────────────────────────────────────────┘
                    │
                    ▼ Spec
┌─────────────────────────────────────────┐
│           Generator engine              │
│   Registry.ForSpec → Apply → []FileOp   │
└─────────────────────────────────────────┘
                    │
                    ▼ []FileOp
┌─────────────────────────────────────────┐
│           FileOp pipeline               │
│   collect → resolve conflicts → write   │
└─────────────────────────────────────────┘
                    │
                    ▼
             project on disk
           + .dot/config.json
           + .dot/manifest.json
```

**Input layers** produce a `Spec`. They know nothing about generators.

**Generator engine** consumes a `Spec`, matches generators, and returns `[]FileOp`. It knows nothing about how the Spec was produced.

**FileOp pipeline** executes the ops. It knows nothing about generators or input layers — it just writes files.

This separation is what makes dot extensible. A new input layer (web UI, MCP server) needs zero changes to the engine or generators. A new generator needs zero changes to the CLI.

---

## The central invariant

> dot must leave the project in a better state than it found it, or leave it exactly as it found it. Never worse.

In practice: every file operation is assembled in memory first. Nothing touches disk until every op has been validated. If anything fails at any point in collect or resolve, the run aborts and the project on disk is unchanged.

A partial write that leaves the project in a broken state is worse than a clear failure message. The pipeline is designed around this.

---

## What dot is NOT

**Not a compiler.** dot generates source files. It does not build them. Run `go build` yourself.

**Not a build tool.** `make`, `gradle`, `cargo` — those are yours. dot writes the Makefile, not what's in it.

**Not a package manager.** dot writes your initial `go.mod` or `package.json`. Updating dependencies after that is your job.

**Not AI.** Everything dot generates comes from deterministic generator code. Same Spec, same files, every time. This is a feature — you can read the generator code and know exactly what will be written.

**Not a runtime.** dot does nothing at runtime. It scaffolds. The developer takes over after that.

---

## Key design decisions

**Spec as the single contract.** Input layers and generators never talk to each other. The Spec struct is the only handshake point. This keeps the two sides independently testable and replaceable.

**Generator interface as the extension point.** Adding support for a new language or framework means writing a new struct that implements `Generator`. You do not touch the engine. You do not fork dot. You register a generator.

**FileOp pipeline as the single write path.** Generators return descriptions of what to write (`[]FileOp`). They do not call `os.WriteFile`. The pipeline is the only thing that writes to disk, and it only does so atomically.

**`.dot/` committed to git.** `.dot/config.json` is not a build artifact — it is project state. `dot new` cannot run without it. `dot add module` cannot detect conflicts without `manifest.json`. Both files belong in git alongside your source code.
