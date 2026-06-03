# Generator: `tanstack_router`

Adds TanStack Router (file-based routing) to a React+Vite app — installs the Vite plugin, wires the root route, and creates a typed index route.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `tanstack_router` |
| Version | `0.2.2` |
| Package | `generators/tanstack_router` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `react_app` | Requires the Vite+React base (vite.config.ts and src/main.tsx to overwrite) |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `vite.config.ts` | Overwritten — adds `TanStackRouterVite` plugin |
| `src/main.tsx` | Overwritten — mounts `RouterProvider` with TanStack router |
| `src/routes/__root.tsx` | Root route layout |
| `src/routes/index.tsx` | Index route (renders at `/`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@tanstack/react-router`; `devDependencies`: `@tanstack/router-plugin` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/routes/__root.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs TanStack Router packages |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `react_router_v7` | Both rewrite src/main.tsx as the router entry point |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
