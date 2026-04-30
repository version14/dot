# Generator: `postgres_env_example`

Appends `DATABASE_URL` (PostgreSQL connection string) to the `.env.example` file created by `express_server_entrypoint`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `postgres_env_example` |
| Version | `0.1.0` |
| Package | `generators/postgres_env_example` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `express_server_entrypoint` | `.env.example` seed file must exist before appending |

---

## Answers consumed

None. Reads `Spec.Metadata.ProjectName` to build the default DB URL.

---

## Files written

| Path | Description |
|------|-------------|
| `.env.example` | Appends `DATABASE_URL=postgresql://postgres:postgres@localhost:5432/<project>` |

---

## Validators

None.

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.
