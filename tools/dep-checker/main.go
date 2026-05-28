// dep-checker scans generator source files for hardcoded dependency versions,
// queries their respective registries, and reports outdated or deprecated
// packages. It can also patch a single generator file in-place.
//
// Usage:
//
//	go run ./tools/dep-checker scan [--output=dep-report.json]
//	go run ./tools/dep-checker patch --generator=<name> --package=<pkg> --current=<ver> --latest=<ver>
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: dep-checker <scan|patch> [flags]")
		fmt.Fprintln(os.Stderr, "  scan  --output=dep-report.json")
		fmt.Fprintln(os.Stderr, "  patch --generator=<name> --package=<pkg> --current=<ver> --latest=<ver>")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		fs := flag.NewFlagSet("scan", flag.ExitOnError)
		output := fs.String("output", "dep-report.json", "path to write the JSON report")
		fs.Parse(os.Args[2:])
		if err := runScan(*output); err != nil {
			fmt.Fprintln(os.Stderr, "dep-checker scan:", err)
			os.Exit(1)
		}

	case "patch":
		fs := flag.NewFlagSet("patch", flag.ExitOnError)
		generator := fs.String("generator", "", "generator name (e.g. express_server_typescript_deps)")
		pkg := fs.String("package", "", "package name to bump")
		current := fs.String("current", "", "current version string as it appears in generator.go (e.g. ^4.21.0)")
		latest := fs.String("latest", "", "latest version from registry (e.g. 5.1.0)")
		fs.Parse(os.Args[2:])
		if err := runPatch(*generator, *pkg, *current, *latest); err != nil {
			fmt.Fprintln(os.Stderr, "dep-checker patch:", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "dep-checker: unknown command %q\n", os.Args[1])
		os.Exit(2)
	}
}

// writeReport serialises entries into a Report and writes it to path.
func writeReport(path string, entries []DepEntry) error {
	report := Report{
		GeneratedAt: time.Now().UTC(),
		Entries:     entries,
	}
	if entries == nil {
		report.Entries = []DepEntry{}
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	outdated := 0
	deprecated := 0
	for _, e := range entries {
		if e.Outdated {
			outdated++
		}
		if e.Deprecated {
			deprecated++
		}
	}
	fmt.Printf("dep-checker: scanned %d deps — %d outdated, %d deprecated → %s\n",
		len(entries), outdated, deprecated, path)
	return nil
}
