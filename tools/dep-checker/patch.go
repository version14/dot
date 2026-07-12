package main

import (
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"regexp"

	"github.com/version14/dot/internal/deps"
)

// catalogPath is the one file dep-checker is allowed to write.
var catalogPath = filepath.Join("internal", "deps", "npm.go")

// runPatch rewrites a Pin in the Catalog.
//
// This is the whole patcher. It touches one machine-owned file with one
// canonical shape, and the value it writes is the same value the scan read.
//
// The old patcher had to reverse-map a version it had observed in generated JSON
// back to a Go string literal somewhere in one of 40 hand-written generator
// files, which it located by regex. Nothing guaranteed such a literal existed:
// auth_clerk_frontend held its version in a variable (`version := "^5.61.3"`), so
// the regex never matched, the patch errored, the bash driver `break`ed out of
// its loop and ran `git checkout main` *without resetting the working tree* —
// leaving every already-patched generator dirty for the next iteration's
// `git add generators/` to sweep into an unrelated commit. That is how commit
// 997e489, titled "bump @clerk/clerk-react", came to contain 13 files of other
// packages and no clerk changes at all.
//
// None of that machinery survives. There is no reverse-mapping, so there is
// nothing to fail, so there is no dirty tree to leak.
func runPatch(pkg, proposed string) error {
	if pkg == "" || proposed == "" {
		return fmt.Errorf("--package and --version are required")
	}
	if !deps.Has(pkg) {
		return fmt.Errorf("%q is not in the Catalog", pkg)
	}

	content, err := os.ReadFile(catalogPath)
	if err != nil {
		return fmt.Errorf("read %s: %w", catalogPath, err)
	}

	// The Catalog is gofmt'd and machine-owned, so every entry is exactly
	// `"name": "version",`. Anchoring on the quoted key makes this exact — unlike
	// the old patcher, which was guessing at the shape of hand-written code.
	pattern := regexp.MustCompile(`(?m)^(\s*` + regexp.QuoteMeta(fmt.Sprintf("%q", pkg)) + `:\s*)"[^"]*"(,)$`)

	// "The key isn't there" and "the key is already at this version" are
	// different facts and must not be conflated. Comparing content before and
	// after would report a no-op patch as a missing package — and a no-op is
	// entirely normal: re-running the bot after a Rollup has merged asks it to
	// patch pins that are already current. Erroring there would abort the run.
	if !pattern.Match(content) {
		return fmt.Errorf("no Pin for %q in %s (Catalog and source disagree)", pkg, catalogPath)
	}

	updated := pattern.ReplaceAll(content, []byte(`${1}"`+proposed+`"${2}`))
	if string(updated) == string(content) {
		fmt.Printf("%s: %s is already at %s (no-op)\n", catalogPath, pkg, proposed)
		return nil
	}

	// Re-align the map literal: a longer version string shifts gofmt's column.
	formatted, err := format.Source(updated)
	if err != nil {
		return fmt.Errorf("gofmt %s after patching %q: %w", catalogPath, pkg, err)
	}

	if err := os.WriteFile(catalogPath, formatted, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", catalogPath, err)
	}

	fmt.Printf("patched %s: %s → %s\n", catalogPath, pkg, proposed)
	return nil
}
