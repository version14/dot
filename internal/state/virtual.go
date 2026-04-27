package state

import (
	"fmt"
	"sort"
	"time"

	"github.com/version14/dot/internal/fileutils"
	"github.com/version14/dot/internal/spec"
)

// FileNode represents a single file in the virtual project tree.
type FileNode struct {
	Path            string
	Content         []byte
	ContentType     ContentType
	CreatedBy       string
	Transformations []string
	ModifiedAt      time.Time
}

// VirtualProjectState is an in-memory filesystem. Generators write here;
// disk write happens once after post-execution validation succeeds.
type VirtualProjectState struct {
	Files    map[string]*FileNode
	Metadata spec.ProjectMetadata

	currentGenerator string
}

func NewVirtualProjectState(meta spec.ProjectMetadata) *VirtualProjectState {
	return &VirtualProjectState{
		Files:    map[string]*FileNode{},
		Metadata: meta,
	}
}

// SetCurrentGenerator tags subsequent file operations with the given name
// so transformation history attributes correctly.
func (s *VirtualProjectState) SetCurrentGenerator(name string) {
	s.currentGenerator = name
}

// CreateFile adds a new raw file. Returns an error if the path already exists.
func (s *VirtualProjectState) CreateFile(path string, content []byte) error {
	path = fileutils.Normalize(path)
	if _, exists := s.Files[path]; exists {
		return fmt.Errorf("state: file %q already exists", path)
	}
	s.Files[path] = &FileNode{
		Path:        path,
		Content:     append([]byte(nil), content...),
		ContentType: ContentRaw,
		CreatedBy:   s.currentGenerator,
		ModifiedAt:  time.Now(),
	}
	return nil
}

// WriteFile overwrites or creates a file with the given content type.
func (s *VirtualProjectState) WriteFile(path string, content []byte, ct ContentType) {
	path = fileutils.Normalize(path)
	existing, ok := s.Files[path]
	if !ok {
		s.Files[path] = &FileNode{
			Path:        path,
			Content:     append([]byte(nil), content...),
			ContentType: ct,
			CreatedBy:   s.currentGenerator,
			ModifiedAt:  time.Now(),
		}
		return
	}
	existing.Content = append([]byte(nil), content...)
	existing.ContentType = ct
	existing.ModifiedAt = time.Now()
	existing.Transformations = append(existing.Transformations, s.currentGenerator)
}

func (s *VirtualProjectState) GetFile(path string) (*FileNode, bool) {
	f, ok := s.Files[fileutils.Normalize(path)]
	return f, ok
}

func (s *VirtualProjectState) FileExists(path string) bool {
	_, ok := s.Files[fileutils.Normalize(path)]
	return ok
}

func (s *VirtualProjectState) DeleteFile(path string) {
	delete(s.Files, fileutils.Normalize(path))
}

// Paths returns every file path in deterministic (sorted) order.
func (s *VirtualProjectState) Paths() []string {
	out := make([]string, 0, len(s.Files))
	for k := range s.Files {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// UpdateJSON loads or initializes a JSON document at path, applies fn, and
// stores the result. The file is created if missing.
func (s *VirtualProjectState) UpdateJSON(path string, fn func(*JSONDoc) error) error {
	path = fileutils.Normalize(path)

	doc := NewJSONDoc()
	if existing, ok := s.Files[path]; ok && len(existing.Content) > 0 {
		if err := doc.Load(existing.Content); err != nil {
			return fmt.Errorf("state: load JSON %s: %w", path, err)
		}
	}
	if err := fn(doc); err != nil {
		return fmt.Errorf("state: update JSON %s: %w", path, err)
	}
	out, err := doc.Marshal()
	if err != nil {
		return fmt.Errorf("state: marshal JSON %s: %w", path, err)
	}
	s.WriteFile(path, out, ContentJSON)
	return nil
}

// UpdateYAML mirrors UpdateJSON for YAML documents.
func (s *VirtualProjectState) UpdateYAML(path string, fn func(*YAMLDoc) error) error {
	path = fileutils.Normalize(path)

	doc := NewYAMLDoc()
	if existing, ok := s.Files[path]; ok && len(existing.Content) > 0 {
		if err := doc.Load(existing.Content); err != nil {
			return fmt.Errorf("state: load YAML %s: %w", path, err)
		}
	}
	if err := fn(doc); err != nil {
		return fmt.Errorf("state: update YAML %s: %w", path, err)
	}
	out, err := doc.Marshal()
	if err != nil {
		return fmt.Errorf("state: marshal YAML %s: %w", path, err)
	}
	s.WriteFile(path, out, ContentYAML)
	return nil
}

// UpdateGoMod loads or initializes a go.mod, applies fn, and re-serializes.
func (s *VirtualProjectState) UpdateGoMod(fn func(*GoMod) error) error {
	const path = "go.mod"
	mod := NewGoMod()
	if existing, ok := s.Files[path]; ok && len(existing.Content) > 0 {
		if err := mod.Load(existing.Content); err != nil {
			return fmt.Errorf("state: load go.mod: %w", err)
		}
	}
	if err := fn(mod); err != nil {
		return fmt.Errorf("state: update go.mod: %w", err)
	}
	s.WriteFile(path, mod.Marshal(), ContentGoMod)
	return nil
}
