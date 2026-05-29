# Generator: `panda_css`

Sets up Panda CSS v1 with a config file. Registers `panda codegen` as a `prepare` npm script so the styled-system is regenerated on every `pnpm install`. Also adds `styled-system/` to `.prettierignore` (and to `biome.json` files.ignore when Biome is the linter).

---

## Identity

| Field | Value |
|-------|-------|
| Name | `panda_css` |
| Version | `0.2.0` |
| Package | `generators/panda_css` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `linter` | string | No | Init flow linter key — `"biome"` adds `styled-system/**` to biome.json |
| `frontend-linter` | string | No | Frontend flow linter key — same effect as `linter` |

---

## Files written

| Path | Description |
|------|-------------|
| `panda.config.ts` | Panda CSS configuration |
| `postcss.config.mjs` | PostCSS config pointing to Panda's PostCSS plugin |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `devDependencies`: `@pandacss/dev`; `scripts.prepare`: `panda codegen` (appended) |
| `.prettierignore` | Appends `styled-system/` |
| `biome.json` | `files.ignore`: `styled-system/**` (biome only) |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `panda.config.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Panda CSS; triggers `prepare` → `panda codegen` |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `tailwind_v4` | Conflicting utility-first CSS systems |
| `shadcn_ui` | shadcn manages its own CSS setup |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
