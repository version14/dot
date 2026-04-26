# DOT: Scaffolding CLI — Architecture & Design Document

**Version:** 0.3
**Last Updated:** April 2026
**Status:** Pre-implementation — defines architecture before coding

---

## 1. PURPOSE & PROBLEM STATEMENT

### What is DOT?

`dot` is a **generative project scaffolding engine** that eliminates repetitive boilerplate by
automating project setup through a hook-based question flow and composable, order-dependent
generators.

`dot` is not just a CLI — it is a **core engine** that can be exposed through multiple interfaces:
- A **CLI** (primary, v1)
- An **MCP** (AI-driven scaffolding)
- A **Web App** (visual project builder)
- An **API** (programmatic scaffolding)

All interfaces consume the same core engine. The CLI is the first implementation.

### The Problem It Solves

**User:** Developers at Décrypté initially, then the open-source community.

**Before dot:**
- Copy-paste from templates (inconsistent, outdated)
- Git clone starter repos (bloated, not customizable)
- Manual edits to configuration files across layers
- Setting up a monorepo with 10 services, 2 apps, migrations, permissions, CI — takes hours and is
  error-prone

**After dot:**
```bash
dot scaffold new
# Answer questions interactively
# → Full project generated, dependencies installed, ready to run
```

---

### Scope: What DOT Does & Doesn't Do

#### ✅ In Scope (v1)
- Interactive CLI prompts (option select, multi-select, text, confirm, loop, conditional)
- Hook-based question flow engine with plugin injection points
- Generate projects from templates and custom Go functions
- Intelligent file editing (JSON, YAML, code AST manipulation)
- Fetch base templates from remote GitHub repositories
- Composable generators with explicit dependencies and conflict declarations
- Post-generation commands (dependency installation: `pnpm i`, `go mod download`, etc.)
- Validation of final project structure before writing
- Community generators via external repositories
- Generator versioning with semver constraints

#### 🔜 In Scope (Future Versions)
- CI/CD provisioning (base workflows: PR merging, branch protection, etc.)
- Docker-based deployment scaffolding (dev server + production)
- `dot add <feature>` sub-commands (add a route, add a service, etc.)
- MCP interface for AI-driven scaffolding
- Web app for visual project building

#### ❌ Permanently Out of Scope
- Live code editing after generation (not an IDE plugin)
- Package manager (`apt`, `brew`, etc.)
- Auto-upgrade of generated projects when generators change
- Rollback on failure (v1 — manual recovery)

---

## 2. THE FLOW ENGINE

### 2.1 Overview

The flow engine is a **directed graph traversal system** built in Go. It drives the interactive
question session and produces the `ProjectSpec` — the complete description of what to generate.

**Key principles:**
- Flows are **Go code** (not YAML/JSON) — type safety, compiler validation, full expressiveness
- The graph is **composed at startup** from the core skeleton + all loaded plugins
- Each node (question) can route to different branches based on the user's answer
- Flow fragments are **contextual variables** — they resolve differently depending on which plugins
  are loaded
- Generators are resolved **post-traversal** — after the full graph is walked

---

### 2.2 The `Answer` Type

An `Answer` represents the value returned by a single question node. Because the answer tree is
recursive (loops produce nested maps, multi-select produces slices), `Answer` is defined as:

```go
// internal/flow/answer.go

// Answer is the value returned by a single question node.
// Concrete types: string, bool, int, []string (multi-select)
type Answer = interface{}

// AnswerNode is a recursive tree node in the ProjectSpec.
// string | bool | int | []string | map[string]AnswerNode | []map[string]AnswerNode
type AnswerNode = interface{}
```

**Rationale for `interface{}`:** A richer wrapper type (e.g. a struct with Scalar/Children/Loop
fields) is tempting but premature. The real usage pain points only become clear when writing the
first generators. Start with `interface{}`, add a typed wrapper in v0.2 if the assertion noise
becomes a problem.

---

### 2.3 Question Node Types

Each question is a typed Go struct. The `Question` interface is intentionally minimal — it
describes what the **engine** needs, not what Huh needs (that's the adapter's job):

```go
// internal/flow/question.go

// Question is the engine-facing interface for a flow node.
// It does NOT know about Huh or any terminal library.
type Question interface {
    ID() string
    // Next returns the edge to follow given a user's answer.
    // Called by the engine after the adapter collects the answer.
    Next(answer Answer) *Next
}
```

Note: `Ask()` is **not** on this interface. Asking is the adapter's responsibility (see Section 2.6).
The engine calls `Next()` to route; the adapter calls Huh to render.

**Built-in question types:**

```go
// Option / Select (single or multi-select)
type OptionQuestion struct {
    ID_         string
    Label       string
    Description string
    Multiple    bool      // true = multi-select (e.g. databases)
    Options     []*Option
    Next_       *Next     // used when Multiple=true: same continuation for all selections
}

type Option struct {
    Label string
    Value string
    Next  *Next // branch to follow if this option is chosen (single-select only)
}

// Text input
type TextQuestion struct {
    ID_         string
    Label       string
    Description string
    Default     string
    Validate    func(string) error
    Next_       *Next
}

// Confirm (yes/no)
type ConfirmQuestion struct {
    ID_     string
    Label   string
    Default bool
    Then    *Next // branch if true
    Else    *Next // branch if false
}

// Loop — asks count first ("How many services?"), then renders N groups
type LoopQuestion struct {
    ID_      string
    Label    string          // e.g. "How many services?"
    Body     []Question      // questions to repeat for each iteration
    Continue *Next           // where to go after all iterations complete
}

// If — no user input, pure routing based on accumulated answers
type IfQuestion struct {
    ID_       string
    Condition func(ctx *FlowContext) bool
    Then      *Next
    Else      *Next
}
```

---

### 2.4 Graph Traversal & Routing

The `Next` struct is the **edge** in the graph:

