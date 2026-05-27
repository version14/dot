# Generator: `storybook_setup`

Configures Storybook v8 — framework-aware: uses `@storybook/react-vite` for Vite projects and `@storybook/nextjs` for Next.js.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `storybook_setup` |
| Version | `0.1.0` |
| Package | `generators/storybook_setup` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `framework` | string | No | `"next"` → uses `@storybook/nextjs`; otherwise uses `@storybook/react-vite` |

---

## Files written

| Path | Description |
|------|-------------|
| `.storybook/main.ts` | Storybook main config (framework, addons) |
| `.storybook/preview.ts` | Global decorators and parameters |
| `src/stories/Button.stories.tsx` | Example Button story |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `devDependencies`: `storybook`, framework-specific `@storybook/*` adapter, `@storybook/addon-essentials`, `@storybook/addon-interactions` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `.storybook/main.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Storybook packages |

## Test commands

No TestCommands.

---

## Conflicts

None.

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
