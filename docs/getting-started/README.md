# Getting Started

This guide walks you through setting up Scaffold CLI for local development.

---

## Prerequisites

| Tool | Version | Install                               |
|------|---------|---------------------------------------|
| go   | 1.26+  | [Install](https://go.dev/doc/install) |
| git  | Latest  | [Install](https://git-scm.com/)       |

---

## Setup

1. **Clone the repository**

   ```bash
   git clone https://github.com/version14/scaffold-cli.git
   cd scaffold-cli
   ```

2. **Activate the commit-msg hook** (one-time, after cloning)

   ```bash
   git config core.hooksPath .githooks
   ```

3. **Download dependencies**

   ```bash
   go mod download
   ```

4. **Run the CLI**

   ```bash
   go run ./cmd/scaffold new
   ```

   This starts an interactive questionnaire that will scaffold a new project.

---

## Project Structure

Here's what you'll work with:

```
scaffold-cli/
├── cmd/scaffold/           # CLI entrypoint
├── internal/
│   ├── survey/            # Interactive questionnaire
│   ├── spec/              # Project specification
│   ├── generators/        # Composable generators
│   ├── template/          # Template rendering
│   └── merge/             # Smart file merging
├── templates/             # Template files
└── go.mod                 # Module definition
```

For details, see the [Architecture Documentation](../../.claude/ressources/Architecture.md).

---

## Common Commands

```bash
# Build the CLI binary
go build -o scaffold ./cmd/scaffold

# Run the interactive CLI
go run ./cmd/scaffold new

# Run all tests
go test ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Format all Go code
go fmt ./...

# Lint code (requires golangci-lint)
golangci-lint run ./...

# Run linter and tests together
go fmt ./... && golangci-lint run ./... && go test ./...
```

---

## Troubleshooting

**Go version mismatch**

Make sure your Go version matches the one listed in [Prerequisites](#prerequisites):

```bash
go version
```

**Dependency issues**

If you encounter dependency problems, try:

```bash
go mod tidy
go mod download
go mod verify
```

**Tests failing**

Run tests with verbose output to see what's failing:

```bash
go test -v ./...
```

**Build errors**

Ensure all dependencies are installed:

```bash
go mod download
go build ./...
```

For other issues, check the [FAQ](../developer-guide/faq.md) or open a [Discussion](../../../discussions).