```go
// internal/flow/next.go
type Next struct {
    Question *Question // go directly to this question node
    Fragment string    // resolve a named fragment (may inject plugin questions)
    End      bool      // end the flow
}
```

**Traversal loop (engine.go):**
1. Start at the root question node
2. If `IfQuestion` — evaluate condition, follow `Then` or `Else` without asking
3. Otherwise — call adapter to ask the question, get an `Answer`
4. Call `question.Next(answer)` to get the next edge
5. If `Next.Fragment` is set — resolve the fragment (may inject plugin questions)
6. Record visited node ID in `FlowContext.VisitedNodes`
7. Continue until `Next.End == true`

**Example flow (simplified):**
```go
// flows/monorepo.go
var ProjectTypeQuestion = &OptionQuestion{
    ID_:   "type",
    Label: "Project type?",
    Options: []*Option{
        {Label: "Monorepo", Value: "monorepo", Next: &Next{Question: &MonorepoSetupQuestion}},
        {Label: "Monolith", Value: "monolith", Next: &Next{Question: &MonolithSetupQuestion}},
    },
}

var LinterQuestion = &OptionQuestion{
    ID_:   "linter",
    Label: "Linter?",
    Options: []*Option{
        {Label: "Biome",    Value: "biome",    Next: &Next{Fragment: "after-linter"}},
        {Label: "Prettier", Value: "prettier", Next: &Next{Fragment: "after-linter"}},
        {Label: "None",     Value: "none",     Next: &Next{Fragment: "after-linter"}},
    },
}
```

---

### 2.5 Huh Integration — The Adapter Pattern

**The core tension:** Huh wants a fully-declared form upfront. The flow graph is dynamic —
branches aren't known until answers arrive.

**Resolution:** One group = one question. The engine drives traversal; the adapter translates
each question node into a single `huh.NewForm(huh.NewGroup(...))` call and runs it immediately.

```
Engine loop:          node → adapter.Ask(node) → answer → node.Next(answer) → next node
Huh per question:     huh.NewForm(huh.NewGroup(oneField)).Run()
```

This means:
- **Back navigation within a question** (e.g. editing a multi-select before confirming) — works
  natively via Huh's internal navigation.
- **Back navigation across questions** — handled by the engine re-running the previous form. The
  engine maintains a `history []questionSnapshot` stack for this.
- **Branch changes on back** — if user goes back to `linter` and changes from `biome` to
  `prettier`, the engine pops the `FlowContext` to that point and re-traverses forward. The
  `after-linter` fragment resolves differently on the new path.

```go
// internal/cli/prompt.go

// Adapter — translates a Question node into a Huh form and runs it.
// The engine calls Ask(); it never imports Huh directly.
type HuhAdapter struct{}

func (a *HuhAdapter) Ask(q flow.Question, ctx *flow.FlowContext) (flow.Answer, error) {
    switch q := q.(type) {

    case *flow.OptionQuestion:
        var result string
        opts := make([]huh.Option[string], len(q.Options))
        for i, o := range q.Options {
            opts[i] = huh.NewOption(o.Label, o.Value)
        }
        form := huh.NewForm(huh.NewGroup(
            huh.NewSelect[string]().
                Title(q.Label).
                Description(q.Description).
                Options(opts...).
                Value(&result),
        ))
        if err := form.Run(); err != nil {
            return nil, err
        }
        return result, nil

    case *flow.TextQuestion:
        var result string
        form := huh.NewForm(huh.NewGroup(
            huh.NewInput().
                Title(q.Label).
                Description(q.Description).
                Placeholder(q.Default).
                Validate(q.Validate).
                Value(&result),
        ))
        if err := form.Run(); err != nil {
            return nil, err
        }
        return result, nil

    case *flow.ConfirmQuestion:
        var result bool
        form := huh.NewForm(huh.NewGroup(
            huh.NewConfirm().
                Title(q.Label).
                Value(&result),
        ))
        if err := form.Run(); err != nil {
            return nil, err
        }
        return result, nil

    case *flow.LoopQuestion:
        // Handled separately — see Section 2.7
        return a.AskLoop(q, ctx)
    }

    return nil, fmt.Errorf("unknown question type: %T", q)
}
```

**The engine never imports Huh.** `prompt.go` is the only file that knows about Huh. When the
MCP or web app adapter is built, it replaces `prompt.go` with a different adapter — the engine
and flow graph are unchanged.

---

### 2.6 Loop Questions — Count First, Then N Groups

Loops use a **count-first** model:
1. Ask "How many services?" → user answers `3`
2. Render groups `1/3`, `2/3`, `3/3` sequentially
3. Each group asks the loop body questions for that iteration
4. Progress indicator (`1/3`) is rendered by Lipgloss in `output.go`

```go
// internal/cli/prompt.go
func (a *HuhAdapter) AskLoop(q *flow.LoopQuestion, ctx *flow.FlowContext) (flow.Answer, error) {
    // Step 1: ask count
    var count int
    form := huh.NewForm(huh.NewGroup(
        huh.NewInput().
            Title(q.Label).         // e.g. "How many services?"
            Validate(validateInt).
            Value(&countStr),
    ))
    if err := form.Run(); err != nil {
        return nil, err
    }
    count = parseCount(countStr)

    // Step 2: ask each iteration
    results := make([]map[string]Answer, count)
    for i := 0; i < count; i++ {
        output.PrintProgress(i+1, count, q.Label) // "Service 1/3"
        iterAnswers := map[string]Answer{}

        for _, bodyQuestion := range q.Body {
            answer, err := a.Ask(bodyQuestion, ctx)
            if err != nil {
                return nil, err
            }
            iterAnswers[bodyQuestion.ID()] = answer
        }
        results[i] = iterAnswers
    }

    return results, nil
}
```

---

### 2.7 Flow Fragments (Contextual Variables)

