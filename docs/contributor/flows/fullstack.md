# Flow: `fullstack`

Full-stack application scaffold. Targets TypeScript frontends (with optional React) and an optional Go backend. Shorter than the monorepo flow — no monorepo style question — and always includes TypeScript as the base.

---

## Identity

| Field | Value |
|-------|-------|
| ID | `fullstack` |
| Title | Fullstack Application |
| File | `flows/fullstack.go` |
| Root question | `project_name` |

---

## Questions

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `project_name` | Text | "Project name" | Default: `"fullstack-app"` |
| `stack` | Option | "Stack" | `typescript`, `polyglot` |
| `ui_library` | Option | "UI library" | `react`, `none` |
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
  └── stack  (typescript | polyglot)
        └── ui_library  (react | none)
              └── use_biome
                    └── [biome_extras] biome_extras.strict_mode
                          └── confirm_generate → (end)
```

Both `stack` options lead to `ui_library`; both `ui_library` options lead to `use_biome`.

---

## Generator resolution

| Condition | Generators added |
|-----------|-----------------|
| Always | `base_project`, `typescript_base` |
| `ui_library` = `react` | `react_app` |
| `use_biome` = `true` | `biome_config` |
| _(biome_extras plugin)_ `use_biome` = `true` **and** `biome_extras.strict_mode` = `true` | `biome_extras.strict_writer` |

---

## Fixture examples

**React + Biome** (`tools/test-flow/testdata/fullstack_react.json`):

```json
{
  "name": "fullstack_react",
  "flow_id": "fullstack",
  "answers": {
    "project_name": "fs-app",
    "stack": "polyglot",
    "ui_library": "react",
    "use_biome": true,
    "biome_extras.strict_mode": false,
    "confirm_generate": true
  },
  "expected_visited": [
    "project_name", "stack", "ui_library",
    "use_biome", "biome_extras.strict_mode", "confirm_generate"
  ]
}
```

**No UI** (`tools/test-flow/testdata/fullstack_no_ui.json`):

```json
{
  "name": "fullstack_no_ui",
  "flow_id": "fullstack",
  "answers": {
    "project_name": "api-only",
    "stack": "typescript",
    "ui_library": "none",
    "use_biome": false,
    "confirm_generate": true
  },
  "expected_visited": ["project_name", "stack", "ui_library", "use_biome", "confirm_generate"]
}
```

---

## Source

`flows/fullstack.go`

## See also

- [docs/generators/typescript_base.md](../generators/typescript_base.md)
- [docs/generators/react_app.md](../generators/react_app.md)
- [docs/generators/biome_config.md](../generators/biome_config.md)
- [docs/plugins/biome_extras.md](../plugins/biome_extras.md)
