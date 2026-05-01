# Generator: `auth_better_auth_schema`

BetterAuth Drizzle schema. Creates the four tables BetterAuth requires — `user`, `session`, `account`, `verification` — in `src/db/schema/auth.schema.ts` and re-exports them from the schema index.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_better_auth_schema` |
| Version | `0.1.0` |
| Package | `generators/auth_better_auth_schema` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `drizzle_postgres_adapter` | Drizzle + Postgres must be configured; the schema index at `src/db/schema/index.ts` must exist |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/db/schema/auth.schema.ts` | `user`, `session`, `account`, `verification` tables matching the BetterAuth schema spec |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `src/db/schema/index.ts` | Appends `export * from './auth.schema';` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/db/schema/auth.schema.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.

## See also

- [auth_better_auth.md](./auth_better_auth.md) — BetterAuth server config that references these tables