A **flow fragment** is a named slot in the graph that resolves to different question sub-graphs
depending on context (which plugins are loaded, which answers were given).

```go
// internal/flow/fragment.go

// FlowFragment is a named contextual resolver.
// Not a fixed sub-graph — a function that returns the appropriate Next.
type FlowFragment struct {
    ID      string
    Resolve func(ctx *FlowContext) *Next
}
```

**Example:**
- Fragment `"after-linter"` with biome plugin loaded → injects `BiomeConfigQuestion`
- Fragment `"after-linter"` with no plugin → passes through to next core question
- Fragment `"after-linter"` after user went back and changed linter → resolves differently

---

### 2.8 Plugin Hook Injection

The core flow defines **named hooks** — extension points where plugins inject questions:

```go
// internal/flow/hook.go
const (
    HookAfterProjectType    = "after-project-type"
    HookAfterLinter         = "after-linter"
    HookAfterBackendSetup   = "after-backend-setup"
    HookAfterFrontendSetup  = "after-frontend-setup"
    HookAfterServiceDefined = "after-service-defined"
    HookBeforeFinalize      = "before-finalize"
)
```

Plugin registration:
```go
func (p *BiomePlugin) Register(engine *FlowEngine) {
    engine.Hook(HookAfterLinter, func(ctx *FlowContext) *Next {
        if ctx.Answers["linter"] != "biome" {
            return nil // no-op
        }
        return &Next{Question: &BiomeConfigQuestion}
    })
    engine.RegisterGenerator("biome-config", BiomeGenerator)
}
```

**Graph composition at startup:**
1. Load core flow skeleton
2. Load all plugins (built-in + installed community plugins)
3. Each plugin taps into named hooks
4. Engine merges all hook registrations into the graph
5. Final graph is ready for traversal

---

### 2.9 FlowContext & Scoping

```go
// internal/flow/context.go
type FlowContext struct {
    Answers       map[string]AnswerNode // recursive answer tree (built during traversal)
    LoopStack     []LoopFrame           // scope stack — one frame per active loop level
    VisitedNodes  []string              // traversal path (used for generator resolution)
    LoadedPlugins []string              // which plugins contributed to the flow
    History       []HistoryEntry        // for back navigation — snapshots of previous states
}

type LoopFrame struct {
    QuestionID string                // which LoopQuestion we're inside
    Index      int                   // current iteration (0-based)
    Answers    map[string]AnswerNode // answers for this iteration
}

type HistoryEntry struct {
    NodeID  string
    Context FlowContext // snapshot — popped to when user goes back
}
```

**Scoping rules:**
- **Global answers**: top-level `Answers` map, accessible by all generators
- **Loop scope**: each iteration lives in `LoopStack`. Generators inside a loop see current frame
  answers + all parent frame answers + global answers (scope chain, deeper wins on conflicts)

---

## 3. PROJECTSPEC

### 3.1 Answer Structure

The spec uses a **recursive answer tree** — the structure mirrors the flow traversal exactly.

**Rules:**
- Key = question node `ID()` (always, 1:1)
- Only visited nodes appear — unvisited branches don't exist in the spec
- Scalar answer = `string`, `bool`, `int`
- Multi-select answer = `[]string`
- Loop answer = `[]map[string]AnswerNode` (array of nested answer maps)
- Nesting is unlimited and recursive

**Go type:**
```go
// AnswerNode is recursive — scalar | map | array of maps
type AnswerNode = interface{}
```

**Example — monorepo with nested loops:**
```json
{
  "project_name": "my-stack",
  "type": "monorepo",
  "linter": "biome",
  "biome_strict": true,
  "apps": [
    {
      "name": "frontend",
      "framework": "react",
      "routes": [
        { "path": "/home",    "auth": false },
        { "path": "/profile", "auth": true  }
      ]
    },
    {
      "name": "dashboard",
      "framework": "react"
    }
  ],
  "services": [
    {
      "name": "auth",
      "has_database": true,
      "db_type": "postgres",
      "tables": [
        {
          "name": "users",
          "columns": [
            { "name": "id",    "type": "uuid" },
            { "name": "email", "type": "text" }
          ]
        }
      ]
    },
    {
      "name": "user",
      "has_database": true,
      "db_type": "postgres"
    }
  ]
}
```

Note what's **absent**: `prettier` branch wasn't visited — no `prettier_*` key. `dashboard` has no
`routes` — that loop wasn't entered. The spec only contains what the user actually answered.

### 3.2 Full ProjectSpec Definition

```go
// internal/spec/spec.go
type ProjectSpec struct {
    FlowID               string                `json:"flow_id"`
    CreatedAt            time.Time             `json:"created_at"`
    Metadata             ProjectMetadata       `json:"metadata"`
    Answers              map[string]AnswerNode `json:"answers"`        // recursive tree
    VisitedNodes         []string              `json:"visited_nodes"`  // traversal audit trail
    LoadedPlugins        []string              `json:"loaded_plugins"`
    GeneratorConstraints map[string]string     `json:"generator_constraints"`
}

type ProjectMetadata struct {
    ProjectName string `json:"project_name"`
    ToolVersion string `json:"tool_version"`
}
```

No separate `loop_answers` field — loops are nested arrays inside `Answers`.

### 3.3 Generator Access to Answers (Scope Chain)

The executor flattens the recursive answer tree into a scoped flat map per generator invocation.
This logic lives in `internal/flow/scope.go`.

```go
// internal/flow/scope.go

// FlattenScope walks the LoopStack from outermost to innermost,
// merging answers at each level. Deeper scopes win on key conflicts.
func FlattenScope(global map[string]AnswerNode, stack []LoopFrame) map[string]interface{} {
    result := map[string]interface{}{}
    // Start with global
    for k, v := range global {
        result[k] = v
    }
    // Each loop frame overrides
    for _, frame := range stack {
        for k, v := range frame.Answers {
            result[k] = v
        }
    }
    return result
}
```

