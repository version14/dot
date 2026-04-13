# FAQ — dot

---

## General

**Q: Where do I start?**
See [Getting Started](../getting-started/README.md) for the full setup guide.

**Q: I found a bug. How do I report it?**
Open a [Bug Report issue](../../../issues/new/choose) using the provided template.

**Q: I want to add a feature. Where do I begin?**
Open a [Feature Request issue](../../../issues/new/choose) first to discuss the idea. See [Adding a Generator](#adding-a-generator) below for implementation guidance.

**Q: How does dot work?**
1. `dot init` launches a huh TUI survey → user choices become a typed `Spec`
2. `Registry.ForSpec` finds generators matching the spec's language + modules
3. Each generator's `Apply(spec)` returns `[]FileOp` (create, template, append, patch)
4. The pipeline collects all ops, resolves conflicts, and writes atomically to disk
5. `.dot/config.json` and `.dot/manifest.json` are written to the project root

See [getting-started/README.md](../getting-started/README.md#project-structure) for a structural overview.

---

## Development

**Q: What's the easiest way to build and run?**
```bash
make dev     # Build and run with colored output
make run     # Run without building
make build   # Just build to bin/dot
```

**Q: How do I run a specific test?**
```bash
go test -v ./internal/generator/... -run TestRegistryForSpec
go test -v ./internal/pipeline/... -run TestPatchImportBlock
```

Or run everything:
```bash
make test
```

**Q: How do I debug a generator?**
Add print statements or use Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/dot
```

**Q: Tests are failing locally but passing in CI (or vice versa).**
- Ensure your Go version matches the `go` directive in `go.mod`
- Run `go mod tidy && go mod download` to sync dependencies

**Q: The build fails with module not found errors.**
```bash
go mod tidy
go mod download
```

**Q: How do I validate all my changes before submitting a PR?**
```bash
make validate
```
Runs in sequence: formatting → vet → lint → tests.

---

## Adding a Generator

**Q: How do I add a new generator (e.g., Redis)?**

1. Create `generators/go/redis.go`
2. Implement the `Generator` interface from `internal/generator/generator.go`:
   ```go
   package gogen

   import (
       "github.com/version14/dot/internal/generator"
       "github.com/version14/dot/internal/spec"
   )

   type GoRedisGenerator struct{}

   func (g *GoRedisGenerator) Name() string      { return "go-redis" }
   func (g *GoRedisGenerator) Language() string  { return "go" }
   func (g *GoRedisGenerator) Modules() []string { return []string{"redis"} }

   func (g *GoRedisGenerator) Apply(s spec.Spec) ([]generator.FileOp, error) {
       // Return FileOps for Redis setup files
   }

   func (g *GoRedisGenerator) Commands() []generator.CommandDef {
       return []generator.CommandDef{
           {
               Name:        "new cache-key",
               Args:        []string{"<name>"},
               Description: "Generate a new Redis cache key helper",
               Action:      "redis.new-cache-key",
               Generator:   "go-redis",
           },
       }
   }

   func (g *GoRedisGenerator) RunAction(action string, args []string, s spec.Spec) ([]generator.FileOp, error) {
       // Handle post-creation commands
   }
   ```
3. Register it in `cmd/dot/build.go`:
   ```go
   func buildRegistry() *generator.Registry {
       reg := &generator.Registry{}
       must(reg.Register(&gogen.GoRestAPIGenerator{}))
       must(reg.Register(&gogen.GoRedisGenerator{}))  // add here
       return reg
   }
   ```
4. Write tests in `generators/go/redis_test.go`

See `generators/go/rest_api.go` for a complete working example.

---

## Distribution

**Q: How do I install dot on a new machine?**

```bash
# Homebrew (macOS/Linux)
brew install version14/tap/dot

# curl (no Go required)
curl -fsSL https://raw.githubusercontent.com/version14/dot/main/install.sh | sh

# go install
go install github.com/version14/dot/cmd/dot@latest
```

**Q: How do I update dot?**
```bash
dot self-update
```

**Q: How do I uninstall dot?**
```bash
# Homebrew
brew uninstall dot

# curl / go install / from source
curl -fsSL https://raw.githubusercontent.com/version14/dot/main/uninstall.sh | sh
```

**Q: How do I cut a new release?**
```bash
git tag v0.2.0
git push origin v0.2.0
```
GoReleaser handles the rest. See [CI_CD.md](../CI_CD.md#how-to-cut-a-release) for details.

---

## Contributing

**Q: How large should a PR be?**
Aim for PRs reviewable in under 30 minutes. Split larger changes.

**Q: Do I need to write tests for every change?**
Yes for new features and bug fixes. Documentation-only PRs are exempt.

**Q: Who merges PRs?**
Maintainers merge once there is one approving review and all CI checks are green.

**Q: What's the PR submission checklist?**
```bash
make validate   # fmt → vet → lint → test
```
Then verify documentation is updated if the change affects user-facing behaviour.

---

Still stuck? Open a [Discussion](../../../discussions).
