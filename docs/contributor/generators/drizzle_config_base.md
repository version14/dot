# Generator: `drizzle_config_base`

Creates `drizzle.config.ts` (pointing at `src/db/schema/index.ts`, outputting to `drizzle/`) and an empty schema barrel file.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `drizzle_config_base` |
| Version | `0.1.0` |
| Package | `generators/drizzle_config_base` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | `package.json` must exist to receive Drizzle deps later |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `drizzle.config.ts` | Drizzle-kit config: PostgreSQL dialect, schema path, output path |
| `src/db/schema/index.ts` | Empty barrel — users export their table definitions here |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `drizzle.config.ts` | `file_exists` | — |
| `src/db/schema/index.ts` | `file_exists` | — |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.

---

## See also

- [`drizzle_typescript_deps`](drizzle_typescript_deps.md) — drizzle-orm + drizzle-kit deps + db scripts
- [`drizzle_postgres_adapter`](drizzle_postgres_adapter.md) — postgres driver + `src/db/index.ts`
