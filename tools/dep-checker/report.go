package main

import "time"

// Ecosystem identifies which package registry a dependency lives in.
type Ecosystem string

const (
	EcosystemNPM   Ecosystem = "npm"
	EcosystemGo    Ecosystem = "go"
	EcosystemCargo Ecosystem = "cargo"
	EcosystemMaven Ecosystem = "maven"
)

// DepEntry is one dependency found in a generator, with its registry status.
type DepEntry struct {
	Generator         string    `json:"generator"`
	File              string    `json:"file"`
	Ecosystem         Ecosystem `json:"ecosystem"`
	Package           string    `json:"package"`
	Current           string    `json:"current"`
	Latest            string    `json:"latest"`
	UpdateType        string    `json:"update_type,omitempty"`
	Outdated          bool      `json:"outdated"`
	Deprecated        bool      `json:"deprecated"`
	DeprecationNotice string    `json:"deprecation_notice,omitempty"`
}

// Report is the top-level JSON output of the scan subcommand.
type Report struct {
	GeneratedAt time.Time  `json:"generated_at"`
	Entries     []DepEntry `json:"entries"`
}
