# Repository Structure

---

## Current layout (v0.1)

```
dot/
├── cmd/dot/                  ← CLI entry point
│   ├── main.go               ← thin: run(os.Args[1:]) → os.Exit
│   ├── build.go              ← buildVersion(), buildRegistry()
│   ├── styles.go             ← lipgloss styles and ASCII banner
│   ├── cmd_init.go           ← dot init (huh TUI → Spec → generators)
│   ├── cmd_new.go            ← dot new <type> <name>
│   ├── cmd_help.go           ← dot help (reads .dot/config.json)
│   └── cmd_selfupdate.go     ← dot self-update
│
├── internal/
│   ├── spec/                 ← Spec, ProjectSpec, CoreConfig, ModuleSpec
│   ├── generator/            ← Generator interface, Registry, FileOp, CommandDef
│   ├── project/              ← Context, Load, Save (.dot/config.json + manifest.json)
│   └── pipeline/             ← FileOp collect → resolve → write
│
├── generators/
│   ├── go/                   ← official Go generators (package gogen)
│   └── common/               ← language-agnostic generators (CI, Docker — v0.2+)
│
├── templates/                ← files embedded via go:embed
│
├── .goreleaser.yaml          ← multi-platform release config
├── install.sh                ← curl installer
├── uninstall.sh              ← curl uninstaller
├── makefile                  ← dev commands
└── go.mod
```

**`cmd/dot/` is thin.** It parses `os.Args`, calls into `internal/`, and prints. No business logic lives here.

**`internal/` is the engine.** These packages have no knowledge of the CLI. They are independently testable and will become their own module when a second consumer (Dashboard, MCP) needs them.

**`generators/` is the content.** Official generators live here. Community generators implement the same interface from anywhere.

---

## The one rule

Nothing in `internal/` ever imports from `cmd/`.

`cmd/dot` depends on `internal/`. Never the reverse. Go's import system enforces this — the build will fail if a circular import sneaks in. This boundary is what makes the progressive split plan feasible.

---

## Progressive split plan

dot starts as a single repo. It will split when there is a real reason to split — not before.

**Phase 1 — now (v0.1 through v0.5)**

Everything in `version14/dot`. One module, one repo. The `internal/` packages are the stable core, but they are not exposed externally yet.

**Phase 2 — extract `dot-core` (triggered by a second consumer)**

When the Dashboard or MCP server needs the generator engine, `internal/` becomes `version14/dot-core`. `cmd/dot` becomes a thin consumer of `dot-core`. Other consumers import `dot-core` directly.

Trigger: a second executable needs the engine. Not before.

**Phase 3 — extract `dot-std` (triggered by community generator authors)**

When enough community generators exist that they need a stable, versioned interface to import (rather than copying types), `dot-std` becomes a separate module exposing just `Generator`, `Spec`, `FileOp`, and `CommandDef`.

Trigger: community generator authors asking for a stable import path.

**Phase 4 — extract `dot-registry` (triggered by community volume)**

A separate service for publishing, discovering, and installing community generators. Like pkg.go.dev for generators.

Trigger: enough community generators to justify a registry service.

Each split is driven by a real need, not speculation. The right amount of indirection is the minimum that solves the actual problem.
