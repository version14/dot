# Generator: `express_error_middleware`

Global Express error-handling middleware. Creates `src/shared/middlewares/error.middleware.ts` and injects the `errorMiddleware` import and `app.use(errorMiddleware)` call into `src/app.ts` as the last middleware before `export default app`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `express_error_middleware` |
| Version | `0.1.0` |
| Package | `generators/express_error_middleware` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | `src/app.ts` must exist to inject the middleware registration |
| `express_shared_errors` | `AppError` must exist at `src/shared/errors/app.error` |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/shared/middlewares/error.middleware.ts` | `errorMiddleware(err, req, res, next)` — maps `AppError` to JSON responses; falls back to HTTP 500 |

Also modifies:

| Path | Change |
|------|--------|
| `src/app.ts` | Prepends the import; appends `app.use(errorMiddleware)` before `export default app` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/shared/middlewares/error.middleware.ts` | `file_exists` | — |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.

## See also

- [express_shared_errors.md](./express_shared_errors.md) — error classes consumed by this middleware
