// dep-checker keeps the dependency Catalog (internal/deps) in step with the
// package registries.
//
// It reads the Catalog, asks each registry what the newest version is, and
// classifies the drift by the Pin's own constraint: an update that satisfies the
// constraint is a Rollup (batched, bot-owned), one that breaks it is a Migration
// (one PR, human-owned). Deprecations get an issue and never a PR.
//
// It executes no generators and writes exactly one file. See ADR-0002.
//
// Usage:
//
//	dep-checker scan   [--output=dep-report.json]
//	dep-checker report [--input=dep-report.json] [--output=-]
//	dep-checker patch  --package=<pkg> --version=<constraint>
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "scan":
		fs := flag.NewFlagSet("scan", flag.ExitOnError)
		output := fs.String("output", "dep-report.json", "path to write the JSON report")
		_ = fs.Parse(os.Args[2:])
		fail(runScan(*output))

	case "report":
		fs := flag.NewFlagSet("report", flag.ExitOnError)
		input := fs.String("input", "dep-report.json", "path to read the JSON report")
		output := fs.String("output", "-", "path to write markdown ('-' for stdout)")
		_ = fs.Parse(os.Args[2:])
		fail(runReport(*input, *output))

	case "patch":
		fs := flag.NewFlagSet("patch", flag.ExitOnError)
		pkg := fs.String("package", "", "package name as it appears in the Catalog")
		version := fs.String("version", "", "new pin, with its constraint (e.g. ^5.62.0)")
		_ = fs.Parse(os.Args[2:])
		fail(runPatch(*pkg, *version))

	default:
		fmt.Fprintf(os.Stderr, "dep-checker: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: dep-checker <scan|report|patch> [flags]")
	fmt.Fprintln(os.Stderr, "  scan    --output=dep-report.json")
	fmt.Fprintln(os.Stderr, "  report  --input=dep-report.json --output=-")
	fmt.Fprintln(os.Stderr, "  patch   --package=<pkg> --version=<constraint>")
}

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "dep-checker:", err)
		os.Exit(1)
	}
}
