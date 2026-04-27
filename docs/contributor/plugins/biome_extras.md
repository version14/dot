# Plugin: `biome_extras`

Built-in demo plugin. Adds an optional Biome strict-mode question and a generator that promotes selected lint rules to errors. It is the canonical in-tree example of the complete plugin contract: namespacing, `InsertAfter` injection, generator registration, and `ResolveExtras` gating.

---

## Identity

| Field | Value |
|-------|-------|
| Plugin ID | `biome_extras` |
| Package | `plugins/biome_extras` |
| Type | Built-in (imported in `cmd/dot/main.go`) |

---

## Injections

### `InsertAfter` on `use_biome`

After the flow's `use_biome` confirm question, the plugin inserts:

| Injected question ID | Type | Label | Default |
|---------------------|------|-------|---------|
| `biome_extras.strict_mode` | `ConfirmQuestion` | "Enable Biome strict mode? (catches more issues but is noisier)" | `false` |

**When is this visible?** Always, when the flow includes a `use_biome` question and the user is shown that question. If `use_biome` is not in the flow, the injection silently has no effect.

**Flows affected:**
- `monorepo` (via `use_biome`)
- `fullstack` (via `use_biome`)

---

## Generators contributed

| Generator name | Version | Description |
|----------------|---------|-------------|
| `biome_extras.strict_writer` | `0.1.0` | Promotes selected Biome lint rules to errors |

### `biome_extras.strict_writer`

**Dependencies:** `biome_config` (must run first — overwrites `biome.json` which this generator updates).

**Files modified:** `biome.json` — merges the following keys under `linter.rules`:

```json
{
  "linter": {
    "rules": {
      "style": {
        "useImportType": "error",
        "noNonNullAssertion": "error"
      },
      "suspicious": {
        "noExplicitAny": "error"
      }
    }
  }
}
```

**Validators:** Checks that `linter.rules.style` key exists in `biome.json` after generation.

---

## ResolveExtras logic

`biome_extras.strict_writer` only runs when **both** conditions are true:

1. The core flow answer `use_biome` is `true`.
2. The injected answer `biome_extras.strict_mode` is `true`.

If either is `false` or missing, `ResolveExtras` returns `nil` and the generator is not included in the invocation set.

---

## Answers added to spec

| Key | Type | Added by |
|-----|------|----------|
| `biome_extras.strict_mode` | bool | InsertAfter injection on `use_biome` |

---

## Fixture requirement

Any fixture for a flow that includes `use_biome` must also provide `biome_extras.strict_mode`:

```json
{
  "answers": {
    "use_biome": true,
    "biome_extras.strict_mode": false
  }
}
```

Omitting it causes: `test-flow: no scripted answer for question "biome_extras.strict_mode"`.

---

## Source

`plugins/biome_extras/plugin.go`
