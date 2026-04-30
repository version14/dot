# Generator: `drizzle_postgres_adapter`

Adds the `postgres` npm driver to `package.json` and creates `src/db/index.ts` — the Drizzle client singleton that connects to PostgreSQL using `DATABASE_URL`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `drizzle_postgres_adapter` |
| Version | `0.1.0` |
| Package | `generators/drizzle_postgres_adapter` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `drizzle_typescript_deps` | `drizzle-orm` must be declared before adding the driver |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/db/index.ts` | Drizzle singleton: connects via `postgres-js`, exports `db` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies.postgres` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/db/index.ts` | `file_exists` | — |
| `dependencies.postgres` in `package.json` | `json_key_exists` | — |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Deduped |

## Test commands

No TestCommands.

---

## Conflicts

None.
