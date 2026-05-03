# Flow: `test-flow`

A test flow to validate the dot-flow skill.

---

## Identity

| Field | Value |
|-------|-------|
| ID | `test-flow` |
| Title | Test Flow |
| File | `flows/test_flow.go` |
| Root question | `project_name` |

---

## Questions

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `project_name` | Text | "Project name" | Default: `""` |

---

## Question graph

```
project_name
  └── (end)
```

---

## Generator resolution

| Condition | Generators added |
|-----------|-----------------|
| Always | _(none — TODO: implement resolveTestFlowFlowGenerators)_ |

---

## Fixture examples

**Happy path** (`tools/test-flow/testdata/test_flow_full.json`):

```json
{
  "name": "test_flow_full",
  "flow_id": "test-flow",
  "answers": {
    "project_name": "test-value"
  },
  "expected_visited": [
    "project_name"
  ],
  "skip_post_commands": true,
  "skip_test_commands": true
}
```

---

## Source

`flows/test_flow.go`
