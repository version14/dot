# Generator: `tanstack_query`

Installs TanStack Query v5, writes a `QueryProvider` wrapper component with devtools, and injects it into `src/main.tsx`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `tanstack_query` |
| Version | `0.1.0` |
| Package | `generators/tanstack_query` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `react_app` | Requires `src/main.tsx` and `package.json` to exist |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/providers/query-provider.tsx` | `QueryProvider` component wrapping children in `QueryClientProvider` with `ReactQueryDevtools` |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@tanstack/react-query^5`, `@tanstack/react-query-devtools^5` |
| `src/main.tsx` | Injects `QueryProvider` import and wraps StrictMode children |

The `src/main.tsx` injection uses string replacement and is compatible with all router variants (no router, `react_router_v7`, `tanstack_router`).

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/providers/query-provider.tsx` | `file_exists` | File is present after generation |
| `dependencies.@tanstack/react-query` in `package.json` | `json_key_exists` | Dependency entry exists |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs TanStack Query packages |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |
| `pnpm exec vite build` | No | — | Verifies production build succeeds |

---

## Conflicts

None.

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
