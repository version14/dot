package cli

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/fingerprint"
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/pkg/dotapi"
)

// fixture is the subset of a test-flow testdata case that fingerprinting needs:
// which flow to run, and the answers to run it with.
type fixture struct {
	Name     string                 `json:"name"`
	FlowID   string                 `json:"flow_id"`
	Answers  map[string]flow.Answer `json:"answers"`
	Disabled bool                   `json:"disabled"`
}

// runGenFingerprint prints one Fingerprint per generator: a hash of what that
// generator contributes to a scaffolded project, across every fixture scenario.
//
// CI compares this output between the base branch and the head of a PR:
//
//   - A generator whose Fingerprint moved scaffolds something different, so its
//     Manifest Version must be bumped and its doc row synced (rule 1).
//   - A generator whose Fingerprint is unchanged scaffolds the same thing, so a
//     reformat or a comment fix demands nothing.
//
// It is entirely in-memory apart from the scratch dir Scaffold writes to — no
// pnpm install, no network. The whole fixture suite runs in seconds.
func runGenFingerprint(args []string) int {
	fs := flag.NewFlagSet("gen-fingerprint", flag.ContinueOnError)
	dir := fs.String("dir", filepath.Join("tools", "test-flow", "testdata"), "directory of flow fixtures")
	asJSON := fs.Bool("json", false, "emit JSON instead of a text table")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	prints, err := computeFingerprints(*dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dot gen-fingerprint:", err)
		return 1
	}

	if *asJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(prints); err != nil {
			fmt.Fprintln(os.Stderr, "dot gen-fingerprint:", err)
			return 1
		}
		return 0
	}

	names := make([]string, 0, len(prints))
	for name := range prints {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%-40s %s\n", name, prints[name])
	}
	return 0
}

// computeFingerprints replays every fixture and returns generator → Fingerprint.
func computeFingerprints(dir string) (map[string]string, error) {
	fixtures, err := loadFixtures(dir)
	if err != nil {
		return nil, err
	}
	if len(fixtures) == 0 {
		return nil, fmt.Errorf("no enabled fixtures found in %s", dir)
	}

	registry, err := DefaultGeneratorRegistry()
	if err != nil {
		return nil, fmt.Errorf("build registry: %w", err)
	}

	scratch, err := os.MkdirTemp("", "dot-fingerprint-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(scratch)

	collector := fingerprint.NewCollector()
	for _, fx := range fixtures {
		if err := replay(fx, registry, scratch, collector); err != nil {
			return nil, err
		}
	}
	return collector.Fingerprints(), nil
}

// loadFixtures reads every enabled fixture in dir.
func loadFixtures(dir string) ([]fixture, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read fixtures: %w", err)
	}

	var out []fixture
	for _, e := range entries {
		if e.IsDir() || filepath.Ext(e.Name()) != ".json" {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", e.Name(), err)
		}
		var fx fixture
		if err := json.Unmarshal(raw, &fx); err != nil {
			return nil, fmt.Errorf("parse %s: %w", e.Name(), err)
		}
		if fx.Disabled || fx.FlowID == "" {
			continue
		}
		out = append(out, fx)
	}
	return out, nil
}

// replay scaffolds one fixture with the collector attached, so it observes what
// each generator contributed.
func replay(fx fixture, registry *generator.Registry, scratch string, collector *fingerprint.Collector) error {
	def, ok := flows.Default().Get(fx.FlowID)
	if !ok {
		return fmt.Errorf("%s: unknown flow %q", fx.Name, fx.FlowID)
	}

	// Scaffold writes the tree to disk; we only care about the in-memory state
	// the Observer saw, so the scratch dir is thrown away.
	out := filepath.Join(scratch, fx.Name)
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}

	if _, err := Scaffold(context.Background(), ScaffoldOptions{
		Flow:        def,
		Registry:    registry,
		OutputDir:   out,
		ToolVersion: "gen-fingerprint",
		Logger:      dotapi.DiscardLogger{},
		Runner:      flow.NewScriptedRunner(fx.Answers, nil, nil),
		Observer:    collector,
	}); err != nil {
		return fmt.Errorf("%s: scaffold: %w", fx.Name, err)
	}
	return nil
}
