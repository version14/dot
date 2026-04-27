package dotdir

import (
	"path/filepath"

	"github.com/version14/dot/internal/spec"
)

const (
	DirName      = ".dot"
	SpecFile     = "spec.json"
	ManifestFile = "manifest.json"
	IgnoreFile   = ".gitignore"
)

// SpecPath returns the absolute path to the spec file under root.
func SpecPath(root string) string {
	return filepath.Join(root, DirName, SpecFile)
}

// LoadSpec reads .dot/spec.json from a project root.
func LoadSpec(root string) (*spec.ProjectSpec, error) {
	return spec.Load(SpecPath(root))
}

// SaveSpec writes .dot/spec.json under a project root.
func SaveSpec(root string, s *spec.ProjectSpec) error {
	return spec.Save(SpecPath(root), s)
}
