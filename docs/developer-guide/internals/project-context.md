# ProjectContext

Implementation: `internal/project/context.go`

---

## What .dot/ contains

Two files, both committed to git:

```
.dot/
├── config.json    ← project spec + available commands
└── manifest.json  ← SHA-256 hash of every generated file at generation time
```

**`config.json` is required for `dot new`.** Without it, dot cannot find which generator handles `"new route"` or what Spec was used at creation time.

**`manifest.json` is required for `dot add module`.** Without it, dot cannot tell whether a file has been user-modified since generation — so it cannot safely inject new content.

Both files must be committed to git alongside source code.

---

## config.json

```go
type Context struct {
    DotVersion  string                // "0.1.0" — semver of dot that created this
    SpecVersion int                   // schema version (1 for v0.1)
    Spec        spec.Spec             // the full Spec used at dot init time
    Commands    map[string]CommandRef // key = CommandDef.Name, e.g. "new route"
}

type CommandRef struct {
    Generator string // matches Generator.Name()
    Action    string // passed to generator.RunAction()
}
```

Example `config.json` for a Go REST API project:

```json
{
  "dot_version": "0.1.0",
  "spec_version": 1,
  "spec": {
    "project": { "name": "my-api", "language": "go", "type": "api" },
    "modules": [{ "name": "rest-api" }],
    "config": { "linter": "golangci-lint", "ci": "github-actions" }
  },
  "available_commands": {
    "new route":   { "generator": "go-rest-api", "action": "rest-api.new-route" },
    "new handler": { "generator": "go-rest-api", "action": "rest-api.new-handler" }
  }
}
```

**`Load(startDir)`** traverses from `startDir` up to the nearest `.git` root looking for `.dot/config.json`. This is what makes `dot new` work from any subdirectory of the project.

**`Save(root, ctx, manifest)`** writes both files to `<root>/.dot/`. Called by `dot init` after the pipeline completes successfully.

---

## manifest.json

```go
type Manifest struct {
    Files map[string]FileRecord // path → record
}

type FileRecord struct {
    Hash      string // "sha256:<hex>"
    Generator string // "go-rest-api"
}
```

Example:

```json
{
  "files": {
    "main.go":          { "hash": "sha256:abc123...", "generator": "go-rest-api" },
    "routes/routes.go": { "hash": "sha256:def456...", "generator": "go-rest-api" },
    "go.mod":           { "hash": "sha256:789ghi...", "generator": "go-rest-api" }
  }
}
```

Only `Create` and `Template` ops are tracked — files that dot fully owns. `Append` and `Patch` ops are not tracked because they modify shared files.

---

## dot new dispatch flow

Step by step for `dot new route UserController`:

1. `project.Load(".")` — traverse up from `$PWD`, find `.dot/config.json`
2. Key = `"new route"`, args = `["UserController"]`
3. `ctx.Commands["new route"]` → `{Generator: "go-rest-api", Action: "rest-api.new-route"}`
4. `registry.Get("go-rest-api")` → the generator instance
5. `generator.RunAction("rest-api.new-route", ["UserController"], ctx.Spec)` → `[]FileOp`
6. `pipeline.Run(ops)` → `routes/UserController.go` written to disk
