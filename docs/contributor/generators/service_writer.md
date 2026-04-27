# Generator: `service_writer`

Writes one Go/TypeScript service skeleton per loop iteration. Used exclusively by the `microservices` flow via a `LoopQuestion`. The generator is invoked once per service the user defines; each invocation receives a different scoped answer set.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `service_writer` |
| Version | `0.1.0` |
| Package | `generators/service_writer` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `base_project` | Must run first to establish the project root |

---

## Answers consumed

These keys come from the loop frame (scoped per iteration), not from global answers:

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `name` | string | **Yes** | Service name, used as the subdirectory (`services/<name>/`) and in log messages. Returns an error if empty. |
| `port` | int or float64 | No | HTTP listen port. Defaults to `3000` if missing. JSON numbers unmarshal as `float64`, so both `int` and `float64` are handled. |

> **Note:** In the `microservices` flow's LoopQuestion body, the keys are `service_name` and `service_port` — but the loop frame in the fixture uses `"name"` and `"port"` as the keys the generator reads. Verify the flow's `LoopFrame.Answers` mapping when debugging scoping issues.

---

## Files written

Files are written under `services/<name>/` — each service gets its own subdirectory:

| Path | Description |
|------|-------------|
| `services/<name>/src/main.ts` | Minimal Node.js HTTP server listening on `<port>`, responding with `{"service":"<name>","ok":true}` |
| `services/<name>/package.json` | `name`, `private`, `type: module`, `scripts.dev/start` |

---

## Validators

No validators. Outputs are dynamic (based on `name`) so static path checks cannot be declared in the manifest.

---

## Commands

No `PostGenerationCommands`. No `TestCommands`.

---

## Conflicts

None.

---

## Loop usage

`service_writer` is designed to run multiple times in a single scaffold. Each invocation has a distinct `LoopFrame` in its `LoopStack`:

```go
// From flows/microservices.go (simplified):
for i, iter := range servicesRaw {
    invs = append(invs, flows.Invocation{
        Name: "service_writer",
        LoopStack: []flow.LoopFrame{
            {QuestionID: "services", Index: i, Answers: iter},
        },
    })
}
```

The resolver preserves duplicate names in the explicit invocation list (unlike auto-added deps which are deduplicated) so all three `service_writer` invocations survive the topological sort.

---

## Fixture example

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
  }
}
```
