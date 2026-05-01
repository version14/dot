# Express Backend Generator Guide

This guide documents the conventions, design rules, and pitfalls discovered while building the Express/Node.js generator family. Read it before adding a new Express generator, changing the monorepo flow, or debugging a failing test-flow run.

---

## Table of Contents

- [Generator family overview](#generator-family-overview)
- [Design rules for Express generators](#design-rules-for-express-generators)
- [How the monorepo flow is structured](#how-the-monorepo-flow-is-structured)
  - [Flow branching rules](#flow-branching-rules)
- [Writing tests for generated code](#writing-tests-for-generated-code)
  - [Test types](#test-types)
  - [vitest configuration rules](#vitest-configuration-rules)
  - [Database test rules](#database-test-rules)
- [TestCommands conventions](#testcommands-conventions)
- [Common pitfalls](#common-pitfalls)

---

## Generator family overview

The Express backend is scaffolded by a set of composable generators that run in topological order. Each generator owns one concern.

```
Base
  typescript_base
    express_server_typescript_deps
      express_server_entrypoint
        express_node_tsconfig
        express_shared_errors
          express_error_middleware
          express_auth_validators
        express_rate_limit
        express_test_setup         ← vitest config + test deps
        backend_architecture_mvc   ← or clean_architecture
        [biome_config | prettier_*]

Database (optional)
  postgres_docker_compose
  postgres_env_example
  drizzle_typescript_deps
    drizzle_config_base
      drizzle_postgres_adapter

Auth (requires database)
  auth_jwt_vanilla               ← signToken / signRefreshToken / verifyToken
    auth_jwt_users_schema        ← users + refresh_tokens tables
      auth_jwt_mvc_route         ← controller + routes + test files (MVC path)
      auth_jwt_clean_arch_module ← use-cases + infra + controller (clean-arch path)
  auth_better_auth               ← better-auth alternative
    auth_better_auth_schema
```

Each generator's full manifest lives in its own `manifest.go`. Its purpose is documented in `docs/contributor/generators/<name>.md`.

---

## Design rules for Express generators

**1. One concern per generator.**  
Each generator should handle one concern. If 2 generator need the same generation it should be splitted in a unique generator.

**2. Use `ctx.PreviousGens` to decide, not `ctx.Answers` alone.**  
When you need to know "is drizzle available?", check `slices.Contains(ctx.PreviousGens, "drizzle_postgres_adapter")`. This is more reliable than checking an answer key, because a dependency might have been pulled in transitively without a matching answer.

**3. Never overwrite another generator's file — inject instead.**  
When a generator needs to add an import or middleware call to `src/app.ts`, it reads the current content, inserts at a known anchor string, and writes it back. The anchor is a unique line that the upstream generator guarantees to produce (e.g. `app.use(errorMiddleware)` or `export default app;`).

```go
if f, ok := ctx.State.GetFile("src/app.ts"); ok {
    content := string(f.Content)
    content = strings.Replace(content, "export default app;", myUse+"\nexport default app;", 1)
    ctx.State.WriteFile("src/app.ts", []byte(content), state.ContentRaw)
}
```

**4. Guard injection with an idempotency check.**  
Before injecting into a file, check that the thing you're adding is not already there. This allows `Generate()` to be called safely in test runs that replay the generator.

```go
if !strings.Contains(content, "authRouter") {
    content = authRouteImport + content
    // ...
}
```

**5. Export all schema tables from `src/db/schema/index.ts`.**  
drizzle-kit reads the schema from `index.ts`. Every generator that adds a table must also export it. Strip the placeholder `export {};` before appending the real export, because some tools treat `export {};` as a separate empty-export statement that can interact badly with `export * from`.

```go
existing = strings.ReplaceAll(existing, "export {};\n", "")
existing = strings.ReplaceAll(existing, "export {};", "")
updated := existing + "export * from './my-table';\n"
```

**6. Use proper typed errors instead of `new Error(...)`.**  
Use-cases and controllers must throw from the error hierarchy defined in `express_shared_errors`, not bare `Error`. The error middleware maps these to HTTP status codes.

| Class | HTTP |
|-------|------|
| `NotFoundError` | 404 |
| `UnauthorizedError` | 401 |
| `ForbiddenError` | 403 |
| `ConflictError` | 409 |
| `ValidationError` | 422 |

**7. Refresh tokens must carry a `jti` claim.**  
`signRefreshToken` must add `jti: crypto.randomUUID()` to the payload. Without it, two calls within the same second (e.g. register then login in a test) produce identical JWTs and violate the `refresh_tokens_token_unique` constraint.

---

## How the monorepo flow is structured

The monorepo flow lives in `flows/monorepo.go`. The question graph for the Express backend looks like this:

```
projectName
  → packageManager
    → backend (yes/no)
      → backendLanguage (typescript)
        → backendArchitecture (mvc | clean-architecture)
          → formatter (biome | prettier)
            → enableDb (yes/no)
              yes → dbType (postgres)
                     → dbORM (drizzle)
                       → enableAuth (yes/no)
                         yes → authMethod (jwt | better-auth)
                                → confirmGenerate
                         no  → confirmGenerate
              no  → confirmGenerate   ← skips auth entirely
```

### Flow branching rules

**No database → no auth.** Auth generators require a database. When `enableDb` is answered "no", the `Else` branch points directly to `confirmGenerate`. Never ask about auth when there is no DB.

```go
enableDb := &flow.ConfirmQuestion{
    // ...
    Then: &flow.Next{Question: dbType},
    Else: &flow.Next{Question: confirmGenerate}, // ← NOT enableAuth
}
```

**MVC and clean-architecture reach the same subsequent questions.** Both architecture options converge on the same `formatter` question after selection. Do not add any questions that appear only in one path unless they are genuinely architecture-specific.

**Build the graph bottom-up.** Declare terminal questions first so earlier questions can reference them by pointer. Declaring in the wrong order produces a compile error or a nil-pointer panic at traversal time.

---

## Writing tests for generated code

### Test types

Every generated project uses four test types, distinguished by file-name suffix and run separately:

| Suffix | Vitest command | What it tests | DB needed |
|--------|---------------|---------------|-----------|
| `.unit.test.ts` | `vitest run unit` | Pure logic: error classes, JWT helpers, Zod validators | No |
| `.feature.test.ts` | `vitest run feature` | Mini Express app assembled in-test with supertest | No |
| `.e2e.test.ts` | `vitest run e2e` | Full `app` from `src/app.ts` via supertest | No (lazy DB conn) |
| `.db.test.ts` | `vitest run db` | Real Postgres via Docker; full lifecycle tests | Yes |

The `include` pattern in `vitest.config.ts` covers all four:

```typescript
include: ['src/**/*.{unit,feature,e2e,db}.test.ts'],
```

The `passWithNoTests: true` option is required — not all generated projects produce tests of every type.

### vitest configuration rules

**`fileParallelism: false` is required.**  
By default vitest runs test files in parallel. Multiple `.db.test.ts` files share the same real Postgres instance and the same `db` singleton; running them concurrently causes lock contention and test timeouts. `fileParallelism: false` serializes file execution within each `vitest run` invocation.

**Set JWT variables in `test.env`, not in `.env`.**  
vitest's `test.env` sets `process.env` variables *before module loading*, which means they override dotenv. JWT tests need `JWT_SECRET` set before `src/shared/services/jwt.ts` is imported.

```typescript
env: {
  JWT_SECRET: 'test-secret-vitest',
  JWT_EXPIRES_IN: '1h',
  JWT_REFRESH_EXPIRES_IN: '7d',
},
```

**Do not add `DATABASE_URL` to `test.env`.**  
`DATABASE_URL` comes from `.env`, which is populated by `.env.example`. Setting it in `test.env` would hardcode a DB name that may not match the Docker container. The e2e `TestCommand` copies `.env.example` to `.env` before running vitest.

### Database test rules

**Use Docker, not a local postgres.**  
Database tests spin up a fresh Docker container using the project's own `docker-compose.yml`. They do not assume anything is running locally. This makes CI and dev environments identical.

**Clean volumes before each test run.**  
Docker named volumes persist across `docker compose down` (volumes are not removed by default). Stale volumes can contain schema state from a previous run that conflicts with the current schema. Always run `docker compose down -v` before starting the test container.

**TestCommand sequence in `auth_jwt_users_schema`:**
```
test -f .env || cp .env.example .env      # idempotent: never clobber an existing .env
docker compose down -v 2>/dev/null || true # clean stale volume
docker compose up -d && sleep 5           # start fresh container
pnpm exec drizzle-kit push --force        # apply schema
bash -c 'pnpm exec vitest run db; EXIT_CODE=$?; docker compose down -v; exit $EXIT_CODE'
```

The final command captures the exit code so docker compose cleanup runs even when tests fail.

**`cp -n` is not portable.**  
On macOS, `cp -n` (no-clobber) exits with status 1 when it skips copying. Use `test -f .env || cp .env.example .env` instead, which exits 0 in both cases.

**`sleep 5` after `docker compose up -d`.**  
PostgreSQL needs a few seconds to finish initializing. 4 seconds was observed to be insufficient on some machines; 5 is the safe minimum. Do not use a loop health-check in a `TestCommand` since the command output is not streamed live.

---

## TestCommands conventions

TestCommands run in topological generator order, which means commands from earlier generators run before commands from later ones. Design for this:

- `express_test_setup` (position ~9) runs unit/feature/e2e tests before any DB generator runs
- `auth_jwt_users_schema` (position ~19) starts Docker and runs db tests last

The e2e tests must work without a running database. postgres.js creates a lazy connection — the client object is safe to instantiate with a placeholder URL at module load time; queries only fail if no server is reachable.

The e2e `TestCommand` in `express_test_setup` copies `.env.example` to `.env` first so postgres.js gets a real `DATABASE_URL` at module load time (even though no DB is actually running):

```
test -f .env.example && cp -n .env.example .env; pnpm exec vitest run e2e
```

---

## Common pitfalls

| Symptom | Root cause | Fix |
|---------|-----------|-----|
| `3D000: database does not exist` | Docker container not running or `DATABASE_URL` points to wrong DB name | Ensure `docker compose up -d` runs before `drizzle-kit push`; check that `DATABASE_URL` matches `POSTGRES_DB` in docker-compose |
| `duplicate key value violates unique constraint "refresh_tokens_token_unique"` | `signRefreshToken` called twice in the same second → same JWT payload → same token string | Add `jti: crypto.randomUUID()` to the refresh token payload |
| DB test timeout (5000ms) | Two `.db.test.ts` files ran in parallel; one held a lock while the other waited | Add `fileParallelism: false` to `vitest.config.ts` |
| `cp -n .env.example .env` exits 1 | macOS `cp -n` returns non-zero when the target already exists | Use `test -f .env || cp .env.example .env` |
| `docker compose up -d && sleep 4` passes but DB not ready | Postgres init takes longer than 4s on slower machines | Use `sleep 5` |
| Port 5432 conflict with local postgres | Docker cannot bind the port; drizzle-kit push connects to the local instance instead | Project uses port 5433 → 5432 in docker-compose to avoid conflict |
| `refresh_tokens` table not created by `drizzle-kit push` | `index.ts` still contains `export {};` alongside `export * from './users.table'` — some versions of drizzle-kit skip the table | Strip `export {};` from `index.ts` before appending table exports |
| E2E tests fail with "DATABASE_URL not set" | `src/db/index.ts` throws on load when `DATABASE_URL` is absent | Copy `.env.example` to `.env` before `vitest run e2e` in `express_test_setup` TestCommands |
| Flow ends early when choosing MVC or prettier | `Else` branch of an earlier question pointed to `confirmGenerate` instead of the shared next question | Trace the question graph: both MVC and clean-arch must reach the same `formatter` question; both biome and prettier must reach the same `enableDb` question |
| Auth questions appear even when no DB is selected | `enableDb.Else` incorrectly pointed to `enableAuth` | Point `enableDb.Else` directly to `confirmGenerate` |
