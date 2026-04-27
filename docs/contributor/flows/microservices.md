# Flow: `microservices`

Scaffold a project containing N independently-named service folders. The canonical example of a `LoopQuestion` in DOT: the user defines each service's name and port interactively, and the generator runs once per service.

---

## Identity

| Field | Value |
|-------|-------|
| ID | `microservices` |
| Title | Microservices Platform |
| File | `flows/microservices.go` |
| Root question | `project_name` |

---

## Questions

| ID | Type | Label | Options / Default |
|----|------|-------|-------------------|
| `project_name` | Text | "Project name" | Default: `"platform"` |
| `services` | Loop | "services" | _(see loop body)_ |
| `confirm_generate` | Confirm | "Generate the project now?" | Default: `true` |

### Loop body: `services`

Each iteration collects one service's details. The loop repeats until the user declines to add another service.

| ID | Type | Label | Default |
|----|------|-------|---------|
| `name` | Text | "Service name" | `"svc"` |
| `port` | Text | "Port" | `"3000"` |

> **Scoping note:** Inside the loop body the question IDs are `name` and `port` (not `service_name` / `service_port`). These are the keys the `service_writer` generator reads from its scoped `ctx.Answers`.

---

## Question graph

```
project_name
  └── services  (LoopQuestion)
        Body (repeated per iteration):
          name
            └── port → (end of iteration)
        ↓ (after all iterations)
        confirm_generate → (end)
```

---

## Generator resolution

| Condition | Generators added |
|-----------|-----------------|
| Always | `base_project` |
| Per loop iteration in `services` | `service_writer` (one invocation per iteration, each with its own `LoopFrame`) |

The resolver in `resolveMicroservicesGenerators` reads `s.Answers["services"]` as a `[]map[string]AnswerNode` (or `[]interface{}` from JSON) and emits one `Invocation{Name: "service_writer", LoopStack: [{...}]}` per entry.

---

## Fixture examples

**Three services** (`tools/test-flow/testdata/microservices_three.json`):

```json
{
  "name": "microservices_three",
  "flow_id": "microservices",
  "answers": {
    "project_name": "platform",
    "services": [
      {"name": "auth",    "port": 3001},
      {"name": "users",   "port": 3002},
      {"name": "billing", "port": 3003}
    ],
    "confirm_generate": true
  },
  "expected_visited": ["project_name", "services", "confirm_generate"]
}
```

Loop answers are an **array of objects**. Each object provides the answers for one iteration of the loop body. The length of the array determines how many times `service_writer` runs.

---

## Source

`flows/microservices.go`

## See also

- [docs/generators/service_writer.md](../generators/service_writer.md)
- [docs/authoring-flows.md](../authoring-flows.md#loops) — loop authoring guide
