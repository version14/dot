# Generator: `auth_clerk_frontend`

Integrates Clerk authentication — framework-aware: `@clerk/nextjs` for Next.js, `@clerk/clerk-react` for Vite. Writes a provider, a `useAuth` hook, and a `.env.example`.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_clerk_frontend` |
| Version | `0.2.0` |
| Package | `generators/auth_clerk_frontend` |

---

## Dependencies

| Generator | Why |
|-----------|-----|
| `typescript_base` | Requires package.json to merge dependencies into |

---

## Answers consumed

| Key | Type | Required | Notes |
|-----|------|----------|-------|
| `framework` | string | No | `"next"` → installs `@clerk/nextjs`; otherwise `@clerk/clerk-react` |

---

## Files written

| Path | Description |
|------|-------------|
| `src/providers/ClerkProvider.tsx` | Clerk provider wrapper |
| `src/hooks/useAuth.ts` | Hook re-exporting Clerk's `useUser` / `useAuth` |
| `.env.example` | Required env vars (`VITE_CLERK_PUBLISHABLE_KEY` or `NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY`) |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `@clerk/nextjs` or `@clerk/clerk-react` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/hooks/useAuth.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Clerk SDK |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `auth_better_auth_frontend` | Both provide auth provider and useAuth hook |
| `auth_vanilla_frontend` | Both provide auth provider and useAuth hook |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
