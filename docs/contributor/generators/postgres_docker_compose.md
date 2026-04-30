# Generator: `postgres_docker_compose`

Creates `docker-compose.yml` with a PostgreSQL 16 service for the local development environment. Uses the project name as the container name and database name.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `postgres_docker_compose` |
| Version | `0.1.0` |
| Package | `generators/postgres_docker_compose` |

---

## Dependencies

None.

---

## Answers consumed

None. Reads `Spec.Metadata.ProjectName` for container/DB naming.

---

## Files written

| Path | Description |
|------|-------------|
| `docker-compose.yml` | PostgreSQL 16-alpine service, port 5432, named volume for persistence |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `docker-compose.yml` | `file_exists` | — |

---

## Post-generation commands

No PostGenerationCommands.

## Test commands

No TestCommands.

---

## Conflicts

None.
