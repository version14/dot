package main

import "testing"

func TestClassify(t *testing.T) {
	tests := []struct {
		name   string
		pin    string
		latest string
		want   Kind
	}{
		{"nothing newer", "^19.2.6", "19.2.6", KindCurrent},
		{"registry behind the pin", "^19.2.6", "19.2.5", KindCurrent},

		{"patch within caret", "^19.2.6", "19.2.7", KindRollup},
		{"minor within caret", "^19.2.6", "19.3.0", KindRollup},
		{"major breaks caret", "^19.2.6", "20.0.0", KindMigration},

		// The bug that would have shipped a breaking Drizzle release as a routine
		// batched bump. "^0.45.2" means >=0.45.2 <0.46.0 — under semver a
		// zero-major minor IS a breaking change, and the caret says so. The old
		// checker compared digits by hand, saw the major stay at 0, and called
		// this "minor". drizzle-orm and drizzle-kit are both pinned at 0.x.
		{"zero-major minor is breaking", "^0.45.2", "0.46.0", KindMigration},
		{"zero-major patch is safe", "^0.45.2", "0.45.3", KindRollup},
		{"zero-major major is breaking", "^0.45.2", "1.0.0", KindMigration},

		// Tilde locks the minor, so a minor bump breaks it even at major >= 1.
		{"tilde allows patch", "~1.2.3", "1.2.9", KindRollup},
		{"tilde refuses minor", "~1.2.3", "1.3.0", KindMigration},

		// An exact pin promises nothing, so every update is a Migration. This is
		// why @version14/ui was moved to a caret — as an exact pin, every single
		// release of it would have demanded a human-owned PR.
		{"exact pin admits nothing", "1.2.3", "1.2.4", KindMigration},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := classify(tc.pin, tc.latest)
			if err != nil {
				t.Fatalf("classify(%q, %q): %v", tc.pin, tc.latest, err)
			}
			if got != tc.want {
				t.Errorf("classify(%q, %q) = %q, want %q", tc.pin, tc.latest, got, tc.want)
			}
		})
	}
}

func TestClassifyRejectsGarbage(t *testing.T) {
	if _, err := classify("^1.2.3", "not-a-version"); err == nil {
		t.Error("expected an error for an unparseable latest version")
	}
	if _, err := classify("~~broken", "1.2.3"); err == nil {
		t.Error("expected an error for an unparseable pin")
	}
}

func TestDeprecatedNeverEntersRollup(t *testing.T) {
	// A deprecated package can still have a newer version under the same name,
	// so it can classify as KindRollup. It must never be batched: every
	// deprecation in this repo needs a rename or a removal, and bumping to the
	// latest deprecated version resolves none of them.
	e := Entry{Kind: KindRollup, Deprecated: true}
	if e.InRollup() {
		t.Error("a deprecated package must never enter the Rollup")
	}
	if !e.Actionable() {
		t.Error("a deprecated package is actionable — it needs an issue")
	}

	ok := Entry{Kind: KindRollup, Deprecated: false}
	if !ok.InRollup() {
		t.Error("a non-deprecated rollup entry belongs in the Rollup")
	}
}
