# ADR-0002: Dependency Catalog

**Status:** Accepted
**Date:** 2026-07-11
**Supersedes:** [ADR-0001](./0001-template-dep-checker-architecture.md)

## Context

ADR-0001 built `tools/dep-checker` on a principle it called **real-state extraction**:
instantiate each generator with an empty `VirtualProjectState`, call `Generate()`, and read
back whatever `package.json` it wrote. The ADR claimed this was *"strictly accurate"* and
*"can never drift"* from real output.

It is neither. In production the bot produced commit `997e489` — titled
`chore(deps): bump @clerk/clerk-react in templates`, containing 13 files touching ark-ui,
better-auth, vitest, next, react, prettier, vite, posthog and react-router, and **zero clerk
changes**. It was merged to `main`.

Three flaws, each fatal on its own.

**1. Reading back an empty-Answers run is not real state.** Generators branch on
`ctx.Answers`; 15 of them do. Calling `Generate()` with `Answers: map[string]interface{}{}`
executes one arbitrary branch. `auth_clerk_frontend` scaffolds `@clerk/clerk-react` when no
framework is chosen and `@clerk/nextjs` when it is — so `@clerk/nextjs ^7.4.2` was **never
checked once** in the tool's lifetime. The scan did not observe real state; it observed one
default slice of it.

**2. Scanning and patching disagreed about what a dependency is.** The scanner read *values*
out of generated JSON. The patcher had to write back to a *Go string literal in source*, which
it located by regex (`patch.go:96`). Nothing guarantees a value observed in the output exists
as a literal in the source. `auth_clerk_frontend` holds its version in a variable
(`version := "^5.61.3"`), so the patch could never match, and errored every run. The tool
could see a dependency it was structurally incapable of patching.

**3. Failure was not contained.** The patch error set `PATCH_FAILED`, broke the loop, and ran
`git checkout main` **without resetting the working tree** — leaving every generator patched
before clerk dirty on disk. `git add generators/` in the next iteration swept them into an
unrelated commit. That is `997e489`.

The `Deprecated:` fields, the ecosystem table, and the registry checkers all worked. The
failure was entirely in *where versions live*.

Two further defects share the same root. Versions living in 40 hand-written files means the
same package can be pinned twice at different versions, and it was: `vitest` sat at `^4.1.7`
and `^4.1.8`, `bcryptjs` at `^2.4.3` and `^3.0.3` — drift the bot itself caused by dying
mid-run. And `dep-checker` reimplemented semver comparison from scratch
(`checkers.go:46-105`), getting zero-major wrong: `^0.45.2 → 0.46.0` is a **breaking** change
under semver, but it was classified `"minor"` and batched as safe. `drizzle-orm` is pinned at
`^0.45.2`.

## Decision

**Pinned versions move out of generators and into a Catalog: `internal/deps/`.**

There is exactly one Pin per dependency, repo-wide. Generators name the packages they need;
the Catalog supplies the version.

```go
// internal/deps/npm.go — machine-owned; the only file the bot rewrites
var npm = map[string]string{
    "@clerk/clerk-react": "^5.61.3",
    "@clerk/nextjs":      "^7.4.2",
    "vitest":             "^4.1.8",
}

// generators/react_app/generator.go
"dependencies":    deps.NPM("react", "react-dom"),
"devDependencies": deps.NPM("vite", "@types/react"),
```

Ecosystem becomes a property the Pin **declares**, not something inferred from which file a
generator happened to write. Adding Rust means adding a `Registry` implementation, not
teaching a scanner a new filename.

This collapses the failure modes rather than patching them:

- **Scan** is a Go import of the Catalog. No `Generate()`, no Answers, no branches to be blind
  to. `@clerk/nextjs` becomes visible because a Catalog has no `if` statements.
- **Patch** rewrites one machine-owned file with one canonical shape. The value the scanner
  read *is* the literal the patcher writes. The reverse-mapping problem ceases to exist.
- **Divergent pins become unrepresentable.** One key, one version.
- **Classification** delegates to `internal/versioning`, which already implements caret and
  zero-major correctly and is tested (`TestConstraint_Caret_ZeroMajor`). A Pin update is
  low-risk **iff the new version satisfies the existing Pin's constraint** — which is what the
  caret already means. `^0.45.2` does not permit `0.46.0`, so Drizzle is correctly treated as
  breaking, with no special case.

### Proposing changes

| Kind | Vehicle | Ownership |
|---|---|---|
| Satisfies the existing constraint | **Rollup** — one batched PR, all packages | Bot. Force-rebuilt every run. Excludes any package holding an open dep PR of its own. |
| Breaks the existing constraint | **Migration** — one PR per package | Human, from the moment it is raised. Bot never touches it again. |
| Deprecated | Issue only, **never a PR** | Human. |

