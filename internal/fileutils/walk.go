package fileutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Walk recursively lists every file under root, returning paths relative to root
// in forward-slash form. Directories are not included.
func Walk(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, relErr := filepath.Rel(root, path)
		if relErr != nil {
			return relErr
		}
		paths = append(paths, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("fileutils: walk %s: %w", root, err)
	}
	return paths, nil
}

// Exists returns true when path refers to any filesystem entry.
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
