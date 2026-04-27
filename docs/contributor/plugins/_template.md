<!-- ─────────────────────────────────────────────────────────────────────────
  PLUGIN DOC TEMPLATE
  ─────────────────────────────────────────────────────────────────────────
  How to use:
    1. Copy this file to docs/plugins/<id>.md  (id = Provider.ID(), no dots)
    2. Replace every <!-- TODO --> comment and every _placeholder_ value
    3. Delete this header block before committing

  Add the new file to the table in docs/README.md (plugins section).
  Also update docs/authoring-plugins.md "Reference implementations" table
  if this is an in-tree or example plugin.
  ───────────────────────────────────────────────────────────────────────── -->

# Plugin: `_plugin_id_`

<!-- TODO: one sentence — what does this plugin add and when would someone use it?
     If it is a reference/example plugin, note that explicitly. -->

---

## Identity

| Field | Value |
|-------|-------|
| Plugin ID | `_plugin_id_` |
| Package | `plugins/_plugin_id_` |
| Type | `Built-in` / `Example` / `Community` |

<!-- Type:
  "Built-in"  — lives in plugins/, imported in cmd/dot/main.go
  "Example"   — lives in examples/, for reference only
  "Community" — published as a separate git repo
-->

---

## Injections

<!-- TODO: one subsection per Injection in Provider.Injections().
     If this plugin registers no injections, replace this section with: None. -->

### `_InjectKind_` on `_target_question_id_`

<!-- TODO: describe what is injected and why. -->

| Injected question ID | Type | Label | Default |
|---------------------|------|-------|---------|
| `_plugin_id_._question_` | `_ConfirmQuestion / TextQuestion / OptionQuestion_` | "_label_" | `_default_` |

**Flows affected:** <!-- TODO: list which built-in flows contain the target question. -->

---

## Generators contributed

<!-- TODO: one row per Entry in Provider.Generators().
     If this plugin contributes no generators, replace with: None. -->

| Generator name | Version | Description |
|----------------|---------|-------------|
| `_plugin_id_._gen_` | `_0.1.0_` | _one-line description_ |

### `_plugin_id_._gen_`

<!-- TODO: fill in the sub-details for each contributed generator. -->

**Dependencies:** <!-- list DependsOn generators -->

**Files written:**

| Path | Description |
|------|-------------|
| `_path_` | _description_ |

**Validators:** <!-- list Checks; or write "None." -->

---

## ResolveExtras logic

<!-- TODO: explain the conditions under which each generator invocation is returned.
     Be explicit about which answer keys are checked and what values trigger it. -->

`_plugin_id_._gen_` runs only when:

1. Answer `_key_` equals `_value_`.
2. _(add more conditions as needed)_

If any condition is false, `ResolveExtras` returns `nil`.

---

## Answers added to spec

<!-- TODO: list every question ID this plugin injects (from Injections above).
     These are the keys that will appear in .dot/spec.json. -->

| Key | Type | Added by |
|-----|------|----------|
| `_plugin_id_._question_` | _bool / string_ | InsertAfter injection on `_target_` |

---

## Fixture requirement

<!-- TODO: show the minimum snippet every test-flow fixture must include
     when running a flow that this plugin is active on.
     Delete this section if the plugin injects no questions. -->

Any fixture for a flow that includes `_target_question_id_` must also provide:

```json
{
  "answers": {
    "_target_question_id_": _value_,
    "_plugin_id_._question_": _value_
  }
}
```

Omitting it causes: `test-flow: no scripted answer for question "_plugin_id_._question_"`.

---

## Source

`_plugins/plugin_id_/plugin.go_`
