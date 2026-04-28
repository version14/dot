# Generator: `plugin_repo_skeleton`

Scaffolds a complete, publishable DOT plugin repository. Used exclusively by the `plugin-template` flow. Unlike most generators it does **not** depend on `base_project` — plugin repos have their own README shape and no need for the generic project boilerplate.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `plugin_repo_skeleton` |
| Version | `0.1.6` |
| Package | `generators/plugin_repo_skeleton` |

---

## Dependencies

None.

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `project_name` | string | **Yes** | Used as the plugin ID. Must not contain `.`. Also used as the Go package name (hyphens and underscores stripped). |
| `module_path` | string | **Yes** | Go module path (e.g. `github.com/you/my-plugin`). Written into `go.mod` and the install URL in `README.md`. |
| `plugin_description` | string | No | One-line description. Default: `"A DOT plugin"`. |
| `plugin_author` | string | No | Author name for `LICENSE`. Default: `"Anonymous"`. |
| `plugin_year` | string | No | Copyright year for `LICENSE`. Default: `"2026"`. |
| `plugin_include_injection` | bool | No | When `true`, generates a sample `InsertAfter` injection in `plugin.go`. |
| `plugin_include_generator` | bool | No | When `true`, generates a sample generator with manifest in `plugin.go`. |

---

## Files written

| Path | Description |
|------|-------------|
| `go.mod` | Module declaration pinned to `github.com/version14/dot v0.1.6` |
| `plugin.json` | Identity file: `id`, `version`, `description`, `entry_point` |
| `plugin.go` | `Provider` implementation + `init()` + optional injection + optional generator |
| `README.md` | Install instructions + description of what the plugin does |
| `LICENSE` | MIT License using `plugin_author` and `plugin_year` |
| `.gitignore` | Common ignores for Go build artifacts |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `plugin.go` exists | `file_exists` | — |
| `plugin.json` exists | `file_exists` | — |
| `go.mod` exists | `file_exists` | — |
| `id` key in `plugin.json` | `json_key_exists` | Plugin identity is present |
| `version` key in `plugin.json` | `json_key_exists` | Version field is present |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `go mod tidy` | project root | Resolves and pins dependencies after scaffold |
| `git init` | project root | Initializes a new git repository |

## Test commands

No test commands — plugin repositories are tested by their own CI after publishing.

---

## Conflicts

None.
