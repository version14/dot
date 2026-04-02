# FAQ

Frequently asked questions about developing Scaffold CLI.

---

## General

**Q: Where do I start?**
See [Getting Started](../getting-started/README.md) for the full setup guide.

**Q: I found a bug. How do I report it?**
Open a [Bug Report issue](../../../issues/new/choose) using the provided template.

**Q: I want to add a feature. Where do I begin?**
Open a [Feature Request issue](../../../issues/new/choose) first to discuss the idea. See [Adding a Generator](#adding-a-generator) for implementation guidance.

**Q: How does Scaffold CLI work?**
It uses a **generator-based architecture**:
1. User answers survey questions
2. Answers build a project specification (JSON)
3. Independent generators read the spec and produce files
4. Multiple generators can safely modify the same file via merge strategies
5. All files are written to disk

See [Architecture Documentation](../../.claude/ressources/Architecture.md) for details.

---

## Development

**Q: How do I run a specific test?**
```bash
go test -v ./internal/generators -run TestAPIGenerator
```

**Q: How do I debug a generator?**
Add print statements or use a debugger like Delve:
```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/scaffold
```

**Q: Tests are failing locally but passing in CI (or vice versa).**
- Ensure your Go version matches the one in [Prerequisites](../getting-started/README.md#prerequisites)
- Run `go mod tidy && go mod download` to sync dependencies
- Check if your system's temp directory has enough space

**Q: The build fails with module not found errors.**
Pull the latest `main` and run `go mod download` — a dependency may have been updated.

**Q: How do I add a new generator?**
See [Adding a Generator](#adding-a-generator) below.

---

## Adding a Generator

**Q: How do I add a new generator (e.g., Redis caching)?**

1. Create `internal/generators/redis.go`
2. Implement the `Generator` interface:
   ```go
   type RedisGenerator struct{}

   func (g *RedisGenerator) Name() string {
       return "Redis Generator"
   }

   func (g *RedisGenerator) Generate(spec *spec.ProjectSpec) ([]generators.File, error) {
       // Return files for Redis setup
   }
   ```
3. Add it to the generator list in `cmd/scaffold/main.go`
4. Write tests in `internal/generators/redis_test.go`
5. Submit a PR with an example template if needed

See `internal/generators/api.go` for a complete example.

---

## Contributing

**Q: How large should a PR be?**
Aim for PRs that can be reviewed in under 30 minutes. Split larger changes into multiple PRs if possible.

**Q: Do I need to write tests for every change?**
Yes for new features and bug fixes. Documentation-only PRs are exempt.

**Q: Who merges PRs?**
Maintainers merge PRs once they have one approving review and all CI checks are green.

**Q: What's the PR submission checklist?**
- [ ] Tests pass: `go test ./...`
- [ ] Code is formatted: `go fmt ./...`
- [ ] Linter passes: `golangci-lint run ./...`
- [ ] Documentation is updated
- [ ] Commit messages follow conventions (see [Code Style](guidelines/code-style.md))

---

Still stuck? Open a [Discussion](../../../discussions).
