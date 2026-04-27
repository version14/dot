package dotdir

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Manifest is the resolved-fact counterpart to ProjectSpec. Where the spec
// stores intent (constraints), the manifest stores what actually ran.
type Manifest struct {
	ToolVersion        string              `json:"tool_version"`
	LastExecutedAt     time.Time           `json:"last_executed_at"`
	ExecutionTimeMs    int64               `json:"execution_time_ms"`
	GeneratorsExecuted []ExecutedGenerator `json:"generators_executed"`
}

type ExecutedGenerator struct {
	Name              string    `json:"name"`
	VersionConstraint string    `json:"version_constraint"`
	ResolvedVersion   string    `json:"resolved_version"`
	ExecutedAt        time.Time `json:"executed_at"`
	InvocationCount   int       `json:"invocation_count"`
	ContentHash       string    `json:"content_hash"`
}

func ManifestPath(root string) string {
	return filepath.Join(root, DirName, ManifestFile)
}

func LoadManifest(root string) (*Manifest, error) {
	data, err := os.ReadFile(ManifestPath(root))
	if err != nil {
		return nil, fmt.Errorf("dotdir: read manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("dotdir: parse manifest: %w", err)
	}
	return &m, nil
}

func SaveManifest(root string, m *Manifest) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("dotdir: marshal manifest: %w", err)
	}
	path := ManifestPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("dotdir: mkdir: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("dotdir: write manifest: %w", err)
	}
	return nil
}

// IgnoreContents is the exact contents required for .dot/.gitignore by spec §9.
const IgnoreContents = `/*
!spec.json
!manifest.json
!.gitignore
`

func WriteIgnore(root string) error {
	path := filepath.Join(root, DirName, IgnoreFile)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("dotdir: mkdir: %w", err)
	}
	if err := os.WriteFile(path, []byte(IgnoreContents), 0o644); err != nil {
		return fmt.Errorf("dotdir: write gitignore: %w", err)
	}
	return nil
}
