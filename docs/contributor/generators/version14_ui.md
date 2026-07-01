# Generator: `version14_ui`

Adds `@version14/ui` — the company's WIP component library built with Panda CSS — and writes an example component.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `version14_ui` |
| Version | `0.4.0` |
| Package | `generators/version14_ui` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/components/V14Example.tsx` | Renders `Button` imported from `@version14/ui/button`, demonstrating the per-component subpath import |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@version14/ui@0.8.0`, `@ark-ui/react@^5.37.2` |

`@ark-ui/react` is a required peer of `@version14/ui` (every component wraps an Ark UI primitive), so it's added explicitly here rather than relying on the package manager to auto-install peers.

Also injects `import "@version14/ui/styles.css";` into the app's real entry point — `src/main.tsx` when `react_app` ran, `src/app/layout.tsx` when `nextjs_base` ran. This is done here rather than in `panda_css` because the component stylesheet is required no matter which styling engine (Tailwind, CSS Modules, Panda, or none) the project uses; `panda_css`'s own `injectV14Fonts`/`injectV14Preset` only run when Panda is picked as the styling engine, and only wire the optional font faces and preset, not the component CSS itself.

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `dependencies.@version14/ui` | `json_key_exists` | Key is present in `package.json` |
| `dependencies.@ark-ui/react` | `json_key_exists` | Key is present in `package.json` |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs @version14/ui from npm |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `shadcn_ui` | Both provide a component library |
| `ark_ui` | Both provide a component library |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
