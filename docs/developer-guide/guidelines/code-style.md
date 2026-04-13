# Go Code Style Guide — dot

We follow Go idioms and best practices from [Effective Go](https://golang.org/doc/effective_go). Consistency matters more than any individual rule — when in doubt, follow existing patterns in the codebase.

---

## Tooling

```bash
# Run ALL checks in sequence (recommended before every push)
make validate

# Individual checks
make fmt      # Format code (gofmt)
make vet      # Go static analysis
make lint     # golangci-lint
make test     # Tests with race detector
```

**Raw Go commands:**

```bash
go fmt ./...
go vet ./...
golangci-lint run ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

All checks must pass before a PR can be merged.

---

## General Principles

- **Clarity over cleverness** — write code for the next reader
- **Explicit over implicit** — avoid magic; name things for what they do
- **Small functions** — each function does one thing
- **No dead code** — remove commented-out code before committing

---

## Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Files | `snake_case` | `rest_api.go`, `cmd_init.go` |
| Packages | `lowercase` | `spec`, `generator`, `pipeline` |
| Exported types | `PascalCase` | `ProjectSpec`, `FileOp` |
| Exported functions | `PascalCase` | `Apply()`, `ForSpec()` |
| Unexported functions | `camelCase` | `resolveConflicts()`, `findImportBlock()` |
| Constants | `PascalCase` (exported) or `camelCase` (unexported) | `AnchorMainFunc`, `defaultPriority` |
| Interfaces | `PascalCase`, usually ending in `er` | `Generator` |
| Sentinel errors | Start with `Err` | `ErrUnsupportedImportForm`, `ErrNotDotProject` |

---

## Formatting

`gofmt` enforces these automatically:

- **Indentation:** tabs
- **Line length:** no hard limit; ~100 chars when practical
- **Braces:** same line as declaration — `func foo() {`
- **Blank lines:** one between top-level declarations; use sparingly within functions

---

## Import Order

Three groups, alphabetically sorted within each:

```go
import (
    // 1. Standard library
    "encoding/json"
    "fmt"
    "os"

    // 2. Third-party packages
    "github.com/charmbracelet/huh"
    "github.com/charmbracelet/lipgloss"

    // 3. Internal packages
    "github.com/version14/dot/internal/generator"
    "github.com/version14/dot/internal/spec"
)
```

`goimports` formats this automatically: `goimports -w ./...`

---

## Error Handling

- Functions that can fail return `error` as the last return value
- Always check errors immediately: `if err != nil { return err }`
- Never swallow errors silently; use `_ =` with a comment only for deferred closes where the error is genuinely unactionable:
  ```go
  defer func() { _ = f.Close() }()
  ```
- Use sentinel errors for cases the caller needs to branch on:
  ```go
  var ErrUnsupportedImportForm = errors.New("unsupported import form")
  ```
- Wrap errors with context:
  ```go
  return fmt.Errorf("pipeline: write %s: %w", op.Path, err)
  ```
- Validate at system boundaries (CLI args, JSON parsing); trust internal code

---

## Testing Conventions

- Test files live in the same package as the code under test, named `*_test.go`
- Use table-driven tests for multiple scenarios:
  ```go
  tests := []struct {
      name    string
      src     string
      want    string
      wantErr error
  }{
      {"block import — add new pkg", "...", "...", nil},
      {"duplicate skipped",          "...", "...", nil},
  }
  for _, tt := range tests {
      t.Run(tt.name, func(t *testing.T) {
          t.Parallel()
          // ...
      })
  }
  ```
- Mark independent tests and subtests with `t.Parallel()`
- Each test must be independent — no shared mutable state
- Use `t.Helper()` in assertion helpers
- Run with race detector: `go test -race ./...`

---

## Go-Specific Best Practices

**Interfaces:**
- Keep interfaces small (1-3 methods in most cases)
- `Generator` is larger by design — it's the core extensibility point
- Accept interfaces, return concrete types

**Concurrency:**
- The FileOp pipeline is intentionally single-threaded (collect → resolve → write atomically)
- Don't introduce goroutines without a clear need
- If you do, pass `context.Context` for cancellation

**Dependencies:**
- Keep `go.mod` minimal; only add packages you actually use
- Run `go mod tidy` before committing
- Current direct dependencies: `bubbletea`, `bubbles`, `huh`, `lipgloss` (charmbracelet suite)
- No framework for CLI dispatch — plain `os.Args` switch in `cmd/dot/main.go`

**Memory:**
- Use `strings.Builder` for string concatenation in hot paths
- The pipeline collects all FileOps in memory before writing — intentional, enables conflict detection before any disk writes

---

## Running the Full Validation Suite

```bash
make validate
```

Executes in sequence:
1. `go fmt` — format
2. `go vet` — static analysis
3. `golangci-lint` — linting
4. `go test -race` — tests with race detector

**Check coverage:**

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out    # Opens in browser
go tool cover -func=coverage.out    # Prints percentages
```