**Example — deeply nested:**
```go
// Generator for a table column sees:
ctx.Answers = {
    "name":         "id",       // current: column iteration
    "type":         "uuid",
    "table_name":   "users",    // parent: table iteration
    "service_name": "auth",     // grandparent: service iteration
    "project_name": "my-stack", // global
    "linter":       "biome",
}
```

---

## 4. VIRTUAL PROJECT STATE

### 4.1 Overview

Generators do **not** write directly to disk. They manipulate an **in-memory virtual filesystem**:

```go
// internal/state/virtual.go
type VirtualProjectState struct {
    Files    map[string]*FileNode
    Metadata ProjectMetadata
}

// internal/state/file.go
type FileNode struct {
    Path            string
    Content         []byte
    ContentType     ContentType  // Raw, JSON, YAML, GoMod
    CreatedBy       string       // generator name
    Transformations []string     // audit trail of edits
    ModifiedAt      time.Time
}
```

**Why virtual state?**
- Generators read the latest version of any file (set by previous generators)
- Validation runs on the complete project before writing to disk
- Dry-run / preview mode is free
- Clean conflict detection

### 4.2 File Operations API

```go
// Create
state.CreateFile("services/auth/main.go", content)

// Read (always returns latest version)
file := state.GetFile("package.json")

// Update — typed operations (internal/state/json.go, yaml.go, gomod.go)
state.UpdateJSON("package.json", func(doc *JSONDoc) {
    doc.SetNested("dependencies.axios", "^0.21.0")
})

state.UpdateYAML("docker-compose.yml", func(doc *YAMLDoc) {
    doc.Append("services", newService)
})

// go.mod has its own format — not JSON/YAML (internal/state/gomod.go)
state.UpdateGoMod(func(mod *GoMod) {
    mod.AddRequire("github.com/lib/pq", "v1.10.0")
})

// Template rendering (internal/render/template.go)
state.RenderTemplate(
    "generators/react-app/templates/src/main.tsx",
    map[string]string{"AppName": spec.Answers["project_name"]},
    "src/main.tsx",
)

// Check existence
state.FileExists("package.json")
```

### 4.3 Conflict Resolution

When two generators edit the same file, the second always **sees the modified version** from the
first — composition, not overwrite. Execution order is controlled by explicit dependency
declarations.

---

## 5. GENERATOR MODEL

### 5.1 What Is a Generator?

A generator is a **Go package** that:
- Declares its metadata, version, dependencies, conflicts, commands, and validators
- Implements `Generate(ctx *generator.Context) error`
- May be called **multiple times** if multiple flow paths trigger it (once per loop iteration)

### 5.2 Generator Manifest

```go
// generators/go-microservice/manifest.go
var Manifest = generator.Manifest{
    Name:        "go-microservice",
    Version:     "1.1.0",
    Description: "Generate a Go microservice with optional database",

    DependsOn:     []string{"base-project"},
    ConflictsWith: []string{},

    // Declared outputs — used by structural validator on re-run
    Outputs: []string{
        "services/{name}/main.go",
        "go.mod",
    },

    // Run after all files are written to disk. Deduped by (Cmd + WorkDir).
    PostGenerationCommands: []generator.Command{
        {Cmd: "go mod download", WorkDir: "."},
    },

    // Run by the flow tester to validate THIS generator's output.
    // {name} is interpolated per loop invocation.
    // Integration-level checks belong in the test case, not here.
    TestCommands: []generator.Command{
        {Cmd: "go vet ./services/{name}/...", WorkDir: "."},
    },

    // Used on re-run to check project structure is intact
    Validators: []generator.Validator{
        {
            Name: "Service entry point",
            Checks: []generator.Check{
                {Type: "file_exists", Path: "services/{name}/main.go"},
                {Type: "file_exists", Path: "go.mod"},
            },
        },
    },
}
```

### 5.3 Generator Implementation

```go
// generators/go-microservice/generator.go
func Generate(ctx *generator.Context) error {
    name  := ctx.Answers["name"].(string)
    hasDB := ctx.Answers["has_database"].(bool)

    ctx.State.RenderTemplate(
        "generators/go-microservice/templates/main.go.tmpl",
        map[string]interface{}{"Name": name, "HasDB": hasDB},
        fmt.Sprintf("services/%s/main.go", name),
    )

    ctx.State.UpdateGoMod(func(mod *GoMod) {
        if hasDB {
            mod.AddRequire("github.com/lib/pq", "v1.10.0")
        }
    })

    return nil
}
```

### 5.4 GenerationContext

```go
// pkg/dotapi/context.go
type Context struct {
    Spec         *ProjectSpec           // Full project spec
    Answers      map[string]interface{} // Scoped flat answers (from scope.FlattenScope)
    State        *VirtualProjectState   // In-memory filesystem
    PreviousGens []string               // Generators already executed
    Logger       Logger
}
```

### 5.5 Post-Generation Commands

Collected from all executed generator manifests, deduped by `(Cmd + WorkDir)`, then executed in
dependency order after disk write.

```
Monorepo with 2 apps + 10 services (before dedup):
  base-project:          git init .
  typescript-base:       pnpm install (root)
  react-app:             pnpm install (apps/frontend)
  react-dashboard:       pnpm install (apps/dashboard)
  go-microservice [×10]: go mod download (.)
  biome-config:          pnpm biome check .

After dedup:
  git init .
  pnpm install (root)
  pnpm install (apps/frontend)
  pnpm install (apps/dashboard)
  go mod download (.)     ← declared 10 times, runs once
  pnpm biome check .
```

---

## 6. GENERATOR SOURCES & VERSIONING

### 6.1 Three Sources

**1. Built-in** — compiled into the binary, maintained in the official `dot` repo.

