# Generator Authoring Guide

This guide walks you through writing a new generator for dot.

---

## Before you start

A generator is the right tool when:
- You want to scaffold files for a new language, framework, or module
- The output is deterministic — same Spec always produces the same files
- The module is structural (route files, config files, Docker setup) rather than business logic

A generator is NOT the right tool for:
- One-off project-specific code
- Runtime behavior or orchestration
- Anything that requires understanding the user's domain logic

When in doubt, ask: would this generator be useful to at least 10 different projects? If yes, it probably belongs in dot.

---

## Step 1: Create the file

**Official generators** (shipping with dot):
```
generators/go/        ← Go generators, package gogen
generators/common/    ← language-agnostic (CI, Docker), package commgen
generators/<lang>/    ← future language generators
```

Naming convention: `<language>_<module>.go`, e.g. `go_redis.go`, `go_postgres.go`.

**Community generators** can live anywhere — implement the same interface and register with the registry.

---

## Step 2: Implement the interface

Here is a complete minimal generator for a Go Redis module:

```go
package gogen

import (
    "fmt"

    "github.com/version14/dot/internal/generator"
    "github.com/version14/dot/internal/spec"
)

type GoRedisGenerator struct{}

func (g *GoRedisGenerator) Name() string      { return "go-redis" }
func (g *GoRedisGenerator) Language() string  { return "go" }
func (g *GoRedisGenerator) Modules() []string { return []string{"redis"} }

func (g *GoRedisGenerator) Apply(s spec.Spec) ([]generator.FileOp, error) {
    return []generator.FileOp{
        {
            Kind:      generator.Create,
            Path:      "internal/cache/redis.go",
            Generator: g.Name(),
            Priority:  0,
            Content:   fmt.Sprintf(`package cache

import "github.com/redis/go-redis/v9"

func New() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })
}
`),
        },
        {
            Kind:      generator.Patch,
            Path:      "main.go",
            Anchor:    generator.AnchorImportBlock,
            Generator: g.Name(),
            Content:   "github.com/redis/go-redis/v9",
        },
    }, nil
}

func (g *GoRedisGenerator) Commands() []generator.CommandDef {
    return []generator.CommandDef{
        {
            Name:        "new cache-key",
            Args:        []string{"<name>"},
            Description: "Generate a new Redis cache key helper",
            Action:      "redis.new-cache-key",
            Generator:   g.Name(),
        },
    }
}

func (g *GoRedisGenerator) RunAction(action string, args []string, s spec.Spec) ([]generator.FileOp, error) {
    if len(args) == 0 {
        return nil, fmt.Errorf("name argument required")
    }
    name := args[0]

    switch action {
    case "redis.new-cache-key":
        return []generator.FileOp{
            {
                Kind:      generator.Create,
                Path:      fmt.Sprintf("internal/cache/%s.go", name),
                Generator: g.Name(),
                Content:   fmt.Sprintf("package cache\n\nconst Key%s = \"%s\"\n", name, name),
            },
        }, nil
    default:
        return nil, fmt.Errorf("unknown action %q", action)
    }
}
```

Walk through the choices:
- `Name()` returns `"go-redis"` — unique, stable, used as a key in `.dot/config.json`
- `Modules()` returns `["redis"]` — the Registry matches this against `spec.Modules[].Name`
- `Apply()` creates a Redis client file and patches `main.go`'s imports
- `Commands()` registers one post-creation command
- `RunAction()` handles the command by creating a new cache key file

---

## Step 3: Register it

Add to `buildRegistry()` in `cmd/dot/build.go`:

```go
func buildRegistry() *generator.Registry {
    reg := &generator.Registry{}
    must(reg.Register(&gogen.GoRestAPIGenerator{}))
    must(reg.Register(&gogen.GoRedisGenerator{})) // add here
    return reg
}
```

`must()` panics on error. A registration conflict (two generators claiming the same language+module) is a programming error — it should fail loudly at startup during development.

---

## Step 4: Write tests

Create `generators/go/redis_test.go`. Test `Apply()` and `RunAction()` exhaustively:

```go
package gogen_test

import (
    "testing"

    "github.com/version14/dot/generators/go"
    "github.com/version14/dot/internal/generator"
    "github.com/version14/dot/internal/spec"
)

func TestGoRedisGeneratorApply(t *testing.T) {
    t.Parallel()
    g := &gogen.GoRedisGenerator{}
    s := spec.Spec{
        Project: spec.ProjectSpec{Name: "my-api", Language: "go"},
        Modules: []spec.ModuleSpec{{Name: "redis"}},
    }
    ops, err := g.Apply(s)
    if err != nil {
        t.Fatalf("Apply: %v", err)
    }
    // Assert expected file paths and kinds
    wantPaths := map[string]generator.FileOpKind{
        "internal/cache/redis.go": generator.Create,
        "main.go":                 generator.Patch,
    }
    for _, op := range ops {
        kind, ok := wantPaths[op.Path]
        if !ok {
            t.Errorf("unexpected op for path %q", op.Path)
            continue
        }
        if op.Kind != kind {
            t.Errorf("path %q: got kind %q, want %q", op.Path, op.Kind, kind)
        }
    }
}
```

Use table-driven tests. No shared mutable state. Always `t.Parallel()`.

---

## Common mistakes

**1. Emitting AnchorImportBlock for unsupported import forms.**
Check `patch-strategies.md` for what forms are safe. If the file might have blank imports or build tags, don't emit the op — return an error instead.

**2. Writing outside your module's concern.**
A Redis generator should not create `main.go` from scratch. It can patch it. It should not create route files. Stay focused on your module's responsibility.

**3. Non-deterministic output.**
Don't use `time.Now()`, `rand`, or map iteration order in your generator code. Same Spec must always produce identical FileOps.

**4. Using Priority=0 on a Create op that another generator also targets.**
If you expect another generator might create the same file, use a higher priority, or switch to `Append`/`Patch`. Check `official-generators.md` to see what paths other generators own.

**5. Forgetting to handle the empty-args case in RunAction.**
`dot new cache-key` with no name argument should return a clear error, not panic.
