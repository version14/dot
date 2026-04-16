# Input Layers

---

## What an input layer is

An input layer is anything that produces a `Spec`.

The generator engine does not know or care which layer was used. It receives a `Spec` and returns `[]FileOp`. This means adding a new input layer — web UI, MCP, config file — requires zero changes to generators or the pipeline.

---

## CLI TUI (v0.1 — ships today)

The default input layer. `dot init` launches an interactive survey using the `huh` library. User choices map to `Spec` fields:

| TUI question | Spec field |
|---|---|
| Project name | `Spec.Project.Name` |
| Language | `Spec.Project.Language` |
| Project type | `Spec.Project.Type` |
| Modules | `Spec.Modules[].Name` |
| CI provider | `Spec.Config.CI` |
| Linter | `Spec.Config.Linter` |
| Deployment | `Spec.Config.Deployment` |

The survey is in `cmd/dot/cmd_init.go` (`surveySpec()`). It returns a fully populated `Spec` that is passed directly to `registry.ForSpec()`.

---

## dot.yaml — Project as Code (v0.5)

A declarative alternative to the TUI. Instead of answering questions, you write a `dot.yaml` file at the repo root and run `dot plan` + `dot apply`.

Same engine, different input. The `dot.yaml` parser reads the file and produces the same `Spec` the TUI would have produced.

```yaml
meta:
  name: my-platform
  type: monorepo

defaults:
  config:
    ci: github-actions
    linter: golangci-lint

apps:
  api:
    language: go
    type: rest-api
    modules: [rest-api, postgres, auth-jwt]
    config:
      deployment: docker
```

`dot plan` diffs the `dot.yaml` against `.dot/config.json` and shows what would change.
`dot apply` runs generators for the delta, piping through the conflict-aware pipeline.

See [roadmap/v0.5.md](../roadmap/v0.5.md) for the full design.

---

## Future layers

**MCP server (v1.1)** — A standalone server that lets AI agents scaffold and extend projects via the MCP protocol. The agent calls `dot_init` with a project description; the MCP layer parses it into a Spec and passes it to the engine. No new generators needed.

**Web Dashboard (v1.x)** — A team-friendly UI that produces the same Spec as the CLI. Useful for onboarding developers who prefer not to use the terminal, or for visualising a monorepo's structure.

Both are input layers. They share every line of generator code. The engine does not change.
