# Context

The domain language of `dot`. This file is a glossary, not a spec — it says what
words mean, never how anything is built.

## Scaffolding

**Generator**
A unit of project scaffolding. Given a set of Answers, it contributes files and
configuration to a Project. Generators compose: a single Project is assembled from
many of them.

**Flow**
An ordered set of Questions presented to a user, producing the Answers that select
and configure Generators.

**Answers**
The user's responses to a Flow. Generators may behave differently depending on
Answers — the same Generator can contribute different Dependencies on different runs.

**Project**
The scaffolded output: the tree of files a user ends up with.

**Manifest Version**
A Generator's declared semver. It describes the *behaviour* of the Generator — the
shape of what it scaffolds and the Answers it responds to. A Project records the
Manifest Version it was scaffolded with, so drift can be detected later.

## Dependencies

**Dependency**
A third-party package that a Generator causes a Project to depend on — `react`,
`vitest`, `bcryptjs`. Distinct from a Generator's own Go imports, which are ordinary
source dependencies and are not part of this domain.

**Ecosystem**
The package universe a Dependency belongs to: npm, Go, Cargo, Maven. An Ecosystem
determines where a Dependency's versions are published, how its version strings are
written, and what kinds of update are mechanically safe. Ecosystem is a property the
Dependency *declares*, not something inferred from the files a Generator happens to
write.

**Pin**
The exact version of a Dependency that Generators scaffold, expressed in that
Ecosystem's native notation. A Pin is a deliberate choice, not a resolved value.

**Catalog**
The single place every Pin lives. There is exactly one Pin per Dependency: a
Dependency cannot be pinned to one version in one Generator and another version
elsewhere. Generators name the Dependencies they need; the Catalog says which
version they get.

**Drift**
Two things wearing one word — keep them apart:

- *Pin drift* — the Catalog's Pin has fallen behind what the Ecosystem publishes.
  This is expected, continuous, and is what the dependency check exists to find.
- *Generator drift* — a Project was scaffolded with a Manifest Version that no
  longer satisfies the constraint it recorded. This is what `dot doctor` reports.

**Deprecation**
An Ecosystem's signal that a Dependency should no longer be used. Deprecation is a
judgement call for a human — unlike Pin drift, it has no mechanical fix, because the
replacement (if any) is a different Dependency, not a different version.

## Closing Pin drift

Not all Pin drift is the same kind of work, and the two kinds are not
interchangeable.

**Rollup**
The single aggregate change carrying every *low-risk* Pin update — those where the
Ecosystem's own versioning promises compatibility. A Rollup is disposable and
machine-owned: it is rebuilt from scratch on every run to reflect whatever is
currently behind, and no human ever edits it. There is at most one Rollup in flight.

**Migration**
A *single* high-risk Pin update — one where the Ecosystem signals a break — which may
require changing Generator code, not just a version string. A Migration is a piece of
engineering work, not a version bump. It becomes human-owned the moment it is raised,
and may stay unresolved indefinitely; an unresolved Migration is the standing record
that the decision has been deferred.

The distinction is about *who owns the change*. The machine proposes Rollups and
rewrites them at will. It proposes a Migration once, then never touches it again.

## Knowing when a Generator has changed

**Contribution**
What a single Generator adds to a Project when it runs — the files it creates and the
edits it makes to files other Generators own. A Project is the sum of its Generators'
Contributions.

**Fingerprint**
An identity for a Generator's Contribution, taken across every scenario the Generator
is exercised in. Two Generators with the same Fingerprint scaffold the same thing.

A Fingerprint answers the only question that matters when deciding whether a Manifest
Version must move: *did what this Generator scaffolds actually change?* Reformatting
its source does not change the Fingerprint. Editing a template does. Moving a Pin the
Generator names does too — which is why a Pin change is a change to every Generator
that names it, even though no Generator's source was touched.

The Fingerprint is what makes Manifest Version meaningful: a Generator's version moves
when, and only when, its Contribution does.