**2. Community** — separate GitHub repos, naming convention `dot-gen-{name}`:
- Installed via `dot generator install github.com/user/dot-gen-firebase-auth`
- Loaded from `~/.dot/generators/`

**3. Local** — private generators on disk in `~/.dot/generators/`

### 6.2 Versioning Strategy

Generators use **semantic versioning** (MAJOR.MINOR.PATCH):
- **MAJOR** — breaking changes (input contract changed)
- **MINOR** — new features, backward compatible
- **PATCH** — bug fixes only

**Constraint syntax:** `^1.2.0` (compatible), `~1.2.0` (patch only), `1.2.3` (exact)

### 6.3 Version Resolution

**First run:** resolve latest matching versions → store ranges in `.dot/spec.json` (intent),
resolved exact versions in `.dot/manifest.json` (fact).

**Re-run:** always use exact versions from `.dot/manifest.json` — no drift.

**Local cache:** `~/.dot/cache/generators/{name}/{version}/`

### 6.4 Upgrading

- `dot update` — upgrade all within spec constraints
- `dot upgrade-generator go-microservice@2.0.0` — selective, validates compat first
- Manual edit of `.dot/manifest.json` + re-run — full control, error-prone

`^1.x` never auto-upgrades to `2.0.0`. MAJOR bumps require explicit user intent.

### 6.5 Reproducibility

Cloning a project and running `dot scaffold` fetches the exact generator versions from
`.dot/manifest.json` — identical output across machines.

---

## 7. EXECUTION MODEL

### 7.1 Full Lifecycle

```
┌──────────────────────────────────────────────────────────────┐
│ 1. STARTUP                                                   │
│    Load plugins → compose flow graph (skeleton + hooks)      │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 2. FLOW TRAVERSAL (interactive CLI)                          │
│    Engine walks graph node by node                           │
│    Calls HuhAdapter.Ask() per node → collects Answer         │
│    Follows Next edges, resolves fragments                    │
│    Builds FlowContext (recursive answer tree + visited nodes) │
│    Maintains history stack for back navigation               │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 3. PROJECTSPEC CONSTRUCTION                                  │
│    FlowContext → ProjectSpec                                  │
│    Attach generator constraints (defaults + overrides)       │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 4. GENERATOR RESOLUTION (post-traversal)                     │
│    Input: answers + visited nodes + loaded plugins           │
│    Match generators, topological sort, detect conflicts      │
│    Expand loops → one invocation per iteration               │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 5. PRE-EXECUTION VALIDATION (internal/generator/constraints) │
│    All dependencies satisfied? No cycles? No conflicts?      │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 6. GENERATOR EXECUTION (ordered)                             │
│    For each invocation (topological order):                  │
│      scope.FlattenScope() → scoped ctx.Answers               │
│      Generate(ctx) → updates VirtualProjectState             │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 7. POST-EXECUTION VALIDATION (internal/generator/structural) │
│    VirtualProjectState vs declared Validators                │
│    Required files exist? Structure matches expectations?     │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 8. WRITE TO DISK                                             │
│    Write VirtualProjectState → target directory              │
│    Save .dot/spec.json + .dot/manifest.json                  │
└──────────────────────┬───────────────────────────────────────┘
                       │
┌──────────────────────▼───────────────────────────────────────┐
│ 9. POST-GENERATION COMMANDS                                  │
│    Collect from all generator manifests → dedup → execute    │
│    pnpm i, go mod download, docker build, etc.               │
└──────────────────────────────────────────────────────────────┘
```

### 7.2 Failure Model

- **Generator fails** → halt immediately, no files written to disk (virtual state only)
- **Post-generation command fails** → partial files on disk, manual recovery required
- `.dot/` not updated on failure — spec and manifest remain at last successful state

---

## 8. IDEMPOTENCY & RE-RUNNING

### 8.1 Re-running in an Existing Project

```bash
cd my-project && dot scaffold
```

1. Detect `.dot/spec.json`
2. Load spec + exact generator versions from `.dot/manifest.json`
3. Run structural validators against current project state
4. Valid → "Already generated, nothing to do"
5. Invalid → report failed validators, suggest fixes
6. Spec modified → ask confirmation before re-generating

### 8.2 Validators

Declared in each generator manifest. Used on re-run to check project structure:

```go
Validators: []generator.Validator{
    {
        Name: "React entry point",
        Checks: []generator.Check{
            {Type: "file_exists",     Path: "src/main.tsx"},
            {Type: "json_key_exists", Path: "package.json", Key: "dependencies.react"},
        },
    },
},
```

---

## 9. STORED STATE: `.dot/` DIRECTORY

### `.dot/spec.json` — Intent (ranges + answers)

```json
{
  "flow_id": "monorepo",
  "created_at": "2026-04-25T10:30:00Z",
  "metadata": { "project_name": "my-stack", "tool_version": "0.3.0" },
  "answers": {
    "project_name": "my-stack",
    "type": "monorepo",
    "linter": "biome",
    "biome_strict": true,
    "apps": [
      { "name": "frontend",  "framework": "react" },
      { "name": "dashboard", "framework": "react" }
    ],
    "services": [
      { "name": "auth", "has_database": true, "db_type": "postgres" },
      { "name": "user", "has_database": true, "db_type": "postgres" }
    ]
  },
  "visited_nodes": ["type", "linter", "biome_strict", "apps", "services"],
  "loaded_plugins": ["biome-plugin"],
  "generator_constraints": {
    "base-project":    "^1.0.0",
    "typescript-base": "^1.2.0",
    "react-app":       "^2.3.0",
    "go-microservice": "^1.1.0",
    "biome-config":    "^1.0.0"
  }
}
```

### `.dot/manifest.json` — Fact (resolved versions)

