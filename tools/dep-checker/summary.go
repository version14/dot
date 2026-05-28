package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
)

type reportSummary struct {
	Total      int
	Outdated   int
	Deprecated int
}

type depRow struct {
	Generator  string
	Ecosystem  string
	Package    string
	Current    string
	Latest     string
	Update     string
	Outdated   bool
	Deprecated bool
}

func runReport(inputPath, outputPath string) error {
	report, err := loadReport(inputPath)
	if err != nil {
		return err
	}

	rows := make([]depRow, 0, len(report.Entries))
	stats := reportSummary{}
	for _, entry := range report.Entries {
		stats.Total++
		if entry.Outdated {
			stats.Outdated++
		}
		if entry.Deprecated {
			stats.Deprecated++
		}
		rows = append(rows, depRow{
			Generator:  entry.Generator,
			Ecosystem:  string(entry.Ecosystem),
			Package:    entry.Package,
			Current:    entry.Current,
			Latest:     entry.Latest,
			Update:     entry.UpdateType,
			Outdated:   entry.Outdated,
			Deprecated: entry.Deprecated,
		})
	}

	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Ecosystem != rows[j].Ecosystem {
			return rows[i].Ecosystem < rows[j].Ecosystem
		}
		if rows[i].Package != rows[j].Package {
			return rows[i].Package < rows[j].Package
		}
		return rows[i].Generator < rows[j].Generator
	})

	lines := make([]string, 0, len(rows)+6)
	lines = append(lines, "# Dependency scan report")
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("- Total dependencies: %d", stats.Total))
	lines = append(lines, fmt.Sprintf("- Outdated: %d", stats.Outdated))
	lines = append(lines, fmt.Sprintf("- Deprecated: %d", stats.Deprecated))
	lines = append(lines, "")
	lines = append(lines, "| Generator | Ecosystem | Package | Current | Latest | Update | Outdated | Deprecated |")
	lines = append(lines, "| --- | --- | --- | --- | --- | --- | --- | --- |")

	for _, row := range rows {
		lines = append(lines, fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %s | %s |",
			row.Generator,
			row.Ecosystem,
			row.Package,
			row.Current,
			row.Latest,
			row.Update,
			boolLabel(row.Outdated),
			boolLabel(row.Deprecated),
		))
	}

	output := strings.Join(lines, "\n") + "\n"
	if outputPath == "-" {
		fmt.Print(output)
		return nil
	}
	return os.WriteFile(outputPath, []byte(output), 0644)
}

func loadReport(path string) (Report, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Report{}, fmt.Errorf("read report: %w", err)
	}

	var report Report
	if err := json.Unmarshal(data, &report); err != nil {
		return Report{}, fmt.Errorf("parse report: %w", err)
	}
	return report, nil
}

func boolLabel(value bool) string {
	if value {
		return "yes"
	}
	return "no"
}
