# Generator: `biome_config`

Biome formatter and linter configuration. Writes `biome.json` and merges `lint`/`format` scripts into `package.json`. The `biome_extras` plugin extends this generator with a strict-mode overlay.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `biome_config` |
| Version | `0.1.0` |
| Package | `generators/biome_config` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires `package.json` to exist for script merging |

---

## Answers consumed

None. `biome_config` writes the same baseline configuration for every project.

---

## Files written

| Path | Description |
|------|-------------|
| `biome.json` | Biome schema, import organizer, linter (recommended rules), formatter (space, 2-space indent) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `scripts.lint = "biome check ."`, `scripts.format = "biome format --write ."`, `devDependencies.@biomejs/biome` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `biome.json` exists | `file_exists` | — |
| `linter.enabled` in `biome.json` | `json_key_exists` | Linter config present |
| `scripts.lint` in `package.json` | `json_key_exists` | Lint script present |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Deduplicated with `typescript_base` |

## Test commands

| Command | Background | Notes |
|---------|-----------|-------|
| `pnpm install` | No | Deduplicated |
| `pnpm exec biome check .` | No | Lints the entire project |

---

## Conflicts

None.

---

## Plugin extension

The `biome_extras` plugin adds an optional strict-mode overlay on top of this generator. See [docs/plugins/biome_extras.md](../plugins/biome_extras.md).
