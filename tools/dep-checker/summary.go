package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// runReport renders a scan Report as markdown, grouped the way the work is
// actually done: what the bot will batch, what needs a human, and what needs a
// decision.
//
// The old report was one flat table of every (generator, package) pair with
// "outdated: yes/no" columns, which told you nothing about what would happen next.
func runReport(inputPath, outputPath string) error {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("read report: %w", err)
	}
	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("parse report: %w", err)
	}

	var b strings.Builder
	b.WriteString("# Dependency scan\n\n")

	rollup := report.Rollup()
	migrations := report.Migrations()
	deprecations := report.Deprecations()

	fmt.Fprintf(&b, "%d pins in the Catalog — **%d rollup**, **%d migration**, **%d deprecated**.\n\n",
		len(report.Entries), len(rollup), len(migrations), len(deprecations))

	if len(rollup) == 0 && len(migrations) == 0 && len(deprecations) == 0 {
		b.WriteString("Every Pin is current. Nothing to do.\n")
		return emit(outputPath, b.String())
	}

	if len(rollup) > 0 {
		b.WriteString("## Rollup\n\n")
		b.WriteString("Satisfy their existing constraint, so the package itself promises compatibility. ")
		b.WriteString("These go into one batched PR.\n\n")
		writeTable(&b, rollup)
	}

	if len(migrations) > 0 {
		b.WriteString("## Migrations\n\n")
		b.WriteString("**Break their existing constraint.** Each needs its own PR and may need generator ")
		b.WriteString("code changes — this is engineering work, not a version bump.\n\n")
		writeTable(&b, migrations)
	}

	if len(deprecations) > 0 {
		b.WriteString("## Deprecated\n\n")
		b.WriteString("The registry says stop using these. A deprecation is **never** fixed by a version ")
		b.WriteString("bump — the replacement is a different package, or none. Issue only.\n\n")
		b.WriteString("| Package | Pin | Notice |\n| --- | --- | --- |\n")
		for _, e := range deprecations {
			fmt.Fprintf(&b, "| `%s` | `%s` | %s |\n", e.Package, e.Pin, oneLine(e.Notice))
		}
		b.WriteString("\n")
	}

	return emit(outputPath, b.String())
}

func writeTable(b *strings.Builder, entries []Entry) {
	b.WriteString("| Package | Pin | Latest | Proposed |\n| --- | --- | --- | --- |\n")
	for _, e := range entries {
		fmt.Fprintf(b, "| `%s` | `%s` | `%s` | `%s` |\n", e.Package, e.Pin, e.Latest, e.Proposed)
	}
	b.WriteString("\n")
}

func oneLine(s string) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "|", "\\|")
	if len(s) > 160 {
		s = s[:157] + "..."
	}
	return s
}

func emit(path, content string) error {
	if path == "-" {
		fmt.Print(content)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
