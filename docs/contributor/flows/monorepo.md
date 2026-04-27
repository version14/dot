# Flow: `monorepo`

The general-purpose project wizard. Walks the user from a project name through monorepo style, language stack, and optional tooling (React, Biome). Suitable for single apps, Turborepo workspaces, or Nx monorepos in TypeScript, Go, or both.

---

## Identity

| Field | Value |
|-------|-------|
| ID | `monorepo` |
| Title | Monorepo / Project Wizard |
| File | `flows/monorepo.go` |
| Root question | `project_name` |

---

## Questions

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `project_name` | Text | "Project name" | Default: `"my-project"` |
| `monorepo_type` | Option | "Monorepo style" | `single`, `turborepo`, `nx` |
| `stack` | Option | "Primary language stack" | `typescript`, `go`, `polyglot` |
| `use_react` | Confirm | "Set up a React app?" | Default: `true` |
| `use_biome` | Confirm | "Add Biome (lint + format)?" | Default: `true` |
| `confirm_generate` | Confirm | "Generate the project now?" | Default: `true` |

### Plugin-injected questions

| ID | Plugin | Target | Type |
|----|--------|--------|------|
| `biome_extras.strict_mode` | `biome_extras` | InsertAfter `use_biome` | Confirm |

---

## Question graph

```
project_name
  └── monorepo_type  (single | turborepo | nx)
        └── stack
              ├── [typescript | polyglot] → use_react
              │                               └── use_biome
              │                                     └── [biome_extras] biome_extras.strict_mode
              │                                           └── confirm_generate → (end)
              └── [go] → confirm_generate → (end)
```

`use_react` and `use_biome` always lead to `confirm_generate` regardless of their value — the branches converge.

---

## Generator resolution

| Condition | Generators added |
|-----------|-----------------|
| Always | `base_project` |
| `stack` = `typescript` or `polyglot` | `typescript_base` |
| `stack` = `typescript` or `polyglot` **and** `use_react` = `true` | `react_app` |
| `stack` = `typescript` or `polyglot` **and** `use_biome` = `true` | `biome_config` |
| _(biome_extras plugin)_ `use_biome` = `true` **and** `biome_extras.strict_mode` = `true` | `biome_extras.strict_writer` |

---

## Fixture examples

**TypeScript + React + Biome** (`tools/test-flow/testdata/turborepo_ts_react.json`):

```json
{
  "name": "turborepo_ts_react",
  "flow_id": "monorepo",
  "answers": {
    "project_name": "demo-app",
    "monorepo_type": "turborepo",
    "stack": "typescript",
    "use_react": true,
    "use_biome": true,
    "biome_extras.strict_mode": false,
    "confirm_generate": true
  },
  "expected_visited": [
    "project_name", "monorepo_type", "stack",
    "use_react", "use_biome", "biome_extras.strict_mode", "confirm_generate"
  ]
}
```

**Go only, no tooling** (`tools/test-flow/testdata/single_go.json`):

```json
{
  "name": "single_go",
  "flow_id": "monorepo",
  "answers": {
    "project_name": "go-svc",
    "monorepo_type": "single",
    "stack": "go",
    "confirm_generate": true
  },
  "expected_visited": ["project_name", "monorepo_type", "stack", "confirm_generate"]
}
```

---

## Source

`flows/monorepo.go`

## See also

- [docs/generators/typescript_base.md](../generators/typescript_base.md)
- [docs/generators/react_app.md](../generators/react_app.md)
- [docs/generators/biome_config.md](../generators/biome_config.md)
- [docs/plugins/biome_extras.md](../plugins/biome_extras.md)
