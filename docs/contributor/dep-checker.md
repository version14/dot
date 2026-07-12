# Dependency Catalog & dep-checker

Every third-party package that `dot` scaffolds into a generated project is pinned in
one place: the **Catalog**, at `internal/deps/`. `tools/dep-checker` keeps the Catalog
in step with the package registries.

> Superseded design: this replaces the generator-scanning checker described in
> [ADR-0001](../adr/0001-template-dep-checker-architecture.md). Read
> [ADR-0002](../adr/0002-dependency-catalog.md) for why.

---

## The Catalog

```go
// internal/deps/npm.go — machine-owned; dep-checker rewrites this and nothing else
var npm = map[string]string{
    "@clerk/clerk-react": "^5.61.3",
    "@clerk/nextjs":      "^7.4.2",
    "react":              "^19.2.7",
    "vitest":             "^4.1.8",
    // ...
}
```

Generators **name** the packages they need. The Catalog says which version they get:

```go
// generators/react_app/generator.go
d.Merge(map[string]interface{}{
    "dependencies":    deps.NPM("react", "react-dom"),
    "devDependencies": deps.NPM("vite", "@types/react"),
})
```

`deps.NPM()` panics on an unknown package, so a typo fails when the generator runs and
`test-flows` catches it.

**There is exactly one Pin per package, repo-wide.** Two generators cannot pin the same
package to different versions — that is unrepresentable, not merely discouraged. It used
to happen: `vitest` sat at `^4.1.7` in one generator and `^4.1.8` in another, and
`bcryptjs` at `^2.4.3` and `^3.0.3`, because the old bot patched some generators and
died before reaching the rest.

### Adding a package

1. Add the Pin to `internal/deps/npm.go`.
2. Name it from a generator: `deps.NPM("your-package")`.

That's it. The scan picks it up automatically.

### Adding an ecosystem

Implement `Registry` in `tools/dep-checker/registry.go` and add one line to
`registryFor()`. Ecosystem is a property a Pin **declares** — it is never inferred from
which file a generator happens to write.

Only npm is live today. Cargo, Maven and Go are absent on purpose: no generator scaffolds
a `Cargo.toml`, `pom.xml` or `go.mod`, and the previous tool carried unreachable code for
all three.

---

## How drift is classified

> **A Pin update is low-risk if and only if the new version satisfies the Pin's existing
> constraint.**

The caret *is* the compatibility promise. `^19.2.6` permits `19.3.0` and refuses `20.0.0`,
so the Pin already contains the answer — there is no separate notion of "major" to get
wrong. Classification delegates to `internal/versioning`, which implements this correctly,
including the case that matters most here:

`^0.45.2` means `>=0.45.2 <0.46.0`. Under semver a **zero-major minor is a breaking
change**, so `drizzle-orm 0.45.2 → 0.46.0` is a Migration, not a routine bump. The old
checker compared version digits by hand, saw the major stay at `0`, and classified it
`"minor"` — it would have shipped a breaking Drizzle release inside a batch. Both
`drizzle-orm` and `drizzle-kit` are pinned at 0.x.

| Kind | Meaning |
|---|---|
| `current` | Nothing newer is published. |
| `rollup` | Satisfies the existing constraint. Batched. |
| `migration` | Breaks the existing constraint. Engineering work. |

---

## What the bot does

```
scan the Catalog  (a map read — no generator is executed)
 │
 ├─ rollup     → ONE batched PR on deps/npm/rollup
 ├─ migration  → one PR per package on deps/npm/<pkg>
 └─ deprecated → an issue. Never a PR.
```

**Rollup** — bot-owned and disposable. Rebuilt from scratch on every run, so it always
reflects what is currently behind. *Never push to it; the bot force-pushes over it.* It
excludes any package that has an open dep PR of its own.

**Migration** — human-owned from the moment it is raised. The bot opens it once and never
touches that branch again. Leaving the PR open is how you tell the bot "I have this" — and
it keeps that package out of the Rollup.

**Deprecation** — an issue, never a PR. A deprecation is not fixed by a version bump: the
replacement is a *different package*, or none. Every deprecation in this repo today needs a
rename or a removal — `@clerk/clerk-react` → `@clerk/react`, `@vercel/flags` → `flags`,
`@types/bcryptjs` → delete. Bumping to the latest *deprecated* version resolves none of them.

### When the Rollup's CI goes red

One package in the batch broke the templates. Its own semver said it was safe; it wasn't.
That makes it a Migration, and the domain already has a vehicle for that:

1. Eject it into its own `deps/npm/<pkg>` PR.
2. The next Rollup rebuild drops it automatically — via the same open-PR check that
   protects Migrations.

The bot does not bisect. Finding the culprit costs a CI log; bisecting costs ~5 full
`test-flows` runs (37 fixtures, real `pnpm install`).

---

## Generator versions

A generator's **Manifest Version** describes its *behaviour*. It must move when — and only
when — what the generator scaffolds actually changes.

That question is answered by a **Fingerprint**: a hash of the generator's *Contribution* —
the files it adds, and the edits it makes to files other generators own — taken across every
fixture that exercises it. A Contribution is a **diff**, not a file hash, so that bumping
`react` moves the Fingerprint of the two generators that name react, and not of the twenty
that merely merge something else into the same `package.json`.

```bash
dot gen-fingerprint          # generator → fingerprint
dot gen-check --docs         # rule 2
dot gen-check --bumped --base=fingerprints.json   # rules 1 and 3
```

| Rule | When | What |
|---|---|---|
| **1** | every PR touching `generators/**` | Fingerprint moved ⇒ manifest bumped, doc row synced |
| **2** | every PR | every doc version row equals its manifest Version |
| **3** | release (`make release-prep`) | every generator whose Fingerprint moved since the last tag, bumped once |

A comment fix or a `gofmt` moves no Fingerprint and demands nothing. A template edit does.

### Why dependency PRs don't bump manifests

Rule 1 deliberately does **not** fire on a diff that only touches `internal/deps/`.

A version bump computed inside a PR is a *relative* operation — `0.8.0 → 0.9.0` — worked out
against `main` at PR-creation time and applied at *merge* time, when `main` has moved. With
several dependency PRs open, the resulting version is a function of merge order. And it is
fatal for long-lived Migrations: a React Migration open for two months would collide with
every weekly Rollup that touches the same manifest, and rot under perpetual rebasing.

Deriving the bump at release, once, from the final state makes merge order **provably**
irrelevant. Five PRs in any sequence produce the identical manifest. A Migration PR touches
only its own line of the Catalog, so it rebases cleanly for months.

---

## Running it locally

```bash
go run ./tools/dep-checker scan --output=dep-report.json
go run ./tools/dep-checker report --input=dep-report.json      # markdown summary
go run ./tools/dep-checker patch --package=react --version=^19.3.0
```

`patch` writes `internal/deps/npm.go` and nothing else.

## CI

| Workflow | Runs | Does |
|---|---|---|
| `ci.yml` → `dep-scan` | every PR | scan + step summary |
| `ci.yml` → `generator-versions` | every PR | rules 1 and 2 |
| `dep-checker.yml` | Wednesdays 09:00 UTC, or manually | Rollup, Migrations, deprecation issues |

The bot does **not** run `test-flows`. The PR's own CI does, and duplicating it would double
the most expensive job in CI for no added signal.