A deprecated package is excluded from the Rollup entirely. Every deprecation in the repo today
(`@clerk/clerk-react` → `@clerk/react`, `@vercel/flags` → `flags`, `@types/bcryptjs` → delete,
`plausible-tracker` → replace) needs a **rename or a removal**. ADR-0001's rule — *"open the
version-bump PR if a newer version exists under the same name"* — resolves none of them:
bumping to the latest *deprecated* version fixes nothing.

A Rollup whose CI goes red is not bisected by the bot. A human ejects the offending package
into its own Migration PR, and the next Rollup rebuild drops it via the same open-PR check
that protects Migrations. One rule, both cases.

### Manifest bumps are derived at release, not stored in the PR

ADR-0001 bumped `manifest.go` inside the dep PR. That stores a **relative** operation
(`0.8.0 → 0.9.0`) computed against `main` at PR-creation time and applies it at *merge* time,
when `main` has moved. With several dep PRs open, the resulting version is a function of merge
order. Long-lived Migrations make this unsurvivable: a React Migration open for two months
collides with every weekly Rollup that touches the same manifest.

Manifest bumps therefore move to **release**, derived from the final state:

```
dot gen-bump   # every generator whose Fingerprint moved since the last tag, bumped once
```

Merge order becomes provably irrelevant, because the bump is no longer a stored delta. A
Migration PR touches only its own Catalog line, so it rebases cleanly for months — the Rollup
edits *other lines of the same map*.

### Fingerprints

A **Fingerprint** is a hash of a generator's Contribution — the files it introduces — taken
across every fixture that invokes it, computed by diffing `VirtualProjectState` before and
after the generator runs. It is entirely in-memory: no `pnpm install`, milliseconds.

The Fingerprint answers the only question that matters: *did what this generator scaffolds
actually change?* Reformatting source does not move it. Editing a template does. Moving a Pin
the generator names does too.

Three CI rules follow:

1. **PR-time** — if the diff touches `generators/**` and a Fingerprint moved, `manifest.go`
   must be bumped and the doc's version row must match.
2. **Always** — every generator's doc version row equals its `manifest.go` Version.
3. **Release-time** — every generator whose Fingerprint moved since the last tag has a bumped
   manifest and a synced doc.

Rule 1 deliberately **does not fire** on a diff touching only `internal/deps/**`. Enforcing a
manifest bump inside a dep PR is precisely the merge-order bug this ADR exists to remove.

## Alternatives considered

| Option | Reason rejected |
|---|---|
| Keep inline literals; ban variable-held versions with a lint | Fixes the clerk crash, but leaves the scanner blind to conditional deps and the regex patcher load-bearing. Treats the symptom. |
| Run generators under many Answer permutations to flush out hidden deps | Machinery to recover information a Catalog never loses. Enumerating permutations is a combinatorial guess; a Catalog is a list. |
| Per-generator `Dependencies` field in `manifest.go` | ADR-0001 rejected this as *"a second source of truth that can drift"* — correctly, but it drew the wrong conclusion. Drift came from versions living in **40** places. The fix is **one** place, not one-per-generator. |
| Bot bisects a red Rollup automatically | ~5 extra full `test-flows` runs (37 fixtures, real `pnpm install`) per failure, to tell a human what the CI log already says. |
| Allow multiple pinned variants (`bcryptjs@2`, `bcryptjs@3`) | All four divergent pins in the repo today are bugs, not intent. No legacy templates are planned. Add the escape hatch when a use for it exists. |

## Consequences

- **Generators no longer contain versions.** A contributor adding a package adds a Pin to the
  Catalog and names it in the generator. `deps.NPM()` panics on an unknown name, so a typo
  fails at generate time and `test-flows` catches it.
- **One version per package, repo-wide.** If a genuine need to pin an old major appears, it
  requires a deliberate new Catalog key — it cannot happen by accident, which is how the
  current `vitest` and `bcryptjs` divergence happened.
- **Between releases, `main` carries generators whose Pin moved but whose manifest has not.**
  This is the price of order-proof bumps. No user observes it: users consume released
  binaries, and `doctor` compares against the released manifest.
- **The bot never runs `test-flows`.** The PR's own CI already does. ADR-0001's promised gate
  (run `test-flows`, open an issue on failure) would double the most expensive job in CI for
  no added signal. It was never implemented, and it should not be.
- **Registry traffic drops to one request per package**, instead of one per
  (generator, package) pair.
- **Cargo, Maven and Go are dead code today** — no generator writes `Cargo.toml`, `pom.xml`, or
  `go.mod`. The `Registry` interface keeps them cheap to revive; the unreachable extractors go.
- `tools/dep-checker` keeps its own semver logic **nowhere**. `parseSemver`, `isOutdated`,
  `updateType`, `stripConstraintPrefix` and `depBumpIsMajor` are deleted in favour of
  `internal/versioning`.
