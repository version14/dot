# Generator: `auth_vanilla_frontend`

Scaffolds a custom JWT auth layer from scratch — in-memory token management, a React context provider, and a `useAuth` hook. No third-party auth library.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_vanilla_frontend` |
| Version | `0.1.0` |
| Package | `generators/auth_vanilla_frontend` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires the base project structure |

---

## Answers consumed

None.

---

## Files written

| Path | Description |
|------|-------------|
| `src/lib/auth.ts` | Token storage and HTTP helpers (in-memory, no external deps) |
| `src/providers/AuthProvider.tsx` | React context provider for auth state |
| `src/hooks/useAuth.ts` | Hook returning user state and auth actions |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/hooks/useAuth.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs base deps |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `auth_clerk_frontend` | Both provide auth provider and useAuth hook |
| `auth_better_auth_frontend` | Both provide auth provider and useAuth hook |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
