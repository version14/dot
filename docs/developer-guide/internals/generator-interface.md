# Generator Interface

Implementation: `internal/generator/generator.go`

---

## The interface

```go
type Generator interface {
    Name() string
    Language() string
    Modules() []string
    Apply(s spec.Spec) ([]FileOp, error)
    Commands() []CommandDef
    RunAction(action string, args []string, s spec.Spec) ([]FileOp, error)
}
```

**`Name()`** — A unique, stable identifier. Used as the key in `.dot/config.json` commands. Once a generator ships, its name must never change. Example: `"go-rest-api"`.

**`Language()`** — The language this generator targets. Use `"*"` for language-agnostic generators (e.g. GitHub Actions CI). The Registry uses this to filter generators against `spec.Project.Language`.

**`Modules()`** — The module names this generator handles. Example: `[]string{"rest-api"}`. The Registry uses this to match generators against `spec.Modules[].Name`. Two generators must not claim the same `(Language, Module)` pair.

**`Apply(spec)`** — Called once during `dot init`. Receives the full project Spec. Returns the `[]FileOp` needed to scaffold the generator's files. Must be deterministic: same Spec, same FileOps, every time.

**`Commands()`** — Returns `[]CommandDef` describing the post-creation commands this generator registers. These are written to `.dot/config.json` after `dot init` so that `dot new` can dispatch them later.

**`RunAction(action, args, spec)`** — Called by `dot new`. `action` matches `CommandDef.Action`. `args` are the positional arguments (e.g. `["UserController"]`). Returns `[]FileOp` to inject into the existing project.

---

## Rules every generator must follow

**Apply() must be deterministic.** Same Spec always produces the same FileOps. No random IDs, no timestamps, no map iteration in templates that could vary between runs.

**RunAction() must be safe on an existing project.** It returns FileOps, and the pipeline will apply them to files that already exist. Prefer `Append` and `Patch` over `Create` for RunAction — don't overwrite files the user may have modified.

**Stay inside your module's concern.** A Redis generator should not write to `main.go` except via `Patch`. It should not claim paths that belong to another generator. If you need to modify a file another generator owns, use `Append` or `Patch` with a clearly documented anchor.

**Return errors, don't panic.** If the file's import block is in an unsupported form, return `ErrUnsupportedImportForm`. If a required argument is missing, return a descriptive error. The pipeline handles errors gracefully; a panic takes down the whole process.

---

## CommandDef format

```go
type CommandDef struct {
    Name        string   // "new route"
    Args        []string // ["<name>"] — for dot help display only
    Description string   // shown in dot help
    Action      string   // passed to RunAction, e.g. "rest-api.new-route"
    Generator   string   // matches Generator.Name()
}
```

**`Name` format: `"verb noun"`** where noun is a single hyphenated word with no spaces.

Correct: `"new route"`, `"new handler"`, `"add migration"`

Wrong: `"new rest api"` (space in noun), `"generate route"` (non-standard verb)

Why: `dot new` splits on the first space after `"new"`. `dot new route UserController` maps to key `"new route"` with args `["UserController"]`. A space in the noun breaks this lookup.

Use hyphens for multi-word nouns: `"new rest-endpoint"` not `"new rest endpoint"`.
