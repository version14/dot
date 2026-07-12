package main

import (
	"time"

	"github.com/version14/dot/internal/deps"
)

// Entry is one Catalog Pin with its registry status.
//
// There is no Generator field, and that is the point. A Pin belongs to the
// Catalog, not to a generator — the old report carried one row per
// (generator, package) pair, which is how the same package ended up pinned to two
// versions and how a single PR title came to list twenty packages.
type Entry struct {
	Ecosystem deps.Ecosystem `json:"ecosystem"`
	Package   string         `json:"package"`

	// Pin is the constraint as written in the Catalog, e.g. "^19.2.6".
	Pin string `json:"pin"`

	// Latest is the newest published version, e.g. "20.0.1".
	Latest string `json:"latest"`

	// Proposed is Latest rendered in the Pin's notation, e.g. "^20.0.1".
	// This is what would be written back to the Catalog.
	Proposed string `json:"proposed"`

	// Kind is how this drift must be closed: current, rollup, or migration.
	Kind Kind `json:"kind"`

	Deprecated bool   `json:"deprecated"`
	Notice     string `json:"deprecation_notice,omitempty"`
}

// Actionable reports whether this entry needs a PR or an issue.
func (e Entry) Actionable() bool {
	return e.Kind != KindCurrent || e.Deprecated
}

// InRollup reports whether this entry belongs in the batched Rollup PR.
//
// A deprecated package never does, even when a newer version exists under the
// same name. Every deprecation in this repo needs a rename or a removal
// (@clerk/clerk-react → @clerk/react, @vercel/flags → flags, @types/bcryptjs →
// delete), and bumping to the latest *deprecated* version resolves none of them.
// Deprecation is a human decision: issue only, never a PR.
func (e Entry) InRollup() bool {
	return e.Kind == KindRollup && !e.Deprecated
}

// Report is the JSON output of `dep-checker scan`.
type Report struct {
	GeneratedAt time.Time `json:"generated_at"`
	Entries     []Entry   `json:"entries"`
}

// Rollup returns the entries that belong in the batched Rollup PR.
func (r Report) Rollup() []Entry {
	return r.filter(Entry.InRollup)
}

// Migrations returns the entries needing one human-owned PR each.
func (r Report) Migrations() []Entry {
	return r.filter(func(e Entry) bool {
		return e.Kind == KindMigration && !e.Deprecated
	})
}

// Deprecations returns the entries needing a tracking issue.
func (r Report) Deprecations() []Entry {
	return r.filter(func(e Entry) bool { return e.Deprecated })
}

func (r Report) filter(keep func(Entry) bool) []Entry {
	var out []Entry
	for _, e := range r.Entries {
		if keep(e) {
			out = append(out, e)
		}
	}
	return out
}
