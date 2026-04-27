# Plugin: `example` (example-plugin)

Reference implementation for community plugin authors. Demonstrates two injection kinds (`InsertAfter` and `AddOption`) in a single plugin, alongside a real generator. Not included in any production flow — exists purely as a learning resource.

**Source:** `examples/example-plugin/plugin.go`

---

## Identity

| Field | Value |
|-------|-------|
| Plugin ID | `example` |
| Package | `examples/example-plugin` |
| Type | Example (not imported by the main binary) |

---

## Injections

### `InsertAfter` on `use_biome`

| Injected question ID | Type | Label | Default |
|---------------------|------|-------|---------|
| `example.add_editorconfig` | `ConfirmQuestion` | "Add .editorconfig for cross-IDE consistency?" | `true` |

### `AddOption` on `stack`

Appends a new option to the `stack` `OptionQuestion`:

| Option value | Label | Next edge |
|-------------|-------|-----------|
| `example.vscode_only` | "VSCode workspace only" | `End` (short-circuits the rest of the flow) |

This illustrates how a plugin can add a complete new branch to an existing question — selecting this option stops the flow at that point.

---

## Generators contributed

| Generator name | Version | Description |
|----------------|---------|-------------|
| `example.editorconfig_writer` | `0.1.0` | Writes a sensible `.editorconfig` at project root |

### `example.editorconfig_writer`

**Dependencies:** `base_project`.

**Files written:**

| Path | Description |
|------|-------------|
| `.editorconfig` | Root-level `.editorconfig` with LF line endings, UTF-8, 2-space indent, trim whitespace |

**Validators:** Checks that `.editorconfig` exists after generation.

---

## ResolveExtras logic

`example.editorconfig_writer` runs only when `example.add_editorconfig` is `true`. If the answer is `false` or missing, `ResolveExtras` returns `nil`.

---

## How to use as a reference

When writing a new plugin, read this file alongside the source at `examples/example-plugin/plugin.go`. It covers:

- The full `Provider` interface implementation.
- Combining two injection kinds in `Injections()`.
- Gating a generator on an injected answer in `ResolveExtras()`.
- The naming convention for all contributed IDs.

To scaffold a new plugin repository with a similar structure:

```bash
dot scaffold plugin-template
```
