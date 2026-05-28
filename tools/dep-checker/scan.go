package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/version14/dot/internal/cli"
	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

// depFiles are the filenames whose contents the scanner extracts deps from.
// Keyed by filename → Ecosystem.
var depFiles = map[string]Ecosystem{
	"package.json": EcosystemNPM,
	"go.mod":       EcosystemGo,
	"Cargo.toml":   EcosystemCargo,
	"pom.xml":      EcosystemMaven,
}

// ScanResult holds everything found during a single generator scan pass,
// before registry queries happen.
type ScanResult struct {
	Generator string
	File      string
	Ecosystem Ecosystem
	Package   string
	Current   string
}

// runScan walks all registered generators, invokes each with a fresh
// VirtualProjectState, extracts declared dependencies, queries their
// respective registries, and writes a JSON Report to outputPath.
func runScan(outputPath string) error {
	registry, err := cli.DefaultGeneratorRegistry()
	if err != nil {
		return fmt.Errorf("build registry: %w", err)
	}

	var raw []ScanResult
	for _, entry := range registry.All() {
		results, err := scanGenerator(entry.Manifest.Name, entry.Generator)
		if err != nil {
			fmt.Printf("warn: skipping %s: %v\n", entry.Manifest.Name, err)
			continue
		}
		raw = append(raw, results...)
	}

	entries, err := checkAll(raw)
	if err != nil {
		return fmt.Errorf("registry check: %w", err)
	}

	return writeReport(outputPath, entries)
}

// scanGenerator calls Generate() on g with an empty VirtualProjectState and
// returns every dep-file entry it wrote.
func scanGenerator(name string, g dotapi.Generator) ([]ScanResult, error) {
	st := state.NewVirtualProjectState(spec.ProjectMetadata{ProjectName: "dep-checker-probe"})

	ctx := &dotapi.Context{
		Spec:    &spec.ProjectSpec{},
		Answers: map[string]interface{}{},
		State:   st,
		Logger:  dotapi.DiscardLogger{},
	}

	if err := g.Generate(ctx); err != nil {
		return nil, err
	}

	var results []ScanResult
	for filename, eco := range depFiles {
		node, ok := st.GetFile(filename)
		if !ok {
			continue
		}
		deps, err := extractDeps(filename, eco, node.Content)
		if err != nil {
			fmt.Printf("warn: %s: extract %s: %v\n", name, filename, err)
			continue
		}
		for pkg, ver := range deps {
			results = append(results, ScanResult{
				Generator: name,
				File:      filename,
				Ecosystem: eco,
				Package:   pkg,
				Current:   ver,
			})
		}
	}
	return results, nil
}

// extractDeps parses the content of a dependency file and returns
// {packageName: versionConstraint} for that ecosystem.
func extractDeps(filename string, eco Ecosystem, content []byte) (map[string]string, error) {
	switch eco {
	case EcosystemNPM:
		return extractNPMDeps(content)
	case EcosystemGo:
		return extractGoDeps(content)
	case EcosystemCargo:
		return extractCargoDeps(content)
	default:
		return nil, fmt.Errorf("no extractor for %s", eco)
	}
}

func extractNPMDeps(content []byte) (map[string]string, error) {
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, err
	}
	out := make(map[string]string)
	for k, v := range pkg.Dependencies {
		out[k] = v
	}
	for k, v := range pkg.DevDependencies {
		out[k] = v
	}
	return out, nil
}

func extractGoDeps(content []byte) (map[string]string, error) {
	mod := state.NewGoMod()
	if err := mod.Load(content); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(mod.Requires))
	for _, r := range mod.Requires {
		// Skip indirect deps — they are managed by `go mod tidy`, not generators.
		if !strings.HasSuffix(r.Version, "// indirect") {
			out[r.Path] = r.Version
		}
	}
	return out, nil
}

func extractCargoDeps(content []byte) (map[string]string, error) {
	// Minimal TOML [dependencies] parser — only handles simple
	// name = "version" lines inside the [dependencies] section.
	out := make(map[string]string)
	inDeps := false
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "[dependencies]" || line == "[dev-dependencies]" {
			inDeps = true
			continue
		}
		if strings.HasPrefix(line, "[") {
			inDeps = false
			continue
		}
		if !inDeps || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		name := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		// Only handle simple string values: name = "version"
		if strings.HasPrefix(val, `"`) {
			ver := strings.Trim(val, `"`)
			out[name] = ver
		}
	}
	return out, nil
}

// checkAll queries registries for every ScanResult and returns DepEntries.
func checkAll(raw []ScanResult) ([]DepEntry, error) {
	// Cache one checker per ecosystem to reuse HTTP clients.
	checkers := map[string]RegistryChecker{}

	var entries []DepEntry
	for _, r := range raw {
		key := string(r.Ecosystem)
		checker, ok := checkers[key]
		if !ok {
			checker, ok = checkerFor(r.File)
			if !ok {
				continue
			}
			checkers[key] = checker
		}

		info, err := checker.Check(r.Package, r.Current)
		if err != nil {
			fmt.Printf("warn: %s/%s: registry: %v\n", r.Generator, r.Package, err)
			continue
		}

		entries = append(entries, DepEntry{
			Generator:         r.Generator,
			File:              r.File,
			Ecosystem:         r.Ecosystem,
			Package:           r.Package,
			Current:           r.Current,
			Latest:            info.Latest,
			Outdated:          isOutdated(r.Current, info.Latest),
			Deprecated:        info.Deprecated,
			DeprecationNotice: info.Notice,
		})
	}
	return entries, nil
}
