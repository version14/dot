package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Load reads a ProjectSpec from disk for re-run / dot update flows.
func Load(path string) (*ProjectSpec, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("spec: read %s: %w", path, err)
	}
	var s ProjectSpec
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("spec: parse %s: %w", path, err)
	}
	return &s, nil
}

// Save writes a ProjectSpec to disk as pretty-printed JSON, creating the
// parent directory if needed.
func Save(path string, s *ProjectSpec) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("spec: marshal: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("spec: mkdir: %w", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("spec: write %s: %w", path, err)
	}
	return nil
}
