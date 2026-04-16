# FileOp Reference

Implementation: `internal/generator/fileop.go`

---

## FileOp struct

```go
type FileOp struct {
    Kind      FileOpKind // what operation to perform
    Path      string     // relative path from project root, e.g. "routes/users.go"
    Content   string     // file content, template source, or content to append/patch
    Anchor    string     // for Patch ops: which anchor to target
    Generator string     // name of the generator that produced this op
    Priority  int        // conflict resolution order — higher wins
}
```

`Path` is always relative to the project root. Do not use absolute paths or `..` traversal.

`Generator` should match `Generator.Name()`. The pipeline uses it for conflict error messages.

---

## FileOpKind values

### Create

Write a new file with the literal content of `Content`.

If another generator also emits a `Create` op for the same `Path`:
- Higher priority wins; the lower-priority op is skipped.
- Same priority → pipeline aborts before any writes, naming both generators.

Use `Create` for files your generator fully owns — `main.go`, `go.mod`, module-specific config files.

```go
generator.FileOp{
    Kind:      generator.Create,
    Path:      "main.go",
    Generator: g.Name(),
    Priority:  0,
    Content:   "package main\n\nfunc main() {}\n",
}
```

### Template

Same as `Create`, but `Content` is a Go `text/template` string. The template is rendered before writing.

The template receives `nil` as its data value in v0.1. To pass Spec data to a template, construct the template content string in `Apply()` using `fmt.Sprintf` or string concatenation rather than relying on template data.

Conflict rules are identical to `Create`.

```go
generator.FileOp{
    Kind:      generator.Template,
    Path:      "config/app.go",
    Generator: g.Name(),
    Priority:  0,
    Content:   "package config\n\nconst AppName = \"{{.Name}}\"\n",
}
```

### Append

Add `Content` to the end of an existing file. If the file does not yet exist in the pipeline's memory (not on disk — in the in-memory accumulator), the content becomes the file's initial content.

No conflict — all `Append` ops for the same path are applied in priority order. Multiple generators can append to the same file.

Use `Append` for registration lists, route tables, and other files where multiple generators each contribute a few lines.

```go
generator.FileOp{
    Kind:      generator.Append,
    Path:      "Makefile",
    Generator: g.Name(),
    Content:   "\n.PHONY: migrate\nmigrate:\n\tgoose up\n",
}
```

### Patch

Insert `Content` at a named anchor point inside an existing file.

Requires `Anchor` to be one of the defined anchor constants. See [patch-strategies.md](patch-strategies.md) for supported anchors and their constraints.

No conflict — all `Patch` ops are applied in order.

```go
generator.FileOp{
    Kind:      generator.Patch,
    Path:      "main.go",
    Anchor:    generator.AnchorImportBlock,
    Generator: g.Name(),
    Content:   "\"database/sql\"",
}
```

---

## Priority guidelines

| Priority | Use case |
|---|---|
| `10` | Framework-level ops that must win over everything (rare) |
| `5` | Default for official generators |
| `0` | Optional additions, lowest precedence |

Two generators at the same priority emitting `Create` ops for the same path → pipeline error. If you expect another generator might target the same file, use a higher priority or switch to `Append`/`Patch`.
