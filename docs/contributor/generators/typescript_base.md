# Generator: `typescript_base`

TypeScript foundation for any project. Writes `package.json` and `tsconfig.json` with strict settings and a `pnpm install` post-gen step. All TypeScript-aware generators depend on this one.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `typescript_base` |
| Version | `0.1.0` |
| Package | `generators/typescript_base` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `base_project` | Must run first to create the project root |

---

## Answers consumed

| Key | Type | Notes |
|-----|------|-------|
| `project_name` | string | Written as the `"name"` field in `package.json`. Falls back to `spec.Metadata.ProjectName`, then `"app"`. |

---

## Files written

| Path | Description |
|------|-------------|
| `package.json` | Sets `name`, `version`, `private`, `type: module`, `scripts.build`, `devDependencies.typescript` |
| `tsconfig.json` | ES2022 target, ESNext module, Bundler resolution, `strict: true`, `esModuleInterop: true`, `outDir: dist`, `include: ["src"]` |

Both files are written via `UpdateJSON` — downstream generators can merge additional keys without conflict.

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `package.json` exists | `file_exists` | — |
| `tsconfig.json` exists | `file_exists` | — |
| `devDependencies.typescript` key in `package.json` | `json_key_exists` | TypeScript is listed as a dev dep |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install` | project root | Installs all dependencies. Deduplicated — only runs once even when `react_app` and `biome_config` also declare it. |

## Test commands

| Command | Background | Notes |
|---------|-----------|-------|
| `pnpm install` | No | Install deps |
| `pnpm exec tsc --noEmit` | No | Type-check |

---

## Conflicts

None.
