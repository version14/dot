// Package spec defines the core types that flow through the dot pipeline.
// Every input layer (CLI TUI, dot.yaml, future MCP) produces a Spec.
// The generator engine consumes it. Neither side knows about the other.
package spec

// ProjectSpec holds the top-level identity of a project.
// Both Language and Type are open strings — no closed enums — so community
// plugins can introduce new values without touching core.
type ProjectSpec struct {
	Name     string `json:"name"`
	Language string `json:"language"` // e.g. "go", "typescript", "python"
	Type     string `json:"type"`     // e.g. "frontend", "api", "mobile"
}

// Spec is the authoritative description of a project. It is produced by
// input layers and consumed by generators.
// Extensions is the universal namespace: all generators — official or
// community — store their answers here, namespaced by plugin
// (e.g. "react.architecture", "prisma.provider").
type Spec struct {
	Project    ProjectSpec    `json:"project"`
	Extensions map[string]any `json:"extensions,omitempty"`
}
