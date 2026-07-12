package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/version14/dot/internal/deps"
)

// scanConcurrency bounds in-flight registry requests. The Catalog is ~75 pins;
// this keeps a full scan under a couple of seconds without hammering npm.
const scanConcurrency = 8

// runScan reads the Catalog, asks each Pin's registry what the newest version
// is, classifies the drift, and writes a JSON Report.
//
// The scan is a map read. It does not execute a single generator.
//
// The previous scanner ran Generate() on every generator with empty Answers and
// read back whatever package.json fell out, on the theory (ADR-0001) that this
// was "strictly accurate". It was not: generators branch on Answers, so an empty
// Answers map executes one arbitrary branch. Three generators
// (auth_clerk_frontend, sentry_frontend, storybook_setup) pick their package
// *name* from the framework answer, so @clerk/nextjs, @sentry/nextjs and
// @storybook/nextjs were never checked once in that tool's entire lifetime.
//
// A Catalog has no branches. Every Pin is visible because every Pin is a map key.
func runScan(outputPath string) error {
	pins := deps.All()

	registries := map[deps.Ecosystem]Registry{}
	for _, pin := range pins {
		if _, ok := registries[pin.Ecosystem]; ok {
			continue
		}
		reg, ok := registryFor(pin.Ecosystem)
		if !ok {
			return fmt.Errorf("no registry for ecosystem %q", pin.Ecosystem)
		}
		registries[pin.Ecosystem] = reg
	}

	var (
		mu      sync.Mutex
		entries []Entry
		failed  []string
		wg      sync.WaitGroup
		sem     = make(chan struct{}, scanConcurrency)
	)

	for _, pin := range pins {
		wg.Add(1)
		go func(pin deps.Pin) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			entry, err := checkPin(registries[pin.Ecosystem], pin)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				failed = append(failed, fmt.Sprintf("%s: %v", pin.Name, err))
				return
			}
			entries = append(entries, entry)
		}(pin)
	}
	wg.Wait()

	// A registry that cannot be reached is a scan failure, not a silent skip.
	// The old checker printed a warning and dropped the package — so a flaky npm
	// response quietly turned into "this dependency is up to date".
	if len(failed) > 0 {
		sort.Strings(failed)
		return fmt.Errorf("registry lookups failed for %d package(s):\n  %s",
			len(failed), strings.Join(failed, "\n  "))
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Ecosystem != entries[j].Ecosystem {
			return entries[i].Ecosystem < entries[j].Ecosystem
		}
		return entries[i].Package < entries[j].Package
	})

	return writeReport(outputPath, entries)
}

func checkPin(reg Registry, pin deps.Pin) (Entry, error) {
	info, err := reg.Latest(pin.Name)
	if err != nil {
		return Entry{}, err
	}

	kind, err := classify(pin.Version, info.Latest)
	if err != nil {
		return Entry{}, err
	}

	return Entry{
		Ecosystem:  pin.Ecosystem,
		Package:    pin.Name,
		Pin:        pin.Version,
		Latest:     info.Latest,
		Proposed:   reg.Render(pin.Version, info.Latest),
		Kind:       kind,
		Deprecated: info.Deprecated,
		Notice:     info.Notice,
	}, nil
}

func writeReport(path string, entries []Entry) error {
	report := Report{GeneratedAt: time.Now().UTC(), Entries: entries}
	if entries == nil {
		report.Entries = []Entry{}
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return err
	}

	fmt.Printf("dep-checker: %d pins — %d rollup, %d migration, %d deprecated → %s\n",
		len(report.Entries),
		len(report.Rollup()),
		len(report.Migrations()),
		len(report.Deprecations()),
		path)
	return nil
}
