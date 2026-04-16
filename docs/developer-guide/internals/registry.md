# Registry

Implementation: `internal/generator/registry.go`

---

## What the Registry does

The Registry holds all registered generators. At startup, `buildRegistry()` in `cmd/dot/build.go` creates a `Registry` and calls `Register` for each official generator. The resulting registry is passed through the CLI to every command that needs it.

At runtime, `ForSpec` resolves which generators apply to a given Spec. `Get` looks up a specific generator by name for `dot new` dispatch.

---

## Registration

```go
func (r *Registry) Register(g Generator) error
```

Returns an error if two generators claim the same `(Language, Module)` pair. The conflict matrix:

| Generator A language | Generator B language | Same module? | Conflict? |
|---|---|---|---|
| `"go"` | `"go"` | yes | yes |
| `"go"` | `"python"` | yes | no |
| `"go"` | `"*"` | yes | yes |
| `"*"` | `"*"` | yes | yes |
| any | any | no | no |

Language `"*"` conflicts with any language-specific generator on the same module. This prevents a language-agnostic generator and a language-specific one from both trying to write the same module's files.

Registration errors are caught at startup before any user action runs. In `cmd/dot/build.go`, registration uses `must()` which panics on error. A registration conflict is a programming error — it should fail loudly during development, not silently at runtime.

---

## Matching (ForSpec)

```go
func (r *Registry) ForSpec(s spec.Spec) []Generator
```

Returns generators where:
- `generator.Language() == spec.Project.Language` OR `generator.Language() == "*"`
- AND at least one of `generator.Modules()` appears in `spec.Modules[].Name`

Eight-case truth table:

| Language match | Module match | Included? |
|---|---|---|
| exact match | exact match | yes |
| exact match | no match | no |
| `"*"` | exact match | yes |
| `"*"` | no match | no |
| different language | exact match | no |
| different language | no match | no |
| exact match | empty modules | no |
| empty registry | any | none returned |

---

## Dispatch (Get)

```go
func (r *Registry) Get(name string) (Generator, bool)
```

Used by `dot new` to look up a generator by name. The flow:

```
dot new route UserController
  → Load(".") → ctx.Commands["new route"]
  → CommandRef{Generator: "go-rest-api", Action: "rest-api.new-route"}
  → registry.Get("go-rest-api") → generator
  → generator.RunAction("rest-api.new-route", ["UserController"], ctx.Spec)
  → []FileOp → pipeline.Run
```

If `Get` returns `false`, the command key was in `.dot/config.json` but the generator is not registered. This means the project was initialized with a version of dot that had a generator the current binary doesn't have — dot surfaces this as an error with a clear message.