```json
{
  "tool_version": "0.3.0",
  "last_executed_at": "2026-04-25T10:30:20Z",
  "execution_time_ms": 2340,
  "generators_executed": [
    {
      "name": "base-project",
      "version_constraint": "^1.0.0",
      "resolved_version": "1.2.5",
      "executed_at": "2026-04-25T10:30:15Z",
      "invocation_count": 1,
      "content_hash": "sha256:abc123..."
    },
    {
      "name": "go-microservice",
      "version_constraint": "^1.1.0",
      "resolved_version": "1.2.0",
      "executed_at": "2026-04-25T10:30:18Z",
      "invocation_count": 2,
      "content_hash": "sha256:def456..."
    }
  ]
}
```

`invocation_count` — how many times a generator ran (once per loop iteration).

### `.dot/.gitignore`
```
/*
!spec.json
!manifest.json
!.gitignore
```

---

## 10. TESTING: FLOW TESTER

### 10.1 Purpose

Validates end-to-end that a complete flow:
1. Accepts scripted answers
2. Generates all expected files
3. Passes generator-level quality checks (from manifests)
4. Passes integration-level quality checks (from test case)

### 10.2 Test Command Split

Test commands come from **two sources**, collected and deduped by the runner:

**Generator manifest** — validates *this generator's output* in isolation. Runs for every test
case that includes this generator. Variables like `{name}` are interpolated per loop invocation.

```go
// generators/typescript_base/manifest.go
TestCommands: []generator.Command{
    {Cmd: "pnpm tsc --noEmit", WorkDir: "."},
},

// generators/biome_config/manifest.go
TestCommands: []generator.Command{
    {Cmd: "pnpm biome check .", WorkDir: "."},
},

// generators/go_microservice/manifest.go
TestCommands: []generator.Command{
    // runs once per service — {name} interpolated from scoped answers
    {Cmd: "go vet ./services/{name}/...", WorkDir: "."},
},
```

**Test case** — validates the *combined result* of the full flow. Integration checks that only
make sense once all generators have run together.

```json
"test_commands": [
  { "cmd": "go build ./...", "workdir": "." },
  { "cmd": "pnpm build",     "workdir": "apps/frontend" }
]
```

**Runner collection logic:**
```
Final test commands for a monorepo test case:
  pnpm tsc --noEmit          ← typescript_base manifest
  pnpm biome check .         ← biome_config manifest
  go vet ./services/auth/... ← go_microservice manifest (invocation 1)
  go vet ./services/user/... ← go_microservice manifest (invocation 2)
  go build ./...             ← test case (integration)
  pnpm build                 ← test case (integration)
```

### 10.3 Test Case Format

```json
{
  "name": "monorepo-react-go-services",
  "flow": "monorepo",

  "answers": {
    "project_name": "test-app",
    "type": "monorepo",
    "linter": "biome",
    "apps": [
      { "name": "frontend", "framework": "react" }
    ],
    "services": [
      { "name": "auth", "has_database": true, "db_type": "postgres" },
      { "name": "user", "has_database": true, "db_type": "postgres" }
    ]
  },

  "expected_files": [
    "package.json",
    "tsconfig.json",
    "biome.json",
    "go.mod",
    "apps/frontend/src/main.tsx",
    "services/auth/main.go",
    "services/auth/migrations/001_create_users.sql",
    "services/user/main.go",
    "services/gateway/main.go"
  ],

  "test_commands": [
    { "cmd": "go build ./...", "workdir": "." },
    { "cmd": "pnpm build",     "workdir": "apps/frontend" }
  ],

  "cache_invalidators": [
    "generators/base_project/",
    "generators/typescript_base/",
    "generators/react_app/",
    "generators/biome_config/",
    "generators/go_microservice/",
    "generators/postgres_migrations/",
    "generators/gateway_setup/"
  ]
}
```

Empty `test_commands: []` is valid — generator manifests may cover everything already.

### 10.4 Caching

```
Cache key = sha256(answers + resolved_generator_versions + invalidator_file_hashes)

On run:
  cache hit  && no invalidator changed → skip, report cached result
  cache miss || invalidator changed    → full run, store result

Stored at: tools/test-flow/.cache/{hash}.json
```

### 10.5 Directory Structure

```
tools/test-flow/
├── main.go           # Entry point for the tester binary
├── runner.go         # Executes one test case end-to-end
├── cache.go          # Cache key computation + read/write
├── reporter.go       # Table output: name | status | duration | cache hit
└── testdata/
    ├── 202604231948-ts-backend-express-clean.json
    ├── 202604231949-ts-backend-express-hexagonal.json
    ├── 202604231950-ts-backend-express-mvc.json
    ├── 202604231951-ts-frontend-react-feature-sliced.json
    ├── 202604232112-ts-frontend-react-atomic-design.json
    └── 202604232114-ts-frontend-react-atomic-container-presentational.json
```

Naming: `{timestamp}-{descriptor}.json` — timestamp preserves creation order.

---

## 11. DIRECTORY STRUCTURE

