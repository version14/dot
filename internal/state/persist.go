package state

import (
	"fmt"
	"path/filepath"

	"github.com/version14/dot/internal/fileutils"
)

// Persist writes every file in the virtual state to disk under root, using
// atomic writes. Returns the count of files written.
func Persist(s *VirtualProjectState, root string) (int, error) {
	if s == nil {
		return 0, fmt.Errorf("state: nil VirtualProjectState")
	}
	count := 0
	for _, path := range s.Paths() {
		node := s.Files[path]
		full := filepath.Join(root, filepath.FromSlash(path))
		if err := fileutils.SafeWrite(full, node.Content, 0o644); err != nil {
			return count, fmt.Errorf("state: persist %s: %w", path, err)
		}
		count++
	}
	return count, nil
}
