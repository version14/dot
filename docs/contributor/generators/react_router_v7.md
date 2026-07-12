# Generator: `react_router_v7`

Adds React Router v7 to a React+Vite app — wires up a `RouterProvider`, creates the router definition, and adds a Home page.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `react_router_v7` |
| Version | `0.4.0` |
| Package | `generators/react_router_v7` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `react_app` | Requires the Vite+React base (src/main.tsx to overwrite) |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/main.tsx` | Overwritten — mounts `RouterProvider` |
| `src/router.tsx` | Router definition with a single `/` route |
| `src/pages/Home.tsx` | Home page component |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `react-router`, `react-router-dom` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/router.tsx` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs React Router packages |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `tanstack_router` | Both rewrite src/main.tsx as the router entry point |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
