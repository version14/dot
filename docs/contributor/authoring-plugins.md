# Authoring Plugins

Plugins extend DOT without modifying its source code. A plugin can:

- **Contribute generators** — add new generators to the registry that any flow can invoke.
- **Inject into flows** — add options to existing questions, insert new questions after existing ones, or replace questions entirely.
- **Conditionally inject generators** — decide at scaffold time (based on the user's answers) whether to add extra generator invocations.

This guide covers the full lifecycle: writing a plugin, testing it locally, publishing it to GitHub, and having others install it.

---

## Table of Contents

- [Plugin anatomy](#plugin-anatomy)
- [The Provider interface](#the-provider-interface)
- [Naming rules](#naming-rules)
- [Injection kinds](#injection-kinds)
- [ResolveExtras](#resolveextras)
- [pkg/dotplugin reference](#pkgdotplugin-reference)
- [Scaffolding a plugin repository](#scaffolding-a-plugin-repository)
- [Testing a plugin locally](#testing-a-plugin-locally)
- [Publishing a plugin](#publishing-a-plugin)
- [Installing a plugin](#installing-a-plugin)
- [In-tree vs installed plugins](#in-tree-vs-installed-plugins)

---

## Plugin anatomy

A plugin is a Go package that:

1. Implements `plugin.Provider` (re-exported as `dotplugin.Provider`).
2. Calls `dotplugin.RegisterBuiltin(p)` in its `init()` function.
3. Exports a `plugin.json` at the repository root for the installer to read.

Minimum file structure:

```
my-plugin/
├── plugin.json      ← identity: id, version, description
├── plugin.go        ← Provider implementation + init()
└── go.mod           ← module github.com/you/my-plugin
```

---

## The Provider interface

```go
// pkg/dotplugin/dotplugin.go (re-exported from internal/plugin)
type Provider interface {
    ID()             dotplugin.PluginID
    Generators()     []dotplugin.Entry
    Injections()     []dotapi.Injection  // actually flow.Injection via alias
    ResolveExtras(s *dotplugin.ProjectSpec) []dotplugin.Invocation
}
```

| Method | Purpose |
|--------|---------|
| `ID()` | Returns the plugin's unique identifier. Must not contain `.`. |
| `Generators()` | Returns generator entries contributed by this plugin. |
| `Injections()` | Returns flow injections this plugin wants to apply. |
| `ResolveExtras(spec)` | Called after the flow resolver; returns additional generator invocations based on user answers. |

---

## Naming rules

Every ID a plugin contributes **must** start with `"<pluginID>."`. This prefix is enforced at registration time.

| Contributed item | Must start with |
|-----------------|----------------|
| Question ID (InsertAfter / Replace) | `"<pluginID>."` |
| Option value (AddOption) | `"<pluginID>."` |
| Generator name | `"<pluginID>."` |

Violating this rule causes `HookRegistry.Inject` to return an error and the plugin's injections will not be registered.

---

## Injection kinds

Injections are declared in `Injections()` and applied by the flow engine when it visits the targeted question.

### InjectAddOption

Append an option to an existing `OptionQuestion`. The option's `Next` edge controls where selecting it leads.

```go
dotplugin.Injection{
    Plugin:   "my-plugin",
    TargetID: "css_framework",            // existing question to extend
    Kind:     dotplugin.InjectAddOption,
    Option: &dotplugin.Option{
        Label: "My CSS Framework",
        Value: "my-plugin.my-css",        // must have "my-plugin." prefix
        Next:  &dotplugin.Next{End: true},
    },
}
```

### InjectInsertAfter

Splice a new question (or chain of questions) after the targeted question. Once the inserted chain reaches `Next{End: true}`, the engine resumes from the original question's next edge.

```go
dotplugin.Injection{
    Plugin:   "my-plugin",
    TargetID: "project_name",
    Kind:     dotplugin.InjectInsertAfter,
    Question: &dotplugin.TextQuestion{
        QuestionBase: dotplugin.QuestionBase{
            ID_:   "my-plugin.extra_config",
            Next_: &dotplugin.Next{End: true}, // resume original flow
        },
        Label:   "Extra config value",
        Default: "default",
    },
}
```

### InjectReplace

Replace a question entirely. The replacement's `Next` chain is used instead of the original question's.

```go
dotplugin.Injection{
    Plugin:      "my-plugin",
    TargetID:    "project_name",
    Kind:        dotplugin.InjectReplace,
    Replacement: &dotplugin.TextQuestion{
        QuestionBase: dotplugin.QuestionBase{
            ID_:   "my-plugin.project_name",
            Next_: &dotplugin.Next{End: true},
        },
        Label:    "Project name (with extra validation)",
        Validate: myStrictValidator,
    },
}
```

Use Replace sparingly — it breaks other plugins that target the same question ID.

---

## Fragment registry

A `FlowFragment` is a named, context-aware routing function. Fragments let plugins contribute reusable decision points — question chains whose exit edge depends on runtime answers — without hardcoding a specific `*Next` pointer.

```go
// internal/flow/fragment.go
type FlowFragment struct {
    ID      string
    Resolve func(ctx *FlowContext) *Next
}
```

### When to use fragments

Use a fragment when a plugin wants to offer a conditional sub-graph that **other plugins or flows can invoke by name**, rather than by embedding a direct `*Next` pointer. This is useful when the exact destination question is not known at plugin init time.

### Registering a fragment

Fragments are registered on the `FragmentRegistry` that lives on the `Runtime` and is injected into the `FlowEngine`:

```go
// In your plugin's setup (or a flow's init):
fr := flow.NewFragmentRegistry()
fr.Register(&flow.FlowFragment{
    ID: "my-plugin.choose_output",
    Resolve: func(ctx *flow.FlowContext) *flow.Next {
        if fast, _ := ctx.Answers["my-plugin.fast_mode"].(bool); fast {
            return &flow.Next{Question: lightweightQuestion}
        }
        return &flow.Next{Question: fullQuestion}
    },
})
```

### Using a fragment from a question graph

A fragment's resolved `*Next` can be returned from any `Question.Next()` implementation. The engine calls `FragmentRegistry.Resolve(id, ctx)` and follows the returned edge.

> **Current status:** The fragment registry is wired through the Runtime and engine but no built-in flow currently uses it. It is ready for plugin authors who need runtime-conditional routing that cannot be expressed with `IfQuestion`.

---

## ResolveExtras

`ResolveExtras(spec)` is called after the flow's own `Generators` resolver has run. Return a slice of additional `Invocation` values to add to the pipeline.

This is the right place to express "if the user picked X, also run my generator":

```go
func (p *MyPlugin) ResolveExtras(s *dotplugin.ProjectSpec) []dotplugin.Invocation {
    // Check an answer the plugin injected or one from the core flow:
    if v, _ := s.Answers["my-plugin.strict_mode"].(bool); v {
        return []dotplugin.Invocation{{Name: "my-plugin.strict_gen"}}
    }
    return nil
}
```

The engine appends these invocations to the flow's invocations before the topological sort. They follow the same `DependsOn` / `ConflictsWith` rules.

---

## pkg/dotplugin reference

External plugins import two packages:

```go
import (
    "github.com/version14/dot/pkg/dotapi"     // Generator, Context, Manifest, Logger
    "github.com/version14/dot/pkg/dotplugin"  // Provider, Question types, Injection types
)
```

`pkg/dotplugin` re-exports every internal type a plugin needs. Never import from `internal/*` — the Go toolchain rejects cross-module internal imports.

**Key re-exports**

| Symbol | Source package | Notes |
|--------|---------------|-------|
| `PluginID` | `internal/flow` | `type PluginID string` |
| `Provider` | `internal/plugin` | The interface to implement |
| `RegisterBuiltin` | `internal/plugin` | Call from `init()` |
| `Question`, `TextQuestion`, `ConfirmQuestion`, … | `internal/flow` | Full question DSL |
| `Injection`, `InjectionKind` | `internal/flow` | Injection declarations |
| `InjectReplace`, `InjectAddOption`, `InjectInsertAfter` | `internal/flow` | Kind constants |
| `Entry`, `Invocation` | `internal/generator` | For `Generators()` and `ResolveExtras()` |
| `ProjectSpec`, `ProjectMetadata` | `internal/spec` | For `ResolveExtras` parameter |
| `VirtualProjectState`, `JSONDoc`, `YAMLDoc`, `GoMod` | `internal/state` | For generators |
| `ContentRaw`, `ContentJSON`, `ContentYAML`, `ContentGoMod` | `internal/state` | Content type constants |

---

## Scaffolding a plugin repository

The fastest way to start a plugin is to scaffold it with DOT itself:

```bash
dot scaffold plugin-template
```

The `plugin-template` flow asks for:

| Question | Notes |
|----------|-------|
| Plugin id | Lowercase, no dots. Used as the namespace prefix. |
| Go module path | e.g. `github.com/you/dot-plugin-my-stack` |
| One-line description | Shown in `dot plugin list` |
| Author name | Used in LICENSE |
| Copyright year | Used in LICENSE |
| Include sample injection? | Generates InsertAfter + AddOption examples |
| Include sample generator? | Generates a minimal Generator + Manifest |

The result is a ready-to-publish Go module with:

```
dot-plugin-my-stack/
├── plugin.json
├── plugin.go           ← init() + Provider implementation
├── go.mod
├── go.sum
├── README.md
└── LICENSE
```

---

## Testing a plugin locally

### Install from a local path

```bash
dot plugin install -from ./dot-plugin-my-stack
```

This copies the directory into `~/.dot/plugins/my-stack/`.

### Rebuild dot with the plugin

Because plugins use Go's `init()` mechanism, you need to import the plugin's package in `cmd/dot/main.go` and rebuild:

```go
import _ "github.com/version14/dot/plugins/my_stack"
```

Then:

```bash
make build
./bin/dot scaffold
```

### Writing test fixtures

Add a JSON fixture to `tools/test-flow/testdata/` that includes answers for any questions your plugin injects. See [test-flow.md](test-flow.md).

---

## Publishing a plugin

1. Push your plugin repository to GitHub (or any git host).
2. Tag a release: `git tag v0.1.0 && git push --tags`.
3. Update `plugin.json` to match the tag version.
4. Users can then install it:

```bash
dot plugin install github.com/you/dot-plugin-my-stack@v0.1.0
```

---

## Installing a plugin

End users install plugins with `dot plugin install`. See [getting-started.md](../user/getting-started.md#manage-plugins) for the full reference.

---

## In-tree vs installed plugins

| Aspect | In-tree (`plugins/`) | Installed (`~/.dot/plugins/`) |
|--------|----------------------|-------------------------------|
| Location | `plugins/<id>/` in the DOT repo | `~/.dot/plugins/<id>/` |
| Registration | Blank import in `cmd/dot/main.go` | `plugin.Load()` at startup |
| Rebuild required | Yes (it's part of the binary) | Yes (currently; dynamic loading is planned) |
| Intended for | Official / bundled extras | Community plugins |
| Example | `plugins/biome_extras` | `examples/example-plugin` |

In-tree plugins import only `pkg/dotapi` and `pkg/dotplugin` — exactly as community plugins do. The in-tree vs. installed distinction is purely about where the code lives and how it is imported, not about what APIs are available.

---

## Reference implementations

| Plugin | What it shows | Docs |
|--------|--------------|------|
| `plugins/biome_extras` | Full contract: InsertAfter injection, generator, ResolveExtras gating | [docs/plugins/biome_extras.md](plugins/biome_extras.md) |
| `examples/example-plugin` | Two injection kinds (InsertAfter + AddOption) in one plugin | [docs/plugins/example_plugin.md](plugins/example_plugin.md) |

Read the source alongside the plugin docs — the code is annotated to explain each decision.
