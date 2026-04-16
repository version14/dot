# dot — Product Brief

**Last updated:** 2026-04-16
**Status:** Living document — update when the vision shifts, not when code ships.

---

## The problem

Starting a project has three broken paths:

1. **Opinionated starters** — fast, but you spend hours removing what you didn't ask for.
2. **Template repos** — someone else's decisions, and you spend 30 minutes untangling them.
3. **From scratch** — full control, but 200 lines of boilerplate before a single line of real code.

Extending an existing project is worse. You add a module, something breaks. You add Docker, the CI config is wrong. You add auth, the import graph is a mess. Every addition is manual surgery.

The gap: no tool lets you describe exactly what you want and extends your project safely over time.

---

## What dot is

dot is a CLI for developers who know what they want.

You describe your project once. dot generates a clean base. Later, when you need to add something, dot extends the project without breaking what's already there.

It is not a framework. It is not a template runner. It is a project companion — it knows your project's structure and uses that knowledge to keep things consistent as the project grows.

---

## Who it is for

Developers with strong opinions about how their projects should be structured. People who know they want Go + golangci-lint + GitHub Actions + Docker, and don't want to wire it all together by hand every time.

Initially: the version14 team. Eventually: any developer who works this way.

---

## The core loop

```
dot init
  └── TUI survey (language, type, modules, config)
        └── Spec (a typed description of the project)
              └── Generator engine
                    └── FileOp pipeline → project on disk + .dot/config.json

dot new <type> <name>
  └── Looks up .dot/config.json
        └── Finds the right generator
              └── FileOp pipeline → new artifact injected into the existing project
```

After `dot init`, the project has a `.dot/config.json`. This file is what makes `dot new` safe — dot knows what it generated, where it put things, and how to extend without colliding.

---

## What "consistent" means

The central promise of dot is **safe extension**.

When you add a module after creation, the existing project must not break. If dot cannot inject safely (conflicting files, unsupported import form, unknown structure), it must:

1. Stop before writing anything
2. Tell the developer exactly what the conflict is
3. Provide instructions to resolve it manually

dot never silently corrupts a project. A partial write that leaves the project in a broken state is the worst possible outcome. It is better to do nothing and explain why.

---

## What dot does NOT do

- **Business logic.** dot generates structure, not behavior. A CRUD scaffold is the limit — the actual query logic is yours.
- **Domain-specific patterns.** dot does not know what a "UserController" should do in your app. It knows where to put the file and how to wire it in.
- **Highly specific one-off things.** If you need something that only makes sense in your product, dot is not the right tool. dot handles the 80% that is the same across projects.
- **Runtime operations.** dot does not deploy, run, or monitor your project.

---

## The generator model

Everything dot generates comes from a generator. A generator is a Go struct that implements a simple interface:

- It declares what language and modules it handles
- Given a Spec, it returns a list of file operations (create, append, patch)
- It registers commands that become available after `dot init`

The official generators live in this repo. Community generators can implement the same interface and be loaded locally. A public registry is a future concern.

This model means: adding support for a new language or framework never changes the core engine. You write a generator.

---

## v0.1 scope (now)

What exists and works today:

| Command | What it does |
|---------|-------------|
| `dot init` | TUI survey → generates a Go REST API project |
| `dot new route <name>` | Adds a route to an existing Go REST API project |
| `dot new handler <name>` | Adds a handler |
| `dot help` | Lists available commands for the current project |
| `dot version` | Prints the current version |
| `dot self-update` | Updates dot to the latest release |

One generator exists: `GoRestAPIGenerator`. It is the proof that the engine works.

---

## What comes next (not committed, not scheduled)

These are directions, not promises:

- **More Go generators** — Postgres, Redis, Docker, GitHub Actions
- **Python generators** — FastAPI, CLI tools
- **dot.yaml (Project as Code)** — describe a full project in a YAML file instead of the TUI
- **Community generator loading** — load generators from local paths
- **MCP server** — let AI agents scaffold projects via the MCP protocol

None of these change the core engine. They are all additive.

---

## The limits dot must respect

**dot is not AI.** It generates from explicit, deterministic rules. What you get is predictable and reproducible. This is a feature, not a limitation — you can read the generator code and know exactly what will be written.

**dot is not a build tool.** It does not compile, run, or test your project. It generates the scaffolding. The developer takes over from there.

**dot is not a package manager.** It does not manage dependencies at runtime. It writes the initial `go.mod` or `package.json`, but updating dependencies is your job after that.

**dot cannot handle every project.** If your project has been heavily customized after `dot init`, some `dot new` operations may produce conflicts. That is expected. dot documents the conflict and stops.

---

## The one rule

> dot must leave the project in a better state than it found it, or leave it exactly as it found it. Never worse.

If an operation cannot complete safely, it does nothing. This is the invariant everything else is built around.
