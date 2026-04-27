<!-- ─────────────────────────────────────────────────────────────────────────
  FLOW DOC TEMPLATE
  ─────────────────────────────────────────────────────────────────────────
  How to use:
    1. Copy this file to docs/flows/<flow-id>.md  (flow-id = FlowDef.ID)
    2. Replace every <!-- TODO --> comment and every _placeholder_ value
    3. Delete this header block before committing

  Add the new file to the table in docs/README.md (flows section).
  Also add the flow to the table in docs/authoring-flows.md (Built-in flows).
  Register it in flows/registry.go Default().
  ───────────────────────────────────────────────────────────────────────── -->

# Flow: `_flow-id_`

<!-- TODO: one-paragraph description of what this flow scaffolds and who it's for. -->

---

## Identity

| Field | Value |
|-------|-------|
| ID | `_flow-id_` |
| Title | _Flow Title_ |
| File | `flows/_flow_id_.go` |
| Root question | `_root_question_id_` |

---

## Questions

<!-- TODO: one row per question node. Include all IDs, including injected ones
     that plugins add (mark them with the plugin ID).
     Loop body questions are listed under a "Loop body" sub-header. -->

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `_id_` | Text | "_label_" | Default: `"_default_"` |
| `_id_` | Option | "_label_" | `_opt1_`, `_opt2_` |
| `_id_` | Confirm | "_label_" | Default: `true` |

<!-- TODO: if the flow has a LoopQuestion, add:

### Loop body: `_loop_id_`

| ID | Type | Label | Default |
|----|------|-------|---------|
| `_id_` | Text | "_label_" | `"_default_"` |
-->

<!-- TODO: if active plugins inject questions, add:

### Plugin-injected questions

| ID | Plugin | Target | Type |
|----|--------|--------|------|
| `_plugin_id_._question_` | `_plugin_id_` | InsertAfter `_target_` | Confirm |
-->

---

## Question graph

<!-- TODO: draw the question graph as an ASCII diagram.
     Show every branch; converging paths can share a node label. -->

```
_root_question_id_
  └── _next_question_id_
        ├── [option A] → _branch_a_
        │                 └── _shared_question_
        └── [option B] → _shared_question_
                           └── confirm_generate
                                 └── (end)
```

---

## Generator resolution

<!-- TODO: table mapping answer conditions to generator invocations.
     "Always" means unconditional. Loop rows emit one invocation per iteration. -->

| Condition | Generators added |
|-----------|-----------------|
| Always | `base_project` |
| `_answer_key_` = `_value_` | `_generator_name_` |
| _(plugin)_ `_plugin_key_` = `true` | `_plugin_id_._gen_` |

---

## Fixture examples

<!-- TODO: at least one fixture snippet showing the happy path.
     Reference the full fixture files in testdata/. -->

**Happy path** (`tools/test-flow/testdata/_fixture_name_.json`):

```json
{
  "name": "_fixture_name_",
  "flow_id": "_flow-id_",
  "answers": {
    "_root_question_id_": "_value_",
    "confirm_generate": true
  },
  "expected_visited": [
    "_root_question_id_",
    "confirm_generate"
  ],
  "skip_post_commands": true,
  "skip_test_commands": true
}
```

<!-- TODO: add a second snippet for a notable branch (e.g. with optional features off). -->

---

## Source

`flows/_flow_id_.go`

## See also

<!-- TODO: link to related generator and plugin docs. Delete if unused. -->

- [docs/generators/_generator_.md](../generators/_generator_.md)
- [docs/plugins/_plugin_.md](../plugins/_plugin_.md)
