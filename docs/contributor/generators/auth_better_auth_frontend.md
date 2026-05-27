# Generator: `auth_better_auth_frontend`

Integrates Better Auth — writes an auth client, a React context provider, and a `useAuth` hook.

---

## Identity

| Field | Value |
|-------|-------|
| Name | `auth_better_auth_frontend` |
| Version | `0.1.0` |
| Package | `generators/auth_better_auth_frontend` |

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
| `src/lib/authClient.ts` | Better Auth client instance |
| `src/providers/AuthProvider.tsx` | React context provider for session state |
| `src/hooks/useAuth.ts` | Hook returning session and sign-in/out helpers |

Also merges into:

| Path | Keys added / updated |
|------|---------------------|
| `package.json` | `dependencies`: `better-auth^1` |

---

## Validators

| Check | Type | Passes when |
|-------|------|-------------|
| `src/hooks/useAuth.ts` | `file_exists` | File is present after generation |

---

## Post-generation commands

| Command | WorkDir | Notes |
|---------|---------|-------|
| `pnpm install --dangerously-allow-all-builds` | project root | Installs Better Auth |

## Test commands

| Command | Background | Ready delay | Notes |
|---------|-----------|-------------|-------|
| `pnpm exec tsc --noEmit` | No | — | Verifies TypeScript compiles |

---

## Conflicts

| Generator | Reason |
|-----------|--------|
| `auth_clerk_frontend` | Both provide auth provider and useAuth hook |
| `auth_vanilla_frontend` | Both provide auth provider and useAuth hook |

---

## See also

- [docs/flows/frontend.md](../flows/frontend.md)
