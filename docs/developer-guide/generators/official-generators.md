# Official Generators

Official generators ship with dot and are registered in `cmd/dot/build.go`.

---

## Go generators (`generators/go/`)

| Generator | Name | Language | Modules | Status |
|-----------|------|----------|---------|--------|
| `GoRestAPIGenerator` | `go-rest-api` | `go` | `rest-api` | v0.1 — ships today |
| `GoPostgresGenerator` | `go-postgres` | `go` | `postgres` | planned v0.2 |
| `GoRedisGenerator` | `go-redis` | `go` | `redis` | planned v0.2 |
| `GoAuthJWTGenerator` | `go-auth-jwt` | `go` | `auth-jwt` | planned v0.2 |

---

## Common generators (`generators/common/`)

Language-agnostic generators use `Language() = "*"` and match any project language.

| Generator | Name | Language | Modules | Status |
|-----------|------|----------|---------|--------|
| `GitHubActionsGenerator` | `common-github-actions` | `*` | `github-actions` | planned v0.2 |
| `DockerGenerator` | `common-docker` | `*` | `docker` | planned v0.2 |
| `DockerComposeGenerator` | `common-docker-compose` | `*` | `docker-compose` | planned v0.2 |

---

## What each generator does

### GoRestAPIGenerator

**Files created by `Apply()`:**

| File | What it is |
|---|---|
| `main.go` | HTTP server entry point with a `/health` route |
| `go.mod` | Go module declaration using `github.com/<project-name>` |
| `routes/routes.go` | Route registration stub |

**Commands registered:**

| Command | Args | What it does |
|---|---|---|
| `new route` | `<name>` | Creates `routes/<name>.go` with an HTTP handler function |
| `new handler` | `<name>` | Creates `handlers/<name>.go` with an `http.Handler` struct |

**`RunAction` actions:**

- `rest-api.new-route` — creates a route file in `routes/`
- `rest-api.new-handler` — creates a handler struct in `handlers/`

See `generators/go/rest_api.go` for the full implementation.