```
dot/
├── cmd/
│   └── dot/
│       └── main.go                      # Entry point — wires cobra + plugins + engine. No logic.
│
├── internal/
│   ├── cli/
│   │   ├── command.go                   # Cobra command definitions (scaffold, update, upgrade-generator, generator install)
│   │   ├── prompt.go                    # HuhAdapter — translates Question nodes into Charm Huh form calls
│   │   └── output.go                    # Lipgloss styles — progress (1/3), section headers, errors
│   │
│   ├── flow/
│   │   ├── engine.go                    # Graph traversal loop + history stack for back navigation
│   │   ├── question.go                  # Question interface + all types (Option, Text, Confirm, Loop, If)
│   │   ├── answer.go                    # Answer and AnswerNode type definitions
│   │   ├── next.go                      # Next struct — graph edge (Question | Fragment | End)
│   │   ├── fragment.go                  # Fragment registry — named contextual resolvers
│   │   ├── hook.go                      # Hook constants + hook registry for plugin injection
│   │   ├── context.go                   # FlowContext — answer tree, loop stack, visited nodes, history
│   │   └── scope.go                     # FlattenScope — resolves scope chain into flat map for generators
│   │
│   ├── spec/
│   │   ├── spec.go                      # ProjectSpec definition (answer tree + constraints)
│   │   ├── builder.go                   # FlowContext → ProjectSpec
│   │   └── loader.go                    # Load .dot/spec.json for re-run / dot update
│   │
│   ├── generator/
│   │   ├── manifest.go                  # Manifest struct (version, deps, conflicts, commands, validators)
│   │   ├── registry.go                  # Discovers + loads all generators (built-in + community + local)
│   │   ├── resolver.go                  # Post-traversal: matches generators to spec, expands loop invocations
│   │   ├── sorter.go                    # Topological sort by DependsOn + circular dep / conflict detection
│   │   ├── executor.go                  # Runs invocations in order, calls scope.FlattenScope per invocation
│   │   ├── constraints.go               # Pre-execution: deps satisfied, no cycles, no conflicts
│   │   ├── structural.go                # Post-execution: project structure vs declared Validators
│   │   └── errors.go                    # Typed errors: ErrCircularDep, ErrConflict, ErrMissingDep, etc.
│   │
│   ├── state/
│   │   ├── virtual.go                   # VirtualProjectState — in-memory filesystem (map[path]*FileNode)
│   │   ├── file.go                      # FileNode — content, ContentType, creator, transformation trail
│   │   ├── json.go                      # JSONDoc operations: SetNested, Merge, AddDep, DeleteKey
│   │   ├── yaml.go                      # YAMLDoc operations: Append, Merge, SetKey
│   │   ├── gomod.go                     # GoMod operations: AddRequire, AddModule, RemoveRequire
│   │   └── persist.go                   # Writes VirtualProjectState to disk via fileutils.SafeWrite
│   │
│   ├── render/
│   │   ├── template.go                  # Go text/template parsing + rendering for .tmpl files
│   │   └── remote.go                    # Fetch template files from remote URLs (GitHub or other)
│   │
│   ├── plugin/
│   │   ├── loader.go                    # Loads community + local plugins at startup, calls plugin.Register(engine)
│   │   ├── installed.go                 # List of installed plugins in ~/.dot/generators/ (not generator/registry.go)
│   │   └── installer.go                 # dot generator install <repo> — fetches plugin to ~/.dot/generators/
│   │
│   ├── dotdir/
│   │   ├── spec.go                      # Read/write .dot/spec.json (was split across spec/ and versioning/)
│   │   └── manifest.go                  # Read/write .dot/manifest.json
│   │
│   ├── versioning/
│   │   ├── semver.go                    # Semver parsing + constraint resolution (^, ~, exact)
│   │   └── cache.go                     # ~/.dot/cache/generators/{name}/{version}/ — local version cache
│   │
│   ├── commands/
│   │   ├── runner.go                    # Executes post-generation commands, streams output via cli/output.go
│   │   └── dedup.go                     # Dedup commands by (Cmd + WorkDir)
│   │
│   └── fileutils/
│       ├── walk.go                      # Recursive directory traversal (used by structural validator)
│       ├── safe_write.go                # Atomic write: tmp file → os.Rename (prevents partial writes)
│       └── path.go                      # Path normalization, join, clean — consistent across OS
│
├── flows/                               # Built-in flow definitions (Go code, compiled into binary)
│   ├── registry.go                      # Maps flow ID strings → root question nodes (used by engine on re-run)
│   ├── monorepo.go                      # Monorepo flow root graph
│   ├── monolith.go                      # Monolith flow root graph
│   ├── ts_backend.go                    # TypeScript backend-only flow
│   ├── react_app.go                     # React app-only flow
│   └── fragments/
│       ├── linter.go                    # "after-linter" fragment — biome / prettier / none branch
│       ├── database.go                  # "after-backend-setup" fragment — db type selection
│       └── service.go                   # Reusable service-definition sub-graph (used in loop body)
│
├── generators/                          # Built-in generators (Go packages, compiled into binary)
│   ├── base_project/
│   │   ├── manifest.go                  # Declares: deps, outputs, post-gen commands, test commands, validators
│   │   ├── generator.go                 # Implements Generate(ctx)
│   │   └── templates/                   # .tmpl files rendered by render/template.go
│   ├── typescript_base/
│   ├── react_app/
│   ├── go_microservice/
│   ├── biome_config/
│   ├── postgres_migrations/
│   └── gateway_setup/                   # services/gateway/main.go — API gateway + route registration
│
├── pkg/                                 # Public API — imported by community plugins
│   └── dotapi/
│       ├── generator.go                 # Generator interface: Generate(*Context) error + Name() + Version()
│       ├── context.go                   # Context struct (Spec, Answers, State, Logger)
│       ├── manifest.go                  # Manifest struct + Command + Validator + Check types
│       └── flow.go                      # Question interface + hook registration API for plugin authors
│
├── tools/
│   └── test-flow/
│       ├── main.go
│       ├── runner.go
│       ├── cache.go
│       ├── reporter.go
│       └── testdata/
│
└── go.mod                               # Single module root — covers all packages in the repo
```

---

## 12. LIMITATIONS & KNOWN CONSTRAINTS

### Hard Limitations (v1)

1. **No rollback on generator failure** — Execution halts. No files written to disk on generator
   failure (virtual state only). Post-generation command failures leave partial files. Manual
   recovery required.

2. **No live editing** — DOT generates once. Not an incremental code editor.

3. **No auto-upgrade** — Generated projects don't sync with generator changes. Explicit
   `dot update` required.

4. **Go-only flows** — Flow definitions require Go. Non-Go contributors need to use the plugin
   API or fork.

