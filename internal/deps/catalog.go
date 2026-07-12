// Package deps is the dependency Catalog: the single place where the versions
// of third-party packages scaffolded into generated projects are pinned.
//
// Generators name the packages they need; the Catalog says which version they
// get. Nothing else in the repo hardcodes a version.
//
//	// generators/react_app/generator.go
//	d.Merge(map[string]interface{}{
//	    "dependencies":    deps.NPM("react", "react-dom"),
//	    "devDependencies": deps.NPM("vite", "@types/react"),
//	})
//
// There is exactly one Pin per package, repo-wide. This is deliberate: when
// versions lived inline in 40 generator files, the same package drifted to two
// different versions without anyone noticing (see ADR-0002). One key, one
// version, and divergence is unrepresentable.
//
// This package is imported by generators, so it stays free of I/O — no network,
// no filesystem. Registry lookups live in tools/dep-checker, which imports this.
package deps

import (
	"fmt"
	"sort"
)

// Ecosystem is the package universe a Pin belongs to. It determines where the
// package's versions are published, how its version strings are written, and
// which kinds of update are mechanically safe.
//
// Ecosystem is declared by the Pin, never inferred from which file a generator
// happens to write. The previous checker guessed it from the filename, which is
// why adding an ecosystem meant teaching a scanner a new filename.
type Ecosystem string

const (
	EcosystemNPM Ecosystem = "npm"
)

// Pin is one package at the version generators scaffold, in that Ecosystem's
// native notation ("^5.61.3" for npm, "v1.9.0" for Go, "1.0.4" for Cargo).
//
// The constraint is not decoration. It is the compatibility promise the package
// makes, and dep-checker respects it: an update is low-risk if and only if the
// new version satisfies it. "^0.45.2" does not permit 0.46.0 — under semver a
// zero-major minor is a breaking change — so such an update is a Migration, not
// a Rollup, with no special-casing needed.
type Pin struct {
	Ecosystem Ecosystem
	Name      string
	Version   string
}

// NPM returns the {name: version} map for the given npm packages, ready to merge
// into a package.json "dependencies" or "devDependencies" block.
//
// Whether a package is a prod or dev dependency is the generator's business —
// it decides which block to merge into. The Catalog only says which version.
//
// NPM panics on an unknown package. A typo is a programming error, not a runtime
// condition: it surfaces the moment the generator runs, and test-flows catches
// it. Returning a zero version would silently scaffold a broken package.json.
func NPM(names ...string) map[string]interface{} {
	out := make(map[string]interface{}, len(names))
	for _, name := range names {
		version, ok := npm[name]
		if !ok {
			panic(fmt.Sprintf(
				"deps: %q is not in the npm Catalog — add it to internal/deps/npm.go", name))
		}
		out[name] = version
	}
	return out
}

// NPMVersion returns the pinned version of a single npm package.
// Use it where a raw version string is needed rather than a merge-ready map.
func NPMVersion(name string) string {
	version, ok := npm[name]
	if !ok {
		panic(fmt.Sprintf(
			"deps: %q is not in the npm Catalog — add it to internal/deps/npm.go", name))
	}
	return version
}

// Has reports whether name is pinned in the npm Catalog.
func Has(name string) bool {
	_, ok := npm[name]
	return ok
}

// All returns every Pin in the Catalog, sorted by ecosystem then name.
//
// This is how dep-checker scans. It is a plain map read: no generator is
// executed, no Answers are supplied, and no conditional branch can hide a
// package from it. The previous scanner ran Generate() with empty Answers and
// therefore only ever saw one arbitrary branch — which is why @clerk/nextjs was
// never checked once in the tool's lifetime.
func All() []Pin {
	pins := make([]Pin, 0, len(npm))
	for name, version := range npm {
		pins = append(pins, Pin{Ecosystem: EcosystemNPM, Name: name, Version: version})
	}
	sort.Slice(pins, func(i, j int) bool {
		if pins[i].Ecosystem != pins[j].Ecosystem {
			return pins[i].Ecosystem < pins[j].Ecosystem
		}
		return pins[i].Name < pins[j].Name
	})
	return pins
}
