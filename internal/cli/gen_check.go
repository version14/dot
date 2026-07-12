package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

// generatorsDir is the root every generator lives under.
const generatorsDir = "generators"

// runGenCheck enforces the invariants that keep a Generator's Manifest Version
// honest. It is what CI runs; the flags map one-to-one onto the rules in ADR-0002.
//
//	dot gen-check --docs
//	    Rule 2: every Generator's doc version row equals its manifest Version.
//	    Diff-independent — catches desync from any cause.
//
//	dot gen-check --bumped --base=<fingerprints.json>
//	    Rule 1: every Generator whose Fingerprint moved since base has had its
//	    Manifest Version bumped. Run on PRs that touch generators/**.
//
// Rule 3 (a Catalog Pin moved → bump every Generator that names it) is the same
// check as rule 1 with the release tag as the base, so `--bumped` serves both.
// The difference is only *when* it runs: rule 1 on every PR, rule 3 at release.
// That separation is the whole reason dependency PRs are order-proof — see
// ADR-0002. Do not run --bumped on a PR that only touches internal/deps/.
func runGenCheck(args []string) int {
	fs := flag.NewFlagSet("gen-check", flag.ContinueOnError)
	docs := fs.Bool("docs", false, "rule 2: doc version rows match manifest versions")
	bumped := fs.Bool("bumped", false, "rule 1/3: generators whose Fingerprint moved have a bumped version")
	base := fs.String("base", "", "path to the base fingerprints JSON (from: dot gen-fingerprint --json)")
	dir := fs.String("dir", filepath.Join("tools", "test-flow", "testdata"), "fixture directory")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if !*docs && !*bumped {
		fmt.Fprintln(os.Stderr, "dot gen-check: pass --docs and/or --bumped")
		return 2
	}

	failures := 0
	if *docs {
		failures += checkDocsMatchManifests()
	}
	if *bumped {
		if *base == "" {
			fmt.Fprintln(os.Stderr, "dot gen-check: --bumped requires --base")
			return 2
		}
		n, err := checkFingerprintsBumped(*base, *dir)
		if err != nil {
			fmt.Fprintln(os.Stderr, "dot gen-check:", err)
			return 1
		}
		failures += n
	}

	if failures > 0 {
		fmt.Fprintf(os.Stderr, "\ndot gen-check: %d problem(s)\n", failures)
		return 1
	}
	fmt.Println("dot gen-check: ok")
	return 0
}

// checkDocsMatchManifests is rule 2.
func checkDocsMatchManifests() int {
	names, err := listGeneratorDirs(generatorsDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dot gen-check:", err)
		return 1
	}

	bad := 0
	for _, name := range names {
		manifestVersion, err := readManifestVersion(name)
		if err != nil {
			problem(name, err)
			bad++
			continue
		}

		docPath := filepath.Join("docs", "contributor", "generators", name+".md")
		content, err := os.ReadFile(docPath)
		if err != nil {
			// A generator with no doc page is a documentation gap, not a version
			// desync. Report it, but don't conflate the two.
			fmt.Fprintf(os.Stderr, "  %s: no doc page at %s\n", name, docPath)
			bad++
			continue
		}

		m := docVersionRowRe.FindStringSubmatch(string(content))
		if m == nil {
			fmt.Fprintf(os.Stderr, "  %s: no `| Version | ... |` row in %s\n", name, docPath)
			bad++
			continue
		}
		if m[1] != manifestVersion {
			fmt.Fprintf(os.Stderr,
				"  %s: doc says %s, manifest says %s — run `dot gen-bump --name %s --set %s`\n",
				name, m[1], manifestVersion, name, manifestVersion)
			bad++
		}
	}
	return bad
}

// checkFingerprintsBumped is rules 1 and 3.
//
// A Generator whose Fingerprint moved scaffolds something different, so its
// Manifest Version must move too. A Generator whose Fingerprint is unchanged
// scaffolds the same thing — a reformat or a comment fix demands nothing.
func checkFingerprintsBumped(basePath, fixtureDir string) (int, error) {
	raw, err := os.ReadFile(basePath)
	if err != nil {
		return 0, fmt.Errorf("read base fingerprints: %w", err)
	}
	var baseFP map[string]string
	if err := json.Unmarshal(raw, &baseFP); err != nil {
		return 0, fmt.Errorf("parse base fingerprints: %w", err)
	}

	headFP, err := computeFingerprints(fixtureDir)
	if err != nil {
		return 0, err
	}

	var moved []string
	for name, fp := range headFP {
		if old, existed := baseFP[name]; existed && old != fp {
			moved = append(moved, name)
		}
	}
	sort.Strings(moved)

	bad := 0
	for _, name := range moved {
		headVersion, err := readManifestVersion(name)
		if err != nil {
			problem(name, err)
			bad++
			continue
		}
		baseVersion, err := gitShowManifestVersion(name)
		if err != nil {
			problem(name, err)
			bad++
			continue
		}
		if headVersion == baseVersion {
			fmt.Fprintf(os.Stderr,
				"  %s: scaffolds something different, but its version is still %s\n"+
					"      → run `dot gen-bump --name %s --bump patch` (or minor, if this is a behaviour change)\n",
				name, headVersion, name)
			bad++
		}
	}

	if len(moved) > 0 && bad == 0 {
		fmt.Printf("dot gen-check: %d generator(s) changed and were correctly bumped: %v\n", len(moved), moved)
	}
	return bad, nil
}

// readManifestVersion reads the Version field from a generator's manifest.go.
func readManifestVersion(name string) (string, error) {
	path := manifestPath(name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read %s: %w", path, err)
	}
	m := manifestVersionRe.FindStringSubmatch(string(content))
	if m == nil {
		return "", fmt.Errorf("no Version field in %s", path)
	}
	return m[2], nil
}

// gitShowManifestVersion reads a generator's Version as it stands on the merge
// base, so "was this bumped in *this* PR" is answerable.
func gitShowManifestVersion(name string) (string, error) {
	ref := os.Getenv("DOT_BASE_REF")
	if ref == "" {
		ref = "origin/main"
	}
	out, err := exec.Command("git", "show", ref+":"+manifestPath(name)).Output()
	if err != nil {
		// The generator does not exist on the base branch — it is new in this PR,
		// so there is nothing to bump.
		return "", nil
	}
	m := manifestVersionRe.FindStringSubmatch(string(out))
	if m == nil {
		return "", nil
	}
	return m[2], nil
}

// manifestPath is where a generator's manifest lives.
func manifestPath(name string) string {
	return filepath.Join(generatorsDir, name, "manifest.go")
}

// problem reports one generator-level failure to stderr.
func problem(name string, err error) {
	fmt.Fprintf(os.Stderr, "  %s: %v\n", name, err)
}