5. **No semantic validation at generation time** — DOT does not verify generated code compiles.
   The flow tester handles this at test time.

### Known Open Questions (Deferred)

- **Community plugin loading** — Go packages require compilation. How are community generators
  loaded without rebuilding the binary? (Go plugins `.so` are painful. Subprocess model?
  `go run`? TBD.)
- **Registry format** — How does the CLI discover available community generators? JSON file,
  GitHub topics, dedicated registry service? TBD.
- **Offline mode** — Generator not in local cache → fail immediately or prompt? TBD.
- **Version deprecation notices** — CLI warn when a generator version reaches EOL? TBD.

---

## 13. FUTURE MODULES

**CI/CD Provisioning** — Base GitHub Actions workflows (PR checks, branch protection). Generated
as `.github/workflows/*.yaml`, declared by generators.

**Docker Deployment Scaffolding** — `Dockerfile` per service, `docker-compose.yml` for dev,
production-ready compose for staging/prod.

**Sub-commands (`dot add`)** — `dot add route`, `dot add service`, `dot add migration`. Reads
`.dot/spec.json` to understand existing project structure before generating.

**MCP Interface** — Same core engine, different input adapter. Receives AI-structured answers
instead of Huh forms. Flow graph and hook system unchanged.

**Web App** — Visual question flow builder, community generator browser, live project structure
preview. Serializes to same `ProjectSpec` format.

---

## 14. APPENDIX: Full Example Walkthrough

### Scenario: Monorepo with 2 Services

**Flow traversal (one Huh form per question):**
```
[type]         → monorepo                         huh.Select
[project_name] → my-stack                         huh.Input
[linter]       → biome                            huh.Select
  fragment: after-linter → biome plugin injects:
  [biome_strict] → yes                            huh.Confirm
[apps] (loop)
  "How many apps?" → 2                            huh.Input
  App 1/2:
    [name]      → frontend                        huh.Input
    [framework] → react                           huh.Select
  App 2/2:
    [name]      → dashboard                       huh.Input
    [framework] → react                           huh.Select
[services] (loop)
  "How many services?" → 2                        huh.Input
  Service 1/2:
    [name]         → auth                         huh.Input
    [has_database] → yes                          huh.Confirm
    [db_type]      → postgres                     huh.Select
  Service 2/2:
    [name]         → user                         huh.Input
    [has_database] → yes                          huh.Confirm
    [db_type]      → postgres                     huh.Select
```

**Resulting `.dot/spec.json` answers tree:**
```json
{
  "type": "monorepo",
  "project_name": "my-stack",
  "linter": "biome",
  "biome_strict": true,
  "apps": [
    { "name": "frontend",  "framework": "react" },
    { "name": "dashboard", "framework": "react" }
  ],
  "services": [
    { "name": "auth", "has_database": true, "db_type": "postgres" },
    { "name": "user", "has_database": true, "db_type": "postgres" }
  ]
}
```

**Resolved generators (topological order):**
```
1. base-project@1.2.5
2. typescript-base@1.3.0
3. biome-config@1.0.2        ← injected by biome-plugin hook
4. react-app@2.3.1
5. go-microservice@1.2.0     ← invoked twice (auth, user)
6. postgres-migrations@1.0.5 ← invoked twice (auth, user)
7. gateway-setup@1.1.0
```

**Post-generation commands (after dedup):**
```
git init .
pnpm install (root)
pnpm install (apps/frontend)
go mod download (.)          ← declared ×2 by go-microservice, runs once
pnpm biome check .
```

**Final generated structure:**
```
my-stack/
├── .github/workflows/
├── .gitignore
├── README.md
├── package.json
├── tsconfig.json
├── biome.json
├── go.mod
├── apps/
│   └── frontend/src/main.tsx
├── vite.config.ts
├── services/
│   ├── auth/
│   │   ├── main.go
│   │   └── migrations/001_create_users.sql
│   ├── user/
│   │   ├── main.go
│   │   └── migrations/001_create_users.sql
│   └── gateway/main.go
└── .dot/
    ├── spec.json
    └── manifest.json
```

---

## 15. IMPLEMENTATION ROADMAP

### Phase 1 — Core Engine (v0.1)
- [ ] Answer + AnswerNode types
- [ ] Question types + FlowContext + scope.go
- [ ] Engine traversal loop + HuhAdapter (one form per question)
- [ ] Loop questions (count-first model + Lipgloss progress)
- [ ] Hook registry (skeleton, no plugins yet)
- [ ] ProjectSpec builder
- [ ] VirtualProjectState + JSON/YAML/GoMod operations
- [ ] Generator manifest + registry + executor
- [ ] dotdir read/write (spec.json + manifest.json)

### Phase 2 — First Generators (v0.2)
- [ ] `base-project`, `typescript-base`, `react-app`
- [ ] Template rendering + remote fetcher
- [ ] Post-generation commands with dedup
- [ ] Flow tester: basic (file existence + manifest TestCommands)
- [ ] reporter.go

### Phase 3 — Full Feature Set (v0.3)
- [ ] Topological sort + constraints.go + structural.go
- [ ] `go-microservice`, `postgres-migrations`, `gateway-setup`
- [ ] Flow tester: test case test_commands (integration checks)
- [ ] Test caching + cache invalidation

### Phase 4 — Versioning & Plugins (v0.4)
- [ ] Semver resolution + local cache
- [ ] `dot update` + `dot upgrade-generator`
- [ ] Plugin loader + hook injection from plugins
- [ ] flows/registry.go + re-run path

### Phase 5 — Polish & Community (v1.0)
- [ ] Full test coverage
- [ ] Documentation
- [ ] CLI UX polish (Lipgloss theme)
- [ ] CI/CD base workflow generation

---

*Document status: v0.3 — complete for Phase 1 implementation. Open questions noted in Section 12.*
