# Flow: `plugin-template`

Scaffolds a complete, publishable DOT plugin repository. The output is a Go module with `plugin.json`, `plugin.go`, `go.mod`, `README.md`, `LICENSE`, and optionally sample injection and generator code. Intended for community plugin authors.

---

## Identity

| Field | Value |
|-------|-------|
| ID | `plugin-template` |
| Title | Plugin Repository Template |
| File | `flows/plugin_template.go` |
| Root question | `project_name` |

---

## Questions

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `project_name` | Text | "Plugin id (lowercase, no dots)" | Default: `"my-plugin"` |
| `module_path` | Text | "Go module path" | Default: `"github.com/your-org/my-plugin"` |
| `plugin_description` | Text | "One-line description" | Default: `"A DOT plugin"` |
| `plugin_author` | Text | "Author name (used in LICENSE)" | Default: `"Anonymous"` |
| `plugin_year` | Text | "Copyright year" | Default: `"2026"` |
| `plugin_include_injection` | Confirm | "Include a sample InsertAfter injection?" | Default: `true` |
| `plugin_include_generator` | Confirm | "Include a sample generator?" | Default: `true` |
| `confirm_generate` | Confirm | "Scaffold the plugin repo now?" | Default: `true` |

**Validation rules:**

| Question | Rule |
|----------|------|
| `project_name` | No `.` (reserved for namespacing), no spaces or path separators |
| `module_path` | Must contain `/` (basic host/owner/repo check) |

---

## Question graph

```
project_name
  └── module_path
        └── plugin_description
              └── plugin_author
                    └── plugin_year
                          └── plugin_include_injection
                                └── plugin_include_generator
                                      └── confirm_generate → (end)
```

This flow is fully linear — no branching. All questions are always visited.

---

## Generator resolution

| Condition | Generators added |
|-----------|-----------------|
| Always | `plugin_repo_skeleton` |

The resolver (`resolvePluginTemplateGenerators`) unconditionally returns one `Invocation{Name: "plugin_repo_skeleton"}`. The generator itself reads `plugin_include_injection` and `plugin_include_generator` to decide which code sections to emit.

---

## Fixture examples

**Full scaffold with injection + generator** (`tools/test-flow/testdata/plugin_template_full.json`):

```json
{
  "name": "plugin_template_full",
  "flow_id": "plugin-template",
  "answers": {
    "project_name": "my-plugin",
    "module_path": "github.com/example/my-plugin",
    "plugin_description": "Adds extra polish to scaffolded projects",
    "plugin_author": "Test Author",
    "plugin_year": "2026",
    "plugin_include_injection": true,
    "plugin_include_generator": true,
    "confirm_generate": true
  },
  "expected_visited": [
    "project_name", "module_path", "plugin_description",
    "plugin_author", "plugin_year",
    "plugin_include_injection", "plugin_include_generator",
    "confirm_generate"
  ],
  "skip_post_commands": false,
  "skip_test_commands": false
}
```

---

## Source

`flows/plugin_template.go`

## See also

- [docs/generators/plugin_repo_skeleton.md](../generators/plugin_repo_skeleton.md)
- [docs/authoring-plugins.md](../authoring-plugins.md) — full plugin authoring guide
