package main

import (
	"fmt"

	"github.com/version14/dot/internal/versioning"
)

// Kind is how a Pin's drift must be closed. See CONTEXT.md.
type Kind string

const (
	// KindCurrent — nothing newer is published. No action.
	KindCurrent Kind = "current"

	// KindRollup — the new version satisfies the Pin's existing constraint, so
	// the package itself promises compatibility. Batched into the single
	// bot-owned Rollup PR.
	KindRollup Kind = "rollup"

	// KindMigration — the new version breaks the Pin's constraint. This is
	// engineering work, not a version bump: it may need generator code changes.
	// One human-owned PR per package.
	KindMigration Kind = "migration"
)

// classify decides how a Pin's drift must be closed, from the Pin's own
// constraint.
//
// The rule is one line: an update is low-risk if and only if the new version
// satisfies the existing constraint. The caret *is* the compatibility promise —
// "^19.2.6" permits 19.3.0 and refuses 20.0.0 — so the Pin already tells us the
// answer and we need no separate notion of "major".
//
// This is why the old checker mis-shipped Drizzle. It compared version numbers
// by hand and called 0.45.2 → 0.46.0 a "minor" because the major digit didn't
// move. But "^0.45.2" means >=0.45.2 <0.46.0: under semver a zero-major minor is
// a *breaking* change, and npm knows it. Delegating to the constraint gets 0.x
// right for free, with no special case — which matters, because drizzle-orm and
// drizzle-kit are both pinned at 0.x.
func classify(pin, latest string) (Kind, error) {
	constraint, err := versioning.ParseConstraint(pin)
	if err != nil {
		return "", fmt.Errorf("parse pin %q: %w", pin, err)
	}

	latestVer, err := versioning.Parse(latest)
	if err != nil {
		return "", fmt.Errorf("parse latest %q: %w", latest, err)
	}

	// Nothing newer? Then there is no drift, regardless of the constraint.
	// Allows() alone cannot answer this: it rejects a version *below* the anchor
	// for the same reason it rejects one above the permitted range.
	if anchor, ok := constraint.Anchor(); ok && latestVer.Compare(anchor) <= 0 {
		return KindCurrent, nil
	}

	if constraint.Allows(latestVer) {
		return KindRollup, nil
	}
	return KindMigration, nil
}
